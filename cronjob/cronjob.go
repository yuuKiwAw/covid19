package cronjob

import (
	"io"
	"log"
	"os"

	"covid19/covid19"

	"github.com/robfig/cron"
)

// 马上执行一次
func run_immediately() {
	covid19.GetCovid19info()
}

// 主定时器
func Main_cron() {
	run_immediately()

	c := cron.New()
	// spec := "*/5 * * * * ?"
	spec2 := "0 0/40 * * * ? "
	c.AddFunc(spec2, func() {
		// 定时业务存放
		covid19.GetCovid19info()
	})
	log.Println("Start CronJob get covid19 info")
	c.Start()

	select {}
}

func init() {
	log_infoPath := "./log/logs.log"

	// 保存日志信息
	logFile, err := os.OpenFile(log_infoPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}

	// 同时显示在终端并写入到log文件
	writers := []io.Writer{
		logFile,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)

	log.SetOutput(fileAndStdoutWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
