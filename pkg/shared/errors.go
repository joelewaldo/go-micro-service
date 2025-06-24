package shared

import (
	"github.com/sirupsen/logrus"
)

func (e CustomError) Log() {
	logrus.WithFields(logrus.Fields{
		"error_type": e.ErrorType,
		"error_id":   e.Data.ID,
		"timestamp":  e.Data.TimeStamp,
		"message":    e.Data.Message,
		"related_id": e.Data.RelatedRecordID,
	}).Error(e.ErrorLogMsg)
}

