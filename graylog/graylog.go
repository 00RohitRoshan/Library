package graylog

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"reflect"
	"time"
)

type Config struct {
	Adr      string
	Protocol string
	LogLevel string
}
type Graylog struct {
	con net.Conn
}

var logLevels = map[string]int{
	"TRACE":   1,
	"DEBUG":   2,
	"INFO":    3,
	"WARN":    4,
	"ERROR":   5,
	"FATAL":   6,
	"PANIC":   7,
}

var logLevel string
//Initialize and return a preferred connection for graylog
// "TRACE":   1,
// "DEBUG":   2,
// "INFO":    3,
// "WARN":    4,
// "ERROR":   5,
// "FATAL":   6,
// "PANIC":   7,
// Logs Below the defined level won't be sent. Higher the number Higher the level
func InitGraylog(c Config) *Graylog {
	conn, err := net.Dial(c.Protocol, c.Adr)
	if err != nil {
		panic("Cannot Dial Graylog Adress")
	}
	if level, exists := logLevels[c.LogLevel]; exists && level >= 1 && level <= 7 {
		logLevel = c.LogLevel
	} else {
		panic("Invalid log level")
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
//Resrves the static value to be set with every log when explicitly not mentioned
func (g *Graylog) SetStatic(m Log) {
	logStatic = m
}

func (g *Graylog) setStatic(m *Log) {
	srcVal := reflect.ValueOf(logStatic)
	destVal := reflect.ValueOf(m).Elem() // Pointer required to modify

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		destField := destVal.Field(i)

		if destField.CanSet() && destField.String() == "" { // Ensure field is settable
			destField.Set(srcField)
		}
	}
}


var ip string
// Checks and formats the fields which should not be empty
// Adds IP if not set
// Adds hostname if not set
// Adds N/A for anything else
func (g *Graylog) checkMustHave(m *Log) {
	if m.IPAddress == "" {
		if ip == "" {
			addrs, err := net.InterfaceAddrs()
			if err != nil {
				// return "", err
			}

			for _, addr := range addrs {
				// Check if the address is an IP address (not a MAC address)
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					// Filter only IPv4 addresses
					if ipNet.IP.To4() != nil {
						// return ipNet.IP.String(), nil
						m.IPAddress = ipNet.IP.String()
						ip = ipNet.IP.String()
					}
				}
			}
			// return ip.Query
		} else {
			m.IPAddress = ip

		}
	}
	if m.HostName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println("Error:" + err.Error())
			return
		}
		m.HostName = hostname
	}

	if m.ParamA == "" {
		m.ParamA = "N/A"
	}
	if m.ParamB == "" {
		m.ParamB = "N/A"
	}
	if m.ParamC == "" {
		m.ParamC = "N/A"
	}
	if m.BankCode == "" {
		m.BankCode = "N/A"
	}
	if m.CFTrID == "" {
		m.CFTrID = "N/A"
	}
	if m.Channel == "" {
		m.Channel = "N/A"
	}
	if m.DeviceInfo == "" {
		m.DeviceInfo = "N/A"
	}
	if m.Message == "" {
		m.Message = "N/A"
	}
	if m.PublishID == "" {
		m.PublishID = "N/A"
	}
	if m.RRN == "" {
		m.RRN = "N/A"
	}
	if m.ReferenceID == "" {
		m.ReferenceID = "N/A"
	}
	if m.Timestamp == "" {
		m.Timestamp = "N/A"
	}
	if m.TrID == "" {
		m.TrID = "N/A"
	}

}
//It writes m to the g.conn After seting the constants and checking for nil or empty values
func (g *Graylog) log(m Log) {
	if logLevels[logLevel] >= logLevels[m.Level]{
		return
	}
	m.Timestamp = time.Now().String()
	g.setStatic(&m)
	g.checkMustHave(&m)
	jsonData, err := json.Marshal(m)
	if err != nil {
		fmt.Println("err :" + err.Error())
		return
	}
	jsonBytes := []byte(jsonData)
	// defer g.con.Close()
	_, err = g.con.Write(jsonBytes)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Gray log")
}
//Writes logs with INFO level
func (g *Graylog) Info(m Log) {
	m.Level = "INFO"
	g.log(m)
}
//Writes logs with INFO level
func (g *Graylog) Debug(m Log) {
	m.Level = "DEBUG"
	g.log(m)
}
//Writes logs with INFO level
func (g *Graylog) Error(m Log) {
	m.Level = "ERROR"
	g.log(m)
}
//Writes logs with INFO level
func (g *Graylog) Fatal(m Log) {
	m.Level = "FATAL"
	g.log(m)
}
//Writes logs with INFO level
func (g *Graylog) Trace(m Log) {
	m.Level = "TRACE"
	g.log(m)
}
//Writes logs with INFO level
func (g *Graylog) Warn(m Log) {
	m.Level = "WARN"
	g.log(m)
}
//Writes logs with INFO level
func (g *Graylog) Panic(m Log) {
	m.Level = "PANIC"
	g.log(m)
}
