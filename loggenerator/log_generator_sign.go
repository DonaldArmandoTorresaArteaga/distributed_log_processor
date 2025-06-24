package loggenerator

import "distributed_log_processor/configuration"

type ILogGenerator interface {
	GenerateLogRecords() error
}
type LogGenerator struct {
	configuration *configuration.Configuration
}

func NewLogGenerator(configuration *configuration.Configuration) ILogGenerator {
	return &LogGenerator{configuration: configuration}
}
