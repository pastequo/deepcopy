package deepcopy

import (
	"fmt"
	"reflect"
	"testing"
)

type A struct {
	B *B
}

type B struct {
	A *A
}

type C struct {
	D1 *D
	D2 *D
}

type D struct {
	Plouf string
}

func TestDeepCopyDouble(t *testing.T) {

	c1 := C{}
	d1 := D{Plouf: "plouf"}
	c1.D1 = &d1
	c1.D2 = &d1

	clone := DeepCopy(c1).(C)

	fmt.Printf("%v (%p): %v -> {%v %v}\n", "c1   ", &c1, c1, *c1.D1, *c1.D2)
	fmt.Printf("%v (%p): %v -> {%v %v}\n", "clone", &clone, clone, clone.D1, clone.D2)

	fmt.Println("--------------------------------")

	c1.D1.Plouf = "plif"
	clone.D1.Plouf = "yolo"

	fmt.Printf("%v (%p): %v -> {%v %v}\n", "c1   ", &c1, c1, *c1.D1, *c1.D2)
	fmt.Printf("%v (%p): %v -> {%v %v}\n", "clone", &clone, clone, *clone.D1, *clone.D2)

	fmt.Println("--------------------------------")

	t.Fatalf("Done")
}

func TestDeepCopySimple(t *testing.T) {

	a1 := A{}
	b1 := B{}
	a1.B = &b1
	b1.A = &a1

	a2 := A{}
	b2 := B{}
	a2.B = &b2
	b2.A = &a2

	fmt.Println("--------------------------------")

	fmt.Printf("%v(%p): %v\n", "a1", &a1, a1)
	fmt.Printf("%v(%p): %v\n", "b1", &b1, b1)
	fmt.Printf("%v(%p): %v\n", "a2", &a2, a2)
	fmt.Printf("%v(%p): %v\n", "b2", &b2, b2)

	fmt.Println("--------------------------------")

	aClone := (DeepCopy(&a1)).(*A)

	fmt.Printf("%v(%p): %v\n", "aClone", aClone, *aClone)
	fmt.Printf("%v(%p): %v\n", "aClone.B", aClone.B, *aClone.B)
	fmt.Printf("%v(%p): %v\n", "aClone.B.A", aClone.B.A, *aClone.B.A)
	fmt.Printf("%v(%p): %v\n", "aClone.B.A.B", aClone.B.A.B, *aClone.B.A.B)

	fmt.Println("--------------------------------")

	t.Fatalf("Done")
}

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
	}
	ref2.Mynested = &nested{NestedString: "inside", NestedArray: [3]bool{true, false, true}, NestedTest: &ref2}
	//*/

	// ref is indeed equals to ref2
	if !reflect.DeepEqual(ref, ref2) {
		dumpTest("ref", &ref, true)
		dumpTest("ref2", &ref2, true)
		t.Fatalf("ref is not equals to ref2")
	}

	yoloRef := yolo{Test: &ref}
	yoloRef2 := yolo{Test: &ref2}
	yoloClone := DeepCopy(&yoloRef).(*yolo)

	fmt.Printf("\n#########################\n\n")

	// First check that ref hadn't been changed.
	if !reflect.DeepEqual(yoloRef, yoloRef2) {
		dumpTest("ref", &ref, true)
		dumpTest("ref2", &ref2, true)
		t.Fatalf("ref had been modified")
	}

	// Then check that clone is equal to ref.
	if !reflect.DeepEqual(yoloRef, *yoloClone) {
		dumpTest("ref", &ref, true)
		dumpTest("clone", yoloClone.Test, true)
		t.Fatalf("clone is not equal to ref")
	}

	// Modify ref and validate that it didn't change clone.
	ref.Myint = 2343
	ref.Mystring = "turbolol"
	*ref.Myint8p = 4
	ref.Myslice = append(ref.Myslice, "d")
	ref.Mymap["swag"] = true
	ref.Mynested.NestedString = "insidemod"

	if !reflect.DeepEqual(yoloRef2, *yoloClone) {
		dumpTest("ref2", &ref2, true)
		dumpTest("clone", yoloClone.Test, true)
		t.Fatalf("clone had been changed")
	}

	dumpTest("ref", &ref, false)
	dumpTest("ref2", &ref2, false)
	dumpTest("clone", yoloClone.Test, false)

	t.Fatalf("yolo,\n%v\n%v\n%v", yoloRef, yoloRef2, *yoloClone)
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
