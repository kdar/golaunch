// automatically generated, do not modify

package flatapi

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Result struct {
	_tab flatbuffers.Table
}

func (rcv *Result) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Result) ResultType() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Result) Result(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func ResultStart(builder *flatbuffers.Builder) { builder.StartObject(2) }
func ResultAddResultType(builder *flatbuffers.Builder, resultType byte) { builder.PrependByteSlot(0, resultType, 0) }
func ResultAddResult(builder *flatbuffers.Builder, result flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(result), 0) }
func ResultEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
