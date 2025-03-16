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
	GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error)
	SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error)
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

func (uc *ReviewUsecase) GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Infof("[biz] GetReview: %v", reviewID)
	return uc.repo.GetReview(ctx, reviewID)
}

func (uc *ReviewUsecase) CreateReply(ctx context.Context, param *ReplyParam) (*model.ReviewReplyInfo, error) {
	uc.log.WithContext(ctx).Infof("[biz] CreateReply: %v", param)
	// 参数逻辑校验
	review, err := uc.repo.GetReview(ctx, param.ReviewID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("GetReplyByReviewID failed: %v", err)
		return nil, v1.ErrorDbFailed("GetReplyByReviewID failed")
	}
	if review.HasReply == 1 {
		uc.log.WithContext(ctx).Errorf("reply already exists")
		return nil, v1.ErrorReplyExisted("review has already reply")
	}
	// 水平越权校验
	if review.StoreID != param.StoreID {
		uc.log.WithContext(ctx).Errorf("storeID not match")
		return nil, v1.ErrorStoreNotMatch("storeID not match")
	}

	reply := &model.ReviewReplyInfo{
		ReplyID:   snowflake.GenID(),
		ReviewID:  param.ReviewID,
		Content:   param.Content,
		StoreID:   param.StoreID,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	}
	return uc.repo.SaveReply(ctx, reply)
}
