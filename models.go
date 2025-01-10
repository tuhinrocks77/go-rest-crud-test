package main

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint64 `gorm:"primary_key;autoIncrement:true"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
}

type Status string

const (
	Pending    Status = "pending"
	InProgress Status = "in-progress"
	Completed  Status = "completed"
)

type Task struct {
	// https://medium.com/@amrilsyaifa_21001/how-to-use-uuid-in-gorm-golang-74be997d7087
	BaseModel
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	// TODO: implement default value & remove logic from controller
	Status Status
}

// TODO add total count
func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page <= 0 {
			page = 1
		}

		limit, _ := strconv.Atoi(q.Get("limit"))
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}
