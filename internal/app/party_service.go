package app

import (
	"encoding/base64"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/database/repositories"
	"go-rest-api/internal/infra/filesystem"
	"log"
	"strconv"
)

type PartyService interface {
	Find(id uint64) (domain.Party, error)
	FindById(id uint64) (domain.Party, error)
	FindByCreatorId(creatorId uint64, page, limit int32) (domain.Parties, error)
	GetParties(page, limit int32) (domain.Parties, error)
	Save(party domain.Party) (domain.Party, error)
	Update(party domain.Party) (domain.Party, error)
	Delete(id uint64) error
}

type partyService struct {
	partyRepo    repositories.PartyRepository
	userService  UserService
	imageService filesystem.ImageStorageService
}

func NewPartyService(partyRepo repositories.PartyRepository, imageServ filesystem.ImageStorageService, userServ UserService) PartyService {
	return partyService{
		partyRepo:    partyRepo,
		imageService: imageServ,
		userService:  userServ,
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
		return domain.Party{}, err
	}

	amountToSpend := party.Price
	if amountToSpend < 10 {
		amountToSpend = 10
	}

	_, err = p.userService.UpdateUserBalance(user, amountToSpend*(-1))
	if err != nil {
		return domain.Party{}, err
	}

	createdParty, err := p.partyRepo.Save(party)
	if err != nil {
		return domain.Party{}, err
	}

	if party.Image != "" {
		partyWithNormalImg, err := p.Update(createdParty)
		if err != nil {
			err := p.Delete(createdParty.Id)
			if err != nil {
				return domain.Party{}, err
			}
			return domain.Party{}, err
		}
		return partyWithNormalImg, nil
	}
	return createdParty, nil
}

func (p partyService) Update(party domain.Party) (domain.Party, error) {
	currentParty, err := p.partyRepo.FindById(party.Id)
	if err != nil {
		return domain.Party{}, err
	}

	imageExist, imgErr := p.imageService.FileIsExist(party.Image)
	if imgErr != nil {
		return domain.Party{}, err
	}

	if !imageExist && currentParty.Image != "" {
		id := strconv.FormatUint(party.Id, 10)
		creatorId := strconv.FormatUint(currentParty.CreatorId, 10)
		imageFileName := "party_" + id + "_by_user_" + creatorId + ".jpg"
		err = p.saveImage(imageFileName, party.Image)
		if err != nil {
			return domain.Party{}, err
		}
		party.Image = imageFileName
	}

	party.CreatorId = currentParty.CreatorId
	updatedParty, err := p.partyRepo.Update(party)
	if err != nil {
		return domain.Party{}, err
	}
	return updatedParty, nil
}

func (p partyService) Delete(id uint64) error {
	deletedParty, err := p.FindById(id)
	if err != nil {
		return err
	}

	err = p.imageService.RemoveImage(deletedParty.Image)
	if err != nil {
		log.Println("Error in delete party img")
	}

	err = p.partyRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (p partyService) saveImage(imageFileName, imageStringData string) error {
	imageData, err := base64.StdEncoding.DecodeString(imageStringData)
	if err != nil {
		return err
	}

	err = p.imageService.SaveImage(imageFileName, imageData)
	if err != nil {
		return err
	}

	return nil
}
