package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Jonss/book-lending/app/usecases"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserRestHandler struct {
	findUserUsecase usecases.FindUserUsecase
}

func NewUserRestHandler(findUserUsecase usecases.FindUserUsecase) UserRestHandler {
	return UserRestHandler{findUserUsecase}
}

func (h UserRestHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loggedId := vars["logged_user_id"]

	loggedUserID := uuid.MustParse(loggedId)
	user, err := h.findUserUsecase.FindUserByID(loggedUserID)
	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
	} else {
		writeResponse(w, http.StatusOK, user)
	}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
