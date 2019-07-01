package deepcopy

import (
	. "reflect"
	"unsafe"
)

type visit struct {
	addr unsafe.Pointer
	typ  Type
}

// DeepCopy copies recursively the objct given as input.
//
// Cycle:
// They should be preserve. If the cycle starts at the first level of the input, assure that the input is a pointer.
// Otherwise the cycle won't start at the first level.
func DeepCopy(input interface{}) interface{} {
	vObj := ValueOf(input)

	mem := make(map[visit]*Value)
	tmp := deepcopy(vObj, mem)

	return tmp.Interface()
}

func deepcopy(obj Value, done map[visit]*Value) Value {

	if !obj.IsValid() {
		return obj
	}

	var n Value

	switch obj.Kind() {
	case Array, Struct, Bool, Complex64, Complex128, Float32, Float64, Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, String, Interface:
		n = New(obj.Type()).Elem()
	case Map:
		n = MakeMap(obj.Type())
	case Ptr:
		n = New(obj.Elem().Type())
	case Slice:
		n = MakeSlice(obj.Type(), obj.Len(), obj.Cap())
	default:
		panic("unsupported kind")
	}

	if obj.Kind() == Ptr {

		key := visit{
			addr: unsafe.Pointer(obj.Elem().UnsafeAddr()),
			typ:  obj.Type(),
		}

		if pv, ok := done[key]; ok {
			return *pv
		}

		done[key] = &n
	}

	switch obj.Kind() {
	case Array:
		for i := 0; i < obj.Len(); i++ {
			n.Index(i).Set(deepcopy(obj.Index(i), done))
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
		n.Set(obj)
		n.Elem().Set(deepcopy(obj.Elem(), done))
	case Map:
		for _, k := range obj.MapKeys() {
			nkey := deepcopy(k, done)
			nval := deepcopy(obj.MapIndex(k), done)
			n.SetMapIndex(nkey, nval)
		}
	case Ptr:
		if obj.IsNil() {
			return obj
		}

		tmp := deepcopy(obj.Elem(), done)
		n.Elem().Set(tmp)
	case Slice:
		for i := 0; i < obj.Len(); i++ {
			n.Index(i).Set(deepcopy(obj.Index(i), done))
		}
	case Struct:
		for i := 0; i < n.NumField(); i++ {
			n.Field(i).Set(deepcopy(obj.Field(i), done))
		}
	}

	return n
}
