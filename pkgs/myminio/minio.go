package myminio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/qobbysam/filesunlimited/pkgs/config"
	"github.com/qobbysam/filesunlimited/pkgs/executor"
)

type MinioStruct struct {
	// Endpoint        string
	// AccessKey       string
	// SecretAccessKey string
	// UseSSL          bool

	Client   *minio.Client
	Exec     *executor.Executor
	Location string
	Ctx      context.Context
	//Buckets []string
}

func (ms *MinioStruct) DoSave(bucketname, objectname string, file []byte, size int64, opts minio.PutObjectOptions) error {

	reader := bytes.NewReader(file)
	uploadInfo, err := ms.Client.PutObject(context.Background(), bucketname, objectname, reader, size, opts)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Successfully uploaded bytes: ", uploadInfo)

	return err
}

func (ms *MinioStruct) RetrieveFile(name string, bucketname string, opts *minio.GetObjectOptions) (*GetFileResponse, error) {

	file, err := ms.Client.GetObject(ms.Ctx, bucketname, name, *opts)

	if err != nil {
		log.Println("This file does not exist")

		return nil, err
	}

	info, err := file.Stat()

	if err != nil {
		log.Println("This file info not exist:  ", err)

		return nil, err
	}

	bytesfile, err := ioutil.ReadAll(file)

	if err != nil {
		log.Println("This failed to info not exist:  ", err)
	}

	response := GetFileResponse{
		Name:  name,
		Size:  info.Size,
		Bytes: bytesfile,
	}

	return &response, nil

}
func (ms *MinioStruct) GetBucket(type_ string) string {

	switch type_ {

	case "pdf":
		return ms.Exec.Buckets.PDF
	case "img":
		return ms.Exec.Buckets.IMG
	case "csv":
		return ms.Exec.Buckets.CSV
	case "txt":
		return ms.Exec.Buckets.Txt
	default:
		return ms.Exec.Buckets.Txt
	}
}

func (ms *MinioStruct) DeleteFile(file DeleteFileArg) (*DeleteFileResponse, error) {

	name := file.Name
	bucketname := ms.GetBucket(file.Type)
	opts := minio.RemoveObjectOptions{}
	err := ms.Client.RemoveObject(context.Background(), bucketname, name, opts)

	if err != nil {
		fmt.Println(err)
		return &DeleteFileResponse{Good: false, Name: name}, err
	}

	response := DeleteFileResponse{
		Name: name,
		Good: true,
	}

	return &response, nil

}

func (ms *MinioStruct) SaveFile(tosave *SaveFileArg) (string, error) {

	//ctx := context.Background()
	opts := minio.PutObjectOptions{ContentType: "application/octet-stream"}

	switch tosave.Type {
	case "txt":

		name := ms.Exec.Buckets.GenerateTXT()

		err := ms.DoSave(ms.Exec.Buckets.Txt, name, tosave.Bytes, tosave.Size, opts)

		if err != nil {

			fmt.Println(err)

			return "", err
		}
		return name, nil
	case "pdf":
		//opts := minio.PutObjectOptions{ContentType: "application/pdf"}
		name := ms.Exec.Buckets.GeneratePDF()

		err := ms.DoSave(ms.Exec.Buckets.PDF, name, tosave.Bytes, tosave.Size, opts)

		if err != nil {

			fmt.Println(err)

			return "", err
		}

		return name, nil
	case "img":
		//opts := minio.PutObjectOptions{ContentType: "application/png"}
		name := ms.Exec.Buckets.GenerateIMG()

		err := ms.DoSave(ms.Exec.Buckets.IMG, name, tosave.Bytes, tosave.Size, opts)

		if err != nil {

			fmt.Println(err)

			return "", err
		}
		return name, nil
	case "csv":
		//	opts := minio.PutObjectOptions{ContentType: "application/csv"}
		name := ms.Exec.Buckets.GenerateCSV()

		err := ms.DoSave(ms.Exec.Buckets.CSV, name, tosave.Bytes, tosave.Size, opts)

		if err != nil {

			fmt.Println(err)

			return "", err
		}

		return name, nil

	default:
		return "", errors.New("unkown object received")

	}

	//return errors.New("unkown object received")

}

func (ms *MinioStruct) Init() error {

	for _, v := range ms.Exec.OutBuckets() {

		min := ms.Client
		err := min.MakeBucket(ms.Ctx, v, minio.MakeBucketOptions{Region: ms.Location})

		//err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			// Check to see if we already own this bucket (which happens if you run this twice)
			exists, errBucketExists := min.BucketExists(ms.Ctx, v)
			if errBucketExists == nil && exists {
				log.Printf("We already own %s\n", v)
			} else {
				return err
				//log.Fatalln(err)
			}
		} else {
			log.Printf("Successfully created %s\n", v)

		}
	}

	//	GlobalMin :=
	return nil

}

func NewMinioStruct(cfg *config.BigConfig, exec *executor.Executor) (*MinioStruct, error) {

	out := MinioStruct{}

	ctx := context.Background()

	// endpoint_strintg := "127.0.0.1:9000"
	// accesskey := "someaccesskey"
	// secretkey := "somesecretkey"
	endpoint_strintg := cfg.MinioConfig.AccessUrl
	accesskey := cfg.MinioConfig.Accesspoint
	secretkey := cfg.MinioConfig.PrivateKey
	useSSl := cfg.MinioConfig.UseSSL

	// bucketName := "mymusic"
	// location := "us-east-1"
	// fmt.Println("")

	creds := credentials.NewStaticV4(accesskey, secretkey, "")
	opt := minio.Options{
		Creds:  creds,
		Secure: useSSl,
	}

	min, err := minio.New(endpoint_strintg, &opt)

	if err != nil {
		msg := fmt.Sprint("error creation min, ", err)

		return nil, errors.New(msg)
	}

	// err = min.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})

	// //err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	// if err != nil {
	// 	// Check to see if we already own this bucket (which happens if you run this twice)
	// 	exists, errBucketExists := min.BucketExists(ctx, bucketName)
	// 	if errBucketExists == nil && exists {
	// 		log.Printf("We already own %s\n", bucketName)
	// 	} else {
	// 		log.Fatalln(err)
	// 	}
	// } else {
	// 	log.Printf("Successfully created %s\n", bucketName)

	// }

	// ex, _ := os.Executable()

	// basedir := filepath.Dir(ex)

	// // Upload the zip file
	// objectName := "test.pdf"
	// outname := "out.pdf"
	// filePath := filepath.Join(basedir, "data", "sfiles", objectName)
	// outPath := filepath.Join(basedir, "data", "sfiles", outname)
	// contentType := "application/pdf"

	// // Upload the zip file with FPutObject
	// info, err := min.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Println(info.ETag)
	// log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	// log.Println(info.Key)

	// err = min.FGetObject(context.Background(), bucketName, objectName, outPath, minio.GetObjectOptions{})
	// if err != nil {
	// 	fmt.Println(err)

	// }

	// log.Println("fileget sucess")

	out.Client = min
	out.Exec = exec
	out.Location = cfg.MinioConfig.Location
	out.Ctx = ctx

	return &out, nil
}
