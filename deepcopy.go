package deepcopy

import (
	. "reflect"
)

// DeepCopy copies the object provided as input.
func DeepCopy(obj interface{}) interface{} {

	return deepcopy(ValueOf(obj)).Interface()
}

func deepcopy(obj Value) Value {

	if !obj.IsValid() {
		return obj
	}

	n := Zero(obj.Type())

	switch obj.Kind() {
	case Array:
		n = Zero(ArrayOf(obj.Len(), obj.Type()))
		for i := 0; i < obj.Len(); i++ {
			n.Index(i).Set(deepcopy(obj.Index(i)))
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
		var tmp uint8
		n2 := ValueOf(&tmp)
		n = n2.Elem()
		n.SetUint(obj.Uint())
	case String:
		var tmp string
		n2 := ValueOf(&tmp)
		n = n2.Elem()
		n.SetString(obj.String())
	case Interface:
		n.Set(obj)
		n.Elem().Set(deepcopy(obj.Elem()))
	case Map:
		n = MakeMap(obj.Type())
		for _, k := range obj.MapKeys() {
			nkey := deepcopy(k)
			nval := deepcopy(obj.MapIndex(k))
			n.SetMapIndex(nkey, nval)
		}
	case Ptr:
		if obj.IsNil() {
			return obj
		}
		n = New(obj.Elem().Type())
		tmp := deepcopy(obj.Elem())
		n.Elem().Set(tmp)
	case Slice:
		n = MakeSlice(obj.Type(), obj.Len(), obj.Cap())
		for i := 0; i < obj.Len(); i++ {
			n.Index(i).Set(deepcopy(obj.Index(i)))
		}
	case Struct:
		/*
			        var tmp XXXX
					n = ValueOf(&tmp)
					for i := 0; i < obj.NumField(); i++ {
						n.Elem().Field(i).Set(deepcopy(obj.Field(i)))
					}
					n = n.Elem()
		*/
	}
	return n
}
