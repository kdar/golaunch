// automatically generated, do not modify

package flatapi

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type QueryResult struct {
	_tab flatbuffers.Table
}

func GetRootAsQueryResult(buf []byte, offset flatbuffers.UOffsetT) *QueryResult {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &QueryResult{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *QueryResult) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *QueryResult) Image() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *QueryResult) Title() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *QueryResult) Subtitle() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *QueryResult) Score() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *QueryResult) Query() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func QueryResultStart(builder *flatbuffers.Builder) { builder.StartObject(5) }
func QueryResultAddImage(builder *flatbuffers.Builder, image flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(image), 0) }
func QueryResultAddTitle(builder *flatbuffers.Builder, title flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(title), 0) }
func QueryResultAddSubtitle(builder *flatbuffers.Builder, subtitle flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(subtitle), 0) }
func QueryResultAddScore(builder *flatbuffers.Builder, score int32) { builder.PrependInt32Slot(3, score, 0) }
func QueryResultAddQuery(builder *flatbuffers.Builder, query flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(query), 0) }
func QueryResultEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
