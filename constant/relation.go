package constant

const (
	RelationApplyTypeFriend = 1
	RelationApplyTypeGroup  = 2
)

const (
	RelationApplyStatusProcess = 0 // 处理中
	RelationApplyStatusPass    = 1 // 接收方同意
	RelationApplyStatusReject  = 2 // 接收方拒绝
	RelationApplyStatusCancel  = 3 // 发起方取消
)

const (
	RelationTypeNone   = 0 // 无关系
	RelationTypeFriend = 1 // 好友
	RelationTypeBlack  = 2 // 拉黑
)
