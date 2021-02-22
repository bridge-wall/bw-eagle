package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	_ "github.com/rhonin-cd/rhonin-eagle/utils"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 日志字段默认值
const (
	LogValueDepartment = "bw"
	LogValueVersion    = "bw-go-1.0"
	LogValueExtra      = "{\"app\":\"default\",\"env\":\"default\",\"mode\":\"default\"}"
)

// 日志配置
type LogConfig struct {
	AppPath string // 业务日志完整路径
	Extra   LogExtra
}

// 日志字段extra
type LogExtra struct {
	AppName   string `json:"app"`
	Env       string `json:"env"`
	Mode      string `json:"mode"`
	Version   string `json:"version"`
	Namespace string `json:"namespace"`
	PodName   string `json:"pod_name"`
	PodIp     string `json:"pod_ip"`
	Uri       string `json:"uri"`
}

func NewLog(config *LogConfig) (*Logger, error) {
	if logger != nil {
		return logger, nil
	}

	logExtra = &LogExtra{
		AppName:   config.Extra.AppName,
		Env:       config.Extra.Env,
		Mode:      config.Extra.Mode,
		Version:   config.Extra.Version,
		Namespace: config.Extra.Namespace,
		PodName:   config.Extra.PodName,
		PodIp:     config.Extra.PodIp,
	}

	logger = &Logger{}
	var err error
	logger.appLog, err = logger.createLogger(config.AppPath, zapcore.DebugLevel)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// 业务日志中间件
func HandlerBizLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

type Logger struct {
	appLog *zap.Logger
}

var (
	logger    *Logger
	logExtra  *LogExtra
	logHeader []string
)

func LogDebug(ctx context.Context, a interface{}) {
	if logger == nil {
		return
	}

	msg := fmt.Sprint(a)
	list := listAppField(ctx, zapcore.DebugLevel)
	logger.appLog.Debug(msg, list...)
}

func LogInfo(ctx context.Context, a interface{}) {
	if logger == nil {
		return
	}

	msg := fmt.Sprint(a)
	list := listAppField(ctx, zapcore.InfoLevel)
	logger.appLog.Info(msg, list...)
}

func LogWarn(ctx context.Context, a interface{}) {
	if logger == nil {
		return
	}

	msg := fmt.Sprint(a)
	list := listAppField(ctx, zapcore.WarnLevel)
	logger.appLog.Warn(msg, list...)
}

func LogError(ctx context.Context, a interface{}) {
	if logger == nil {
		return
	}

	msg := fmt.Sprint(a)
	list := listAppField(ctx, zapcore.ErrorLevel)
	logger.appLog.Error(msg, list...)
}

func LogCritical(ctx context.Context, a interface{}) {
	if logger == nil {
		return
	}

	msg := fmt.Sprint(a)
	list := listAppField(ctx, zapcore.FatalLevel)
	logger.appLog.Fatal(msg, list...)
}

func LogAddCustomHeader(header string) {
	if logger == nil {
		return
	}

	logHeader = append(logHeader, header)
}

func (log *Logger) createDir(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			dirName, fileName := filepath.Split(filePath)
			if dirName == "" || fileName == "" {
				return errors.New("dir or file is empty")
			}

			err = os.MkdirAll(dirName, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (log *Logger) createLogger(logPath string, level zapcore.Level) (*zap.Logger, error) {
	fileDir, fileName := filepath.Split(logPath)
	if fileDir == "" || fileName == "" {
		return nil, errors.New("file path error")
	}

	err := log.createFileDir(fileDir)
	if err != nil {
		return nil, err
	}

	fileExt := filepath.Ext(fileName)
	// 日志文件名优先使用项目名
	fileNameOnly := logExtra.AppName
	if logExtra.AppName == "" {
		fileNameOnly = strings.TrimSuffix(fileName, fileExt)
	}

	filePath := fileDir + "/" + fileNameOnly + "_%Y%m%d.log"
	hook, err := rotatelogs.New(
		filePath,
		// rotatelogs.WithLinkName("/path/to/access_log"),
		// rotatelogs.WithMaxAge(24*time.Hour),
		// rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		return nil, err
	}
	w := zapcore.AddSync(hook)

	encoderCfg := zap.NewProductionEncoderConfig()
	//encoderCfg.MessageKey = LogKeyMsg
	//encoderCfg.LevelKey = LogKeyLevel
	//encoderCfg.TimeKey = LogKeyTimestamp
	encoderCfg.NameKey = ""
	encoderCfg.CallerKey = ""
	encoderCfg.StacktraceKey = ""
	encoderCfg.LineEnding = ""
	encoderCfg.EncodeLevel = log.levelEncoder
	encoderCfg.EncodeTime = log.timeEncoder
	encoderCfg.EncodeDuration = nil
	encoderCfg.EncodeCaller = nil
	encoderCfg.EncodeName = nil

	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			w,
			level,
		),
		// zap.AddCaller(),
		// zap.AddCallerSkip(1),
	)
	defer logger.Sync()

	return logger, nil
}

func (log *Logger) createFileDir(fileDir string) error {
	_, err := os.Stat(fileDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(fileDir, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (log *Logger) levelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	// enc.AppendString("[" + l.CapitalString() + "]")
	// enc.AppendString(l.String())
	enc.AppendString(getLevelName(l))
}

func (log *Logger) timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
	enc.AppendInt64(t.Unix())
}

func (log *Logger) durationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + d.String() + "]")
}

func (log *Logger) callerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

func (log *Logger) nameEncoder(name string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(name)
}

func getLevelName(level zapcore.Level) string {
	switch level {
	case zapcore.DebugLevel:
		return "debug"
	case zapcore.InfoLevel:
		return "info"
	case zapcore.WarnLevel:
		return "warning"
	case zapcore.ErrorLevel:
		return "error"
	case zapcore.DPanicLevel:
		return "panic"
	case zapcore.PanicLevel:
		return "panic"
	case zapcore.FatalLevel:
		return "fatal"
	}

	return "unknown"
}

func listAppField(ctx context.Context, level zapcore.Level) []zap.Field {
	var list []zap.Field
	//list = append(list, zap.String(LogKeyModule, "app."+logExtra.AppName))
	return list
}

func getFieldExtra(ctx context.Context) string {
	result := LogValueExtra
	logExtra.Uri = ""
	if ctx == nil {
		extra, err := json.Marshal(logExtra)
		if err == nil {
			result = string(extra)
		}

		return result
	}

	extra, err := json.Marshal(logExtra)
	if err == nil {
		result = string(extra)
	}

	return result
}

func getFieldFileLine(skip int) (string, string) {
	caller := zapcore.NewEntryCaller(runtime.Caller(skip))
	fileLine := strings.Split(caller.TrimmedPath(), ":")
	if len(fileLine) == 2 {
		return fileLine[0], fileLine[1]
	} else if len(fileLine) > 0 {
		return fileLine[0], ""
	}

	return "", ""
}

func getCustomHeader(c *gin.Context) string {
	if c.Request == nil {
		return ""
	}

	header := make(map[string]string)
	for _, k := range logHeader {
		v := c.GetHeader(k)
		if v != "" {
			header[k] = v
		}
	}

	result, err := json.Marshal(header)
	if err == nil {
		return string(result)
	}

	return ""
}
