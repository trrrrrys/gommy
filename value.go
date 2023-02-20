package main

import "math/rand"

type GommyValue struct {
	Const string `yaml:"Const"`

	Choice string   `yaml:"Choice"`
	Min    int      `yaml:"Min"`
	Max    int      `yaml:"Max"`
	In     []string `yaml:"In"`
}

func (v *GommyValue) ChoiceRange(n int) int {
	if v.Min != 0 || v.Max != 0 {
		r := v.Max - v.Min + 1
		switch v.Choice {
		case ChoiceTypeRandom:
			return rand.Intn(r) + v.Min
		case ChoiceTypeOrder:
			return n%r + v.Min
		case ChoiceTypeReverse:
			return v.Max - n%r
		default:
			return v.Min
		}
	}
	return -1
}

func (v *GommyValue) ChoiceIn(n int) string {
	l := len(v.In)
	switch v.Choice {
	case ChoiceTypeRandom:
		return v.In[rand.Int()%l]
	case ChoiceTypeOrder:
		return v.In[n%l]
	case ChoiceTypeReverse:
		return v.In[(l-1)-n%l]
	default:
		return v.In[0]
	}
}
