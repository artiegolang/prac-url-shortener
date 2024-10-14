package logger

import "go.uber.org/zap"

// Logger is an interface for logging
type Logger interface {
	Info(args ...interface{})
}

type logger struct {
	sugarLogger *zap.SugaredLogger
}

// NewLogger creates a new logger
func NewLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	// Не вызывайте defer logger.Sync() здесь, так как это завершит работу логгера сразу после выхода из функции
	return logger.Sugar()
}

func (l *logger) Info(format string, args ...interface{}) {
	l.sugarLogger.Infof(format, args...)
}
