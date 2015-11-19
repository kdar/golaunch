// automatically generated, do not modify

package flatapi

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Query struct {
	_tab flatbuffers.Table
}

func (rcv *Query) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Query) Query() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func QueryStart(builder *flatbuffers.Builder) { builder.StartObject(1) }
func QueryAddQuery(builder *flatbuffers.Builder, query flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(query), 0) }
func QueryEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
