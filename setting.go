package main

import (
	"fmt"
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

type SQLSetting struct {
	TableName string `yaml:"TableName"`
}

const (
	ChoiceTypeRandom  = "random"
	ChoiceTypeOrder   = "order"
	ChoiceTypeReverse = "reverse"
)

func (s *GommySetting) Verify() error {
	if !s.Repeat {
		s.Times = 1
	}
	for _, v := range s.Columns {
		if v.V != nil && v.V.Max < v.V.Min {
			return fmt.Errorf("invalid format setting")
		}
	}
	return nil
}
