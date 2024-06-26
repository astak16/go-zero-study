package handler

import (
	"net/http"

	"go-zero-study/api/internal/logic"
	"go-zero-study/api/internal/svc"
	"go-zero-study/api/internal/types"
	"go-zero-study/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func createUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCreateUserLogic(r.Context(), svcCtx)
		resp, err := l.CreateUser(&req)
		// if err != nil {
		// 	httpx.ErrorCtx(r.Context(), w, err)
		// } else {
		// 	httpx.OkJsonCtx(r.Context(), w, resp)
		// }
		response.Response(r, w, resp, err)
	}
}
