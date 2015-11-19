// automatically generated, do not modify

package flatapi

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Request struct {
	_tab flatbuffers.Table
}

func GetRootAsRequest(buf []byte, offset flatbuffers.UOffsetT) *Request {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Request{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Request) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Request) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Request) DataType() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Request) Data(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func RequestStart(builder *flatbuffers.Builder) { builder.StartObject(3) }
func RequestAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0) }
func RequestAddDataType(builder *flatbuffers.Builder, dataType byte) { builder.PrependByteSlot(1, dataType, 0) }
func RequestAddData(builder *flatbuffers.Builder, data flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(data), 0) }
func RequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
