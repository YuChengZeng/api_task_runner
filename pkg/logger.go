package logger

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

var (
    Logger    *zap.SugaredLogger
    once      sync.Once
    logCloser func() error
)

type Config struct {
    LogDir      string
    FileName    string
    MaxKeepDays int
    Level       string
    EnableFile  bool
}

func Initialize(cfg Config) {
    once.Do(func() {
        if cfg.LogDir == "" {
            cfg.LogDir = "./logs"
        }
        if cfg.FileName == "" {
            cfg.FileName = "app.log"
        }
        if cfg.MaxKeepDays <= 0 {
            cfg.MaxKeepDays = 30
        }
        if cfg.Level == "" {
            cfg.Level = "info"
        }

        _ = os.MkdirAll(cfg.LogDir, os.ModePerm)

        level := parseLogLevel(cfg.Level)

        encoderConfig := zapcore.EncoderConfig{
            TimeKey:        "time",
            LevelKey:       "level",
            NameKey:        "logger",
            CallerKey:      "caller",
            MessageKey:     "msg",
            StacktraceKey:  "stacktrace",
            LineEnding:     zapcore.DefaultLineEnding,
            EncodeLevel:    zapcore.LowercaseLevelEncoder,
            EncodeTime:     localTimeEncoder,
            EncodeDuration: zapcore.StringDurationEncoder,
            EncodeCaller:   zapcore.ShortCallerEncoder,
        }

        encoder := zapcore.NewJSONEncoder(encoderConfig)

        var writeSyncers []zapcore.WriteSyncer
        writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))

        if cfg.EnableFile {
            logPath := filepath.Join(cfg.LogDir, cfg.FileName)
            logPattern := strings.TrimSuffix(logPath, filepath.Ext(logPath)) + "-%Y-%m-%d" + filepath.Ext(logPath)
            writer, err := rotatelogs.New(
                logPattern,
                rotatelogs.WithLinkName(logPath),
                rotatelogs.WithRotationTime(24*time.Hour),
                rotatelogs.WithMaxAge(time.Duration(cfg.MaxKeepDays)*24*time.Hour),
            )
            if err != nil {
                fmt.Println("Failed to initialize file logger:", err)
                os.Exit(1)
            }
            writeSyncers = append(writeSyncers, zapcore.AddSync(writer))
            logCloser = writer.Close
        }

        core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncers...), level)
        Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()

        Logger.Infof("Logger initialized. Level: %s, File Output: %v, Log Path: %s", cfg.Level, cfg.EnableFile, cfg.LogDir)
    })
}

func Close() {
    if logCloser != nil {
        _ = logCloser()
    }
    if Logger != nil {
        _ = Logger.Sync()
    }
}

func parseLogLevel(level string) zapcore.Level {
    switch strings.ToLower(level) {
    case "debug":
        return zapcore.DebugLevel
    case "warn":
        return zapcore.WarnLevel
    case "error":
        return zapcore.ErrorLevel
    default:
        return zapcore.InfoLevel
    }
}

func localTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
    zone, offset := t.Zone()
    enc.AppendString(t.Local().Format("2006-01-02T15:04:05") + fmt.Sprintf(" %s%+02d:%02d", zone, offset/3600, (offset%3600)/60))
}