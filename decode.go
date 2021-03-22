package ru_doc_code

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const tagName string = "code"

const lengthTagParam = "length"

type InvalidUnmarshalError struct {
	Err error
}

func (i *InvalidUnmarshalError) Error() string {
	return "code: failed unmarshal" + i.Err.Error()
}


func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{errors.New("unmarshal non-pointer")}
	}

	t := reflect.TypeOf(v)

	var pos int

	fmt.Println(t.NumField())
	for i:=0; i< t.NumField(); i++ {
		field := t.Field(i)
		fmt.Println(field)
		tag := field.Tag.Get(tagName)
		nums := strings.Split(tag, "=")
		fmt.Println(field)
		if len(nums) == 0 && len(nums) > 2 {
			continue
		}

		fmt.Println(lengthTagParam)
		// does not supported tag params
		if nums[0] != lengthTagParam {
			continue
		}

		fmt.Println(field.Name, "tag", tag, field.Type)

		//todo
		shift, err := strconv.Atoi(nums[1])
		if err != nil {
			return err
		}

		//rnv := reflect.New(field.Type)
		//fmt.Println(field.Type.ConvertibleTo(reflect.TypeOf(string(data[pos:pos+shift]))), data[pos:pos+shift], string(data[pos:pos+shift]))
		//fmt.Println(field.Type)
		//reflect.ValueOf()
		sl := data[pos:pos+shift]

		var rnvValue reflect.Value
		switch field.Type.String() {
		case "string":
			rnvValue = reflect.ValueOf(string(sl))
		case "ru_doc_code.TaxRegionCode":
			rnvValue = reflect.ValueOf(TaxRegionCode(sl))
		case "ru_doc_code.SerialNumber":
			rnvValue = reflect.ValueOf(SerialNumber(sl))
		case "ru_doc_code.RegistrationReasonCode":
			rnvValue = reflect.ValueOf(RegistrationReasonCode(sl))
		default:
			return &InvalidUnmarshalError{
				Err: errors.New("does not compare input and data"),
			}
		}
		pos += shift

		fmt.Println(rnvValue, pos)
		fmt.Println("rnvValue", rnvValue, "rv", rv.Elem().Field(i))
		//rv.Elem().Field(i).Set(rnvValue)
	}

	//todo: warning
	if len(data) != pos {
	}

	return nil
}
