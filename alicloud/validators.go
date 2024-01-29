package alicloud

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"gopkg.in/yaml.v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
	}

	return
}

func validateVpnCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	cidrs := strings.Split(value, ",")
	for _, cidr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid CIDR, got error parsing: %s", k, err))
			return
		}

		if ipnet == nil || cidr != ipnet.String() {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid network CIDR, expected %q, got %q",
				k, ipnet, cidr))
			return
		}
	}

	return
}

func validateSwitchCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
		return
	}

	mark, _ := strconv.Atoi(strings.Split(ipnet.String(), "/")[1])
	if mark < 16 || mark > 29 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a network CIDR which mark between 16 and 29",
			k))
	}

	return
}

func validateAllowedSplitStringValue(ss []string, splitStr string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		existed := false
		tsList := strings.Split(value, splitStr)

		for _, ts := range tsList {
			existed = false
			for _, s := range ss {
				if ts == s {
					existed = true
					break
				}
			}
		}
		if !existed {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid string value should in %#v, got %q",
				k, ss, value))
		}
		return

	}
}

func validateStringConvertInt64() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		if value, ok := v.(string); ok {
			_, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				errors = append(errors, fmt.Errorf(
					"%q should be convert to int64, got %q", k, value))
			}
		} else {
			errors = append(errors, fmt.Errorf(
				"%q should be convert to string, got %q", k, value))
		}

		return
	}
}

func validateOssBucketDateTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as date YYYY-MM-DD Format", value))
	}
	return
}

func validateOnsGroupId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !(strings.HasPrefix(value, "GID-") || strings.HasPrefix(value, "GID_")) {
		errors = append(errors, fmt.Errorf("%q is invalid, it must start with 'GID-' or 'GID_'", k))
	}
	if reg := regexp.MustCompile(`^[\w\-]{7,64}$`); !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("%q length is limited to 7-64 and only characters such as letters, digits, '_' and '-' are allowed", k))
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

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func normalizeYamlString(yamlString interface{}) (string, error) {
	var j interface{}

	if yamlString == nil || yamlString.(string) == "" {
		return "", nil
	}

	s := yamlString.(string)

	err := yaml.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	// The error is intentionally ignored here to allow empty policies to passthrough validation.
	// This covers any interpolated values
	bytes, _ := yaml.Marshal(j)

	return string(bytes[:]), nil
}

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func normalizeJsonString(jsonString interface{}) (string, error) {
	var j interface{}

	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}

	s := jsonString.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	// The error is intentionally ignored here to allow empty policies to passthrough validation.
	// This covers any interpolated values
	bytes, _ := json.Marshal(j)

	return string(bytes[:]), nil
}

func validateYamlString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := normalizeYamlString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid YAML: %s", k, err))
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

// Validate length(2~128) and prefix of the name.
func validateNormalName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%s cannot be longer than 128 characters", k))
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}
	return
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

var resourceSchemaValidationSkipped bool

func skipResourceSchemaValidation() bool {
	if resourceSchemaValidationSkipped {
		return resourceSchemaValidationSkipped
	}
	if os.Getenv("TF_SKIP_RESOURCE_SCHEMA_VALIDATION") == "true" {
		resourceSchemaValidationSkipped = true
	}
	return resourceSchemaValidationSkipped
}

const skipResourceSchemaValidationWarning = "\n[NOTE] set env variable TF_SKIP_RESOURCE_SCHEMA_VALIDATION to true can skip the error and get a warning"

// IntBetween returns a SchemaValidateFunc which tests if the provided value
// is of type int and is between min and max (inclusive)
func IntBetween(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		if v < min || v > max {
			if skipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			} else {
				es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d %s", k, min, max, v, skipResourceSchemaValidationWarning))
			}
			return
		}

		return
	}
}

// IntAtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type int and is at least min (inclusive)
func IntAtLeast(min int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		if v < min {
			if skipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected %s to be at least (%d), got %d", k, min, v))
			} else {
				es = append(es, fmt.Errorf("expected %s to be at least (%d), got %d %s", k, min, v, skipResourceSchemaValidationWarning))
			}
			return
		}

		return
	}
}

// IntAtMost returns a SchemaValidateFunc which tests if the provided value
// is of type int and is at most max (inclusive)
func IntAtMost(max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		if v > max {
			if skipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected %s to be at most (%d), got %d", k, max, v))
			} else {
				es = append(es, fmt.Errorf("expected %s to be at most (%d), got %d %s", k, max, v, skipResourceSchemaValidationWarning))
			}
			return
		}

		return
	}
}

// IntInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type int and matches the value of an element in the valid slice
func IntInSlice(valid []int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be an integer", k))
			return
		}

		for _, validInt := range valid {
			if v == validInt {
				return
			}
		}

		if skipResourceSchemaValidation() {
			s = append(s, fmt.Sprintf("expected %s to be one of %v, got %d", k, valid, v))
		} else {
			es = append(es, fmt.Errorf("expected %s to be one of %v, got %d %s", k, valid, v, skipResourceSchemaValidationWarning))
		}
		return
	}
}

// StringInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and matches the value of an element in the valid slice
// will test with in lower case if ignoreCase is true
func StringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		for _, str := range valid {
			if v == str || (ignoreCase && strings.ToLower(v) == strings.ToLower(str)) {
				return
			}
		}

		if skipResourceSchemaValidation() {
			s = append(s, fmt.Sprintf("expected %s to be one of %v, got %s", k, valid, v))
		} else {
			es = append(es, fmt.Errorf("expected %s to be one of %v, got %s %s", k, valid, v, skipResourceSchemaValidationWarning))
		}
		return
	}
}

