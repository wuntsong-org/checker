package checker

import (
	"encoding/json"
	errors "github.com/wuntsong/wterrors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const TagIgnore = "c-ignore"

const TagBoolMust = "cb-must" // bool必须为指定值

const TagIntIgnore = "ci-ignore" // int忽略范围检查
const TagIntMax = "ci-max"       // int最大值
const TagIntMin = "ci-min"       // int最小值
const TagIntZero = "ci-zero"     // int允许为zero
const TagIntCheck = "ci-checker" // 检查函数
const TagIntMust = "ci-must"     // 检查函数

const TagStringJsonNumber = "cs-json-number"
const TagStringLengthMin = "cs-min"
const TagStringLengthMax = "cs-max"
const TagStringLength = "cs-length"
const TagStringZero = "cs-zero"
const TagStringIgnore = "cs-ignore"
const TagStringChecker = "cs-checker"
const TagStringMust = "cs-must"
const TagStringRegex = "cs-regex"

const TagSliceLengthMin = "csl-min"
const TagSliceLengthMax = "csl-max"
const TagSliceLength = "csl-length"
const TagSliceZero = "csl-zero"
const TagSliceIgnore = "csl-ignore"
const TagSliceChecker = "csl-checker"

const TagMapLengthMin = "cm-min"
const TagMapLengthMax = "cm-max"
const TagMapLength = "cm-length"
const TagMapZero = "cm-zero"
const TagMapIgnore = "cm-ignore"
const TagMapChecker = "cm-checker"

func (c *Checker) checkBool(vi bool, field *reflect.StructField) errors.WTError {
	if field == nil {
		return nil
	}

	if field.Tag.Get(TagIgnore) == "true" {
		return nil
	}

	boolMust := field.Tag.Get(TagBoolMust)
	if boolMust == "true" && !vi {
		return ReturnFieldError(field, "bool must be true")
	} else if boolMust == "false" && vi {
		return ReturnFieldError(field, "bool must be true")
	}

	return nil
}

func (c *Checker) checkInt64(vi int64, field *reflect.StructField) errors.WTError {
	if field == nil {
		return nil
	}

	if field.Tag.Get(TagIgnore) == "true" {
		return nil
	}

	intZero := field.Tag.Get(TagIntZero)
	intIgnore := field.Tag.Get(TagIntIgnore)

	if intZero == "notcheck" && vi == 0 { // 如果是零值，忽略后面检查
		return nil // 忽略后续检查
	} else if (intZero == "ignore" && vi == 0) || intIgnore != "true" {
		// 忽略内置检查
	} else {
		mustString := field.Tag.Get(TagIntMust)
		if mustString != "" {
			must, err := strconv.ParseInt(mustString, 10, 64)
			if err != nil {
				return ReturnFieldError(field, err.Error())
			}

			if vi != must {
				return ReturnFieldError(field, "must be %d", must)
			}
		} else {
			tagMax := field.Tag.Get(TagIntMax)
			tagMin := field.Tag.Get(TagIntMin)

			if tagMax != "" {
				intMax, err := strconv.ParseInt(tagMax, 10, 64)
				if err != nil {
					return ReturnFieldError(field, err.Error())
				}

				if vi > intMax {
					return ReturnFieldError(field, "too big")
				}
			}

			if tagMin != "" {
				intMin, err := strconv.ParseInt(tagMin, 10, 64)
				if err != nil {
					return ReturnFieldError(field, err.Error())
				}

				if vi <= intMin {
					return ReturnFieldError(field, "too small")
				}
			}
		}
	}

	checker := strings.Split(field.Tag.Get(TagIntCheck), ",")
	for _, ch := range checker {
		ch = strings.TrimSpace(ch)
		if len(ch) == 0 {
			continue
		}

		fn, ok := c.intChecker[ch]
		if !ok {
			return ReturnFieldError(field, "checker %s not found", ch)
		}

		err := fn(field, vi)
		if err != nil {
			return ReturnFieldError(field, "checker %s: %s", ch, err.Error())
		}
	}

	return nil
}

func (c *Checker) checkString(vi string, field *reflect.StructField) errors.WTError {
	if field == nil {
		return nil
	}

	if field.Tag.Get(TagIgnore) == "true" {
		return nil
	}

	if field.Tag.Get(TagStringJsonNumber) == "true" {
		jsonNumber, err := json.Number(vi).Int64()
		if err != nil {
			return ReturnFieldError(field, err.Error())
		}

		return c.checkInt64(jsonNumber, field)
	}

	stringZero := field.Tag.Get(TagStringZero)
	stringIgnore := field.Tag.Get(TagStringIgnore)

	if stringZero == "notcheck" && vi == "" { // 如果是零值，忽略后面检查
		return nil // 忽略后续检查
	} else if (stringZero == "ignore" && vi == "") || stringIgnore != "true" {
		// 忽略内置检查
	} else {
		must := field.Tag.Get(TagStringMust)
		if must != "" {
			if vi != must {
				return ReturnFieldError(field, "must be %s", must)
			}
		} else {
			tagMax := field.Tag.Get(TagStringLengthMax)
			tagMin := field.Tag.Get(TagStringLengthMin)
			tagLength := field.Tag.Get(TagStringLength)
			tagRegex := field.Tag.Get(TagStringRegex)

			if tagLength != "" {
				stringLength, err := strconv.ParseInt(tagLength, 10, 64)
				if err != nil {
					return ReturnFieldError(field, err.Error())
				}

				if len(vi) != int(stringLength) {
					return ReturnFieldError(field, "bad length")
				}
			}

			if tagMax != "" {
				stringLengthMax, err := strconv.ParseInt(tagMax, 10, 64)
				if err != nil {
					return ReturnFieldError(field, err.Error())
				}

				if len(vi) > int(stringLengthMax) {
					return ReturnFieldError(field, "too long")
				}
			}

			if tagMin != "" {
				stringLengthMin, err := strconv.ParseInt(tagMin, 10, 64)
				if err != nil {
					return ReturnFieldError(field, err.Error())
				}

				if len(vi) <= int(stringLengthMin) {
					return ReturnFieldError(field, "too short")
				}
			}

			if tagRegex != "" {
				r, err := regexp.Compile(tagRegex)
				if err != nil {
					return ReturnFieldError(field, err.Error())
				}

				if !r.MatchString(vi) {
					return ReturnFieldError(field, "regex not match")
				}
			}
		}
	}

	checker := strings.Split(field.Tag.Get(TagStringChecker), ",")
	for _, ch := range checker {
		ch = strings.TrimSpace(ch)
		if len(ch) == 0 {
			continue
		}

		fn, ok := c.stringChecker[ch]
		if !ok {
			return ReturnFieldError(field, "checker %s not found", ch)
		}

		err := fn(field, vi)
		if err != nil {
			return ReturnFieldError(field, "checker %s: %s", ch, err.Error())
		}
	}

	return nil
}

func (c *Checker) checkSlice(s any, field *reflect.StructField) errors.WTError {
	if field != nil && field.Tag.Get(TagIgnore) == "true" {
		return nil
	}

	if s == nil {
		return nil
	}

	st := reflect.TypeOf(s)
	if st.Kind() != reflect.Slice {
		return errors.Errorf("not a slice")
	}

	sv := reflect.ValueOf(s)
	if !sv.CanInterface() {
		return errors.Errorf("can not export")
	}

	sliceZero := field.Tag.Get(TagSliceZero)
	sliceIgnore := field.Tag.Get(TagSliceIgnore)

	if sliceZero == "notcheck" && sv.Len() == 0 { // 如果是零值，忽略后面检查
		return nil // 忽略后续检查
	} else if (sliceZero == "ignore" && sv.Len() == 0) || sliceIgnore != "true" {
		// 忽略内置检查
	} else {
		tagMax := field.Tag.Get(TagSliceLengthMax)
		tagMin := field.Tag.Get(TagSliceLengthMin)
		tagLength := field.Tag.Get(TagSliceLength)

		if tagLength != "" {
			sliceLength, err := strconv.ParseInt(tagLength, 10, 64)
			if err != nil {
				return ReturnFieldError(field, err.Error())
			}

			if sv.Len() != int(sliceLength) {
				return ReturnFieldError(field, "bad length")
			}
		}

		if tagMax != "" {
			sliceLengthMax, err := strconv.ParseInt(tagMax, 10, 64)
			if err != nil {
				return ReturnFieldError(field, err.Error())
			}

			if sv.Len() > int(sliceLengthMax) {
				return ReturnFieldError(field, "too long")
			}
		}

		if tagMin != "" {
			sliceLengthMin, err := strconv.ParseInt(tagMin, 10, 64)
			if err != nil {
				return ReturnFieldError(field, err.Error())
			}

			if sv.Len() <= int(sliceLengthMin) {
				return ReturnFieldError(field, "too short")
			}
		}
	}

	checker := strings.Split(field.Tag.Get(TagSliceChecker), ",")
	for _, ch := range checker {
		ch = strings.TrimSpace(ch)
		if len(ch) == 0 {
			continue
		}

		fn, ok := c.sliceChecker[ch]
		if !ok {
			return ReturnFieldError(field, "checker %s not found", ch)
		}

		err := fn(field, s)
		if err != nil {
			return ReturnFieldError(field, "checker %s: %s", ch, err.Error())
		}
	}

	for i := 0; i < sv.Len(); i++ {
		v := sv.Index(i)
		if !v.CanInterface() {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			vi := v.Interface()
			err := c.checkStruct(vi, field)
			if err != nil {
				return err
			}
		case reflect.Slice:
			if c.strictSlice {
				return errors.Errorf("not allow slice in slice")
			}
			// 忽略
		case reflect.Map:
			vi := v.Interface()
			err := c.checkMap(vi, field)
			if err != nil {
				return err
			}
		case reflect.Interface, reflect.Pointer:
			vi := v.Interface()
			err := c.checkInterfaceOrPointer(vi, field)
			if err != nil {
				return err
			}
		case reflect.Bool:
			vi, ok := v.Interface().(bool)
			if !ok {
				return errors.Errorf("not a bool")
			}
			err := c.checkBool(vi, field)
			if err != nil {
				return err
			}
		case reflect.Int64:
			vi, ok := v.Interface().(int64)
			if !ok {
				return errors.Errorf("not a int64")
			}
			err := c.checkInt64(vi, field)
			if err != nil {
				return err
			}
		case reflect.String:
			vi, ok := v.Interface().(string)
			if !ok {
				return errors.Errorf("not a string")
			}
			err := c.checkString(vi, field)
			if err != nil {
				return err
			}
		default:
			return errors.Errorf("unknown type")
		}
	}

	return nil
}

func (c *Checker) checkMap(s any, field *reflect.StructField) errors.WTError {
	if field != nil && field.Tag.Get(TagIgnore) == "true" {
		return nil
	}

	if s == nil {
		return nil
	}

	st := reflect.TypeOf(s)
	if st.Kind() != reflect.Map {
		return errors.Errorf("not a map")
	}

	sv := reflect.ValueOf(s)
	if !sv.CanInterface() {
		return errors.Errorf("can not export")
	}

	mapZero := field.Tag.Get(TagMapZero)
	mapIgnore := field.Tag.Get(TagMapIgnore)

	if mapZero == "notcheck" && sv.Len() == 0 { // 如果是零值，忽略后面检查
		return nil // 忽略后续检查
	} else if (mapZero == "ignore" && sv.Len() == 0) || mapIgnore != "true" {
		// 忽略内置检查
	} else {
		tagMax := field.Tag.Get(TagMapLengthMax)
		tagMin := field.Tag.Get(TagMapLengthMin)
		tagLength := field.Tag.Get(TagMapLength)

		if tagLength != "" {
			mapLength, err := strconv.ParseInt(tagLength, 10, 64)
			if err != nil {
				return ReturnFieldError(field, err.Error())
			}

			if sv.Len() != int(mapLength) {
				return ReturnFieldError(field, "bad length")
			}
		}

		if tagMax != "" {
			mapLengthMax, err := strconv.ParseInt(tagMax, 10, 64)
			if err != nil {
				return ReturnFieldError(field, err.Error())
			}

			if sv.Len() > int(mapLengthMax) {
				return ReturnFieldError(field, "too long")
			}
		}

		if tagMin != "" {
			mapLengthMin, err := strconv.ParseInt(tagMin, 10, 64)
			if err != nil {
				return ReturnFieldError(field, err.Error())
			}

			if sv.Len() <= int(mapLengthMin) {
				return ReturnFieldError(field, "too short")
			}
		}
	}

	checker := strings.Split(field.Tag.Get(TagMapChecker), ",")
	for _, ch := range checker {
		ch = strings.TrimSpace(ch)
		if len(ch) == 0 {
			continue
		}

		fn, ok := c.mapChecker[ch]
		if !ok {
			return ReturnFieldError(field, "checker %s not found", ch)
		}

		err := fn(field, s)
		if err != nil {
			return ReturnFieldError(field, "checker %s: %s", ch, err.Error())
		}
	}

	for _, k := range sv.MapKeys() {
		if k.Kind() != reflect.String {
			return ReturnFieldError(field, "key must be string")
		}

		v := sv.MapIndex(k)
		if !v.CanInterface() {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			vi := v.Interface()
			err := c.checkStruct(vi, field)
			if err != nil {
				return err
			}
		case reflect.Interface, reflect.Pointer:
			vi := v.Interface()
			err := c.checkMap(vi, field)
			if err != nil {
				return err
			}
		case reflect.Map:
			if c.strictMap {
				return errors.Errorf("not allow map in map")
			}
		case reflect.Bool, reflect.Int64, reflect.String, reflect.Slice:
			if c.strictMap {
				return errors.Errorf("only struct can in map")
			}
		default:
			return errors.Errorf("unknown type")
		}
	}

	return nil
}

func (c *Checker) checkStruct(s any, field *reflect.StructField) errors.WTError {
	if field != nil && field.Tag.Get(TagIgnore) == "true" {
		return nil
	}

	if s == nil {
		return nil
	}

	st := reflect.TypeOf(s)
	if st.Kind() != reflect.Struct {
		return errors.Errorf("not a struct")
	}

	sv := reflect.ValueOf(s)
	if !sv.CanInterface() {
		return errors.Errorf("can not export")
	}

	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		v := sv.Field(i)

		if !f.IsExported() || !v.CanInterface() {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			vi := v.Interface()
			err := c.checkStruct(vi, &f)
			if err != nil {
				return err
			}
		case reflect.Slice:
			vi := v.Interface()
			err := c.checkSlice(vi, &f)
			if err != nil {
				return err
			}
		case reflect.Map:
			vi := v.Interface()
			err := c.checkMap(vi, field)
			if err != nil {
				return err
			}
		case reflect.Interface, reflect.Pointer:
			vi := v.Interface()
			err := c.checkInterfaceOrPointer(vi, field)
			if err != nil {
				return err
			}
		case reflect.Bool:
			vi, ok := v.Interface().(bool)
			if !ok {
				return errors.Errorf("not a bool")
			}
			err := c.checkBool(vi, &f)
			if err != nil {
				return err
			}
		case reflect.Int64:
			vi, ok := v.Interface().(int64)
			if !ok {
				return errors.Errorf("not a int64")
			}
			err := c.checkInt64(vi, &f)
			if err != nil {
				return err
			}
		case reflect.String:
			vi, ok := v.Interface().(string)
			if !ok {
				return errors.Errorf("not a string")
			}
			err := c.checkString(vi, &f)
			if err != nil {
				return err
			}
		default:
			return errors.Errorf("unknown type")
		}

	}

	return nil
}

