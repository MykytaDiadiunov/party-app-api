package repositories

import (
	"database/sql"
	"go-rest-api/internal/domain"
)

type user struct {
	Id       uint64 `db:"id, omitempty"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Points   int32  `db:"points"`
}

type UserRepository interface {
	UpdateUserBalance(user domain.User, amount int32) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
	FindById(id uint64) (domain.User, error)
	Save(user domain.User) (domain.User, error)
	Delete(id uint64) error
}
type userRepository struct {
	db *sql.DB
}

func NewUserRepository(database *sql.DB) UserRepository {
	return &userRepository{db: database}
}

func (ur userRepository) UpdateUserBalance(user domain.User, amount int32) (domain.User, error) {
	userModel := ur.domainToModel(user)
	sqlCommand := `UPDATE users SET points = points + $1 WHERE id = $2`

	_, err := ur.db.Exec(sqlCommand, amount, userModel.Id)
	if err != nil {
		return domain.User{}, err
	}

	userModel.Points = userModel.Points + amount

	return ur.modelToDomain(userModel), nil
}

func (ur userRepository) FindByEmail(email string) (domain.User, error) {
	userModel := user{}
	sqlCommand := `SELECT * FROM users WHERE email=$1`

	err := ur.db.QueryRow(
		sqlCommand,
		email,
	).Scan(
		&userModel.Id,
		&userModel.Name,
		&userModel.Email,
		&userModel.Password,
		&userModel.Points,
	)
	if err != nil {
		return domain.User{}, err
	}

	return ur.modelToDomain(userModel), nil
}

func (ur userRepository) FindById(id uint64) (domain.User, error) {
	userModel := user{}
	sqlCommand := `SELECT * FROM users WHERE id=$1`
	err := ur.db.QueryRow(
		sqlCommand,
		id,
	).Scan(
		&userModel.Id,
		&userModel.Name,
		&userModel.Email,
		&userModel.Password,
		&userModel.Points,
	)
	if err != nil {
		return domain.User{}, err
	}

	return ur.modelToDomain(userModel), nil
}

func (ur userRepository) Save(user domain.User) (domain.User, error) {
	userModel := ur.domainToModel(user)
	sqlCommand := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, points`

	err := ur.db.QueryRow(
		sqlCommand,
		userModel.Name,
		userModel.Email,
		userModel.Password,
	).Scan(
		&userModel.Id,
		&userModel.Points,
	)
	if err != nil {
		return domain.User{}, err
	}

	return ur.modelToDomain(userModel), nil
}

func (ur userRepository) Delete(id uint64) error {
	sqlCommand := `DELETE FROM users WHERE id=$1`
	_, err := ur.db.Exec(sqlCommand, id)
	if err != nil {
		return err
	}
	return nil
}

func (ur userRepository) modelToDomain(u user) domain.User {
	return domain.User{
		Id:       u.Id,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Points:   u.Points,
	}
}

func (ur userRepository) domainToModel(u domain.User) user {
	return user{
		Id:       u.Id,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Points:   u.Points,
	}
}
