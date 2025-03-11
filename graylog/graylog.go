package graylog

import (
	"encoding/json"
	"fmt"
	"net"

	Mylibrary "github.com/00RohitRoshan/Rohit"
)

type Config struct {
	Adr      string
	Protocol string
}
type Graylog struct {
	con net.Conn
}

func InitGraylog(c Config) *Graylog {
	conn, err := net.Dial(c.Protocol, c.Adr)
	if err != nil {
		panic("Cannot Dial Graylog Adress")
	}
	return &Graylog{
		con: conn,
	}
}

func (g *Graylog) Log(m map[string]interface{}) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	jsonBytes := []byte(jsonData)
	defer g.con.Close()
	g.con.Write(jsonBytes)
	Mylibrary.Rohit("Graylog log")
	return
}
