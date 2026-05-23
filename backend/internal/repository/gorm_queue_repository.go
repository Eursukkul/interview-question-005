package repository

import (
	"context"
	"time"

	"example.com/interview-question-005/backend/internal/model"

	"gorm.io/gorm"
)

type GormQueueRepository struct {
	db *gorm.DB
}

func NewGormQueueRepository(db *gorm.DB) *GormQueueRepository {
	return &GormQueueRepository{db: db}
}

func (r *GormQueueRepository) EnsureState(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO queue_state (id, current_letter, current_number, current_queue, updated_at)
		VALUES (1, NULL, NULL, '00', NOW())
		ON CONFLICT (id) DO NOTHING
	`).Error
}

func (r *GormQueueRepository) GetCurrent(ctx context.Context) (model.QueueState, error) {
	var state model.QueueState
	err := r.db.WithContext(ctx).
		Raw(`
			SELECT s.id, s.current_letter, s.current_number, s.current_queue, s.updated_at,
			       (SELECT h.created_at
			          FROM queue_history h
			         WHERE h.queue_number = s.current_queue
			         ORDER BY h.id DESC
			         LIMIT 1) AS issued_at
			FROM queue_state s
			WHERE s.id = 1
		`).
		Scan(&state).Error
	return state, err
}

func (r *GormQueueRepository) Reset(ctx context.Context) (model.QueueState, error) {
	var state model.QueueState
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Raw(`
			UPDATE queue_state
			SET current_letter = NULL,
			    current_number = NULL,
			    current_queue = '00',
			    updated_at = NOW()
			WHERE id = 1
			RETURNING id, current_letter, current_number, current_queue, updated_at
		`).Scan(&state).Error
	})
	return state, err
}

func (r *GormQueueRepository) Next(ctx context.Context, calculate func(model.QueueState) (model.QueueState, error)) (model.QueueState, error) {
	var next model.QueueState
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current model.QueueState
		if err := tx.Raw(`
			SELECT id, current_letter, current_number, current_queue, updated_at
			FROM queue_state
			WHERE id = 1
			FOR UPDATE
		`).Scan(&current).Error; err != nil {
			return err
		}

		calculated, err := calculate(current)
		if err != nil {
			return err
		}

		if err := tx.Raw(`
			UPDATE queue_state
			SET current_letter = ?,
			    current_number = ?,
			    current_queue = ?,
			    updated_at = NOW()
			WHERE id = 1
			RETURNING id, current_letter, current_number, current_queue, updated_at
		`, calculated.CurrentLetter, calculated.CurrentNumber, calculated.CurrentQueue).Scan(&next).Error; err != nil {
			return err
		}

		var inserted struct {
			CreatedAt time.Time
		}
		if err := tx.Raw(`
			INSERT INTO queue_history (queue_number, created_at)
			VALUES (?, NOW())
			RETURNING created_at
		`, next.CurrentQueue).Scan(&inserted).Error; err != nil {
			return err
		}
		next.IssuedAt = &inserted.CreatedAt
		return nil
	})
	return next, err
}
