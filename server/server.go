package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"wheep-server/hub"
	"wheep-server/security"
	"wheep-server/user"
)

type Server struct {
	*mux.Router
}

func (s *Server) HandleFunc(path string, f func(http.ResponseWriter, *http.Request) error) *mux.Route {
	return s.Router.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		err := f(writer, request)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			jsonErr := json.NewEncoder(writer).Encode(struct {
				Message string    `json:"message"`
				Date    time.Time `json:"date"`
				Path    string    `json:"path"`
			}{
				Message: err.Error(),
				Date:    time.Now(),
				Path:    path,
			})
			if jsonErr != nil {
				log.Println(jsonErr)
			}
		}
	})
}

func (s *Server) HandleFuncAuthorized(path string, f func(user.Model, http.ResponseWriter, *http.Request) error) *mux.Route {
	return s.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) error {
		u, err := security.HandleAuthorize(writer, request)
		if err != nil {
			return err
		}
		return f(u, writer, request)
	})
}

func StartServer() error {
	server := Server{mux.NewRouter()}
	server.Use(jsonMW)
	server.HandleFunc("/login", security.HandleLogin).Methods("POST")
	server.HandleFunc("/user", user.HandleAdd).Methods("POST")
	server.HandleFuncAuthorized("/createIndexes", security.HandleCreateIndexes).Methods("GET")
	server.HandleFuncAuthorized("/user/me", security.HandleMe).Methods("GET")
	server.HandleFuncAuthorized("/user/me/hubs", hub.HandleFindMyHubs).Methods("GET")
	server.HandleFuncAuthorized("/hub", hub.HandleAdd).Methods("POST")
	server.HandleFuncAuthorized("/hub", hub.HandleGet).Methods("GET")
	server.HandleFuncAuthorized("/hub", hub.HandleDelete).Methods("DELETE")
	server.HandleFuncAuthorized("/hub/rename", hub.HandleRename).Methods("POST")
	server.HandleFuncAuthorized("/hub/users/add", hub.HandleAddUsers).Methods("POST")
	server.HandleFuncAuthorized("/hub/users/remove", hub.HandleRemoveUsers).Methods("POST")
	return http.ListenAndServe(":8080", server)
}

func jsonMW(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}