func (c *Checker) checkInterfaceOrPointer(s any, field *reflect.StructField) errors.WTError {
	if field != nil && field.Tag.Get(TagIgnore) == "true" {
		return nil
	}

	if s == nil {
		return nil
	}

	fst := reflect.TypeOf(s)
	if fst.Kind() != reflect.Interface && fst.Kind() != reflect.Pointer {
		return errors.Errorf("not a struct")
	}

	fv := reflect.ValueOf(s)
	if fv.IsNil() || fv.IsZero() {
		return nil
	}

	v := fv.Elem()
	if !v.CanInterface() {
		return errors.Errorf("can not export")
	}

	switch v.Kind() {
	case reflect.Struct:
		vi := v.Interface()
		err := c.checkStruct(vi, field)
		if err != nil {
			return err
		}
	case reflect.Slice:
		vi := v.Interface()
		err := c.checkSlice(vi, field)
		if err != nil {
			return err
		}
	case reflect.Map:
		vi := v.Interface()
		err := c.checkMap(vi, field)
		if err != nil {
			return err
		}
	case reflect.Interface, reflect.Pointer:
		vi := v.Interface()
		err := c.checkInterfaceOrPointer(vi, field)
		if err != nil {
			return err
		}
	case reflect.Bool:
		vi, ok := v.Interface().(bool)
		if !ok {
			return errors.Errorf("not a bool")
		}
		err := c.checkBool(vi, field)
		if err != nil {
			return err
		}
	case reflect.Int64:
		vi, ok := v.Interface().(int64)
		if !ok {
			return errors.Errorf("not a int64")
		}
		err := c.checkInt64(vi, field)
		if err != nil {
			return err
		}
	case reflect.String:
		vi, ok := v.Interface().(string)
		if !ok {
			return errors.Errorf("not a string")
		}
		err := c.checkString(vi, field)
		if err != nil {
			return err
		}
	default:
		return errors.Errorf("unknown type")
	}

	return nil

}
