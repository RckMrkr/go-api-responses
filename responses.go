package responses

import (
	"encoding/json"
	"net/http"
)

func respond(w http.ResponseWriter, s interface{}, code int) {
	w.WriteHeader(code)
	if s != nil {
		json.NewEncoder(w).Encode(s)
	}
}

func respondError(w http.ResponseWriter, m string, code int) {
	s := struct {
		Error string `json:"error"`
	}{
		Error: m,
	}

	respond(w, s, code)
}

// RespondSuccess marshals s to json (if not nil) and set status code to 200
func RespondSuccess(w http.ResponseWriter, s interface{}) {
	respond(w, s, http.StatusOK)
}

// RespondNotFound marshals s to json (if not nil) and set status code to 404
func RespondNotFound(w http.ResponseWriter) {
	respond(w, nil, http.StatusNotFound)
}

// RespondUnauthorizedBearerJWT inserts www-authenticate header and sets status code to 401
func RespondUnauthorizedBearerJWT(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer token_type="JWT"`)
	respondError(w, "You are not authorized", http.StatusUnauthorized)
}

// RespondInternalError responds with the string in envelope sets status code to 500
func RespondInternalError(w http.ResponseWriter, m string) {
	respondError(w, m, http.StatusInternalServerError)
}

// RespondBadRequest responds with the string in envelope sets status code to 400
func RespondBadRequest(w http.ResponseWriter, m string) {
	respondError(w, m, http.StatusBadRequest)
}
