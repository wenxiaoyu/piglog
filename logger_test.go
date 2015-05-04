// logger_test
package piglog

import (
	"testing"
)

func TestLog(t *testing.T) {
	err := Log.Debug("access", "out put my test log .带中文.")
	err = Log.Trace("access", "out put my test log -->")
	err = Log.Info("access", "out put my test log -->")
	err = Log.Warn("access", "out put my test log -->")
	err = Log.Error("access", "out put my test log -->")

	err = Log.Tracef("access", "test access log by format %d", 3)
	err = Log.Debugf("access", "test access log by format %d", 3)
	err = Log.Infof("access", "test access log by format %d", 3)
	err = Log.Warnf("access", "test access log by format %d", 3)
	err = Log.Errorf("access", "test access log by format %d", 3)

	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkLog(b *testing.B) {
	b.StopTimer()
	b.StartTimer()

	for i := 1; i < b.N; i++ {
		err := Log.Debug("access", "out put my test log .带中文.")
		err = Log.Trace("access", "out put my test log -->")
		err = Log.Info("access", "out put my test log -->")
		err = Log.Warn("access", "out put my test log -->")
		err = Log.Error("access", "out put my test log -->")

		err = Log.Tracef("access", "test access log by format %d", 3)
		err = Log.Debugf("access", "test access log by format %d", 3)
		err = Log.Infof("access", "test access log by format %d", 3)
		err = Log.Warnf("access", "test access log by format %d", 3)
		err = Log.Errorf("access", "test access log by format %d", 3)
		//Log.Debug("access", "测试日志输出，压力测试测试日志输出测试日志输出，压力测试测试日志输出测试日志输出，压力测试测试日志输出测试日志输出，压力测试测试日志输出测试日志输出，压力测试测试日志输出测试日志输出，压力测试测试日志输出，压力测试测试日志输出，压力测试测试日志输出，压力测试测试日志输出，压力测试测试日志输出，压力测试测试日志输出，压力测试测试日志输出，压力测试测试日志输出，压力测试"+strconv.Itoa(i))
		if err != nil {
			b.Fail()
		}
	}
}
