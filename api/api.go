package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/airbrake/gobrake/v4"
	"github.com/apex/log"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rbg/simplekv/airbrake"
	"github.com/rbg/simplekv/store"
	"github.com/spf13/viper"
)

const (
	apiEp = ":7800"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Infof("404 NOT FOUND: %s", r.RequestURI)
	log.Infof("VARS: %+#v", vars)
	http.Error(w, "No matching request found", http.StatusNotFound)
	n := gobrake.NewNotice("Not Found", r, 0)
	airbrake.GoBrake.Notify(n, nil)
}

// Server is the Rest instance
type Server struct {
	addr string
	s    *negroni.Negroni
	be   store.Store
}

//Run startup and serves
func (ap *Server) Run() {

	defer airbrake.GoBrake.Close()
	defer airbrake.GoBrake.NotifyOnPanic()

	notice := gobrake.NewNotice("skv starts", nil, 0)
	airbrake.GoBrake.Notify(notice, nil)

	log.Infof("Starting: %s", ap.addr)

	http.ListenAndServe(ap.addr, ap.s)

	log.Infof("exit run: %s", ap.addr)

}

type abResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newABResponseWriter(w http.ResponseWriter) *abResponseWriter {
	return &abResponseWriter{w, http.StatusOK}
}

func (a *abResponseWriter) WriteHeader(code int) {
	a.statusCode = code
	a.ResponseWriter.WriteHeader(code)
}

func perf(route string, h http.HandlerFunc) (string, http.HandlerFunc) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx, routeMetric := gobrake.NewRouteMetric(ctx, req.Method, route)
		arw := newABResponseWriter(w)

		h.ServeHTTP(arw, req)

		routeMetric.StatusCode = arw.statusCode
		airbrake.GoBrake.Routes.Notify(ctx, routeMetric) // Stops the timing and reports
		log.Infof("code: %v, method: %v, route: %v", arw.statusCode, req.Method, route)
	})

	return route, handler
}

//New will generate a new server instance for the given store back-end
func New(be store.Store, ep string) *Server {
	// validate required arguments:
	if len(viper.GetString("ab_key")) == 0 {
		fmt.Printf("Airbrake 'ab_key' is missing")
		os.Exit(1)
	}
	log.Infof("Using Airbrake:")
	log.Infof("\tkey\t%s", viper.GetString("ab_key"))

	if len(viper.GetString("ab_env")) == 0 {
		fmt.Printf("Airbrake 'ab_env' is missing")
		os.Exit(1)
	}
	log.Infof("\tenv\t%s", viper.GetString("ab_env"))

	if len(viper.GetString("ab_url")) == 0 {
		fmt.Printf("Airbrake 'ab_url' is missing")
		os.Exit(1)
	}
	log.Infof("\turl\t%s", viper.GetString("ab_url"))
	if viper.GetInt64("ab_proj") == 0 {
		fmt.Printf("Airbrake 'ab_poj' is missing")
		os.Exit(1)
	}
	log.Infof("\tproj\t%v", viper.GetInt64("ab_proj"))
	airbrake.Init()
	ap := &Server{
		be:   be,
		addr: apiEp,
	}

	if len(ep) == 0 {
		ap.addr = ep
	}

	restapi := mux.NewRouter()
	restapi.NotFoundHandler = http.HandlerFunc(notFound)
	restapi.HandleFunc(perf("/api/kv/keys/", ap.KVkeys)).Methods("GET")
	restapi.HandleFunc(perf("/api/kv/keys/{key}", ap.KVrequest)).Methods("GET", "PUT", "POST", "DELETE")

	ap.s = negroni.New()
	ap.s.UseHandler(restapi)
	return ap
}
