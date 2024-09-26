package handler

import (
	"net/http"

	"gozeroImSystem/imapi/internal/logic"
	"gozeroImSystem/imapi/internal/svc"
	"gozeroImSystem/imapi/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ImapiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewImapiLogic(r.Context(), svcCtx)
		resp, err := l.Imapi(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
