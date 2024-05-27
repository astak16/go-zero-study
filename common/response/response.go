package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func Response(r *http.Request, w http.ResponseWriter, res any, err error) {
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, &Body{
			Code: 500,
			Data: nil,
			Msg:  err.Error(),
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, &Body{
		Code: 200,
		Data: res,
		Msg:  "成功",
	})
}
