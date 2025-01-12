package api

import (
	"net/http"

	"github.com/papidb/aqua/pkg/internal"
	"github.com/rs/zerolog"
)

// Success sends a JSend success message with status code 200. It logs the response
// if a zerolog.Logger is attached to the request.
func Success(r *http.Request, w http.ResponseWriter, v interface{}) {
	log := zerolog.Ctx(r.Context())
	raw := internal.ToJSON(v)

	send(w, http.StatusOK, raw)

	log.Info().
		Int("status", http.StatusOK).
		Int("length", len(raw)).
		Interface("response_headers", internal.ToLowerKeys(w.Header())).
		RawJSON("response", raw).
		Msg("")
}

// Error sends a Json error message. It logs the response if a zerolog.Logger is attached to the request.
func Error(r *http.Request, w http.ResponseWriter, err AppErr) {
	log := zerolog.Ctx(r.Context())
	raw := internal.ToJSON(err)

	send(w, err.Code, raw)

	log.Err(err).
		Int("status", err.Code).
		Int("length", len(raw)).
		Interface("response_headers", internal.ToLowerKeys(w.Header())).
		RawJSON("response", raw).
		Msg("")
}
