package main

import (
	"errors"
	"strconv"
	"time"
)

// Some masking values
const (
	ByteMask byte = 0xFF
	HighBit  byte = 0x80
)

// Constants for Contexts
const (
	CtxTest        byte = iota
	CtxNodePayout  byte = iota
	CtxEaiTiming   byte = iota
	CtxNodeQuality byte = iota
	CtxMarketPrice byte = iota
)

// Constants for time
const (
	EpochStart      = "2018-01-01T00:00:00Z"
	TimestampFormat = "2006-01-02T15:04:05Z"
)

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

// Node is the fundamental unit that the parser manipulates (it builds a structure of nodes).
//Each node can emit itself as an array of bytes, or nil.
type Node interface {
	bytes() []byte
}

// Script is the highest level node in the system
type Script struct {
	preamble Node
	opcodes  []Node
}

func (n Script) bytes() []byte {
	b := append([]byte{}, n.preamble.bytes()...)
	for _, op := range n.opcodes {
		b = append(b, op.bytes()...)
	}
	return b
}

func newScript(p interface{}, opcodes interface{}) (Script, error) {
	preamble, ok := p.(PreambleNode)
	if !ok {
		return Script{}, errors.New("not a preamble node")
	}
	sl := toIfaceSlice(opcodes)
	ops := []Node{}
	for _, v := range sl {
		if n, ok := v.(Node); ok {
			ops = append(ops, n)
		}
	}
	return Script{preamble: preamble, opcodes: ops}, nil
}

// PreambleNode expresses the information in the preamble (which for now is just a context byte)
type PreambleNode struct {
	context byte
}

func (n PreambleNode) bytes() []byte {
	return []byte{n.context}
}

func newPreambleNode(ctx byte) (PreambleNode, error) {
	return PreambleNode{context: ctx}, nil
}

// UnitaryOpcode is for opcodes that cannot take arguments
type UnitaryOpcode struct {
	opcode byte
}

func (n UnitaryOpcode) bytes() []byte {
	return []byte{n.opcode}
}

func newUnitaryOpcode(b byte) (UnitaryOpcode, error) {
	return UnitaryOpcode{opcode: b}, nil
}

// BinaryOpcode is for opcodes that take one single-byte argument
type BinaryOpcode struct {
	opcode byte
	value  byte
}

func (n BinaryOpcode) bytes() []byte {
	return []byte{n.opcode, n.value}
}

func newBinaryOpcode(b byte, v string) (BinaryOpcode, error) {
	n, err := strconv.ParseUint(v, 0, 8)
	if err != nil {
		return BinaryOpcode{}, err
	}
	return BinaryOpcode{opcode: b, value: byte(n)}, nil
}

// toBytes returns an array of 8 bytes encoding n as a uint in little-endian form
func toBytesU(n uint64) []byte {
	b := []byte{}
	a := n
	for nbytes := 0; nbytes < 8; nbytes++ {
		b = append(b, byte(a)&ByteMask)
		a >>= 8
	}
	return b
}

// toBytes returns an array of 8 bytes encoding n as a signed value in little-endian form
func toBytes(n int64) []byte {
	b := []byte{}
	a := n
	for nbytes := 0; nbytes < 8; nbytes++ {
		b = append(b, byte(a)&ByteMask)
		a >>= 8
	}
	return b
}

// PushOpcode constructs push operations with the appropriate number of bytes to express
// the specified value. It has special cases for the special opcodes zero, one, and neg1.
type PushOpcode struct {
	arg int64
}

// This function builds a sequence of bytes consisting of either:
//   A ZERO, ONE, or NEG1 opcode
// OR
//   A PushN opcode followed by N bytes, where N is a value from 1-8.
//   The bytes are a representation of the value in little-endian order (low
//   byte first). The highest bit is the sign bit.
func (n PushOpcode) bytes() []byte {
	switch n.arg {
	case 0:
		return []byte{OpZero}
	case 1:
		return []byte{OpOne}
	case -1:
		return []byte{OpNeg1}
	default:
		b := toBytes(n.arg)
		var suppress byte
		if n.arg < 0 {
			suppress = byte(0xFF)
		}
		for b[len(b)-1] == suppress {
			b = b[:len(b)-1]
		}
		nbytes := byte(len(b))
		op := OpPushN | nbytes
		b = append([]byte{op}, b...)
		return b
	}
}

func newPushOpcode(s string) (PushOpcode, error) {
	v, err := strconv.ParseInt(s, 0, 64)
	return PushOpcode{arg: v}, err
}

// Push64 is a 64-bit unsigned value
type Push64 struct {
	u uint64
}

func newPush64(s string) (Push64, error) {
	v, err := strconv.ParseUint(s, 0, 64)
	return Push64{u: v}, err
}

func (n Push64) bytes() []byte {
	return append([]byte{OpPush64}, toBytesU(n.u)...)
}

// PushTimestamp is a 64-bit representation of the time since the start of the epoch in microseconds
type PushTimestamp struct {
	t uint64
}

func newPushTimestamp(s string) (PushTimestamp, error) {
	epoch, err := time.Parse(TimestampFormat, EpochStart)
	if err != nil {
		panic("Epoch isn't a valid timestamp!")
	}
	ts, err := time.Parse(TimestampFormat, s)
	if err != nil {
		return PushTimestamp{}, err
	}
	return PushTimestamp{uint64(ts.Sub(epoch).Nanoseconds() / 1000)}, nil // durations are in nanoseconds but we want microseconds
}

func (n PushTimestamp) bytes() []byte {
	return append([]byte{OpPushT}, toBytesU(n.t)...)
}
