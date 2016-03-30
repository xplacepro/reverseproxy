package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/xplacepro/reverseproxy/controllers"
	"github.com/xplacepro/rpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func ReloadEnv(env *rpc.Env, config map[string]string) {
	env.Auth = rpc.BasicAuthorization{config["auth.user"], config["auth.password"]}
	env.ClientAuth = rpc.ClientBasicAuthorization{config["client_auth.user"], config["client_auth.password"]}
}

func main() {
	var ConfigPath = flag.String("config", "config.ini", "Path to configuration file")
	flag.Parse()

	var config map[string]string

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	env := &rpc.Env{}

	go func() {
		for sign := range c {
			rpc.ParseConfiguration(*ConfigPath, &config)
			log.Printf("Reloading configuration, %s\n", sign)
			ReloadEnv(env, config)
		}
	}()

	rpc.ParseConfiguration(*ConfigPath, &config)

	ReloadEnv(env, config)

	r := mux.NewRouter()
	r.StrictSlash(false)
	r.Handle("/api/v1/domains", rpc.Handler{env, controllers.PostListDomainHandler}).Methods("POST")
	r.Handle("/api/v1/domains/{domain:[a-zA-Z0-9-.]+}", rpc.Handler{env, controllers.DeleteDomainHandler}).Methods("DELETE")
	http.Handle("/", r)
	log.Printf("Started server on %s", config["listen"])
	panic(http.ListenAndServe(config["listen"], nil))
}
