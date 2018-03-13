package fsquota

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/hashicorp/go-version"
)

var expectedKernelVersion = version.Must(version.NewVersion("4.6.0"))

func isKernel46OrLater() (is46OrLater bool, err error) {
	var utsname syscall.Utsname
	if err = syscall.Uname(&utsname); err != nil {
		return
	}

	// get kernel release string
	kernelReleaseBytes := make([]byte, 0, cap(utsname.Release))
	for _, b := range utsname.Release {
		if b == '-' {
			// Replace the first dash with a null-byte
			b = 0x0
		}

		if b == 0x0 {
			break
		}

		kernelReleaseBytes = append(kernelReleaseBytes, byte(b))
	}

	var kernelRelease *version.Version
	if kernelRelease, err = version.NewVersion(string(kernelReleaseBytes)); err != nil {
		return
	}

	is46OrLater = kernelRelease.GreaterThan(expectedKernelVersion) || kernelRelease.Equal(expectedKernelVersion)
	return
}

const passwdFile = "/etc/passwd"
const groupFile = "/etc/group"

func getIDsFromUserOrGroupFile(path string) (ids []uint32, err error) {
	var f *os.File

	if f, err = os.Open(path); err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, ":") {
			// Skip malformed line
			continue
		}
		lineParts := strings.SplitN(line, ":", 4)
		if len(lineParts) != 4 {
			continue
		}

		var id uint64
		var parseErr error
		if id, parseErr = strconv.ParseUint(lineParts[2], 10, 32); parseErr == nil {
			ids = append(ids, uint32(id))
		}
	}

	return
}
