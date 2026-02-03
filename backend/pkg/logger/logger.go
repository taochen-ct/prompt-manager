package logger

import (
	"backend/pkg/common"
	"backend/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

func Init(conf *config.Config) (*lumberjack.Logger, *zap.Logger) {
	var level zapcore.Level  // zap 日志等级
	var options []zap.Option // zap 配置项

	logFileDir := conf.Log.RootDir
	if !filepath.IsAbs(logFileDir) {
		wd, err := os.Getwd()
		if err != nil {
			wd = ""
		}
		logFileDir = filepath.Join(wd, logFileDir)
	}
	if !common.IsExist(logFileDir) {
		_ = common.CreateDir(logFileDir)
	}

	switch conf.Log.Level {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(conf.Server.Env + "." + l.String())
	}

	loggerWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logFileDir, conf.Log.Filename),
		MaxSize:    conf.Log.MaxSize,
		MaxBackups: conf.Log.MaxBackups,
		MaxAge:     conf.Log.MaxAge,
		Compress:   conf.Log.Compress,
	}
	//fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(loggerWriter), level)
	//return loggerWriter, zap.New(zapcore.NewTee(fileCore), options...)

	consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), zapcore.ErrorLevel)
	fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(loggerWriter), level)
	core := zapcore.NewTee(consoleCore, fileCore)
	return loggerWriter, zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}
