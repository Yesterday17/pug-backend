package error

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	NoError = Error{200, ""}

	ErrNoUser    = Error{500, "用户不存在"}
	ErrUserExist = Error{501, "用户已存在"}

	ErrInputInvalid  = Error{520, "输入验证失败"}
	ErrInputNotFound = Error{521, "未找到对应输入"}

	ErrFailTokenGen = Error{540, "凭据创建失败"}

	ErrDBRead   = Error{550, "数据读取失败"}
	ErrDBWrite  = Error{551, "数据写入失败"}
	ErrDBDelete = Error{552, "数据删除失败"}

	ErrModuleNotFound          = Error{560, "模块不存在"}
	ErrPipeNotFound            = Error{561, "管道不存在"}
	ErrConstructedPipeNotFound = Error{562, "构建管道不存在"}
	ErrCannotRestrictSelf      = Error{563, "限制权限不可超过自身权限等级"}

	ErrPermissionDeny           = Error{600, "权限不足"}
	ErrCannotVisitPrivate       = Error{601, "无法访问他人 Private 内容"}
	ErrCannotDeleteNotOwnedPipe = Error{602, "无法删除他人的管道构建"}
)

func (e *Error) Error() string {
	return e.Message
}
