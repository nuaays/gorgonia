package tensori

import (
	"testing"

	"github.com/chewxy/gorgonia/tensor/types"
	"github.com/stretchr/testify/assert"
)

func TestAt(t *testing.T) {
	backing := RangeInt(0, 6)
	T := NewTensor(WithShape(2, 3), WithBacking(backing))
	zeroone := T.At(0, 1)
	assert.Equal(t, int(1), zeroone)

	oneone := T.At(1, 1)
	assert.Equal(t, int(4), oneone)

	fail := func() {
		T.At(1, 2, 3)
	}
	assert.Panics(t, fail, "Expected too many coordinates to panic")

	backing = RangeInt(0, 24)
	T = NewTensor(WithShape(2, 3, 4), WithBacking(backing))
	/*
		T = [0, 1, 2, 3]
			[4, 5, 6, 7]
			[8, 9, 10, 11]

			[12, 13, 14, 15]
			[16, 17, 18, 19]
			[20, 21, 22, 23]
	*/
	oneoneone := T.At(1, 1, 1)
	assert.Equal(t, int(17), oneoneone)
	zthreetwo := T.At(0, 2, 2)
	assert.Equal(t, int(10), zthreetwo)
	onetwothree := T.At(1, 2, 3)
	assert.Equal(t, int(23), onetwothree)

	fail = func() {
		T.At(0, 3, 2)
	}
	assert.Panics(t, fail)
}

func TestT_transposeIndex(t *testing.T) {
	assert := assert.New(t)
	var T *Tensor

	T = NewTensor(WithShape(2, 2), WithBacking(RangeInt(0, 4)))

	correct := []int{0, 2, 1, 3}
	for i, v := range correct {
		assert.Equal(v, T.transposeIndex(i, []int{1, 0}, []int{2, 1}))
	}
}

var transposeTests = []struct {
	name          string
	shape         types.Shape
	transposeWith []int

	correctShape    types.Shape
	correctStrides  []int // after .T()
	correctStrides2 []int // after .Transpose()
	correctData     []int // after .Transpose()
}{
	{"c.T()", types.Shape{4, 1}, nil, types.Shape{1, 4}, []int{1}, []int{1}, RangeInt(0, 4)},
	{"r.T()", types.Shape{1, 4}, nil, types.Shape{4, 1}, []int{1}, []int{1}, RangeInt(0, 4)},
	{"v.T()", types.Shape{4}, nil, types.Shape{4}, []int{1}, []int{1}, RangeInt(0, 4)},
	{"M.T()", types.Shape{2, 3}, nil, types.Shape{3, 2}, []int{1, 3}, []int{2, 1}, []int{0, 3, 1, 4, 2, 5}},
	{"M.T(0,1) (NOOP)", types.Shape{2, 3}, []int{0, 1}, types.Shape{2, 3}, []int{3, 1}, []int{3, 1}, RangeInt(0, 6)},
	{"3T.T()", types.Shape{2, 3, 4}, nil, types.Shape{4, 3, 2}, []int{1, 4, 12}, []int{6, 2, 1}, []int{0, 12, 4, 16, 8, 20, 1, 13, 5, 17, 9, 21, 2, 14, 6, 18, 10, 22, 3, 15, 7, 19, 11, 23}},
	{"3T.T(2, 1, 0) (Same as .T())", types.Shape{2, 3, 4}, []int{2, 1, 0}, types.Shape{4, 3, 2}, []int{1, 4, 12}, []int{6, 2, 1}, []int{0, 12, 4, 16, 8, 20, 1, 13, 5, 17, 9, 21, 2, 14, 6, 18, 10, 22, 3, 15, 7, 19, 11, 23}},
	{"3T.T(0, 2, 1)", types.Shape{2, 3, 4}, []int{0, 2, 1}, types.Shape{2, 4, 3}, []int{12, 1, 4}, []int{12, 3, 1}, []int{0, 4, 8, 1, 5, 9, 2, 6, 10, 3, 7, 11, 12, 16, 20, 13, 17, 21, 14, 18, 22, 15, 19, 23}},
	{"3T.T{1, 0, 2)", types.Shape{2, 3, 4}, []int{1, 0, 2}, types.Shape{3, 2, 4}, []int{4, 12, 1}, []int{8, 4, 1}, []int{0, 1, 2, 3, 12, 13, 14, 15, 4, 5, 6, 7, 16, 17, 18, 19, 8, 9, 10, 11, 20, 21, 22, 23}},
	{"3T.T{1, 2, 0)", types.Shape{2, 3, 4}, []int{1, 2, 0}, types.Shape{3, 4, 2}, []int{4, 1, 12}, []int{8, 2, 1}, []int{0, 12, 1, 13, 2, 14, 3, 15, 4, 16, 5, 17, 6, 18, 7, 19, 8, 20, 9, 21, 10, 22, 11, 23}},
	{"3T.T{2, 0, 1)", types.Shape{2, 3, 4}, []int{2, 0, 1}, types.Shape{4, 2, 3}, []int{1, 12, 4}, []int{6, 3, 1}, []int{0, 4, 8, 12, 16, 20, 1, 5, 9, 13, 17, 21, 2, 6, 10, 14, 18, 22, 3, 7, 11, 15, 19, 23}},
	{"3T.T{0, 1, 2} (NOOP)", types.Shape{2, 3, 4}, []int{0, 1, 2}, types.Shape{2, 3, 4}, []int{12, 4, 1}, []int{12, 4, 1}, RangeInt(0, 24)},
}

