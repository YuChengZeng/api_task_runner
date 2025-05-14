package worker

import (
    "context"
    "encoding/json"
    "os"
    "time"
    "log"
    "api_task_runner/configs"
    "api_task_runner/internal/api"
    "api_task_runner/internal/db"
    "api_task_runner/internal/models"
)

func Run() {
    cfg := configs.LoadConfig()
    progressDB := db.NewProgressDB(cfg.MongoURI, cfg.MongoDBName, cfg.MongoColl)
    ctx := context.Background()

    // 讀取 JSON 檔
    file, _ := os.ReadFile("data.json")
    var records []models.Record
    json.Unmarshal(file, &records)

    ticker := time.NewTicker(time.Second / time.Duration(cfg.RateLimit))
    defer ticker.Stop()

    headers := map[string]string{"api-key": cfg.IntelligenceKeyCIB}

    for _, r := range records {
        keyProcessed, _ := progressDB.IsProcessed(ctx, r.Blockchain, r.Address)
        if keyProcessed {
            continue
        }

        <-ticker.C

        params := map[string]string{
            "chain_name":    r.Blockchain,
            "address":       r.Address,
            "source_list_code": "010",
            "search_flag":   "false",
            "quick_mode":    "false",
        }

        response, err := api.MakeRequest(ctx, cfg.IntelligenceHost, headers, params)
        if err == nil {
            progressDB.MarkAsDone(ctx, r.Blockchain, r.Address, response)
        } else {
            log.Printf("Failed: %s %s. Error: %v", r.Blockchain, r.Address, err)
            progressDB.MarkAsFailed(ctx, r.Blockchain, r.Address, err.Error())
        }
    }
}