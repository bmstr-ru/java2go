package httphandler

import (
	"encoding/json"
	java2go "github.com/bmstr-ru/java2go/go"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

func GetClientSummary(exposureService java2go.TotalExposureService) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		clientIdStr := p.ByName("clientId")
		log.Info().Str("clientId", clientIdStr).Msg("Received request")
		clientId, err := strconv.ParseInt(clientIdStr, 10, 0)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Invalid clientId"))
			return
		}
		exposure, err := exposureService.GetClientsTotalExposure(clientId)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		body, err := json.Marshal(exposure)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(body)
	}
}
