package worker

import (
	"context"
	"sync"
	"testing"
	"time"

	"event-coming/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockSchedulerService is a mock implementation of service.SchedulerService
type MockSchedulerService struct {
	mock.Mock
}

func (m *MockSchedulerService) Create(ctx context.Context, input *domain.CreateSchedulerInput, orgID uuid.UUID) (*domain.Scheduler, error) {
	args := m.Called(ctx, input, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Scheduler), args.Error(1)
}

func (m *MockSchedulerService) GetByID(ctx context.Context, id uuid.UUID, orgID uuid.UUID) (*domain.Scheduler, error) {
	args := m.Called(ctx, id, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Scheduler), args.Error(1)
}

func (m *MockSchedulerService) Cancel(ctx context.Context, id uuid.UUID, orgID uuid.UUID) error {
	args := m.Called(ctx, id, orgID)
	return args.Error(0)
}

func (m *MockSchedulerService) ProcessPendingTasks(ctx context.Context, batchSize int) (int, error) {
	args := m.Called(ctx, batchSize)
	return args.Int(0), args.Error(1)
}

func TestNewSchedulerWorker(t *testing.T) {
	logger := zap.NewNop()
	mockService := new(MockSchedulerService)

	tests := []struct {
		name              string
		interval          time.Duration
		batchSize         int
		expectedInterval  time.Duration
		expectedBatchSize int
	}{
		{
			name:              "valid parameters",
			interval:          1 * time.Minute,
			batchSize:         50,
			expectedInterval:  1 * time.Minute,
			expectedBatchSize: 50,
		},
		{
			name:              "zero batch size defaults to 100",
			interval:          1 * time.Minute,
			batchSize:         0,
			expectedInterval:  1 * time.Minute,
			expectedBatchSize: 100,
		},
		{
			name:              "negative batch size defaults to 100",
			interval:          1 * time.Minute,
			batchSize:         -10,
			expectedInterval:  1 * time.Minute,
			expectedBatchSize: 100,
		},
		{
			name:              "zero interval defaults to 30 seconds",
			interval:          0,
			batchSize:         50,
			expectedInterval:  30 * time.Second,
			expectedBatchSize: 50,
		},
		{
			name:              "negative interval defaults to 30 seconds",
			interval:          -1 * time.Second,
			batchSize:         50,
			expectedInterval:  30 * time.Second,
			expectedBatchSize: 50,
		},
		{
			name:              "both defaults applied",
			interval:          0,
			batchSize:         0,
			expectedInterval:  30 * time.Second,
			expectedBatchSize: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := NewSchedulerWorker(mockService, logger, tt.interval, tt.batchSize)

			assert.NotNil(t, worker)
			assert.Equal(t, mockService, worker.schedulerService)
			assert.Equal(t, logger, worker.logger)
			assert.Equal(t, tt.expectedInterval, worker.interval)
			assert.Equal(t, tt.expectedBatchSize, worker.batchSize)
			assert.NotNil(t, worker.stopCh)
		})
	}
}

func TestSchedulerWorker_StartAndStop(t *testing.T) {
	logger := zap.NewNop()
	mockService := new(MockSchedulerService)

	// Expect at least one call to ProcessPendingTasks (immediate call on start)
	mockService.On("ProcessPendingTasks", mock.Anything, 100).Return(0, nil)

	worker := NewSchedulerWorker(mockService, logger, 100*time.Millisecond, 100)

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		worker.Start(ctx)
	}()

	// Let it run for a short time
	time.Sleep(150 * time.Millisecond)

	// Stop the worker via context cancellation
	cancel()

	// Wait for the worker to stop
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success - worker stopped
	case <-time.After(1 * time.Second):
		t.Fatal("Worker did not stop in time")
	}

	mockService.AssertExpectations(t)
}

