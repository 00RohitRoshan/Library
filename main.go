package main

import (
	"github.com/00RohitRoshan/Library/Mylibrary"
	"github.com/00RohitRoshan/Library/graylog"
)

func main() {
	Mylibrary.SetName("Rohit")
	Logger := graylog.InitGraylog(graylog.Config{
		Adr:      "172.16.5.51:12201",
		Protocol: "udp",
	})

	for range 100 {
		Logger.Info(graylog.Log{})
	}
}
