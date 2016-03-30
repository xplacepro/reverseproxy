package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/xplacepro/reverseproxy/controllers"
	"github.com/xplacepro/rpc"
	"log"
	"net/http"
	"os"
)

func Usage() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	var ConfigPath = flag.String("config", "config.ini", "Path to configuration file")
	flag.Parse()

	config, err := rpc.ParseConfiguration(*ConfigPath)

	if err != nil {
		panic(err)
	}

	env := &rpc.Env{Auth: rpc.BasicAuthorization{config["auth.user"], config["auth.password"]},
		ClientAuth: rpc.ClientBasicAuthorization{config["client_auth.user"], config["client_auth.password"]}}

	r := mux.NewRouter()
	r.StrictSlash(false)
	r.Handle("/api/v1/domains", rpc.Handler{env, controllers.PostListDomainHandler}).Methods("POST")
	r.Handle("/api/v1/domains/{domain:[a-zA-Z0-9-.]+}", rpc.Handler{env, controllers.DeleteDomainHandler}).Methods("DELETE")
	http.Handle("/", r)
	log.Printf("Started server on %s", config["listen"])
	panic(http.ListenAndServe(config["listen"], nil))
}
