package vm

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Opcode) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 byte
		zb0001, err = dc.ReadByte()
		if err != nil {
			return
		}
		(*z) = Opcode(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Opcode) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteByte(byte(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Opcode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendByte(o, byte(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Opcode) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 byte
		zb0001, bts, err = msgp.ReadByteBytes(bts)
		if err != nil {
			return
		}
		(*z) = Opcode(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Opcode) Msgsize() (s int) {
	s = msgp.ByteSize
	return
}