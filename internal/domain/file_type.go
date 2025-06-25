package domain

type FileType string

const (
	FileTypeFile   FileType = "file"
	FileTypeFolder FileType = "folder"
)

func (ft FileType) IsValid() bool {
	return ft == FileTypeFile || ft == FileTypeFolder
}

func (ft FileType) IsFile() bool {
	return ft == FileTypeFile
}

func (ft FileType) IsFolder() bool {
	return ft == FileTypeFolder
}

func (ft FileType) String() string {
	return string(ft)
}
