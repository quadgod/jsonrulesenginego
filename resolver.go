package pathresolver

import (
	"errors"
	"reflect"
	"unsafe"
)

func getValueByPath(pw *pathPathWalker, possiblePtrData *reflect.Value) any {
	data := unrefValue(possiblePtrData)

	switch {
	case !pw.IsCurrentNodeIndex() && data.Kind() == reflect.Struct:
		fVal := data.FieldByName(pw.StringValue())
		unwrappedValue := unrefValue(&fVal)
		if unwrappedValue == nil || !unwrappedValue.IsValid() {
			return nil
		}

		if pw.MoveToNextNode() {
			return getValueByPath(pw, unwrappedValue)
		}

		if unwrappedValue.CanInterface() {
			return unwrappedValue.Interface()
		}

		return reflect.NewAt(unwrappedValue.Type(), unsafe.Pointer(unwrappedValue.UnsafeAddr())).Elem().Interface()
	case !pw.IsCurrentNodeIndex() && data.Kind() == reflect.Map:
		keys := data.MapKeys()
		for _, k := range keys {
			if k.Kind() == reflect.String && k.Interface() == pw.StringValue() {
				fVal := data.MapIndex(k)
				unwrappedValue := unrefValue(&fVal)
				if unwrappedValue == nil || !unwrappedValue.IsValid() {
					return nil
				}

				if pw.MoveToNextNode() {
					return getValueByPath(pw, unwrappedValue)
				}

				if unwrappedValue.CanInterface() {
					return unwrappedValue.Interface()
				}

				return reflect.NewAt(unwrappedValue.Type(), unsafe.Pointer(unwrappedValue.UnsafeAddr())).Elem().Interface()
			}
		}
		return nil
	case pw.IsCurrentNodeIndex() && (data.Kind() == reflect.Array || data.Kind() == reflect.Slice):
		if data.Len() == 0 {
			return nil
		}

		for i := 0; i < data.Len(); i++ {
			if i == pw.IndexValue() {
				rawArrOrSliceItem := data.Index(i)
				arrOrSliceItem := unrefValue(&rawArrOrSliceItem)

				if arrOrSliceItem == nil || !arrOrSliceItem.IsValid() {
					return nil
				}

				if pw.MoveToNextNode() {
					return getValueByPath(pw, arrOrSliceItem)
				}

				if arrOrSliceItem.CanInterface() {
					return arrOrSliceItem.Interface()
				}

				return reflect.NewAt(arrOrSliceItem.Type(), unsafe.Pointer(arrOrSliceItem.UnsafeAddr())).Elem().Interface()
			}
		}

		return nil
	default:
		return nil
	}
}

func TryGetValueByPath(path string, data any) (any, error) {
	if data == nil {
		return nil, errors.New("invalid data arg")
	}

	value := reflect.ValueOf(data)
	if !value.IsValid() {
		return nil, errors.New("invalid data arg")
	}

	for {
		if value.Kind() == reflect.UnsafePointer {
			return nil, errors.New("invalid data arg")
		}

		if value.Kind() == reflect.Pointer {
			value = value.Elem()
			continue
		}

		break
	}

	if !(value.Kind() == reflect.Struct ||
		value.Kind() == reflect.Map ||
		value.Kind() == reflect.Array ||
		value.Kind() == reflect.Slice) {
		return nil, errors.New("data arg must be a struct, map, array or slice")
	}

	pw, err := newPathWalker(path)
	if err != nil {
		return nil, err
	}

	return getValueByPath(pw, &value), nil
}
