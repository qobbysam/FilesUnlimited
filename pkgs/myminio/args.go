package myminio

type SaveFileArg struct {
	//Name   string
	Type  string
	Size  int64
	Bytes []byte
}

type GetFileResponse struct {
	Name  string
	Size  int64
	Bytes []byte
}

type DeleteFileResponse struct {
	Good bool
	Name string
}

type DeleteFileArg struct {
	Type string
	Name string
}
