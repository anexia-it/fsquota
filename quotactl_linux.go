package fsquota

import (
	"os"
	"syscall"
	"unsafe"
)

type quotaCtlCmd uintptr

const (
	// Q_GETQUOTA
	cmdGetQuota quotaCtlCmd = 0x00800007
	// Q_SETQUOTA
	cmdSetQuota = 0x00800008
	// Q_GETNEXTQUOTA
	cmdGetNextQuota = 0x00800009
)

type quotaCtlType uintptr

const (
	// USRQUOTA
	userQuota quotaCtlType = 0
	// GRPQUOTA
	groupQuota = 1
)

const (
	// SUBCMDSHIFT
	quotaSubCmdShift = 8
	// SUBCMDMASK
	quotaSubCmdMask = 0x00ff
)

func getQuotaCommand(cmd quotaCtlCmd, typ quotaCtlType) uintptr {
	// Implementation of Q_CMD macro logic
	return uintptr(cmd<<quotaSubCmdShift) | uintptr(typ&quotaSubCmdMask)
}

const (
	// QIF_BLIMITS
	qifBLimits uint32 = 1
	// QIF_SPACE
	qifSpace = 2
	// QIF_ILIMITS
	qifILimits = 4
	// QIF_INODES
	qifInodes = 8
	// QIF_BTIME
	qifBTime = 16
	// QIF_ITIME
	qifITime = 32

	// QIF_ALL
	qifAll = qifBLimits | qifSpace | qifILimits | qifInodes | qifBTime | qifITime
)

func bytesToDqBlocks(bytes uint64) uint64 {
	return bytes / 1024
}

func dqBlocksToBytes(blocks uint64) uint64 {
	return blocks * 1024
}

type dqblk struct {
	dqbBHardlimit uint64
	dqbBSoftlimit uint64
	dqbCurSpace   uint64
	dqbIHardlimit uint64
	dqbISoftlimit uint64
	dqbCurInodes  uint64
	dqbBTime      uint64
	dqbITime      uint64
	dqbValid      uint32
}

func (d dqblk) toInfo() (info *Info) {
	info = &Info{
		Limits: Limits{
			Files: Limit{
				hard: &d.dqbIHardlimit,
				soft: &d.dqbISoftlimit,
			},
			Bytes: Limit{
				hard: &d.dqbBHardlimit,
				soft: &d.dqbBSoftlimit,
			},
		},
		FilesUsed: d.dqbCurInodes,
		BytesUsed: d.dqbCurSpace,
	}

	//info.BytesUsed = info.BytesUsed
	*info.Limits.Bytes.hard = dqBlocksToBytes(*info.Limits.Bytes.hard)
	*info.Limits.Bytes.soft = dqBlocksToBytes(*info.Limits.Bytes.soft)

	return
}

func dqblkFromLimits(limits *Limits) (quotas *dqblk) {
	quotas = &dqblk{}

	// Process bytes limits
	if bytesHard, bytesSoft, haveBytesLimits := limits.Bytes.getValues(); haveBytesLimits {
		// Set flag indicating the block limit fields are valid
		quotas.dqbValid |= qifBLimits

		// Convert bytes to blocks and set the fields
		quotas.dqbBHardlimit = bytesToDqBlocks(bytesHard)
		quotas.dqbBSoftlimit = bytesToDqBlocks(bytesSoft)
	}

	if inodesHard, inodesSoft, haveInodesLimits := limits.Files.getValues(); haveInodesLimits {
		// Set flag indicating the inode limit fields are valid
		quotas.dqbValid |= qifILimits

		// Set the inode limit fields
		quotas.dqbIHardlimit = inodesHard
		quotas.dqbISoftlimit = inodesSoft
	}

	// Ensure only known flags have been set
	quotas.dqbValid = quotas.dqbValid & qifAll

	return
}

func quotactl(cmd quotaCtlCmd, typ quotaCtlType, device string, id uint32, target unsafe.Pointer) (err error) {
	// Thin wrapper around SYS_QUOTACTL syscall
	fullCommand := getQuotaCommand(cmd, typ)

	var deviceNamePtr *byte
	if deviceNamePtr, err = syscall.BytePtrFromString(device); err != nil {
		return
	}

	if _, _, errno := syscall.RawSyscall6(syscall.SYS_QUOTACTL, fullCommand,
		uintptr(unsafe.Pointer(deviceNamePtr)), uintptr(id), uintptr(target), 0, 0); errno != 0 {
		err = os.NewSyscallError("quotactl", errno)
	}

	return
}
