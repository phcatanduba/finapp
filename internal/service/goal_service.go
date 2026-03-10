package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"finapp/internal/model"
	"finapp/internal/repository"
)

type goalService struct {
	goals repository.GoalRepository
}

func NewGoalService(goals repository.GoalRepository) GoalService {
	return &goalService{goals: goals}
}

func (s *goalService) Create(ctx context.Context, userID uuid.UUID, req model.GoalRequest) (*model.Goal, error) {
	if req.Name == "" || req.TargetAmount <= 0 {
		return nil, fmt.Errorf("name and target_amount are required")
	}
	g := &model.Goal{
		ID:            uuid.New(),
		UserID:        userID,
		Name:          req.Name,
		Description:   req.Description,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: req.CurrentAmount,
		Deadline:      req.Deadline,
		Type:          req.Type,
	}
	if err := s.goals.Create(ctx, g); err != nil {
		return nil, fmt.Errorf("create goal: %w", err)
	}
	return g, nil
}

func (s *goalService) List(ctx context.Context, userID uuid.UUID) ([]model.Goal, error) {
	return s.goals.FindByUserID(ctx, userID)
}

func (s *goalService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Goal, error) {
	g, err := s.goals.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, ErrNotFound
	}
	return g, nil
}

func (s *goalService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req model.GoalRequest) (*model.Goal, error) {
	g, err := s.goals.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, ErrNotFound
	}

	g.Name = req.Name
	g.Description = req.Description
	g.TargetAmount = req.TargetAmount
	g.CurrentAmount = req.CurrentAmount
	g.Deadline = req.Deadline
	g.Type = req.Type
	g.IsCompleted = g.CurrentAmount >= g.TargetAmount

	if err := s.goals.Update(ctx, g); err != nil {
		return nil, fmt.Errorf("update goal: %w", err)
	}
	return g, nil
}

func (s *goalService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	g, err := s.goals.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	if g == nil {
		return ErrNotFound
	}
	return s.goals.Delete(ctx, id, userID)
}

func (s *goalService) GetProgress(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.GoalProgress, error) {
	g, err := s.goals.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, ErrNotFound
	}

	remaining := g.TargetAmount - g.CurrentAmount
	pct := 0.0
	if g.TargetAmount > 0 {
		pct = (g.CurrentAmount / g.TargetAmount) * 100
	}

	progress := &model.GoalProgress{
		Goal:       *g,
		Remaining:  remaining,
		Percentage: pct,
		OnTrack:    false,
	}

	if g.Deadline != nil {
		now := time.Now()
		days := int(g.Deadline.Sub(now).Hours() / 24)
		progress.DaysToDeadline = &days

		// On track: progress percentage >= time elapsed percentage
		totalDays := g.Deadline.Sub(g.CreatedAt).Hours() / 24
		if totalDays > 0 {
			elapsedDays := now.Sub(g.CreatedAt).Hours() / 24
			timeElapsedPct := elapsedDays / totalDays * 100
			progress.OnTrack = pct >= timeElapsedPct
		}
	}

	return progress, nil
}
