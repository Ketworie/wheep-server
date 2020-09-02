package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
	"wheep-server/chat"
	"wheep-server/hub"
	"wheep-server/notebook"
	"wheep-server/security"
	"wheep-server/storage"
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
	jsonServer.HandleFuncAuthorized("/user", user.HandleGet).Methods("GET")
	jsonServer.HandleFuncAuthorized("/user/list", user.HandleGetList).Methods("POST")
	jsonServer.HandleFuncAuthorized("/me", security.HandleMe).Methods("GET")
	jsonServer.HandleFuncAuthorized("/me/hubs", hub.HandleFindMyHubs).Methods("GET")
	jsonServer.HandleFuncAuthorized("/avatar/update", user.HandleUpdateAvatar).Methods("POST")
	jsonServer.HandleFuncAuthorized("/contact/list", notebook.HandleGetContacts).Methods("GET")
	jsonServer.HandleFuncAuthorized("/contact/add", notebook.HandleAddContact).Methods("GET")
	jsonServer.HandleFuncAuthorized("/contact/remove", notebook.HandleRemoveContact).Methods("GET")
	jsonServer.HandleFuncAuthorized("/hub", hub.HandleAdd).Methods("POST")
	jsonServer.HandleFuncAuthorized("/hub", hub.HandleGet).Methods("GET")
	jsonServer.HandleFuncAuthorized("/hub", hub.HandleDelete).Methods("DELETE")
	jsonServer.HandleFuncAuthorized("/hub/send", chat.HandleSend).Methods("POST")
	jsonServer.HandleFuncAuthorized("/hub/last", chat.HandleLast).Methods("GET")
	jsonServer.HandleFuncAuthorized("/hub/prev", chat.HandlePrev).Methods("GET")
	jsonServer.HandleFuncAuthorized("/hub/next", chat.HandleNext).Methods("GET")
	jsonServer.HandleFuncAuthorized("/hub/rename", hub.HandleRename).Methods("POST")
	jsonServer.HandleFuncAuthorized("/hub/users/add", hub.HandleAddUsers).Methods("POST")
	jsonServer.HandleFuncAuthorized("/hub/users/remove", hub.HandleRemoveUsers).Methods("POST")
	jsonServer.HandleFuncAuthorized("/upload/image", storage.HandleUploadImage).Methods("POST")
	server.PathPrefix("/wayne/{?:\\w{24}}/{?:[\\w\\.]+}").Handler(http.StripPrefix("/wayne/", http.FileServer(http.Dir(storage.ResourceRoot))))
	return http.ListenAndServe(":8080", server)
}

func jsonMW(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}
