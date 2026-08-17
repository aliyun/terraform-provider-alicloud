package arn

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// Prefix is the first section of every Alibaba Cloud ARN. AWS varies the
	// equivalent section by partition (aws, aws-cn, aws-us-gov), which is why its
	// ARN type carries a Partition field; Alibaba Cloud has exactly one, so here it
	// is a constant and not a field of ARN.
	Prefix = "acs"

	// Delimiter separates the sections of an ARN.
	Delimiter = ":"
)

// Section indices within an ARN. sections closes the block, so it holds the number
// of sections and stays correct if one is ever inserted above it.
const (
	sectionPrefix = iota
	sectionService
	sectionRegion
	sectionAccountID
	sectionResource
	sections
)

// Reasons Parse rejects a string. Parse wraps them with the offending value, so a
// caller that needs to tell them apart uses errors.Is rather than matching text.
var (
	ErrSections = fmt.Errorf("want %s%s<service>%s<region>%s<account_id>%s<resource>", Prefix, Delimiter, Delimiter, Delimiter, Delimiter)
	ErrPrefix   = fmt.Errorf("first section must be %q", Prefix)
	ErrResource = errors.New("resource is empty")
)

// ARN is a parsed Alibaba Cloud Resource Name. The zero value is not a valid ARN.
type ARN struct {
	// Service is the RAM code of the service that owns the resource, such as "ecs",
	// "oss" or "ram".
	Service string

	// Region is the region of the resource, such as "cn-hangzhou". It is empty for
	// services whose resources are not regional, such as RAM, and "*" in a policy
	// that matches every region.
	Region string

	// AccountID is the ID of the account that owns the resource. It is empty or "*"
	// in a policy that matches every account.
	AccountID string

	// Resource is the service-specific resource path, such as
	// "instance/i-bp1234567890abcdef". It is the last section, so it may itself
	// contain Delimiter.
	Resource string
}

// Parse splits s into its sections.
func Parse(s string) (ARN, error) {
	// The resource is the last section and may contain Delimiter, so the split has to
	// stop counting there — strings.Split would shred "project/p:logstore/l".
	parts := strings.SplitN(s, Delimiter, sections)

	if len(parts) != sections {
		return ARN{}, fmt.Errorf("%q is not an ARN: %w", s, ErrSections)
	}
	if parts[sectionPrefix] != Prefix {
		return ARN{}, fmt.Errorf("%q is not an ARN: %w", s, ErrPrefix)
	}
	if parts[sectionResource] == "" {
		return ARN{}, fmt.Errorf("%q is not an ARN: %w", s, ErrResource)
	}

	return ARN{
		Service:   parts[sectionService],
		Region:    parts[sectionRegion],
		AccountID: parts[sectionAccountID],
		Resource:  parts[sectionResource],
	}, nil
}

// String joins the sections back into an ARN. It is the inverse of Parse for every
// ARN that Parse returned, but it validates nothing itself, so a hand-built ARN
// with a colon in any section but Resource formats into a string Parse will not
// read back the same way.
func (a ARN) String() string {
	return strings.Join([]string{Prefix, a.Service, a.Region, a.AccountID, a.Resource}, Delimiter)
}
