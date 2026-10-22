package models

import (
	"encoding/json"
	"time"
)

type HasUserId struct {
	UserId int `json:"user_id" gorm:"index;not nul"`
}

type HasBalanceId struct {
	BalanceId int `json:"balance_id" gorm:"index;not null"`
}

type HasIsActive struct {
	IsActive bool `json:"is_active" gorm:"index;not null;default:true"`
}

func (h HasIsActive) GetIsActive() bool {
	return h.IsActive
}

type MyDate struct {
	time.Time `json:"date" gorm:"index"`
}

func (d *MyDate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	date, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = date
	return nil
}

func (d *MyDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format("2006-01-02"))
}
