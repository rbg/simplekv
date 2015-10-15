package simplekv

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
)

var (
	Mstore = make(map[string][]byte)
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
	var (
		err      error
		jsondata []byte
	)
	slog.Printf("Request: req %#v", r)
	vars := mux.Vars(r)
	keyID := vars["key"]

	if len(keyID) < 1 {
		http.Error(w, "Key name missing", http.StatusBadRequest)
		return
	}

	slog.Println("Key is: ", keyID)

	switch r.Method {
	case "GET":
		var val []byte
		val, err = ap.be.Get(keyID)
		if err == nil {
			err = json.Unmarshal(val, &jsondata)
		}

	case "DELETE":
		err = ap.be.Delete(keyID)

	case "PUT":
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := "Unable to read url body"
			slog.ErrPrintf("%s: %s", msg, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if payload, err := json.Marshal(bs); err != nil {
			msg := "Unable to jsonify body"
			slog.ErrPrintf("%s: %s", msg, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			err = ap.be.Put(keyID, payload)
		}

	default:
		http.Error(w, "Method not supported", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(jsondata)
	w.Write([]byte("\n"))
	slog.Printf("Sending: %#v", jsondata)
}
