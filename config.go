// config
//<piglogs>
//	<piglog>
//    <event>access|stat|debug|info|trace|warn|error</event>
//	  <switch>1|0</switch>
//	  <path>/home/logs/access.log</path>
//	  <split>
//		<bytime>D|H</bytime>
//		<bysize>30M</bysize>
//	  </split>
//	</piglog
//</piglogs>
package piglog

import (
	md5 "crypto/md5"
	hex "encoding/hex"
	xml "encoding/xml"
	"fmt"
	io "io/ioutil"
	http "net/http"
	"os"
	exec "os/exec"
	filepath "path/filepath"
	"strconv"
	"strings"
)

var configFileName = "piglog.xml"

type LogLeve int

const (
	Trace LogLeve = 1
	Debug         = Trace + 1
	Info          = Debug + 1
	Warn          = Info + 1
	Error         = Warn + 1
)

type Config struct {
	Piglog []PigLogConfig
	Remote RemoteConfig
}

type PigLogConfig struct {
	Event   string
	Onoff   int
	Console int
	Level   LogLeve
	Path    string
	Split   SplitConfig
}
type SplitConfig struct {
	Bytime string
	Bysize int
}

type RemoteConfig struct {
	Ip        string
	Port      int
	SecretKey string
}

var configmap = make(map[string]PigLogConfig)
var config Config

func init() {
	err := loadConfig()
	if err != nil {
		panic(err)
	}
}

//load an xml of log config
func loadConfig() (err error) {
	file, err := os.Open(configFileName)

	if err == nil {
		b, e := io.ReadAll(file)

		if e == nil {
			xml.Unmarshal(b, &config)
			for _, plog := range config.Piglog {
				configmap[plog.Event] = plog
			}
			if config.Remote.Ip != "" && config.Remote.Port > 0 {
				fmt.Print("start http service : ")
				fmt.Println(config.Remote)
				startHttpServer(config.Remote)
			}
		} else {
			return e
		}
	} else {
		return err
	}
	return nil
}

//if the event exisit ,then replaced.
func ReplaceConfig(configfile string) error {
	configFileName = configfile
	return loadConfig()
}

func GetOneConfig(event string) (PigLogConfig, bool) {
	plog, ok := configmap[event]
	return plog, ok
}
func GetCurrentPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret
}

func startHttpServer(remote RemoteConfig) {
	http.HandleFunc("/log", httpLogServer)
	iport := remote.Ip + ":" + strconv.Itoa(remote.Port)
	err := http.ListenAndServe(iport, nil)

	if err != nil {
		fmt.Errorf("ListenAndServe:", err)
	}
}

// handler log server.
func httpLogServer(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	event := req.FormValue("event")
	log := req.FormValue("log")
	clientSign := req.FormValue("sign")
	logtype := req.FormValue("type")
	if event == "" || log == "" || clientSign == "" {
		rw.Write([]byte("params losed."))
	}
	if validateSign(event, log, clientSign) {
		switch logtype {
		case "trace":
			Log.Trace(event, log)
		case "debug":
			Log.Debug(event, log)
		case "info":
			Log.Info(event, log)
		case "warn":
			Log.Warn(event, log)
		case "error":
			Log.Error(event, log)
		}
		rw.Write([]byte("ok"))
	} else {
		rw.Write([]byte("sign error."))
	}
}

//validate the sign,signed by event log params
func validateSign(event string, log string, clientSign string) bool {
	configSecretkey := config.Remote.SecretKey
	hasher := md5.New()
	hasher.Write([]byte(event + log + configSecretkey))
	computerSign := hex.EncodeToString(hasher.Sum(nil))
   
	if computerSign == clientSign{
		return true
	}
	return false
}
