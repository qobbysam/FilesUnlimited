package internal

import (
	"log"
	"sync"

	"github.com/qobbysam/filesunlimited/pkgs/app"
)

func StartApp(startstring string, pathtoconfig string) {

	app, err := app.NewApp(startstring, pathtoconfig)
	// cfg, err := config.NewConfig(path_config)

	if err != nil {
		msg := "failed to create app"
		log.Println(msg, err)

		panic(msg)
	}

	errchan := make(chan error, 2)

	donechan := make(chan struct{}, 3)

	var wg sync.WaitGroup

	app.StartApp(wg, donechan, errchan)

	//wg.Wait()
}
