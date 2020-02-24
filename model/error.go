package model

import "fmt"

const (
	MarshalErrorType                     = "Marshal"
	UnMarshalErrorType                   = "UnMarshal"
	ConvertErrorType                     = "Convert"
	PutStateErrorType                    = "PutState"
	GetStateErrorType                    = "GetState"
	SetEventErrorType                    = "SetEvent"
	CreateCompositeKeyErrorType          = "CreateCompositeKey"
	GetStatePartialCompositeKeyErrorType = "GetStatePartialCompositeKey"
	SpliteCompositeKeyErrorType          = "SpliteCompositeKey"
)

type CustomError struct {
	ErrorType string
	TypeName  string
	Message   string
}

func NewCustomError(errorType, typeName, message string) *CustomError {
	return &CustomError{
		ErrorType: errorType,
		TypeName:  typeName,
		Message:   message,
	}
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("failed to %s %s, error: %s", e.ErrorType, e.TypeName, e.Message)
}
