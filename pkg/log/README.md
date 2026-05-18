# pkg/log

Structured file logging for Go services, built on [logrus](https://github.com/sirupsen/logrus) with automatic log rotation via [lumberjack](https://github.com/natefinch/lumberjack). Timezone is resolved from the `TZ` environment variable (default: `Asia/Jakarta`).

---

## Constructors

### `NewLogger() *Logger`

Creates a Logger with the default log directory (`/var/log` on Linux, `C:\Logs` on Windows) and no application name.

```go
logger := log.NewLogger()
logger.Say("service started")
```

### `NewLoggerWithFilename(AppName string) *Logger`

Creates a Logger that includes `AppName` in the log file name and every log line.

```go
logger := log.NewLoggerWithFilename("my-service")
logger.Say("service started")
```

### `NewLoggerWithPath(AppName, dirPath string) *Logger`

Creates a Logger writing to a custom directory. Useful for tests and services that need a configurable log directory.

```go
logger := log.NewLoggerWithPath("my-service", "/data/logs")
logger.Say("service started")
```

---

## Log file naming

```
logs_<AppName>_<YYYY-MM-DD>.log
```

Examples:
- `logs_my-service_2024-01-15.log`
- `logs_2024-01-15.log` (when AppName is empty)

---

## Rotation configuration

| Parameter   | Value  | Description                              |
|-------------|--------|------------------------------------------|
| MaxSize     | 50 MB  | Maximum size before rotation             |
| MaxBackups  | 7      | Number of old log files to retain        |
| MaxAge      | 28 days| Maximum age of old log files             |
| Compress    | false  | Rotated files are not compressed         |

---

## Log format

```
[AppName] [datetime] [level] - message - Data: fields
```

Example output:
```
[my-service] [2024-01-15 10:30:00] [info] - service started - Data: map[]
[my-service] [2024-01-15 10:30:01] [error] - connection failed - Data: map[host:db.example.com]
```

---

## Available methods

| Method                                          | Level |
|-------------------------------------------------|-------|
| `Say(msg string)`                               | Info  |
| `Sayf(fmt string, args ...interface{})`         | Info  |
| `SayWithField(msg, key string, val interface{})` | Info  |
| `SayWithFields(msg string, fields map[string]interface{})` | Info  |
| `SayError(msg string)`                          | Error |
| `SayErrorf(fmt string, args ...interface{})`    | Error |
| `SayFatal(msg string)`                          | Fatal |
| `SayFatalf(fmt string, args ...interface{})`    | Fatal |
