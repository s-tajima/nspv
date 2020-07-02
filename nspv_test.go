package nspv

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func activateMock() {
	httpmock.Activate()

	mockUrl := fmt.Sprintf("=~^%s/range/.*$", HIBP_API_BASE_URL)
	mockResponse := "96E32BEF9763E32668B444F64A6E591B462:10\r\n1C4ECB46CCD8BCB1FD0E9F07D859F554D76:0"

	httpmock.RegisterResponder(
		"GET",
		mockUrl,
		httpmock.NewStringResponder(200, mockResponse),
	)
}

func deactivateMock() {
	httpmock.DeactivateAndReset()
}

func TestValidatorDefault(t *testing.T) {
	activateMock()
	defer deactivateMock()

	cases := []struct {
		password string
		result   Result
	}{
		{"okpassword", Ok},
		{strings.Repeat("a", 7), ViolateMinLengthCheck},
		{strings.Repeat("a", 8), Ok},
		{strings.Repeat("a", 64), Ok},
		{strings.Repeat("a", 65), ViolateMaxLengthCheck},
		{"password", ViolateDictCheck},
		{"p@ssword", ViolateDictCheck},
		{"p@s3w0rd", Ok},
		{"pwnedpwned", ViolateBibpCheck},
		{"pwnedqwn3d", Ok},
	}

	c := NewValidator()
	c.SetDict([]string{"password"})

	for _, tt := range cases {
		result, err := c.Validate(tt.password)
		if err != nil {
			t.Error(err)
		}

		if result != tt.result {
			t.Errorf("for %s ... got: %d, want: %d", tt.password, result, tt.result)
		}
	}
}

func TestValidatorFixedLength(t *testing.T) {
	activateMock()
	defer deactivateMock()

	cases := []struct {
		password string
		result   Result
	}{
		{strings.Repeat("a", 2), ViolateMinLengthCheck},
		{strings.Repeat("a", 3), Ok},
		{strings.Repeat("a", 4), Ok},
		{strings.Repeat("a", 6), Ok},
		{strings.Repeat("a", 7), Ok},
		{strings.Repeat("a", 8), ViolateMaxLengthCheck},
	}

	c := NewValidator()
	c.SetMinLength(3)
	c.SetMaxLength(7)

	for _, tt := range cases {
		result, err := c.Validate(tt.password)
		if err != nil {
			t.Error(err)
		}

		if result != tt.result {
			t.Errorf("for %s ... got: %d, want: %d", tt.password, result, tt.result)
		}
	}
}

func TestValidatorWithContextTimeout(t *testing.T) {
	cases := []struct {
		password string
		err      error
	}{
		{"weakpass", context.DeadlineExceeded},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	c := NewValidator()
	c.SetHibpClientContext(ctx)

	for _, tt := range cases {
		_, err := c.Validate(tt.password)
		if err == nil || !errors.Is(err, tt.err) {
			t.Errorf("for %s ... something wrong with context. (%s)", tt.password, err)
		}
	}
}

func TestValidatorHibpError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockUrl := fmt.Sprintf("=~^%s/range/.*$", HIBP_API_BASE_URL)
	httpmock.RegisterResponder(
		"GET",
		mockUrl,
		httpmock.NewStringResponder(500, ""),
	)

	c := NewValidator()
	result, _ := c.Validate("password")
	if result != Error {
		t.Errorf("for password got:%s, want: %s", result, Error)
	}

	c.SetIgnoreHibpError(true)
	result, _ = c.Validate("password")
	if result != Ok {
		t.Errorf("for password got:%s, want: %s", result, Ok)
	}
}
