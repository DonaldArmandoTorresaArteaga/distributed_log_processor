package loggenerator

import (
	"fmt"
	"maps"
	"math/rand"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type log struct {
	Date   string `fake:"{date}"`
	Type   string `fake:"{type}"`
	Email  string `fake:"{email}"`
	Action string `fake:"{action}"`
}

type LogTypes int

const (
	INFORMATION LogTypes = iota
	WARNING
	ERROR
)

var LogTypesValues = map[LogTypes]string{
	INFORMATION: "INFO",
	WARNING:     "WARN",
	ERROR:       "ERROR",
}

type LogActions int

const (
	LOGIN LogActions = iota
	LOGOUT
	UPLOAD
	DOWNLOAD
	DELETE
)

var LogActionsValues = map[LogActions]string{
	LOGIN:    "Login",
	LOGOUT:   "Logout",
	UPLOAD:   "Upload",
	DOWNLOAD: "Download",
	DELETE:   "Delete",
}

func (lg *LogGenerator) GenerateLogRecords() error {

	amountUsersPool := lg.configuration.EnvironmentVariables.GetAmountUsersPool()
	emailPool := make([]string, 0, amountUsersPool)

	for range amountUsersPool {
		emailPool = append(emailPool, gofakeit.Email())
	}

	gofakeit.AddFuncLookup("email", gofakeit.Info{
		Category:    "custom",
		Description: "Random email",
		Example:     "email@email.com",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {
			return f.RandomString(emailPool), nil
		},
	})

	var l log
	var wg sync.WaitGroup
	var mu sync.Mutex

	slyv := slices.Collect(maps.Values(LogTypesValues))
	slav := slices.Collect(maps.Values(LogActionsValues))
	clt := lg.configuration.EnvironmentVariables.GetConcurrentLogTasks()
	alr := lg.configuration.EnvironmentVariables.GetAmountLogRecords()
	channel := make(chan string, (alr / (alr / clt)))

	for range alr / (alr / clt) {
		wg.Add(1)
		go func(ch chan string, wg *sync.WaitGroup, mu *sync.Mutex) {
			var logs string
			for range alr / clt {
				gofakeit.Struct(&l)
				now := time.Now()
				l.Date = gofakeit.DateRange(now.AddDate(0, -1, 0), now).Format(time.RFC3339)
				mu.Lock()
				logs = logs + fmt.Sprintf(
					"%s | %s | User:%s | Action:%s \n",
					l.Date, slyv[rand.Intn(len(slyv))],
					l.Email,
					slav[rand.Intn(len(slav))],
				)
				mu.Unlock()
			}
			channel <- logs
			defer wg.Done()
		}(channel, &wg, &mu)
	}

	wg.Wait()
	close(channel)

	var logs string
	for v := range channel {
		logs += v
	}
	err := os.WriteFile(lg.configuration.EnvironmentVariables.GetLogPath(), []byte(logs), os.FileMode(0644))
	if err != nil {
		return err
	}

	return nil
}
