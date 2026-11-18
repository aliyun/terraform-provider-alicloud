package client

import (
	"errors"
	"fmt"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	sdkCredentials "github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/stretchr/testify/assert"
)

// mockCredential implements credential.Credential for testing purposes.
// When err is non-nil, GetCredential returns an error so AuthCredential
// does not override the Config's AccessKey/SecretKey/SecurityToken fields.
type mockCredential struct {
	accessKeyId     string
	accessKeySecret string
	securityToken   string
	err             error
}

func (m *mockCredential) GetAccessKeyId() (*string, error) {
	return tea.String(m.accessKeyId), m.err
}

func (m *mockCredential) GetAccessKeySecret() (*string, error) {
	return tea.String(m.accessKeySecret), m.err
}

func (m *mockCredential) GetSecurityToken() (*string, error) {
	return tea.String(m.securityToken), m.err
}

func (m *mockCredential) GetBearerToken() *string {
	return tea.String("")
}

func (m *mockCredential) GetType() *string {
	return tea.String("mock")
}

func (m *mockCredential) GetCredential() (*credentials.CredentialModel, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &credentials.CredentialModel{
		AccessKeyId:     tea.String(m.accessKeyId),
		AccessKeySecret: tea.String(m.accessKeySecret),
		SecurityToken:   tea.String(m.securityToken),
		Type:            tea.String("mock"),
	}, nil
}

func TestUnitCommonRefreshAuthCredential_ECS_Role(t *testing.T) {
	client := NewTestClient(t)
	err := client.config.RefreshAuthCredential()

	assert.NoError(t, err)
	assert.NotEmpty(t, client.config.AccessKey)
	assert.NotEmpty(t, client.config.SecretKey)
}

func TestUnitCommonRefreshAuthCredential_OIDC(t *testing.T) {
	client := NewTestClient(t)
	config := client.config
	err := config.RefreshAuthCredential()

	assert.NoError(t, err)
	assert.NotEmpty(t, config.AccessKey)
	assert.NotEmpty(t, config.SecretKey)
}

func TestUnitCommonRefreshAuthCredential_AssumeRole(t *testing.T) {
	client := NewTestClient(t)
	config := client.config
	err := config.RefreshAuthCredential()

	assert.NoError(t, err)
}

func TestUnitCommonRefreshAuthCredential_Error(t *testing.T) {
	config := &Config{
		AssumeRoleWithOidc: &AssumeRoleWithOidc{
			RoleARN: "invalid-arn",
		},
	}

	err := config.RefreshAuthCredential()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refreshing credential failed")
}

