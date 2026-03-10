package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"finapp/internal/model"
	"finapp/internal/repository"
)

type categoryService struct {
	categories repository.CategoryRepository
}

func NewCategoryService(categories repository.CategoryRepository) CategoryService {
	return &categoryService{categories: categories}
}

func (s *categoryService) Create(ctx context.Context, userID uuid.UUID, req model.CategoryRequest) (*model.Category, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	cat := &model.Category{
		ID:       uuid.New(),
		UserID:   &userID,
		Name:     req.Name,
		Color:    req.Color,
		Icon:     req.Icon,
		ParentID: req.ParentID,
		IsSystem: false,
	}
	if err := s.categories.Create(ctx, cat); err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}
	return cat, nil
}

func (s *categoryService) List(ctx context.Context, userID uuid.UUID) ([]model.Category, error) {
	return s.categories.FindByUserID(ctx, userID)
}

func (s *categoryService) GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	cat, err := s.categories.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrNotFound
	}
	return cat, nil
}

func (s *categoryService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req model.CategoryRequest) (*model.Category, error) {
	cat, err := s.categories.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrNotFound
	}
	if cat.IsSystem {
		return nil, fmt.Errorf("cannot modify system categories")
	}
	if cat.UserID == nil || *cat.UserID != userID {
		return nil, ErrForbidden
	}

	cat.Name = req.Name
	cat.Color = req.Color
	cat.Icon = req.Icon
	cat.ParentID = req.ParentID

	if err := s.categories.Update(ctx, cat); err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}
	return cat, nil
}

func (s *categoryService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	cat, err := s.categories.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return ErrNotFound
	}
	if cat.IsSystem {
		return fmt.Errorf("cannot delete system categories")
	}
	if cat.UserID == nil || *cat.UserID != userID {
		return ErrForbidden
	}
	return s.categories.Delete(ctx, id, userID)
}
