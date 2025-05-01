package entity

type Competition struct {
	CompetitionID   int    `json:"competition_id" gorm:"type:int;primaryKey"`
	CompetitionName string `json:"competition_name" gorm:"type:varchar(70)"`
	Description     string `json:"description" gorm:"type:text"`

	Registrations  []Registration `gorm:"foreignKey:CompetitionID"`
	TeamProgresses []TeamProgress `gorm:"foreignKey:CompetitionID"`
	Announcements  []Announcement `gorm:"foreignKey:CompetitionID"`
}
