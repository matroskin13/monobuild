package errors

import "fmt"

type RichError struct {
	Text      string
	CodeError error
}

func NewRichError(text string, codeError error) *RichError {
	return &RichError{Text: text, CodeError: codeError}
}

func (r RichError) Error() string {
	res := r.Text

	if r.CodeError != nil {
		return fmt.Sprintf("%s. \r\n \r\nSystem error: %s", res, r.CodeError)
	}

	return res
}
