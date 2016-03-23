package controllers

import (
	"encoding/json"
	"errors"
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

func ValidatePostListDomain(c CreateDomainParams) map[string]interface{} {
	validationErrors := make(map[string]interface{})
	if strings.Trim(c.Domain, " ") == "" {
		validationErrors["domain"] = "domain is required"
	}
	if strings.Trim(c.Config, " ") == "" {
		validationErrors["config"] = "config is required"
	}
	return validationErrors
}

func PostListDomainHandler(env *rpc.Env, w http.ResponseWriter, r *http.Request) (rpc.Response, int, error) {
	var create_params CreateDomainParams

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, rpc.StatusError{Err: err}
	}

	err = json.Unmarshal(data, &create_params)
	if err != nil {
		return nil, http.StatusBadRequest, rpc.StatusError{Err: err}
	}

	validation_errors := ValidatePostListDomain(create_params)
	if len(validation_errors) > 0 {
		return nil, http.StatusBadRequest, rpc.StatusError{Err: errors.New("Validation error"), MetadataMap: validation_errors}
	}

	domain := nginx.Domain{Domain: create_params.Domain, Config: create_params.Config}
	if err := domain.Create(); err != nil {
		return nil, http.StatusInternalServerError, rpc.StatusError{Err: err}
	}

	if err := nginx.Test(); err != nil {
		log.Printf("Error testing nginx configuration for domain %s, %s", domain.Domain, err)
		return nil, http.StatusInternalServerError, rpc.StatusError{Err: err}
	}

	if err := nginx.Reload(); err != nil {
		return nil, http.StatusInternalServerError, rpc.StatusError{Err: err}
	}

	return rpc.SyncResponse{"Success", http.StatusOK, map[string]interface{}{}}, http.StatusOK, nil

}
