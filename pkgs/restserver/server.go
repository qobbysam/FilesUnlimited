package restserver

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/qobbysam/filesunlimited/pkgs/config"
	"github.com/qobbysam/filesunlimited/pkgs/myminio"
)

type RestServer struct {
	Mux         chi.Router
	Port        string
	MountString string
}

func (rs *RestServer) Init() error {

	err := rs.BuildRoutes()

	return err
}

func (rs *RestServer) StartServer(donechan chan struct{}, errchan chan error) {

	err := rs.Init()

	if err != nil {
		log.Println("failed to Init rest server")
	}

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(rs.Mux, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
	log.Println("Serving server onn port :  ", rs.Port)
	err = http.ListenAndServe(rs.Port, rs.Mux)

	if err != nil {
		log.Println("failed to serve server rest server")
	}

	log.Println("Served server onn port :  ", rs.Port)
}

func NewRestServer(cfg *config.BigConfig, min *myminio.MinioStruct) (*RestServer, error) {
	GlobalMin = &RestMin{Min: min}
	out := RestServer{}
	out.MountString = cfg.RestServer.Mountstring
	out.Port = cfg.RestServer.Port
	return &out, nil
}
