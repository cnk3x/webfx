package config

import (
	"encoding/json"
	"os"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// Select returns the first non-zero arguments. Arguments must be comparable.
func Select[T comparable](values ...T) (out T) {
	out, _ = lo.Coalesce(values...)
	return
}

// IifF is a 1 line if/else statement whose options are functions
func Iif[T any](condition bool, ifFunc, elseFunc T) T {
	return lo.Ternary(condition, ifFunc, elseFunc)
}

// IifF is a 1 line if/else statement whose options are functions
func IifF[T any](condition bool, ifFunc, elseFunc func() T) T {
	return lo.TernaryF(condition, ifFunc, elseFunc)
}

// ReadYAMLFile 从文件中解析配置
func ReadYAMLFile[T any](path string, out *T) (err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return
	}
	return yaml.Unmarshal(data, out)
}

// ReadJSONFile 从文件中解析配置
func ReadJSONFile(path string, out any) (err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return
	}
	return json.Unmarshal(data, out)
}

// WriteYAMLFile 将配置写入文件
func WriteYAMLFile(path string, in any) (err error) {
	var data []byte
	if data, err = yaml.Marshal(in); err != nil {
		return
	}
	os.WriteFile(path, data, 0o644)
	return
}

// WriteJSONFile 将配置写入文件
func WriteJSONFile(path string, in any) (err error) {
	var data []byte
	if data, err = json.Marshal(in); err != nil {
		return
	}
	os.WriteFile(path, data, 0o644)
	return
}

// UnmarshalYAMLString 从字符串中解析配置
func UnmarshalYAMLString[T any](data string, out *T) (err error) {
	return yaml.Unmarshal([]byte(data), out)
}

// YamlFileToJSON read yaml file, unmarshal it and convert to gjson value
func YamlFileToJSON(path string) (r Value) {
	var dataAny any
	if err := ReadYAMLFile(path, &dataAny); err == nil {
		data, _ := json.Marshal(dataAny)
		r = ParseJSON(data)
	}
	return
}
