package repository

import (
	"context"
	"database/sql"

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
		"EXEC dbo.paIniciarSesion @Correo = @p1",
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
		"EXEC dbo.paObtenerPerfil @IdUsuario = @p1",
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
		"EXEC dbo.paActualizarPerfil @IdUsuario = @p1, @Nombre = @p2, @Correo = @p3, @ContrasenaActual = @p4, @ContrasenaNueva = @p5",
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

// CrearUsuario ejecuta paCrearUsuario y devuelve el identificador generado.
func (r *UserRepository) CrearUsuario(
	ctx context.Context,
	roleID int32,
	name string,
	password string,
	email *string,
) (int32, error) {
	stmt, err := r.db.PrepareContext(ctx,
		"EXEC dbo.paCrearUsuario @IdRol = @p1, @Nombre = @p2, @Contrasena = @p3, @Correo = @p4")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var userID int32
	err = stmt.QueryRowContext(ctx, roleID, name, password, email).Scan(&userID)
	return userID, err
}

// ObtenerUsuarios devuelve los usuarios sin exponer sus contraseñas.
func (r *UserRepository) ObtenerUsuarios(ctx context.Context) ([]models.UserAdmin, error) {
	stmt, err := r.db.PrepareContext(ctx, "EXEC dbo.paObtenerUsuarios")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.UserAdmin, 0)
	for rows.Next() {
		var user models.UserAdmin
		var email sql.NullString
		if err := rows.Scan(&user.UserID, &user.RoleID, &user.RoleName, &user.Name, &email); err != nil {
			return nil, err
		}
		if email.Valid {
			user.Email = &email.String
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// ObtenerUsuario devuelve un usuario por ID sin incluir la contraseña.
func (r *UserRepository) ObtenerUsuario(ctx context.Context, userID int32) (*models.UserAdmin, error) {
	stmt, err := r.db.PrepareContext(ctx,
		"EXEC dbo.paObtenerUsuario @IdUsuario = @p1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user models.UserAdmin
	var email sql.NullString
	err = stmt.QueryRowContext(ctx, userID).Scan(
		&user.UserID, &user.RoleID, &user.RoleName, &user.Name, &email)
	if err != nil {
		return nil, err
	}
	if email.Valid {
		user.Email = &email.String
	}
	return &user, nil
}

// ActualizarUsuario puede cambiar contraseña; nil significa conservar la actual.
func (r *UserRepository) ActualizarUsuario(
	ctx context.Context,
	userID int32,
	roleID int32,
	name string,
	password *string,
	email *string,
) (*models.UserAdmin, error) {
	stmt, err := r.db.PrepareContext(ctx,
		"EXEC dbo.paActualizarUsuario @IdUsuario = @p1, @IdRol = @p2, @Nombre = @p3, @Contrasena = @p4, @Correo = @p5")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user models.UserAdmin
	var userEmail sql.NullString
	err = stmt.QueryRowContext(ctx, userID, roleID, name, password, email).Scan(
		&user.UserID, &user.RoleID, &user.RoleName, &user.Name, &userEmail)
	if err != nil {
		return nil, err
	}
	if userEmail.Valid {
		user.Email = &userEmail.String
	}
	return &user, nil
}

func (r *UserRepository) EliminarUsuario(ctx context.Context, userID int32) error {
	stmt, err := r.db.PrepareContext(ctx,
		"EXEC dbo.paEliminarUsuario @IdUsuario = @p1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, userID)
	return err
}
