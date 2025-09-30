package gateway

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"api-gateway/internal/logic/gateway"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"
)

func ProxyGetGoodsDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetGoodsDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := gateway.NewProxyGetGoodsDetailLogic(r.Context(), svcCtx)
		resp, err := l.ProxyGetGoodsDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
