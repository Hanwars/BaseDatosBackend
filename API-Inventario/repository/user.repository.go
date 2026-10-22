package repository

import (
	"context"

	"rest/models"
)

type UserRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Login(ctx context.Context, email string) (*models.User, error) {

	// Preparar el procedimiento almacenado.
	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paIniciarSesion @Correo = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ejecutar el procedimiento.
	row := stmt.QueryRowContext(ctx, email)

	// Modelo que recibirá los datos del usuario.
	var user models.User

	err = row.Scan(
		&user.UserID,
		&user.RoleID,
		&user.Name,
		&user.Password,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ObtenerPerfil ejecuta el procedimiento almacenado paObtenerPerfil.
func (r *UserRepository) ObtenerPerfil(ctx context.Context, userID int32) (*models.User, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paObtenerPerfil @IdUsuario = @p1",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user models.User

	err = stmt.QueryRowContext(ctx, userID).Scan(
		&user.UserID,
		&user.Name,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ActualizarPerfil ejecuta el procedimiento almacenado paActualizarPerfil.
// Si newPassword viene vacío, la contraseña no cambia.
func (r *UserRepository) ActualizarPerfil(
	ctx context.Context,
	userID int32,
	name string,
	email *string,
	currentPassword *string,
	newPassword *string,
) (*models.User, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		"EXEC paActualizarPerfil @IdUsuario = @p1, @Nombre = @p2, @Correo = @p3, @ContrasenaActual = @p4, @ContrasenaNueva = @p5",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user models.User

	err = stmt.QueryRowContext(
		ctx,
		userID,
		name,
		email,
		currentPassword,
		newPassword,
	).Scan(
		&user.UserID,
		&user.Name,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
