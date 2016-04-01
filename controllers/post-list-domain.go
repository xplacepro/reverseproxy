package controllers

import (
	"encoding/json"
	"github.com/xplacepro/reverseproxy/nginx"
	"github.com/xplacepro/rpc"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type CreateDomainParams struct {
	Domain string
	Config string
}

func ValidatePostListDomain(c CreateDomainParams) bool {
	if strings.Trim(c.Domain, " ") == "" {
		return false
	}
	if strings.Trim(c.Config, " ") == "" {
		return false
	}
	return true
}

func PostListDomainHandler(env *rpc.Env, w http.ResponseWriter, r *http.Request) rpc.Response {
	var create_params CreateDomainParams

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return rpc.BadRequest(err)
	}

	err = json.Unmarshal(data, &create_params)
	if err != nil {
		return rpc.BadRequest(err)
	}

	if !ValidatePostListDomain(create_params) {
		return rpc.BadRequest(ValidationError)
	}

	domain := nginx.Domain{Domain: create_params.Domain, Config: create_params.Config}
	if err := domain.Create(); err != nil {
		return rpc.InternalError(err)
	}

	if err := nginx.Test(); err != nil {
		log.Printf("Error testing nginx configuration for domain %s, %s", domain.Domain, err)
		return rpc.InternalError(err)
	}

	if err := nginx.Reload(); err != nil {
		return rpc.InternalError(err)
	}

	return rpc.SyncResponse(nil)

}
