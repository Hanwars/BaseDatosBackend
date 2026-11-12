package repository

import (
	"context"

	"rest/models"
)

type SupplierRepository struct {
	db DBTX
}

func NewSupplierRepository(db DBTX) *SupplierRepository {
	return &SupplierRepository{
		db: db,
	}
}

// CrearProveedor ejecuta el procedimiento almacenado paCrearProveedor.
func (r *SupplierRepository) CrearProveedor(
	ctx context.Context,
	companyName string,
	companyPhone int32,
	companyEmail string,
	contactName *string,
	contactPhone *int32,
	contactEmail *string,
) (*models.Supplier, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paCrearProveedor @NombreEmpresa = @p1, @TelefonoEmpresa = @p2, @CorreoEmpresa = @p3, @NombreContacto = @p4, @TelefonoContacto = @p5, @CorreoContacto = @p6",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var supplierID int32

	err = stmt.QueryRowContext(
		ctx,
		companyName,
		companyPhone,
		companyEmail,
		contactName,
		contactPhone,
		contactEmail,
	).Scan(&supplierID)

	if err != nil {
		return nil, err
	}

	return &models.Supplier{
		SupplierID:   supplierID,
		CompanyName:  companyName,
		CompanyPhone: companyPhone,
		CompanyEmail: companyEmail,
		ContactName:  contactName,
		ContactPhone: contactPhone,
		ContactEmail: contactEmail,
	}, nil
}

// ObtenerProveedores ejecuta el procedimiento almacenado paObtenerProveedores.
func (r *SupplierRepository) ObtenerProveedores(
	ctx context.Context,
) ([]models.Supplier, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paObtenerProveedores",
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

	var suppliers []models.Supplier

	for rows.Next() {

		var supplier models.Supplier

		err := rows.Scan(
			&supplier.SupplierID,
			&supplier.CompanyName,
			&supplier.CompanyPhone,
			&supplier.CompanyEmail,
			&supplier.ContactName,
			&supplier.ContactPhone,
			&supplier.ContactEmail,
		)

		if err != nil {
			return nil, err
		}

		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return suppliers, nil
}

// ObtenerProveedor ejecuta el procedimiento almacenado paObtenerProveedor.
func (r *SupplierRepository) ObtenerProveedor(
	ctx context.Context,
	supplierID int32,
) (*models.Supplier, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paObtenerProveedor @IdProveedor = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var supplier models.Supplier

	err = stmt.QueryRowContext(
		ctx,
		supplierID,
	).Scan(
		&supplier.SupplierID,
		&supplier.CompanyName,
		&supplier.CompanyPhone,
		&supplier.CompanyEmail,
		&supplier.ContactName,
		&supplier.ContactPhone,
		&supplier.ContactEmail,
	)

	if err != nil {
		return nil, err
	}

	return &supplier, nil
}

// ActualizarProveedor ejecuta el procedimiento almacenado paActualizarProveedor.
func (r *SupplierRepository) ActualizarProveedor(
	ctx context.Context,
	supplierID int32,
	companyName string,
	companyPhone int32,
	companyEmail string,
	contactName *string,
	contactPhone *int32,
	contactEmail *string,
) (*models.Supplier, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paActualizarProveedor @IdProveedor = @p1, @NombreEmpresa = @p2, @TelefonoEmpresa = @p3, @CorreoEmpresa = @p4, @NombreContacto = @p5, @TelefonoContacto = @p6, @CorreoContacto = @p7",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var supplier models.Supplier

	err = stmt.QueryRowContext(
		ctx,
		supplierID,
		companyName,
		companyPhone,
		companyEmail,
		contactName,
		contactPhone,
		contactEmail,
	).Scan(
		&supplier.SupplierID,
		&supplier.CompanyName,
		&supplier.CompanyPhone,
		&supplier.CompanyEmail,
		&supplier.ContactName,
		&supplier.ContactPhone,
		&supplier.ContactEmail,
	)

	if err != nil {
		return nil, err
	}

	return &supplier, nil
}

// EliminarProveedor ejecuta el procedimiento almacenado paEliminarProveedor.
func (r *SupplierRepository) EliminarProveedor(
	ctx context.Context,
	supplierID int32,
) error {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paEliminarProveedor @IdProveedor = @p1",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, supplierID)

	return err
}
