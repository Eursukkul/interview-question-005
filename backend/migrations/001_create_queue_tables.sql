CREATE TABLE IF NOT EXISTS queue_state (
  id INTEGER PRIMARY KEY,
  current_letter TEXT NULL CHECK (current_letter IS NULL OR current_letter ~ '^[A-Z]$'),
  current_number INTEGER NULL CHECK (current_number IS NULL OR current_number BETWEEN 0 AND 9),
  current_queue TEXT NOT NULL DEFAULT '00',
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT queue_state_single_row CHECK (id = 1),
  CONSTRAINT queue_state_consistency CHECK (
    (current_queue = '00' AND current_letter IS NULL AND current_number IS NULL)
    OR
    (current_queue ~ '^[A-Z][0-9]$' AND current_letter IS NOT NULL AND current_number IS NOT NULL)
  )
);

INSERT INTO queue_state (id, current_letter, current_number, current_queue, updated_at)
VALUES (1, NULL, NULL, '00', NOW())
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS queue_history (
  id BIGSERIAL PRIMARY KEY,
  queue_number TEXT NOT NULL CHECK (queue_number ~ '^[A-Z][0-9]$'),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_queue_history_created_at ON queue_history (created_at DESC);
