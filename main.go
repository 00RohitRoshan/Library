package main

import (
	"github.com/00RohitRoshan/Rohit/Mylibrary"
	"github.com/00RohitRoshan/Rohit/graylog"
)

func main() {
	Mylibrary.SetName("Rohit")
	Logger := graylog.InitGraylog(graylog.Config{
		Adr:      "172.16.5.51:12201",
		Protocol: "udp",
	})

	for range 100 {
		Logger.Log(graylog.Log{
			
		})
	}
}
