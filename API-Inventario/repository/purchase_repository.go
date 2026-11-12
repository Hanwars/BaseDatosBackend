package repository

import (
	"context"
	"database/sql"
	"time"

	"rest/models"
)

// DBTX permite utilizar tanto *sql.DB como *sql.Tx.
type DBTX interface {
	PrepareContext(context.Context, string) (*sql.Stmt, error)
}

type PurchaseRepository struct {
	db DBTX
}

func NewPurchaseRepository(db DBTX) *PurchaseRepository {
	return &PurchaseRepository{
		db: db,
	}
}

// CrearCompra ejecuta el procedimiento almacenado paCrearCompra.
func (r *PurchaseRepository) CrearCompra(
	ctx context.Context,
	userID int32,
	supplierID int32,
	date *time.Time,
	iva float64,
) (*models.Purchase, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paCrearCompra @IdUsuario = @p1, @IdProveedor = @p2, @Fecha = @p3, @IVA = @p4",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var purchaseID int32

	err = stmt.QueryRowContext(
		ctx,
		userID,
		supplierID,
		date,
		iva,
	).Scan(&purchaseID)

	if err != nil {
		return nil, err
	}

	purchaseDate := time.Now()

	if date != nil {
		purchaseDate = *date
	}

	return &models.Purchase{
		PurchaseID: purchaseID,
		UserID:     userID,
		SupplierID: supplierID,
		Date:       purchaseDate,
		Subtotal:   0,
		Discount:   0,
		IVA:        iva,
		Total:      0,
	}, nil
}

// ObtenerCompras ejecuta el procedimiento almacenado paObtenerCompras.
func (r *PurchaseRepository) ObtenerCompras(
	ctx context.Context,
) ([]models.Purchase, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paObtenerCompras",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []models.Purchase

	for rows.Next() {

		var purchase models.Purchase

		err := rows.Scan(
			&purchase.PurchaseID,
			&purchase.UserID,
			&purchase.UserName,
			&purchase.SupplierID,
			&purchase.SupplierName,
			&purchase.Date,
			&purchase.Subtotal,
			&purchase.Discount,
			&purchase.IVA,
			&purchase.Total,
		)

		if err != nil {
			return nil, err
		}

		purchases = append(purchases, purchase)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return purchases, nil
}

// ObtenerCompra ejecuta el procedimiento almacenado paObtenerCompra.
func (r *PurchaseRepository) ObtenerCompra(
	ctx context.Context,
	purchaseID int32,
) (*models.Purchase, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paObtenerCompra @IdCompra = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var purchase models.Purchase

	err = stmt.QueryRowContext(
		ctx,
		purchaseID,
	).Scan(
		&purchase.PurchaseID,
		&purchase.UserID,
		&purchase.UserName,
		&purchase.SupplierID,
		&purchase.SupplierName,
		&purchase.Date,
		&purchase.Subtotal,
		&purchase.Discount,
		&purchase.IVA,
		&purchase.Total,
	)

	if err != nil {
		return nil, err
	}

	return &purchase, nil
}

// ActualizarCompra ejecuta el procedimiento almacenado paActualizarCompra.
func (r *PurchaseRepository) ActualizarCompra(
	ctx context.Context,
	purchaseID int32,
	userID int32,
	supplierID int32,
	date *time.Time,
	iva float64,
) (*models.Purchase, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paActualizarCompra @IdCompra = @p1, @IdUsuario = @p2, @IdProveedor = @p3, @Fecha = @p4, @IVA = @p5",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var purchase models.Purchase

	err = stmt.QueryRowContext(
		ctx,
		purchaseID,
		userID,
		supplierID,
		date,
		iva,
	).Scan(
		&purchase.PurchaseID,
		&purchase.UserID,
		&purchase.UserName,
		&purchase.SupplierID,
		&purchase.SupplierName,
		&purchase.Date,
		&purchase.Subtotal,
		&purchase.Discount,
		&purchase.IVA,
		&purchase.Total,
	)

	if err != nil {
		return nil, err
	}

	return &purchase, nil
}

// EliminarCompra ejecuta el procedimiento almacenado paEliminarCompra.
func (r *PurchaseRepository) EliminarCompra(
	ctx context.Context,
	purchaseID int32,
) error {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paEliminarCompra @IdCompra = @p1",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, purchaseID)

	return err
}
