package logger

import (
	"log"
	"os"

	"github.com/aclgo/grpc-mail/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	StartLogger() error
	Degug(args ...any)
	Debugf(template string, args ...any)
	Info(args ...any)
	Infof(template string, args ...any)
	Warn(args ...any)
	Warnf(template string, args ...any)
	Error(args ...any)
	Errorf(template string, args ...any)
	Fatal(args ...any)
	Fatalf(template string, args ...any)
}

type apiLogger struct {
	config        *config.Config
	sugaredLogger *zap.SugaredLogger
}

func NewapiLogger(config *config.Config) *apiLogger {
	apiLogger := &apiLogger{
		config: config,
	}

	if err := apiLogger.StartLogger(); err != nil {
		log.Printf("NewapiLogger.StartLogger: %v", err)
	}

	return apiLogger
}

var mapLogLevel = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, ok := mapLogLevel[cfg.Logger.Level]
	if !ok && cfg.Logger.Level == "" {
		return zapcore.DebugLevel
	}

	return level
}

func (l *apiLogger) StartLogger() error {
	loggerLevel := getLoggerLevel(l.config)

	loggerWriter := zapcore.AddSync(os.Stderr)

	encoderConfig := zapcore.EncoderConfig{}

	if l.config.Logger.ServerMode == "dev" {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	encoderConfig.LevelKey = "LEVEL"
	encoderConfig.CallerKey = "CALLER"
	encoderConfig.TimeKey = "TIME"
	encoderConfig.NameKey = "NAME"
	encoderConfig.MessageKey = "MESSAGE"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder

	if l.config.Logger.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, loggerWriter, zap.NewAtomicLevelAt(loggerLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugaredLogger = logger.Sugar()
	if err := l.sugaredLogger.Sync(); err != nil {
		log.Printf("StartLogger.Sync: %v", err)
	}

	return nil
}

func (l *apiLogger) Degug(args ...any) {
	l.sugaredLogger.Debug(args...)
}

func (l *apiLogger) Debugf(template string, args ...any) {
	l.sugaredLogger.Debugf(template, args...)
}

func (l *apiLogger) Info(args ...any) {
	l.sugaredLogger.Info(args...)
}

func (l *apiLogger) Infof(template string, args ...any) {
	l.sugaredLogger.Infof(template, args...)
}

func (l *apiLogger) Warn(args ...any) {
	l.sugaredLogger.Warn(args...)
}

func (l *apiLogger) Warnf(template string, args ...any) {
	l.sugaredLogger.Warnf(template, args...)
}

func (l *apiLogger) Error(args ...any) {
	l.sugaredLogger.Error(args...)
}

func (l *apiLogger) Errorf(template string, args ...any) {
	l.sugaredLogger.Errorf(template, args...)
}

func (l *apiLogger) Fatal(args ...any) {
	l.sugaredLogger.Fatal(args...)
}

func (l *apiLogger) Fatalf(template string, args ...any) {
	l.sugaredLogger.Fatalf(template, args...)
}
