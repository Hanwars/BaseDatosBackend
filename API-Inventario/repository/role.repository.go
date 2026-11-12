package repository

import (
	"context"
	"database/sql"

	"rest/models"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

// CrearRol ejecuta el procedimiento almacenado paCrearRol.
func (r *RoleRepository) CrearRol(ctx context.Context, name string) (*models.Role, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paCrearRol @Nombre = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var role models.Role

	err = stmt.QueryRowContext(ctx, name).Scan(
		&role.RoleID,
	)

	if err != nil {
		return nil, err
	}

	role.Name = name

	return &role, nil
}

// ObtenerRoles ejecuta el procedimiento almacenado paObtenerRoles.
func (r *RoleRepository) ObtenerRoles(ctx context.Context) ([]models.Role, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paObtenerRoles",
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

	var roles []models.Role

	for rows.Next() {

		var role models.Role

		err := rows.Scan(
			&role.RoleID,
			&role.Name,
		)

		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

// ObtenerRol ejecuta el procedimiento almacenado paObtenerRol.
func (r *RoleRepository) ObtenerRol(ctx context.Context, roleID int32) (*models.Role, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paObtenerRol @IdRol = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var role models.Role

	err = stmt.QueryRowContext(ctx, roleID).Scan(
		&role.RoleID,
		&role.Name,
	)

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// ActualizarRol ejecuta el procedimiento almacenado paActualizarRol.
func (r *RoleRepository) ActualizarRol(
	ctx context.Context,
	roleID int32,
	name string,
) (*models.Role, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paActualizarRol @IdRol = @p1, @Nombre = @p2",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var role models.Role

	err = stmt.QueryRowContext(
		ctx,
		roleID,
		name,
	).Scan(
		&role.RoleID,
		&role.Name,
	)

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// EliminarRol ejecuta el procedimiento almacenado paEliminarRol.
func (r *RoleRepository) EliminarRol(ctx context.Context, roleID int32) error {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC dbo.paEliminarRol @IdRol = @p1",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, roleID)

	return err
}
