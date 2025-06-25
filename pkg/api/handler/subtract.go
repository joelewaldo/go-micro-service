package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/joelewaldo/go-micro-service/pkg/shared"
)

type SubtractRequest struct {
	Minuend    float64 `json:"minuend"`
	Subtrahend float64 `json:"subtrahend"`
}

type SubtractResponse struct {
	Result float64 `json:"result"`
}

func SubtractHandler(w http.ResponseWriter, r *http.Request) {
	minuendStr := r.PathValue("minuend")
	minuend, err := strconv.ParseFloat(minuendStr, 64)
	if err != nil {
		writeCustomError(w,
			shared.CustomError{
				ErrorType:   shared.EnumErrorTypes.BadRequest,
				ErrorLogMsg: "failed to parse minuend",
				Data: shared.ErrorData{
					ID:              uuid.NewString(),
					TimeStamp:       time.Now().Format(time.RFC3339),
					Message:         "Invalid minuend: must be a number",
					RelatedRecordID: ""},
			},
			http.StatusBadRequest,
		)
		return
	}

	subtrahendStr := r.PathValue("subtrahend")
	subtrahend, err := strconv.ParseFloat(subtrahendStr, 64)
	if err != nil {
		writeCustomError(w,
			shared.CustomError{
				ErrorType:   shared.EnumErrorTypes.BadRequest,
				ErrorLogMsg: "failed to parse subtrahend",
				Data: shared.ErrorData{
					ID:              uuid.NewString(),
					TimeStamp:       time.Now().Format(time.RFC3339),
					Message:         "Invalid subtrahend: must be a number",
					RelatedRecordID: "",
				},
			},
			http.StatusBadRequest,
		)
		return
	}

	result := minuend - subtrahend

	resp := SubtractResponse{Result: result}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeCustomError(w,
			shared.CustomError{
				ErrorType:   shared.EnumErrorTypes.Internal,
				ErrorLogMsg: "json encoding failed",
				Data: shared.ErrorData{
					ID:              uuid.NewString(),
					TimeStamp:       time.Now().Format(time.RFC3339),
					Message:         "Internal server error",
					RelatedRecordID: "",
				},
			},
			http.StatusInternalServerError,
		)
		return
	}
}

func writeCustomError(w http.ResponseWriter, ce shared.CustomError, status int) {
	ce.Log()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ce)
}
