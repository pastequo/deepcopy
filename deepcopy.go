package deepcopy

import (
	"errors"
	. "reflect"
	"unsafe"
)

// All errors that can be returned.
var (
	ErrNilInput              = errors.New("input must be non-nil")
	ErrNilDestination        = errors.New("destination must be non-nil")
	ErrDestinationNotPointer = errors.New("destination must be a pointer")
	ErrTypeMismatch          = errors.New("destination is not a pointer of the input type")
	ErrUnexpectedKind        = errors.New("unexpected kind")
)

// DeepCopy copies recursively the object given as input into the destination.
//
//   - Channel are not copied
//   - Cycle should be preserved. If the cycle starts at the first level of the input, assure that the input is a pointer of that cycle.
//     Otherwise the cycle won't start at the first level. (Go passing the parameters by copy)
func DeepCopy(input, destination interface{}) error {
	if input == nil {
		return ErrNilInput
	}

	if destination == nil {
		return ErrNilDestination
	}

	destValue := ValueOf(destination)

	destType := destValue.Type()
	if destType.Kind() != Pointer {
		return ErrDestinationNotPointer
	}

	inputValue := ValueOf(input)
	if !sameType(inputValue, destValue.Elem()) {
		return ErrTypeMismatch
	}

	mem := make(map[visit]*Value)

	tmp, err := deepcopy(inputValue, mem)
	if err != nil {
		return err
	}

	destValue.Elem().Set(tmp)

	return nil
}

func sameType(input, destination Value) bool {
	return input.Type() == destination.Type()
}

type visit struct {
	addr unsafe.Pointer
	typ  Type
}

func deepcopy(obj Value, done map[visit]*Value) (Value, error) {
	if !obj.IsValid() {
		return obj, nil
	}

	var n Value

	switch obj.Kind() {
	case Array, Struct, Bool, Complex64, Complex128, Float32, Float64, Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, String, Interface:
		n = New(obj.Type()).Elem()
	case Map:
		n = MakeMap(obj.Type())
	case Pointer:
		if !obj.Elem().IsValid() {
			return obj, nil
		}

		n = New(obj.Elem().Type())
	case Slice:
		n = MakeSlice(obj.Type(), obj.Len(), obj.Cap())
	default:
		return n, ErrUnexpectedKind
	}

	if obj.Kind() == Ptr {
		key := visit{
			addr: unsafe.Pointer(obj.Elem().UnsafeAddr()),
			typ:  obj.Type(),
		}

		if pv, ok := done[key]; ok {
			return *pv, nil
		}

		done[key] = &n
	}

	switch obj.Kind() {
	case Array:
		for i := 0; i < obj.Len(); i++ {
			elem, err := deepcopy(obj.Index(i), done)
			if err != nil {
				return n, err
			}

			n.Index(i).Set(elem)
		}
	case Bool:
		n.SetBool(obj.Bool())
	case Complex64, Complex128:
		n.SetComplex(obj.Complex())
	case Float32, Float64:
		n.SetFloat(obj.Float())
	case Int, Int8, Int16, Int32, Int64:
		n.SetInt(obj.Int())
	case Uint, Uint8, Uint16, Uint32, Uint64:
		n.SetUint(obj.Uint())
	case String:
		n.SetString(obj.String())
	case Interface:
		elem, err := deepcopy(obj.Elem(), done)
		if err != nil {
			return n, err
		}

		n.Set(obj)
		n.Elem().Set(elem)
	case Map:
		for _, k := range obj.MapKeys() {
			nkey, err := deepcopy(k, done)
			if err != nil {
				return n, err
			}

			nval, err := deepcopy(obj.MapIndex(k), done)
			if err != nil {
				return n, err
			}

			n.SetMapIndex(nkey, nval)
		}
	case Ptr:
		if obj.IsNil() {
			return obj, nil
		}

		elem, err := deepcopy(obj.Elem(), done)
		if err != nil {
			return n, err
		}

		n.Elem().Set(elem)
	case Slice:
		for i := 0; i < obj.Len(); i++ {
			elem, err := deepcopy(obj.Index(i), done)
			if err != nil {
				return n, err
			}

			n.Index(i).Set(elem)
		}
	case Struct:
		for i := 0; i < n.NumField(); i++ {
			elem, err := deepcopy(obj.Field(i), done)
			if err != nil {
				return n, err
			}

			n.Field(i).Set(elem)
		}
	}

	return n, nil
}
