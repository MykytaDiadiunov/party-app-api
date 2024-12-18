package app

import (
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/database/repositories"
	"go-rest-api/internal/infra/filesystem"
	"log"
	"strconv"
	"time"
)

type PartyService interface {
	Find(id uint64) (domain.Party, error)
	FindById(id uint64) (domain.Party, error)
	FindByCreatorId(creatorId uint64, page, limit int32) (domain.Parties, error)
	FindPartiesByLikerId(likerId uint64, page, limit int32) (domain.Parties, error)
	GetParties(page, limit int32) (domain.Parties, error)
	Save(party domain.Party) (domain.Party, error)
	Update(party domain.Party) (domain.Party, error)
	Delete(id uint64) error
}

type partyService struct {
	partyRepo         repositories.PartyRepository
	userService       UserService
	cloudinaryService *filesystem.CloudinaryService
}

func NewPartyService(partyRepo repositories.PartyRepository, cloudinaryService *filesystem.CloudinaryService, userServ UserService) PartyService {
	return partyService{
		partyRepo:         partyRepo,
		userService:       userServ,
		cloudinaryService: cloudinaryService,
	}
}

func (p partyService) Find(id uint64) (domain.Party, error) {
	party, err := p.FindById(id)
	if err != nil {
		return domain.Party{}, err
	}
	return party, nil
}

func (p partyService) FindById(id uint64) (domain.Party, error) {
	party, err := p.partyRepo.FindById(id)
	if err != nil {
		return domain.Party{}, err
	}
	return party, nil
}

func (p partyService) FindByCreatorId(creatorId uint64, page, limit int32) (domain.Parties, error) {
	parties, err := p.partyRepo.FindByCreatorId(creatorId, page, limit)
	if err != nil {
		return domain.Parties{}, err
	}
	return parties, nil
}

func (p partyService) FindPartiesByLikerId(likerId uint64, page, limit int32) (domain.Parties, error) {
	parties, err := p.partyRepo.FindPartiesByLikerId(likerId, page, limit)
	if err != nil {
		return domain.Parties{}, err
	}
	return parties, nil
}

func (p partyService) GetParties(page, limit int32) (domain.Parties, error) {
	parties, err := p.partyRepo.GetParties(page, limit)
	if err != nil {
		return domain.Parties{}, err
	}
	return parties, nil
}

func (p partyService) Save(party domain.Party) (domain.Party, error) {
	user, err := p.userService.FindById(party.CreatorId)
	if err != nil {
		log.Printf("Party service Save.UserFineById: %s", err)
		return domain.Party{}, err
	}

	amountToSpend := party.Price
	if amountToSpend < 10 {
		amountToSpend = 10
	}

	if party.Image != "" {
		imageFileName := "file_" + strconv.FormatInt(time.Now().UnixNano(), 32)

		imageUrl, err := p.cloudinaryService.SaveImageToCloudinary(party.Image, imageFileName)
		if err != nil {
			log.Printf("Party service Update.SaveImageToCloud: %s", err)
			return domain.Party{}, err
		}
		party.Image = imageUrl
	}

	createdParty, err := p.partyRepo.Save(party)
	if err != nil {
		log.Printf("Party service Save.RepoSave: %s", err)
		return domain.Party{}, err
	}

	_, err = p.userService.UpdateUserBalance(user, amountToSpend*(-1))
	if err != nil {
		log.Printf("Party service Save.UpdateUserBalance: %s", err)
		return domain.Party{}, err
	}

	return createdParty, nil
}

func (p partyService) Update(party domain.Party) (domain.Party, error) {
	partyFromDb, err := p.partyRepo.FindById(party.Id)
	if err != nil {
		log.Printf("Party service Update.FindByIdFromRepo: %s", err)
		return domain.Party{}, err
	}

	if party.Image == "" {
		party.Image = partyFromDb.Image
	}

	imageExists := partyFromDb.Image == party.Image

	if !imageExists {
		if partyFromDb.Image != "" {
			err = p.cloudinaryService.DeleteImage(partyFromDb.Image)
			if err != nil {
				return domain.Party{}, err
			}
		}

		imageFileName := "file_" + strconv.FormatInt(time.Now().UnixNano(), 32)

		imageUrl, err := p.cloudinaryService.SaveImageToCloudinary(party.Image, imageFileName)
		if err != nil {
			log.Printf("Party service Update.SaveImageToCloud: %s", err)
			return domain.Party{}, err
		}
		party.Image = imageUrl
	}

	updatedParty, err := p.partyRepo.Update(party)
	if err != nil {
		log.Printf("Party service Update.RepoUpdate: %s", err)
		return domain.Party{}, err
	}
	return updatedParty, nil
}

func (p partyService) Delete(id uint64) error {
	deletedParty, err := p.FindById(id)
	if err != nil {
		return err
	}

	if deletedParty.Image != "" {
		err = p.cloudinaryService.DeleteImage(deletedParty.Image)
		if err != nil {
			return err
		}
	}

	err = p.partyRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
