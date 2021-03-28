package ru_doc_code

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

const tagName string = "code"

const lengthTagParam = "length"

type InvalidUnmarshalError struct {
	Err error
}

func (i *InvalidUnmarshalError) Error() string {
	return "code: failed unmarshal" + i.Err.Error()
}

func Unmarshal(data []byte, src interface{}) error {
	typ := reflect.TypeOf(src)
	if typ.Kind() != reflect.Ptr {
		return &InvalidUnmarshalError{errors.New("source is not pointer")}
	}

	if typ = typ.Elem(); typ.Kind() != reflect.Struct {
		return &InvalidUnmarshalError{errors.New("source is not struct")}
	}
	if onlyPrivateFields(typ) {
		return &InvalidUnmarshalError{errors.New("source struct has not exported fields")}
	}

	var pos int
	srcRv := reflect.ValueOf(src).Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		tag := field.Tag.Get(tagName)
		nums := strings.Split(tag, "=")

		//todo: add 4,5,6 or 4|5 length value
		if len(nums) != 2 {
			continue
		}

		if nums[0] != lengthTagParam {
			continue
		}

		shift, err := strconv.Atoi(nums[1])
		if err != nil {
			return err
		}

		sl := data[pos : pos+shift]
		fv := srcRv.Field(i)

		switch field.Type.Kind() {
		case reflect.String:
			str := string(sl)
			v := reflect.NewAt(fv.Type(), unsafe.Pointer(reflect.ValueOf(&str).Elem().UnsafeAddr())).Elem()
			fv.Set(v)
		default:
		}
		pos += shift
	}

	//todo: warning message but not correct
	if len(data) != pos {
	}

	return nil
}

func onlyPrivateFields(typ reflect.Type) bool {
	for i := 0; i < typ.NumField(); i++ {
		tf := typ.Field(i).Type
		if typ.Field(i).Type.Kind() == reflect.Ptr {
			tf = typ.Field(i).Type.Elem()
		}
		if tf.Kind() == reflect.Struct && !isPrivateField(typ.Field(i)) {
			if !onlyPrivateFields(tf) {
				return false
			}
		} else if !isPrivateField(typ.Field(i)) {
			return false
		}
	}
	return true
}

func isPrivateField(sf reflect.StructField) bool {
	return len(sf.PkgPath) != 0 && !sf.Anonymous
}
