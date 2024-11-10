package string_utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Global map that map a type/kind to a method that converts string to that type
var fromstrings = make(map[reflect.Kind]func(string, []reflect.Kind) reflect.Value)

const (
	errorValue = "Failed to convert string to instance:"
)

type StructInstantiationProvider interface {
	NewInstance(string) (interface{}, error)
}

var Provider StructInstantiationProvider

// initialize the map
func init() {
	fromstrings[reflect.String] = stringFromString
	fromstrings[reflect.Int] = intFromString
	fromstrings[reflect.Int8] = int8FromString
	fromstrings[reflect.Int16] = int16FromString
	fromstrings[reflect.Int32] = int32FromString
	fromstrings[reflect.Int64] = int64FromString
	fromstrings[reflect.Uint] = uintFromString
	fromstrings[reflect.Uint8] = uint8FromString
	fromstrings[reflect.Uint16] = uint16FromString
	fromstrings[reflect.Uint32] = uint32FromString
	fromstrings[reflect.Uint64] = uint64FromString
	fromstrings[reflect.Float32] = float32FromString
	fromstrings[reflect.Float64] = float64FromString
	fromstrings[reflect.Bool] = boolFromString
	fromstrings[reflect.Ptr] = ptrFromString
	fromstrings[reflect.Slice] = sliceFromString
	fromstrings[reflect.Map] = mapFromString
	fromstrings[reflect.Interface] = interfaceFromString
	fromstrings[reflect.Struct] = structFromString
}

// Comvert string to string
func stringFromString(str string, kinds []reflect.Kind) reflect.Value {
	return reflect.ValueOf(str)
}

// Convert string to int
func intFromString(str string, kinds []reflect.Kind) reflect.Value {
	if str == "" {
		return reflect.ValueOf(0)
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		er(err, "int")
	}
	return reflect.ValueOf(i)
}

// Convert string to int8
func int8FromString(str string, kinds []reflect.Kind) reflect.Value {
	i, err := strconv.Atoi(str)
	if err != nil {
		er(err, "int8")
	}
	return reflect.ValueOf(int8(i))
}

// Convert string to int16
func int16FromString(str string, kinds []reflect.Kind) reflect.Value {
	i, err := strconv.Atoi(str)
	if err != nil {
		er(err, "int16")
	}
	return reflect.ValueOf(int16(i))
}

// Convert string to int32
func int32FromString(str string, kinds []reflect.Kind) reflect.Value {
	if str == "" {
		return reflect.ValueOf(int32(0))
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		er(err, "int32")
	}
	return reflect.ValueOf(int32(i))
}

// Convert string to int64
func int64FromString(str string, kinds []reflect.Kind) reflect.Value {
	i, err := strconv.Atoi(str)
	if err != nil {
		er(err, "int64")
	}
	return reflect.ValueOf(int64(i))
}

// Convert string to uint
func uintFromString(str string, kinds []reflect.Kind) reflect.Value {
	i, err := strconv.Atoi(str)
	if err != nil {
		er(err, "uint")
	}
	return reflect.ValueOf(uint(i))
}

// Convert string to uint8
func uint8FromString(str string, kinds []reflect.Kind) reflect.Value {
	if str == "" {
		return reflect.ValueOf(byte(0))
	}
	return reflect.ValueOf([]byte(str))
}

// Convert string to uint16
func uint16FromString(str string, kinds []reflect.Kind) reflect.Value {
	i, err := strconv.Atoi(str)
	if err != nil {
		er(err, "uint16")
	}
	return reflect.ValueOf(uint16(i))
}

// Convert string to uint32
func uint32FromString(str string, kinds []reflect.Kind) reflect.Value {
	i, err := strconv.Atoi(str)
	if err != nil {
		er(err, "uint32")
	}
	return reflect.ValueOf(uint32(i))
}

// Convert string to uint64
func uint64FromString(str string, kinds []reflect.Kind) reflect.Value {
	i, err := strconv.Atoi(str)
	if err != nil {
		er(err, "uint64")
	}
	return reflect.ValueOf(uint64(i))
}

// Convert string to bool
func boolFromString(str string, kinds []reflect.Kind) reflect.Value {
	i, err := strconv.ParseBool(str)
	if err != nil {
		er(err, "bool")
	}
	return reflect.ValueOf(i)
}

// Convert string to float32
func float32FromString(str string, kinds []reflect.Kind) reflect.Value {
	f, err := strconv.ParseFloat(str, 32)
	if err != nil {
		er(err, "float32")
	}
	return reflect.ValueOf(float32(f))
}

// Convert string to float64
func float64FromString(str string, kinds []reflect.Kind) reflect.Value {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		er(err, "float64")
	}
	return reflect.ValueOf(float64(f))
}

// Convert string to pointer
func ptrFromString(str string, kinds []reflect.Kind) reflect.Value {
	f := fromstrings[kinds[0]]
	if f != nil {
		v := f(str, kinds[1:])
		if !v.IsValid() {
			return v
		}
		newPtr := reflect.New(v.Type())
		newPtr.Elem().Set(v)
		return newPtr
	}
	er(errors.New("Pointer cloud not be created for kind "+kinds[0].String()), "ptr")
	return reflect.ValueOf(nil)
}

