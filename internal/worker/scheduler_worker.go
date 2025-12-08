package worker

import (
	"context"
	"sync"
	"time"

	"event-coming/internal/service"

	"go.uber.org/zap"
)

// SchedulerWorker processa tasks agendadas periodicamente
type SchedulerWorker struct {
	schedulerService service.SchedulerService
	logger           *zap.Logger
	interval         time.Duration
	batchSize        int
	stopCh           chan struct{}
	wg               sync.WaitGroup
}

// NewSchedulerWorker cria um novo worker de scheduler
func NewSchedulerWorker(
	schedulerService service.SchedulerService,
	logger *zap.Logger,
	interval time.Duration,
	batchSize int,
) *SchedulerWorker {
	if batchSize <= 0 {
		batchSize = 100
	}
	if interval <= 0 {
		interval = 30 * time.Second
	}

	return &SchedulerWorker{
		schedulerService: schedulerService,
		logger:           logger,
		interval:         interval,
		batchSize:        batchSize,
		stopCh:           make(chan struct{}),
	}
}

// Start inicia o loop de processamento
func (w *SchedulerWorker) Start(ctx context.Context) {
	w.wg.Add(1)
	defer w.wg.Done()

	w.logger.Info("Scheduler worker started",
		zap.Duration("interval", w.interval),
		zap.Int("batch_size", w.batchSize),
	)

	// Processar imediatamente ao iniciar
	w.processScheduledTasks(ctx)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Scheduler worker stopping (context cancelled)")
			return
		case <-w.stopCh:
			w.logger.Info("Scheduler worker stopping (stop signal)")
			return
		case <-ticker.C:
			w.processScheduledTasks(ctx)
		}
	}
}

// Stop para o worker gracefully
func (w *SchedulerWorker) Stop() {
	close(w.stopCh)
	w.wg.Wait()
	w.logger.Info("Scheduler worker stopped")
}

// processScheduledTasks processa as tasks pendentes
func (w *SchedulerWorker) processScheduledTasks(ctx context.Context) {
	start := time.Now()

	processed, err := w.schedulerService.ProcessPendingTasks(ctx, w.batchSize)
	if err != nil {
		w.logger.Error("Failed to process scheduled tasks", zap.Error(err))
		return
	}

	if processed > 0 {
		w.logger.Info("Processed scheduled tasks",
			zap.Int("count", processed),
			zap.Duration("duration", time.Since(start)),
		)
	}
}
