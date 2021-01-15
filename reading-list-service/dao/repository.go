package dao

import "github.com/KumKeeHyun/medium-rare/reading-list-service/domain"

type ReadingRepository interface {
	FindViewedsByUserID(userID int) ([]domain.Viewed, error)
	// first -> create, on conflict -> update
	SaveViewed(viewed domain.Viewed) (domain.Viewed, error)

	FindSavedsByUserID(userID int) ([]domain.Saved, error)
	// only create
	SaveSaved(saved domain.Saved) (domain.Saved, error)
}
