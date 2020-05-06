package error

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrNoUser    = Error{500, "用户不存在"}
	ErrUserExist = Error{501, "用户已存在"}

	ErrInputValid    = Error{520, "输入验证失败"}
	ErrInputNotFound = Error{521, "未找到对应输入"}

	ErrFailTokenGen = Error{540, "凭据创建失败"}

	ErrDBRead  = Error{550, "数据读取失败"}
	ErrDBWrite = Error{551, "数据写入失败"}
)
