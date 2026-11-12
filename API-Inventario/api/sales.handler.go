package api

import (
	"net/http"
	"rest/repository"
	"rest/security"
	"time"

	"github.com/gin-gonic/gin"
)

// Las ventas son operaciones de trabajo para cualquier usuario autenticado.
// El usuario que registra la venta se toma del token, nunca del JSON.

type createSaleRequest struct {
	ClientID int32      `json:"client_id" binding:"required"`
	Date     *time.Time `json:"date"`
	IVA      float64    `json:"iva"`
}

type updateSaleRequest struct {
	SaleID   int32      `json:"sale_id" binding:"required"`
	ClientID int32      `json:"client_id" binding:"required"`
	Date     *time.Time `json:"date"`
	IVA      float64    `json:"iva"`
}

type deleteSaleRequestJSON struct {
	SaleID int32 `json:"sale_id" binding:"required"`
}

type deleteSaleRequestURI struct {
	SaleID int32 `uri:"sale_id" binding:"required"`
}

type getSaleByIDRequest struct {
	SaleID int32 `uri:"sale_id" binding:"required"`
}

type createSaleDetailRequest struct {
	SaleID    int32   `json:"sale_id" binding:"required"`
	SupplyID  int32   `json:"supply_id" binding:"required"`
	Quantity  int32   `json:"quantity" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
	Discount  float64 `json:"discount"`
}

type updateSaleDetailRequest struct {
	DetailID  int32   `json:"detail_id" binding:"required"`
	SupplyID  int32   `json:"supply_id" binding:"required"`
	Quantity  int32   `json:"quantity" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
	Discount  float64 `json:"discount"`
}

type deleteSaleDetailRequestJSON struct {
	DetailID int32 `json:"detail_id" binding:"required"`
}

type deleteSaleDetailRequestURI struct {
	DetailID int32 `uri:"detail_id" binding:"required"`
}

type getSaleDetailByIDRequest struct {
	DetailID int32 `uri:"detail_id" binding:"required"`
}

// ========================= VENTAS =========================

func (server *Server) createSale(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)

	var req createSaleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sale, err := repository.NewSaleRepository(tx).CrearVenta(
		ctx, authPayload.UserID, req.ClientID, req.Date, req.IVA)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"generated_id": sale.SaleID})
}

func (server *Server) updateSale(ctx *gin.Context) {
	var req updateSaleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	currentSale, err := server.saleRepository.ObtenerVenta(ctx, req.SaleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = repository.NewSaleRepository(tx).ActualizarVenta(
		ctx, req.SaleID, currentSale.UserID, req.ClientID, req.Date, req.IVA)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Venta actualizada correctamente"})
}

func (server *Server) deleteSaleJSON(ctx *gin.Context) {
	var req deleteSaleRequestJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	server.deleteSale(ctx, req.SaleID)
}

func (server *Server) deleteSaleURI(ctx *gin.Context) {
	var req deleteSaleRequestURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	server.deleteSale(ctx, req.SaleID)
}

func (server *Server) deleteSale(ctx *gin.Context, saleID int32) {
	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = repository.NewSaleRepository(tx).EliminarVenta(ctx, saleID)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Venta eliminada correctamente"})
}

func (server *Server) getAllSales(ctx *gin.Context) {
	sales, err := server.saleRepository.ObtenerVentas(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, sales)
}

func (server *Server) getSaleByID(ctx *gin.Context) {
	var req getSaleByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sale, err := server.saleRepository.ObtenerVenta(ctx, req.SaleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, sale)
}

// ===================== DETALLES DE VENTA =====================

func (server *Server) createSaleDetail(ctx *gin.Context) {
	var req createSaleDetailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	detail, err := repository.NewSaleDetailRepository(tx).CrearDetalleVenta(
		ctx, req.SaleID, req.SupplyID, req.Quantity, req.UnitPrice, req.Discount)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"generated_id": detail.DetailID})
}

func (server *Server) updateSaleDetail(ctx *gin.Context) {
	var req updateSaleDetailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = repository.NewSaleDetailRepository(tx).ActualizarDetalleVenta(
		ctx, req.DetailID, req.SupplyID, req.Quantity, req.UnitPrice, req.Discount)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Detalle de venta actualizado correctamente"})
}

func (server *Server) deleteSaleDetailJSON(ctx *gin.Context) {
	var req deleteSaleDetailRequestJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	server.deleteSaleDetail(ctx, req.DetailID)
}

func (server *Server) deleteSaleDetailURI(ctx *gin.Context) {
	var req deleteSaleDetailRequestURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	server.deleteSaleDetail(ctx, req.DetailID)
}

func (server *Server) deleteSaleDetail(ctx *gin.Context, detailID int32) {
	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = repository.NewSaleDetailRepository(tx).EliminarDetalleVenta(ctx, detailID)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Detalle de venta eliminado correctamente"})
}

func (server *Server) getSaleDetails(ctx *gin.Context) {
	var req getSaleByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	details, err := server.saleDetailRepository.ObtenerDetallesVenta(ctx, req.SaleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, details)
}

func (server *Server) getSaleDetailByID(ctx *gin.Context) {
	var req getSaleDetailByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	detail, err := server.saleDetailRepository.ObtenerDetalleVenta(ctx, req.DetailID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, detail)
}
