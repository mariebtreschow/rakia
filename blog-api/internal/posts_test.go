package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTitle(t *testing.T) {
	cases := []struct {
		title string
		want  error
	}{
		{"Valid Title", nil},
		{"", ErrTitleEmpty},
		{"Too long title exceeding the maximum character limit", ErrTitleInvalid},
		{"!!", ErrTitleInvalidChars},
	}

	for _, tc := range cases {
		got := validateTitle(tc.title)
		assert.Equal(t, tc.want, got)
	}
}

func TestValidateContent(t *testing.T) {

	cases := []struct {
		content string
		want    error
	}{
		{"Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore.", nil},
		{"", ErrContentEmpty},
		{"Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore.", ErrContentInvalid},
		{"&&&&$##$#$#$!$!$!$$!$!$$!$!$!$$!$!$$!$!$$!$!$!$$!$!$!$&&&&$##$#$#$!$!$!$$!$!$$!$!$!$$!$!$$!$!$$!$!$!$$!$!$!$", ErrContentInvalid},
	}

	for _, tc := range cases {
		got := validateContent(tc.content)
		assert.Equal(t, tc.want, got)
	}
}
