package zlog

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/DavidGamba/go-getoptions"
	"github.com/rs/zerolog"
	"golang.org/x/exp/maps"

	"github.com/DRVTiny/lib-unilog/unilog"
)

type Zloger struct {
	optsLogLevel      string
	optsLogFile       string
	optsHumanReadable bool
}

func NewWithOpts(opts *getoptions.GetOpt) (zl *Zloger) {
	zl = &Zloger{}
	zl.AddGetoptKeys(opts)
	return
}

func New(logLevel string, logFile string, humanize bool) (zl *Zloger) {
	zl = &Zloger{
		optsLogLevel:      logLevel,
		optsLogFile:       logFile,
		optsHumanReadable: humanize,
	}

	return
}

type ZluniLog zerolog.Logger

func GetUnilog(zll zerolog.Logger) ZluniLog {
	return ZluniLog(zll)
}

func (zl ZluniLog) Fatal(args ...any) {
	zul := zerolog.Logger(zl)
	zul.Fatal().Msg(unilog.Anys2LogString(args))
}

func (zl ZluniLog) Fatalf(fstr string, args ...any) {
	zul := zerolog.Logger(zl)
	zul.Fatal().Msgf(fstr, args...)
}

func (zl ZluniLog) Error(args ...any) {
	zul := zerolog.Logger(zl)
	zul.Error().Msg(unilog.Anys2LogString(args))
}

func (zl ZluniLog) Errorf(fstr string, args ...any) {
	zul := zerolog.Logger(zl)
	zul.Error().Msgf(fstr, args...)
}

func (zl ZluniLog) Warn(args ...any) {
	zul := zerolog.Logger(zl)
	zul.Warn().Msg(unilog.Anys2LogString(args))
}

func (zl ZluniLog) Warnf(fstr string, args ...any) {
	zul := zerolog.Logger(zl)
	zul.Warn().Msgf(fstr, args...)
}

func (zl ZluniLog) Info(args ...any) {
	zul := zerolog.Logger(zl)
	zul.Info().Msg(unilog.Anys2LogString(args))
}

func (zl ZluniLog) Infof(fstr string, args ...any) {
	zul := zerolog.Logger(zl)
	zul.Info().Msgf(fstr, args...)
}

func (zl ZluniLog) Debug(args ...any) {
	zul := zerolog.Logger(zl)
	zul.Debug().Msg(unilog.Anys2LogString(args))
}

func (zl ZluniLog) Debugf(fstr string, args ...any) {
	zul := zerolog.Logger(zl)
	zul.Debug().Msgf(fstr, args...)
}

func (zl ZluniLog) Print(args ...any) {
	zl.Debug(args...)
}

func (zl ZluniLog) Printf(fstr string, args ...any) {
	zl.Debugf(fstr, args...)
}

func (zl ZluniLog) Println(args ...any) {
	zul := zerolog.Logger(zl)
	zul.Debug().Msg(unilog.Anys2LogString(args) + "\n")
}

var loglevels = map[string]zerolog.Level{
	"fatal": zerolog.FatalLevel,
	"error": zerolog.ErrorLevel,
	"warn":  zerolog.WarnLevel,
	"info":  zerolog.InfoLevel,
	"debug": zerolog.DebugLevel,
	"trace": zerolog.TraceLevel,
}

type knownLogLevels struct {
	ll []string
	m  *sync.Mutex
}

var knownLL = knownLogLevels{m: &sync.Mutex{}}

func KnownLogLevels() []string {
	if len(knownLL.ll) > 0 {
		return knownLL.ll
	} else {
		knownLL.m.Lock()
		defer knownLL.m.Unlock()
		knownLL.ll = maps.Keys(loglevels)
		n := len(knownLL.ll)
		for i := 0; i < n; i++ {
			knownLL.ll = append(knownLL.ll, strings.ToUpper(knownLL.ll[i]))
		}

	}
	return knownLL.ll
}

func (zl *Zloger) AddGetoptKeys(opt *getoptions.GetOpt) {
	opt.StringVar(
		&zl.optsLogLevel, "log-level", "error",
		opt.Alias("L"), opt.Description("set maximum level of emited log messages"),
		opt.ValidValues(KnownLogLevels()...),
	)

	opt.StringVar(
		&zl.optsLogFile, "log-file", "-",
		opt.Alias("l"), opt.Description(`file to write logs. use "-" to output to console (STDERR)`),
	)

	opt.BoolVar(
		&zl.optsHumanReadable, "log-for-humans", false,
		opt.Description(`write logs in human-readable format. this option only applicable of you uses console output specified with -l - (default)`),
	)
}

func (zl *Zloger) GetLogger() (zll zerolog.Logger, err error) {
	var logFile *os.File

	switch zl.optsLogFile {
	case "-", "STDERR":
		logFile = os.Stderr
	case "STDOUT":
		logFile = os.Stdout
	case "":
		return zerolog.Logger{}, fmt.Errorf(`log file path can not be represented as empty string. maybe, you need to use "-" instead?`)
	default:
		logFile, err = os.OpenFile(zl.optsLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return zerolog.Logger{}, fmt.Errorf("failed to open log file %s for append: %v", zl.optsLogFile, err)
		}
	}

	return GetLogger(zl.optsLogLevel, logFile, zl.optsHumanReadable)
}

func (zl *Zloger) GetUniLogger() (unilog.UniLogger, error) {
	zll, err := zl.GetLogger()
	if err != nil {
		return nil, err
	} else {
		return ZluniLog(zll), nil
	}
}

func LogLevel_Str2V(ll string) (*zerolog.Level, error) {
	if zll, ok := loglevels[strings.ToLower(ll)]; ok {
		return &zll, nil
	} else {
		return nil, fmt.Errorf("invalid/unknown log level specified: %s", ll)
	}
}

func GetLogger(logLevel string, logFile *os.File, humanize bool) (zlog zerolog.Logger, err error) {
	zlogLevel, err := LogLevel_Str2V(logLevel)
	if err != nil {
		return
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano

	zerolog.SetGlobalLevel(*zlogLevel)

	if logFile == nil {
		logFile = os.Stderr
	}

	if humanize {
		output := zerolog.ConsoleWriter{Out: logFile, TimeFormat: "2006-01-02 15:04:05.000"}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		zlog = zerolog.New(output).With().Timestamp().Logger()
	} else {
		zlog = zerolog.New(logFile).With().Timestamp().Logger()
	}

	return
}
