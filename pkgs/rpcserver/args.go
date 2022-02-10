package rpcserver

type UploadOneFileArg struct {
	Type string
	Size int64
	File []byte
}

type DeleteOneFileArg struct {
	Type string
	Name string
	//	File []byte
}
