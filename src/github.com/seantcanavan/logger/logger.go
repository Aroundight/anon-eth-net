package logger

import (
	"bufio"
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"time"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/utils"
)

const LOG_EXTENSION = ".log"

// Logger allows for aggressive log management in scenarios where disk space
// might be limited. You can limit based on log message count or duration and
// also prune log files when too many are saved on disk. I don't like the
// current implementation - it makes too many assumptions about how many
// simultaneous loggers are running and never checks available disk space.
// I need to consider keeping a list of all loggers currently executing and
// keep the option to manually prune and report logs as disk space runs low.
// This could be managed via a watchdog go routine. Ideas here:
// https://stackoverflow.com/questions/20108520/get-amount-of-free-disk-space-using-go
// This would essentially be the nail in the coffin for windows builds if
// implemented. RIP Bill Gates.
type Logger struct {
	MaxLogFileCount    uint64        // The maximum number of log files saved to disk before pruning occurs
	MaxLogMessageCount uint64        // The maximum number of bytes a log file can take up before it's cut off and a new one is created
	MaxLogDuration     uint64        // The maximum number of seconds a log can exist for before it's cut off and a new one is created
	baseLogName        string        // The beginning text to append to this log instance for naming and management purposes
	logFileCount       uint64        // The current number of logs that have been created
	logFileNames       list.List     // The list of log files we're currently holding on to
	logMessageCount    uint64        // The current number of messages that have been logged
	logDuration        uint64        // The duration, in seconds, that this log has been logging for
	logStamp           uint64        // The time when this log was last written to in unix time
	log                *os.File      // The file that we're logging to
	writer             *bufio.Writer // our writer we use to log to the current log file
}

func FromVolatilityValue(logBaseName string) (*Logger, error) {
	volatility := config.Cfg.LogVolatility
	switch volatility {
	case 0:
		return HoardingLogger(logBaseName)
	case 1:
		return AnticonservativeLogger(logBaseName)
	case 2:
		return ConservativeLogger(logBaseName)
	case 3:
		return MinimalLogger(logBaseName)
	default:
		return nil, fmt.Errorf("The value you gave: %d does not have a logging map available for it. Please check logger/logger.go for valid logging mapping values", volatility)
	}
}

func CustomLogger(logBaseName string, maxFileCount uint64, maxMessageCount uint64, maxDuration uint64) (*Logger, error) {

	sl := Logger{}
	// public variables
	sl.MaxLogFileCount = maxFileCount
	sl.MaxLogMessageCount = maxMessageCount
	sl.MaxLogDuration = maxDuration

	err := sl.initLogger(logBaseName)
	if err != nil {
		return nil, err
	}

	return &sl, nil
}

// MinimalLogger will return a Logger struct which will make sure to be as
// minimal as possible and only hold on to the most recent of log files for a
// short period of time. Recommended for systems with 100GB or less of overall
// storage. If the logs are not checked or reported via email daily it's
// possible that data could be missed.
func MinimalLogger(logBaseName string) (*Logger, error) {

	sl := Logger{}
	// public variables
	sl.MaxLogFileCount = 10
	sl.MaxLogMessageCount = 1000
	sl.MaxLogDuration = 86400 // one day in seconds

	err := sl.initLogger(logBaseName)
	if err != nil {
		return nil, err
	}

	return &sl, nil
}

// ConservativeLogger will return a Logger struct which will hold on to a
// respectful number of log files and messages. Recommended for systems with
// at least 250GB of overall storage. If the logs are not checked or reported
// via email at least every three days it's possible that data could be missed.
func ConservativeLogger(logBaseName string) (*Logger, error) {

	sl := Logger{}
	// public variables
	sl.MaxLogFileCount = 100
	sl.MaxLogMessageCount = 5000
	sl.MaxLogDuration = 259200 // three days in seconds

	err := sl.initLogger(logBaseName)
	if err != nil {
		return nil, err
	}

	return &sl, nil
}

// AnticonservativeLogger is not a political message. It will return a
// Logger struct which will hold a large number of log files and messages.
// Recommended for systems with at least 500GB of overall storage. If the logs
// are not checked or reported via email at least every five days it's possible
// that data could be missed.
func AnticonservativeLogger(logBaseName string) (*Logger, error) {

	sl := Logger{}
	// public variables
	sl.MaxLogFileCount = 1000
	sl.MaxLogMessageCount = 10000
	sl.MaxLogDuration = 432000 // five days in seconds

	err := sl.initLogger(logBaseName)
	if err != nil {
		return nil, err
	}

	return &sl, nil
}

