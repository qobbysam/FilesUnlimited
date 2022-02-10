package app

import (
	"errors"
	"log"
	"sync"

	"github.com/qobbysam/filesunlimited/pkgs/config"
	"github.com/qobbysam/filesunlimited/pkgs/executor"
	"github.com/qobbysam/filesunlimited/pkgs/myminio"
	"github.com/qobbysam/filesunlimited/pkgs/restserver"
	"github.com/qobbysam/filesunlimited/pkgs/rpcserver"
)

type App struct {
	RestServer  *restserver.RestServer
	RpcObject   *rpcserver.RpcObject
	MinioStruct *myminio.MinioStruct
	StartString string
}

//NewApp builds a new application.
//We have the option for a rest server , rpc server or both
func NewApp(startstring, filepath string) (*App, error) {

	out := App{}

	out.StartString = startstring

	cfg, err := config.NewConfig(filepath)

	if err != nil {
		log.Println("failed to get config")
	}
	exec_ := executor.NewExecutor(cfg)

	miniostruct, err := myminio.NewMinioStruct(cfg, exec_)
	if err != nil {
		log.Println("failed to minio ")
	}

	switch startstring {
	case "all":
		restserver_, err := restserver.NewRestServer(cfg, miniostruct)

		if err != nil {
			log.Println("failed to rest server ")
		}

		rpc_object, err := rpcserver.NewRpcObject(cfg, miniostruct)
		if err != nil {
			log.Println("failed to rest server ")
		}
		out.MinioStruct = miniostruct
		out.RestServer = restserver_
		out.RpcObject = rpc_object

		return &out, nil

	case "rest":
		restserver_, err := restserver.NewRestServer(cfg, miniostruct)

		if err != nil {
			log.Println("failed to rest server ")
		}
		out.MinioStruct = miniostruct
		out.RestServer = restserver_

		return &out, nil

	case "rpc":
		rpc_object, err := rpcserver.NewRpcObject(cfg, miniostruct)
		if err != nil {
			log.Println("failed to rest server ")
		}
		out.MinioStruct = miniostruct
		//out.RestServer = restserver_
		out.RpcObject = rpc_object

		return &out, nil

	default:
		return nil, errors.New("failed to create an app")

	}

}

func (ap *App) StartApp(wg sync.WaitGroup, donechan chan struct{}, errchan chan error) {
	//var wg sync.WaitGroup

	err := ap.MinioStruct.Init()

	if err != nil {
		log.Println("failed to init minio")
	}
	log.Println("minio init successful")

	switch ap.StartString {

	case "all":

		err = ap.RpcObject.Init()

		if err != nil {
			log.Println("failed to init rpc")
		}
		log.Println("rpc init successful")

		wg.Add(1)
		go func() {
			ap.RestServer.StartServer(donechan, errchan)
		}()

		wg.Add(1)
		go func() {
			ap.RpcObject.RpcServe(donechan, errchan)
		}()

	case "rest":

		wg.Add(1)
		go func() {
			ap.RestServer.StartServer(donechan, errchan)
		}()

	case "rpc":

		err = ap.RpcObject.Init()

		if err != nil {
			log.Println("failed to init rpc")
		}
		log.Println("rpc init successful")

		wg.Add(1)
		go func() {
			ap.RpcObject.RpcServe(donechan, errchan)
		}()

	}
	wg.Wait()
}
