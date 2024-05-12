package api

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/op/go-logging"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log = logging.MustGetLogger("go_logger")

var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05.999} PID:%{pid} ▶ %{level:.4s} %{id:03x} [%{shortfile}][%{shortfunc}]%{color:reset} %{message}`,
)

func InitLogger(logPath string, logDebug bool) {
	var backends []logging.Backend

	dirname := filepath.Dir(logPath)
	error := os.MkdirAll(dirname, 0644)
	if error != nil {
		fmt.Printf("Failed to create '%s' directory: %s\n", logPath, error)
		os.Exit(100)
	}

	logf := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, /* megabytes */
		MaxBackups: 10,
		MaxAge:     30, /* days */
		Compress:   true,
	}

	stderrBackend := logging.NewLogBackend(os.Stderr, "", 0)
	formattedStderrBackend := logging.NewBackendFormatter(stderrBackend, format)
	backends = append(backends, formattedStderrBackend)

	fileBackend := logging.NewLogBackend(logf, "", 0)
	formatterFileBackend := logging.NewBackendFormatter(fileBackend, format)
	backends = append(backends, formatterFileBackend)

	logging.SetBackend(backends...)
	logging.SetLevel(logging.INFO, "")

	if logDebug {
		logging.SetLevel(logging.INFO, "")

	}
}
