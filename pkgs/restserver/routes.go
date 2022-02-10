package restserver

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/minio/minio-go/v7"
	"github.com/qobbysam/filesunlimited/pkgs/myminio"
)

type RestMin struct {
	Min *myminio.MinioStruct
}

var GlobalMin *RestMin

//Build Routes builds the rest routes with the config values
func (rs *RestServer) BuildRoutes() error {

	//router := chi.NewRouter()
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		//AllowedOrigins:   []string{"https://*", "http://*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Get("/img", GetImgFuncHandler)

	router.Get("/pdf", GetPdfFuncHandler)

	//router.Get("/txt", Get)

	out_router := chi.NewRouter()

	out_router.Mount(rs.MountString, router)

	rs.Mux = out_router

	return nil
}

//const GetHandler  = http.HandleFunc("/", GetHandlerFunc)

var GetImgFuncHandler = func(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	opts := minio.GetObjectOptions{}
	bucketname := GlobalMin.Min.Exec.Buckets.IMG

	file, err := GlobalMin.Min.RetrieveFile(path, bucketname, &opts)

	fileBytes := file.Bytes
	if err != nil {
		panic(err)
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.Write(fileBytes)

}

var GetPdfFuncHandler = func(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	opts := minio.GetObjectOptions{}
	bucketname := GlobalMin.Min.Exec.Buckets.PDF

	file, err := GlobalMin.Min.RetrieveFile(path, bucketname, &opts)

	fileBytes := file.Bytes
	if err != nil {
		panic(err)
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.Write(fileBytes)
}