// StringLenBetween returns a SchemaValidateFunc which tests if the provided value
// is of type string and has length between min and max (inclusive)
func StringLenBetween(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		if len(v) < min || len(v) > max {
			if skipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected length of %s to be in the range (%d - %d), got %s", k, min, max, v))
			} else {
				es = append(es, fmt.Errorf("expected length of %s to be in the range (%d - %d), got %s %s", k, min, max, v, skipResourceSchemaValidationWarning))
			}
		}
		return
	}
}

// StringMatch returns a SchemaValidateFunc which tests if the provided value
// matches a given regexp. Optionally an error message can be provided to
// return something friendlier than "must match some globby regexp".
func StringMatch(r *regexp.Regexp, message string) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
		}

		if ok := r.MatchString(v); !ok {
			if message != "" {
				if skipResourceSchemaValidation() {
					s = append(s, fmt.Sprintf("invalid value for %s (%s)", k, message))
				} else {
					es = append(es, fmt.Errorf("invalid value for %s (%s) %s", k, message, skipResourceSchemaValidationWarning))
				}
				return

			}
			if skipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected value of %s to match regular expression %q", k, r))
			} else {
				es = append(es, fmt.Errorf("expected value of %s to match regular expression %q %s", k, r, skipResourceSchemaValidationWarning))
			}
			return
		}
		return nil, nil
	}
}

// StringDoesNotMatch returns a SchemaValidateFunc which tests if the provided value
// does not match a given regexp. Optionally an error message can be provided to
// return something friendlier than "must not match some globby regexp".
func StringDoesNotMatch(r *regexp.Regexp, message string) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
		}

		if ok := r.MatchString(v); ok {
			if message != "" {
				if skipResourceSchemaValidation() {
					s = append(s, fmt.Sprintf("invalid value for %s (%s)", k, message))
				} else {
					es = append(es, fmt.Errorf("invalid value for %s (%s) %s", k, message, skipResourceSchemaValidationWarning))
				}
				return

			}
			if skipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected value of %s to not match regular expression %q", k, r))
			} else {
				es = append(es, fmt.Errorf("expected value of %s to not match regular expression %q %s", k, r, skipResourceSchemaValidationWarning))
			}
			return
		}
		return
	}
}

// FloatBetween returns a SchemaValidateFunc which tests if the provided value
// is of type float64 and is between min and max (inclusive).
func FloatBetween(min, max float64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(float64)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be float64", k))
			return
		}

		if v < min || v > max {
			if skipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected %s to be in the range (%f - %f), got %f", k, min, max, v))
			} else {
				es = append(es, fmt.Errorf("expected %s to be in the range (%f - %f), got %f %s", k, min, max, v, skipResourceSchemaValidationWarning))
			}
			return
		}

		return
	}
}

// FloatAtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type float and is at least min (inclusive)
func FloatAtLeast(min float64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(float64)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be float", k))
			return
		}

		if v < min {
			if skipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected %s to be at least (%f), got %f", k, min, v))
			} else {
				es = append(es, fmt.Errorf("expected %s to be at least (%f), got %f %s", k, min, v, skipResourceSchemaValidationWarning))
			}
			return
		}

		return
	}
}

// FloatAtMost returns a SchemaValidateFunc which tests if the provided value
// is of type float and is at most max (inclusive)
func FloatAtMost(max float64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(float64)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be float", k))
			return
		}

		if v > max {
			if skipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected %s to be at most (%f), got %f", k, max, v))
			} else {
				es = append(es, fmt.Errorf("expected %s to be at most (%f), got %f %s", k, max, v, skipResourceSchemaValidationWarning))
			}
			return
		}

		return
	}
}

// StringDoesNotContainAny returns a SchemaValidateFunc which validates that the
// provided value does not contain any of the specified Unicode code points in chars.
func StringDoesNotContainAny(chars string) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if strings.ContainsAny(v, chars) {
			if skipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected value of %s to not contain any of %q", k, chars))
			} else {
				es = append(es, fmt.Errorf("expected value of %s to not contain any of %q %s", k, chars, skipResourceSchemaValidationWarning))
			}
			return
		}

		return
	}
}

// ValidateRFC3339TimeString is a ValidateFunc that ensures a string parses
// as time.RFC3339 format
func ValidateRFC3339TimeString(allowEmpty bool) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		if v == "" && allowEmpty {
			return
		}
		if _, err := time.Parse(time.RFC3339, v); err != nil {
			if skipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("%q: invalid RFC3339 timestamp", k))
			} else {
				es = append(es, fmt.Errorf("%q: invalid RFC3339 timestamp %s", k, skipResourceSchemaValidationWarning))
			}
		}
		return
	}
}

// StringLenAtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type string and has length at least min (inclusive)
func StringLenAtLeast(min int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {

		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		valueLen := len(strings.TrimSpace(v))
		if valueLen < min {
			if skipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected length of %s to be at least (%d), got (%d)", k, min, valueLen))
			} else {
				es = append(es, fmt.Errorf("expected length of %s to be at least (%d), got (%d) %s", k, min, valueLen, skipResourceSchemaValidationWarning))
			}
		}
		return
	}
}

func validateRedisConfig(v interface{}, k string) (ws []string, errors []error) {
	value, _ := v.(map[string]interface{})

	if len(value) < 1 {
		errors = append(errors, fmt.Errorf("invalid value for %s (%s)", k, value))
	}

	return
}
