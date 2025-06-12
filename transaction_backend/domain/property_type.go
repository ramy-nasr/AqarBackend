package domain

import "errors"

type PropertyType string

const (
	Land      PropertyType = "Land"
	Villa     PropertyType = "Villa"
	Apartment PropertyType = "Apartment"
)

func (t PropertyType) IsValid() bool {
	switch t {
	case Land, Villa, Apartment:
		return true
	default:
		return false
	}
}

func NewPropertyType(value string) (PropertyType, error) {
	pt := PropertyType(value)
	if !pt.IsValid() {
		return "", errors.New("invalid property type")
	}
	return pt, nil
}
