package repositories

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"go-rest-api/internal/domain"
)

type session struct {
	UserId uint64    `db:"user_id"`
	UUID   uuid.UUID `db:"uuid"`
}

type SessionRepository interface {
	Save(sess domain.Session) error
	Exists(sess domain.Session) error
	Delete(sess domain.Session) error
}

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (sr sessionRepository) Save(sess domain.Session) error {
	s := sr.domainToModel(sess)
	sqlCommand := `INSERT INTO sessions (uuid, user_id) VALUES ($1, $2)`
	_, err := sr.db.Exec(sqlCommand, s.UUID, s.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (sr sessionRepository) Exists(sess domain.Session) error {
	s := sr.domainToModel(sess)
	sqlCommand := `SELECT * FROM sessions WHERE uuid = $1 AND user_id = $2`
	rows, err := sr.db.Query(sqlCommand, s.UUID, sess.UserId)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return errors.New("session does not exist")
	}
	return nil
}

func (sr sessionRepository) Delete(sess domain.Session) error {
	s := sr.domainToModel(sess)
	sqlCommand := `DELETE FROM sessions WHERE uuid = $1 AND user_id = $2`
	_, err := sr.db.Exec(sqlCommand, s.UUID, sess.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (sr sessionRepository) domainToModel(sess domain.Session) session {
	return session{
		UserId: sess.UserId,
		UUID:   sess.UUID,
	}
}