// Convert string to interface
func interfaceFromString(str string, kinds []reflect.Kind) reflect.Value {
	f := fromstrings[kinds[0]]
	if f != nil {
		v := f(str, kinds[1:])
		newInt := reflect.New(v.Type())
		newInt.Set(v)
		return newInt
	}
	er(errors.New("Pointer cloud not be created for kind "+kinds[0].String()), "interface")
	return reflect.ValueOf(nil)
}

// Convert string to map
func mapFromString(str string, kinds []reflect.Kind) reflect.Value {
	str = strings.TrimSpace(str)
	str = str[1 : len(str)-1]
	items := strings.Split(str, ",")
	var newMap *reflect.Value
	for _, item := range items {
		index := strings.Index(item, "=")
		if index == -1 {
			er(errors.New("Item '"+item+"' does not contain a '=' sign"), "map")
			continue
		}
		keyStr := strings.TrimSpace(item[0:index])
		valueStr := strings.TrimSpace(item[index+1:])
		keyF := fromstrings[kinds[0]]
		valueF := fromstrings[kinds[1]]
		if keyF == nil || valueF == nil {
			er(errors.New("Item '"+item+"' cannot find either the key type or the value type converter"), "map")
			continue
		}
		keyV := keyF(keyStr, kinds[2:])
		valueV := valueF(valueStr, kinds[2:])
		if newMap == nil {
			m := reflect.MakeMap(reflect.MapOf(keyV.Type(), valueV.Type()))
			newMap = &m
		}
		newMap.SetMapIndex(keyV, valueV)
	}
	return *newMap
}

// Convert string to slice
func sliceFromString(str string, kinds []reflect.Kind) reflect.Value {
	str = strings.TrimSpace(str)
	// if it is byte array, it will not have square brackets
	if len(str) > 1 && str[0] == '[' {
		str = str[1 : len(str)-1]
	}
	items := strings.Split(str, ",")

	itemF := fromstrings[kinds[0]]
	if itemF == nil {
		er(errors.New("Cannot find converter item kind "+kinds[0].String()), "slice")
		return reflect.ValueOf(nil)
	}

	defaultValue := itemF("", kinds[1:])

	if str == "" {
		return reflect.MakeSlice(reflect.SliceOf(defaultValue.Type()), 0, 0)
	}

	//Special case for byte array
	if defaultValue.Kind() == reflect.Uint8 {
		newSlice := reflect.MakeSlice(reflect.SliceOf(defaultValue.Type()), len(str), len(str))
		for i, v := range str {
			newSlice.Index(i).Set(reflect.ValueOf(byte(v)))
		}
		return newSlice
	}

	newSlice := reflect.MakeSlice(reflect.SliceOf(defaultValue.Type()), len(items), len(items))

	for i, item := range items {
		v := itemF(item, kinds[1:])
		newSlice.Index(i).Set(v)
	}
	return newSlice
}

// Convert string to an instance
func InstanceOf(str string) interface{} {
	v := FromString(str)
	if v.IsValid() {
		return v.Interface()
	}
	return nil
}

func structFromString(str string, kinds []reflect.Kind) reflect.Value {
	if Provider == nil {
		New("No struct instantiation provider available to instantiate ", str).Panic()
	}
	if str == "<Nil>" {
		return reflect.ValueOf(nil)
	}
	v, e := Provider.NewInstance(str)
	if e != nil {
		panic("Failed to instantiate struct " + str + ", please check that you registered it in the registry. " + e.Error())
	}
	return reflect.ValueOf(v)
}

// Conver string to a reflect.value
func FromString(str string) reflect.Value {
	if str == "" || str == "{0}" {
		return reflect.ValueOf(nil)
	}
	v, k := parseStringForKinds(str)
	f := fromstrings[k[0]]
	if f == nil {
		New("no converter was found for kind ", k[0].String()).Panic()
	}
	return f(v, k[1:])
}

// Extract the kinds from the prefix of the string
func parseStringForKinds(str string) (string, []reflect.Kind) {
	if len(str) < 3 {
		New("'", str, "'lenght is less than 3, which means it is not in the correct format of {kind}...").Panic()
	}
	if str[0] != '{' {
		New("'", str, "' does not start with '{'").Panic()
	}
	index := strings.Index(str, "}")
	if index == -1 {
		New("'", str, "'does not have a closing '}'").Panic()
	}
	types := str[1:index]
	result := str[index+1:]
	k := parseKinds(types)
	return result, k
}

// extract the kinds to a list of reflect.Kind
func parseKinds(types string) []reflect.Kind {
	split := strings.Split(types, ",")
	kinds := make([]reflect.Kind, len(split))
	for i, v := range split {
		k, e := strconv.Atoi(v)
		if e != nil {
			New("Error parsing kind:", v).Panic()
		}
		kinds[i] = reflect.Kind(k)
	}
	return kinds
}

func er(err error, tag string) {
	fmt.Println(errorValue, tag, ":", err)
}
