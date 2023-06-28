package chatgptplugin

type Logger interface {
	Printf(format string, v ...any)
}

type dummyLogger struct{}

func (dummyLogger) Printf(format string, v ...any) {}

type PrintfFunc func(format string, v ...any)

func DummyPrintf(format string, v ...any) {}