func TestTranspose(t *testing.T) {
	assert := assert.New(t)
	var T *Tensor
	var err error

	// standard transposes
	for _, tts := range transposeTests {
		T = NewTensor(WithShape(tts.shape...), WithBacking(RangeInt(0, tts.shape.TotalSize())))
		if err = T.T(tts.transposeWith...); err != nil {
			t.Errorf("%v - %v", tts.name, err)
			continue
		}

		assert.True(tts.correctShape.Eq(T.Shape()), "Transpose %v Expected shape: %v. Got %v", tts.name, tts.correctShape, T.Shape())
		assert.Equal(tts.correctStrides, T.Strides())
		T.Transpose()
		assert.True(tts.correctShape.Eq(T.Shape()), "Transpose %v Expected shape: %v. Got %v", tts.name, tts.correctShape, T.Shape())
		assert.Equal(tts.correctStrides2, T.Strides())
		assert.Equal(tts.correctData, T.data)
	}

	// test stacked .T() calls

	// column vector
	T = NewTensor(WithShape(4, 1), WithBacking(RangeInt(0, 4)))
	if err = T.T(); err != nil {
		t.Errorf("Stacked .T() #1 for vector. Error: %v", err)
		goto matrev
	}
	if err = T.T(); err != nil {
		t.Errorf("Stacked .T() #1 for vector. Error: %v", err)
		goto matrev
	}
	assert.Nil(T.old)
	assert.Nil(T.transposeWith)
	assert.True(T.IsColVec())

matrev:
	// matrix, reversed
	T = NewTensor(WithShape(2, 3), WithBacking(RangeInt(0, 6)))
	if err = T.T(); err != nil {
		t.Errorf("Stacked .T() #1 for matrix reverse. Error: %v", err)
		goto matnorev
	}
	if err = T.T(); err != nil {
		t.Errorf("Stacked .T() #2 for matrix reverse. Error: %v", err)
		goto matnorev
	}
	assert.Nil(T.old)
	assert.Nil(T.transposeWith)
	assert.True(types.Shape{2, 3}.Eq(T.Shape()))

matnorev:
	// 3-tensor, non reversed
	T = NewTensor(WithShape(2, 3, 4), WithBacking(RangeInt(0, 24)))
	if err = T.T(); err != nil {
		t.Fatalf("Stacked .T() #1 for tensor with no reverse. Error: %v", err)
	}
	if err = T.T(2, 0, 1); err != nil {
		t.Fatalf("Stacked .T() #2 for tensor with no reverse. Error: %v", err)
	}
	correctData := []int{0, 12, 4, 16, 8, 20, 1, 13, 5, 17, 9, 21, 2, 14, 6, 18, 10, 22, 3, 15, 7, 19, 11, 23}
	assert.Equal(correctData, T.data)
	assert.Equal([]int{2, 0, 1}, T.transposeWith)
	assert.NotNil(T.old)

}

