package simplekv

import "net/http"

var (
	store = make(map[string][]byte)
)

func (ap *ApiServer) KVkeys(w http.ResponseWriter, r *http.Request) {
	slog.Println("Keys: req ", r)
}

func (ap *ApiServer) KVrequest(w http.ResponseWriter, r *http.Request) {
	slog.Println("Request: req ", r)
}
