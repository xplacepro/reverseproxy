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
	var User = flag.String("user", "", "Basic auth user")
	var Password = flag.String("password", "", "Basic auth password")

	var ClientUser = flag.String("client-user", "", "Basic auth user for dashboard client")
	var ClientPassword = flag.String("client-password", "", "Basic auth user password for dashboard client")

	var Listen = flag.String("listen", ":8080", "Interface and port to listen, :8080")

	flag.Parse()

	if *User == "" || *Password == "" {
		Usage()
	}

	env := &rpc.Env{Auth: rpc.BasicAuthorization{*User, *Password}, ClientUser: *ClientUser, ClientPassword: *ClientPassword}
	r := mux.NewRouter()
	r.StrictSlash(false)
	r.Handle("/api/v1/domains", rpc.Handler{env, controllers.PostListDomainHandler}).Methods("POST")
	r.Handle("/api/v1/domains/{domain:[a-zA-Z0-9-.]+}", rpc.Handler{env, controllers.DeleteDomainHandler}).Methods("DELETE")
	http.Handle("/", r)
	log.Printf("Started server on %s", *Listen)
	http.ListenAndServe(*Listen, nil)
}
