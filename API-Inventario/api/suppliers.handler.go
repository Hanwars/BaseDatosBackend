package api

import (
	"net/http"
	"rest/repository"

	"github.com/gin-gonic/gin"
)

type createSupplierRequest struct {
	CompanyName  string  `json:"company_name" binding:"required"`
	CompanyPhone int32   `json:"company_phone" binding:"required"`
	CompanyEmail string  `json:"company_email" binding:"required"`
	ContactName  *string `json:"contact_name"`
	ContactPhone *int32  `json:"contact_phone"`
	ContactEmail *string `json:"contact_email"`
}

type updateSupplierRequest struct {
	CompanyName  string  `json:"company_name" binding:"required"`
	CompanyPhone int32   `json:"company_phone" binding:"required"`
	CompanyEmail string  `json:"company_email" binding:"required"`
	ContactName  *string `json:"contact_name"`
	ContactPhone *int32  `json:"contact_phone"`
	ContactEmail *string `json:"contact_email"`
	SupplierID   int32   `json:"supplier_id" binding:"required"`
}

type deleteSupplierRequestJSON struct {
	SupplierID int32 `json:"supplier_id" binding:"required"`
}

type deleteSupplierRequestURI struct {
	SupplierID int32 `uri:"supplier_id" binding:"required"`
}

type getSupplierRequest struct {
	SupplierID int32 `uri:"supplier_id" binding:"required"`
}

// createSupplier es el handler para crear un nuevo proveedor.
func (server *Server) createSupplier(ctx *gin.Context) {

	var req createSupplierRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	supplierRepository := repository.NewSupplierRepository(server.conn)

	supplier, err := supplierRepository.CrearProveedor(
		ctx,
		req.CompanyName,
		req.CompanyPhone,
		req.CompanyEmail,
		req.ContactName,
		req.ContactPhone,
		req.ContactEmail,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"generated_id": supplier.SupplierID,
	})
}

// updateSupplier es el handler para actualizar un proveedor.
func (server *Server) updateSupplier(ctx *gin.Context) {

	var req updateSupplierRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	supplierRepository := repository.NewSupplierRepository(server.conn)

	_, err := supplierRepository.ActualizarProveedor(
		ctx,
		req.SupplierID,
		req.CompanyName,
		req.CompanyPhone,
		req.CompanyEmail,
		req.ContactName,
		req.ContactPhone,
		req.ContactEmail,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Proveedor actualizado correctamente",
	})
}

// deleteSupplierJSON elimina un proveedor por JSON.
func (server *Server) deleteSupplierJSON(ctx *gin.Context) {

	var req deleteSupplierRequestJSON

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	supplierRepository := repository.NewSupplierRepository(server.conn)

	err := supplierRepository.EliminarProveedor(
		ctx,
		req.SupplierID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Proveedor eliminado correctamente",
	})
}

// deleteSupplierURI elimina un proveedor por URI.
func (server *Server) deleteSupplierURI(ctx *gin.Context) {

	var req deleteSupplierRequestURI

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	supplierRepository := repository.NewSupplierRepository(server.conn)

	err := supplierRepository.EliminarProveedor(
		ctx,
		req.SupplierID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Proveedor eliminado correctamente",
	})
}

// getAllSuppliers es el handler para obtener todos los proveedores.
func (server *Server) getAllSuppliers(ctx *gin.Context) {

	supplierRepository := repository.NewSupplierRepository(server.conn)

	suppliers, err := supplierRepository.ObtenerProveedores(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, suppliers)
}

// getSuppliersByID es el handler para obtener un proveedor por su ID.
func (server *Server) getSuppliersByID(ctx *gin.Context) {

	var req getSupplierRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	supplierRepository := repository.NewSupplierRepository(server.conn)

	supplier, err := supplierRepository.ObtenerProveedor(
		ctx,
		req.SupplierID,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, supplier)
}
