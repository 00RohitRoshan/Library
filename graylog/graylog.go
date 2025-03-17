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
	File  	 *os.File
}
type Graylog struct {
	con net.Conn
}

var logLevels = map[string]int{
	"TRACE": 1,
	"DEBUG": 2,
	"INFO":  3,
	"WARN":  4,
	"ERROR": 5,
	"FATAL": 6,
	"PANIC": 7,
}

var logLevel string
var file *os.File

// Initialize and return a preferred connection for graylog
//	"TRACE":   1,
//	"DEBUG":   2,
//	"INFO":    3,
//	"WARN":    4,
//	"ERROR":   5,
//	"FATAL":   6,
//	"PANIC":   7,
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
	if c.File == nil {
		panic("No file writer provided in config c `InitGraylog(c Config) *Graylog`")
	}
	file = c.File
	if err != nil {
		panic("Cannot creat file\n "+err.Error())
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
	// TrID        string `json:"tr_id"`
	// Channel     string `json:"channel"`
	// BankCode    string `json:"bank_code"`
	// ReferenceID string `json:"reference_id"`
	// RRN         string `json:"rrn"`
	// PublishID   string `json:"publish_id"`
	// CFTrID      string `json:"cf_trid"`
	// DeviceInfo  string `json:"device_info"`
	// ParamA      string `json:"param_a"`
	// ParamB      string `json:"param_b"`
	// ParamC      string `json:"param_c"`
}

var logStatic Log

//  Resrves the static value to be set with every log when explicitly not mentioned.       
//  Also can be used to override the default value for empty field from `N/A` to anything of Your choice
func (g *Graylog) SetStatic(m Log) {
	logStatic = m
	//  return logStatic
}

func (m *Log) setStatic() {
	srcVal := reflect.ValueOf(logStatic)
	destVal := reflect.ValueOf(m).Elem() //  Pointer required to modify

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		destField := destVal.Field(i)

		if destField.CanSet() && destField.String() == "" { //  Ensure field is settable
			destField.Set(srcField)
		}
	}
}

func (g *Graylog) SetDefaultEmpty(s string) {
	v := reflect.ValueOf(logStatic).Elem() //  Get the underlying struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			field.SetString(s)
		}
	}
}

var setMustHave = map[string]func() string{
	"IPAddress":   getIPAddress,
	"HostName":    getHostName,
	// "ParamA":      func() string { return "N/A" },
	// "ParamB":      func() string { return "N/A" },
	// "ParamC":      func() string { return "N/A" },
	// "BankCode":    func() string { return "N/A" },
	// "CFTrID":      func() string { return "N/A" },
	// "Channel":     func() string { return "N/A" },
	// "DeviceInfo":  func() string { return "N/A" },
	// "Message":     func() string { return "N/A" },
	// "PublishID":   func() string { return "N/A" },
	// "RRN":         func() string { return "N/A" },
	// "ReferenceID": func() string { return "N/A" },
	// "Timestamp":   func() string { return "N/A" },
	// "TrID":        func() string { return "N/A" },
}

// Function to check and set required fields dynamically
// Checks and formats the fields which should not be empty
// Adds IP if not set
// Adds hostname if not set
// Adds N/A for anything else
// Map to store field-specific setter functions
func (m *Log) checkMustHave() {
	v := reflect.ValueOf(m).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		//  If the field is empty and a set function exists, use it
		if field.Kind() == reflect.String && field.String() == "" {
			if setter, exists := setMustHave[fieldName]; exists {
				field.SetString(setter())
			}
		}
	}
}

// Function to add a new setter function dynamically
func (g *Graylog) MustHaveFuncAdd(field string, setter func() string) {
	setMustHave[field] = setter
}

// Function to remove a setter function dynamically
func (g *Graylog) MustHaveFuncRemove(field string) {
	delete(setMustHave, field)
}

// Function to update an existing setter function dynamically
func (g *Graylog) MustHaveFuncUpdate(field string, setter func() string) {
	if _, exists := setMustHave[field]; exists {
		setMustHave[field] = setter
	} else {
		fmt.Fprintf(file, "Field %s does not exist in setMustHave\n", field)
	}
}

// Function to get the IP Address
func getIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "Unknown IP"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "Unknown IP"
}

// Function to get the Host Name
func getHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Unknown Host"
	}
	return hostname
}

// It writes m to the g.conn After seting the constants and checking for nil or empty values
func (g *Graylog) log(m Log) {
	if logLevels[logLevel] > logLevels[m.Level] {
		return
	}
	m.Timestamp = time.Now().String()
	m.setStatic()
	m.checkMustHave()
	jsonData, err := json.Marshal(m)
	if err != nil {
		fmt.Fprintln(file,"err :" + err.Error())
		return
	}
	jsonBytes := []byte(jsonData)
	//  defer g.con.Close()
	_, err = g.con.Write(jsonBytes)
	if err != nil {
		fmt.Fprintln(file,err.Error())
	}
	fmt.Fprintln(file,"Gray log",m)
}

// Writes logs with INFO level
func (g *Graylog) Info(m Log) {
	m.Level = "INFO"
	g.log(m)
}

// Writes logs with DEBUG level
func (g *Graylog) Debug(m Log) {
	m.Level = "DEBUG"
	g.log(m)
}

// Writes logs with ERROR level
func (g *Graylog) Error(m Log) {
	m.Level = "ERROR"
	g.log(m)
}

// Writes logs with FATAL level
func (g *Graylog) Fatal(m Log) {
	m.Level = "FATAL"
	g.log(m)
}

// Writes logs with TRACE level
func (g *Graylog) Trace(m Log) {
	m.Level = "TRACE"
	g.log(m)
}

// Writes logs with WARN level
func (g *Graylog) Warn(m Log) {
	m.Level = "WARN"
	g.log(m)
}

// Writes logs with PANIC level
func (g *Graylog) Panic(m Log) {
	m.Level = "PANIC"
	g.log(m)
}
// get Field By index
//  1 :		"IPAddress"   ,
//  2:		"HostName"    ,
//  3:		"ParamA"      ,
//  4:		"ParamB"      ,
//  5:		"ParamC"      ,
//  6:		"BankCode"    ,
//  7:		"CFTrID"      ,
//  8:		"Channel"     ,
//  9:		"DeviceInfo"  ,
//  10:		"Message"     ,
//  11:		"PublishID"   ,
//  12:		"RRN"         ,
//  13:		"ReferenceID" ,
//  14:		"Timestamp"   ,
//  15:		"TrID"        ,
var Fields = map[int]string{
	1 :		"IPAddress"   ,
	2:		"HostName"    ,
	// 3:		"ParamA"      ,
	// 4:		"ParamB"      ,
	// 5:		"ParamC"      ,
	// 6:		"BankCode"    ,
	// 7:		"CFTrID"      ,
	// 8:		"Channel"     ,
	// 9:		"DeviceInfo"  ,
	// 10:		"Message"     ,
	// 11:		"PublishID"   ,
	// 12:		"RRN"         ,
	// 13:		"ReferenceID" ,
	// 14:		"Timestamp"   ,
	// 15:		"TrID"        ,
}
// Get Log levels by int Key
//  1:"TRACE" ,
//  2:"DEBUG" ,
//  3:"INFO"  ,
//  4:"WARN"  ,
//  5:"ERROR" ,
//  6:"FATAL" ,
//  7:"PANIC" ,
var LogLevels = map[int]string{
	1:"TRACE" ,
	2:"DEBUG" ,
	3:"INFO"  ,
	4:"WARN"  ,
	5:"ERROR" ,
	6:"FATAL" ,
	7:"PANIC" ,
}