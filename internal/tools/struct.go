package tools

import (
	"reflect"

	"gorm.io/gorm"
)

func ForEachStruct(data interface{}, callback func(key string, val interface{})) {
	var (
		paramsType  = reflect.TypeOf(data)
		paramsValue = reflect.ValueOf(data)
	)
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)
		value := paramsValue.Field(i)
		switch value.Interface().(type) {
		case gorm.Model:
			for ii := 0; ii < value.NumField(); ii++ {
				x := field.Type.Field(ii)
				y := value.Field(ii)
				callback(x.Name, y.Interface())
			}
		default:
			callback(field.Name, value.Interface())
		}
	}
}
