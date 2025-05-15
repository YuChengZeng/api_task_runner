package worker

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"api_task_runner/configs"
	"api_task_runner/internal/api"
	"api_task_runner/internal/db"
	"api_task_runner/internal/models"
	"api_task_runner/pkg/logger"
)

func Run() {
    cfg := configs.LoadConfig()
    progressDB := db.NewProgressDB(cfg.MongoURI, cfg.MongoDBName, cfg.MongoColl)
    ctx := context.Background()

    log := logger.Logger.With("module", "worker")

    file, err := os.ReadFile("data/data.json")
    if err != nil {
        log.Errorf("Failed to read data.json: %v", err)
        return
    }

    var records []models.Record
    if err := json.Unmarshal(file, &records); err != nil {
        log.Errorf("Failed to parse data.json: %v", err)
        return
    }

    ticker := time.NewTicker(2 * time.Second / time.Duration(cfg.RateLimit))
    defer ticker.Stop()

    headers := map[string]string{"api-key": cfg.IntelligenceKeyCIB}

    processedCount := 0
    skippedCount := 0
    failedCount := 0

    for _, r := range records {
        processed, _ := progressDB.IsProcessed(ctx, r.Blockchain, r.Address)
        if processed {
            log.Infof("Already processed: %s %s", r.Blockchain, r.Address)
            skippedCount++
            continue
        }

        <-ticker.C

        params := map[string]string{
            "chain_name":       r.Blockchain,
            "address":          r.Address,
            "source_list_code": "010",
            "search_flag":      "false",
            "quick_mode":       "false",
        }

        response, err := api.MakeRequest(ctx, cfg.IntelligenceHost, headers, params)

        if err == nil {
            processedCount++
            completedCount := processedCount + skippedCount
            progressDB.MarkAsDone(ctx, r.Blockchain, r.Address, response)
            log.Infof("Processed successfully: %s %s | Completed: %d | Failed: %d",
                r.Blockchain, r.Address, completedCount, failedCount)
        } else {
            failedCount++
            completedCount := processedCount + skippedCount
            progressDB.MarkAsFailed(ctx, r.Blockchain, r.Address, err.Error())
            log.Warnf("Failed processing: %s %s, Error: %v | Completed: %d | Failed: %d",
                r.Blockchain, r.Address, err, completedCount, failedCount)
        }
    }

    log.Infof("Processing complete. Total processed: %d, Skipped (already done): %d, Failed: %d, Total records: %d",
        processedCount, skippedCount, failedCount, len(records))
}