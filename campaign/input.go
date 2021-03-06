package campaign

import "bwastartup/user"

type GetCampaignDetailInput struct {
	ID int `uri:"id" binding:"required"`
}

type CreateCampaignInput struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	Perks            string `json:"perks" binding:"required"`
	GoalAmount       int    `json:"goal_amount" binding:"required"`
	User             user.User
}

type UpdateCampaignInput struct {
	Name             string `json:"name" binding:"omitempty"`
	ShortDescription string `json:"short_description" binding:"omitempty"`
	Description      string `json:"description" binding:"omitempty"`
	Perks            string `json:"perks" binding:"omitempty"`
	GoalAmount       int    `json:"goal_amount" binding:"omitempty"`
	User             user.User
}

type CreateCampaignImageInput struct {
	CampaignID int  `form:"campaign_id" binding:"required"`
	IsPrimary  bool `form:"is_primary"`
	User       user.User
}
