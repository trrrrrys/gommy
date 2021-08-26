package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

var (
	f = flag.String("f", "./gommy.yaml", "yaml file path")
)

const (
	ValueTypeUUID4  = "uuid4"
	ValueTypeString = "string"
	ValueTypeBool   = "bool"
	ValueTypeInt    = "int"
	ValueTypeUint   = "uint"
)

type Setting struct {
	Name        string          `yaml:"Name"`
	Description string          `yaml:"Description"`
	Repeat      bool            `yaml:"Repeat"`
	Times       int             `yaml:"Times"`
	TableName   string          `yaml:"TableName"`
	Columns     []SettingColumn `yaml:"Columns"`
}

func (s *Setting) Verify() error {
	for _, v := range s.Columns {
		if v.V.Max < v.V.Min {
			return fmt.Errorf("")
		}
	}
	return nil
}

type SettingColumn struct {
	Column string          `yaml:"Column"`
	Type   string          `yaml:"Type"`
	V      SettingRowValue `yaml:"Value"`
}

func (c *SettingColumn) Value() string {
	v := c.V
	switch c.Type {
	case ValueTypeUUID4:
		return `"` + uuid.New().String() + `"`
	case ValueTypeString:
		if v.Const != "" {
			return `"` + v.Const + `"`
		}
		switch v.Choice {
		case "random":
			return `"` + v.In[rand.Int()%len(v.In)] + `"`
		default:
			return `"` + v.In[0] + `"`
		}
	default:
		if v.Const != "" {
			return v.Const
		}
		if v.Min != 0 || v.Max != 0 {
			switch v.Choice {
			case "random":
				return strconv.Itoa(rand.Intn(v.Max-v.Min) + v.Min)
			default:
				return strconv.Itoa(v.Min)
			}
		}
		switch v.Choice {
		case "random":
			return v.In[rand.Int()%len(v.In)]
		default:
			return v.In[0]
		}
	}
}

type SettingRowValue struct {
	Const string `yaml:"Const"`

	Min int `yaml:"Min"`
	Max int `yaml:"Max"`

	Choice string   `yaml:"Choice"`
	In     []string `yaml:"In"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.Lshortfile)
}

func main() {
	statusCode := 0
	err := runner(context.Background())
	if err != nil {
		fmt.Println(err)
		statusCode = 1
	}
	os.Exit(statusCode)
}

func runner(ctx context.Context) error {
	flag.Parse()
	f, err := os.ReadFile(*f)
	if err != nil {
		return err
	}
	_ = f
	var setting Setting
	if err := yaml.Unmarshal(f, &setting); err != nil {
		return err
	}
	if !setting.Repeat {
		setting.Times = 1
	}
	q := "insert into %s ( %s ) values %s;"
	var columns, valuesQuery string
	for _, v := range setting.Columns {
		columns += v.Column + ","
	}
	columns = columns[:len(columns)-1]
	for i := 0; i < setting.Times; i++ {
		valuesQuery += "("
		for _, v := range setting.Columns {
			valuesQuery += v.Value() + ","
		}
		valuesQuery = valuesQuery[:len(valuesQuery)-1]
		valuesQuery += "),"
	}
	valuesQuery = valuesQuery[:len(valuesQuery)-1]
	fmt.Printf(q, setting.TableName, columns, valuesQuery)
	return nil
}
