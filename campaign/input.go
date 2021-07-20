package campaign

type CreateCampaignInput struct {
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description"`
	Perks            string    `json:"perks"`
	GoalAmount int `json:"goal_amount"`
	UserId             int
}