// config_test
package piglog

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {

	if config, ok := GetOneConfig("access"); ok {
		if config.Onoff != 1 {
			t.Fatal("access onoff accept : 1")
		}

		if config.Split.Bysize != 1 {
			t.Fatal("access split size accept : 1")
		}
		if config.Split.Bytime != `D` {
			t.Fatal("access split time accept : D")
		}
	} else {
		t.Fatal("load access config fail.")
	}
}
