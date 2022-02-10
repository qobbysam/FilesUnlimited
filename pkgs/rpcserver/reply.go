package rpcserver

type UploadOneFileResponse struct {
	Err  error
	Good bool
	Path string
}

type DeleteFileResponse struct {
	Err  error
	Good bool
	Name string
}
