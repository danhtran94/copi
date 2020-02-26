package copi

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
)

type Option struct {
	SQLScanner   bool
	DriverValuer bool
}

func Dup(from interface{}, to interface{}) error {
	defaultOpt := Option{
		SQLScanner:   false,
		DriverValuer: false,
	}

	return DupWithOpt(from, to, defaultOpt)
	// return copy(reflect.ValueOf(from), reflect.ValueOf(to), defaultOpt)
}

func DupWithOpt(from interface{}, to interface{}, opt Option) error {
	return copy(reflect.ValueOf(from), reflect.ValueOf(to), opt)
}

func copiError(err error) error {
	return fmt.Errorf("copi: %s", err)
}

// initNilValue initial poiter value if it nil
func initNilValue(v reflect.Value) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		init := reflect.New(v.Type().Elem())
		init.Elem().Set(reflect.Zero(init.Elem().Type()))
		v.Set(init)
	}
}

func copy(from, to reflect.Value, opt Option) error {
	if from.Kind() == reflect.Ptr || to.Kind() == reflect.Ptr {
		if from.Kind() == reflect.Ptr && from.IsNil() {
			return nil
		}
		initNilValue(to)
		return copy(reflect.Indirect(from), reflect.Indirect(to), opt)
	}

	if to.CanSet() {
		if from.Kind() == reflect.Invalid {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}

		if from.Type().AssignableTo(to.Type()) {
			to.Set(from)
			return nil
		}

		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
			return nil
		}

		if val, ok := from.Interface().(driver.Valuer); ok && opt.DriverValuer {
			val, err := val.Value()
			if err != nil {
				return copiError(err)
			}
			return copy(reflect.ValueOf(val), to, opt)
		}

		if to.CanAddr() {
			if scanner, ok := to.Addr().Interface().(sql.Scanner); ok && opt.SQLScanner {
				err := scanner.Scan(from.Interface())
				if err != nil {
					return copiError(err)
				}
				return nil
			}
		}

		srcTags := scanTags(from.Type())

		switch to.Type().Kind() {
		case reflect.Struct:
			for _, dstFieldMeta := range deepFields(to.Type()) {
				dstFieldVal := to.FieldByName(dstFieldMeta.Name)

				if !dstFieldVal.CanSet() {
					return nil
				}

				if from.Type().Kind() == reflect.Struct {
					var srcFieldVal reflect.Value
					if byTag := dstFieldMeta.Tag.Get("copi"); byTag != "" {
						srcFieldVal = from.FieldByName(byTag)
					} else if srcFieldName, avail := srcTags[dstFieldMeta.Name]; avail {
						srcFieldVal = from.FieldByName(srcFieldName)
					} else {
						srcFieldVal = from.FieldByName(dstFieldMeta.Name)
					}
					if srcFieldVal.IsValid() {
						copy(srcFieldVal, dstFieldVal, opt)
					}
				}
			}
		case reflect.Slice:
			dstSliceLen := to.Len()

			if from.Type().Kind() == reflect.Slice {
				srcSliceLen := from.Len()

				if !from.IsNil() && to.IsNil() {
					to.Set(reflect.MakeSlice(to.Type(), 0, 0))
				}

				for i := 0; i < srcSliceLen; i++ {
					srcElemVal := from.Index(i)
					if i < dstSliceLen {
						dstElemVal := to.Index(i)
						copy(srcElemVal, dstElemVal, opt)
					} else {
						to.Set(reflect.Append(to, reflect.Zero(to.Type().Elem())))
						dstElemVal := to.Index(i)
						copy(srcElemVal, dstElemVal, opt)
					}
				}
			}
		case reflect.Map:
			if to.IsNil() {
				to.Set(reflect.MakeMap(to.Type()))
			}

			dstKeyType := to.Type().Key()
			if from.Type().Kind() == reflect.Map {
				srcKeyType := from.Type().Key()

				convert := srcKeyType.ConvertibleTo(dstKeyType)
				assign := srcKeyType.AssignableTo(dstKeyType)

				for _, srcElemKey := range from.MapKeys() {
					srcElemVal := from.MapIndex(srcElemKey)
					if assign {
						init := reflect.New(to.Type().Elem())
						copy(srcElemVal, init, opt)
						to.SetMapIndex(srcElemKey, init.Elem())
					} else if convert {
						init := reflect.New(to.Type().Elem())
						copy(srcElemVal, init, opt)
						to.SetMapIndex(srcElemKey.Convert(dstKeyType), init.Elem())
					} else {
						return nil
					}
				}
			}
		default:
		}
	}

	return nil
}

func scanTags(reflectType reflect.Type) map[string]string {
	srcTags := map[string]string{}

	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if tag := v.Tag.Get("copi-to"); tag != "" {
				srcTags[tag] = v.Name
			}
		}
	}

	return srcTags
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}