// HoardingLogger will return a Logger struct which will hoard a massive
// amount of logs and messages. Recommended for systems with at least 1TB of
// overall storage. If the logs are not checked or reported via email at least
// every week it's possible that data could be missed.
func HoardingLogger(logBaseName string) (*Logger, error) {

	sl := Logger{}
	// public variables
	sl.MaxLogFileCount = 5000
	sl.MaxLogMessageCount = 10000
	sl.MaxLogDuration = 604800

	err := sl.initLogger(logBaseName)
	if err != nil {
		return nil, err
	}

	return &sl, nil
}

// Flush manually flushes the IO buffer to the local disk to ensure that any
// cached log messages are permanently stored onto the local disk. Useful to
// call before a logger goes out of scope.
func (sl *Logger) Flush() {
	sl.writer.Flush()
}

func (sl *Logger) CurrentLogContents() ([]byte, error) {
	sl.writer.Flush()

	fileBytes, readErr := ioutil.ReadFile(sl.log.Name())
	if readErr != nil {
		return nil, readErr
	}
	return fileBytes, nil
}

func (sl *Logger) CurrentLogName() (string, error) {
	fileInfo, statErr := sl.log.Stat()
	if statErr != nil {
		return "", statErr
	}
	return fileInfo.Name(), nil
}

func (sl *Logger) initLogger(logBaseName string) error {

	logFileName := utils.TimeStampFileName(logBaseName, LOG_EXTENSION)

	filePtr, err := os.Create(logFileName)
	if err != nil {
		return err
	}

	// private variables
	sl.baseLogName = logBaseName
	sl.logFileCount = 0
	sl.logDuration = 0
	sl.logStamp = uint64(time.Now().Unix())
	sl.log = filePtr
	sl.writer = bufio.NewWriter(sl.log)
	sl.logFileNames.PushBack(logFileName)
	return nil
}

// LogMessage will write the given string to the log file. It will then perform
// all the necessary checks to make sure that the max number of messages, the
// max duration of the log file, and the maximum number of overall log files
// has not been reached. If any of the above parameters have been tripped,
// log cleanup will occur.
func (sl *Logger) LogMessage(formatString string, values ...interface{}) {

	// what time is it right now?
	now := uint64(time.Now().Unix())

	// write the logging message to the current log file
	fmt.Fprintln(sl.writer, fmt.Sprintf(formatString, values...))
	// write the logging message to std.out for local watchers
	fmt.Println(fmt.Sprintf(formatString, values...))

	sl.logMessageCount++
	sl.logDuration += now - sl.logStamp
	sl.logStamp = now

	if sl.logMessageCount >= sl.MaxLogMessageCount ||
		sl.logDuration >= sl.MaxLogDuration {
		sl.newFile()
	}
}

// newFile generates a new log file to store the log messages within. It
// intelligently keeps track of the number of log files that have already been
// created so that you don't overload your disk with logs and can 'prune' extra
// logs as necessary.
func (sl *Logger) newFile() error {

	logFileName := utils.TimeStampFileName(sl.baseLogName, LOG_EXTENSION)

	filePtr, err := os.Create(logFileName)
	if err != nil {
		return err
	}

	sl.writer.Flush()
	sl.log.Close()

	sl.log = filePtr
	sl.writer = bufio.NewWriter(sl.log)

	sl.logMessageCount = 0
	sl.logFileCount++
	sl.logFileNames.PushBack(logFileName)

	if sl.logFileCount >= sl.MaxLogFileCount {
		if err := sl.pruneFile(); err != nil {
			return err
		}
	}

	return nil
}

// pruneFile will remove the oldest file handle from the queue and delete the
// file from the local file system.
func (sl *Logger) pruneFile() error {

	oldestLog := sl.logFileNames.Remove(sl.logFileNames.Front())
	logFileName := reflect.ValueOf(oldestLog).String()

	fmt.Println(fmt.Sprintf("Deleting oldest log file: %v", logFileName))

	return os.Remove(logFileName)
}
