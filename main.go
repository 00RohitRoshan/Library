package main

import "fmt"

func main() {
	// Mylibrary.SetName("Rohit")
	// file, _ := os.Create("./graylog.log")
	// defer file.Close()
	// Logger := graylog.InitGraylog(graylog.Config{
	// 	Adr:      "172.16.5.51:12201",
	// 	Protocol: "udp",
	// 	LogLevel: graylog.LogLevels[1],
	// 	File:     file,
	// })

	// Logger.SetStatic(graylog.Log{AppName: "TestLogApplication"})
	// Logger.MustHaveFuncAdd(graylog.Fields[1], func() string {
	// 	return " my Custom Ip address"
	// })

	// for range 100 {
	// 	Logger.Info(graylog.Log{Message: "wertyuiop"})
	// }

	var e = ConnectionPayload{
		EventName: EventNames(DIRECT_MESSAGE),
	}
	fmt.Println("e", e)
	e.SetEventName("rtyui")
	fmt.Println("e", e)

}

type EventNames string

const (
	NEW_USER       EventNames = "NEW_USER"
	DIRECT_MESSAGE EventNames = "DIRECT_MESSAGE"
	DISCONNECT     EventNames = "DISCONNECT"
)

type ConnectionPayload struct {
	EventName    EventNames `json:"eventName"`
	EventPayload any        `json:"eventPayload"`
}

func (c *ConnectionPayload) SetEventName(e EventNames) {
	c.EventName = e
}
