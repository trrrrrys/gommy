package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/google/uuid"
)

type GommySetting struct {
	Name        string         `yaml:"Name"`
	Description string         `yaml:"Description"`
	Repeat      bool           `yaml:"Repeat"`
	Times       int            `yaml:"Times"`
	Output      string         `yaml:"Output"`
	SQLSetting  *SQLSetting    `yaml:"SQL"`
	Columns     []*GommyColumn `yaml:"Columns"`
}

func (s *GommySetting) Verify() error {
	if !s.Repeat {
		s.Times = 1
	}
	for _, v := range s.Columns {
		if v.V != nil && v.V.Max < v.V.Min {
			return fmt.Errorf("")
		}
	}
	return nil
}

type GommyColumn struct {
	Column    string      `yaml:"Column"`
	Type      string      `yaml:"Type"`
	V         *GommyValue `yaml:"Value"`
	enclosure string
}

type SQLSetting struct {
	TableName string `yaml:"TableName"`
}

func (c *GommyColumn) Value() string {
	v := c.V
	switch c.Type {
	// uuid4
	case ValueTypeUUID4:
		return c.enclosure + uuid.New().String() + c.enclosure

	// string value
	case ValueTypeString:
		if v.Const != "" {
			return c.enclosure + v.Const + c.enclosure
		}
		switch v.Choice {
		case ChoiceTypeRandom:
			return c.enclosure + v.In[rand.Int()%len(v.In)] + c.enclosure
		case ChoiceTypeOrder:
			return c.enclosure + v.In[rand.Int()%len(v.In)] + c.enclosure
		case ChoiceTypeReverse:
			return c.enclosure + v.In[rand.Int()%len(v.In)] + c.enclosure
		default:
			return c.enclosure + v.In[0] + c.enclosure
		}
	case ValueTypeBool:
		if v.Const != "" {
			return v.Const
		}
		return "false"

	case ValueTypeInt, ValueTypeInt64, ValueTypeUint, ValueTypeUint64:
		if v.Const != "" {
			return v.Const
		}
		if v.Min != 0 || v.Max != 0 {
			switch v.Choice {
			case ChoiceTypeRandom:
				return strconv.Itoa(rand.Intn(v.Max-v.Min) + v.Min)
			default:
				return strconv.Itoa(v.Min)
			}
		}
		switch v.Choice {
		case ChoiceTypeRandom:
			return v.In[rand.Int()%len(v.In)]
		default:
			return v.In[0]
		}
	default:
		return ""
	}
}

const (
	ChoiceTypeRandom  = "random"
	ChoiceTypeOrder   = "order"
	ChoiceTypeReverse = "reverse"
)

type GommyValue struct {
	Const string `yaml:"Const"`

	Min int `yaml:"Min"`
	Max int `yaml:"Max"`

	Choice string   `yaml:"Choice"`
	In     []string `yaml:"In"`
}
