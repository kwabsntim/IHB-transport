package services

import (
	"fmt"
	"ihb-transport/internal/models"
	"ihb-transport/internal/repository"
	"ihb-transport/utils"

	"github.com/google/uuid"
)

type reviewService struct {
	reviewRepo repository.ReviewsInterface
}

func NewReviewService(reviewRepo repository.ReviewsInterface) ReviewServiceInterface {
	return &reviewService{reviewRepo: reviewRepo}
}

// CreateReview validates and persists a review
func (s *reviewService) CreateReview(review *models.Reviews) error {
	if err := utils.ValidateRequired("client name", review.ClientName); err != nil {
		return err
	}
	if err := utils.ValidateRequired("content", review.Content); err != nil {
		return err
	}

	review.ID = uuid.Nil
	if err := s.reviewRepo.CreateReview(review); err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}
	return nil
}

// GetReviewByID fetches a review by its string UUID
func (s *reviewService) GetReviewByID(id string) (*models.Reviews, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid review id: %w", err)
	}
	rev, err := s.reviewRepo.FindByID(uid)
	if err != nil {
		return nil, fmt.Errorf("review not found: %w", err)
	}
	return rev, nil
}

// GetAllReviews returns all reviews
func (s *reviewService) GetAllReviews() ([]models.Reviews, error) {
	return s.reviewRepo.FindAll()
}

// DeleteReview deletes a review by ID
func (s *reviewService) DeleteReview(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid review id: %w", err)
	}
	// Check if review exists
	_, err = s.reviewRepo.FindByID(uid)
	if err != nil {
		return fmt.Errorf("review not found: %w", err)
	}
	return s.reviewRepo.Delete(uid)
}
