package deepcopy

import (
	"fmt"
	"reflect"
	"testing"
)

type yolo struct {
	Test *test
}

type test struct {
	Myint    int
	Mystring string
	Myint8p  *int8
	My2int8p *int8
	Myslice  []string
	Mymap    map[string]bool
	Mynested *nested
}

type nested struct {
	NestedString string
	NestedArray  [3]bool
	NestedTest   *test
}

func TestDeepCopy(t *testing.T) {

	//*
	a := int8(13)
	ref := test{
		Myint:    42,
		Mystring: "test1",
		Myint8p:  &a,
		My2int8p: &a,
		Myslice:  []string{"a", "b", "c"},
		Mymap:    map[string]bool{"lol": true, "yolo": false},
	}
	ref.Mynested = &nested{NestedString: "inside", NestedArray: [3]bool{true, false, true}, NestedTest: &ref}

	b := int8(13)
	ref2 := test{
		Myint:    42,
		Mystring: "test1",
		Myint8p:  &b,
		My2int8p: &b,
		Myslice:  []string{"a", "b", "c"},
		Mymap:    map[string]bool{"lol": true, "yolo": false},
		Mynested: &nested{NestedString: "inside", NestedArray: [3]bool{true, false, true}},
	}
	ref2.Mynested = &nested{NestedString: "inside", NestedArray: [3]bool{true, false, true}, NestedTest: &ref2}
	//*/

	//ref := "kjo"
	//ref2 := "kjo"

	yoloRef := yolo{Test: &ref}
	yoloRef2 := yolo{Test: &ref2}

	yoloClone := yolo{}
	DeepCopy(&yoloClone, &yoloRef)

	fmt.Printf("\n#########################\n\n")

	// First check that ref hadn't been changed.
	if !reflect.DeepEqual(yoloRef, yoloRef2) {
		dumpTest("ref", &ref, true)
		dumpTest("ref2", &ref2, true)
		t.Fatalf("ref had been modified")
	}

	// Then check that clone is equal to ref.
	if !reflect.DeepEqual(yoloRef, yoloClone) {
		dumpTest("ref", &ref, true)
		dumpTest("clone", yoloRef2.Test, true)
		t.Fatalf("clone is not equal to ref")
	}

	// Modify ref and validate that it didn't change clone.
	//*
	ref.Myint = 2343
	ref.Mystring = "turbolol"
	*ref.Myint8p = 4
	ref.Myslice = append(ref.Myslice, "d")
	ref.Mymap["swag"] = true
	ref.Mynested.NestedString = "insidemod"
	//*/
	//ref = "lol"

	if !reflect.DeepEqual(yoloRef2, yoloClone) {
		dumpTest("ref2", &ref2, true)
		dumpTest("clone", yoloClone.Test, true)
		t.Fatalf("clone had been changed")
	}

	dumpTest("ref", &ref, false)
	dumpTest("ref2", &ref2, false)
	dumpTest("clone", yoloClone.Test, false)

	t.Fatalf("yolo,\n%v\n%v\n%v", yoloRef, yoloRef2, yoloClone)
}

func dumpTest(name string, input *test, recursive bool) {
	fmt.Printf("%v (%p): %v\n", name, input, *input)
	fmt.Println("\tMyint8p:", *input.Myint8p)
	fmt.Println("\tMy2int8p:", *input.My2int8p)
	fmt.Printf("\tMynested (%p): %v\n", input.Mynested, *input.Mynested)
	fmt.Printf("\tMynested.NestedTest (%p): %v\n", input.Mynested.NestedTest, *(input.Mynested.NestedTest))

	fmt.Println()

	if recursive {
		dumpTest("rec_"+name, input.Mynested.NestedTest, false)
	}
}
