package runLog

import (
	"gin-web/init/config"
	"os"
	"path/filepath"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var RunningLog *zap.Logger

func InitRunLog(conf config.Config) error {
	// 确保日志目录存在
	if err := os.MkdirAll(conf.APP.RunLog, os.ModePerm); err != nil {
		return err
	}

	// 日志编码配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 🔥 使用 file-rotatelogs 按天轮转
	logFileName := filepath.Join(conf.APP.RunLog, "app-%Y-%m-%d.log")
	//每次utc时间轮询
	rotator, err := rotatelogs.New(
		logFileName, // 轮转后的文件名，%Y-%m-%d 会被替换为日期
		rotatelogs.WithLinkName(filepath.Join(conf.APP.RunLog, "app.log")), // 软链接指向最新日志
		rotatelogs.WithMaxAge(30*24*time.Hour),                             // 30天过期
		rotatelogs.WithRotationTime(24*time.Hour),                          // 每24小时轮转一次
	)
	if err != nil {
		return err
	}

	// 日志级别
	var level zapcore.Level
	if conf.APP.Mode == "debug" {
		level = zapcore.DebugLevel // 开启所有日志，包括 Debug
	} else {
		level = zapcore.InfoLevel //只屏蔽 Debug,输出 Info 及以上
	}

	// 🔥 将 rotator (它实现了 WriteSyncer) 传给 zap
	core := zapcore.NewCore(encoder, zapcore.AddSync(rotator), level)
	RunningLog = zap.New(core, zap.AddCaller())
	return nil
}
