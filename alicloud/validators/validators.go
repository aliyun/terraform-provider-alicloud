package validators

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
)

// Generic schema validators moved from the root package
// (alicloud/validators.go); package alicloud re-exports them as var
// aliases so the flat layout keeps compiling unchanged, and service
// packages import them here directly. Product-specific validators stay
// in the root package and migrate with their products.

// ValidateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func ValidateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
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

func ValidateVpnCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
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

func ValidateSwitchCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
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

func ValidateAllowedSplitStringValue(ss []string, splitStr string) schema.SchemaValidateFunc {
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

func ValidateStringConvertInt64() schema.SchemaValidateFunc {
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

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func NormalizeYamlString(yamlString interface{}) (string, error) {
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
func NormalizeJsonString(jsonString interface{}) (string, error) {
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

func ValidateYamlString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := NormalizeYamlString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid YAML: %s", k, err))
	}

	return
}

// Validate length(2~128) and prefix of the name.
func ValidateNormalName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%s cannot be longer than 128 characters", k))
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}
	return
}

var resourceSchemaValidationSkipped bool

func SkipResourceSchemaValidation() bool {
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
			if SkipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			} else {
				es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d %s", k, min, max, v, skipResourceSchemaValidationWarning))
			}
			return
		}

		return
	}
}

// ValidateNullableIntBetween returns a SchemaValidateFunc for a nullable int
// field declared as TypeString: it accepts "" (meaning "not set") or a
// canonical decimal integer string (e.g. "10") that parses to an int64
// between min and max (inclusive). Non-canonical forms like "+10", "010"
// or "-0" are always rejected because they would be normalized by Read and
// cause a permanent diff. A range violation is downgraded to a warning when
// SkipResourceSchemaValidation() is true; format errors are never downgraded.
func ValidateNullableIntBetween(min, max int64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		value, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		if value == "" {
			return
		}
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			es = append(es, fmt.Errorf("expected %s to be empty or an integer string, got %s", k, value))
			return
		}
		if strconv.FormatInt(v, 10) != value {
			es = append(es, fmt.Errorf("expected %s to be a canonical decimal integer string (e.g. \"10\"), got %s", k, value))
			return
		}
		if v < min || v > max {
			if SkipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			} else {
				es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d %s", k, min, max, v, skipResourceSchemaValidationWarning))
			}
			return
		}
		return
	}
}

// ValidateNullableBool is a SchemaValidateFunc for a nullable bool field
// declared as TypeString: it accepts "" (meaning "not set"), or the lowercase
// "true" / "false" (case sensitive). An invalid value is downgraded to a
// warning when SkipResourceSchemaValidation() is true; a non-string type is
// always a hard error.
func ValidateNullableBool(i interface{}, k string) (s []string, es []error) {
	value, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	if value != "" && value != "true" && value != "false" {
		if SkipResourceSchemaValidation() {
			s = append(s, fmt.Sprintf("expected %s to be an empty string or lowercase \"true\"/\"false\", got %s", k, value))
		} else {
			es = append(es, fmt.Errorf("expected %s to be an empty string or lowercase \"true\"/\"false\", got %s %s", k, value, skipResourceSchemaValidationWarning))
		}
	}
	return
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
			if SkipResourceSchemaValidation() {
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
			if SkipResourceSchemaValidation() {
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

		if SkipResourceSchemaValidation() {
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

		if SkipResourceSchemaValidation() {
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
			if SkipResourceSchemaValidation() {
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
				if SkipResourceSchemaValidation() {
					s = append(s, fmt.Sprintf("invalid value for %s (%s)", k, message))
				} else {
					es = append(es, fmt.Errorf("invalid value for %s (%s) %s", k, message, skipResourceSchemaValidationWarning))
				}
				return

			}
			if SkipResourceSchemaValidation() {
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
				if SkipResourceSchemaValidation() {
					s = append(s, fmt.Sprintf("invalid value for %s (%s)", k, message))
				} else {
					es = append(es, fmt.Errorf("invalid value for %s (%s) %s", k, message, skipResourceSchemaValidationWarning))
				}
				return

			}
			if SkipResourceSchemaValidation() {
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
			if SkipResourceSchemaValidation() {
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
			if SkipResourceSchemaValidation() {
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
			if SkipResourceSchemaValidation() {
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
			if SkipResourceSchemaValidation() {
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
			if SkipResourceSchemaValidation() {
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
			if SkipResourceSchemaValidation() {
				s = append(s, fmt.Sprintf("expected length of %s to be at least (%d), got (%d)", k, min, valueLen))
			} else {
				es = append(es, fmt.Errorf("expected length of %s to be at least (%d), got (%d) %s", k, min, valueLen, skipResourceSchemaValidationWarning))
			}
		}
		return
	}
}
