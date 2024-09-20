package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	sysLog  *Logger
	service = ""
)

type Message struct {
	ChainID string `json:"chainID"`
	Level   string `json:"level"`
	Version string `json:"version"`
	Service string `json:"service"`
	Time    string `json:"time"`
	Msg     string `json:"msg"`
}

func SysLog() *Logger {
	return sysLog
}

func InitSysLog(serviceName, level string) error {
	service = serviceName
	sysLog = newLogger(getZapLevel(level), serviceName)

	sysLog.Info(context.Background(), "Logger init success")
	return nil
}

func newLogger(level zapcore.Level, serviceName string) *Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(level),
		Development:   true,
		Encoding:      "console",
		EncoderConfig: encoderConfig,
		InitialFields: map[string]interface{}{
			"version": "1",
			"service": serviceName,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	return &Logger{logger: logger}
}

type LoggerInterface interface {
	Debug(ctx context.Context, message string)
	Info(ctx context.Context, message string)
	Warn(ctx context.Context, message string)
	Error(ctx context.Context, message string)
	Panic(ctx context.Context, message string)
}

type Logger struct {
	logger *zap.Logger
}

func (lg *Logger) Debug(ctx context.Context, message string) {
	lg.logger.Debug(message, chainID(ctx))
}

func (lg *Logger) Info(ctx context.Context, message string) {
	lg.logger.Info(message, chainID(ctx))
}

func (lg *Logger) Warn(ctx context.Context, message string) {
	lg.logger.Warn(message, chainID(ctx))
}

func (lg *Logger) Error(ctx context.Context, message string) {
	lg.logger.Error(message, chainID(ctx))
}

func (lg *Logger) Panic(ctx context.Context, message string) {
	lg.logger.Panic(message, chainID(ctx))
}

func chainID(ctx context.Context) zap.Field {
	chainID := ""
	if s := ctx.Value("ChainID"); s != nil {
		chainID = s.(string)
	}

	return zap.String("chain", chainID)
}

func getZapLevel(l string) zapcore.Level {
	switch l {
	case zapcore.DebugLevel.String(): // "debug"
		return zapcore.DebugLevel
	case zapcore.InfoLevel.String(): // "info"
		return zapcore.InfoLevel
	case zapcore.WarnLevel.String(): // "warn"
		return zapcore.WarnLevel
	case zapcore.ErrorLevel.String(): // "error"
		return zapcore.ErrorLevel
	case zapcore.DPanicLevel.String(): // "dpanic"
		return zapcore.DPanicLevel
	case zapcore.PanicLevel.String(): // "panic"
		return zapcore.PanicLevel
	case zapcore.FatalLevel.String(): // "fatal"
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
