package repository

import (
	"context"

	"rest/models"
)

type PurchaseDetailRepository struct {
	db DBTX
}

func NewPurchaseDetailRepository(db DBTX) *PurchaseDetailRepository {
	return &PurchaseDetailRepository{
		db: db,
	}
}

// CrearDetalleCompra ejecuta el procedimiento almacenado paCrearDetalleCompra.
func (r *PurchaseDetailRepository) CrearDetalleCompra(
	ctx context.Context,
	purchaseID int32,
	supplyID int32,
	quantity int32,
	unitPrice float64,
	discount float64,
) (*models.PurchaseDetail, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paCrearDetalleCompra @IdCompra = @p1, @IdInsumo = @p2, @Cantidad = @p3, @PrecioUnitario = @p4, @Descuento = @p5",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var detail models.PurchaseDetail

	err = stmt.QueryRowContext(
		ctx,
		purchaseID,
		supplyID,
		quantity,
		unitPrice,
		discount,
	).Scan(
		&detail.DetailID,
		&detail.PurchaseID,
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

// ObtenerDetallesCompra ejecuta el procedimiento almacenado paObtenerDetallesCompra.
func (r *PurchaseDetailRepository) ObtenerDetallesCompra(
	ctx context.Context,
	purchaseID int32,
) ([]models.PurchaseDetail, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paObtenerDetallesCompra @IdCompra = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, purchaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []models.PurchaseDetail

	for rows.Next() {

		var detail models.PurchaseDetail

		err := rows.Scan(
			&detail.DetailID,
			&detail.PurchaseID,
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

// ObtenerDetalleCompra ejecuta el procedimiento almacenado paObtenerDetalleCompra.
func (r *PurchaseDetailRepository) ObtenerDetalleCompra(
	ctx context.Context,
	detailID int32,
) (*models.PurchaseDetail, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paObtenerDetalleCompra @IdDetalle = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var detail models.PurchaseDetail

	err = stmt.QueryRowContext(
		ctx,
		detailID,
	).Scan(
		&detail.DetailID,
		&detail.PurchaseID,
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

// ActualizarDetalleCompra ejecuta el procedimiento almacenado paActualizarDetalleCompra.
func (r *PurchaseDetailRepository) ActualizarDetalleCompra(
	ctx context.Context,
	detailID int32,
	supplyID int32,
	quantity int32,
	unitPrice float64,
	discount float64,
) (*models.PurchaseDetail, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paActualizarDetalleCompra @IdDetalle = @p1, @IdInsumo = @p2, @Cantidad = @p3, @PrecioUnitario = @p4, @Descuento = @p5",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var detail models.PurchaseDetail

	err = stmt.QueryRowContext(
		ctx,
		detailID,
		supplyID,
		quantity,
		unitPrice,
		discount,
	).Scan(
		&detail.DetailID,
		&detail.PurchaseID,
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

// EliminarDetalleCompra ejecuta el procedimiento almacenado paEliminarDetalleCompra.
func (r *PurchaseDetailRepository) EliminarDetalleCompra(
	ctx context.Context,
	detailID int32,
) error {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paEliminarDetalleCompra @IdDetalle = @p1",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, detailID)

	return err
}
