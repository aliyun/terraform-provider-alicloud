// Shared utilities for split-layout service packages, reused from the root
// alicloud package (alicloud/errors.go, alicloud/common.go,
// alicloud/validators.go) which service packages cannot import, mirroring
// NeedRetry and IncrementalWait in this package. IsExpectedErrors is
// reduced to the error shapes the tea-rpc service packages produce
// (errs.ComplexError wrapping and alibaba-cloud-sdk-go ServerError);
// every other helper is verbatim.
package helper

import (
	"fmt"
	"os"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const skipResourceSchemaValidationWarning = "\n[NOTE] set env variable TF_SKIP_RESOURCE_SCHEMA_VALIDATION to true can skip the error and get a warning"

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

// IsExpectedErrors reports whether err carries one of the expected API error
// codes.
func IsExpectedErrors(err error, expectCodes []string) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*errs.ComplexError); ok {
		return IsExpectedErrors(e.Cause, expectCodes)
	}

	if e, ok := err.(*errors.ServerError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	for _, code := range expectCodes {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}

func ExpandSingletonToList(singleton interface{}) []interface{} {
	vs := make([]interface{}, 0)
	vs = append(vs, singleton)
	return vs
}

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
