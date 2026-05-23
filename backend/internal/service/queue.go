package service

import (
	"context"
	"errors"
	"fmt"

	"example.com/interview-question-005/backend/internal/model"
)

var (
	ErrQueueLimitReached = errors.New("Queue limit reached")
	ErrInvalidQueueState = errors.New("invalid queue state")
)

type QueueRepository interface {
	GetCurrent(ctx context.Context) (model.QueueState, error)
	Reset(ctx context.Context) (model.QueueState, error)
	Next(ctx context.Context, calculate func(model.QueueState) (model.QueueState, error)) (model.QueueState, error)
}

type QueueService struct {
	repository QueueRepository
}

func NewQueueService(repository QueueRepository) *QueueService {
	return &QueueService{repository: repository}
}

func (s *QueueService) Current(ctx context.Context) (model.QueueState, error) {
	return s.repository.GetCurrent(ctx)
}

func (s *QueueService) Reset(ctx context.Context) (model.QueueState, error) {
	return s.repository.Reset(ctx)
}

func (s *QueueService) Next(ctx context.Context) (model.QueueState, error) {
	return s.repository.Next(ctx, CalculateNextQueue)
}

func CalculateNextQueue(state model.QueueState) (model.QueueState, error) {
	if state.CurrentQueue == "00" {
		letter := "A"
		number := 0
		state.CurrentLetter = &letter
		state.CurrentNumber = &number
		state.CurrentQueue = "A0"
		return state, nil
	}

	if err := validateState(state); err != nil {
		return model.QueueState{}, err
	}

	letter := *state.CurrentLetter
	number := *state.CurrentNumber

	if letter == "Z" && number == 9 {
		return model.QueueState{}, ErrQueueLimitReached
	}

	if number < 9 {
		number++
	} else {
		letter = string(rune(letter[0]) + 1)
		number = 0
	}

	state.CurrentLetter = &letter
	state.CurrentNumber = &number
	state.CurrentQueue = fmt.Sprintf("%s%d", letter, number)
	return state, nil
}

func validateState(state model.QueueState) error {
	if state.CurrentLetter == nil || state.CurrentNumber == nil {
		return fmt.Errorf("%w: current letter and number are required when current_queue is not 00", ErrInvalidQueueState)
	}

	letter := *state.CurrentLetter
	number := *state.CurrentNumber
	if len(letter) != 1 || letter[0] < 'A' || letter[0] > 'Z' {
		return fmt.Errorf("%w: current_letter must be A-Z", ErrInvalidQueueState)
	}
	if number < 0 || number > 9 {
		return fmt.Errorf("%w: current_number must be 0-9", ErrInvalidQueueState)
	}
	if state.CurrentQueue != fmt.Sprintf("%s%d", letter, number) {
		return fmt.Errorf("%w: current_queue does not match current_letter/current_number", ErrInvalidQueueState)
	}
	return nil
}
