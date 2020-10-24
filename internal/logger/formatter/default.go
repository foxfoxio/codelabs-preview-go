

package formatter

import (
	"bytes"
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/internal/logger/constant"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type DefaultLogFormatter struct{}

func valueToString(value interface{}) string {
	var v interface{} = nil
	switch value := value.(type) {
	case string:
		v = value
	case error:
		v = value.Error()
	case fmt.Stringer:
		v = value.String()
	case fmt.GoStringer:
		v = value.GoString()
	default:
		v = value
	}
	return fmt.Sprintf("%v", v)
}
func printToBuffer(b *bytes.Buffer, value interface{}, defaultValue string) {
	if value != nil {
		b.WriteString(fmt.Sprintf("| %s ", valueToString(value)))
	} else {
		b.WriteString(fmt.Sprintf("| %s ", defaultValue))
	}
}

// Print log from custom field entry
// <RFC date> | <Level> | <Request ID> <URL> <Name> | <Service ID> | <Library information> "<Message>" "...<Field vey>=<Field value>"
func (g *DefaultLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// format data string
	var Data interface{}
	var toStringData = make([]string, 0, len(entry.Data))
	for k, v := range entry.Data {
		if k == constant.FieldData || k == constant.FieldRequestID || k == constant.FieldServiceID || k == constant.FieldServiceInfo || k == constant.FieldURL || k == constant.FieldUserID {
			continue
		}
		toStringData = append(toStringData, fmt.Sprintf("%s=%s", k, valueToString(v)))
	}
	// print data only with value
	if len(toStringData) > 0 {
		Data = strings.Join(toStringData[:], " | ")
	}

	var (
		dateTimeFormat = "2006-01-02T15:04:05.000000-0700"
		RFCDate        = time.Now().Format(dateTimeFormat)
		ServiceID      = entry.Data[constant.FieldServiceID]
		ServiceInfo    = entry.Data[constant.FieldServiceInfo]
		UserID         = entry.Data[constant.FieldUserID]
		RequestID      = entry.Data[constant.FieldRequestID]
		RequestURL     = entry.Data[constant.FieldURL]
	)

	b.WriteString(RFCDate)
	b.WriteString(" ")
	printToBuffer(b, strings.ToUpper(entry.Level.String()), "INFO")
	printToBuffer(b, RequestID, "-")
	printToBuffer(b, UserID, "system")
	if RequestURL != nil && RequestURL != "" {
		b.WriteString(" ")
		b.WriteString(fmt.Sprintf("%s ", RequestURL))
	}
	printToBuffer(b, ServiceID, "-")
	printToBuffer(b, ServiceInfo, "-")
	printToBuffer(b, entry.Message, "-")
	if Data != nil {
		printToBuffer(b, Data, "")
	}
	b.WriteString("\n")
	return b.Bytes(), nil
}
