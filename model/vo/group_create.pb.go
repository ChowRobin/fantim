// Code generated by protoc-gen-go.
// source: group_create.proto
// DO NOT EDIT!

package vo

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type GroupCreateRequest struct {
	Name        string  `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Avatar      string  `protobuf:"bytes,2,opt,name=avatar" json:"avatar,omitempty"`
	Description string  `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
	Members     []int64 `protobuf:"varint,4,rep,packed,name=members" json:"members,omitempty"`
}

func (m *GroupCreateRequest) Reset()         { *m = GroupCreateRequest{} }
func (m *GroupCreateRequest) String() string { return proto.CompactTextString(m) }
func (*GroupCreateRequest) ProtoMessage()    {}

func (m *GroupCreateRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GroupCreateRequest) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

func (m *GroupCreateRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *GroupCreateRequest) GetMembers() []int64 {
	if m != nil {
		return m.Members
	}
	return nil
}

type GroupCreateResponse struct {
	StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode" json:"status_code,omitempty"`
	StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg" json:"status_msg,omitempty"`
	GroupId    int64  `protobuf:"varint,3,opt,name=group_id,json=groupId" json:"group_id,omitempty"`
}

func (m *GroupCreateResponse) Reset()         { *m = GroupCreateResponse{} }
func (m *GroupCreateResponse) String() string { return proto.CompactTextString(m) }
func (*GroupCreateResponse) ProtoMessage()    {}

func (m *GroupCreateResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *GroupCreateResponse) GetStatusMsg() string {
	if m != nil {
		return m.StatusMsg
	}
	return ""
}

func (m *GroupCreateResponse) GetGroupId() int64 {
	if m != nil {
		return m.GroupId
	}
	return 0
}
