package campaign

import "strings"

type CampaignFormatter struct {
	Id               int    `json:"id"`
	UserId           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageUrl         string `json:"image_url"`
	Slug             string `json:"slug"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{}
	campaignFormatter.Id = campaign.Id
	campaignFormatter.UserId = campaign.UserId
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.Slug = campaign.Slug
	campaignFormatter.CurrentAmount = campaign.CurrentAmount
	campaignFormatter.ImageUrl = ""

	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageUrl = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	if len(campaigns) == 0 {
		return []CampaignFormatter{}
	}
	var campaignsFormatter []CampaignFormatter

	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter
}

type campaignDetailFormatter struct {
	Id               int      `json:"id"`
	Name             string   `json:"name"`
	ShortDescription string   `json:"short_description"`
	Description      string   `json:"description"`
	ImageUrl         string   `json:"image_url"`
	GoalAmount       int      `json:"goal_amount"`
	CurrentAmount    int      `json:"current_amount"`
	UserId           int      `json:"user_id"`
	Slug             string   `json:"slug"`
	Perks            []string `json:"perks"`
	User campaignDetailUserFormat `json:"user"`
	Images []campaignDetailImagesFormat `json:"images"`
}

type campaignDetailUserFormat struct {
	Name string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

type campaignDetailImagesFormat struct {
	ImageUrl string `json:"image_url"`
	IsPrimary bool `json:"is_primary"`
}

func FormatCampaignDetail(campaign Campaign) campaignDetailFormatter {
	campaignDetailFormat := campaignDetailFormatter{}
	campaignDetailFormat.Id = campaign.Id
	campaignDetailFormat.Name = campaign.Name
	campaignDetailFormat.ShortDescription = campaign.ShortDescription
	campaignDetailFormat.Description = campaign.Description
	campaignDetailFormat.GoalAmount = campaign.GoalAmount
	campaignDetailFormat.CurrentAmount = campaign.CurrentAmount
	campaignDetailFormat.UserId = campaign.UserId
	campaignDetailFormat.Slug = campaign.Slug
	campaignDetailFormat.ImageUrl = ""
	
	// MAP PERKS TO CONFIG
	if len(campaign.CampaignImages) > 0 {
		campaignDetailFormat.ImageUrl = campaign.CampaignImages[0].FileName
	}

	listOfPerks := strings.Split(campaign.Perks, ",")
	campaignDetailFormat.Perks = listOfPerks

	// MAP USER TO CONFIG
	campaignUser := campaign.User
	campaignUserFormat := campaignDetailUserFormat{}
	campaignUserFormat.Name = campaignUser.Name
	campaignUserFormat.AvatarUrl = campaignUser.AvatarFileName

	// MAP IMAGES TO CONFIG
	campaignImages := campaign.CampaignImages
	campaignImagesFormat := campaignDetailImagesFormat{}
	campaignImagesFormatArr := []campaignDetailImagesFormat{}
	for _, images := range campaignImages {
		campaignImagesFormat.ImageUrl = images.FileName
		if images.IsPrimary == 1 {
			campaignImagesFormat.IsPrimary = true
		} else {
			campaignImagesFormat.IsPrimary = false
		}
		campaignImagesFormatArr = append(campaignImagesFormatArr,campaignImagesFormat)
	}

	campaignDetailFormat.User = campaignUserFormat
	campaignDetailFormat.Images = campaignImagesFormatArr

	return campaignDetailFormat

}