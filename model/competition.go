package model

type GetAllCompetitionsResponse struct {
	CompetitionID   int    `json:"competition_id"`
	CompetitionName string `json:"competition_name"`
	Description     string `json:"description"`
}
