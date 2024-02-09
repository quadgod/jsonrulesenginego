package pathresolver

import "reflect"

func unrefValue(arg *reflect.Value) *reflect.Value {
	if arg == nil || !arg.IsValid() {
		return nil
	}

	for {
		if arg.Kind() == reflect.UnsafePointer {
			return nil
		}

		if arg.Kind() == reflect.Pointer {
			elem := arg.Elem()
			return unrefValue(&elem)
		}

		break
	}

	return arg
}
