package arn

import (
	"errors"
	"strings"
	"testing"
)

// valid ARNs, each paired with the sections Parse must produce. Every entry is also
// a round-trip case: String has to reproduce the input exactly, which is what lets
// arn_parse and arn_build compose.
var validARNs = []struct {
	arn       string
	service   string
	region    string
	accountID string
	resource  string
}{
	// RAM is global, so the region section is empty rather than absent.
	{"acs:ram::123456789012****:role/MyRole", "ram", "", "123456789012****", "role/MyRole"},
	{"acs:ecs:cn-hangzhou:123456789012****:instance/i-bp1234567890abcdef", "ecs", "cn-hangzhou", "123456789012****", "instance/i-bp1234567890abcdef"},
	{"acs:oss:*:*:my-bucket/*", "oss", "*", "*", "my-bucket/*"},
	{"acs:fc:cn-hangzhou:123456789012****:services/foo.LATEST/functions/bar", "fc", "cn-hangzhou", "123456789012****", "services/foo.LATEST/functions/bar"},
	{"acs:mns:cn-hangzhou:123456789012****:/queues/my-queue/messages", "mns", "cn-hangzhou", "123456789012****", "/queues/my-queue/messages"},
	{"acs:kms:*:*:*", "kms", "*", "*", "*"},
	// A wildcard policy names no account, so both the region and the account section
	// are empty or "*". This exact shape appears as an example in the SDK's
	// PrivateLink models.
	{"acs:ram:*::*", "ram", "*", "", "*"},
	// The resource is the last section, so colons inside it are part of it.
	// strings.Split would have shredded this one.
	{"acs:log:cn-hangzhou:123456789012****:project/my-project:logstore/my-logstore", "log", "cn-hangzhou", "123456789012****", "project/my-project:logstore/my-logstore"},
	// Parse validates the shape, not the vocabulary: an unknown or empty service is
	// still a syntactically valid ARN.
	{"acs::cn-hangzhou:123456789012****:instance/i-bp1234567890abcdef", "", "cn-hangzhou", "123456789012****", "instance/i-bp1234567890abcdef"},
}

func TestARNParse(t *testing.T) {
	for _, c := range validARNs {
		got, err := Parse(c.arn)
		if err != nil {
			t.Errorf("Parse(%q) unexpected error: %s", c.arn, err)
			continue
		}

		want := ARN{Service: c.service, Region: c.region, AccountID: c.accountID, Resource: c.resource}
		if got != want {
			t.Errorf("Parse(%q) = %+v, want %+v", c.arn, got, want)
		}
	}
}

func TestARNParseErrors(t *testing.T) {
	cases := []struct {
		arn  string
		want error
	}{
		{"", ErrSections},
		{"acs", ErrSections},
		{"acs:ecs", ErrSections},
		{"acs:ecs:cn-hangzhou:123456789012****", ErrSections},
		// Five sections, wrong first one. An AWS ARN is the likeliest thing a user
		// pastes by mistake, so it has to fail on the prefix rather than parse into
		// nonsense.
		{"arn:aws:s3:::my-bucket", ErrPrefix},
		// The prefix is compared verbatim, so case matters.
		{"ACS:ecs:cn-hangzhou:123456789012****:instance/i-bp1234567890abcdef", ErrPrefix},
		{":ecs:cn-hangzhou:123456789012****:instance/i-bp1234567890abcdef", ErrPrefix},
		// An ARN that names no resource is a truncation, not a wildcard: "*" is how a
		// policy matches everything.
		{"acs:ecs:cn-hangzhou:123456789012****:", ErrResource},
		{"acs::::", ErrResource},
	}

	for _, c := range cases {
		got, err := Parse(c.arn)
		if err == nil {
			t.Errorf("Parse(%q) = %+v, want error %s", c.arn, got, c.want)
			continue
		}

		if !errors.Is(err, c.want) {
			t.Errorf("Parse(%q) error = %q, want it to wrap %q", c.arn, err, c.want)
		}

		// A rejected ARN must not leak half-parsed sections to a caller that ignores
		// the error.
		if got != (ARN{}) {
			t.Errorf("Parse(%q) returned %+v alongside its error, want the zero ARN", c.arn, got)
		}

		// The message has to name the offending value, or a plan with several ARNs in
		// it gives no clue which one is wrong.
		if c.arn != "" && !strings.Contains(err.Error(), c.arn) {
			t.Errorf("Parse(%q) error = %q, want it to quote the input", c.arn, err)
		}
	}
}

func TestARNString(t *testing.T) {
	for _, c := range validARNs {
		a := ARN{Service: c.service, Region: c.region, AccountID: c.accountID, Resource: c.resource}
		if got := a.String(); got != c.arn {
			t.Errorf("ARN%+v.String() = %q, want %q", a, got, c.arn)
		}
	}
}

// TestARNStringDoesNotValidate pins the documented asymmetry between the two
// directions: String is a formatter, exactly like the arn_build function, so it
// emits a string Parse then refuses. Tightening String is a breaking change and has
// to break this case first.
func TestARNStringDoesNotValidate(t *testing.T) {
	zero := ARN{}.String()
	if zero != "acs::::" {
		t.Errorf("ARN{}.String() = %q, want %q", zero, "acs::::")
	}

	if _, err := Parse(zero); !errors.Is(err, ErrResource) {
		t.Errorf("Parse(ARN{}.String()) error = %v, want it to wrap %q", err, ErrResource)
	}
}
