package api

import (
	"net/http"

	"github.com/apex/log"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rbg/simplekv/store"
)

const (
	apiEp = ":7800"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Infof("404 NOT FOUND: %s", r.RequestURI)
	log.Infof("VARS: %+#v", vars)
	http.Error(w, "No matching request found", http.StatusNotFound)
}

// Server is the Rest instance
type Server struct {
	addr string
	s    *negroni.Negroni
	be   store.Store
}

//Run startup and serves
func (ap *Server) Run() {

	log.Infof("Starting: %s", ap.addr)

	http.ListenAndServe(ap.addr, ap.s)

	log.Infof("exit run: %s", ap.addr)

}

//New will generate a new server instance for the given store back-end
func New(be store.Store, ep string) *Server {
	ap := &Server{
		be:   be,
		addr: apiEp,
	}

	if len(ep) == 0 {
		ap.addr = ep
	}

	restapi := mux.NewRouter()
	restapi.NotFoundHandler = http.HandlerFunc(notFound)
	restapi.HandleFunc("/api/kv/keys/", ap.KVkeys).Methods("GET")
	restapi.HandleFunc("/api/kv/keys/{key}", ap.KVrequest).Methods("GET", "PUT", "POST", "DELETE")

	ap.s = negroni.New()
	ap.s.UseHandler(restapi)
	return ap
}
