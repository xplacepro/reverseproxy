package controllers

import (
	"github.com/gorilla/mux"
	"github.com/xplacepro/reverseproxy/nginx"
	"github.com/xplacepro/rpc"
	"net/http"
)

func DeleteDomainHandler(env *rpc.Env, w http.ResponseWriter, r *http.Request) rpc.Response {
	vars := mux.Vars(r)
	domain_name := vars["domain"]

	domain := nginx.Domain{Domain: domain_name}

	if err := domain.Delete(); err != nil {
		return rpc.BadRequest(err)
	}

	return rpc.SyncResponse(nil)
}
