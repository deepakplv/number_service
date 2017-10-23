package plivolog

import (
	"fmt"
	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
	"os"
)

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
type Level int

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

type Facility int

const (
	// Facility.

	// From /usr/include/sys/syslog.h.
	// These are the same up to LOG_FTP on Linux, BSD, and OS X.
	LOG_KERN Facility = iota << 3
	LOG_USER
	LOG_MAIL
	LOG_DAEMON
	LOG_AUTH
	LOG_SYSLOG
	LOG_LPR
	LOG_NEWS
	LOG_UUCP
	LOG_CRON
	LOG_AUTHPRIV
	LOG_FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOG_LOCAL0
	LOG_LOCAL1
	LOG_LOCAL2
	LOG_LOCAL3
	LOG_LOCAL4
	LOG_LOCAL5
	LOG_LOCAL6
	LOG_LOCAL7
)

// Priority
type Priority int

const (
	// Severity.

	// From /usr/include/sys/syslog.h.
	// These are the same on Linux, BSD, and OS X.
	LOG_EMERG Priority = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

// UTCFormatter - Log formatter to print utc time
type UTCFormatter struct {
	logrus.Formatter
}

// Format - Formatter for UTCFormatter
func (u UTCFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}

type PlivoLogger struct {
	*logrus.Logger

	tagname      string
	log_filename string
	priority     Priority
	address      string
}

type LoggerOptions func(*PlivoLogger) error

// Creates a new PlivoLogger object
func New(options ...LoggerOptions) (*PlivoLogger, error) {
	jsonTimeFormatter := UTCFormatter{&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.999999",
		FieldMap: logrus.FieldMap{logrus.FieldKeyMsg: "message"}}}

	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: jsonTimeFormatter,
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	// default values
	log_facility := int(LOG_LOCAL4) | int(LOG_DEBUG)
	plog := &PlivoLogger{logger, "", "", Priority(log_facility), ""}

	// apply option parameters
	for _, op := range options {
		op(plog)
	}

	hook, err := logrus_syslog.NewSyslogHook("", plog.address, syslog.Priority(plog.priority), plog.tagname)
	if err != nil {
		return plog, err
	}
	logger.Hooks.Add(hook)
	return plog, nil
}

func OptionLevel(level Level) func(l *PlivoLogger) error {
	return func(l *PlivoLogger) error {
		l.Logger.Level = logrus.Level(level)
		return nil
	}
}

func OptionSyslogAddress(address string) func(l *PlivoLogger) error {
	return func(l *PlivoLogger) error {
		l.address = address
		return nil
	}
}

func OptionTag(tagname string) func(l *PlivoLogger) error {
	return func(l *PlivoLogger) error {
		l.tagname = tagname
		return nil
	}
}

func OptionPriorty(priority Priority) func(l *PlivoLogger) error {
	return func(l *PlivoLogger) error {
		l.priority = priority
		return nil
	}
}

func OptionLogfile(filename string) func(l *PlivoLogger) error {
	return func(l *PlivoLogger) error {
		if filename != "" {
			f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
			if err != nil {
				fmt.Errorf("PLIVOLOG::Can't open file %s for logging", filename)
				return err
			}
			l.Logger.Out = f
		}
		return nil
	}
}

type AuditLog map[string]interface{}

type AuditLogOptions func(*AuditLog)

func OptionSubject(subject string) func(a *AuditLog) {
	return func(a *AuditLog) {
		(*a)["subject"] = subject
	}
}

func OptionSubjectType(subjectType string) func(a *AuditLog) {
	return func(a *AuditLog) {
		(*a)["subject_type"] = subjectType
	}
}

func OptionCorrelationID(correlationID string) func(a *AuditLog) {
	return func(a *AuditLog) {
		(*a)["correlation_id"] = correlationID
	}
}

func OptionNotes(notes interface{}) func(a *AuditLog) {
	return func(a *AuditLog) {
		(*a)["notes"] = notes
	}
}

// AuditLog - Audit logger function
func (l *PlivoLogger) AuditLog(actor string, action string, actorType string, message string, status string,
	options ...AuditLogOptions) {
	alog := &AuditLog{}
	(*alog)["type"] = "auditlog"
	(*alog)["actor"] = actor
	(*alog)["action"] = action
	(*alog)["actor_type"] = actorType
	(*alog)["status"] = status
	logFields := logrus.Fields(*alog)

	// apply option parameters
	for _, op := range options {
		op(alog)
	}

	l.WithFields(logFields).Info(message)
}
