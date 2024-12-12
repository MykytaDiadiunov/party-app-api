package repositories

import (
	"database/sql"
	"go-rest-api/internal/domain"
	"time"
)

type party struct {
	Id          uint64    `db:"id, omitempty"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Image       string    `db:"image"`
	Price       int32     `db:"price"`
	StartDate   time.Time `db:"start_date"`
	CreatorId   uint64    `db:"creator_id"`
}

type PartyRepository interface {
	FindById(id uint64) (domain.Party, error)
	FindByCreatorId(creatorId uint64, page, limit int32) (domain.Parties, error)
	FindPartiesByLikerId(likerId uint64, page, limit int32) (domain.Parties, error)
	GetParties(page, limit int32) (domain.Parties, error)
	Save(party domain.Party) (domain.Party, error)
	Update(party domain.Party) (domain.Party, error)
	Delete(id uint64) error
}

type partyRepository struct {
	db *sql.DB
}

func NewPartyRepository(db *sql.DB) PartyRepository {
	return partyRepository{db: db}
}

func (p partyRepository) FindById(id uint64) (domain.Party, error) {
	partyModel := party{}
	sqlCommand := `SELECT id, title, description, image, price, start_date, creator_id FROM parties WHERE id = $1;`
	err := p.db.QueryRow(
		sqlCommand,
		id,
	).Scan(
		&partyModel.Id,
		&partyModel.Title,
		&partyModel.Description,
		&partyModel.Image,
		&partyModel.Price,
		&partyModel.StartDate,
		&partyModel.CreatorId,
	)
	if err != nil {
		return domain.Party{}, err
	}
	return p.modelToDomain(partyModel), nil
}

func (p partyRepository) FindByCreatorId(creatorId uint64, page, limit int32) (domain.Parties, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var parties []domain.Party

	sqlCommand := `SELECT id, title, description, image, price, start_date, creator_id FROM parties 
	WHERE creator_id = $1 ORDER BY created_date DESC LIMIT $2 OFFSET $3;`
	rows, err := p.db.Query(sqlCommand, creatorId, limit, offset)
	if err != nil {
		return domain.Parties{}, err
	}
	defer rows.Close()

	for rows.Next() {
		partyModel := party{}
		err := rows.Scan(
			&partyModel.Id,
			&partyModel.Title,
			&partyModel.Description,
			&partyModel.Image,
			&partyModel.Price,
			&partyModel.StartDate,
			&partyModel.CreatorId,
		)
		if err != nil {
			return domain.Parties{}, err
		}
		parties = append(parties, p.modelToDomain(partyModel))
	}

	var total uint64
	totalSqlCommand := `SELECT COUNT(*) FROM parties WHERE creator_id = $1;`
	err = p.db.QueryRow(totalSqlCommand, creatorId).Scan(&total)
	if err != nil {
		return domain.Parties{}, err
	}
	var pages int32
	if total > 0 {
		pages = (int32(total) + limit - 1) / limit
	}

	return domain.Parties{
		Parties:     parties,
		Total:       total,
		CurrentPage: page,
		LastPage:    pages,
	}, nil
}

func (p partyRepository) FindPartiesByLikerId(likerId uint64, page, limit int32) (domain.Parties, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	sqlCommand := `select distinct 
	parties.id, parties.title, parties.description, parties.image, parties.price, parties.start_date, parties.creator_id, parties.created_date
	from likes inner join parties on parties.creator_id = likes.liked_id where likes.liker_id = $1 order by parties.created_date desc LIMIT $2 OFFSET $3;`
	rows, err := p.db.Query(sqlCommand, likerId, limit, offset)
	if err != nil {
		return domain.Parties{}, err
	}

	defer rows.Close()
	var parties []domain.Party
	var cork any
	for rows.Next() {
		var party party
		err := rows.Scan(
			&party.Id,
			&party.Title,
			&party.Description,
			&party.Image,
			&party.Price,
			&party.StartDate,
			&party.CreatorId,
			&cork,
		)
		if err != nil {
			return domain.Parties{}, err
		}
		parties = append(parties, p.modelToDomain(party))
	}

	var total uint64
	totalSqlCommand := `SELECT COUNT(DISTINCT parties.id) 
	FROM likes INNER JOIN parties ON parties.creator_id = likes.liked_id WHERE likes.liker_id = $1;`
	err = p.db.QueryRow(totalSqlCommand, likerId).Scan(&total)
	if err != nil {
		return domain.Parties{}, err
	}
	var pages int32
	if total > 0 {
		pages = (int32(total) + limit - 1) / limit
	}

	return domain.Parties{
		Parties:     parties,
		Total:       total,
		CurrentPage: page,
		LastPage:    pages,
	}, nil
}

func (p partyRepository) GetParties(page, limit int32) (domain.Parties, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	sqlCommand := `SELECT id, title, description, image, price, start_date, creator_id FROM parties ORDER BY created_date DESC LIMIT $1 OFFSET $2`
	rows, err := p.db.Query(sqlCommand, limit, offset)
	if err != nil {
		return domain.Parties{}, err
	}
	defer rows.Close()

	var parties []domain.Party
	for rows.Next() {
		var party party
		err := rows.Scan(
			&party.Id,
			&party.Title,
			&party.Description,
			&party.Image,
			&party.Price,
			&party.StartDate,
			&party.CreatorId,
		)
		if err != nil {
			return domain.Parties{}, err
		}
		parties = append(parties, p.modelToDomain(party))
	}
	var total uint64
	totalSqlCommand := `SELECT COUNT(*) FROM parties;`
	err = p.db.QueryRow(totalSqlCommand).Scan(&total)
	if err != nil {
		return domain.Parties{}, err
	}
	var pages int32
	if total > 0 {
		pages = (int32(total) + limit - 1) / limit
	}

	return domain.Parties{
		Parties:     parties,
		Total:       total,
		CurrentPage: page,
		LastPage:    pages,
	}, nil
}

func (p partyRepository) Save(party domain.Party) (domain.Party, error) {
	partyModel := p.domainToModel(party)

	sqlCommand := `INSERT INTO parties(
                  title, 
                  description, 
                  image, 
                  price, 
                  start_date, 
                  creator_id
			  ) VALUES($1, $2, $3, $4, $5, $6) RETURNING id`

	err := p.db.QueryRow(
		sqlCommand,
		partyModel.Title,
		partyModel.Description,
		partyModel.Image,
		partyModel.Price,
		partyModel.StartDate,
		partyModel.CreatorId,
	).Scan(&partyModel.Id)
	if err != nil {
		return domain.Party{}, err
	}

	return p.modelToDomain(partyModel), nil
}

func (p partyRepository) Update(party domain.Party) (domain.Party, error) {
	partyModel := p.domainToModel(party)
	sqlCommand := `UPDATE parties SET 
                 title = $1,
                 description = $2,
                 image = $3,
                 start_date = $4 WHERE id = $5`

	_, err := p.db.Exec(
		sqlCommand,
		partyModel.Title,
		partyModel.Description,
		partyModel.Image,
		partyModel.StartDate,
		partyModel.Id,
	)
	if err != nil {
		return domain.Party{}, err
	}

	newParty, err := p.FindById(partyModel.Id)
	if err != nil {
		return domain.Party{}, err
	}

	return newParty, nil
}

func (p partyRepository) Delete(id uint64) error {
	sqlCommand := `DELETE FROM parties WHERE id = $1`
	_, err := p.db.Exec(sqlCommand, id)
	if err != nil {
		return err
	}
	return nil
}

func (p partyRepository) domainToModel(domainParty domain.Party) party {
	return party{
		Id:          domainParty.Id,
		Title:       domainParty.Title,
		Description: domainParty.Description,
		Image:       domainParty.Image,
		Price:       domainParty.Price,
		StartDate:   domainParty.StartDate,
		CreatorId:   domainParty.CreatorId,
	}
}

func (p partyRepository) modelToDomain(modelParty party) domain.Party {
	return domain.Party{
		Id:          modelParty.Id,
		Title:       modelParty.Title,
		Description: modelParty.Description,
		Image:       modelParty.Image,
		Price:       modelParty.Price,
		StartDate:   modelParty.StartDate,
		CreatorId:   modelParty.CreatorId,
	}
}
