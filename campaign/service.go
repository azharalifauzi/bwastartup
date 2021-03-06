package campaign

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignDetailInput, inputData UpdateCampaignInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)

		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	campaigns, err := s.repository.FindAll()

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)

	if err != nil {
		return campaign, err
	}

	if campaign.ID == 0 {
		errMessage := fmt.Sprintf("Campaign with that ID of %d was not found", input.ID)
		return campaign, errors.New(errMessage)
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID
	campaign.Perks = input.Perks

	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugCandidate)

	newCampagin, err := s.repository.Save(campaign)

	if err != nil {
		return newCampagin, err
	}

	return newCampagin, nil

}

func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData UpdateCampaignInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(inputID.ID)

	if err != nil {
		return campaign, err
	}

	if campaign.ID == 0 {
		return campaign, errors.New("404")
	}

	if inputData.User.ID != campaign.UserID {
		return campaign, errors.New("401")
	}

	c := reflect.ValueOf(&campaign).Elem()
	ri := reflect.ValueOf(&inputData).Elem()
	typeOfRi := ri.Type()

	for i := 0; i < typeOfRi.NumField(); i++ {
		value := ri.Field(i).Interface()
		field := typeOfRi.Field(i).Name

		str, ok := value.(string)

		if ok && len(str) > 0 {
			c.FieldByName(field).SetString(str)
		}

		integer, ok := value.(int)

		if ok && integer > 0 {
			c.FieldByName(field).SetInt(int64(integer))
		}
	}

	updatedCampaign, err := s.repository.Update(campaign)

	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {

	campaign, err := s.repository.FindByID(input.CampaignID)

	if err != nil {
		return CampaignImage{}, err
	}

	if input.User.ID != campaign.UserID {
		return CampaignImage{}, errors.New("401")
	}

	isPrimary := 0

	if input.IsPrimary {
		isPrimary = 1

		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)

		if err != nil {
			return CampaignImage{}, err
		}
	}

	campaignImage := CampaignImage{}
	campaignImage.FileName = fileLocation
	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = isPrimary

	newCampaignImage, err := s.repository.CreateImage(campaignImage)

	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}
