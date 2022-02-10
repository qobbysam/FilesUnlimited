package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/qobbysam/filesunlimited/pkgs/rpcserver"
	"github.com/smallnest/rpcx/client"
)

var addr = "localhost:10015"

func main() {
	//flag.Parse()
	path, _ := os.Executable()
	dirpath := filepath.Dir(path)
	d, err := client.NewPeer2PeerDiscovery("tcp@"+addr, "")

	if err != nil {

		fmt.Println(err)
	}
	xclient := client.NewXClient("RPCADDFUNC", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	file_ := filepath.Join(dirpath, "data", "sfiles", "out.pdf")

	file, err := os.Open(file_)

	if err != nil {
		fmt.Println("failed to open file")
	}
	bytes_, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println("failed to read file")
	}

	info, _ := file.Stat()
	sendargs := rpcserver.UploadOneFileArg{
		Type: "pdf",
		Size: info.Size(),
		File: bytes_,
	}

	reply := rpcserver.UploadOneFileResponse{}

	err = xclient.Call(context.Background(), "SaveFile", &sendargs, &reply)

	if err != nil {
		fmt.Println("send failed ", err)
	}

	log.Println("Send Successfull")

	log.Println(reply)
	// err := xclient.SendFile(context.Background(), "abc.txt", 0, map[string]string{"SAVEFILE": "bar"})
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("send ok")

}
