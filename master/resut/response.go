package resut

const (
	ErrorCode   = -1
	SuccessCode = 0
)

type Rest struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Page struct {
	Total int64 `json:"total"`
	List  interface{}
}

func FAILMSG(msg string) *Rest {
	return &Rest{Code: ErrorCode, Message: msg, Data: ""}
}

func FAIL() *Rest {
	return &Rest{Code: ErrorCode, Message: "操作失败", Data: ""}
}

func SUCCESS() *Rest {
	return &Rest{Code: SuccessCode, Message: "操作成功", Data: ""}
}

func DATA(data interface{}) *Rest {
	return &Rest{Code: SuccessCode, Message: "操作成功", Data: data}
}

func PAGE(total int64, data interface{}) *Rest {
	return &Rest{Code: SuccessCode, Message: "操作成功", Data: &Page{
		Total: total,
		List:  data,
	}}
}
