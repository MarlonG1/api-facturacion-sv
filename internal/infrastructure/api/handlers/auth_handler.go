package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MarlonG1/api-facturacion-sv/internal/application/auth"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/auth/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/user"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
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

// Login godoc
// @Summary      Login
// @Description  Login with API Key, API Secret and Hacienda Credentials
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security  	 BearerAuth
// @Param auth body models.AuthCredentials true "Auth credentials"
// @Success      200 {object} string "token"
// @Failure      400 {object} response.APIError
// @Failure      401 {object} response.APIError
// @Failure      500 {object} response.APIError
// @Router       /api/v1/auth/login [post]
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

// Register godoc
// @Summary      Register
// @Description  Register a new user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param user body user.User true "User data"
// @Success      201 {object} []user.ListBranchesResponse
// @Failure      400 {object} response.APIError
// @Failure      401 {object} response.APIError
// @Failure      500 {object} response.APIError
// @Router       /api/v1/auth/register [post]
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
