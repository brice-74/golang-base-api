package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

var (
	ErrNotFound       = errors.New("user not found")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type Model struct {
	DB *sql.DB
}

func (m Model) ExistEmail(email string) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM "user_account"
		WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int

	err := m.DB.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (m Model) GetById(id string) (*User, error) {
	return m.getBy("id", id)
}

func (m Model) GetByEmail(email string) (*User, error) {
	return m.getBy("email", email)
}

func (m Model) getBy(column string, value interface{}) (*User, error) {
	query := fmt.Sprintf(`
		SELECT 
			id,
			created_at,
			updated_at,
			deactivated_at,
			email,
			password,
			roles,
			profil_name, 
			short_id
		FROM "user_account"
		WHERE %s = $1`, column)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		user           User
		deactivated_at pq.NullTime
	)

	err := m.DB.QueryRowContext(ctx, query, value).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deactivated_at,
		&user.Email,
		&user.Password,
		pq.Array(&user.Roles),
		&user.ProfilName,
		&user.ShortId,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	user.DeactivatedAt = deactivated_at.Time

	return &user, nil
}

func (m Model) InsertRegisteredUserAccount(user *User) error {
	query := `
		INSERT INTO "user_account" (
			email,
			password,
			roles,
			profil_name, 
			short_id
		) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at, deactivated_at`

	args := []interface{}{
		user.Email,
		user.Password,
		pq.Array(user.Roles),
		user.ProfilName,
		user.ShortId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var deactivated_at pq.NullTime

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &deactivated_at)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "user_account_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	user.DeactivatedAt = deactivated_at.Time

	return nil
}

func (m Model) InsertOrUpdateUserSession(session *Session) error {
	query := `
		INSERT INTO "user_session" (
			id,
			deactivated_at,
			ip,
			name,
			location,
			user_id
		) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET  
			deactivated_at = $2,
			ip = $3,
			name = $4,
			location = $5,
			user_id = $6
		RETURNING created_at, updated_at`

	args := []interface{}{
		session.ID,
		session.DeactivatedAt,
		session.IP,
		session.Name,
		session.Location,
		session.UserID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
