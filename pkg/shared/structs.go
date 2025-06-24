package shared

type CustomError struct {
	ErrorType   string    `json:"-"`
	ErrorLogMsg string    `json:"-"`
	Data        ErrorData `json:"Error"`
}

type ErrorData struct {
	ID              string `json:"ErrorId"`
	TimeStamp       string `json:"ErrorTimeStamp"`
	Message         string `json:"ErrorMessage"`
	RelatedRecordID string `json:"RelatedRecordId"`
}

var EnumErrorTypes = struct {
	Auth,
	BadRequest,
	Communication,
	Internal,
	NotFound string
}{
	"Authorization",
	"BadRequest",
	"Communication",
	"Internal",
	"NotFound",
}