func TestSchedulerWorker_StopSignal(t *testing.T) {
	logger := zap.NewNop()
	mockService := new(MockSchedulerService)

	// Expect at least one call to ProcessPendingTasks (immediate call on start)
	mockService.On("ProcessPendingTasks", mock.Anything, 50).Return(0, nil)

	worker := NewSchedulerWorker(mockService, logger, 100*time.Millisecond, 50)

	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		worker.Start(ctx)
	}()

	// Let it run for a short time
	time.Sleep(50 * time.Millisecond)

	// Stop the worker via stop signal
	worker.Stop()

	// Wait for the worker to stop
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success - worker stopped
	case <-time.After(1 * time.Second):
		t.Fatal("Worker did not stop in time")
	}

	mockService.AssertExpectations(t)
}

func TestSchedulerWorker_ProcessesTasksPeriodically(t *testing.T) {
	logger := zap.NewNop()
	mockService := new(MockSchedulerService)

	callCount := 0
	var mu sync.Mutex

	// Track how many times ProcessPendingTasks is called
	mockService.On("ProcessPendingTasks", mock.Anything, 100).Return(0, nil).Run(func(args mock.Arguments) {
		mu.Lock()
		callCount++
		mu.Unlock()
	})

	worker := NewSchedulerWorker(mockService, logger, 50*time.Millisecond, 100)

	ctx, cancel := context.WithCancel(context.Background())

	go worker.Start(ctx)

	// Let it run for enough time to process multiple times
	time.Sleep(180 * time.Millisecond)

	cancel()
	worker.wg.Wait()

	mu.Lock()
	count := callCount
	mu.Unlock()

	// Should have been called at least 3 times (1 immediate + ~2 from ticker)
	assert.GreaterOrEqual(t, count, 3, "Should process tasks periodically")
}

func TestSchedulerWorker_ReturnsProcessedCount(t *testing.T) {
	logger := zap.NewNop()
	mockService := new(MockSchedulerService)

	// Simulate processing some tasks
	mockService.On("ProcessPendingTasks", mock.Anything, 100).Return(5, nil).Once()
	mockService.On("ProcessPendingTasks", mock.Anything, 100).Return(0, nil)

	worker := NewSchedulerWorker(mockService, logger, 50*time.Millisecond, 100)

	ctx, cancel := context.WithCancel(context.Background())

	go worker.Start(ctx)

	// Let it run
	time.Sleep(30 * time.Millisecond)

	cancel()
	worker.wg.Wait()

	mockService.AssertExpectations(t)
}

func TestSchedulerWorker_HandlesErrors(t *testing.T) {
	logger := zap.NewNop()
	mockService := new(MockSchedulerService)

	// Simulate an error
	mockService.On("ProcessPendingTasks", mock.Anything, 100).Return(0, assert.AnError).Once()
	mockService.On("ProcessPendingTasks", mock.Anything, 100).Return(0, nil)

	worker := NewSchedulerWorker(mockService, logger, 50*time.Millisecond, 100)

	ctx, cancel := context.WithCancel(context.Background())

	go worker.Start(ctx)

	// Let it run
	time.Sleep(100 * time.Millisecond)

	cancel()
	worker.wg.Wait()

	// Should continue running even after errors
	mockService.AssertExpectations(t)
}

func TestSchedulerWorker_ProcessScheduledTasks(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name      string
		processed int
		err       error
	}{
		{
			name:      "no tasks processed",
			processed: 0,
			err:       nil,
		},
		{
			name:      "some tasks processed",
			processed: 10,
			err:       nil,
		},
		{
			name:      "error during processing",
			processed: 0,
			err:       assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSchedulerService)
			mockService.On("ProcessPendingTasks", mock.Anything, 100).Return(tt.processed, tt.err)

			worker := NewSchedulerWorker(mockService, logger, 30*time.Second, 100)

			// Call processScheduledTasks directly
			worker.processScheduledTasks(context.Background())

			mockService.AssertExpectations(t)
		})
	}
}
