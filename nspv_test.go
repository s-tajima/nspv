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

	mockURL := fmt.Sprintf("=~^%s/range/.*$", HibpApiBaseURL)
	mockResponse := "96E32BEF9763E32668B444F64A6E591B462:10\r\n1C4ECB46CCD8BCB1FD0E9F07D859F554D76:0"

	httpmock.RegisterResponder(
		"GET",
		mockURL,
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
		{"pwnedpwned", ViolateHibpCheck},
		{"pwnedqwn3d", Ok},
	}

	v := NewValidator()
	v.SetDict([]string{"password"})

	for _, tt := range cases {
		result, err := v.Validate(tt.password)
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

	v := NewValidator()
	v.SetMinLength(3)
	v.SetMaxLength(7)

	for _, tt := range cases {
		result, err := v.Validate(tt.password)
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

	v := NewValidator()
	v.SetHibpClientContext(ctx)

	for _, tt := range cases {
		_, err := v.Validate(tt.password)
		if err == nil || !errors.Is(err, tt.err) {
			t.Errorf("for %s ... something wrong with context. (%s)", tt.password, err)
		}
	}
}

func TestValidatorHibpError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockURL := fmt.Sprintf("=~^%s/range/.*$", HibpApiBaseURL)
	httpmock.RegisterResponder(
		"GET",
		mockURL,
		httpmock.NewStringResponder(500, ""),
	)

	v := NewValidator()
	result, _ := v.Validate("password")
	if result != Error {
		t.Errorf("for password got:%s, want: %s", result, Error)
	}

	v.SetIgnoreHibpError(true)
	result, _ = v.Validate("password")
	if result != Ok {
		t.Errorf("for password got:%s, want: %s", result, Ok)
	}
}
