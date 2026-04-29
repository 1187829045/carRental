package handler

// DataGridView mirrors com.yeqifu.sys.utils.DataGridView (LayUI table response).
// code default 0, msg default "".
type DataGridView struct {
	Code  int64 `json:"code"`
	Msg   string `json:"msg"`
	Count int64 `json:"count,omitempty"`
	Data  any   `json:"data,omitempty"`
}

func NewData(data any) DataGridView {
	return DataGridView{Code: 0, Msg: "", Data: data}
}

func NewPage(count int64, data any) DataGridView {
	return DataGridView{Code: 0, Msg: "", Count: count, Data: data}
}

// ResultObj mirrors com.yeqifu.sys.utils.ResultObj (operation response).
type ResultObj struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

var (
	AddSuccess = ResultObj{Code: 0, Msg: "添加成功"}
	AddError   = ResultObj{Code: -1, Msg: "添加失败"}

	UpdateSuccess = ResultObj{Code: 0, Msg: "更新成功"}
	UpdateError   = ResultObj{Code: -1, Msg: "更新失败"}

	DeleteSuccess = ResultObj{Code: 0, Msg: "删除成功"}
	DeleteError   = ResultObj{Code: -1, Msg: "删除失败"}

	ResetSuccess = ResultObj{Code: 0, Msg: "重置成功"}
	ResetError   = ResultObj{Code: -1, Msg: "重置失败"}

	DispatchSuccess = ResultObj{Code: 0, Msg: "分配成功"}
	DispatchError   = ResultObj{Code: -1, Msg: "分配失败"}

	StatusTrue  = ResultObj{Code: 0}
	StatusFalse = ResultObj{Code: -1}

	AddSuccessRent = ResultObj{Code: 0, Msg: "添加出租单成功，请等待审核！"}
	AddErrorRent   = ResultObj{Code: -1, Msg: "添加出租单失败"}
	CheckSuccess   = ResultObj{Code: 0, Msg: "审核成功"}
	CheckError     = ResultObj{Code: 0, Msg: "审核失败"}
)