func TestTUT(t *testing.T) {
	assert := assert.New(t)
	var T *Tensor

	T = NewTensor(WithShape(2, 3, 4))
	T.T()
	T.UT()
	assert.Nil(T.old)
	assert.Nil(T.transposeWith)

	T.T(2, 0, 1)
	T.UT()
	assert.Nil(T.old)
	assert.Nil(T.transposeWith)
}

func TestTRepeat(t *testing.T) {
	assert := assert.New(t)
	var T, T2 *Tensor
	var expectedShape types.Shape
	var expectedData []int
	var err error

	// SCALARS

	T = NewTensor(AsScalar(int(3)))
	T2, err = T.Repeat(0, 3)
	if err != nil {
		t.Error(err)
	}

	if T == T2 {
		t.Error("Not supposed to be the same pointer")
	}
	expectedShape = types.Shape{3}
	expectedData = []int{3, 3, 3}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	T2, err = T.Repeat(1, 3)
	if err != nil {
		t.Error(err)
	}

	expectedShape = types.Shape{1, 3}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	// VECTORS

	// These are the rules for vector repeats:
	// 	- Vectors can repeat on axis 0 and 1
	// 	- For vanilla vectors, repeating on axis 0 and 1 is as if it were a colvec
	// 	- For non vanilla vectors, it's as if it were a matrix being repeated

	var backing = []int{1, 2}

	// repeats on axis 1: colvec
	T = NewTensor(WithShape(2, 1), WithBacking(backing))
	T2, err = T.Repeat(1, 3)
	if err != nil {
		t.Error(err)
	}

	expectedShape = types.Shape{2, 3}
	expectedData = []int{1, 1, 1, 2, 2, 2}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	// repeats on axis 1: vanilla vector
	T = NewTensor(WithShape(2), WithBacking(backing))
	T2, err = T.Repeat(1, 3)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	// repeats on axis 1: rowvec
	T = NewTensor(WithShape(1, 2), WithBacking(backing))
	T2, err = T.Repeat(1, 3)
	if err != nil {
		t.Error(err)
	}
	expectedShape = types.Shape{1, 6}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	// repeats on axis 0: vanilla vectors
	T = NewTensor(WithShape(2), WithBacking(backing))
	T2, err = T.Repeat(0, 3)
	if err != nil {
		t.Error(err)
	}
	expectedShape = types.Shape{6}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	// repeats on axis 0: colvec
	T = NewTensor(WithShape(2, 1), WithBacking(backing))
	T2, err = T.Repeat(0, 3)
	if err != nil {
		t.Error(err)
	}
	expectedShape = types.Shape{6, 1}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	// repeats on axis 0: rowvec
	T = NewTensor(WithShape(1, 2), WithBacking(backing))
	T2, err = T.Repeat(0, 3)
	if err != nil {
		t.Error(err)
	}
	expectedData = []int{1, 2, 1, 2, 1, 2}
	expectedShape = types.Shape{3, 2}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	// repeats on -1 : all should have shape of (6)
	T = NewTensor(WithShape(2, 1), WithBacking(backing))
	T2, err = T.Repeat(-1, 3)
	if err != nil {
		t.Error(err)
	}
	expectedData = []int{1, 1, 1, 2, 2, 2}
	expectedShape = types.Shape{6}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	T = NewTensor(WithShape(1, 2), WithBacking(backing))
	T2, err = T.Repeat(-1, 3)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	T = NewTensor(WithShape(2), WithBacking(backing))
	T2, err = T.Repeat(-1, 3)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	// MATRICES

	backing = []int{1, 2, 3, 4}

	/*
		1, 2,
		3, 4
	*/

	T = NewTensor(WithShape(2, 2), WithBacking(backing))
	T2, err = T.Repeat(-1, 1, 2, 1, 1)
	if err != nil {
		t.Error(err)
	}

	expectedShape = types.Shape{5}
	expectedData = []int{1, 2, 2, 3, 4}

	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	/*
		1, 1, 2
		3, 3, 4
	*/
	T2, err = T.Repeat(1, 2, 1)
	if err != nil {
		t.Error(err)
	}
	expectedShape = types.Shape{2, 3}
	expectedData = []int{1, 1, 2, 3, 3, 4}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	/*
		1, 2, 2,
		3, 4, 4
	*/
	T2, err = T.Repeat(1, 1, 2)
	if err != nil {
		t.Error(err)
	}
	expectedShape = types.Shape{2, 3}
	expectedData = []int{1, 2, 2, 3, 4, 4}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	/*
		1, 2,
		3, 4,
		3, 4
	*/
	T2, err = T.Repeat(0, 1, 2)
	if err != nil {
		t.Error(err)
	}
	expectedShape = types.Shape{3, 2}
	expectedData = []int{1, 2, 3, 4, 3, 4}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	/*
		1, 2,
		1, 2,
		3, 4
	*/
	T2, err = T.Repeat(0, 2, 1)
	if err != nil {
		t.Error(err)
	}
	expectedShape = types.Shape{3, 2}
	expectedData = []int{1, 2, 1, 2, 3, 4}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	// MORE THAN 2D!!
	/*
		In:
			1, 2,
			3, 4,
			5, 6,

			7, 8,
			9, 10,
			11, 12
		Out:
			1, 2,
			3, 4
			3, 4
			5, 6

			7, 8,
			9, 10,
			9, 10,
			11, 12
	*/
	T = NewTensor(WithShape(2, 3, 2), WithBacking(RangeInt(1, 2*3*2+1)))
	T2, err = T.Repeat(1, 1, 2, 1)
	if err != nil {
		t.Error(err)
	}
	expectedShape = types.Shape{2, 4, 2}
	expectedData = []int{1, 2, 3, 4, 3, 4, 5, 6, 7, 8, 9, 10, 9, 10, 11, 12}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	// broadcast errors
	T2, err = T.Repeat(0, 1, 2, 1)
	if err == nil {
		t.Error("Expected a broadacast/shapeMismatch error")
	}

	// generic repeat - repeat EVERYTHING by 2
	T2, err = T.Repeat(types.AllAxes, 2)
	if err != nil {
		t.Error(err)
	}
	expectedShape = types.Shape{24}
	expectedData = []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 11, 12, 12}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	// generic repeat, axis specified
	T2, err = T.Repeat(2, 2)
	if err != nil {
		t.Error(err)
	}
	expectedShape = types.Shape{2, 3, 4}
	expectedData = []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 11, 12, 12}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	// repeat scalars!
	T = NewTensor(AsScalar(int(3)))
	T2, err = T.Repeat(0, 5)
	if err != nil {
		t.Error(err)
	}
	expectedData = []int{3, 3, 3, 3, 3}
	expectedShape = types.Shape{5}
	assert.Equal(expectedData, T2.data)
	assert.Equal(expectedShape, T2.Shape())

	/* IDIOTS SECTION */

	// trying to repeat on a nonexistant axis - Vector
	T = NewTensor(WithShape(2, 1), WithBacking([]int{1, 2}))
	fails := func() {
		T.Repeat(2, 3)
	}
	assert.Panics(fails)

	T = NewTensor(WithShape(2, 3), WithBacking([]int{1, 2, 3, 4, 5, 6}))
	fails = func() {
		T.Repeat(3, 3)
	}
	assert.Panics(fails)
}

