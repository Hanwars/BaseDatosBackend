package repository

import (
	"context"

	"rest/models"
)

type SaleDetailRepository struct {
	db DBTX
}

func NewSaleDetailRepository(db DBTX) *SaleDetailRepository {
	return &SaleDetailRepository{
		db: db,
	}
}

// CrearDetalleVenta ejecuta el procedimiento almacenado paCrearDetalleVenta.
func (r *SaleDetailRepository) CrearDetalleVenta(
	ctx context.Context,
	saleID int32,
	supplyID int32,
	quantity int32,
	unitPrice float64,
	discount float64,
) (*models.SaleDetail, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paCrearDetalleVenta @IdVenta = @p1, @IdInsumo = @p2, @Cantidad = @p3, @PrecioUnitario = @p4, @Descuento = @p5",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var detail models.SaleDetail

	err = stmt.QueryRowContext(
		ctx,
		saleID,
		supplyID,
		quantity,
		unitPrice,
		discount,
	).Scan(
		&detail.DetailID,
		&detail.SaleID,
		&detail.SupplyID,
		&detail.SupplyName,
		&detail.Quantity,
		&detail.UnitPrice,
		&detail.Discount,
		&detail.Subtotal,
	)

	if err != nil {
		return nil, err
	}

	return &detail, nil
}

// ObtenerDetallesVenta ejecuta el procedimiento almacenado paObtenerDetallesVenta.
func (r *SaleDetailRepository) ObtenerDetallesVenta(
	ctx context.Context,
	saleID int32,
) ([]models.SaleDetail, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paObtenerDetallesVenta @IdVenta = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []models.SaleDetail

	for rows.Next() {

		var detail models.SaleDetail

		err := rows.Scan(
			&detail.DetailID,
			&detail.SaleID,
			&detail.SupplyID,
			&detail.SupplyName,
			&detail.Quantity,
			&detail.UnitPrice,
			&detail.Discount,
			&detail.Subtotal,
		)

		if err != nil {
			return nil, err
		}

		details = append(details, detail)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return details, nil
}

// ObtenerDetalleVenta ejecuta el procedimiento almacenado paObtenerDetalleVenta.
func (r *SaleDetailRepository) ObtenerDetalleVenta(
	ctx context.Context,
	detailID int32,
) (*models.SaleDetail, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paObtenerDetalleVenta @IdDetalle = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var detail models.SaleDetail

	err = stmt.QueryRowContext(
		ctx,
		detailID,
	).Scan(
		&detail.DetailID,
		&detail.SaleID,
		&detail.SupplyID,
		&detail.SupplyName,
		&detail.Quantity,
		&detail.UnitPrice,
		&detail.Discount,
		&detail.Subtotal,
	)

	if err != nil {
		return nil, err
	}

	return &detail, nil
}

// ActualizarDetalleVenta ejecuta el procedimiento almacenado paActualizarDetalleVenta.
func (r *SaleDetailRepository) ActualizarDetalleVenta(
	ctx context.Context,
	detailID int32,
	supplyID int32,
	quantity int32,
	unitPrice float64,
	discount float64,
) (*models.SaleDetail, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paActualizarDetalleVenta @IdDetalle = @p1, @IdInsumo = @p2, @Cantidad = @p3, @PrecioUnitario = @p4, @Descuento = @p5",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var detail models.SaleDetail

	err = stmt.QueryRowContext(
		ctx,
		detailID,
		supplyID,
		quantity,
		unitPrice,
		discount,
	).Scan(
		&detail.DetailID,
		&detail.SaleID,
		&detail.SupplyID,
		&detail.SupplyName,
		&detail.Quantity,
		&detail.UnitPrice,
		&detail.Discount,
		&detail.Subtotal,
	)

	if err != nil {
		return nil, err
	}

	return &detail, nil
}

// EliminarDetalleVenta ejecuta el procedimiento almacenado paEliminarDetalleVenta.
func (r *SaleDetailRepository) EliminarDetalleVenta(
	ctx context.Context,
	detailID int32,
) error {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paEliminarDetalleVenta @IdDetalle = @p1",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, detailID)

	return err
}
