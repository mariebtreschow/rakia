package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTitle(t *testing.T) {
	cases := []struct {
		title string
		want  error
		test  string
	}{
		{"Title 2", nil, "valid title"},
		{"", ErrTitleEmpty, "empty title"},
		{"Too long title exceeding the maximum character limit is too long", ErrTitleInvalid, "too long title"},
		{"!!&&&&& Title", ErrTitleInvalidChars, "title with invalid characters"},
		{"Title with       excessive whitespace", ErrTitleFormat, "title with excessive whitespace"},
		{"Title  with multiple consecutive spaces", ErrTitleFormat, "title with multiple consecutive spaces"},
		{"Buy now at big discount", ErrTitleSpammy, "title with spammy phrase"},
		{"title with lowercase first letter", ErrTitleCapitalization, "title with lowercase first letter"},
		{"TITLE WITH ALL CAPS", ErrTitleCapitalization, "title with all caps"},
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
		test    string
	}{
		{"Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore.", nil, "valid content"},
		{"", ErrContentEmpty, "empty content"},
		{"Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore. Amet quiquia sed ut velit eius. Etincidunt non consectetur porro velit neque. Quiquia est dolorem dolore quiquia dolore eius quisquam. Dolor tempora dolor magnam dolor sed quiquia consectetur. Quiquia quaerat numquam consectetur neque. Dolor amet modi modi. Voluptatem adipisci etincidunt quiquia dolor etincidunt. Est velit etincidunt ipsum dolor. Sit etincidunt neque quaerat voluptatem dolorem dolor dolore.", ErrContentInvalid, "content too long"},
		{"&&&&$##$#$#$!$!$!$$!$!$$!$!$!$$!$!$$!$!$$!$!$!$$!$!$!$&&&&$##$#$#$!$!$!$$!$!$$!$!$!$$!$!$$!$!$$!$!$!$$!$!$!$", ErrContentInvalid, "content with too many special characters"},
		{"Aaaaa detta e sa bra eller hur det e liksom fantasikst hur man kan skriva sa har langt for att de maste vara en valid content", ErrContentConsecutiveChar, "content with excessive consecutive identical characters"},
	}

	for _, tc := range cases {
		got := validateContent(tc.content)
		assert.Equal(t, tc.want, got)
	}
}

func TestValidateAuthor(t *testing.T) {

	cases := []struct {
		author string
		want   error
		test   string
	}{
		{"Author 1", nil, "valid author"},
		{"", ErrAuthorEmpty, "empty author"},
		{"M", ErrAuthorNameInvalid, "author too short"},
	}

	for _, tc := range cases {
		got := validateAuthor(tc.author)
		assert.Equal(t, tc.want, got)
	}

}
