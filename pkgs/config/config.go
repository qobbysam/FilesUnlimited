package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type BigConfig struct {
	BucketConfig *Buckets         `json:"buckets"`
	MinioConfig  *MinioConnection `json:"minio"`
	RpcConfig    *RPC             `json:"rpc"`
	RestServer   *RestConfig      `json:"restconfig"`
}

type RestConfig struct {
	Mountstring string `json:"mountstring"`
	Port        string `json:"port"`
}
type RPC struct {
	Address      string `json:"address"`
	AddFileFunc  string `json:"addfunc"`
	ServerString string `json:"serverstring"`
}
type Buckets struct {
	Txt string `json:"txt"`
	PDF string `json:"pdf"`
	CSV string `json:"csv"`
	IMG string `json:"img"`
}

type MinioConnection struct {
	AccessUrl   string `json:"accessurl"`
	Accesspoint string `json:"accesspoint"`
	PrivateKey  string `json:"privatekey"`
	UseSSL      bool   `json:"usessl"`
	Location    string `json:"location"`
}

func NewConfig(filepath string) (*BigConfig, error) {

	out := BigConfig{}

	file, err := os.Open(filepath)

	if err != nil {
		log.Println("failed to load file")
	}

	bytes_, _ := io.ReadAll(file)

	//var Cfg BigConfig

	err = json.Unmarshal(bytes_, &out)

	if err != nil {
		log.Println("failed to marshall config file", err)
	}
	return &out, nil
}
