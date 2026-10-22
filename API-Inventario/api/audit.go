package api

import (
	"context"
	"database/sql"

	"rest/security"

	"github.com/gin-gonic/gin"
)

/*
Auditoría con el usuario autenticado.
Los triggers de auditoría (trg_AuditoriaClientes, trg_AuditoriaInsumos y
trg_AuditoriaUsuarios) leen SESSION_CONTEXT('idUsuario'). Ese valor debe
establecerse en la MISMA conexión SQL que hace la operación, por eso se usa
una transacción (*sql.Tx), que el repositorio recibe gracias a DBTX.
*/

// beginAuditTx inicia una transacción y guarda en SESSION_CONTEXT el id del
// usuario que viene en el token.
func (server *Server) beginAuditTx(ctx *gin.Context) (*sql.Tx, error) {
	authPayload := ctx.MustGet("authorized").(*security.Payload)

	tx, err := server.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(
		ctx,
		"EXEC sys.sp_set_session_context @key = N'idUsuario', @value = @p1",
		authPayload.UserID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return tx, nil
}

// endAuditTx limpia SESSION_CONTEXT antes de devolver la conexión al pool.
// Si la operación falló (opErr) revierte la transacción; si no, la confirma.
func endAuditTx(ctx context.Context, tx *sql.Tx, opErr error) error {
	_, clearErr := tx.ExecContext(
		ctx,
		"EXEC sys.sp_set_session_context @key = N'idUsuario', @value = NULL",
	)

	if opErr != nil {
		tx.Rollback()
		return opErr
	}

	if clearErr != nil {
		tx.Rollback()
		return clearErr
	}

	return tx.Commit()
}
