package cdek_pay

import (
	"fmt"
	"reflect"
)

// FlattenStructToMap Функция для преобразования структуры в плоскую map[string]interface{}
func FlattenStructToMap(data interface{}, prefix string) map[string]interface{} {
	flatMap := make(map[string]interface{})
	value := reflect.ValueOf(data)
	dataType := reflect.TypeOf(data)

	// Проверяем, является ли value nil
	if !value.IsValid() || (value.Kind() == reflect.Ptr && value.IsNil()) {
		return flatMap
	}

	// Если это указатель, получаем значение, на которое он указывает
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		dataType = dataType.Elem()
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := dataType.Field(i)
		key := prefix + fieldType.Tag.Get("json")

		// Пропускаем поля со значением nil
		if !field.IsValid() || (field.Kind() == reflect.Ptr && field.IsNil()) {
			continue
		}

		// Если поле — указатель, получаем значение
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		switch field.Kind() {
		case reflect.Struct:
			// Рекурсивный вызов для вложенных структур
			for k, v := range FlattenStructToMap(field.Interface(), key+".") {
				flatMap[k] = v
			}
		case reflect.Slice:
			// Обработка срезов
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j).Interface()
				for k, v := range FlattenStructToMap(elem, fmt.Sprintf("%s.%d.", key, j)) {
					flatMap[k] = v
				}
			}
		default:
			// Прямое добавление полей, если значение не nil
			flatMap[key] = field.Interface()
		}
	}
	return flatMap
}
