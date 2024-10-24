package httphandler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Health(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
