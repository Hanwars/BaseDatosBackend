package repository

import (
	"context"

	"rest/models"
)

type ClientRepository struct {
	db DBTX
}

func NewClientRepository(db DBTX) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

// CrearCliente ejecuta el procedimiento almacenado paCrearCliente.
func (r *ClientRepository) CrearCliente(
	ctx context.Context,
	name string,
	contactPhone int32,
	contactEmail string,
) (*models.Client, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paCrearCliente @Nombre = @p1, @TelefonoContacto = @p2, @CorreoContacto = @p3",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var clientID int32

	err = stmt.QueryRowContext(
		ctx,
		name,
		contactPhone,
		contactEmail,
	).Scan(&clientID)

	if err != nil {
		return nil, err
	}

	return &models.Client{
		ClientID:     clientID,
		Name:         name,
		ContactPhone: contactPhone,
		ContactEmail: contactEmail,
	}, nil
}

// ObtenerClientes ejecuta el procedimiento almacenado paObtenerClientes.
func (r *ClientRepository) ObtenerClientes(
	ctx context.Context,
) ([]models.Client, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paObtenerClientes",
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

	var clients []models.Client

	for rows.Next() {

		var client models.Client

		err := rows.Scan(
			&client.ClientID,
			&client.Name,
			&client.ContactPhone,
			&client.ContactEmail,
		)

		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

// ObtenerCliente ejecuta el procedimiento almacenado paObtenerCliente.
func (r *ClientRepository) ObtenerCliente(
	ctx context.Context,
	clientID int32,
) (*models.Client, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paObtenerCliente @IdCliente = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var client models.Client

	err = stmt.QueryRowContext(
		ctx,
		clientID,
	).Scan(
		&client.ClientID,
		&client.Name,
		&client.ContactPhone,
		&client.ContactEmail,
	)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

// ActualizarCliente ejecuta el procedimiento almacenado paActualizarCliente.
func (r *ClientRepository) ActualizarCliente(
	ctx context.Context,
	clientID int32,
	name string,
	contactPhone int32,
	contactEmail string,
) (*models.Client, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paActualizarCliente @IdCliente = @p1, @Nombre = @p2, @TelefonoContacto = @p3, @CorreoContacto = @p4",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var client models.Client

	err = stmt.QueryRowContext(
		ctx,
		clientID,
		name,
		contactPhone,
		contactEmail,
	).Scan(
		&client.ClientID,
		&client.Name,
		&client.ContactPhone,
		&client.ContactEmail,
	)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

// EliminarCliente ejecuta el procedimiento almacenado paEliminarCliente.
func (r *ClientRepository) EliminarCliente(
	ctx context.Context,
	clientID int32,
) error {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paEliminarCliente @IdCliente = @p1",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, clientID)

	return err
}
