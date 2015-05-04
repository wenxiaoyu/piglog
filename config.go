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

var configmap = make(map[string]PigLogConfig)

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
			var config Config
			xml.Unmarshal(b, &config)
			for _, plog := range config.Piglog {
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