var sliceTests = []struct {
	name string

	shape  types.Shape
	slices []types.Slice

	correctShape  types.Shape
	correctStride []int
	correctData   []int
}{
	{"a[0]", types.Shape{5}, []types.Slice{ss(0)}, types.ScalarShape(), nil, []int{0}},
	{"a[0:2]", types.Shape{5}, []types.Slice{makeRS(0, 2)}, types.Shape{2}, []int{1}, []int{0, 1}},
	{"a[1:5:2]", types.Shape{5}, []types.Slice{makeRS(0, 5, 2)}, types.Shape{2}, []int{2}, []int{0, 1, 2, 3, 4}},

	// colvec
	{"c[0]", types.Shape{5, 1}, []types.Slice{ss(0)}, types.ScalarShape(), nil, []int{0}},
	{"c[0:2]", types.Shape{5, 1}, []types.Slice{makeRS(0, 2)}, types.Shape{2, 1}, []int{1}, []int{0, 1}},
	{"c[1:5:2]", types.Shape{5, 1}, []types.Slice{makeRS(0, 5, 2)}, types.Shape{2, 1}, []int{2}, []int{0, 1, 2, 3, 4}},

	// rowvec
	{"r[0]", types.Shape{1, 5}, []types.Slice{ss(0)}, types.Shape{1, 5}, []int{1}, []int{0, 1, 2, 3, 4}},
	{"r[0:2]", types.Shape{1, 5}, []types.Slice{makeRS(0, 2)}, types.Shape{1, 5}, []int{1}, []int{0, 1, 2, 3, 4}},
	{"r[0:5:2]", types.Shape{1, 5}, []types.Slice{makeRS(0, 5, 2)}, types.Shape{1, 5}, []int{1}, []int{0, 1, 2, 3, 4}},
	{"r[:, 0]", types.Shape{1, 5}, []types.Slice{nil, ss(0)}, types.ScalarShape(), nil, []int{0}},
	{"r[:, 0:2]", types.Shape{1, 5}, []types.Slice{nil, makeRS(0, 2)}, types.Shape{1, 2}, []int{1}, []int{0, 1}},
	{"r[:, 1:5:2]", types.Shape{1, 5}, []types.Slice{nil, makeRS(1, 5, 2)}, types.Shape{1, 2}, []int{2}, []int{1, 2, 3, 4}},

	// matrix
	{"A[0]", types.Shape{2, 3}, []types.Slice{ss(0)}, types.Shape{1, 3}, []int{1}, RangeInt(0, 3)},
	{"A[0:10]", types.Shape{4, 5}, []types.Slice{makeRS(0, 2)}, types.Shape{2, 5}, []int{5, 1}, RangeInt(0, 10)},
	{"A[0, 0]", types.Shape{4, 5}, []types.Slice{ss(0), ss(0)}, types.ScalarShape(), nil, []int{0}},
	{"A[0, 1:5]", types.Shape{4, 5}, []types.Slice{ss(0), makeRS(1, 5)}, types.Shape{4}, []int{1}, RangeInt(1, 5)},
	{"A[0, 1:5:2]", types.Shape{4, 5}, []types.Slice{ss(0), makeRS(1, 5, 2)}, types.Shape{1, 2}, []int{2}, RangeInt(1, 5)},
	{"A[:, 0]", types.Shape{4, 5}, []types.Slice{nil, ss(0)}, types.Shape{4, 1}, []int{5}, RangeInt(0, 16)},
	{"A[:, 1:5]", types.Shape{4, 5}, []types.Slice{nil, makeRS(1, 5)}, types.Shape{4, 4}, []int{5, 1}, RangeInt(1, 20)},
	{"A[:, 1:5:2]", types.Shape{4, 5}, []types.Slice{nil, makeRS(1, 5, 2)}, types.Shape{4, 2}, []int{5, 2}, RangeInt(1, 20)},
}

