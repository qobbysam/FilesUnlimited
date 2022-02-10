package executor

import (
	"github.com/lithammer/shortuuid"
	"github.com/qobbysam/filesunlimited/pkgs/config"
)

//Executor will route requests to the correct buckets
type ExBuckets struct {
	Txt string `json:"txt"`
	PDF string `json:"pdf"`
	CSV string `json:"csv"`
	IMG string `json:"img"`
}

//this will create a uniqie name for the file to be saved
func (ex *ExBuckets) GenerateAName() string {

	return shortuuid.New()
}

func (ex *ExBuckets) GenerateTXT() string {

	name := ex.GenerateAName() + ".txt"
	return name
}

func (ex *ExBuckets) GenerateIMG() string {

	name := ex.GenerateAName() + ".png"
	return name
}

func (ex *ExBuckets) GeneratePDF() string {

	name := ex.GenerateAName() + ".pdf"
	return name
}

func (ex *ExBuckets) GenerateCSV() string {

	name := ex.GenerateAName() + ".csv"
	return name
}

type Executor struct {
	Buckets *ExBuckets
}

func (ex *Executor) OutBuckets() []string {

	return []string{ex.Buckets.CSV, ex.Buckets.PDF, ex.Buckets.Txt, ex.Buckets.IMG}
}

func NewExecutor(cfg *config.BigConfig) *Executor {

	ebuckets := ExBuckets{
		Txt: cfg.BucketConfig.Txt,
		IMG: cfg.BucketConfig.IMG,
		PDF: cfg.BucketConfig.PDF,
		CSV: cfg.BucketConfig.CSV,
	}

	out := Executor{Buckets: &ebuckets}

	GlobalExecutor = &out
	return &out

}

var GlobalExecutor *Executor
