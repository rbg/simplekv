package simplekv

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rbg/simplekv/store"
	"github.com/stackengine/selog"
)

const (
	api_endpoint = ":1964"
)

var (
	slog = selog.Register("simplekv", 0)
)

func notFound(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slog.Println("404 NOT FOUND: ", r.RequestURI)
	slog.Println("VARS: ", vars)
	http.Error(w, "No matching request found", http.StatusNotFound)
}

type ApiServer struct {
	addr string
	s    *negroni.Negroni
	be   store.Store
}

func (api_server *ApiServer) Run() {
	slog.Println("Starting: ", api_server.addr)
	http.ListenAndServe(api_server.addr, api_server.s)
	slog.Println("exit run: ", api_server.addr)
}

func NewServer(be store.Store) *ApiServer {
	api_server := &ApiServer{be: be}

	api := mux.NewRouter()
	api.NotFoundHandler = http.HandlerFunc(notFound)
	api.HandleFunc("/api/kv/keys/", api_server.KVkeys).Methods("GET")
	api.HandleFunc("/api/kv/{key}", api_server.KVrequest).Methods("GET", "PUT", "POST", "DELETE")
	api_server.addr = api_endpoint
	api_server.s = negroni.New()
	api_server.s.UseHandler(api)
	return api_server
}
