package logreader

import (
	"bufio"
	"distributed_log_processor/loggenerator"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"
)

type log struct {
	Date   string
	Type   string
	Email  string
	Action string
}

type LogLevels struct {
	INFO  int
	WARN  int
	ERROR int
}

type CounterResult struct {
	LogLevels    *LogLevels     `json:"log_levels"`
	UniqueUsers  int            `json:"unique_users"`
	UsersActions map[string]int `json:"user_actions"`
}

func (lr *LogReader) LogDataCounter() (*CounterResult, error) {

	ch, err := logLoader(
		lr.configuration.EnvironmentVariables.GetLogPath(),
		int(lr.configuration.EnvironmentVariables.GetAmountLogRecords()),
	)
	close(ch)

	if err != nil {
		return nil, err
	}

	bch := make(chan log, lr.configuration.EnvironmentVariables.GetConcurrentLogReaderTasks())
	mllc := map[string]int{
		loggenerator.LogTypesValues[loggenerator.INFORMATION]: 0,
		loggenerator.LogTypesValues[loggenerator.ERROR]:       0,
		loggenerator.LogTypesValues[loggenerator.WARNING]:     0,
	}
	muu := make(map[string]int)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for v := range ch {
		bch <- v
		wg.Add(1)
		go func(l log, bch chan log, mllc map[string]int, mu *sync.Mutex, wg *sync.WaitGroup) {
			defer wg.Done()
			mu.Lock()
			mllc[l.Type] = mllc[l.Type] + 1
			if _, ok := muu[l.Email]; ok {
				muu[l.Email] = muu[l.Email] + 1
			} else {
				muu[l.Email] = 1
			}
			mu.Unlock()
			<-bch
		}(v, bch, mllc, &mu, &wg)
	}
	wg.Wait()
	return &CounterResult{
		&LogLevels{
			INFO:  mllc[loggenerator.LogTypesValues[loggenerator.INFORMATION]],
			WARN:  mllc[loggenerator.LogTypesValues[loggenerator.WARNING]],
			ERROR: mllc[loggenerator.LogTypesValues[loggenerator.ERROR]],
		},
		len(slices.Collect(maps.Keys(muu))),
		muu,
	}, nil
}

func logLoader(logPath string, amountLogRecords int) (chan log, error) {

	file, err := os.Open(logPath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	ch := make(chan log, amountLogRecords)

	for scanner.Scan() {
		sp := strings.Split(scanner.Text(), "|")
		ch <- log{
			Date:   sp[0],
			Type:   strings.TrimSpace(sp[1]),
			Email:  strings.TrimSpace(strings.Split(sp[2], ":")[1]),
			Action: strings.TrimSpace(strings.Split(sp[3], ":")[1]),
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ch, nil
}
