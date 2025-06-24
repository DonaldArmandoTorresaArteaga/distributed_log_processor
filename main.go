package main

import (
	"distributed_log_processor/configuration"
	"distributed_log_processor/loggenerator"
	"distributed_log_processor/logreader"
	"distributed_log_processor/logwriter"
	"log"
)

func main() {
	configuration, err := configuration.UpConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	err = loggenerator.NewLogGenerator(configuration).GenerateLogRecords()
	if err != nil {
		log.Fatal(err)
	}

	ldc, err := logreader.NewLogReader(configuration).LogDataCounter()
	if err != nil {
		log.Fatal(err)
	}

	err = logwriter.NewLogWritter(configuration, ldc).Write()
	if err != nil {
		log.Fatal(err)
	}

}
