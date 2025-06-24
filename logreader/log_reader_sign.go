package logreader

import "distributed_log_processor/configuration"

type ILogReader interface {
	LogDataCounter() (*CounterResult, error)
}

type LogReader struct {
	configuration *configuration.Configuration
}

func NewLogReader(configuration *configuration.Configuration) ILogReader {
	return &LogReader{configuration: configuration}
}
