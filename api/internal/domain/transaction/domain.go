package transaction

import (
	"api/internal/model"
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type domain struct {
	transactionRepo repoProvider
}

type repoProvider interface {
	Create(ctx context.Context, in *model.Transaction) error
}

// New is to initialize transaction domain instance.
func New(db *gorm.DB) *domain {
	repo := newRepository(db)
	return &domain{
		transactionRepo: repo,
	}
}

// Create is to create a single transaction entry
func (d *domain) Create(ctx context.Context, in *model.CreateTransaction) (*model.CommonResponse, error) {
	err := d.transactionRepo.Create(ctx, &model.Transaction{
		Id:          ksuid.New().String(),
		UserId:      in.UserId,
		Status:      "pending",
		Description: in.Description,
		Remarks:     "",
		Amount:      in.Amount,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	return &model.CommonResponse{
		Status:  true,
		Message: "ok",
	}, err
}
