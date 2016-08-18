// Code generated by protoc-gen-gogo.
// source: plugin.proto
// DO NOT EDIT!

/*
	Package plugin is a generated protocol buffer package.

	It is generated from these files:
		plugin.proto

	It has these top-level messages:
		TLSAuthorization
*/
package plugin

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

import strings "strings"
import github_com_gogo_protobuf_proto "github.com/gogo/protobuf/proto"
import sort "sort"
import strconv "strconv"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.GoGoProtoPackageIsVersion1

type TLSAuthorization struct {
	// Roles contains the acceptable TLS OU roles for the handler.
	Roles []string `protobuf:"bytes,1,rep,name=roles" json:"roles,omitempty"`
	// Insecure is set to true if this method does not require
	// authorization. NOTE: Specifying both "insecure" and a nonempty
	// list of roles is invalid. This would fail at codegen time.
	Insecure         *bool  `protobuf:"varint,2,opt,name=insecure" json:"insecure,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *TLSAuthorization) Reset()                    { *m = TLSAuthorization{} }
func (*TLSAuthorization) ProtoMessage()               {}
func (*TLSAuthorization) Descriptor() ([]byte, []int) { return fileDescriptorPlugin, []int{0} }

var E_Deepcopy = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         70000,
	Name:          "docker.protobuf.plugin.deepcopy",
	Tag:           "varint,70000,opt,name=deepcopy,def=1",
}

var E_TlsAuthorization = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*TLSAuthorization)(nil),
	Field:         73626345,
	Name:          "docker.protobuf.plugin.tls_authorization",
	Tag:           "bytes,73626345,opt,name=tls_authorization,json=tlsAuthorization",
}

func init() {
	proto.RegisterType((*TLSAuthorization)(nil), "docker.protobuf.plugin.TLSAuthorization")
	proto.RegisterExtension(E_Deepcopy)
	proto.RegisterExtension(E_TlsAuthorization)
}
func (this *TLSAuthorization) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&plugin.TLSAuthorization{")
	if this.Roles != nil {
		s = append(s, "Roles: "+fmt.Sprintf("%#v", this.Roles)+",\n")
	}
	if this.Insecure != nil {
		s = append(s, "Insecure: "+valueToGoStringPlugin(this.Insecure, "bool")+",\n")
	}
	if this.XXX_unrecognized != nil {
		s = append(s, "XXX_unrecognized:"+fmt.Sprintf("%#v", this.XXX_unrecognized)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringPlugin(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func extensionToGoStringPlugin(e map[int32]github_com_gogo_protobuf_proto.Extension) string {
	if e == nil {
		return "nil"
	}
	s := "map[int32]proto.Extension{"
	keys := make([]int, 0, len(e))
	for k := range e {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	ss := []string{}
	for _, k := range keys {
		ss = append(ss, strconv.Itoa(k)+": "+e[int32(k)].GoString())
	}
	s += strings.Join(ss, ",") + "}"
	return s
}
func (m *TLSAuthorization) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *TLSAuthorization) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Roles) > 0 {
		for _, s := range m.Roles {
			data[i] = 0xa
			i++
			l = len(s)
			for l >= 1<<7 {
				data[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			data[i] = uint8(l)
			i++
			i += copy(data[i:], s)
		}
	}
	if m.Insecure != nil {
		data[i] = 0x10
		i++
		if *m.Insecure {
			data[i] = 1
		} else {
			data[i] = 0
		}
		i++
	}
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeFixed64Plugin(data []byte, offset int, v uint64) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	data[offset+4] = uint8(v >> 32)
	data[offset+5] = uint8(v >> 40)
	data[offset+6] = uint8(v >> 48)
	data[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Plugin(data []byte, offset int, v uint32) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintPlugin(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}
func (m *TLSAuthorization) Size() (n int) {
	var l int
	_ = l
	if len(m.Roles) > 0 {
		for _, s := range m.Roles {
			l = len(s)
			n += 1 + l + sovPlugin(uint64(l))
		}
	}
	if m.Insecure != nil {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovPlugin(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozPlugin(x uint64) (n int) {
	return sovPlugin(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *TLSAuthorization) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&TLSAuthorization{`,
		`Roles:` + fmt.Sprintf("%v", this.Roles) + `,`,
		`Insecure:` + valueToStringPlugin(this.Insecure) + `,`,
		`XXX_unrecognized:` + fmt.Sprintf("%v", this.XXX_unrecognized) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringPlugin(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *TLSAuthorization) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPlugin
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TLSAuthorization: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TLSAuthorization: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Roles", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlugin
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPlugin
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Roles = append(m.Roles, string(data[iNdEx:postIndex]))
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Insecure", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlugin
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			b := bool(v != 0)
			m.Insecure = &b
		default:
			iNdEx = preIndex
			skippy, err := skipPlugin(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPlugin
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipPlugin(data []byte) (n int, err error) {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPlugin
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
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
					return 0, ErrIntOverflowPlugin
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if data[iNdEx-1] < 0x80 {
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
					return 0, ErrIntOverflowPlugin
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthPlugin
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowPlugin
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := data[iNdEx]
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
				next, err := skipPlugin(data[start:])
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
	ErrInvalidLengthPlugin = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPlugin   = fmt.Errorf("proto: integer overflow")
)

var fileDescriptorPlugin = []byte{
	// 259 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0xc8, 0x29, 0x4d,
	0xcf, 0xcc, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x4b, 0xc9, 0x4f, 0xce, 0x4e, 0x2d,
	0x82, 0xf0, 0x92, 0x4a, 0xd3, 0xf4, 0x20, 0xb2, 0x52, 0x0a, 0xe9, 0xf9, 0xf9, 0xe9, 0x39, 0xa9,
	0xfa, 0x30, 0x71, 0xfd, 0x94, 0xd4, 0xe2, 0xe4, 0xa2, 0xcc, 0x82, 0x92, 0x7c, 0xa8, 0x5a, 0x25,
	0x17, 0x2e, 0x81, 0x10, 0x9f, 0x60, 0xc7, 0xd2, 0x92, 0x8c, 0xfc, 0xa2, 0xcc, 0xaa, 0xc4, 0x92,
	0xcc, 0xfc, 0x3c, 0x21, 0x11, 0x2e, 0xd6, 0xa2, 0xfc, 0x9c, 0xd4, 0x62, 0x09, 0x46, 0x05, 0x66,
	0x0d, 0xce, 0x20, 0x08, 0x47, 0x48, 0x8a, 0x8b, 0x23, 0x33, 0xaf, 0x38, 0x35, 0xb9, 0xb4, 0x28,
	0x55, 0x82, 0x49, 0x81, 0x51, 0x83, 0x23, 0x08, 0xce, 0xb7, 0x72, 0xe6, 0xe2, 0x48, 0x49, 0x4d,
	0x2d, 0x48, 0xce, 0x2f, 0xa8, 0x14, 0x92, 0xd7, 0x83, 0x58, 0x8a, 0x70, 0x8c, 0x6f, 0x6a, 0x71,
	0x71, 0x62, 0x7a, 0xaa, 0x7f, 0x01, 0xc8, 0xf4, 0x62, 0x89, 0x0f, 0x8b, 0x58, 0x40, 0xda, 0xad,
	0x58, 0x4a, 0x8a, 0x4a, 0x53, 0x83, 0xe0, 0x1a, 0xad, 0x2a, 0xb8, 0x04, 0x4b, 0x72, 0x8a, 0xe3,
	0x13, 0x51, 0xdc, 0x22, 0x87, 0xc5, 0xb4, 0x92, 0x8c, 0xfc, 0x14, 0x98, 0x61, 0x2f, 0x9f, 0xf6,
	0x2a, 0x2b, 0x30, 0x6a, 0x70, 0x1b, 0x69, 0xe8, 0x61, 0x0f, 0x03, 0x3d, 0x74, 0xef, 0x05, 0x09,
	0x94, 0xe4, 0x14, 0xa3, 0x88, 0x38, 0xc9, 0x9c, 0x78, 0x28, 0xc7, 0x70, 0xe3, 0xa1, 0x1c, 0xc3,
	0x87, 0x87, 0x72, 0x8c, 0x0d, 0x8f, 0xe4, 0x18, 0x4f, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e,
	0xf1, 0xc1, 0x23, 0x39, 0x46, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x04, 0x4e, 0xf8, 0x38, 0x6b,
	0x01, 0x00, 0x00,
}
