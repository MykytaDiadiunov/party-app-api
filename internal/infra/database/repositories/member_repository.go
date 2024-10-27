package repositories

import (
	"database/sql"
	"errors"
	"go-rest-api/internal/domain"
)

type member struct {
	PartyId uint64 `db:"party_id"`
	UserId  uint64 `db:"user_id"`
}

type memberRepository struct {
	db *sql.DB
}

type MemberRepository interface {
	Save(domainMember domain.Member) error
	Exists(domainMember domain.Member) error
	Delete(domainMember domain.Member) error
	FindByUserId(userId uint64) ([]domain.Member, error)
	FindByPartyId(partyId uint64) ([]domain.Member, error)
}

func NewMemberRepository(db *sql.DB) MemberRepository {
	return memberRepository{db: db}
}

func (m memberRepository) Save(domainMember domain.Member) error {
	memberModel := m.domainToModel(domainMember)
	sqlCommand := `INSERT INTO party_users (party_id, user_id) VALUES ($1, $2)`
	_, err := m.db.Exec(sqlCommand, memberModel.PartyId, memberModel.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (m memberRepository) FindByUserId(userId uint64) ([]domain.Member, error) {
	sqlCommand := `SELECT * FROM party_users WHERE user_id = $1`
	rows, err := m.db.Query(sqlCommand, userId)
	if err != nil {
		return []domain.Member{}, err
	}
	defer rows.Close()

	var members []domain.Member

	for rows.Next() {
		memberModel := member{}

		err := rows.Scan(
			&memberModel.PartyId,
			&memberModel.UserId,
		)
		if err != nil {
			return []domain.Member{}, err
		}

		members = append(members, m.modelToDomain(memberModel))
	}

	return members, nil
}

func (m memberRepository) FindByPartyId(partyId uint64) ([]domain.Member, error) {
	sqlCommand := `SELECT * FROM party_users WHERE party_id = $1`
	rows, err := m.db.Query(sqlCommand, partyId)
	if err != nil {
		return []domain.Member{}, err
	}
	defer rows.Close()

	var members []domain.Member

	for rows.Next() {
		memberModel := member{}

		err := rows.Scan(
			&memberModel.PartyId,
			&memberModel.UserId,
		)
		if err != nil {
			return []domain.Member{}, err
		}

		members = append(members, m.modelToDomain(memberModel))
	}

	return members, nil
}

func (m memberRepository) Delete(domainMemeber domain.Member) error {
	memberModel := m.domainToModel(domainMemeber)
	sqlCommand := `DELETE FROM party_users WHERE user_id = $1 AND party_id = $2`

	_, err := m.db.Exec(sqlCommand, memberModel.UserId, memberModel.PartyId)
	if err != nil {
		return err
	}

	return nil
}

func (m memberRepository) Exists(domainMember domain.Member) error {
	memberModel := m.domainToModel(domainMember)
	sqlCommand := `SELECT * FROM party_users WHERE party_id = $1 AND user_id = $2`
	rows, err := m.db.Query(sqlCommand, memberModel.PartyId, memberModel.UserId)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return errors.New("member does not exist")
	}
	return nil
}

func (m memberRepository) domainToModel(domainMember domain.Member) member {
	return member{
		PartyId: domainMember.PartyId,
		UserId:  domainMember.UserId,
	}
}

func (m memberRepository) modelToDomain(modelMember member) domain.Member {
	return domain.Member{
		PartyId: modelMember.PartyId,
		UserId:  modelMember.UserId,
	}
}
