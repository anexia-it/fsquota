package fsquota

// Info contains quota information
type Info struct {
	Limits

	// Byte usage
	BytesUsed uint64
	// File usage
	FilesUsed uint64
}

func (i *Info) isEmpty() bool {
	bytesHard, bytesSoft, _ := i.Bytes.getValues()
	filesHard, filesSoft, _ := i.Files.getValues()

	return bytesSoft == 0 && bytesHard == 0 && i.BytesUsed == 0 &&
		filesHard == 0 && filesSoft == 0 && i.FilesUsed == 0
}
