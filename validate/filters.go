package validate

import "gorm.io/gorm"

func NoFilter(db *gorm.DB) *gorm.DB {
	return db
}

func NewFilter(s string, a ...any) FilterFunc {
	return func(d *gorm.DB) {
		d.Where(s, a...)
	}
}

// func JoinFilters(f1 FilterFunc, fxs ...FilterFunc) FilterFunc {
// 	return func(db *gorm.DB) *gorm.DB {
// 		db = f1(db)
// 		for _, fx := range fxs {
// 			db = fx(db)
// 		}
// 		return db
// 	}
// }
