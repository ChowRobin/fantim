package status

import (
	"encoding/json"
)

type ErrStatus struct {
	Code int32
	Msg  string
}

var (
	Success               = &ErrStatus{Code: 0, Msg: "success"}
	ErrInvalidParam       = &ErrStatus{Code: 5, Msg: "invalid param"}
	ErrServiceInternal    = &ErrStatus{Code: 4, Msg: "service internal error"}
	ErrDuplicateUserId    = &ErrStatus{Code: 6, Msg: "duplicate user_id error"}
	ErrUserNotLogin       = &ErrStatus{Code: 7, Msg: "user not login"}
	ErrInvalidPassword    = &ErrStatus{Code: 8, Msg: "invalid password"}
	ErrInvalidApplyType   = &ErrStatus{Code: 9, Msg: "invalid apply type"}
	ErrInvalidApplyStatus = &ErrStatus{Code: 10, Msg: "invalid apply status"}
	ErrInvalidPageParam   = &ErrStatus{Code: 11, Msg: "invalid page param"}
)

func FillResp(resp interface{}, status *ErrStatus) map[string]interface{} {
	resultData, _ := json.Marshal(resp)
	resultMap := make(map[string]interface{})
	_ = json.Unmarshal(resultData, &resultMap)
	resultMap["status_code"] = status.Code
	resultMap["status_msg"] = status.Msg
	return resultMap
}
