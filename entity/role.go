package entity

type Role struct {
	RoleID   int    `json:"role_id" gorm:"type:int;primaryKey"`
	RoleName string `json:"role_name" gorm:"type:varchar(20)"`
	Users    []User `gorm:"foreignKey:RoleID"`
}
