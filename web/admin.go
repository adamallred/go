package web

import (
	gocontext "context"

	"net/http"
	"time"

	"github.com/kellegous/go/context"
)

type adminHandler struct {
	ctx *context.Context
}

func adminGet(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	p := parseName("/admin/", r.URL.Path)

	if p == "" {
		writeJSONOk(w)
		return
	}

	if p == "dumps" {
		goctx, cancel := gocontext.WithTimeout(gocontext.Background(), time.Minute)
		defer cancel()

		if golinks, err := ctx.GetAll(goctx); err != nil {
			writeJSONBackendError(w, err)
			return
		} else {
			writeJSON(w, golinks, http.StatusOK)
		}
	}

}

func (h *adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		adminGet(h.ctx, w, r)
	default:
		writeJSONError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusOK) // fix
	}
}
