package handler

import (
	"net/http"

	"gozeroImSystem/imrpc/internal/logic"
	"gozeroImSystem/imrpc/internal/svc"
	"gozeroImSystem/imrpc/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ImrpcHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewImrpcLogic(r.Context(), svcCtx)
		resp, err := l.Imrpc(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
