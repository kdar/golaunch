// automatically generated, do not modify

package flatapi

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Response struct {
	_tab flatbuffers.Table
}

func GetRootAsResponse(buf []byte, offset flatbuffers.UOffsetT) *Response {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Response{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Response) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Response) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Response) Result(obj *Result, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
	if obj == nil {
		obj = new(Result)
	}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Response) ResultLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func ResponseStart(builder *flatbuffers.Builder) { builder.StartObject(2) }
func ResponseAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0) }
func ResponseAddResult(builder *flatbuffers.Builder, result flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(result), 0) }
func ResponseStartResultVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func ResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
