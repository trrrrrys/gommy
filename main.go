package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/goccy/go-yaml"
)

var (
	filePath  = flag.String("f", "./gommy.yaml", "yaml file path")
	separator = flag.String("s", ",", "output type")
)

const (
	ValueTypeUUID4  = "uuid4"
	ValueTypeULID   = "ulid"
	ValueTypeString = "string"
	ValueTypeBool   = "bool"
	ValueTypeNumber = "number"
	ValueTypeDate   = "date"
)

const (
	OutputTypeCSV         = "csv"
	OutputTypeMySQLInsert = "mysql_insert"
)

const (
	EnclosureDoubleQuotation = `"`
	EnclosureNone            = ""
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	f, err := os.Create("profile_cpu")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	mf, err := os.Create("profile_memory")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer mf.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(mf); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	f, err := os.ReadFile(*filePath)
	if err != nil {
		return err
	}
	var setting GommySetting
	if err := yaml.Unmarshal(f, &setting); err != nil {
		return err
	}
	if err := setting.Verify(); err != nil {
		return err
	}
	var result string
	switch setting.Output {
	case OutputTypeMySQLInsert:
		result = createMySQLInsert(&setting)
	case OutputTypeCSV:
		result = createCSV(&setting)
	default:
		result = createMySQLInsert(&setting)
	}
	fmt.Fprintln(os.Stdout, result)
	return nil
}

func createMySQLInsert(s *GommySetting) string {
	q := "INSERT INTO %s (%s) VALUES %s;"
	var columns, valuesQuery string
	for _, v := range s.Columns {
		v.enclosure = EnclosureDoubleQuotation
		columns += v.Column + ","
	}
	columns = columns[:len(columns)-1]
	for i := 0; i < s.Times; i++ {
		valuesQuery += "("
		for _, v := range s.Columns {
			valuesQuery += v.Value(i) + ","
		}
		valuesQuery = valuesQuery[:len(valuesQuery)-1]
		valuesQuery += "),"
	}
	valuesQuery = valuesQuery[:len(valuesQuery)-1]
	return fmt.Sprintf(q, s.SQLSetting.TableName, columns, valuesQuery)
}

func createCSV(s *GommySetting) string {
	var header, body string
	for _, v := range s.Columns {
		header += v.Column + *separator
	}
	header = header[:len(header)-1]
	for i := 0; i < s.Times; i++ {
		body += "\n"
		for _, v := range s.Columns {
			body += v.Value(i) + ","
		}
		body = body[:len(body)-1]
	}
	return fmt.Sprintf("%s%s", header, body)
}
