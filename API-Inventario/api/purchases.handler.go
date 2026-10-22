package api

import (
	"net/http"
	"rest/repository"
	"rest/security"
	"time"

	"github.com/gin-gonic/gin"
)

/*
Compras: el encabezado se crea con totales en 0 y luego se le agregan
detalles. El trigger trg_DetallesCompra recalcula subtotal, descuento,
total y el stock de los insumos cada vez que cambia un detalle.
Las operaciones que mueven stock (detalles y eliminar compra) usan
beginAuditTx para que trg_AuditoriaInsumos registre al usuario.
*/

// Requests para compras
type createPurchaseRequest struct {
	SupplierID int32      `json:"supplier_id" binding:"required"`
	Date       *time.Time `json:"date"`
	IVA        float64    `json:"iva"`
}

type updatePurchaseRequest struct {
	PurchaseID int32      `json:"purchase_id" binding:"required"`
	SupplierID int32      `json:"supplier_id" binding:"required"`
	Date       *time.Time `json:"date"`
	IVA        float64    `json:"iva"`
}

type deletePurchaseRequestJSON struct {
	PurchaseID int32 `json:"purchase_id" binding:"required"`
}

type deletePurchaseRequestURI struct {
	PurchaseID int32 `uri:"purchase_id" binding:"required"`
}

type getPurchaseByIDRequest struct {
	PurchaseID int32 `uri:"purchase_id" binding:"required"`
}

// Requests para detalles de compra
type createPurchaseDetailRequest struct {
	PurchaseID int32   `json:"purchase_id" binding:"required"`
	SupplyID   int32   `json:"supply_id" binding:"required"`
	Quantity   int32   `json:"quantity" binding:"required"`
	UnitPrice  float64 `json:"unit_price" binding:"required"`
	Discount   float64 `json:"discount"`
}

type updatePurchaseDetailRequest struct {
	DetailID  int32   `json:"detail_id" binding:"required"`
	SupplyID  int32   `json:"supply_id" binding:"required"`
	Quantity  int32   `json:"quantity" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
	Discount  float64 `json:"discount"`
}

type deletePurchaseDetailRequestJSON struct {
	DetailID int32 `json:"detail_id" binding:"required"`
}

type deletePurchaseDetailRequestURI struct {
	DetailID int32 `uri:"detail_id" binding:"required"`
}

type getPurchaseDetailByIDRequest struct {
	DetailID int32 `uri:"detail_id" binding:"required"`
}

// ========================= COMPRAS =========================

// createPurchase es el handler para crear una nueva compra.
// El usuario de la compra es el del token.
func (server *Server) createPurchase(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req createPurchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.purchaseRepository.CrearCompra(
		ctx, authPayload.UserID, req.SupplierID, req.Date, req.IVA)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"generated_id": result.PurchaseID,
	})
}

// updatePurchase es el handler para actualizar una compra.
// Se conserva el usuario que registró la compra.
func (server *Server) updatePurchase(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req updatePurchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	purchase, err := server.purchaseRepository.ObtenerCompra(ctx, req.PurchaseID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	_, err = server.purchaseRepository.ActualizarCompra(
		ctx, req.PurchaseID, purchase.UserID, req.SupplierID, req.Date, req.IVA)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Compra actualizada correctamente",
	})
}

// deletePurchaseJSON elimina una compra por JSON.
func (server *Server) deletePurchaseJSON(ctx *gin.Context) {
	var req deletePurchaseRequestJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	server.deletePurchase(ctx, req.PurchaseID)
}

// deletePurchaseURI elimina una compra por URI.
func (server *Server) deletePurchaseURI(ctx *gin.Context) {
	var req deletePurchaseRequestURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	server.deletePurchase(ctx, req.PurchaseID)
}

// deletePurchase elimina la compra con sus detalles (el trigger resta el stock).
func (server *Server) deletePurchase(ctx *gin.Context, purchaseID int32) {
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

	err = repository.NewPurchaseRepository(tx).EliminarCompra(ctx, purchaseID)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Compra eliminada correctamente",
	})
}

// getAllPurchases obtiene todas las compras.
func (server *Server) getAllPurchases(ctx *gin.Context) {
	purchases, err := server.purchaseRepository.ObtenerCompras(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, purchases)
}

// getPurchaseByID obtiene una compra por su ID.
func (server *Server) getPurchaseByID(ctx *gin.Context) {
	var req getPurchaseByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	purchase, err := server.purchaseRepository.ObtenerCompra(ctx, req.PurchaseID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, purchase)
}

// ===================== DETALLES DE COMPRA =====================

// createPurchaseDetail agrega un insumo a una compra.
func (server *Server) createPurchaseDetail(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req createPurchaseDetailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	detail, err := repository.NewPurchaseDetailRepository(tx).CrearDetalleCompra(
		ctx, req.PurchaseID, req.SupplyID, req.Quantity, req.UnitPrice, req.Discount)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"generated_id": detail.DetailID,
	})
}

// updatePurchaseDetail actualiza un detalle de compra.
func (server *Server) updatePurchaseDetail(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req updatePurchaseDetailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = repository.NewPurchaseDetailRepository(tx).ActualizarDetalleCompra(
		ctx, req.DetailID, req.SupplyID, req.Quantity, req.UnitPrice, req.Discount)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Detalle de compra actualizado correctamente",
	})
}

// deletePurchaseDetailJSON elimina un detalle de compra por JSON.
func (server *Server) deletePurchaseDetailJSON(ctx *gin.Context) {
	var req deletePurchaseDetailRequestJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	server.deletePurchaseDetail(ctx, req.DetailID)
}

// deletePurchaseDetailURI elimina un detalle de compra por URI.
func (server *Server) deletePurchaseDetailURI(ctx *gin.Context) {
	var req deletePurchaseDetailRequestURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	server.deletePurchaseDetail(ctx, req.DetailID)
}

// deletePurchaseDetail contiene la lógica común de eliminar un detalle.
func (server *Server) deletePurchaseDetail(ctx *gin.Context, detailID int32) {
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

	err = repository.NewPurchaseDetailRepository(tx).EliminarDetalleCompra(ctx, detailID)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Detalle de compra eliminado correctamente",
	})
}

// getPurchaseDetails obtiene los detalles de una compra.
func (server *Server) getPurchaseDetails(ctx *gin.Context) {
	var req getPurchaseByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	details, err := server.purchaseDetailRepository.ObtenerDetallesCompra(ctx, req.PurchaseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, details)
}

// getPurchaseDetailByID obtiene un detalle de compra por su ID.
func (server *Server) getPurchaseDetailByID(ctx *gin.Context) {
	var req getPurchaseDetailByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	detail, err := server.purchaseDetailRepository.ObtenerDetalleCompra(ctx, req.DetailID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, detail)
}
