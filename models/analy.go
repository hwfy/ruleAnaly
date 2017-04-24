package models

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

type analyServer struct {
	Name string
	Data map[string]interface{}
}

var err string = "Rules do not match"

func Run(data []byte) (string, error) {
	as := new(analyServer)

	if err := json.Unmarshal(data, as); err != nil {
		return "", errors.New("The business data type is incorrect: " + err.Error())
	}

	fun := reflect.ValueOf(as).MethodByName(strings.Title(as.Name))
	if !fun.IsValid() {
		return "", errors.New("The parsing service does not exist: " + as.Name)
	}
	param := make([]reflect.Value, 1)
	param[0] = reflect.ValueOf(as.Data)

	replys := fun.Call(param)

	if !replys[1].IsNil() {
		return "", replys[1].Interface().(error)
	}
	return replys[0].String(), nil
}

func contains(ss []interface{}, s interface{}) bool {
	for _, v := range ss {
		if reflect.TypeOf(v).Kind() == reflect.Int {
			v = float64(v.(int))
		}
		if s == v {
			return true
		}
	}
	return false
}
