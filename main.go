package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/qobbysam/fileserver/internal"
)

func main() {

	startString := flag.String("st", "all", "Start String  -st=rest for restserver instance, -st=rpc for just rpc instance, ")

	path, _ := os.Executable()

	basedir := filepath.Dir(path)

	default_config_file := filepath.Join(basedir, "config.json")
	pathConfig := flag.String("c", default_config_file, "-c=/usr/config.json  pass path to config File")

	flag.Parse()

	internal.StartApp(*startString, *pathConfig)
}
