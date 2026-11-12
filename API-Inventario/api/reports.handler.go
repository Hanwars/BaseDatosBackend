package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) getSalesSummaryByClient(ctx *gin.Context) {
	items, err := server.reportRepository.ObtenerResumenVentasPorCliente(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func (server *Server) getCurrentInventoryReport(ctx *gin.Context) {
	items, err := server.reportRepository.ObtenerInventarioActual(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func (server *Server) getDiscountMovementsReport(ctx *gin.Context) {
	items, err := server.reportRepository.ObtenerMovimientosConDescuento(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, items)
}
