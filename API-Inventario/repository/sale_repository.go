package repository

import (
	"context"
	"time"

	"rest/models"
)

type SaleRepository struct {
	db DBTX
}

func NewSaleRepository(db DBTX) *SaleRepository {
	return &SaleRepository{
		db: db,
	}
}

// CrearVenta ejecuta el procedimiento almacenado paCrearVenta.
func (r *SaleRepository) CrearVenta(
	ctx context.Context,
	userID int32,
	clientID int32,
	date *time.Time,
	iva float64,
) (*models.Sale, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paCrearVenta @IdUsuario = @p1, @IdCliente = @p2, @Fecha = @p3, @IVA = @p4",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var saleID int32

	err = stmt.QueryRowContext(
		ctx,
		userID,
		clientID,
		date,
		iva,
	).Scan(&saleID)

	if err != nil {
		return nil, err
	}

	return &models.Sale{
		SaleID:   saleID,
		UserID:   userID,
		ClientID: clientID,
		Date:     timeValue(date),
		Subtotal: 0,
		Discount: 0,
		IVA:      iva,
		Total:    0,
	}, nil
}

// timeValue devuelve una fecha válida para el modelo.
// Si no se proporciona fecha, se utiliza la fecha actual.
func timeValue(date *time.Time) time.Time {
	if date != nil {
		return *date
	}

	return time.Now()
}

// ObtenerVentas ejecuta el procedimiento almacenado paObtenerVentas.
func (r *SaleRepository) ObtenerVentas(ctx context.Context) ([]models.Sale, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paObtenerVentas",
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

	var sales []models.Sale

	for rows.Next() {

		var sale models.Sale

		err := rows.Scan(
			&sale.SaleID,
			&sale.UserID,
			&sale.UserName,
			&sale.ClientID,
			&sale.ClientName,
			&sale.Date,
			&sale.Subtotal,
			&sale.Discount,
			&sale.IVA,
			&sale.Total,
		)

		if err != nil {
			return nil, err
		}

		sales = append(sales, sale)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sales, nil
}

// ObtenerVenta ejecuta el procedimiento almacenado paObtenerVenta.
func (r *SaleRepository) ObtenerVenta(
	ctx context.Context,
	saleID int32,
) (*models.Sale, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paObtenerVenta @IdVenta = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var sale models.Sale

	err = stmt.QueryRowContext(
		ctx,
		saleID,
	).Scan(
		&sale.SaleID,
		&sale.UserID,
		&sale.UserName,
		&sale.ClientID,
		&sale.ClientName,
		&sale.Date,
		&sale.Subtotal,
		&sale.Discount,
		&sale.IVA,
		&sale.Total,
	)

	if err != nil {
		return nil, err
	}

	return &sale, nil
}

// ActualizarVenta ejecuta el procedimiento almacenado paActualizarVenta.
func (r *SaleRepository) ActualizarVenta(
	ctx context.Context,
	saleID int32,
	userID int32,
	clientID int32,
	date *time.Time,
	iva float64,
) (*models.Sale, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paActualizarVenta @IdVenta = @p1, @IdUsuario = @p2, @IdCliente = @p3, @Fecha = @p4, @IVA = @p5",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var sale models.Sale

	err = stmt.QueryRowContext(
		ctx,
		saleID,
		userID,
		clientID,
		date,
		iva,
	).Scan(
		&sale.SaleID,
		&sale.UserID,
		&sale.UserName,
		&sale.ClientID,
		&sale.ClientName,
		&sale.Date,
		&sale.Subtotal,
		&sale.Discount,
		&sale.IVA,
		&sale.Total,
	)

	if err != nil {
		return nil, err
	}

	return &sale, nil
}

// EliminarVenta ejecuta el procedimiento almacenado paEliminarVenta.
func (r *SaleRepository) EliminarVenta(
	ctx context.Context,
	saleID int32,
) error {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paEliminarVenta @IdVenta = @p1",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, saleID)

	return err
}
