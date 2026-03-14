// Package logger provides a flexible logging solution using zerolog and lumberjack for log rotation.
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {
	EnableConsole bool
	FilePath      string
	MaxSize       int  // Megabytes before rotating
	MaxBackups    int  // Number of old log files to keep
	MaxAge        int  // Days to keep old log files
	Compress      bool // Whether to zip old files
}

func New(cfg LogConfig) zerolog.Logger {
	var writers []io.Writer

	// 1. Setup Console (Pretty Print)
	if cfg.EnableConsole {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	// 2. Setup Lumberjack for File Logging
	if cfg.FilePath != "" {
		// Default values if not specified
		if cfg.MaxSize == 0 {
			cfg.MaxSize = 10
		}
		if cfg.MaxBackups == 0 {
			cfg.MaxBackups = 5
		}

		rollingLogger := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		writers = append(writers, rollingLogger)
	}

	multi := zerolog.MultiLevelWriter(writers...)
	return zerolog.New(multi).With().Timestamp().Logger()
}
