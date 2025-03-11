package graylog

import (
	"encoding/json"
	"net"

	"github.com/00RohitRoshan/Rohit/Mylibrary"
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
		Mylibrary.Console("err :" + err.Error())
		return
	}
	jsonBytes := []byte(jsonData)
	// defer g.con.Close()
	_, err = g.con.Write(jsonBytes)
	if err != nil {
		Mylibrary.Console(err.Error())
	}
	Mylibrary.Console("Gray log")
}
