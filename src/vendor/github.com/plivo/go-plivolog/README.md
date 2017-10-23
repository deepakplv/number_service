# PlivoLog Helper in Golang

PlivoLog helper is a Go package that can be integrated with Golang projects to generate standardized logs. This does:

* Centralized logging configuration
* JSON log format for regular logs, which is easily parsable by Scalyr etc
* Audit logs: These are special type of logs that can be used to build a user or business story. This can be used to capture user progress, start/progress/stop status of batch jobs, actions performed by admin users on user accounts etc

A rule of thumb for using audit logs: The log should make sense to a support team member i.e. they are not debug logs.

## Installation and Usage

### Get package (Install)

    go get github.com/plivo/go-plivolog

### Usage

```golang
import "github.com/plivo/go-plivolog/plivolog"

// Defaults - Create a logger for console and syslog and use it
l, _ := plivolog.New()
l.Infof("Test %s", "logging stuff")

// Create an audit log
l.AuditLog("MA1234567890", "MAKE_CALL", "API", "Call API initiated", "SUCCESS")

// Audit log with subject and notes
notes := map[string]interface{}{"card_id": "ca_123456789"}
l.AuditLog("support-user@plivo.com",
	"VERIFY_CC",
	"ADMIN",
	"Verifying credit card for user",
	"SUCCESS",
	plivolog.OptionSubject("user@email.com"),
	plivolog.OptionSubjectType("USER"),
	plivolog.OptionNotes(notes))

// Optional - Set log level to info while logging to a logfile
l, _ = plivolog.New(plivolog.OptionLogfile("/tmp/logfile"))
l.Error("This goes to both syslog and logfile")

// Optional - customize syslog
log_facility := plivolog.Priority(int(plivolog.LOG_LOCAL4) | int(plivolog.LOG_INFO))
l, _ = plivolog.New(plivolog.OptionSyslogAddress("remote.ip.address"),
	plivolog.OptionPriorty(log_facility))
l.Warn("This goes to both syslog on local4 facility")
```

By default all logs will go to syslog. `PlivoLogger.logger` is the standard logrus logger object.

Auditlog takes following fields as parameters:

```golang
actor           string      // Auth Id/Unique Id to identify the service/application from which the log is generated
action          string      // Action that the user is performing. Eg.: VERIFY_CREDIT_CARD
actorType       string      // Type of the actor 

subject         string      // Auth Id/User Id/Unique Id to identify the user performing the action
subjectType     string      // Type of the subject performing the action. Eg.: User account, cron-service etc.

status          string      // Success / fail

message         string      // Log message. Description of the action performed

correlationId    string      // Unique Identification which can be used to tie logs across apps. Eg. api_id from
plivoapi
notes           interface{} // Any extra parameters that might be useful for analysis related to the action performed. Recommended to store a dict/map with key-value pairs. 
```

### Sample Output

Audit log (prettified):

    {
        "type": "auditlog",
        "action": "VERIFY_CC",
        "actor": "MAFSDSFS42MMOL",
        "actor_type": "UI",
        "level": "info",
        "status": "success",
        "message": "This is the message",
        "corelation_id": "long-api-id-or-call-uuid",
        "notes": {
            "key1": "value1",
            "key2": 10,
            "key3": ["fdsf","fsf"]
        },
        "subject": "100",
        "subject_type": "account",
        "time": "2017-04-12 14:00:22"
    }

Regular log:

    {"message": "This is an error", "time": "2017-06-07 14:41:01.479173", "level": "ERROR"}
