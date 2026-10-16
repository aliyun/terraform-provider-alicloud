package helper

import (
	"fmt"
	"strings"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

// Tests for the invoker moved from the root package
// (alicloud/common.go); test names are unchanged from
// alicloud/common_test.go.

func TestUnitCommonNewInvoker(t *testing.T) {
	// Test creating new invoker
	invoker := NewInvoker()

	if invoker.catchers == nil {
		t.Error("Expected non-nil catchers")
	}

	// Verify default catchers were added
	if len(invoker.catchers) != 3 {
		t.Errorf("Expected 3 catchers, got %d", len(invoker.catchers))
	}

	// Verify catcher reasons
	expectedReasons := map[string]bool{
		errs.AliyunGoClientFailure: false,
		"ServiceUnavailable":       false,
		errs.Throttling:            false,
	}

	for _, catcher := range invoker.catchers {
		if _, ok := expectedReasons[catcher.Reason]; ok {
			expectedReasons[catcher.Reason] = true
		}
	}

	for reason, found := range expectedReasons {
		if !found {
			t.Errorf("Expected catcher with reason %s not found", reason)
		}
	}
}

func TestUnitCommonInvokerAddCatcher(t *testing.T) {
	invoker := Invoker{}

	customCatcher := Catcher{
		Reason:           "CustomError",
		RetryCount:       5,
		RetryWaitSeconds: 1,
	}

	invoker.AddCatcher(customCatcher)

	if len(invoker.catchers) != 1 {
		t.Errorf("Expected 1 catcher, got %d", len(invoker.catchers))
	}

	if invoker.catchers[0].Reason != "CustomError" {
		t.Errorf("Expected CustomError, got %s", invoker.catchers[0].Reason)
	}
}

func TestUnitCommonInvokerRun(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() (*Invoker, func() error)
		expectError bool
		errorMsg    string
	}{
		{
			name: "success on first try",
			setupFunc: func() (*Invoker, func() error) {
				invoker := NewInvoker()
				fn := func() error { return nil }
				return &invoker, fn
			},
			expectError: false,
		},
		{
			name: "non-retryable error",
			setupFunc: func() (*Invoker, func() error) {
				invoker := NewInvoker()
				fn := func() error { return fmt.Errorf("NonRetryableError") }
				return &invoker, fn
			},
			expectError: true,
			errorMsg:    "NonRetryableError",
		},
		{
			name: "retryable error - throttling",
			setupFunc: func() (*Invoker, func() error) {
				invoker := Invoker{}
				invoker.AddCatcher(Catcher{
					Reason:           errs.Throttling,
					RetryCount:       2,
					RetryWaitSeconds: 0,
				})
				callCount := 0
				fn := func() error {
					callCount++
					if callCount == 1 {
						// Create a tea.SDKError
						code := errs.Throttling
						return &tea.SDKError{
							Code: &code,
						}
					}
					return nil
				}
				return &invoker, fn
			},
			expectError: false,
		},
		{
			name: "retry timeout",
			setupFunc: func() (*Invoker, func() error) {
				invoker := Invoker{}
				invoker.AddCatcher(Catcher{
					Reason:           "TestError",
					RetryCount:       1,
					RetryWaitSeconds: 0,
				})
				fn := func() error {
					code := "TestError"
					return &tea.SDKError{
						Code: &code,
					}
				}
				return &invoker, fn
			},
			expectError: true,
			errorMsg:    "Retry timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoker, fn := tt.setupFunc()
			err := invoker.Run(fn)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expectError && err != nil && tt.errorMsg != "" {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %s, got %v", tt.errorMsg, err)
				}
			}
		})
	}
}
