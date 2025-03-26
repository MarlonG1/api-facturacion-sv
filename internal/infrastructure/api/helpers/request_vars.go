package helpers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func GetRequestVar(r *http.Request, key string) string {
	vars := mux.Vars(r)
	return vars[key]
}
