package graylog

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"time"

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

type Log struct {
	Timestamp   string `json:"timestamp"`
	Level       string `json:"level"`
	Message     string `json:"message"`
	IPAddress   string `json:"ip_address"`
	AppName     string `json:"appname"`
	HostName    string `json:"hostname"`
	TrID        string `json:"tr_id"`
	Channel     string `json:"channel"`
	BankCode    string `json:"bank_code"`
	ReferenceID string `json:"reference_id"`
	RRN         string `json:"rrn"`
	PublishID   string `json:"publish_id"`
	CFTrID      string `json:"cf_trid"`
	DeviceInfo  string `json:"device_info"`
	ParamA      string `json:"param_a"`
	ParamB      string `json:"param_b"`
	ParamC      string `json:"param_c"`
}

var logStatic Log

func (g *Graylog) SetStatic(m Log) {
	logStatic = m
}

func (g *Graylog) setStatic(m *Log) {
	srcVal := reflect.ValueOf(logStatic)
	destVal := reflect.ValueOf(m).Elem() // Pointer required to modify

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		destField := destVal.Field(i)

		if destField.CanSet() && srcField.String() != "" { // Ensure field is settable
			destField.Set(srcField)
		}
	}
}

type IP struct {
    Query string
}

var ip string
func (g *Graylog) checkMustHave(m *Log) {
	if m.IPAddress == "" {
		if ip == ""{
			req, err := http.Get("http://ip-api.com/json/")
			if err != nil {
				Mylibrary.Console(err.Error())
			}
			defer req.Body.Close()

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				Mylibrary.Console(err.Error())
			}

			var ip IP
			json.Unmarshal(body, &ip)
			m.IPAddress = ip.Query
			// return ip.Query
			}else{
			m.IPAddress = ip

		}
	}
	if m.HostName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			Mylibrary.Console("Error:" + err.Error())
			return
		}
		m.HostName = hostname
	}

}

func (g *Graylog) Log(m Log) {
	m.Timestamp = time.Now().String()
	g.setStatic(&m)
	g.checkMustHave(&m)
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

func (g *Graylog) Info(m Log) {
	m.Level = "Info"
	g.Log(m)
}
func (g *Graylog) Debug(m Log) {
	m.Level = "Debug"
	g.Log(m)
}
func (g *Graylog) Error(m Log) {
	m.Level = "Error"
	g.Log(m)
}
func (g *Graylog) Fatal(m Log) {
	m.Level = "Fatal"
	g.Log(m)
}
func (g *Graylog) Trace(m Log) {
	m.Level = "Trace"
	g.Log(m)
}
func (g *Graylog) Warn(m Log) {
	m.Level = "warn"
	g.Log(m)
}
func (g *Graylog) Panic(m Log) {
	m.Level = "Panic"
	g.Log(m)
}

