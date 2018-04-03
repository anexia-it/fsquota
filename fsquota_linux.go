package fsquota

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/docker/docker/pkg/mount"
	"github.com/speijnik/go-errortree"
	"golang.org/x/sys/unix"
)

func getQuota(t quotaCtlType, path string, idString string) (info *Info, err error) {
	var device string
	var id uint32
	if device, id, err = prepareArguments(path, idString); err != nil {
		return
	}

	return internalGetQuota(t, device, id)
}

func internalGetQuota(t quotaCtlType, device string, id uint32) (info *Info, err error) {
	// Retrieve the quota info again
	quotaInfoStruct := dqblkFromLimits(&Limits{})
	if err = quotactl(cmdGetQuota, t, device, id, unsafe.Pointer(quotaInfoStruct)); err != nil {
		return
	}

	info = quotaInfoStruct.toInfo()
	return
}

func setQuota(t quotaCtlType, path string, idString string, limits *Limits) (info *Info, err error) {
	var device string
	var id uint32

	if device, id, err = prepareArguments(path, idString); err != nil {
		return
	}

	quotaInfoStruct := dqblkFromLimits(limits)

	if err = quotactl(cmdSetQuota, t, device, id, unsafe.Pointer(quotaInfoStruct)); err != nil {
		return
	}

	info, err = internalGetQuota(userQuota, device, id)
	return
}

func setUserQuota(path string, usr *user.User, limits *Limits) (info *Info, err error) {
	return setQuota(userQuota, path, usr.Uid, limits)
}

func getUserInfo(path string, user *user.User) (info *Info, err error) {
	info, err = getQuota(userQuota, path, user.Uid)
	return
}

func setGroupQuota(path string, group *user.Group, limits *Limits) (info *Info, err error) {
	return setQuota(groupQuota, path, group.Gid, limits)
}

func getGroupInfo(path string, group *user.Group) (info *Info, err error) {
	return getQuota(groupQuota, path, group.Gid)
}

func pathToDevice(path string) (device string, err error) {
	if path, err = filepath.EvalSymlinks(path); err != nil {
		// Evaluate symlinks first
		return
	}

	// Call stat on the path, as it may be a device
	var statRes os.FileInfo
	if statRes, err = os.Stat(path); err != nil {
		// os.Stat failed
		return
	}

	fileMode := statRes.Mode()

	if (fileMode & os.ModeCharDevice) != 0 {
		// Char device found
		err = errors.New("target must not be a character device")
		return
	} else if (fileMode & os.ModeDevice) != 0 {
		// Block device found: as expected
		device = path
		return
	}

	// Getting this far means path was not a device, but a regular path
	// We thus need to retrieve the device underlying the path
	var statT *syscall.Stat_t
	var statTOK bool
	if statT, statTOK = statRes.Sys().(*syscall.Stat_t); !statTOK {
		err = errors.New("internal error: could not retrieve Stat_t from stat result")
		return
	}

	// Retrieve major and minor device numbers
	devMajor := unix.Major(statT.Dev)
	devMinor := unix.Minor(statT.Dev)

	// Retrieve mount info
	var mountInfos []*mount.Info
	if mountInfos, err = mount.GetMounts(); err != nil {
		return
	}

	for _, mountInfo := range mountInfos {
		if uint32(mountInfo.Major) == devMajor && uint32(mountInfo.Minor) == devMinor && strings.HasPrefix(mountInfo.Source, "/dev/") {
			// Call pathToDevice again
			device, err = pathToDevice(mountInfo.Source)
			return
		}
	}
	err = errors.New("unable to find mount point for path")

	return
}

func prepareArguments(path string, idString string) (device string, id uint32, err error) {
	// Look up the device beneath the provided path
	var pathErr error
	if device, pathErr = pathToDevice(path); pathErr != nil {
		err = errortree.Add(err, "path", pathErr)
	}

	// Convert ID string to uint32
	if id64, parseErr := strconv.ParseUint(idString, 10, 32); parseErr != nil {
		err = errortree.Add(err, "id", parseErr)
	} else {
		id = uint32(id64)
	}

	return
}

