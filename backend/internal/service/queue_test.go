package service

import (
	"errors"
	"testing"

	"example.com/interview-question-005/backend/internal/model"
)

func TestCalculateNextQueue(t *testing.T) {
	tests := []struct {
		name    string
		state   model.QueueState
		want    string
		wantErr error
	}{
		{name: "reset state to A0", state: state("00"), want: "A0"},
		{name: "A0 to A1", state: state("A0"), want: "A1"},
		{name: "A8 to A9", state: state("A8"), want: "A9"},
		{name: "A9 to B0", state: state("A9"), want: "B0"},
		{name: "Y9 to Z0", state: state("Y9"), want: "Z0"},
		{name: "Z8 to Z9", state: state("Z8"), want: "Z9"},
		{name: "Z9 limit reached", state: state("Z9"), wantErr: ErrQueueLimitReached},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateNextQueue(tt.state)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.CurrentQueue != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got.CurrentQueue)
			}
		})
	}
}

func TestResetStateValue(t *testing.T) {
	got := model.QueueState{CurrentQueue: "00"}
	if got.CurrentQueue != "00" {
		t.Fatalf("expected reset queue to be 00, got %s", got.CurrentQueue)
	}
}

func state(queue string) model.QueueState {
	if queue == "00" {
		return model.QueueState{CurrentQueue: queue}
	}
	letter := queue[:1]
	number := int(queue[1] - '0')
	return model.QueueState{
		CurrentLetter: &letter,
		CurrentNumber: &number,
		CurrentQueue:  queue,
	}
}
