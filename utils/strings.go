package utils

import (
	"reflect"
	"unsafe"
)

func StringToByte(source string) (target []byte) {
	str := *(*reflect.StringHeader)(unsafe.Pointer(&source))
	sli := (*reflect.SliceHeader)(unsafe.Pointer(&target))
	sli.Len, sli.Data, sli.Cap = str.Len, str.Data, 0
	return target
}

func ByteToString(source []byte) string {
	return *(*string)(unsafe.Pointer(&source))
}
