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
	xml "encoding/xml"

	io "io/ioutil"

	"os"
	exec "os/exec"
	filepath "path/filepath"

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
var LogConfig Config

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
			xml.Unmarshal(b, &LogConfig)
			for _, plog := range LogConfig.Piglog {
				configmap[plog.Event] = plog
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
