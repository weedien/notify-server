package logs

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"strings"
)

type CustomFormatter struct {
}

// 定义不同级别日志的颜色
var levelColors = map[logrus.Level]string{
	logrus.DebugLevel: "\033[37m", // 灰色
	logrus.InfoLevel:  "\033[36m", // 蓝色
	logrus.WarnLevel:  "\033[33m", // 黄色
	logrus.ErrorLevel: "\033[31m", // 红色
	logrus.FatalLevel: "\033[35m", // 紫色
	logrus.PanicLevel: "\033[32m", // 青色
}

// 添加重置颜色的ANSI码
var resetColor = "\033[0m"

func (m *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	// 获取日志级别对应的颜色
	levelColor, ok := levelColors[entry.Level]
	if !ok {
		levelColor = "\033[37m" // 默认为白色
	}

	// 将 entry.Level 转换为大写
	level := strings.ToUpper(entry.Level.String())

	// 构建日志信息
	var newLog string
	newLog = fmt.Sprintf("[%s] %s%s%s \"%s\"", timestamp, levelColor, level, resetColor, entry.Message)

	// 对 entry.Data 中的键值对进行排序
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := entry.Data[k]
		newLog += fmt.Sprintf(" %s%s%s=%+v", levelColor, k, resetColor, v)
	}

	b.WriteString(newLog)
	b.WriteByte('\n')
	return b.Bytes(), nil
}
