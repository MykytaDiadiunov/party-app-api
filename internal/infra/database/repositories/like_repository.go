package repositories

import (
	"database/sql"
	"errors"
	"go-rest-api/internal/domain"
)

type like struct {
	LikedId uint64 `db:"liked_id"`
	LikerId uint64 `db:"liker_id"`
}

type likeRepository struct {
	db *sql.DB
}

type LikeRepository interface {
	Save(domainLike domain.Like) error
	FindByLikedId(likedId uint64) ([]domain.Like, error)
	FindByLikerId(likerId uint64) ([]domain.Like, error)
	Delete(domainLike domain.Like) error
	Exists(domainLike domain.Like) error
}

func NewLikeRepository(db *sql.DB) LikeRepository {
	return likeRepository{
		db: db,
	}
}

func (l likeRepository) Save(domainLike domain.Like) error {
	likeModel := l.DomainToModel(domainLike)
	sqlCommand := `INSERT INTO likes (liked_id, liker_id) VALUES ($1, $2)`

	_, err := l.db.Exec(sqlCommand, likeModel.LikedId, likeModel.LikerId)
	if err != nil {
		return err
	}

	return nil
}

func (l likeRepository) FindByLikedId(likedId uint64) ([]domain.Like, error) {
	sqlCommand := `SELECT * FROM likes WHERE liked_id = $1`
	rows, err := l.db.Query(sqlCommand, likedId)
	if err != nil {
		return []domain.Like{}, err
	}
	defer rows.Close()

	likes := []domain.Like{}

	for rows.Next() {
		likeModel := like{}

		err = rows.Scan(
			&likeModel.LikedId,
			&likeModel.LikerId,
		)
		if err != nil {
			return []domain.Like{}, err
		}

		likes = append(likes, l.ModelToDomain(likeModel))
	}

	return likes, nil
}

func (l likeRepository) FindByLikerId(likerId uint64) ([]domain.Like, error) {
	sqlCommand := `SELECT * FROM likes WHERE liker_id = $1`
	rows, err := l.db.Query(sqlCommand, likerId)
	if err != nil {
		return []domain.Like{}, err
	}
	defer rows.Close()

	likes := []domain.Like{}

	for rows.Next() {
		likeModel := like{}

		err = rows.Scan(
			&likeModel.LikedId,
			&likeModel.LikerId,
		)
		if err != nil {
			return []domain.Like{}, err
		}

		likes = append(likes, l.ModelToDomain(likeModel))
	}

	return likes, nil
}

func (l likeRepository) Delete(domainLike domain.Like) error {
	likeModel := l.DomainToModel(domainLike)
	sqlCommand := `DELETE FROM likes WHERE liked_id = $1 AND liker_id = $2`

	_, err := l.db.Exec(sqlCommand, likeModel.LikedId, likeModel.LikerId)
	if err != nil {
		return err
	}

	return nil
}

func (l likeRepository) Exists(domainLike domain.Like) error {
	likeModel := l.DomainToModel(domainLike)
	sqlCommand := `SELECT * FROM likes WHERE liked_id = $1 AND liker_id = $2`

	rows, err := l.db.Query(sqlCommand, likeModel.LikedId, likeModel.LikerId)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return errors.New("like does not exist")
	}

	return nil
}

func (l likeRepository) DomainToModel(domainLike domain.Like) like {
	return like{
		LikedId: domainLike.LikedId,
		LikerId: domainLike.LikerId,
	}
}

func (l likeRepository) ModelToDomain(likeModel like) domain.Like {
	return domain.Like{
		LikedId: likeModel.LikedId,
		LikerId: likeModel.LikerId,
	}
}
