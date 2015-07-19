package simplekv

import (
	"encoding/json"
	"net/http"
	"sort"
)

var (
	store = make(map[string][]byte)
)

func (ap *ApiServer) KVkeys(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		jsondata []byte
	)

	keys, err := ap.be.Keys()
	if err == nil {
		sort.Strings(keys)
		jsondata, err = json.Marshal(keys)
	}

	if err != nil {
		slog.ErrPrintln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.Header().Set("Content-type", "application/json")
		w.Write(jsondata)
	}

}

func (ap *ApiServer) KVrequest(w http.ResponseWriter, r *http.Request) {
	slog.Println("Request: req ", r)
}
