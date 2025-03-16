package service

import (
	"context"
	"fmt"

	pb "review-service/api/review/v1"
	"review-service/internal/biz"
	"review-service/internal/data/model"
)

type ReviewService struct {
	pb.UnimplementedReviewServer
	uc *biz.ReviewUsecase
}

func NewReviewService(uc *biz.ReviewUsecase) *ReviewService {
	return &ReviewService{uc: uc}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewReply, error) {
	fmt.Printf("[service] CreateReview: %v\n", req)
	var anonymous int32
	if req.Anonymous {
		anonymous = 1
	}
	review, err := s.uc.CreateReview(ctx, &model.ReviewInfo{
		UserID:       req.UserID,
		OrderID:      req.OrderID,
		Content:      req.Content,
		Score:        req.Score,
		ServiceScore: req.ServiceScore,
		ExpressScore: req.ExpressScore,
		PicInfo:      req.PicInfo,
		VideoInfo:    req.VideoInfo,
		Anonymous:    anonymous,
		Status:       0,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateReviewReply{ReviewID: review.ReviewID}, nil
}

func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {
	fmt.Printf("[service] GetReview: %v\n", req)
	review, err := s.uc.GetReview(ctx, req.ReviewID)
	if err != nil {
		return nil, err
	}
	return &pb.GetReviewReply{
		Review: &pb.ReviewInfo{
			ReviewID:     review.ReviewID,
			UserID:       review.UserID,
			OrderID:      review.OrderID,
			Content:      review.Content,
			Score:        review.Score,
			ServiceScore: review.ServiceScore,
			ExpressScore: review.ExpressScore,
			PicInfo:      review.PicInfo,
			VideoInfo:    review.VideoInfo,
			Status:       review.Status,
		},
	}, nil
}

func (s *ReviewService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	fmt.Printf("[service] ReplyReview: %v\n", req)
	reply, err := s.uc.CreateReply(ctx, &biz.ReplyParam{
		ReviewID:  req.ReviewID,
		Content:   req.Content,
		StoreID:   req.StoreID,
		PicInfo:   req.PicInfo,
		VideoInfo: req.VideoInfo,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReplyReviewReply{ReplyID: reply.ReplyID}, nil
}

func (s *ReviewService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error) {
	return &pb.AppealReviewReply{}, nil
}

func (s *ReviewService) AuditReview(ctx context.Context, req *pb.AuditReviewRequest) (*pb.AuditReviewReply, error) {
	return &pb.AuditReviewReply{}, nil
}

func (s *ReviewService) AuditAppeal(ctx context.Context, req *pb.AuditAppealRequest) (*pb.AuditAppealReply, error) {
	return &pb.AuditAppealReply{}, nil
}

func (s *ReviewService) ListReviewByUserID(ctx context.Context, req *pb.ListReviewByUserIDRequest) (*pb.ListReviewByUserIDReply, error) {
	return &pb.ListReviewByUserIDReply{}, nil
}
