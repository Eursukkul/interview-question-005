package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"example.com/interview-question-005/backend/internal/handler"
	"example.com/interview-question-005/backend/internal/model"
	"example.com/interview-question-005/backend/internal/repository"
	"example.com/interview-question-005/backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestConcurrentNextRequestsDoNotDuplicateQueues(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	db, cleanup := setupTestDatabase(t, ctx, databaseURL)
	defer cleanup()

	repo := repository.NewGormQueueRepository(db)
	if err := repo.EnsureState(ctx); err != nil {
		t.Fatalf("ensure state: %v", err)
	}

	queueHandler := handler.NewQueueHandler(service.NewQueueService(repo))
	app := fiber.New()
	app.Post("/api/queue/next", queueHandler.Next)

	const requestCount = 40
	var wg sync.WaitGroup
	results := make(chan string, requestCount)
	errors := make(chan error, requestCount)

	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, "/api/queue/next", nil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			if err != nil {
				errors <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errors <- fmt.Errorf("unexpected status %d", resp.StatusCode)
				return
			}

			var payload model.QueueResponse
			if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
				errors <- err
				return
			}
			results <- payload.QueueNumber
		}()
	}

	wg.Wait()
	close(results)
	close(errors)

	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}

	seen := map[string]struct{}{}
	for queueNumber := range results {
		if _, exists := seen[queueNumber]; exists {
			t.Fatalf("duplicate queue number generated: %s", queueNumber)
		}
		seen[queueNumber] = struct{}{}
	}
	if len(seen) != requestCount {
		t.Fatalf("expected %d queue numbers, got %d", requestCount, len(seen))
	}
}

func setupTestDatabase(t *testing.T, ctx context.Context, databaseURL string) (*gorm.DB, func()) {
	t.Helper()

	adminDB, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		t.Fatalf("connect admin db: %v", err)
	}
	adminSQL, err := adminDB.DB()
	if err != nil {
		t.Fatalf("get admin sql db: %v", err)
	}
	defer adminSQL.Close()

	schema := fmt.Sprintf("queue_test_%d", time.Now().UnixNano())
	if err := adminDB.WithContext(ctx).Exec(fmt.Sprintf("CREATE SCHEMA %s", schema)).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}

	schemaURL, err := databaseURLWithSearchPath(databaseURL, schema)
	if err != nil {
		t.Fatalf("build schema database url: %v", err)
	}
	db, err := gorm.Open(postgres.Open(schemaURL), &gorm.Config{})
	if err != nil {
		t.Fatalf("connect schema db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get schema sql db: %v", err)
	}

	migration, err := os.ReadFile("../migrations/001_create_queue_tables.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	if err := db.WithContext(ctx).Exec(string(migration)).Error; err != nil {
		t.Fatalf("run migration: %v", err)
	}

	cleanup := func() {
		_ = sqlDB.Close()
		adminDB, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
		if err != nil {
			return
		}
		adminSQL, err := adminDB.DB()
		if err == nil {
			defer adminSQL.Close()
		}
		_ = adminDB.WithContext(ctx).Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema)).Error
	}
	return db, cleanup
}

func databaseURLWithSearchPath(databaseURL string, schema string) (string, error) {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return "", err
	}

	query := parsed.Query()
	query.Set("search_path", schema)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}
