// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: server/proto/internal.proto

/*
	Package proto is a generated protocol buffer package.

	It is generated from these files:
		server/proto/internal.proto

	It has these top-level messages:
		RaftLog
		CreateStreamOp
		Stream
		RaftJoinRequest
		RaftJoinResponse
		MetadataSnapshot
		ReplicationRequest
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type RaftLog_Op int32

const (
	RaftLog_CREATE_STREAM RaftLog_Op = 0
)

var RaftLog_Op_name = map[int32]string{
	0: "CREATE_STREAM",
}
var RaftLog_Op_value = map[string]int32{
	"CREATE_STREAM": 0,
}

func (x RaftLog_Op) String() string {
	return proto1.EnumName(RaftLog_Op_name, int32(x))
}
func (RaftLog_Op) EnumDescriptor() ([]byte, []int) { return fileDescriptorInternal, []int{0, 0} }

type RaftLog struct {
	Op             RaftLog_Op      `protobuf:"varint,1,opt,name=op,proto3,enum=proto.RaftLog_Op" json:"op,omitempty"`
	CreateStreamOp *CreateStreamOp `protobuf:"bytes,2,opt,name=createStreamOp" json:"createStreamOp,omitempty"`
}

func (m *RaftLog) Reset()                    { *m = RaftLog{} }
func (m *RaftLog) String() string            { return proto1.CompactTextString(m) }
func (*RaftLog) ProtoMessage()               {}
func (*RaftLog) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{0} }

func (m *RaftLog) GetOp() RaftLog_Op {
	if m != nil {
		return m.Op
	}
	return RaftLog_CREATE_STREAM
}

func (m *RaftLog) GetCreateStreamOp() *CreateStreamOp {
	if m != nil {
		return m.CreateStreamOp
	}
	return nil
}

type CreateStreamOp struct {
	Stream *Stream `protobuf:"bytes,1,opt,name=stream" json:"stream,omitempty"`
}

func (m *CreateStreamOp) Reset()                    { *m = CreateStreamOp{} }
func (m *CreateStreamOp) String() string            { return proto1.CompactTextString(m) }
func (*CreateStreamOp) ProtoMessage()               {}
func (*CreateStreamOp) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{1} }

func (m *CreateStreamOp) GetStream() *Stream {
	if m != nil {
		return m.Stream
	}
	return nil
}

type Stream struct {
	Subject           string   `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	Name              string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	ConsumerGroup     string   `protobuf:"bytes,3,opt,name=consumerGroup,proto3" json:"consumerGroup,omitempty"`
	ReplicationFactor int32    `protobuf:"varint,4,opt,name=replicationFactor,proto3" json:"replicationFactor,omitempty"`
	Replicas          []string `protobuf:"bytes,5,rep,name=replicas" json:"replicas,omitempty"`
	Leader            string   `protobuf:"bytes,6,opt,name=leader,proto3" json:"leader,omitempty"`
	Isr               []string `protobuf:"bytes,7,rep,name=isr" json:"isr,omitempty"`
}

func (m *Stream) Reset()                    { *m = Stream{} }
func (m *Stream) String() string            { return proto1.CompactTextString(m) }
func (*Stream) ProtoMessage()               {}
func (*Stream) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{2} }

func (m *Stream) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *Stream) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Stream) GetConsumerGroup() string {
	if m != nil {
		return m.ConsumerGroup
	}
	return ""
}

func (m *Stream) GetReplicationFactor() int32 {
	if m != nil {
		return m.ReplicationFactor
	}
	return 0
}

func (m *Stream) GetReplicas() []string {
	if m != nil {
		return m.Replicas
	}
	return nil
}

func (m *Stream) GetLeader() string {
	if m != nil {
		return m.Leader
	}
	return ""
}

func (m *Stream) GetIsr() []string {
	if m != nil {
		return m.Isr
	}
	return nil
}

// RaftJoinRequest is a request to join a Raft group.
type RaftJoinRequest struct {
	NodeID   string `protobuf:"bytes,1,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
	NodeAddr string `protobuf:"bytes,2,opt,name=nodeAddr,proto3" json:"nodeAddr,omitempty"`
}

func (m *RaftJoinRequest) Reset()                    { *m = RaftJoinRequest{} }
func (m *RaftJoinRequest) String() string            { return proto1.CompactTextString(m) }
func (*RaftJoinRequest) ProtoMessage()               {}
func (*RaftJoinRequest) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{3} }

func (m *RaftJoinRequest) GetNodeID() string {
	if m != nil {
		return m.NodeID
	}
	return ""
}

func (m *RaftJoinRequest) GetNodeAddr() string {
	if m != nil {
		return m.NodeAddr
	}
	return ""
}

// RaftJoinResponse is a response to a RaftJoinRequest.
type RaftJoinResponse struct {
	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (m *RaftJoinResponse) Reset()                    { *m = RaftJoinResponse{} }
func (m *RaftJoinResponse) String() string            { return proto1.CompactTextString(m) }
func (*RaftJoinResponse) ProtoMessage()               {}
func (*RaftJoinResponse) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{4} }

func (m *RaftJoinResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type MetadataSnapshot struct {
	Streams []*Stream `protobuf:"bytes,1,rep,name=streams" json:"streams,omitempty"`
}

func (m *MetadataSnapshot) Reset()                    { *m = MetadataSnapshot{} }
func (m *MetadataSnapshot) String() string            { return proto1.CompactTextString(m) }
func (*MetadataSnapshot) ProtoMessage()               {}
func (*MetadataSnapshot) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{5} }

func (m *MetadataSnapshot) GetStreams() []*Stream {
	if m != nil {
		return m.Streams
	}
	return nil
}

type ReplicationRequest struct {
	ReplicaID     string `protobuf:"bytes,1,opt,name=replicaID,proto3" json:"replicaID,omitempty"`
	HighWatermark int64  `protobuf:"varint,2,opt,name=highWatermark,proto3" json:"highWatermark,omitempty"`
}

func (m *ReplicationRequest) Reset()                    { *m = ReplicationRequest{} }
func (m *ReplicationRequest) String() string            { return proto1.CompactTextString(m) }
func (*ReplicationRequest) ProtoMessage()               {}
func (*ReplicationRequest) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{6} }

func (m *ReplicationRequest) GetReplicaID() string {
	if m != nil {
		return m.ReplicaID
	}
	return ""
}

func (m *ReplicationRequest) GetHighWatermark() int64 {
	if m != nil {
		return m.HighWatermark
	}
	return 0
}

func init() {
	proto1.RegisterType((*RaftLog)(nil), "proto.RaftLog")
	proto1.RegisterType((*CreateStreamOp)(nil), "proto.CreateStreamOp")
	proto1.RegisterType((*Stream)(nil), "proto.Stream")
	proto1.RegisterType((*RaftJoinRequest)(nil), "proto.RaftJoinRequest")
	proto1.RegisterType((*RaftJoinResponse)(nil), "proto.RaftJoinResponse")
	proto1.RegisterType((*MetadataSnapshot)(nil), "proto.MetadataSnapshot")
	proto1.RegisterType((*ReplicationRequest)(nil), "proto.ReplicationRequest")
	proto1.RegisterEnum("proto.RaftLog_Op", RaftLog_Op_name, RaftLog_Op_value)
}
func (m *RaftLog) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RaftLog) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Op != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintInternal(dAtA, i, uint64(m.Op))
	}
	if m.CreateStreamOp != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintInternal(dAtA, i, uint64(m.CreateStreamOp.Size()))
		n1, err := m.CreateStreamOp.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func (m *CreateStreamOp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CreateStreamOp) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Stream != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintInternal(dAtA, i, uint64(m.Stream.Size()))
		n2, err := m.Stream.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func (m *Stream) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Stream) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Subject) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintInternal(dAtA, i, uint64(len(m.Subject)))
		i += copy(dAtA[i:], m.Subject)
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintInternal(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.ConsumerGroup) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintInternal(dAtA, i, uint64(len(m.ConsumerGroup)))
		i += copy(dAtA[i:], m.ConsumerGroup)
	}
	if m.ReplicationFactor != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintInternal(dAtA, i, uint64(m.ReplicationFactor))
	}
	if len(m.Replicas) > 0 {
		for _, s := range m.Replicas {
			dAtA[i] = 0x2a
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	if len(m.Leader) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintInternal(dAtA, i, uint64(len(m.Leader)))
		i += copy(dAtA[i:], m.Leader)
	}
	if len(m.Isr) > 0 {
		for _, s := range m.Isr {
			dAtA[i] = 0x3a
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	return i, nil
}

func (m *RaftJoinRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RaftJoinRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.NodeID) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintInternal(dAtA, i, uint64(len(m.NodeID)))
		i += copy(dAtA[i:], m.NodeID)
	}
	if len(m.NodeAddr) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintInternal(dAtA, i, uint64(len(m.NodeAddr)))
		i += copy(dAtA[i:], m.NodeAddr)
	}
	return i, nil
}

func (m *RaftJoinResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RaftJoinResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Error) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintInternal(dAtA, i, uint64(len(m.Error)))
		i += copy(dAtA[i:], m.Error)
	}
	return i, nil
}

func (m *MetadataSnapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MetadataSnapshot) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Streams) > 0 {
		for _, msg := range m.Streams {
			dAtA[i] = 0xa
			i++
			i = encodeVarintInternal(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *ReplicationRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReplicationRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ReplicaID) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintInternal(dAtA, i, uint64(len(m.ReplicaID)))
		i += copy(dAtA[i:], m.ReplicaID)
	}
	if m.HighWatermark != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintInternal(dAtA, i, uint64(m.HighWatermark))
	}
	return i, nil
}

func encodeVarintInternal(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *RaftLog) Size() (n int) {
	var l int
	_ = l
	if m.Op != 0 {
		n += 1 + sovInternal(uint64(m.Op))
	}
	if m.CreateStreamOp != nil {
		l = m.CreateStreamOp.Size()
		n += 1 + l + sovInternal(uint64(l))
	}
	return n
}

func (m *CreateStreamOp) Size() (n int) {
	var l int
	_ = l
	if m.Stream != nil {
		l = m.Stream.Size()
		n += 1 + l + sovInternal(uint64(l))
	}
	return n
}

func (m *Stream) Size() (n int) {
	var l int
	_ = l
	l = len(m.Subject)
	if l > 0 {
		n += 1 + l + sovInternal(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovInternal(uint64(l))
	}
	l = len(m.ConsumerGroup)
	if l > 0 {
		n += 1 + l + sovInternal(uint64(l))
	}
	if m.ReplicationFactor != 0 {
		n += 1 + sovInternal(uint64(m.ReplicationFactor))
	}
	if len(m.Replicas) > 0 {
		for _, s := range m.Replicas {
			l = len(s)
			n += 1 + l + sovInternal(uint64(l))
		}
	}
	l = len(m.Leader)
	if l > 0 {
		n += 1 + l + sovInternal(uint64(l))
	}
	if len(m.Isr) > 0 {
		for _, s := range m.Isr {
			l = len(s)
			n += 1 + l + sovInternal(uint64(l))
		}
	}
	return n
}

func (m *RaftJoinRequest) Size() (n int) {
	var l int
	_ = l
	l = len(m.NodeID)
	if l > 0 {
		n += 1 + l + sovInternal(uint64(l))
	}
	l = len(m.NodeAddr)
	if l > 0 {
		n += 1 + l + sovInternal(uint64(l))
	}
	return n
}

func (m *RaftJoinResponse) Size() (n int) {
	var l int
	_ = l
	l = len(m.Error)
	if l > 0 {
		n += 1 + l + sovInternal(uint64(l))
	}
	return n
}

func (m *MetadataSnapshot) Size() (n int) {
	var l int
	_ = l
	if len(m.Streams) > 0 {
		for _, e := range m.Streams {
			l = e.Size()
			n += 1 + l + sovInternal(uint64(l))
		}
	}
	return n
}

func (m *ReplicationRequest) Size() (n int) {
	var l int
	_ = l
	l = len(m.ReplicaID)
	if l > 0 {
		n += 1 + l + sovInternal(uint64(l))
	}
	if m.HighWatermark != 0 {
		n += 1 + sovInternal(uint64(m.HighWatermark))
	}
	return n
}

func sovInternal(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozInternal(x uint64) (n int) {
	return sovInternal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RaftLog) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInternal
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RaftLog: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RaftLog: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Op", wireType)
			}
			m.Op = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Op |= (RaftLog_Op(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreateStreamOp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CreateStreamOp == nil {
				m.CreateStreamOp = &CreateStreamOp{}
			}
			if err := m.CreateStreamOp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInternal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInternal
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *CreateStreamOp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInternal
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CreateStreamOp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CreateStreamOp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stream", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Stream == nil {
				m.Stream = &Stream{}
			}
			if err := m.Stream.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInternal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInternal
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Stream) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInternal
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Stream: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Stream: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subject", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Subject = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConsumerGroup", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConsumerGroup = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReplicationFactor", wireType)
			}
			m.ReplicationFactor = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ReplicationFactor |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Replicas", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Replicas = append(m.Replicas, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Leader", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Leader = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Isr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Isr = append(m.Isr, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInternal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInternal
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RaftJoinRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInternal
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RaftJoinRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RaftJoinRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NodeID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NodeAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInternal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInternal
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RaftJoinResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInternal
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RaftJoinResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RaftJoinResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Error = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInternal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInternal
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MetadataSnapshot) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInternal
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MetadataSnapshot: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetadataSnapshot: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Streams", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Streams = append(m.Streams, &Stream{})
			if err := m.Streams[len(m.Streams)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInternal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInternal
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ReplicationRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInternal
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ReplicationRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReplicationRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReplicaID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInternal
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReplicaID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HighWatermark", wireType)
			}
			m.HighWatermark = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HighWatermark |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipInternal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthInternal
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipInternal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowInternal
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowInternal
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthInternal
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowInternal
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipInternal(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthInternal = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowInternal   = fmt.Errorf("proto: integer overflow")
)

func init() { proto1.RegisterFile("server/proto/internal.proto", fileDescriptorInternal) }

var fileDescriptorInternal = []byte{
	// 443 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x50, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xee, 0xc6, 0x8d, 0x43, 0xa6, 0x6a, 0x70, 0x56, 0xfc, 0x58, 0x80, 0x22, 0x63, 0x81, 0xf0,
	0x01, 0xa5, 0x52, 0x38, 0x70, 0x40, 0x1c, 0x42, 0x09, 0x08, 0x44, 0x15, 0x69, 0x53, 0x09, 0x6e,
	0x68, 0x6b, 0x0f, 0x8d, 0x21, 0xf1, 0x2e, 0xb3, 0x1b, 0x9e, 0x80, 0x87, 0xe0, 0x91, 0x38, 0xf6,
	0x11, 0x50, 0x78, 0x11, 0xe4, 0xf5, 0x26, 0x51, 0xe8, 0x69, 0xe7, 0x9b, 0xef, 0x9b, 0x9d, 0x6f,
	0x3e, 0xb8, 0x6f, 0x90, 0x7e, 0x20, 0x9d, 0x68, 0x52, 0x56, 0x9d, 0x94, 0x95, 0x45, 0xaa, 0xe4,
	0x62, 0xe8, 0x20, 0x6f, 0xbb, 0x27, 0xfd, 0xc9, 0xa0, 0x23, 0xe4, 0x17, 0xfb, 0x41, 0x5d, 0xf2,
	0x87, 0xd0, 0x52, 0x3a, 0x66, 0x09, 0xcb, 0x7a, 0xa3, 0x7e, 0x23, 0x1b, 0x7a, 0x6e, 0x38, 0xd5,
	0xa2, 0xa5, 0x34, 0x7f, 0x09, 0xbd, 0x9c, 0x50, 0x5a, 0x9c, 0x59, 0x42, 0xb9, 0x9c, 0xea, 0xb8,
	0x95, 0xb0, 0xec, 0x68, 0x74, 0xdb, 0xcb, 0x4f, 0xf7, 0x48, 0xf1, 0x9f, 0x38, 0xbd, 0x0b, 0xad,
	0xa9, 0xe6, 0x7d, 0x38, 0x3e, 0x15, 0x93, 0xf1, 0xf9, 0xe4, 0xf3, 0xec, 0x5c, 0x4c, 0xc6, 0x67,
	0xd1, 0x41, 0xfa, 0x1c, 0x7a, 0xfb, 0xa3, 0xfc, 0x31, 0x84, 0xc6, 0xd5, 0xce, 0xd0, 0xd1, 0xe8,
	0xd8, 0x6f, 0x68, 0x04, 0xc2, 0x93, 0xe9, 0x15, 0x83, 0xb0, 0x69, 0xf1, 0x18, 0x3a, 0x66, 0x75,
	0xf1, 0x15, 0x73, 0xeb, 0x46, 0xba, 0x62, 0x03, 0x39, 0x87, 0xc3, 0x4a, 0x2e, 0xd1, 0x79, 0xed,
	0x0a, 0x57, 0xf3, 0x47, 0x70, 0x9c, 0xab, 0xca, 0xac, 0x96, 0x48, 0x6f, 0x49, 0xad, 0x74, 0x1c,
	0x38, 0x72, 0xbf, 0xc9, 0x9f, 0x42, 0x9f, 0x50, 0x2f, 0xca, 0x5c, 0xda, 0x52, 0x55, 0x6f, 0x64,
	0x6e, 0x15, 0xc5, 0x87, 0x09, 0xcb, 0xda, 0xe2, 0x3a, 0xc1, 0xef, 0xc1, 0x0d, 0xdf, 0x34, 0x71,
	0x3b, 0x09, 0xb2, 0xae, 0xd8, 0x62, 0x7e, 0x07, 0xc2, 0x05, 0xca, 0x02, 0x29, 0x0e, 0xdd, 0x22,
	0x8f, 0x78, 0x04, 0x41, 0x69, 0x28, 0xee, 0x38, 0x79, 0x5d, 0xa6, 0x13, 0xb8, 0x59, 0xa7, 0xfe,
	0x5e, 0x95, 0x95, 0xc0, 0xef, 0x2b, 0x34, 0xb6, 0x1e, 0xae, 0x54, 0x81, 0xef, 0x5e, 0xfb, 0xcb,
	0x3c, 0xaa, 0x17, 0xd6, 0xd5, 0xb8, 0x28, 0xc8, 0x1f, 0xb7, 0xc5, 0x69, 0x06, 0xd1, 0xee, 0x1b,
	0xa3, 0x55, 0x65, 0x90, 0xdf, 0x82, 0x36, 0x12, 0x29, 0xf2, 0xdf, 0x34, 0x20, 0x7d, 0x01, 0xd1,
	0x19, 0x5a, 0x59, 0x48, 0x2b, 0x67, 0x95, 0xd4, 0x66, 0xae, 0x2c, 0x7f, 0x02, 0x9d, 0x26, 0x61,
	0x13, 0xb3, 0x24, 0xb8, 0x9e, 0xff, 0x86, 0x4d, 0x3f, 0x01, 0x17, 0xbb, 0x20, 0x36, 0x86, 0x1f,
	0x40, 0xd7, 0x5f, 0xbe, 0xf5, 0xbc, 0x6b, 0xd4, 0xd9, 0xcf, 0xcb, 0xcb, 0xf9, 0x47, 0x69, 0x91,
	0x96, 0x92, 0xbe, 0x39, 0xef, 0x81, 0xd8, 0x6f, 0xbe, 0x8a, 0x7e, 0xaf, 0x07, 0xec, 0x6a, 0x3d,
	0x60, 0x7f, 0xd6, 0x03, 0xf6, 0xeb, 0xef, 0xe0, 0xe0, 0x22, 0x74, 0x16, 0x9e, 0xfd, 0x0b, 0x00,
	0x00, 0xff, 0xff, 0x74, 0xc0, 0x96, 0xc3, 0xd9, 0x02, 0x00, 0x00,
}
