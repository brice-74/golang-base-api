package mocks

type Logger struct {
	PrintInfoCalled  bool
	PrintErrorCalled bool
	PrintFatalCalled bool
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) PrintInfo(_ string, _ map[string]string) {
	l.PrintFatalCalled = true
}

func (l *Logger) PrintError(_ error, _ map[string]string) {
	l.PrintErrorCalled = true
}

func (l *Logger) PrintFatal(_ error, _ map[string]string) {
	l.PrintFatalCalled = true
}
