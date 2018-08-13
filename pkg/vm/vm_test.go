package vm

import (
	"math"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildVM(t *testing.T, s string) *ChaincodeVM {
	ops := miniAsm(s)
	bin := ChasmBinary{"test", "TEST", ops}
	vm, err := New(bin)
	assert.Nil(t, err)
	return vm
}

func buildVMfail(t *testing.T, s string) {
	ops := miniAsm(s)
	bin := ChasmBinary{"test", "TEST", ops}
	_, err := New(bin)
	assert.NotNil(t, err)
}

func checkStack(t *testing.T, st *Stack, values ...int64) {
	for i := range values {
		n, err := st.PopAsInt64()
		assert.Nil(t, err)
		assert.Equal(t, values[len(values)-i-1], n)
	}
}

func TestMiniAsmBasics(t *testing.T) {
	ops := miniAsm("neg1 zero one push1 45 push2 01 01 2000-01-01T00:00:00Z")
	bytes := []Opcode{OpNeg1, OpZero, OpOne, OpPush1, 69, OpPush2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0}
	assert.Equal(t, ops, bytes)
}

func TestNop(t *testing.T) {
	vm := buildVM(t, "handler 0 nop enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	assert.Equal(t, vm.Stack().Depth(), 0)
}

func TestPush(t *testing.T) {
	vm := buildVM(t, "handler 0 neg1 zero one push1 45 push2 01 02 ret enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -1, 0, 1, 69, 513)
}

func TestBigPush(t *testing.T) {
	vm := buildVM(t, `handler 0
		push3 1 2 3
		push4 4 0 0 1
		push5 5 0 0 0 1
		push6 6 0 0 0 0 1
		push7 1 2 3 4 5 6 7
		push8 fb ff ff ff ff ff ff ff enddef`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 197121, 16777220, 4294967301, 1099511627782, 1976943448883713, -5)
}

func TestPushB1(t *testing.T) {
	vm := buildVM(t, "handler 0 pushb 09 41 42 43 44 45 46 47 48 49 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	v, err := vm.Stack().Pop()
	assert.Nil(t, err)
	assert.IsType(t, NewBytes(nil), v)
	assert.Equal(t, NewBytes([]byte{65, 66, 67, 68, 69, 70, 71, 72, 73}), v)
}

func TestPushB2(t *testing.T) {
	vm := buildVM(t, `handler 0 pushb "ABCDEFGHI" enddef`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	v, err := vm.Stack().Pop()
	assert.Nil(t, err)
	assert.IsType(t, NewBytes(nil), v)
	assert.Equal(t, NewBytes([]byte{65, 66, 67, 68, 69, 70, 71, 72, 73}), v)
}

func TestPushA(t *testing.T) {
	vm := buildVM(t, `handler 0 pusha ndadprx764ciigti8d8whtw2kct733r85qvjukhqhke3dka4 enddef`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	v, err := vm.Stack().Pop()
	assert.Nil(t, err)
	assert.IsType(t, NewBytes(nil), v)
	assert.Equal(t, NewBytes([]byte{
		0x6e, 0x64, 0x61, 0x64, 0x70, 0x72, 0x78, 0x37,
		0x36, 0x34, 0x63, 0x69, 0x69, 0x67, 0x74, 0x69,
		0x38, 0x64, 0x38, 0x77, 0x68, 0x74, 0x77, 0x32,
		0x6b, 0x63, 0x74, 0x37, 0x33, 0x33, 0x72, 0x38,
		0x35, 0x71, 0x76, 0x6a, 0x75, 0x6b, 0x68, 0x71,
		0x68, 0x6b, 0x65, 0x33, 0x64, 0x6b, 0x61, 0x34}), v)
}

func TestDrop(t *testing.T) {
	vm := buildVM(t, "handler 0 push1 7 nop one zero neg1 drop drop2 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 7)
}

func TestDup(t *testing.T) {
	vm := buildVM(t, "handler 0 one push1 2 dup push1 3 dup2 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 1, 2, 2, 3, 2, 3)
}

func TestSwapOverPickRoll(t *testing.T) {
	vm := buildVM(t, "handler 0 zero one push1 2 push1 3 swap over pick 4 roll 4 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 3, 2, 3, 0, 1)
}

func TestPickRollEdgeCases(t *testing.T) {
	vm := buildVM(t, "handler 0 zero one pick 0 push1 2 roll 0 push1 3 roll 1 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 1, 1, 3, 2)
}

func TestTuck(t *testing.T) {
	vm := buildVM(t, "handler 0 zero one push1 2 push1 3 tuck 0 tuck 1 tuck 1 tuck 3 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3, 0, 1, 2)
}

func TestTuckFail(t *testing.T) {
	vm := buildVM(t, "handler 0 zero one push1 2 push1 3 tuck 4 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestMath(t *testing.T) {
	vm := buildVM(t, "handler 0 push1 55 dup dup add sub push1 7 push1 6 mul dup push1 3 div dup push1 3 mod enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -85, 42, 14, 2)
}

func TestMul(t *testing.T) {
	type a struct {
		in  []int64
		out int64
	}
	args := []a{
		a{[]int64{5, 3}, 15},
		a{[]int64{5, 5}, 25},
		a{[]int64{3, 5}, 15},
		a{[]int64{12, 4}, 48},
		a{[]int64{5, -3}, -15},
		a{[]int64{5, 0}, 0},
		a{[]int64{0, 5}, 0},
		a{[]int64{-12, -4}, 48},
	}
	vm := buildVM(t, "handler 0 mul enddef")

	for a := range args {
		vm.Init(0, NewNumber(args[a].in[0]), NewNumber(args[a].in[1]))
		err := vm.Run(false)
		assert.Nil(t, err)
		checkStack(t, vm.Stack(), args[a].out)
	}
}

func TestDiv(t *testing.T) {
	type a struct {
		in  []int64
		out int64
	}
	args := []a{
		a{[]int64{5, 3}, 1},
		a{[]int64{5, 5}, 1},
		a{[]int64{3, 5}, 0},
		a{[]int64{12, 4}, 3},
		a{[]int64{5, -3}, -1},
		a{[]int64{50, 5}, 10},
		a{[]int64{0, 5}, 0},
		a{[]int64{-12, -4}, 3},
	}
	vm := buildVM(t, "handler 0 div enddef")

	for a := range args {
		vm.Init(0, NewNumber(args[a].in[0]), NewNumber(args[a].in[1]))
		err := vm.Run(false)
		assert.Nil(t, err)
		checkStack(t, vm.Stack(), args[a].out)
	}
}

func TestMod(t *testing.T) {
	type a struct {
		in  []int64
		out int64
	}
	args := []a{
		a{[]int64{5, 3}, 2},
		a{[]int64{5, 5}, 0},
		a{[]int64{3, 5}, 3},
		a{[]int64{12, 4}, 0},
	}
	vm := buildVM(t, "handler 0 mod enddef")
	for a := range args {
		vm.Init(0, NewNumber(args[a].in[0]), NewNumber(args[a].in[1]))
		err := vm.Run(false)
		assert.Nil(t, err)
		checkStack(t, vm.Stack(), args[a].out)
	}
}

func TestDivMod(t *testing.T) {
	vm := buildVM(t, "handler 0 push1 17 push1 7 divmod enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 2, 3)
}

func TestMathErrors(t *testing.T) {
	vm := buildVM(t, "handler 0 push1 55 zero div enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.NotNil(t, err, "divide by zero")

	vm = buildVM(t, "handler 0 push1 55 zero mod enddef")
	vm.Init(0)
	err = vm.Run(false)
	assert.NotNil(t, err, "mod by zero")

	vm = buildVM(t, "handler 0 push1 55 zero divmod enddef")
	vm.Init(0)
	err = vm.Run(false)
	assert.NotNil(t, err, "divmod by zero")

	vm = buildVM(t, "handler 0 push1 55 push1 2 zero muldiv enddef")
	vm.Init(0)
	err = vm.Run(false)
	assert.NotNil(t, err, "muldiv by zero")
}

func TestMathOverflows(t *testing.T) {
	vm := buildVM(t, "handler 0 push8 7a bb cc dd ee ff 99 88 push1 ff mul enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.NotNil(t, err, "mul overflow")
	vm = buildVM(t, "handler 0 push8 7f bb cc dd ee ff 99 88 push8 7f bb cc dd ee ff 99 88 add enddef")
	vm.Init(0)
	err = vm.Run(false)
	assert.NotNil(t, err, "add overflow")
	vm = buildVM(t, "handler 0 push8 7f bb cc dd ee ff 99 78 push8 ff bb cc dd ee ff 99 88 sub enddef")
	vm.Init(0)
	err = vm.Run(false)
	assert.NotNil(t, err, "sub overflow")
}

func TestNotNegIncDec(t *testing.T) {
	vm := buildVM(t, "handler 0 push1 7 not dup not push1 8 neg push1 4 inc push1 6 dec enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 1, -8, 5, 5)
}

func TestIf1(t *testing.T) {
	vm := buildVM(t, "handler 0 zero ifz push1 13 else push1 42 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIf2(t *testing.T) {
	vm := buildVM(t, "handler 0 zero ifnz push1 13 else push1 42 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 66, 17)
}

func TestIf3(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifz push1 13 else push1 42 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 66, 17)
}

func TestIf4(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifnz push1 13 else push1 42 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIf5(t *testing.T) {
	vm := buildVM(t, "handler 0 zero ifz push1 13 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIf6(t *testing.T) {
	vm := buildVM(t, "handler 0 zero ifnz push1 13 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 17)
}

func TestIf7(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifz push1 13 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 17)
}

func TestIf8(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifnz push1 13 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIfNested1(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifnz push1 13 zero ifz push1 15 else push1 13 endif endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 21, 17)
}

func TestIfNested2(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifz push1 13 zero ifz push1 15 else push1 13 endif endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 17)
}

func TestIfNested3(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifnz push1 13 zero ifnz push1 15 else push1 13 endif endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 19, 17)
}

func TestIfNull1(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifnz endif enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack())
}

func TestIfNull2(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifnz else endif enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack())
}

func TestCompares1(t *testing.T) {
	vm := buildVM(t, "handler 0 one neg1 eq one neg1 lt one neg1 gt enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 0, 1)
}

func TestCompares2(t *testing.T) {
	vm := buildVM(t, "handler 0 one one eq one one lt one one gt enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 1, 0, 0)
}

func TestCompares3(t *testing.T) {
	vm := buildVM(t, "handler 0 neg1 one eq neg1 one lt neg1 one gt enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 1, 0)
}

func TestCompares4(t *testing.T) {
	vm := buildVM(t, "handler 0 neg1 pushb 8 1 2 3 4 5 6 7 8 eq enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestCompares5(t *testing.T) {
	vm := buildVM(t, `handler 0 pushb "hello" pushb "hi" dup2 eq pick 2 pick 2 lt enddef`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 1)
}

func TestCompareLists1(t *testing.T) {
	vm := buildVM(t, `handler 0 pushl zero append one append dup dup eq swap dup gt enddef`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 1, 0)
}

func TestCompareLists2(t *testing.T) {
	vm := buildVM(t, `handler 0 pushl zero append one append dup one append dup pick 2 eq swap roll 2 gt enddef`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 1)
}

func TestCompareLists3(t *testing.T) {
	vm := buildVM(t, `handler 0 pushl zero append one append dup one append swap dup pick 2 eq swap roll 2 gt enddef`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 0)
}

func TestCompare7(t *testing.T) {
	vm := buildVM(t, "handler 0 dup zero index pick 1 one index eq enddef")
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(3*i), NewNumber(3*i+1), NewNumber(3*i+2))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0)
}

func TestCompareTimestampGt(t *testing.T) {
	vm := buildVM(t, `
		handler 0
		pusht 2018-07-18T00:00:00Z pusht 2018-01-01T00:00:00Z
		gt
		pusht 2018-01-01T00:00:00Z pusht 2018-07-18T00:00:00Z
		gt
		pusht 2018-07-18T00:00:00Z pusht 2018-07-18T00:00:00Z
		gt
		enddef`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 1, 0, 0)
}

func TestCompareTimestampLt(t *testing.T) {
	vm := buildVM(t, `
		handler 0
		pusht 2018-07-18T00:00:00Z pusht 2018-01-01T00:00:00Z
		lt
		pusht 2018-01-01T00:00:00Z pusht 2018-07-18T00:00:00Z
		lt
		pusht 2018-07-18T00:00:00Z pusht 2018-07-18T00:00:00Z
		lt
		enddef`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 1, 0)
}

func TestCompareTimestampEq(t *testing.T) {
	vm := buildVM(t, `
		handler 0
		pusht 2018-07-18T00:00:00Z pusht 2018-01-01T00:00:00Z
		eq
		pusht 2018-01-01T00:00:00Z pusht 2018-07-18T00:00:00Z
		eq
		pusht 2018-07-18T00:00:00Z pusht 2018-07-18T00:00:00Z
		eq
		enddef`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 0, 1)
}

func TestTimestamp1(t *testing.T) {
	vm := buildVM(t, `
		handler 0
		pusht 2018-07-18T00:00:00Z
		pusht 2018-01-01T00:00:00Z
		sub
		push3 40 42 0f
		div
		push1 3C
		dup
		mul
		push1 18
		mul
		div
		enddef
		`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 198)
}

func TestTimestampNegative(t *testing.T) {
	vm := buildVM(t, `
		handler 0
		pusht 2018-01-01T00:00:00Z
		pusht 2018-07-18T00:00:00Z
		sub
		push3 40 42 0f
		div
		push1 3C
		dup
		mul
		push1 18
		mul
		div
		enddef
		`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -198)
}

func TestTimestampInjectedNow(t *testing.T) {
	vm := buildVM(t, `
		handler 0
		now
		pusht 2018-01-01T00:00:00Z
		sub
		enddef
		`)
	ts, err := ParseTimestamp("2018-01-02T03:04:05Z")
	assert.Nil(t, err)
	now, err := NewCachingNow(ts)
	assert.Nil(t, err)
	vm.SetNow(now)
	vm.Init(0)
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 97445000000)
}

func TestTimestampDefaultNow(t *testing.T) {
	vm := buildVM(t, `
		handler 0
		now
		dup
		pusht 2018-01-01T00:00:00Z
		lt
		swap
		now
		sub
		zero
		eq
		pusht 2022-02-02T22:22:22Z
		now
		gt
		enddef
		`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 0, 1)
}

func TestInjectedRand(t *testing.T) {
	vm := buildVM(t, "handler 0 rand rand eq rand rand eq enddef")
	r, err := NewCachingRand()
	assert.Nil(t, err)
	vm.SetRand(r)
	vm.Init(0)
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 1, 1)
}

func TestDefaultRand(t *testing.T) {
	vm := buildVM(t, "handler 0 rand rand eq rand rand eq enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 0)
}

func TestList1(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl one append push1 7 append dup len swap one index enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 2, 7)
}

func TestExtend(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl one append push1 7 append dup zero append swap extend dup len swap push1 2 index enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 5, 0)
}

func TestSlice(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl zero append one append push1 2 append one push1 3 slice len enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 2)
}

func TestSum(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl zero append one append push1 2 append push1 3 append sum enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 6)
}

type seededRand struct {
	n int64
}

// RandInt implements Randomer for seededRand
func (r seededRand) RandInt() (int64, error) {
	return r.n, nil
}

func TestChoice(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl zero append one append push1 2 append push1 3 append choice enddef")
	r := seededRand{n: 12345}
	vm.SetRand(r)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3)
}

func TestWChoice1(t *testing.T) {
	vm := buildVM(t, "handler 0 wchoice 0 field 0 enddef")
	r := seededRand{n: math.MaxInt64 / 2}
	vm.SetRand(r)

	l := NewList()
	for i := int64(0); i < 6; i++ {
		st := NewStruct(NewNumber(i))
		l = l.Append(st)
	}
	vm.Init(0, l)

	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 4)
}

func TestWChoice2(t *testing.T) {
	vm := buildVM(t, "handler 0 wchoice 0 field 0 enddef")
	r := seededRand{n: math.MaxInt64 / 2}
	vm.SetRand(r)

	l := NewList()
	for i := int64(0); i < 6; i++ {
		st := NewStruct(NewNumber(6 - i))
		l = l.Append(st)
	}
	vm.Init(0, l)

	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 5)
}

func TestWChoiceErr(t *testing.T) {
	vm := buildVM(t, "handler 0 wchoice 0 field 0 enddef")
	r := seededRand{n: math.MaxInt64 / 2}
	vm.SetRand(r)

	l := NewList()
	vm.Init(0, l)

	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestAvg(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl one append push1 7 append push1 16 append avg enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 10)
}

func TestMin(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl one append push1 2 append push1 3 append min enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 1)
}

func TestMax(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl one append push1 2 append push1 3 append max enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3)
}

func TestAvgFail(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl avg enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestField(t *testing.T) {
	vm := buildVM(t, "handler 0 field 2 enddef")
	st := NewStruct(NewNumber(3), NewNumber(9), NewNumber(27))
	vm.Init(0, st)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 27)
}

func TestFieldFail(t *testing.T) {
	vm := buildVM(t, "handler 0 field 9 enddef")
	st := NewStruct(NewNumber(3), NewNumber(9), NewNumber(27))
	vm.Init(0, st)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestFieldL(t *testing.T) {
	vm := buildVM(t, "handler 0 fieldl 2 one index enddef")
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(3*i), NewNumber(3*i+1), NewNumber(3*i+2))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 5)
}

func TestFieldLFail(t *testing.T) {
	vm := buildVM(t, "handler 0 fieldl 9 one index enddef")
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(3*i), NewNumber(3*i+1), NewNumber(3*i+2))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestSortFields(t *testing.T) {
	vm := buildVM(t, "handler 0 sort 2 push1 3 index field 1 enddef")
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(2*i), NewNumber(3*i+1), NewNumber(4*(6-i)))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 4)
}

func TestSortFail(t *testing.T) {
	vm := buildVM(t, "handler 0 sort 6 push1 3 index field 1 enddef")
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(2*i), NewNumber(3*i+1), NewNumber(4*(6-i)))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestNestingFail1(t *testing.T) {
	buildVMfail(t, "def 1 nop enddef")
	buildVMfail(t, "handler 0 nop enddef handler 0 nop enddef")
	buildVMfail(t, "handler 0 nop enddef def 2 nop enddef")
	buildVMfail(t, "handler 0 ifz enddef")
	buildVMfail(t, "handler 0 ifnz enddef")
	buildVMfail(t, "handler 0 enddef enddef")
	buildVMfail(t, "handler 0 ifz else else enddef enddef")
	buildVMfail(t, "handler 0 push8 enddef")
}

func TestCall1(t *testing.T) {
	vm := buildVM(t, "handler 0 one call 0 1 enddef def 0 push1 2 add enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3)
}

func TestCall2(t *testing.T) {
	vm := buildVM(t, `
		handler 0 one call 0 1 enddef
		def 0 push1 2 call 1 2 enddef
		def 1 add enddef
	`)
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3)
}

func TestCallFail1(t *testing.T) {
	vm := buildVM(t, "handler 0 one call 1 1 enddef def 0 push1 2 add enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestCallFail2(t *testing.T) {
	vm := buildVM(t, "handler 0 one call 0 2 enddef def 0 push1 2 add enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestCallFail3(t *testing.T) {
	vm := buildVM(t, "handler 0 one call 0 1 enddef def 0 push1 2 add drop enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestDeco1(t *testing.T) {
	vm := buildVM(t, `
		handler 0 deco 0 0 fieldl 2 sum enddef
		def 0 dup field 0 dup mul swap  field 1 dup mul add enddef
	`)
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(2*i), NewNumber(3*i+1))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 455)
}

func TestStringers(t *testing.T) {
	assert.Equal(t, "Call", OpCall.String())
	vid := NewBytes([]byte("hi"))
	assert.Equal(t, "hi", vid.String())
	vn := NewNumber(123)
	assert.Equal(t, "123", vn.String())
	vt := NewTimestamp(0)
	assert.Equal(t, "2000-01-01T00:00:00Z", vt.String())
	vl := NewList()
	vl = vl.Append(NewBytes([]byte("July"))).Append(NewNumber(18))
	assert.Equal(t, "[July, 18]", vl.String())
	vs := NewStruct(NewBytes([]byte("July")), NewNumber(18))
	assert.Equal(t, "str(0)[July, 18]", vs.String())
}

func TestExerciseStrings(t *testing.T) {
	vm := buildVM(t, "handler 0 sort 6 push1 3 index field 1 enddef")
	vm.Init(0)

	assert.Contains(t, vm.String(), "Sort")
	da, n := vm.Disassemble(4)
	assert.Equal(t, 2, n)
	assert.Contains(t, da, "Push1")
}

func TestLookup1(t *testing.T) {
	vm := buildVM(t, `
		handler 0 lookup 0 0 enddef
		def 0 field 0 push1 4 gt enddef
	`)
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(2*i), NewNumber(3*i+1))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3)
}

func TestLookup2(t *testing.T) {
	vm := buildVM(t, `
		handler 0 lookup 0 0 enddef
		def 0 field 1 push1 4 gt enddef
	`)
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(2*i), NewNumber(3*i+1))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 2)
}

func TestLookupFail1(t *testing.T) {
	vm := buildVM(t, `
		handler 0 lookup 0 0 enddef
		def 0 field 1 push1 FF gt enddef
	`)
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(2*i), NewNumber(3*i+1))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestUnimplemented(t *testing.T) {
	// first make sure that the validation check forbids an invalid opcode
	buildVMfail(t, "handler 0 FF enddef")

	// now let's hack a VM after it passes validation to contain illegal data
	vm := buildVM(t, "handler 0 NOP enddef")
	// replace the nop with FF and try to run it; should still fail
	vm.code[3] = Opcode(0xFF)
	vm.Init(0)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestUnderflows(t *testing.T) {
	p := regexp.MustCompile("[[:space:]]+")
	keywords := p.Split(`drop drop2 dup dup2 swap over
		add sub mul div mod divmod muldiv not neg inc dec
		eq lt gt index len append extend slice sum avg max min`, -1)
	for _, k := range keywords {
		prog := "handler 0 " + k + " enddef"
		vm := buildVM(t, prog)
		vm.Init(0)
		err := vm.Run(false)
		assert.NotNil(t, err)
		correct := strings.HasPrefix(err.Error(), "stack underflow") ||
			strings.HasPrefix(err.Error(), "stack index error")
		assert.True(t, correct, "Keyword=%s msg=%s", k, err)
	}
}

func TestDisableOpcode(t *testing.T) {
	// now let's hack a VM after it passes validation to contain illegal data
	vm := buildVM(t, "handler 0 NOP enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)

	DisableOpcode(OpNop)
	// now the validation check should fail an invalid opcode
	buildVMfail(t, "handler 0 NOP enddef")
	// but we have to re-enable Nop or other tests might fail
	EnabledOpcodes.Set(int(OpNop))
}

func TestNegativeIndex(t *testing.T) {
	prog := `Handler 00
		Neg1 Index
		EndDef`
	vm := buildVM(t, prog)
	vm.Init(0, NewList().Append(NewNumber(1)).Append(NewNumber(2)))
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestIndex2(t *testing.T) {
	// this test is making sure that the 8f embedded into the PushB doesn't
	// cause skipToMatchingBracket to fail
	prog := `Handler 00
		IfZ
		PushB 8 b6 42 59 a3 8f 28 81 70
		EndIf
		EndDef`
	vm := buildVM(t, prog)
	vm.Init(0, NewList().Append(NewNumber(1)).Append(NewNumber(2)))
	err := vm.Run(false)
	assert.Nil(t, err)
}

func TestNoHandlers(t *testing.T) {
	// this test is making sure that if no handlers are defined, the vm won't load
	buildVMfail(t, "Def 00 Nop EndDef")
}

func TestMultipleHandlers(t *testing.T) {
	// this tests that we can define and call multiple handlers
	prog := `handler 1 10 add enddef
		handler 1 8 sub enddef
		handler 1 30 mul enddef
		handler 1 3 div enddef
		`
	vm := buildVM(t, prog)

	vm.Init(16, NewNumber(12), NewNumber(4))
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 16)

	vm.Init(8, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 8)

	vm.Init(48, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 48)

	vm.Init(3, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3)
}

func TestDefaultHandler(t *testing.T) {
	// this tests that we can define a different default handler
	// and also that it gets called if we invoke an event not in our list
	prog := `handler 1 10 add enddef
		handler 1 8 sub enddef
		handler 0 mod enddef
		handler 1 30 mul enddef
		handler 1 3 div enddef
		`
	vm := buildVM(t, prog)

	vm.Init(0, NewNumber(12), NewNumber(4))
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0)

	vm.Init(8, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 8)

	vm.Init(77, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0)
}

func TestMultipleEvents(t *testing.T) {
	// this tests that we can define a different default handler
	// and also that it gets called if we invoke an event not in our list
	prog := `handler 3 10 12 14 add enddef
		handler 2 0 5 mul enddef
		`
	vm := buildVM(t, prog)

	vm.Init(18, NewNumber(12), NewNumber(4))
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 16)

	vm.Init(18, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 16)

	vm.Init(20, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 16)

	vm.Init(0, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 48)

	vm.Init(5, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 48)

	vm.Init(77, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 48)
}

func TestHandlerIDs(t *testing.T) {
	// this tests that the HandlerIDs function works right
	prog := `handler 3 10 12 14 add enddef
		handler 2 0 5 mul enddef
		`
	vm := buildVM(t, prog)
	assert.Equal(t, []int{0, 5, 16, 18, 20}, vm.HandlerIDs())

	prog = `handler 1 1 add enddef`
	vm = buildVM(t, prog)
	assert.Equal(t, []int{1}, vm.HandlerIDs())

	prog = `handler 0 add enddef`
	vm = buildVM(t, prog)
	assert.Equal(t, []int{0}, vm.HandlerIDs())
}
