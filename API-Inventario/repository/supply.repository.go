package repository

import (
	"context"

	"rest/models"
)

type SupplyRepository struct {
	db DBTX
}

func NewSupplyRepository(db DBTX) *SupplyRepository {
	return &SupplyRepository{
		db: db,
	}
}

// CrearInsumo ejecuta el procedimiento almacenado paCrearInsumo.
func (r *SupplyRepository) CrearInsumo(
	ctx context.Context,
	name string,
	supplyType string,
	stock int32,
	categoryID int32,
	price float64,
	image string,
) (*models.Supply, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paCrearInsumo @Nombre = @p1, @Tipo = @p2, @Stock = @p3, @IdCategoria = @p4, @Precio = @p5, @Imagen = @p6",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var supplyID int32

	err = stmt.QueryRowContext(
		ctx,
		name,
		supplyType,
		stock,
		categoryID,
		price,
		image,
	).Scan(&supplyID)

	if err != nil {
		return nil, err
	}

	return &models.Supply{
		SupplyID:   supplyID,
		Name:       name,
		Type:       supplyType,
		Stock:      stock,
		CategoryID: categoryID,
		Price:      price,
		Image:      image,
	}, nil
}

// ObtenerInsumos ejecuta el procedimiento almacenado paObtenerInsumos.
func (r *SupplyRepository) ObtenerInsumos(
	ctx context.Context,
) ([]models.Supply, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paObtenerInsumos",
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

	var insumos []models.Supply

	for rows.Next() {
		var insumo models.Supply

		// Orden devuelto por SQL Server:
		// idInsumo, nombre, tipo, stock, idCategoria,
		// precio, imagen, categoria.
		err := rows.Scan(
			&insumo.SupplyID,
			&insumo.Name,
			&insumo.Type,
			&insumo.Stock,
			&insumo.CategoryID,
			&insumo.Price,
			&insumo.Image,
			&insumo.CategoryName,
		)

		if err != nil {
			return nil, err
		}

		insumos = append(insumos, insumo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return insumos, nil
}

// ObtenerInsumo ejecuta el procedimiento almacenado paObtenerInsumo.
func (r *SupplyRepository) ObtenerInsumo(
	ctx context.Context,
	supplyID int32,
) (*models.Supply, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paObtenerInsumo @IdInsumo = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var insumo models.Supply

	err = stmt.QueryRowContext(
		ctx,
		supplyID,
	).Scan(
		&insumo.SupplyID,
		&insumo.Name,
		&insumo.Type,
		&insumo.Stock,
		&insumo.CategoryID,
		&insumo.Price,
		&insumo.Image,
		&insumo.CategoryName,
	)

	if err != nil {
		return nil, err
	}

	return &insumo, nil
}

// ActualizarInsumo ejecuta el procedimiento almacenado paActualizarInsumo.
func (r *SupplyRepository) ActualizarInsumo(
	ctx context.Context,
	supplyID int32,
	name string,
	supplyType string,
	stock int32,
	categoryID int32,
	price float64,
	image string,
) (*models.Supply, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paActualizarInsumo @IdInsumo = @p1, @Nombre = @p2, @Tipo = @p3, @Stock = @p4, @IdCategoria = @p5, @Precio = @p6, @Imagen = @p7",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var insumo models.Supply
	var categoryName string

	err = stmt.QueryRowContext(
		ctx,
		supplyID,
		name,
		supplyType,
		stock,
		categoryID,
		price,
		image,
	).Scan(
		&insumo.SupplyID,
		&insumo.Name,
		&insumo.Type,
		&insumo.Stock,
		&insumo.CategoryID,
		&insumo.Price,
		&insumo.Image,
		&categoryName,
	)

	if err != nil {
		return nil, err
	}

	insumo.CategoryName = categoryName

	return &insumo, nil
}

// EliminarInsumo ejecuta el procedimiento almacenado paEliminarInsumo.
func (r *SupplyRepository) EliminarInsumo(
	ctx context.Context,
	supplyID int32,
) error {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paEliminarInsumo @IdInsumo = @p1",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, supplyID)

	return err
}