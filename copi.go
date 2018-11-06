package copi

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"

	logrus "github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
}

func Debugging() {
	log.SetLevel(logrus.DebugLevel)
}

func Dup(from interface{}, to interface{}) error {
	return copy(reflect.ValueOf(from), reflect.ValueOf(to))
}

func copiError(err error) error {
	return fmt.Errorf("copi: %s", err)
}

func initNilValue(v reflect.Value) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		init := reflect.New(v.Type().Elem())
		log.Debug("init: ", v.Type().Elem(), init, " for ", v)
		init.Elem().Set(reflect.Zero(init.Elem().Type()))
		v.Set(init)
		log.Debug("init-done: ", v)
	}
}

func copy(from, to reflect.Value) error {
	if from.Kind() == reflect.Ptr || to.Kind() == reflect.Ptr {
		log.Debug("duping: ", from.Kind(), from, to.Kind(), to)
		initNilValue(to)
		return copy(reflect.Indirect(from), reflect.Indirect(to))
	}

	if to.CanSet() {
		if from.Kind() == reflect.Invalid {
			log.Debug("from invalid: ", from)
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

		if val, ok := from.Interface().(driver.Valuer); ok {
			val, err := val.Value()
			if err != nil {
				return copiError(err)
			}
			return copy(reflect.ValueOf(val), to)
		}

		if to.CanAddr() {
			if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
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

				log.Debug("struct:field:", dstFieldMeta.Name, ":", dstFieldMeta.Anonymous)

				if !dstFieldVal.CanSet() {
					return nil
				}

				if from.Type().Kind() == reflect.Struct {
					var srcFieldVal reflect.Value
					if byTag := dstFieldMeta.Tag.Get("copi"); byTag != "" {
						log.Debug("byTag:", byTag)
						srcFieldVal = from.FieldByName(byTag)
					} else if srcFieldName, avail := srcTags[dstFieldMeta.Name]; avail {
						log.Debug("bySrcTag:", srcFieldName)
						srcFieldVal = from.FieldByName(srcFieldName)
					} else {
						log.Debug("byName:", dstFieldMeta.Name)
						srcFieldVal = from.FieldByName(dstFieldMeta.Name)
					}
					if srcFieldVal.IsValid() {
						copy(srcFieldVal, dstFieldVal)
					}
				}
			}
		case reflect.Slice:
			dstSliceLen := to.Len()

			if from.Type().Kind() == reflect.Slice {
				srcSliceLen := from.Len()
				for i := 0; i < srcSliceLen; i++ {
					srcElemVal := from.Index(i)
					if i < dstSliceLen {
						dstElemVal := to.Index(i)
						copy(srcElemVal, dstElemVal)
					} else {
						log.Debug("dstSlice not enough cap")
						to.Set(reflect.Append(to, reflect.Zero(to.Type().Elem())))
						dstElemVal := to.Index(i)
						copy(srcElemVal, dstElemVal)
						log.Debug("len after append: ", to.Len())
					}
				}
			}
		case reflect.Map:
			log.Debug("to: ", to.Kind(), " is Nill ", to.IsNil())
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
						log.Debug("assign: ", srcElemKey)
						init := reflect.New(to.Type().Elem())
						copy(srcElemVal, init)
						to.SetMapIndex(srcElemKey, init.Elem())
					} else if convert {
						log.Debug("convert: ", srcElemKey)
						init := reflect.New(to.Type().Elem())
						copy(srcElemVal, init)
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
