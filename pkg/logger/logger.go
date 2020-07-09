package logger

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	debugRedisAddChannel = "add:log-debug"
	debugRedisRmvChannel = "rmv:log-debug"
)

var opts Options

var loggers = make(map[string]*logrus.Logger)
var loggersMu sync.RWMutex

// Options contains the configuration values of the logger system
type Options struct {
	Level        string
	ReportCaller bool
	Redis        redis.UniversalClient
	Formatter    *struct {
		DisableColors bool `json:"disable_colors"`
	}
	OutPut *struct {
		Filename string `json:"filename"`
		MaxSize  int    `json:"maxsize"`
		MaxAge   int    `json:"maxage"`
	}
}

func NewOptions(level string, rc bool, redis redis.UniversalClient, formatter, output map[string]interface{}) Options {
	opt := Options{
		Level:        level,
		ReportCaller: rc,
		Redis:        redis,
	}
	if formatter != nil {
		if bytes, err := json.Marshal(formatter); err == nil {
			json.Unmarshal(bytes, &opt.Formatter)
		}
	}
	if output != nil && len(output) > 0 {
		if bytes, err := json.Marshal(output); err == nil {
			json.Unmarshal(bytes, &opt.OutPut)
		}
	}
	return opt
}

// Init initializes the logger module with the specified options.
func Init(opt Options) error {
	level := opt.Level
	if level == "" {
		level = "info"
	}
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(logLevel)

	logrus.SetReportCaller(opt.ReportCaller)

	var formatter logrus.Formatter
	if opt.OutPut != nil {
		formatter = &logrus.JSONFormatter{TimestampFormat: "2006/01/02/15:04:05.000Z07"}
	} else {
		var disableColors bool
		if opt.Formatter != nil {
			disableColors = opt.Formatter.DisableColors
		}
		formatter = &logrus.TextFormatter{DisableColors: disableColors, FullTimestamp: true, TimestampFormat: "2006/01/02/15:04:05.000Z07"}
	}
	logrus.SetFormatter(formatter)

	if opt.OutPut != nil {
		rolling := &lumberjack.Logger{
			Filename: opt.OutPut.Filename,
			MaxSize:  opt.OutPut.MaxSize,
			MaxAge:   opt.OutPut.MaxAge,
		}
		logrus.SetOutput(rolling)
	}

	if cli := opt.Redis; cli != nil {
		go subscribeLoggersDebug(cli)
	}
	opts = opt
	return nil
}

// Clone clones a logrus.Logger struct.
func Clone(in *logrus.Logger) *logrus.Logger {
	out := &logrus.Logger{
		Out:       in.Out,
		Hooks:     make(logrus.LevelHooks),
		Formatter: in.Formatter,
		Level:     in.Level,
	}
	for k, v := range in.Hooks {
		out.Hooks[k] = v
	}
	return out
}

// AddDebugDomain adds the specified domain to the debug list.
func AddDebugDomain(domain string) error {
	if cli := opts.Redis; cli != nil {
		return publishLoggersDebug(cli, debugRedisAddChannel, domain)
	}
	addDebugDomain(domain)
	return nil
}

// RemoveDebugDomain removes the specified domain from the debug list.
func RemoveDebugDomain(domain string) error {
	if cli := opts.Redis; cli != nil {
		return publishLoggersDebug(cli, debugRedisRmvChannel, domain)
	}
	removeDebugDomain(domain)
	return nil
}

// WithNamespace returns a logger with the specified nspace field.
func WithNamespace(nspace string) *logrus.Entry {
	return logrus.WithField("nspace", nspace)
}

func LogTime(a ...interface{}) func() {
	var start = time.Now().UTC()
	return func() {
		logrus.WithField("duration", strconv.FormatInt(time.Since(start).Milliseconds(), 10)+"ms").Debugln(a...)
	}
}

// WithDomain returns a logger with the specified domain field.
func WithDomain(domain string) *logrus.Entry {
	loggersMu.RLock()
	defer loggersMu.RUnlock()
	if logger, ok := loggers[domain]; ok {
		return logger.WithField("domain", domain)
	}
	return logrus.WithField("domain", domain)
}

func addDebugDomain(domain string) {
	loggersMu.Lock()
	defer loggersMu.Unlock()
	_, ok := loggers[domain]
	if ok {
		return
	}
	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	loggers[domain] = logger
}

func removeDebugDomain(domain string) {
	loggersMu.Lock()
	defer loggersMu.Unlock()
	delete(loggers, domain)
}

func subscribeLoggersDebug(cli redis.UniversalClient) {
	sub := cli.Subscribe(debugRedisAddChannel, debugRedisRmvChannel)
	for msg := range sub.Channel() {
		domain := msg.Payload
		switch msg.Channel {
		case debugRedisAddChannel:
			addDebugDomain(domain)
		case debugRedisRmvChannel:
			removeDebugDomain(domain)
		}
	}
}

func publishLoggersDebug(cli redis.UniversalClient, channel, domain string) error {
	cmd := cli.Publish(channel, domain)
	return cmd.Err()
}

// IsDebug returns whether or not the debug mode is activated.
func IsDebug(logger *logrus.Entry) bool {
	return logger.Logger.Level == logrus.DebugLevel
}
