package handlers

import (
	"encoding/json"
	"github.com/MarlonG1/api-facturacion-sv/internal/application/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"net/http"
)

type AuthHandler struct {
	authUseCase *auth.AuthUseCase
	respWriter  *response.ResponseWriter
}

func NewAuthHandler(authUseCase *auth.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		respWriter:  response.NewResponseWriter(),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar la solicitud
	var req models.AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.Error(w, http.StatusBadRequest, "Invalid request format", nil)
		return
	}

	// 2. Iniciar sesión
	response, err := h.authUseCase.Login(r.Context(), &req)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Responder con éxito
	h.respWriter.Success(w, http.StatusOK, response, nil)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// 1. Decodear la solicitud
	var req user.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respWriter.Error(w, http.StatusBadRequest, "Invalid request format", nil)
		return
	}

	// 2. Registrar el usuario
	response, err := h.authUseCase.Register(r.Context(), &req)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	// 3. Responder con éxito
	h.respWriter.Success(w, http.StatusCreated, response, nil)
}
