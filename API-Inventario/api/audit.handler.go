package api

import (
	"net/http"
	"strconv"
	"strings"

	"rest/security"

	"github.com/gin-gonic/gin"
)

// getAuditRecords permite consultar los registros escritos por los triggers de auditoría.
// Por contener datos anteriores/nuevos, solo está disponible para administradores.
func (server *Server) getAuditRecords(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Solo administradores pueden consultar la auditoría"})
		return
	}

	limit := 100
	if rawLimit := strings.TrimSpace(ctx.Query("limit")); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed < 1 || parsed > 500 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro limit debe estar entre 1 y 500"})
			return
		}
		limit = parsed
	}

	table := strings.TrimSpace(ctx.Query("table"))
	if table != "" && table != "clientes" && table != "insumos" && table != "usuarios" && table != "auditoriaRespaldos" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "table debe ser clientes, insumos, usuarios o auditoriaRespaldos"})
		return
	}

	operation := strings.ToUpper(strings.TrimSpace(ctx.Query("operation")))
	if operation != "" && operation != "INSERT" && operation != "UPDATE" && operation != "DELETE" && operation != "BACKUP" && operation != "BACKUP_ERROR" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "operation debe ser INSERT, UPDATE, DELETE, BACKUP o BACKUP_ERROR"})
		return
	}

	records, err := server.reportRepository.ObtenerAuditoria(
		ctx.Request.Context(), limit, table, operation)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, records)
}
