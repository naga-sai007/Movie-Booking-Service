package models

type Theatre struct {
	TheatreID   int64   `gorm:"type:BIGINT(7);column:theatre_id;primarykey" json:"theatreId"`
	TheatreName string  `gorm:"column:theatre_name" json:"theatreName"`
	Address     string  `gorm:"column:address" json:"address"`
	Screens     string  `gorm:"column:screens" json:"screens"`
	Shows       []Shows `gorm:"foreignkey:theatre_id" json:"shows"`
}
