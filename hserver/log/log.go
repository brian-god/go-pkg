package log

import (
	"github.com/brian-god/go-pkg/configs"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path"
	"time"
)

func Config(cof *configs.Log) {
	if configs.Mode == configs.Development {
		//Set log output location
		hlog.SetOutput(os.Stdout)
		hlog.SetLevel(hlog.LevelDebug)
	} else {
		// Customizable output directory。
		logFilePath := cof.OutPath
		if logFilePath == "" {
			logFilePath = "./logs/"
		}
		if err := os.MkdirAll(logFilePath, 0o777); err != nil {
			log.Println(err.Error())
			return
		}

		// Set file name to date
		logFileName := cof.FilePrefix + time.Now().Format("2006-01-02") + ".log"
		fileName := path.Join(logFilePath, logFileName)
		if _, err := os.Stat(fileName); err != nil {
			if _, err := os.Create(fileName); err != nil {
				log.Println(err.Error())
				return
			}
		}

		logger := hertzzap.NewLogger()
		maxSize := 20
		if cof.MaxSize != 0 {
			maxSize = int(cof.MaxSize)
		}
		maxBackups := 5
		if cof.MaxBackups != 0 {
			maxSize = int(cof.MaxBackups)
		}
		maxAge := 10
		if cof.MaxAge != 0 {
			maxSize = int(cof.MaxAge)
		}
		// Provide compression and deletion
		lumberjackLogger := &lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    maxSize,      // 一个文件最大可达20M。
			MaxBackups: maxBackups,   // 最多同时保存 5 个文件。
			MaxAge:     maxAge,       // 一个文件最多可以保存 10 天。
			Compress:   cof.Compress, // 用 gzip 压缩。
		}

		logger.SetOutput(lumberjackLogger)
		logger.SetLevel(hlog.Level(cof.Level))
		hlog.SetLogger(logger)
	}
}