func TestUnitCommonNeedRetry(t *testing.T) {
	testCases := []struct {
		name        string
		err         error
		expected    bool
		description string
	}{
		{
			name:        "Post_HTTPS_Error",
			err:         errors.New("Post https://example.com failed"),
			expected:    true,
			description: "Post HTTPS errors should be retried",
		},
		{
			name:        "Post_HTTPS_With_Quotes_Error",
			err:         errors.New("Post \"https://example.com\" failed"),
			expected:    true,
			description: "Post HTTPS errors with quotes should be retried",
		},
		{
			name: "ServiceUnavailable_Error",
			err: &tea.SDKError{
				Code:    tea.String("ServiceUnavailable"),
				Message: tea.String("Service unavailable"),
			},
			expected:    true,
			description: "ServiceUnavailable errors should be retried",
		},
		{
			name: "Rejected_Throttling_Error",
			err: &tea.SDKError{
				Code:    tea.String("Rejected.Throttling"),
				Message: tea.String("Request throttled"),
			},
			expected:    true,
			description: "Rejected.Throttling errors should be retried",
		},
		{
			name: "Throttling_Error",
			err: &tea.SDKError{
				Code:    tea.String("Throttling"),
				Message: tea.String("Request throttled"),
			},
			expected:    true,
			description: "Throttling errors should be retried",
		},
		{
			name: "Client_Timeout_Error",
			err: &tea.SDKError{
				Code:    tea.String("Error"),
				Message: tea.String("Client.Timeout exceeded"),
			},
			expected:    true,
			description: "Client timeout errors should be retried",
		},
		{
			name: "5xx_Error",
			err: &tea.SDKError{
				Code:    tea.String("InternalError"),
				Message: tea.String("code: 500 Internal Server Error"),
			},
			expected:    true,
			description: "5xx errors should be retried",
		},
		{
			name: "Regular_Error",
			err: &tea.SDKError{
				Code:    tea.String("InvalidParameter"),
				Message: tea.String("Invalid parameter"),
			},
			expected:    false,
			description: "Regular errors should not be retried",
		},
		{
			name:        "Other_Error",
			err:         errors.New("Some other error"),
			expected:    false,
			description: "Other errors should not be retried",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := needRetry(tc.err)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonIsExpectedErrors(t *testing.T) {
	testCases := []struct {
		name        string
		err         error
		expectCodes []string
		expected    bool
		description string
	}{
		{
			name:        "Nil_Error",
			err:         nil,
			expectCodes: []string{"NotFound"},
			expected:    false,
			description: "Nil error should return false",
		},
		{
			name: "Tea_SDK_Error_Exact_Match",
			err: &tea.SDKError{
				Code: tea.String("NotFound"),
				Data: tea.String(""),
			},
			expectCodes: []string{"NotFound"},
			expected:    true,
			description: "Exact code match should return true",
		},
		{
			name: "Tea_SDK_Error_Prefix_Match",
			err: &tea.SDKError{
				Code: tea.String("NotFound"),
				Data: tea.String(""),
			},
			expectCodes: []string{"NotFound.Instance"},
			expected:    true,
			description: "Prefix match should return true",
		},
		{
			name: "Tea_SDK_Error_Data_Contains",
			err: &tea.SDKError{
				Code: tea.String("Error"),
				Data: tea.String("ResourceNotFound in data"),
			},
			expectCodes: []string{"ResourceNotFound"},
			expected:    true,
			description: "Error code in data should return true",
		},
		{
			name: "Tea_SDK_Error_No_Match",
			err: &tea.SDKError{
				Code: tea.String("InvalidParameter"),
				Data: tea.String(""),
			},
			expectCodes: []string{"NotFound"},
			expected:    false,
			description: "No match should return false",
		},
		{
			name:        "Regular_Error_Contains_Code",
			err:         errors.New("Error: NotFound - Resource not found"),
			expectCodes: []string{"NotFound"},
			expected:    true,
			description: "Error message containing code should return true",
		},
		{
			name:        "Regular_Error_No_Match",
			err:         errors.New("Some error occurred"),
			expectCodes: []string{"NotFound"},
			expected:    false,
			description: "Error message not containing code should return false",
		},
		{
			name: "Multiple_Expected_Codes",
			err: &tea.SDKError{
				Code: tea.String("Throttling"),
				Data: tea.String(""),
			},
			expectCodes: []string{"NotFound", "Throttling", "InvalidParameter"},
			expected:    true,
			description: "Should match one of multiple expected codes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isExpectedErrors(tc.err, tc.expectCodes)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonGetUserAgent(t *testing.T) {
	testCases := []struct {
		name             string
		terraformVersion string
		configSource     string
		traceId          string
		expectedContains []string
		description      string
	}{
		{
			name:             "Standard_UserAgent",
			terraformVersion: "1.5.0",
			configSource:     "terraform-alicloud-modules/vpc",
			traceId:          "trace-123",
			expectedContains: []string{
				"Terraform/1.5.0",
				"Provider/",
				"Module/terraform-alicloud-modules/vpc",
				"TerraformTraceId/trace-123",
			},
			description: "Should generate standard user agent string",
		},
		{
			name:             "Empty_Values",
			terraformVersion: "",
			configSource:     "",
			traceId:          "",
			expectedContains: []string{
				"Terraform/",
				"Provider/",
				"Module/",
				"TerraformTraceId/",
			},
			description: "Should handle empty values",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{
				TerraformVersion:    tc.terraformVersion,
				ConfigurationSource: tc.configSource,
				TerraformTraceId:    tc.traceId,
			}

			userAgent := config.UserAgent()
			assert.NotEmpty(t, userAgent, tc.description)

			for _, expected := range tc.expectedContains {
				assert.Contains(t, userAgent, expected, fmt.Sprintf("%s - should contain %s", tc.description, expected))
			}
		})
	}
}

func TestUnitCommonGetAuthCredential_StsToken(t *testing.T) {
	// When Credential returns error, fields on Config are preserved.
	config := &Config{
		AccessKey:     "test-ak",
		SecretKey:     "test-sk",
		SecurityToken: "test-sts-token",
		Credential:    &mockCredential{err: fmt.Errorf("skip credential override")},
	}
	result := config.AuthCredential(true)
	assert.NotNil(t, result, "Should return a non-nil credential")
	_, ok := result.(*sdkCredentials.StsTokenCredential)
	assert.True(t, ok, "Should return StsTokenCredential when stsSupported=true and SecurityToken is set")
}

func TestUnitCommonGetAuthCredential_StsNotSupported(t *testing.T) {
	config := &Config{
		AccessKey:     "test-ak",
		SecretKey:     "test-sk",
		SecurityToken: "test-sts-token",
		Credential:    &mockCredential{err: fmt.Errorf("skip credential override")},
	}
	result := config.AuthCredential(false)
	assert.NotNil(t, result)
	_, ok := result.(*sdkCredentials.AccessKeyCredential)
	assert.True(t, ok, "Should return AccessKeyCredential when stsSupported=false (SecurityToken is ignored)")
}

func TestUnitCommonGetAuthCredential_RamRoleArn(t *testing.T) {
	config := &Config{
		AccessKey:          "test-ak",
		SecretKey:          "test-sk",
		RamRoleArn:         "acs:ram::123456789012:role/testrole",
		RamRoleSessionName: "test-session",
		Credential:         &mockCredential{err: fmt.Errorf("skip credential override")},
	}
	result := config.AuthCredential(true)
	assert.NotNil(t, result)
	_, ok := result.(*sdkCredentials.RamRoleArnCredential)
	assert.True(t, ok, "Should return RamRoleArnCredential when RamRoleArn is set (no ExternalId)")
}

func TestUnitCommonGetAuthCredential_RamRoleArnWithExternalId(t *testing.T) {
	config := &Config{
		AccessKey:          "test-ak",
		SecretKey:          "test-sk",
		RamRoleArn:         "acs:ram::123456789012:role/testrole",
		RamRoleSessionName: "test-session",
		RamRoleExternalId:  "ext-12345",
		Credential:         &mockCredential{err: fmt.Errorf("skip credential override")},
	}
	result := config.AuthCredential(true)
	assert.NotNil(t, result)
	_, ok := result.(*sdkCredentials.RamRoleArnCredential)
	assert.True(t, ok, "Should return RamRoleArnCredential when RamRoleArn+ExternalId are set")
}

func TestUnitCommonGetAuthCredential_AccessKey(t *testing.T) {
	config := &Config{
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
		Credential: &mockCredential{err: fmt.Errorf("skip credential override")},
	}
	result := config.AuthCredential(true)
	assert.NotNil(t, result)
	_, ok := result.(*sdkCredentials.AccessKeyCredential)
	assert.True(t, ok, "Should return AccessKeyCredential when only AK/SK are set")
}

func TestUnitCommonGetAuthCredential_EcsRoleName(t *testing.T) {
	config := &Config{
		EcsRoleName: "test-ecs-role",
		Credential:  &mockCredential{err: fmt.Errorf("skip credential override")},
	}
	result := config.AuthCredential(true)
	assert.NotNil(t, result)
	_, ok := result.(*sdkCredentials.EcsRamRoleCredential)
	assert.True(t, ok, "Should return EcsRamRoleCredential when only EcsRoleName is set")
}

func TestUnitCommonGetAuthCredential_EmptyCredentials(t *testing.T) {
	config := &Config{
		Credential: &mockCredential{err: fmt.Errorf("skip credential override")},
	}
	result := config.AuthCredential(true)
	assert.NotNil(t, result)
	_, ok := result.(*sdkCredentials.AccessKeyCredential)
	assert.True(t, ok, "Should return empty AccessKeyCredential when no credentials are set")
}

func TestUnitCommonGetAuthCredential_CredentialOverridesFields(t *testing.T) {
	// When Credential.GetCredential() succeeds, it overrides Config fields.
	// We set up a mock that returns AK/SK with a security token.
	mock := &mockCredential{
		accessKeyId:     "override-ak",
		accessKeySecret: "override-sk",
		securityToken:   "override-token",
	}
	config := &Config{
		AccessKey:     "original-ak",
		SecretKey:     "original-sk",
		SecurityToken: "",
		Credential:    mock,
	}
	result := config.AuthCredential(true)
	assert.NotNil(t, result)
	// After credential override, SecurityToken="override-token" → StsTokenCredential
	_, ok := result.(*sdkCredentials.StsTokenCredential)
	assert.True(t, ok, "Credential override should set SecurityToken and return StsTokenCredential")
}
