package httphandler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Health() func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}
}
