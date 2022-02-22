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
	ErrNotFoundUserAndSession = errors.New("User or user session not found")
	ErrNotFoundSession        = errors.New("User session not found")
	ErrNotFoundUser           = errors.New("User not found")
	ErrDuplicateEmail         = errors.New("Duplicate email")
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
		user          User
		deactivatedAt pq.NullTime
	)

	err := m.DB.QueryRowContext(ctx, query, value).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deactivatedAt,
		&user.Email,
		&user.Password,
		pq.Array(&user.Roles),
		&user.ProfilName,
		&user.ShortId,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFoundUser
		default:
			return nil, err
		}
	}

	user.DeactivatedAt = deactivatedAt.Time

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

	var deactivatedAt pq.NullTime

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &deactivatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "user_account_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	user.DeactivatedAt = deactivatedAt.Time

	return nil
}

func (m Model) InsertOrUpdateUserSession(session *Session) error {
	query := `
		INSERT INTO "user_session" (
			id,
			deactivated_at,
			ip,
			agent,
			user_id
		) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET  
			deactivated_at = $2,
			ip = $3,
			agent = $4,
			user_id = $5
		RETURNING created_at, updated_at`

	args := []interface{}{
		session.ID,
		session.DeactivatedAt,
		session.IP,
		session.Agent,
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

func (m Model) GetSessionByID(id string) (*Session, error) {
	return m.getSessionBy("id", id)
}

func (m Model) getSessionBy(column string, value interface{}) (*Session, error) {
	query := fmt.Sprintf(`
		SELECT 
			id,
			created_at,
			updated_at,
			deactivated_at,
			ip,
			agent,
			user_id
		FROM user_session
		WHERE %s = $1`, column)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		session Session
	)

	err := m.DB.QueryRowContext(ctx, query, value).Scan(
		&session.ID,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.DeactivatedAt,
		&session.IP,
		&session.Agent,
		&session.UserID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFoundSession
		default:
			return nil, err
		}
	}

	return &session, nil
}

func (m Model) GetUserAndSession(userID, sessionID string) (*User, *Session, error) {
	query := `
		SELECT 
			u.id,
			u.created_at,
			u.updated_at,
			u.deactivated_at,
			u.email,
			u.password,
			u.roles,
			u.profil_name, 
			u.short_id,
			s.id,
			s.created_at,
			s.updated_at,
			s.deactivated_at,
			s.ip,
			s.agent,
			s.user_id
		FROM user_account AS u
		INNER JOIN user_session AS s
			ON s.user_id = u.id
		WHERE u.id = $1
		AND s.id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		session           Session
		user              User
		userDeactivatedAt pq.NullTime
	)

	err := m.DB.QueryRowContext(ctx, query, userID, sessionID).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&userDeactivatedAt,
		&user.Email,
		&user.Password,
		pq.Array(&user.Roles),
		&user.ProfilName,
		&user.ShortId,
		&session.ID,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.DeactivatedAt,
		&session.IP,
		&session.Agent,
		&session.UserID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil, ErrNotFoundUserAndSession
		default:
			return nil, nil, err
		}
	}

	user.DeactivatedAt = userDeactivatedAt.Time

	return &user, &session, nil
}
