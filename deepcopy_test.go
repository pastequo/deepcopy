// nolint: exhaustruct
package deepcopy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pastequo/deepcopy"
)

func TestSimpleString(t *testing.T) {
	input := "test"
	var ret string

	err := deepcopy.DeepCopy(input, &ret)

	assert.Nil(t, err)
	assert.Equal(t, input, ret)
}

func TestSimpleInt(t *testing.T) {
	input := 42
	var ret int

	err := deepcopy.DeepCopy(input, &ret)

	assert.Nil(t, err)
	assert.Equal(t, input, ret)
}

func TestSimplePointer(t *testing.T) {
	input := "plouf"
	var output *string

	err := deepcopy.DeepCopy(&input, &output)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "plouf", *output)
	assert.Equal(t, "plouf", input)

	*output = "modified"
	assert.Equal(t, "plouf", input)
}

func TestSimpleNilPointer(t *testing.T) {
	var input1 *string
	var output1 *string

	err := deepcopy.DeepCopy(input1, &output1)

	assert.Nil(t, err)
	assert.Nil(t, input1)
	assert.Nil(t, output1)
}

func TestMultiplePointer(t *testing.T) {
	var input1 ****string
	var output1 ****string

	err := deepcopy.DeepCopy(input1, &output1)

	assert.Nil(t, err)
	assert.Nil(t, input1)
	assert.Nil(t, output1)

	var tmp *string
	tmp1 := &tmp
	tmp2 := &tmp1
	input2 := &tmp2
	var output2 ****string

	err = deepcopy.DeepCopy(input2, &output2)

	assert.Nil(t, err)

	assert.NotNil(t, input2)
	assert.NotNil(t, *input2)
	assert.NotNil(t, **input2)
	assert.Nil(t, ***input2)

	assert.NotNil(t, output2)
	assert.NotNil(t, *output2)
	assert.NotNil(t, **output2)
	assert.Nil(t, ***output2)
}

// func TestStruct(t *testing.T) {
// }

// func TestDeepStruct(t *testing.T) {
// }

// func TestStructWithSamePointedValue(t *testing.T) {
// }

type A struct {
	B *B
}

type B struct {
	A *A
}

type Root struct {
	A *A
}

func TestInsideCycle(t *testing.T) {
	a := A{}
	b := B{A: &a}
	a.B = &b

	input := Root{A: &a}
	output := Root{}

	err := deepcopy.DeepCopy(input, &output)
	assert.Nil(t, err)

	assert.True(t, output.A == output.A.B.A)
}

// nolint: exhaustruct
func TestFirstLevelCycle(t *testing.T) {
	input := A{}
	b := B{A: &input}
	input.B = &b

	output := &A{}

	// deepcopy.DeepCopy(input, output) would fail to detect the very first cycle with input because go is making a copy of input
	err := deepcopy.DeepCopy(&input, &output)
	assert.Nil(t, err)

	assert.True(t, output == output.B.A)
}

func TestNilSlice(t *testing.T) {
	var input []int = nil
	var output []int
	err := deepcopy.DeepCopy(input, &output)
	assert.Nil(t, err)
}

func TestNilInput(t *testing.T) {
	var ret string
	err := deepcopy.DeepCopy(nil, &ret)
	assert.NotNil(t, err)
	assert.Equal(t, deepcopy.ErrNilInput, err)
}

func TestNilDestination(t *testing.T) {
	err := deepcopy.DeepCopy(nil, nil)
	assert.NotNil(t, err)

	err = deepcopy.DeepCopy(1, nil)
	assert.NotNil(t, err)
	assert.Equal(t, deepcopy.ErrNilDestination, err)
}

func TestTypeMismatch(t *testing.T) {
	input1 := "test"
	var ret1 float64

	err := deepcopy.DeepCopy(input1, &ret1)
	assert.NotNil(t, err)
	assert.Equal(t, deepcopy.ErrTypeMismatch, err)

	var input2 ****string
	var ret2 ****uint32

	err = deepcopy.DeepCopy(input2, &ret2)
	assert.NotNil(t, err)
	assert.Equal(t, deepcopy.ErrTypeMismatch, err)
}
