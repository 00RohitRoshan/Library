package main

import (
	"os"

	"github.com/00RohitRoshan/Library/Mylibrary"
	"github.com/00RohitRoshan/Library/graylog"
)

func main() {
	Mylibrary.SetName("Rohit")
	file, _ := os.Create("./graylog.log")
	defer file.Close()
	Logger := graylog.InitGraylog(graylog.Config{
		Adr:      "172.16.5.51:12201",
		Protocol: "udp",
		LogLevel: graylog.LogLevels[1],
		// File:     file,
	})

	Logger.SetStatic(graylog.Log{AppName: "TestLogApplication"})
	Logger.MustHaveFuncAdd(graylog.Fields[1], func() string {
		return " my Custom Ip address"
	})

	for range 100 {
		Logger.Info(graylog.Log{})
	}

}
