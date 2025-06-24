package configuration

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvironmentVariables struct {
	amountUsersPool          int64
	amountLogRecords         int64
	concurrentLogTasks       int64
	concurrentLogReaderTasks int64
	logPath                  string
}

func (e *EnvironmentVariables) GetAmountUsersPool() int64          { return e.amountUsersPool }
func (e *EnvironmentVariables) GetAmountLogRecords() int64         { return e.amountLogRecords }
func (e *EnvironmentVariables) GetConcurrentLogTasks() int64       { return e.concurrentLogTasks }
func (e *EnvironmentVariables) GetConcurrentLogReaderTasks() int64 { return e.concurrentLogReaderTasks }
func (e *EnvironmentVariables) GetLogPath() string                 { return e.logPath }

func loadEnvVariables() (*EnvironmentVariables, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	var aup int64
	if aup, err = strconv.ParseInt(os.Getenv("AMOUNT_USERS_POOL"), 10, 64); err != nil {
		return nil, err
	}

	var alr int64
	if alr, err = strconv.ParseInt(os.Getenv("AMOUNT_LOG_RECORDS"), 10, 64); err != nil {
		return nil, err
	}

	var clt int64
	if clt, err = strconv.ParseInt(os.Getenv("CONCURRENT_LOG_TASK"), 10, 64); err != nil {
		return nil, err
	}

	var clrt int64
	if clrt, err = strconv.ParseInt(os.Getenv("CONCURRENT_LOG_READER_TASK"), 10, 64); err != nil {
		return nil, err
	}

	return &EnvironmentVariables{
		amountUsersPool:          aup,
		amountLogRecords:         alr,
		concurrentLogTasks:       clt,
		concurrentLogReaderTasks: clrt,
		logPath:                  os.Getenv("LOG_PATH"),
	}, nil
}
