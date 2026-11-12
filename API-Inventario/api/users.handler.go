package api

import (
	"database/sql"
	"net/http"
	"rest/repository"
	"rest/security"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
)

// Cosas del Login
type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	AccessToken string  `json:"access_token"`
	Payload     payload `json:"payload"`
}
type payload struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Role int32  `json:"role"`
}

// Requests para usuarios
type createUserRequest struct {
	RoleID   int32   `json:"role_id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Email    *string `json:"email"`
}

type updateUserRequest struct {
	UserID   int32   `json:"user_id" binding:"required"`
	RoleID   int32   `json:"role_id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
}

// UserID se ignora: el perfil que se modifica siempre es el del token.
// CurrentPassword y NewPassword son opcionales (solo para cambiar contraseña).
type updateProfileRequest struct {
	UserID          int32   `json:"user_id"`
	Name            string  `json:"name" binding:"required"`
	Email           *string `json:"email"`
	CurrentPassword *string `json:"current_password"`
	NewPassword     *string `json:"new_password"`
}

type profileResponse struct {
	UserID int32  `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type updateUserRoleRequest struct {
	UserID int32 `uri:"user_id" binding:"required"`
	RoleID int32 `uri:"role_id" binding:"required"`
}

type deleteUserRequest struct {
	UserID int32 `json:"user_id" binding:"required"`
}

type deleteUserRequestURI struct {
	UserID int32 `uri:"user_id" binding:"required"`
}

type getUserByIDRequest struct {
	UserID int32 `uri:"user_id" binding:"required"`
}

/*Funcion de Login*/
func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.userRepository.Login(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if user.Password != req.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Usuario no autorizado"})
		return
	}
	//GENERACION DEL TOKEN
	accessToken, err := server.tokenBuilder.CreateToken(user.UserID, user.RoleID, user.Name, time.Hour*1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := LoginResponse{
		AccessToken: accessToken,
		Payload: payload{
			ID:   user.UserID,
			Name: user.Name,
			Role: user.RoleID,
		},
	}
	ctx.JSON(http.StatusOK, resp)

}

// getProfile obtiene el perfil del usuario autenticado (id tomado del token).
func (server *Server) getProfile(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)

	user, err := server.userRepository.ObtenerPerfil(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, profileResponse{
		UserID: user.UserID,
		Name:   user.Name,
		Email:  user.Email,
	})
}

// updateProfile actualiza nombre, correo y (opcional) contraseña del usuario
// autenticado. Para cambiar la contraseña se exige la actual.
func (server *Server) updateProfile(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)

	var req updateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req.UserID = authPayload.UserID

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = repository.NewUserRepository(tx).ActualizarPerfil(
		ctx, req.UserID, req.Name, req.Email, req.CurrentPassword, req.NewPassword)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Perfil actualizado correctamente",
	})
}

func (server *Server) createUser(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Solo administradores"})
		return
	}

	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	userID, err := repository.NewUserRepository(tx).CrearUsuario(
		ctx, req.RoleID, req.Name, req.Password, req.Email)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"generated_id": userID})
}

func (server *Server) updateUser(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Solo administradores"})
		return
	}

	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.Password != nil && strings.TrimSpace(*req.Password) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password no puede estar vacío; omítalo para conservar la contraseña"})
		return
	}
	if _, err := server.userRepository.ObtenerUsuario(ctx, req.UserID); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	updated, err := repository.NewUserRepository(tx).ActualizarUsuario(
		ctx, req.UserID, req.RoleID, req.Name, req.Password, req.Email)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

func (server *Server) deleteUserJSON(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	server.deleteUser(ctx, req.UserID)
}

func (server *Server) deleteUserURI(ctx *gin.Context) {
	var req deleteUserRequestURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	server.deleteUser(ctx, req.UserID)
}

func (server *Server) deleteUser(ctx *gin.Context, userID int32) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Solo administradores"})
		return
	}
	if userID == authPayload.UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No puede eliminar su propio usuario"})
		return
	}
	if _, err := server.userRepository.ObtenerUsuario(ctx, userID); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = repository.NewUserRepository(tx).EliminarUsuario(ctx, userID)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}

func (server *Server) getAllUsers(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Solo administradores"})
		return
	}
	users, err := server.userRepository.ObtenerUsuarios(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (server *Server) getUserByID(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Solo administradores"})
		return
	}
	var req getUserByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.userRepository.ObtenerUsuario(ctx, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	ctx.JSON(http.StatusOK, user)
}
