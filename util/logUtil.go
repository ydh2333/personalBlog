package util

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// Configuration for logging
type Config struct {
	// Enable console logging
	ConsoleLoggingEnabled bool
	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool
	// Directory to log to to when filelogging is enabled
	Directory string
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a logfile
	MaxAge int
	// Level the zerolog Level
	Level int
}

var Lg zerolog.Logger

// Configure sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func InitLogger(config Config) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var logLevel = zerolog.Level(config.Level)
	if config.Level < -1 || config.Level > 7 {
		logLevel = zerolog.InfoLevel // default to INFO
	}

	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			FieldsExclude: []string{
				"user_agent",
				"git_revision",
				"go_version",
			}})
	}

	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	mw := io.MultiWriter(writers...)

	var gitRevision string

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, v := range buildInfo.Settings {
			if v.Key == "vcs.revision" {
				gitRevision = v.Value
				break
			}
		}
	}

	Lg = zerolog.New(mw).
		Level(zerolog.Level(logLevel)).
		With().
		Str("git_revision", gitRevision).
		Str("go_version", buildInfo.GoVersion).
		Timestamp().
		Logger()

	Lg.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("jsonLogOutput", config.EncodeLogsAsJson).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Msg("logging configured")
}

func newRollingFile(config Config) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		Lg.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		cost := time.Since(start)
		defer Lg.Info().
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Str("ip", c.ClientIP()).
			Str("user-agent", c.Request.UserAgent()).
			Str("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()).
			Dur("cost", cost).Send()

		c.Next()
	}
}

// GinRecovery recovers possible project panic
func GinRecovery(stack bool) gin.HandlerFunc {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					Lg.
						Error().
						Str("path", c.Request.URL.Path).
						Any("error", err).
						Str("request", string(httpRequest)).
						Send()

					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					errors.New(string(debug.Stack()))
					Lg.
						Error().
						Stack().
						Err(errors.New(string(debug.Stack()))).
						Str("error", "[Recovery from panic]").
						Str("request", string(httpRequest)).
						Send()

				} else {
					Lg.
						Error().
						Str("error", "[Recovery from panic]").
						Any("error", err).
						Str("request", string(httpRequest)).
						Send()
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
