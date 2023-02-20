package main

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type GommyColumn struct {
	Column    string      `yaml:"Column"`
	Type      string      `yaml:"Type"`
	V         *GommyValue `yaml:"Value"`
	enclosure string
}

func (c *GommyColumn) Value(i int) string {
	v := c.V
	switch c.Type {
	// uuid4
	case ValueTypeUUID4:
		return c.enclosure + uuid.New().String() + c.enclosure
	// ulid
	case ValueTypeULID:
		return c.enclosure + newULID() + c.enclosure
	// string value
	case ValueTypeString:
		if v.Const != "" {
			return c.enclosure + v.Const + c.enclosure
		}
		return c.enclosure + v.ChoiceIn(i) + c.enclosure
	case ValueTypeBool:
		if v.Const != "" {
			return v.Const
		}
		return "false"
	case ValueTypeNumber:
		if v.Const != "" {
			return v.Const
		}
		rv := v.ChoiceRange(i)
		if rv != -1 {
			return strconv.Itoa(rv)
		}
		return v.ChoiceIn(i)
	case ValueTypeDate:
		if v.Const != "" {
			ut, err := strconv.Atoi(v.Const)
			if err == nil {
				return c.enclosure + time.Unix(int64(ut), 0).Format("2006-01-02 03:04:05") + c.enclosure
			}
		}
		rv := v.ChoiceRange(i)
		if rv != -1 {
			return c.enclosure + time.Unix(int64(rv), 0).Format("2006-01-02 03:04:05") + c.enclosure
		}
		return v.ChoiceIn(i)
	default:
		return ""
	}
}
