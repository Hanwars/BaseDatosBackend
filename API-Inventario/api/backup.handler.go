package api

import (
	"context"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"rest/security"

	"github.com/gin-gonic/gin"
)

func (server *Server) createDatabaseBackup(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Solo administradores pueden generar respaldos"})
		return
	}

	backupDir := strings.TrimSpace(server.backupDir)
	if backupDir == "" {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "BACKUP_DIR no está configurado en app.env"})
		return
	}

	now := time.Now()
	fileName := "inventariodb_" + now.Format("20060102_150405_000000000") + ".bak"
	backupPath := filepath.Clean(filepath.Join(backupDir, fileName))

	backupCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Minute)
	defer cancel()

	_, backupErr := server.conn.ExecContext(backupCtx, `
		BACKUP DATABASE [inventariodb]
		TO DISK = @p1
		WITH COPY_ONLY, CHECKSUM, STATS = 10;`, backupPath)

	status := "OK"
	var details any
	if backupErr != nil {
		status = "ERROR"
		details = backupErr.Error()
	}

	auditCtx, auditCancel := context.WithTimeout(ctx.Request.Context(), 15*time.Second)
	defer auditCancel()
	_, auditErr := server.conn.ExecContext(
		auditCtx,
		`EXEC dbo.paRegistrarAuditoriaRespaldo
			@IdUsuarioAplicacion = @p1,
			@NombreArchivo = @p2,
			@RutaArchivo = @p3,
			@Estado = @p4,
			@Detalle = @p5`,
		authPayload.UserID, fileName, backupPath, status, details,
	)

	if backupErr != nil {
		response := gin.H{
			"error":   "SQL Server no pudo completar el respaldo. Verifique permisos y la ruta configurada.",
			"details": backupErr.Error(),
		}
		if auditErr != nil {
			response["audit_error"] = auditErr.Error()
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	if auditErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":          "El respaldo se creó, pero no fue posible registrar la operación en auditoría.",
			"backup_created": true,
			"file_name":      fileName,
			"backup_path":    backupPath,
			"audit_error":    auditErr.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":     "Respaldo de inventariodb creado y registrado en auditoría",
		"file_name":   fileName,
		"backup_path": backupPath,
		"created_at":  now,
	})
}

func (server *Server) getDatabaseBackupHistory(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Solo administradores pueden consultar los respaldos"})
		return
	}

	limit := 20
	if raw := strings.TrimSpace(ctx.Query("limit")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > 100 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "limit debe estar entre 1 y 100"})
			return
		}
		limit = parsed
	}

	history, err := server.reportRepository.ObtenerHistorialRespaldos(ctx.Request.Context(), limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, history)
}
