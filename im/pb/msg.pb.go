// Code generated by protoc-gen-go.
// source: msg.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	msg.proto

It has these top-level messages:
	ImMsg
	RC
	DsMsg
	KV
	Evn
*/
package pb

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type ImMsg struct {
	I                *string  `protobuf:"bytes,1,opt,name=i" json:"i,omitempty" bson:"_id"`
	S                *string  `protobuf:"bytes,2,opt,name=s" json:"s,omitempty"`
	R                []string `protobuf:"bytes,3,rep,name=r" json:"r,omitempty"`
	T                *uint32  `protobuf:"varint,4,req,name=t" json:"t,omitempty"`
	D                *string  `protobuf:"bytes,5,opt,name=d" json:"d,omitempty"`
	C                []byte   `protobuf:"bytes,6,req,name=c" json:"c,omitempty"`
	A                *string  `protobuf:"bytes,7,opt,name=a" json:"a,omitempty"`
	Time             *int64   `protobuf:"varint,8,opt,name=time" json:"time,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *ImMsg) Reset()         { *m = ImMsg{} }
func (m *ImMsg) String() string { return proto.CompactTextString(m) }
func (*ImMsg) ProtoMessage()    {}

func (m *ImMsg) GetI() string {
	if m != nil && m.I != nil {
		return *m.I
	}
	return ""
}

func (m *ImMsg) GetS() string {
	if m != nil && m.S != nil {
		return *m.S
	}
	return ""
}

func (m *ImMsg) GetR() []string {
	if m != nil {
		return m.R
	}
	return nil
}

func (m *ImMsg) GetT() uint32 {
	if m != nil && m.T != nil {
		return *m.T
	}
	return 0
}

func (m *ImMsg) GetD() string {
	if m != nil && m.D != nil {
		return *m.D
	}
	return ""
}

func (m *ImMsg) GetC() []byte {
	if m != nil {
		return m.C
	}
	return nil
}

func (m *ImMsg) GetA() string {
	if m != nil && m.A != nil {
		return *m.A
	}
	return ""
}

func (m *ImMsg) GetTime() int64 {
	if m != nil && m.Time != nil {
		return *m.Time
	}
	return 0
}

type RC struct {
	R                *string `protobuf:"bytes,1,req,name=r" json:"r,omitempty"`
	C                *string `protobuf:"bytes,2,req,name=c" json:"c,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *RC) Reset()         { *m = RC{} }
func (m *RC) String() string { return proto.CompactTextString(m) }
func (*RC) ProtoMessage()    {}

func (m *RC) GetR() string {
	if m != nil && m.R != nil {
		return *m.R
	}
	return ""
}

func (m *RC) GetC() string {
	if m != nil && m.C != nil {
		return *m.C
	}
	return ""
}

type DsMsg struct {
	M                *ImMsg `protobuf:"bytes,1,req,name=m" json:"m,omitempty"`
	Rc               []*RC  `protobuf:"bytes,2,rep,name=rc" json:"rc,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *DsMsg) Reset()         { *m = DsMsg{} }
func (m *DsMsg) String() string { return proto.CompactTextString(m) }
func (*DsMsg) ProtoMessage()    {}

func (m *DsMsg) GetM() *ImMsg {
	if m != nil {
		return m.M
	}
	return nil
}

func (m *DsMsg) GetRc() []*RC {
	if m != nil {
		return m.Rc
	}
	return nil
}

type KV struct {
	Key              *string `protobuf:"bytes,1,req,name=key" json:"key,omitempty"`
	Val              *string `protobuf:"bytes,2,req,name=val" json:"val,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *KV) Reset()         { *m = KV{} }
func (m *KV) String() string { return proto.CompactTextString(m) }
func (*KV) ProtoMessage()    {}

func (m *KV) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *KV) GetVal() string {
	if m != nil && m.Val != nil {
		return *m.Val
	}
	return ""
}

type Evn struct {
	Uid              *string `protobuf:"bytes,1,req,name=uid" json:"uid,omitempty"`
	Name             *string `protobuf:"bytes,2,req,name=name" json:"name,omitempty"`
	Action           *string `protobuf:"bytes,3,req,name=action" json:"action,omitempty"`
	Time             *int64  `protobuf:"varint,4,req,name=time" json:"time,omitempty"`
	Type             *int32  `protobuf:"varint,5,req,name=type" json:"type,omitempty"`
	Kvs              []*KV   `protobuf:"bytes,6,rep,name=kvs" json:"kvs,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Evn) Reset()         { *m = Evn{} }
func (m *Evn) String() string { return proto.CompactTextString(m) }
func (*Evn) ProtoMessage()    {}

func (m *Evn) GetUid() string {
	if m != nil && m.Uid != nil {
		return *m.Uid
	}
	return ""
}

func (m *Evn) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Evn) GetAction() string {
	if m != nil && m.Action != nil {
		return *m.Action
	}
	return ""
}

func (m *Evn) GetTime() int64 {
	if m != nil && m.Time != nil {
		return *m.Time
	}
	return 0
}

func (m *Evn) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

func (m *Evn) GetKvs() []*KV {
	if m != nil {
		return m.Kvs
	}
	return nil
}

func init() {
}
