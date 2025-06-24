
# Distributed Log Processor

## Overview
This project implements a **concurrent log‑processing pipeline** in Go that simulates the kind of batch analytics system used by large‑scale platforms (e.g., Netflix, Cloudflare) to crunch millions of application logs.

The pipeline executes three main stages:

1. **Log Generation** – concurrently produces **N** synthetic log lines and streams them to disk.
2. **Log Reading & Counting** – reads the generated file with a bounded worker pool, parses each line, and aggregates:
   - total events per log level (`INFO`, `WARN`, `ERROR`)
   - unique user count
   - number of actions per user
3. **Result Writing** – marshals the aggregated counters into a `results.json` file.

Everything is **configurable via environment variables** allowing you to scale the workload without touching code.

---

## Repository Layout

```
.
├── main.go                       # Program entry‑point
├── .env                          # Runtime configuration
├── configuration/                # Loads & validates env vars
│   ├── configuration.go
│   └── env_variables.go
├── loggenerator/                 # Stage 1
│   ├── log_generator.go
│   └── log_generator_sign.go
├── logreader/                    # Stage 2
│   ├── log_reader.go
│   └── log_reader_sign.go
└── logwriter/                    # Stage 3
    ├── log_writer.go
    └── log_writer_sign.go
```

---

## Concurrency Model

| Stage | Concurrency Primitive              | Purpose                                                    |
|-------|------------------------------------|------------------------------------------------------------|
| **Log Generation** | Goroutine pool (`CONCURRENT_LOG_TASK`) + buffered channel | Generates log batches in parallel and pushes them to a shared channel for serialization |
| **Log Reading** | Worker pool (`CONCURRENT_LOG_READER_TASK`) + job channel + `sync.Mutex` | Bounds goroutine count while safely incrementing shared counters (log‑level map, per‑user map) |
| **Result Writing** | Single goroutine | Performs blocking I/O to write the final JSON file |

`sync.WaitGroup` barriers guarantee each stage finishes before the next begins.

---

## Environment Variables (`.env`)

| Variable | Description | Example |
|----------|-------------|---------|
| `AMOUNT_USERS_POOL` | Size of fake user email pool | `10000` |
| `AMOUNT_LOG_RECORDS` | Total log lines to generate | `1000000` |
| `CONCURRENT_LOG_TASK` | Goroutines used for log generation | `200` |
| `CONCURRENT_LOG_READER_TASK` | Workers used for log parsing/aggregation | `100` |
| `LOG_PATH` | Target path for the raw log file (and JSON output) | `./logs` |

Create or edit `.env` in the project root to experiment with different scales.

---

## Building & Running

```bash
# 1. Clone
git clone https://github.com/your‑handle/distributed_log_processor.git
cd distributed_log_processor

# 2. Bootstrap Go modules
go mod download

# 3. Set up .env  (optional – defaults provided)
cp .env.example .env && vim .env

# 4. Run
go run ./...
```

Upon success you will see:

```
Generating 1,000,000 logs…
Parsing logs with 100 workers…
Wrote ./logs          (raw logs)
Wrote ./logs.json     (aggregated results)
Done in 14.3 s
Peak goroutines: 302
```

> **Note:** Actual timings vary based on CPU cores and disk speed.

---

## Sample Output (`logs.json`)

```json
{
  "log_levels": {
    "INFO": 734123,
    "WARN": 193212,
    "ERROR": 72765
  },
  "unique_users": 9974,
  "user_actions": {
    "john.doe@example.com": 102,
    "jane.smith@example.com": 321
  }
}
```

---

## Performance Tuning

| Dial | Effect |
|------|--------|
| Increase `CONCURRENT_LOG_TASK` | Faster generation but higher memory usage |
| Increase `CONCURRENT_LOG_READER_TASK` | Faster parsing up to CPU core limit |
| Decrease `AMOUNT_LOG_RECORDS` | Quick local runs for development |

Use `runtime.NumGoroutine()` or `pprof` (`go tool pprof`) to inspect goroutine counts and hotspots.

---

## Extending the Project

* **Streaming Mode** – swap disk I/O for a message bus (Kafka, NATS) and turn the batch pipeline into a real‑time stream processor.
* **Additional Metrics** – add latency percentiles, per‑action rates, or elastic window aggregation.
* **Resilience** – integrate context cancellation, retries, and graceful shutdown hooks.

---

## License

MIT © 2025 Your Name
