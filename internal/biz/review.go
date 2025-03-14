package biz

import (
	"context"

	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"

	"github.com/go-kratos/kratos/v2/log"
)

type ReviewRepo interface {
	SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewsByOrderID(ctx context.Context, orderID int64) ([]*model.ReviewInfo, error)
}

type ReviewUsecase struct {
	repo ReviewRepo
	log  *log.Helper
}

func NewReviewUsecase(repo ReviewRepo, logger log.Logger) *ReviewUsecase {
	return &ReviewUsecase{repo: repo, log: log.NewHelper(logger)}
}

// 业务逻辑
func (uc *ReviewUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Infof("[biz] CreateReview: %v", review)
	// 参数基础校验

	// 参数逻辑校验
	reviews, err := uc.repo.GetReviewsByOrderID(ctx, review.OrderID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("GetReviewsByOrderID failed: %v", err)
		return nil, v1.ErrorDbFailed("GetReviewsByOrderID failed")
	}
	if len(reviews) > 0 {
		uc.log.WithContext(ctx).Errorf("review already exists")
		return nil, v1.ErrorReviewExisted("order has already review")
	}

	// 生成reviewID
	review.ReviewID = snowflake.GenID()
	return uc.repo.SaveReview(ctx, review)
}
