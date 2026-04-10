package utils

import (
	"io"
	"log/slog"
	"os"
)

func InitLogger() *os.File {
	logDir := "log/"
	logFile := logDir + "logs.log"

	err := os.MkdirAll(logDir,0755)
	if err != nil {
		slog.Error("Failed to create lo directory.","err",err)
		panic(err)
	}

	file,err := os.OpenFile(logFile,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	if err != nil {
		slog.Error("Failed to open file.","err",err,"filename",logFile)
		panic(err)
	}

	writer := io.MultiWriter(os.Stdout,file)	

	slogHanlder := slog.NewTextHandler(writer,&slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				return slog.String("time",t.Format("2006-01-02 15:04:05"))
			}
			return a
		},
	})

	logger := slog.New(slogHanlder)
	slog.SetDefault(logger)
	return file
}
