package gateway

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"api-gateway/internal/logic/gateway"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"
)

func ProxyUpdateGoodsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateGoodsDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := gateway.NewProxyUpdateGoodsLogic(r.Context(), svcCtx)
		resp, err := l.ProxyUpdateGoods(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
