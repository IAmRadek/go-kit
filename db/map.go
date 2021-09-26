package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type MapBool map[string]bool

func (p MapBool) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *MapBool) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]bool)
	if !ok {
		return errors.New("Type assertion .(map[string]bool) failed.")
	}

	return nil
}

type MapString map[string]string

func (p MapString) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *MapString) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]string)
	if !ok {
		return errors.New("Type assertion .(map[string]string) failed.")
	}

	return nil
}
