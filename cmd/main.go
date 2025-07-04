package main

import (
    "api_task_runner/configs"
    "api_task_runner/internal/worker"
    "api_task_runner/pkg/logger"
)

func main() {
    cfg := configs.LoadConfig()

    logger.Initialize(cfg.LoggerConfig)
    defer logger.Close()

    logger.Logger.Info("Application starting...")

    workerLog := logger.Logger.With("module", "worker")
    workerLog.Info("Worker module starting...")

    worker.Run()
}