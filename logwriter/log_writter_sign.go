package logwriter

import (
	"distributed_log_processor/configuration"
	"distributed_log_processor/logreader"
)

type LogWritter struct {
	configuration *configuration.Configuration
	counterResult *logreader.CounterResult
}

type ILogWritter interface {
	Write() error
}

func NewLogWritter(
	configuration *configuration.Configuration,
	counterResult *logreader.CounterResult,
) ILogWritter {
	return &LogWritter{
		configuration: configuration,
		counterResult: counterResult,
	}
}
