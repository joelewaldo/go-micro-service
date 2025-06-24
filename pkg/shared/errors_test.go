package shared

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestCustomErrorLog(t *testing.T) {
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableHTMLEscape: true,
	})
	logrus.SetLevel(logrus.ErrorLevel)

	err := CustomError{
		ErrorType:   EnumErrorTypes.Internal,
		ErrorLogMsg: "Test logging internal error",
		Data: ErrorData{
			ID:              "TEST001",
			TimeStamp:       time.Now().Format(time.RFC3339),
			Message:         "Something went wrong",
			RelatedRecordID: "1234",
		},
	}

	err.Log()
	var logOutput map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logOutput); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	// Assertions
	if logOutput["error_type"] != EnumErrorTypes.Internal {
		t.Errorf("expected error_type %q, got %q", EnumErrorTypes.Internal, logOutput["error_type"])
	}
	if logOutput["error_id"] != "TEST001" {
		t.Errorf("expected error_id 'TEST001', got %q", logOutput["error_id"])
	}
	if logOutput["msg"] != "Test logging internal error" {
		t.Errorf("expected msg 'Test logging internal error', got %q", logOutput["msg"])
	}
}
