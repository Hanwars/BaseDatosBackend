package api

import (
	"net/http"
	"rest/repository"
	"rest/security"

	"github.com/gin-gonic/gin"
)

// Requests para clientes
type createClientRequest struct {
	Name         string `json:"name" binding:"required"`
	ContactPhone int32  `json:"contact_phone" binding:"required"`
	ContactEmail string `json:"contact_email" binding:"required"`
}

type updateClientRequest struct {
	ClientID     int32  `json:"client_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	ContactPhone int32  `json:"contact_phone" binding:"required"`
	ContactEmail string `json:"contact_email" binding:"required"`
}

type deleteClientRequestJSON struct {
	ClientID int32 `json:"client_id" binding:"required"`
}

type deleteClientRequestURI struct {
	ClientID int32 `uri:"client_id" binding:"required"`
}

type getClientByIDRequest struct {
	ClientID int32 `uri:"client_id" binding:"required"`
}

// createClient es el handler para crear un nuevo cliente.
func (server *Server) createClient(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req createClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	client, err := repository.NewClientRepository(tx).CrearCliente(
		ctx, req.Name, req.ContactPhone, req.ContactEmail)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"generated_id": client.ClientID,
	})
}

// updateClient es el handler para actualizar un cliente.
func (server *Server) updateClient(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req updateClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = repository.NewClientRepository(tx).ActualizarCliente(
		ctx, req.ClientID, req.Name, req.ContactPhone, req.ContactEmail)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cliente actualizado correctamente",
	})
}

// deleteClientJSON elimina un cliente por JSON.
func (server *Server) deleteClientJSON(ctx *gin.Context) {
	var req deleteClientRequestJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	server.deleteClient(ctx, req.ClientID)
}

// deleteClientURI elimina un cliente por URI.
func (server *Server) deleteClientURI(ctx *gin.Context) {
	var req deleteClientRequestURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	server.deleteClient(ctx, req.ClientID)
}

// deleteClient contiene la lógica común de deleteClientJSON y deleteClientURI.
func (server *Server) deleteClient(ctx *gin.Context, clientID int32) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = repository.NewClientRepository(tx).EliminarCliente(ctx, clientID)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cliente eliminado correctamente",
	})
}

// getAllClients es el handler para obtener todos los clientes.
func (server *Server) getAllClients(ctx *gin.Context) {
	clients, err := server.clientRepository.ObtenerClientes(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, clients)
}

// getClientByID es el handler para obtener un cliente por su ID.
func (server *Server) getClientByID(ctx *gin.Context) {
	var req getClientByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	client, err := server.clientRepository.ObtenerCliente(ctx, req.ClientID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, client)
}
