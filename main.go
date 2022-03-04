package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

var (
	filePath  = flag.String("f", "./gommy.yaml", "yaml file path")
	separator = flag.String("s", ",", "output type")
)

const (
	ValueTypeUUID4  = "uuid4"
	ValueTypeString = "string"
	ValueTypeBool   = "bool"
	ValueTypeInt    = "int"
	ValueTypeInt64  = "int64"
	ValueTypeUint   = "uint"
	ValueTypeUint64 = "uint64"
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
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.Lshortfile)
}

func main() {
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
	_ = f
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
	q := "insert into %s ( %s ) values %s;"
	var columns, valuesQuery string
	for _, v := range s.Columns {
		v.enclosure = EnclosureDoubleQuotation
		columns += v.Column + ","
	}
	columns = columns[:len(columns)-1]
	for i := 0; i < s.Times; i++ {
		valuesQuery += "("
		for _, v := range s.Columns {
			valuesQuery += v.Value() + ","
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
			body += v.Value() + ","
		}
		body = body[:len(body)-1]
	}
	return fmt.Sprintf("%s%s", header, body)
}