type reportLegacyIDLookupFn func() ([]uint32, error)

func getReport(path string, typ quotaCtlType, idLookupFn reportLegacyIDLookupFn) (report *Report, err error) {
	var device string

	if device, err = pathToDevice(path); err != nil {
		return
	}

	// Look up kernel version to determine which approach to take
	var kernel46OrLater bool
	if kernel46OrLater, err = isKernel46OrLater(); err != nil {
		return
	}

	if !kernel46OrLater {
		// Kernel version < 4.6: use legacy approach via passwd file
		report, err = getReportLegacy(typ, device, idLookupFn)
	} else {
		// Kernel version >= 4.6: use GETNEXTQUOTA approach
		report, err = getReportByNextQuota(typ, device)
	}

	return
}

func getUserReport(path string) (report *Report, err error) {
	return getReport(path, userQuota, func() ([]uint32, error) {
		return getIDsFromUserOrGroupFile(passwdFile)
	})
}

func getGroupReport(path string) (report *Report, err error) {
	return getReport(path, groupQuota, func() ([]uint32, error) {
		return getIDsFromUserOrGroupFile(groupFile)
	})
}

type nextdqblk struct {
	dqbBHardlimit uint64
	dqbBSoftlimit uint64
	dqbCurSpace   uint64
	dqbIHardlimit uint64
	dqbISoftlimit uint64
	dqbCurInodes  uint64
	dqbBTime      uint64
	dqbITime      uint64
	dqbValid      uint32
	dqbId         uint32
}

func (n nextdqblk) toDqblk() *dqblk {
	return &dqblk{
		dqbBHardlimit: n.dqbBHardlimit,
		dqbBSoftlimit: n.dqbBSoftlimit,
		dqbCurSpace:   n.dqbCurSpace,
		dqbIHardlimit: n.dqbIHardlimit,
		dqbISoftlimit: n.dqbISoftlimit,
		dqbCurInodes:  n.dqbCurInodes,
		dqbBTime:      n.dqbBTime,
		dqbITime:      n.dqbITime,
		dqbValid:      n.dqbValid,
	}
}

func getReportByNextQuota(t quotaCtlType, device string) (report *Report, err error) {
	rep := &Report{
		Infos: make(map[string]*Info),
	}

	// Always start at ID 0
	nextId := uint32(0)

	for {
		nextQuotaInfoStruct := &nextdqblk{}

		// Retrieve per-user quota
		if err = quotactl(cmdGetNextQuota, t, device, nextId, unsafe.Pointer(nextQuotaInfoStruct)); err != nil {
			if scErr, isSCErr := err.(*os.SyscallError); isSCErr && os.IsNotExist(scErr.Err) {
				// GetNextQuota will respond ESRCH when no further quotas can be found
				err = nil
			}

			// Break out of our loop
			break
		}

		rep.Infos[fmt.Sprint(nextQuotaInfoStruct.dqbId)] = nextQuotaInfoStruct.toDqblk().toInfo()
		nextId += 1
	}

	if err == nil {
		report = rep
	}

	return
}

func getReportLegacy(t quotaCtlType, device string, idLookupFn reportLegacyIDLookupFn) (report *Report, err error) {
	var ids []uint32
	if ids, err = idLookupFn(); err != nil {
		return
	}

	rep := &Report{
		Infos: make(map[string]*Info, len(ids)),
	}

	for _, id := range ids {
		var info *Info
		if info, err = internalGetQuota(t, device, id); err != nil {
			return
		}

		if info.isEmpty() {
			// Skip empty info objects
			continue
		}

		rep.Infos[fmt.Sprint(id)] = info
	}

	report = rep

	return
}

func quotasSupported(t quotaCtlType, path string) (supported bool, err error) {
	var device string
	if device, _, err = prepareArguments(path, "0"); err != nil {
		return
	}

	if _, err = internalGetQuota(t, device, 0); err == nil {
		supported = true
	}

	return
}

func userQuotasSupported(path string) (supported bool, err error) {
	return quotasSupported(userQuota, path)
}

func groupQuotasSupported(path string) (supported bool, err error) {
	return quotasSupported(groupQuota, path)
}
