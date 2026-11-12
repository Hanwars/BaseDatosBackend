package repository

import (
	"context"
	"database/sql"

	"rest/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) ObtenerResumenVentasPorCliente(ctx context.Context) ([]models.SalesSummary, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT idCliente, nombreCliente, cantidadVentas, totalVentas, ultimaVenta
		FROM dbo.vResumenVentasPorCliente
		ORDER BY totalVentas DESC, nombreCliente ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.SalesSummary, 0)
	for rows.Next() {
		var item models.SalesSummary
		var lastSale sql.NullTime
		if err := rows.Scan(&item.ClientID, &item.ClientName, &item.SalesCount, &item.TotalSales, &lastSale); err != nil {
			return nil, err
		}
		if lastSale.Valid {
			item.LastSaleDate = &lastSale.Time
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *ReportRepository) ObtenerInventarioActual(ctx context.Context) ([]models.InventoryReportItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT idInsumo, nombre, tipo, stock, idCategoria,
		       nombreCategoria, precio, imagen, estadoStock
		FROM dbo.vInventarioActual
		ORDER BY nombreCategoria, nombre`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.InventoryReportItem, 0)
	for rows.Next() {
		var item models.InventoryReportItem
		if err := rows.Scan(
			&item.SupplyID, &item.Name, &item.Type, &item.Stock,
			&item.CategoryID, &item.CategoryName, &item.Price,
			&item.Image, &item.StockStatus,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *ReportRepository) ObtenerMovimientosConDescuento(ctx context.Context) ([]models.DiscountMovement, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT idMovimiento, subtotal, descuento, tipoMovimiento
		FROM dbo.vMovimientosDescuento
		ORDER BY tipoMovimiento, idMovimiento`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.DiscountMovement, 0)
	for rows.Next() {
		var item models.DiscountMovement
		if err := rows.Scan(&item.MovementID, &item.Subtotal, &item.Discount, &item.Type); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *ReportRepository) ObtenerAuditoria(ctx context.Context, limit int, affectedTable, operation string) ([]models.AuditRecord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT TOP (@p1)
		       tablaAuditoria, idAuditoria, tablaAfectada, operacion, idRegistro,
		       fechaHora, loginSQL, idUsuarioAplicacion, valoresAnteriores,
		       valoresNuevos, contrasenaCambiada
		FROM dbo.vAuditoriaSistema
		WHERE (@p2 = N'' OR tablaAfectada = @p2)
		  AND (@p3 = N'' OR operacion = @p3)
		ORDER BY fechaHora DESC, idAuditoria DESC, tablaAuditoria ASC`,
		limit, affectedTable, operation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.AuditRecord, 0)
	for rows.Next() {
		var item models.AuditRecord
		var sqlLogin, previousValues, newValues sql.NullString
		var applicationUser sql.NullInt32
		var passwordChanged sql.NullBool
		if err := rows.Scan(
			&item.AuditTable, &item.AuditID, &item.AffectedTable,
			&item.Operation, &item.RecordID, &item.DateTime,
			&sqlLogin, &applicationUser, &previousValues, &newValues,
			&passwordChanged,
		); err != nil {
			return nil, err
		}
		if sqlLogin.Valid {
			item.SQLLogin = &sqlLogin.String
		}
		if applicationUser.Valid {
			value := applicationUser.Int32
			item.ApplicationUser = &value
		}
		if previousValues.Valid {
			item.PreviousValues = &previousValues.String
		}
		if newValues.Valid {
			item.NewValues = &newValues.String
		}
		if passwordChanged.Valid {
			value := passwordChanged.Bool
			item.PasswordChanged = &value
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *ReportRepository) ObtenerHistorialRespaldos(ctx context.Context, limit int) ([]models.BackupRecord, error) {
	rows, err := r.db.QueryContext(ctx, `
		EXEC dbo.paObtenerHistorialRespaldos @Limite = @p1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.BackupRecord, 0)
	for rows.Next() {
		var item models.BackupRecord
		var details, sqlLogin sql.NullString
		var applicationUser sql.NullInt32
		if err := rows.Scan(
			&item.AuditID, &item.FileName, &item.Path, &item.CreatedAt,
			&item.Status, &details, &applicationUser, &sqlLogin,
		); err != nil {
			return nil, err
		}
		if details.Valid {
			item.Details = &details.String
		}
		if applicationUser.Valid {
			value := applicationUser.Int32
			item.ApplicationUser = &value
		}
		if sqlLogin.Valid {
			item.SQLLogin = &sqlLogin.String
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
