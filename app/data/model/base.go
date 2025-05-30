package model

import (
	"dex/app/data/utility"

	"gorm.io/gorm"
)

type Base struct {
	gorm.Model
}

func (base *Base) BeforeCreate(tx *gorm.DB) error {
	base.ID = uint(utility.Sequence.Next())
	return nil
}
