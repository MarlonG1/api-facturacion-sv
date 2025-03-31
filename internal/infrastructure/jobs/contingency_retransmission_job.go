package jobs

import (
	"context"
	"errors"
	"github.com/MarlonG1/api-facturacion-sv/config/drivers"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/contingency/interfaces"
	"sync/atomic"
	"time"

	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

type RetransmissionJob struct {
	connection         *drivers.DbConnection
	ContingencyService interfaces.ContingencyManager
	IsRunning          atomic.Bool
	MaxExecutionTime   time.Duration
}

func NewRetransmissionJob(contingencyService interfaces.ContingencyManager, connection *drivers.DbConnection) *RetransmissionJob {
	return &RetransmissionJob{
		connection:         connection,
		ContingencyService: contingencyService,
		MaxExecutionTime:   10 * time.Minute,
	}
}

// Execute ejecuta el trabajo de retransmisión de documentos en contingencia.
func (j *RetransmissionJob) Execute() {
	// Evitar ejecuciones concurrentes
	if !j.IsRunning.CompareAndSwap(false, true) {
		logs.Warn("Job already running, skipping execution")
		return
	}
	defer j.IsRunning.Store(false)

	ctx, cancel := context.WithTimeout(context.Background(), j.MaxExecutionTime)
	defer cancel()

	logs.Info("Starting retransmission job", map[string]interface{}{
		"MaxExecutionTime": j.MaxExecutionTime,
		"timestamp":        utils.TimeNow().Format(time.RFC3339),
	})

	sqlDb, err := j.connection.Db.DB()
	if err != nil {
		logs.Error("Error connecting to database", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	sqlDb.Ping()

	if err := j.ContingencyService.RetransmitPendingDocuments(ctx); err != nil {
		j.handleExecutionError(err)
		return
	}

	logs.Info("Retransmission job completed successfully", map[string]interface{}{
		"timestamp": utils.TimeNow().Format(time.RFC3339),
	})
}

func (j *RetransmissionJob) handleExecutionError(err error) {
	if errors.Is(err, context.DeadlineExceeded) {
		logs.Error("Job execution timed out", map[string]interface{}{
			"MaxExecutionTime": j.MaxExecutionTime,
			"error":            err.Error(),
		})
		return
	}

	logs.Error("Job execution failed", map[string]interface{}{
		"error": err.Error(),
	})
}