func TestTSlice(t *testing.T) {
	assert := assert.New(t)
	var T, V *Tensor
	var err error

	for _, sts := range sliceTests {
		T = NewTensor(WithShape(sts.shape...), WithBacking(RangeInt(0, sts.shape.TotalSize())))
		t.Log(sts.name)
		if V, err = T.Slice(sts.slices...); err != nil {
			t.Error(err)
			continue
		}
		assert.True(sts.correctShape.Eq(V.Shape()), "Test: %v - Incorrect Shape. Correct: %v. Got %v", sts.name, sts.correctShape, V.Shape())
		assert.Equal(sts.correctStride, V.Strides(), "Test: %v - Incorrect Stride", sts.name)
		assert.Equal(sts.correctData, V.data, "Test: %v - Incorrect Data", sts.name)
	}

	// And now, ladies and gentlemen, the idiots!

	// too many slices
	_, err = T.Slice(ss(1), ss(2), ss(3), ss(4))
	if err == nil {
		t.Error("Expected a DimMismatchError error")
	}

	// out of range sliced
	_, err = T.Slice(makeRS(20, 5))
	if err == nil {
		t.Error("Expected a IndexError")
	}

	// surely nobody can be this dumb? Having a start of negatives
	_, err = T.Slice(makeRS(-1, 1))
	if err == nil {
		t.Error("Expected a IndexError")
	}

}

