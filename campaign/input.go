package campaign

type CreateCampaignInput struct {
	Name             string `json:"name" validate:"required"`
	ShortDescription string `json:"short_description" validate:"required"`
	Description      string `json:"description" validate:"required"`
	Perks            string `json:"perks" validate:"required"`
	GoalAmount       int    `json:"goal_amount" validate:"required"`
	UserId           int
}

type CreateCampaignImageInput struct {
	CampaignId int `form:"campaign_id" validate:"required"`
	IsPrimary  int `form:"is_primary" validate:"required"`
}