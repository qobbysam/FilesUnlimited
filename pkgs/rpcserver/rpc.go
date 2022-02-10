package rpcserver

import (
	"context"
	"fmt"
	"log"

	"github.com/qobbysam/fileserver/pkgs/config"
	"github.com/qobbysam/fileserver/pkgs/myminio"
	"github.com/smallnest/rpcx/server"
)

//Minio Instance for Rpc
type LocalMin struct {
	Min *myminio.MinioStruct
}

var GlobalMin *LocalMin

type RPCADDFUNC int

func (t *RPCADDFUNC) SaveFile(ctx context.Context, args *UploadOneFileArg, reply *UploadOneFileResponse) error {
	//reply.C = args.A * args.B
	//return nil

	to_save := myminio.SaveFileArg{
		Type:  args.Type,
		Size:  args.Size,
		Bytes: args.File,
	}
	name, err := GlobalMin.Min.SaveFile(&to_save)

	if err != nil {
		reply.Good = false
		reply.Path = ""
		reply.Err = err
		return err
	}

	reply.Good = true
	reply.Path = name
	reply.Err = nil

	return err
}

func (t *RPCADDFUNC) DeleteFile(ctx context.Context, args *DeleteOneFileArg, reply *DeleteFileResponse) error {
	//reply.C = args.A * args.B
	//return nil

	to_delete := myminio.DeleteFileArg{
		Type: args.Type,
		Name: args.Name,
	}
	name, err := GlobalMin.Min.DeleteFile(to_delete)

	if err != nil {
		reply.Good = name.Good
		reply.Name = args.Name
		reply.Err = err
		return err
	}

	reply.Good = name.Good
	reply.Name = ""
	reply.Err = nil

	return err
}

type RpcObject struct {
	Addr             string
	FileTransferAddr string
	ServerString     string
	Rpc              *config.RPC
	Server           *server.Server
}

func (ro *RpcObject) RpcServe(donechan chan struct{}, errchan chan error) {

	log.Println("serving rpc on port :  ", ro.Addr)

	err := ro.Server.Serve(ro.ServerString, ro.Addr)
	if err != nil {
		panic(err)
	}
}

func (ro *RpcObject) Init() error {

	err := ro.Server.RegisterName(ro.Rpc.AddFileFunc, new(RPCADDFUNC), "")

	return err

}

func NewRpcObject(cfg *config.BigConfig, min *myminio.MinioStruct) (*RpcObject, error) {

	GlobalMin = &LocalMin{Min: min}

	out := RpcObject{}

	fmt.Println("")
	s := server.NewServer()

	out.Server = s
	out.Addr = cfg.RpcConfig.Address
	out.FileTransferAddr = cfg.RpcConfig.AddFileFunc
	out.Rpc = cfg.RpcConfig
	out.ServerString = cfg.RpcConfig.ServerString
	return &out, nil

}
