package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type HasUserId struct {
	UserId int `json:"user_id" gorm:"index;not nul"`
}

type HasWalletId struct {
	WalletId int `json:"wallet_id" gorm:"index;not null"`
}

type HasIsActive struct {
	IsActive bool `json:"is_active" gorm:"index;not null;default:true"`
}

func (h HasIsActive) GetIsActive() bool {
	return h.IsActive
}

type MyDate struct {
	time.Time
}

func (d MyDate) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

func (d *MyDate) Scan(value any) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into MyDate", value)
	}
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
