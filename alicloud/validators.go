package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
)

// The generic validators below moved to the leaf package validators; these
// aliases keep the flat layout compiling unchanged. Product-specific
// validators (OTS, Ons, DNS RR, OSS date, DB port, SSL VPN port, Redis
// config) stay in this file and migrate with their products.
var (
	validateCIDRNetworkAddress       = validators.ValidateCIDRNetworkAddress
	validateVpnCIDRNetworkAddress    = validators.ValidateVpnCIDRNetworkAddress
	validateSwitchCIDRNetworkAddress = validators.ValidateSwitchCIDRNetworkAddress
	validateAllowedSplitStringValue  = validators.ValidateAllowedSplitStringValue
	validateStringConvertInt64       = validators.ValidateStringConvertInt64
	normalizeYamlString              = validators.NormalizeYamlString
	normalizeJsonString              = validators.NormalizeJsonString
	validateYamlString               = validators.ValidateYamlString
	validateNormalName               = validators.ValidateNormalName
	skipResourceSchemaValidation     = validators.SkipResourceSchemaValidation

	IntBetween                 = validators.IntBetween
	ValidateNullableIntBetween = validators.ValidateNullableIntBetween
	ValidateNullableBool       = validators.ValidateNullableBool
	IntAtLeast                 = validators.IntAtLeast
	IntAtMost                  = validators.IntAtMost
	IntInSlice                 = validators.IntInSlice
	StringInSlice              = validators.StringInSlice
	StringLenBetween           = validators.StringLenBetween
	StringMatch                = validators.StringMatch
	StringDoesNotMatch         = validators.StringDoesNotMatch
	StringDoesNotContainAny    = validators.StringDoesNotContainAny
	FloatBetween               = validators.FloatBetween
	FloatAtLeast               = validators.FloatAtLeast
	FloatAtMost                = validators.FloatAtMost
	ValidateRFC3339TimeString  = validators.ValidateRFC3339TimeString
	StringLenAtLeast           = validators.StringLenAtLeast

	UpperLowerCaseDiffSuppressFunc = validators.UpperLowerCaseDiffSuppressFunc
)

func validateOssBucketDateTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as date YYYY-MM-DD Format", value))
	}
	return
}

func validateRR(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if strings.HasPrefix(value, "-") || strings.HasSuffix(value, "-") {
		errors = append(errors, fmt.Errorf("RR is invalid, it can not starts or ends with '-'"))
	}

	if len(value) > 253 {
		errors = append(errors, fmt.Errorf("RR can not longer than 253 characters."))
	}

	for _, part := range strings.Split(value, ".") {
		if len(part) > 63 {
			errors = append(errors, fmt.Errorf("Each part of RR split with . can not longer than 63 characters."))
			return
		}
	}
	return
}

func validateDBConnectionPort(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil {
			errors = append(errors, err)
		}
		if port < 1000 || port > 5999 {
			errors = append(errors, fmt.Errorf("%q cannot be less than 3001 and larger than 3999.", k))
		}
	}
	return
}

func validateSslVpnPortValue(is []int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		ws, errors = validation.IntBetween(1, 65535)(v, k)
		if errors != nil {
			return
		}

		value := v.(int)
		for _, i := range is {
			if i == value {
				errors = append(errors, fmt.Errorf(
					"%q must contain a valid int value should not be in array %#v, got %q",
					k, is, value))
				return
			}
		}
		return

	}
}

// below copy/pasta from https://github.com/hashicorp/terraform-plugin-sdk/blob/master/helper/validation/validation.go
// alicloud vendor contains very old version of Terraform which lacks this functions

// IntBetween returns a SchemaValidateFunc which tests if the provided value
// is of type int and is between min and max (inclusive)
func intBetween(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		if v < min || v > max {
			es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			return
		}

		return
	}
}

// The instance name must be composed of a~z, A~Z, 0~9 and a hyphen (-),
// the first character must be a letter and the last character cannot be a hyphen (-),
// the legal length range is 3~16 bytes.
func validateOTSInstanceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	reg := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]{1,14}[a-zA-Z0-9]$")
	if !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("the instance name must consist of a~z, A~Z, 0~9 and a hyphen (-), "+
			"the first character must be a letter and the last character cannot be a hyphen (-), the legal length range is 3~16 bytes"))
	}
	return
}

// The table name must consist of a~z, A~Z, 0~9 and an underscore (_), the first character must be a letter or underscore (_),
// the legal length range is 1~255 bytes.
func validateOTSTableName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	reg := regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]{0,254}$")
	if !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("the table name must consist of a~z, A~Z, 0~9 and an underscore (_), "+
			"the first character must be a letter or underscore (_), the legal length range is 1~255 bytes"))
	}
	return
}

// The tunnel name must consist of a~z, A~Z, 0~9 and an underscore (_), the first character must be a letter or underscore (_),
// the legal length range is 1~255 bytes.
func validateOTSTunnelName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	reg := regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]{0,254}$")
	if !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("the tunnel name must consist of a~z, A~Z, 0~9 and an underscore (_), "+
			"the first character must be a letter or underscore (_), the legal length range is 1~255 bytes"))
	}
	return
}

// The index name must consist of a~z, A~Z, 0~9 and an underscore (_), the first character must be a letter or underscore (_),
// the legal length range is 1~255 bytes.
func validateOTSIndexName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	reg := regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]{0,254}$")
	if !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("the index name must consist of a~z, A~Z, 0~9 and an underscore (_), "+
			"the first character must be a letter or underscore (_), the legal length range is 1~255 bytes"))
	}
	return
}

func validateRedisConfig(v interface{}, k string) (ws []string, errors []error) {
	value, _ := v.(map[string]interface{})

	if len(value) < 1 {
		errors = append(errors, fmt.Errorf("invalid value for %s (%s)", k, value))
	}

	return
}
