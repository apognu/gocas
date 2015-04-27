package cas

import "github.com/gorilla/mux"

func New(r *mux.Router) {
	r.HandleFunc("/login", loginRequestorHandler).Methods("GET")
	r.HandleFunc("/login", loginAcceptorHandler).Methods("POST")
}
