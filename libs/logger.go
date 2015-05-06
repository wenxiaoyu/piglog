// logger
package piglog

import (
	"errors"
	"fmt"
	io "io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//define error
var (
	getConfigError   = errors.New("get the log config fail.")
	nonLogLevelError = errors.New("no define this log level.")
)

var Log Logger

type Logger struct {
	config PigLogConfig
}
//define the log interface
type LoggerInterface interface {
	Trace(event string, log string) error
	Debug(event string, log string) error
	Info(event string, log string) error
	Wran(event string, log string) error
	Error(event string, log string) error

	Tracef(event string, log string, v ...interface{}) error
	Debugf(event string, log string, v ...interface{}) error
	Infof(event string, log string, v ...interface{}) error
	Wranf(event string, log string, v ...interface{}) error
	Errorf(event string, log string, v ...interface{}) error
}

func New(cfg PigLogConfig) *Logger {
	return &Logger{cfg}
}

func (logger *Logger) Trace(event string, log string) error {
	return logger.log(Trace, event, log)
}

func (logger *Logger) Tracef(event string, log string, v ...interface{}) error {
	log = fmt.Sprintf(log, v)
	return logger.log(Trace, event, log)
}

func (logger *Logger) Debug(event string, log string) error {
	return logger.log(Debug, event, log)
}

func (logger *Logger) Debugf(event string, log string, v ...interface{}) error {
	log = fmt.Sprintf(log, v)
	return logger.log(Debug, event, log)
}

func (logger *Logger) Info(event string, log string) error {
	return logger.log(Info, event, log)
}

func (logger *Logger) Infof(event string, log string, v ...interface{}) error {
	log = fmt.Sprintf(log, v)
	return logger.log(Info, event, log)
}

func (logger *Logger) Warn(event string, log string) error {
	return logger.log(Warn, event, log)
}

func (logger *Logger) Warnf(event string, log string, v ...interface{}) error {
	log = fmt.Sprintf(log, v)
	return logger.log(Warn, event, log)
}

func (logger *Logger) Error(event string, log string) error {
	return logger.log(Error, event, log)
}

func (logger *Logger) Errorf(event string, log string, v ...interface{}) error {
	log = fmt.Sprintf(log, v)
	return logger.log(Error, event, log)
}
//log the log,switch different conditions
func (logger *Logger) log(leve LogLeve, event string, log string) error {
	config, ok := GetOneConfig(event)
	if !ok {
		return getConfigError
	}
	logger.config = config
	tstr := logtime()

	switch leve {
	case Trace:
		log ="["+tstr+"]"+"["+event+"]"+"Trace->" +  log
	case Debug:
		log = "["+tstr+"]"+"["+event+"]"+"Debug->" +  log
	case Info:
		log = "["+tstr+"]"+"["+event+"]"+"Info->" +  log
	case Warn:
		log = "["+tstr+"]"+"["+event+"]"+"Warn->" +  log
	case Error:
		log = "["+tstr+"]"+"["+event+"]"+"Error->" +  log
	default:
		fmt.Println(log)
	}
	if config.Console == 1 {
		fmt.Println(log)
	}
	return logger.WriteLog(log + "\n")
}
//writ log to file
func (logger *Logger) WriteLog(blog string) error {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	wf, err := logger.getWriteFile()
	defer wf.Close()
	if err != nil {
		return err
	}
	_, err1 := wf.WriteString(blog)
	return err1
}

//get the current file for log write
func (logger *Logger) getWriteFile() (*os.File, error) {

	cfg := logger.config
	fstr := cfg.Path + cfg.Event + string(os.PathSeparator) + cfg.Event + ".log"
	direrr := os.MkdirAll(cfg.Path+cfg.Event+string(os.PathSeparator), 0666)
	if direrr != nil {
		return nil, direrr
	}
	f, err := os.OpenFile(fstr, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		return nil, err
	}
	finfo, err1 := f.Stat()
	if err1 != nil {
		return nil, err1
	}
	fsize := finfo.Size()
	time := time.Now()
	y := strconv.Itoa(int(time.Year()))
	M := strconv.Itoa(int(time.Month()))
	d := strconv.Itoa(int(time.Day()))
	h := strconv.Itoa(int(time.Hour()))

	newLogName := cfg.Event + ".log"
	if cfg.Split.Bytime == `D` && time.Hour() == 23 && time.Minute() == 59 && time.Second() == 59 {
		newLogName = cfg.Event + ".log." + y + "-" + M + "-" + d
	}

	if cfg.Split.Bytime == `H` && time.Minute() == 59 && time.Second() == 59 {
		newLogName = cfg.Event + ".log." + y + "-" + M + "-" + d + "-" + h
	}

	var maxNum int

	if fsize >= int64(cfg.Split.Bysize*1024*1024) && cfg.Split.Bytime == `D` {
		dri := cfg.Path + cfg.Event
		filter := cfg.Event + ".log." + y + "-" + M + "-" + d
		maxNum = getMaxSpliNum(dri, filter)
		newLogName = filter + "." + strconv.Itoa(maxNum)
	}
	if fsize >= int64(cfg.Split.Bysize*1024*1024) && cfg.Split.Bytime == `H` {
		dri := cfg.Path + cfg.Event
		filter := cfg.Event + ".log." + y + "-" + M + "-" + d + "-" + h
		maxNum = getMaxSpliNum(dri, filter)
		newLogName = filter + "." + strconv.Itoa(maxNum)
	}

	if newLogName == cfg.Event+".log" {
		return f, err
	}

	newpath := cfg.Path + cfg.Event + string(os.PathSeparator) + newLogName
	f.Close()
	reerr := os.Rename(fstr, newpath)

	if reerr == nil {
		f, err = os.OpenFile(fstr, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		return nil, reerr
	}
	return f, err

}
//get the max num of log file
func getMaxSpliNum(dir string, filter string) int {
	files, err := io.ReadDir(dir)
	if err != nil {
		return 0
	}
	var maxNum = -1
	for _, file := range files {
		fname := file.Name()
		lastdot := strings.LastIndex(fname, ".")
		if fname[:lastdot] == filter {
			tmpNum, err1 := strconv.Atoi(fname[lastdot+1:])
			if err1 != nil {
				continue
			}
			if tmpNum > maxNum {
				maxNum = tmpNum
			}
		}
	}
	return maxNum + 1
}
//the log time
func logtime() string {
	tt := time.Now()
	y := strconv.Itoa(int(tt.Year()))
	M := strconv.Itoa(int(tt.Month()))
	d := strconv.Itoa(int(tt.Day()))
	h := strconv.Itoa(int(tt.Hour()))
	m := strconv.Itoa(int(tt.Minute()))
	s := strconv.Itoa(int(tt.Second()))

	ms := strconv.Itoa(int(tt.Nanosecond() / 1e6))

	tstr :=  y + "/" + M + "/" + d + " " + h + ":" + m + ":" + s + "." + ms 
	return tstr
}
