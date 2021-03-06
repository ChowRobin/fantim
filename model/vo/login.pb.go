// Code generated by protoc-gen-go.
// source: login.proto
// DO NOT EDIT!

package vo

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type LoginRequest struct {
	UserId   int64  `protobuf:"varint,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *LoginRequest) Reset()         { *m = LoginRequest{} }
func (m *LoginRequest) String() string { return proto.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()    {}

func (m *LoginRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *LoginRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type LoginResponse struct {
	StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode" json:"status_code,omitempty"`
	StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg" json:"status_msg,omitempty"`
	SessionId  string `protobuf:"bytes,3,opt,name=session_id,json=sessionId" json:"session_id,omitempty"`
	UserInfo   *User  `protobuf:"bytes,4,opt,name=user_info,json=userInfo" json:"user_info,omitempty"`
}

func (m *LoginResponse) Reset()         { *m = LoginResponse{} }
func (m *LoginResponse) String() string { return proto.CompactTextString(m) }
func (*LoginResponse) ProtoMessage()    {}

func (m *LoginResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *LoginResponse) GetStatusMsg() string {
	if m != nil {
		return m.StatusMsg
	}
	return ""
}

func (m *LoginResponse) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *LoginResponse) GetUserInfo() *User {
	if m != nil {
		return m.UserInfo
	}
	return nil
}