func TestT_at_itol(t *testing.T) {
	assert := assert.New(t)
	var err error
	var T *Tensor
	var shape types.Shape

	T = NewTensor(WithBacking(RangeInt(0, 12)), WithShape(3, 4))
	t.Logf("%+v", T)

	shape = T.Shape()
	for i := 0; i < shape[0]; i++ {
		for j := 0; j < shape[1]; j++ {
			coord := []int{i, j}
			idx, err := T.at(coord...)
			if err != nil {
				t.Error(err)
			}

			got, err := T.itol(idx)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(coord, got)
		}
	}

	T = NewTensor(WithBacking(RangeInt(0, 24)), WithShape(2, 3, 4))

	shape = T.Shape()
	for i := 0; i < shape[0]; i++ {
		for j := 0; j < shape[1]; j++ {
			for k := 0; k < shape[2]; k++ {
				coord := []int{i, j, k}
				idx, err := T.at(coord...)
				if err != nil {
					t.Error(err)
				}

				got, err := T.itol(idx)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(coord, got)
			}
		}
	}

	/* Transposes */

	T = NewTensor(WithBacking(RangeInt(0, 6)), WithShape(2, 3))
	t.Logf("%+v", T)
	err = T.T()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v, %v", T.Shape(), T.Shape())
	t.Logf("%v, %v", T.Strides(), T.ostrides())

	shape = T.Shape()
	for i := 0; i < shape[0]; i++ {
		for j := 0; j < shape[1]; j++ {
			coord := []int{i, j}
			idx, err := T.at(coord...)
			if err != nil {
				t.Error(err)
				continue
			}

			got, err := T.itol(idx)
			if err != nil {
				t.Error(err)
				continue
			}

			assert.Equal(coord, got)
		}
	}

	/* IDIOT OF THE WEEK */

	T = NewTensor(WithBacking(RangeInt(0, 24)), WithShape(2, 3, 4))

	_, err = T.at(1, 3, 2) // the 3 is out of range
	if err == nil {
		t.Error("Expected an error")
	}
	t.Log(err)

	_, err = T.itol(24) // 24 is out of range
	if err == nil {
		t.Error("Expected an error")
	}
}

func TestCopyTo(t *testing.T) {
	assert := assert.New(t)
	var T, T2, T3 *Tensor
	var err error

	T = NewTensor(WithShape(2), WithBacking([]int{1, 2}))
	T2 = NewTensor(WithShape(1, 2))

	err = T.CopyTo(T2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(T2.data, T.data)

	// now, modify T1's data
	T.data[0] = 5000
	assert.NotEqual(T2.data, T.data)

	// test views
	T = NewTensor(WithShape(3, 3))
	T2 = NewTensor(WithShape(2, 2))
	T3, _ = T.Slice(makeRS(0, 2), makeRS(0, 2)) // T[0:2, 0:2], shape == (2,2)
	if err = T2.CopyTo(T3); err != nil {
		t.Log(err) // for now it's a not yet implemented error. TODO: FIX THIS
	}

	// dumbass time

	T = NewTensor(WithShape(3, 3))
	T2 = NewTensor(WithShape(2, 2))
	if err = T.CopyTo(T2); err == nil {
		t.Error("Expected an error")
	}

	if err = T.CopyTo(T); err != nil {
		t.Error("Copying a *Tensor to itself should yield no error. ")
	}

}
