package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/apex/log"
	"github.com/gorilla/mux"
)

var (
	Mstore = make(map[string][]byte)
)

func (ap *Server) KVkeys(w http.ResponseWriter, r *http.Request) {
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
		log.Infof("Keys Error; %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.Header().Set("Content-type", "application/json")
		w.Write(jsondata)
	}

}

func (ap *Server) KVrequest(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		jsondata []byte
	)
	log.Debugf("Request: req %#v", r)
	vars := mux.Vars(r)
	keyID := vars["key"]

	if len(keyID) < 1 {
		http.Error(w, "Key name missing", http.StatusBadRequest)
		return
	}

	log.Debugf("KeyID; %s", keyID)

	switch r.Method {
	case "GET":
		var val []byte
		val, err = ap.be.Get(keyID)
		if err == nil {
			err = json.Unmarshal(val, &jsondata)
		}

	case "DELETE":
		err = ap.be.Delete(keyID)

	case "PUT", "POST":
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := "Unable to read url body"
			log.Infof("%s: %s", msg, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		payload, err := json.Marshal(bs)
		if err != nil {
			msg := "Unable to jsonify body"
			log.Infof("%s: %s", msg, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = ap.be.Put(keyID, payload)

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
	log.Debugf("Sending: %#v", jsondata)
}
