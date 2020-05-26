package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
	"wheep-server/hub"
	"wheep-server/security"
	"wheep-server/user"
)

var ResourceRoot = "/resources/"

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

func (s *Server) HandleFuncAuthorized(path string, f func(primitive.ObjectID, http.ResponseWriter, *http.Request) error) *mux.Route {
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
	jsonServer := Server{server.NewRoute().Subrouter()}
	jsonServer.Use(jsonMW)
	jsonServer.HandleFunc("/login", security.HandleLogin).Methods("POST")
	jsonServer.HandleFunc("/user", user.HandleAdd).Methods("POST")
	jsonServer.HandleFuncAuthorized("/createIndexes", security.HandleCreateIndexes).Methods("GET")
	jsonServer.HandleFuncAuthorized("/user/me", security.HandleMe).Methods("GET")
	jsonServer.HandleFuncAuthorized("/user/me/hubs", hub.HandleFindMyHubs).Methods("GET")
	jsonServer.HandleFuncAuthorized("/hub", hub.HandleAdd).Methods("POST")
	jsonServer.HandleFuncAuthorized("/hub", hub.HandleGet).Methods("GET")
	jsonServer.HandleFuncAuthorized("/hub", hub.HandleDelete).Methods("DELETE")
	jsonServer.HandleFuncAuthorized("/hub/rename", hub.HandleRename).Methods("POST")
	jsonServer.HandleFuncAuthorized("/hub/users/add", hub.HandleAddUsers).Methods("POST")
	jsonServer.HandleFuncAuthorized("/hub/users/remove", hub.HandleRemoveUsers).Methods("POST")
	jsonServer.HandleFuncAuthorized("/upload", HandleUpload).Methods("POST")

	server.PathPrefix("/wayne/{?:\\w{24}}/{?:[\\w\\.]+}").Handler(http.StripPrefix("/wayne/", http.FileServer(http.Dir(ResourceRoot))))
	return http.ListenAndServe(":8080", server)
}

func jsonMW(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}
