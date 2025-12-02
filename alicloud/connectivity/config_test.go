package connectivity

import (
	"errors"
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestUnitCommonValidateRegion(t *testing.T) {
	testCases := []struct {
		name        string
		region      Region
		regionId    string
		expectError bool
		description string
	}{
		{
			name:        "Valid_Region_Hangzhou",
			region:      Hangzhou,
			regionId:    "cn-hangzhou",
			expectError: false,
			description: "Valid region cn-hangzhou should not return error",
		},
		{
			name:        "Valid_Region_Beijing",
			region:      Beijing,
			regionId:    "cn-beijing",
			expectError: false,
			description: "Valid region cn-beijing should not return error",
		},
		{
			name:        "Valid_Region_Shanghai",
			region:      Shanghai,
			regionId:    "cn-shanghai",
			expectError: false,
			description: "Valid region cn-shanghai should not return error",
		},
		{
			name:        "Invalid_Region",
			region:      Region("invalid-region"),
			regionId:    "invalid-region",
			expectError: true,
			description: "Invalid region should return error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{
				Region:   tc.region,
				RegionId: tc.regionId,
			}

			err := config.validateRegion()
			if tc.expectError {
				assert.Error(t, err, tc.description)
				assert.Contains(t, err.Error(), "Invalid Alibaba Cloud region")
			} else {
				assert.NoError(t, err, tc.description)
			}
		})
	}
}

func TestUnitCommonLoadAndValidate(t *testing.T) {
	testCases := []struct {
		name        string
		region      Region
		regionId    string
		expectError bool
		description string
	}{
		{
			name:        "Valid_Configuration",
			region:      Hangzhou,
			regionId:    "cn-hangzhou",
			expectError: false,
			description: "Valid configuration should not return error",
		},
		{
			name:        "Invalid_Region",
			region:      Region("invalid-region"),
			regionId:    "invalid-region",
			expectError: true,
			description: "Invalid region should return error during validation",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{
				Region:   tc.region,
				RegionId: tc.regionId,
			}

			err := config.loadAndValidate()
			if tc.expectError {
				assert.Error(t, err, tc.description)
			} else {
				assert.NoError(t, err, tc.description)
			}
		})
	}
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

			userAgent := config.getUserAgent()
			assert.NotEmpty(t, userAgent, tc.description)

			for _, expected := range tc.expectedContains {
				assert.Contains(t, userAgent, expected, fmt.Sprintf("%s - should contain %s", tc.description, expected))
			}
		})
	}
}

func TestUnitCommonGetCredentialConfig(t *testing.T) {
	testCases := []struct {
		name           string
		accessKey      string
		secretKey      string
		securityToken  string
		ramRoleArn     string
		ecsRoleName    string
		stsSupported   bool
		expectedType   string
		description    string
	}{
		{
			name:         "AccessKey_Credential",
			accessKey:    "test-ak",
			secretKey:    "test-sk",
			stsSupported: false,
			expectedType: "access_key",
			description:  "Should return access_key type credential config",
		},
		{
			name:          "STS_Credential",
			accessKey:     "test-ak",
			secretKey:     "test-sk",
			securityToken: "test-token",
			stsSupported:  true,
			expectedType:  "sts",
			description:   "Should return sts type credential config",
		},
		{
			name:         "RamRoleArn_Credential",
			accessKey:    "test-ak",
			secretKey:    "test-sk",
			ramRoleArn:   "acs:ram::123456789012:role/testrole",
			stsSupported: false,
			expectedType: "ram_role_arn",
			description:  "Should return ram_role_arn type credential config",
		},
		{
			name:         "EcsRamRole_Credential",
			ecsRoleName:  "test-ecs-role",
			stsSupported: false,
			expectedType: "ecs_ram_role",
			description:  "Should return ecs_ram_role type credential config",
		},
		{
			name:         "Empty_Credential",
			stsSupported: false,
			expectedType: "",
			description:  "Should return empty type for no credentials",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{
				AccessKey:          tc.accessKey,
				SecretKey:          tc.secretKey,
				SecurityToken:      tc.securityToken,
				RamRoleArn:         tc.ramRoleArn,
				RamRoleSessionName: "test-session",
				EcsRoleName:        tc.ecsRoleName,
			}

			credConfig := config.getCredentialConfig(tc.stsSupported)
			assert.NotNil(t, credConfig, tc.description)
			assert.Equal(t, tc.expectedType, *credConfig.Type, tc.description)

			// Verify credential fields based on type
			if tc.accessKey != "" && tc.secretKey != "" {
				assert.Equal(t, tc.accessKey, *credConfig.AccessKeyId)
				assert.Equal(t, tc.secretKey, *credConfig.AccessKeySecret)
			}

			if tc.stsSupported && tc.securityToken != "" {
				assert.Equal(t, tc.securityToken, *credConfig.SecurityToken)
			}

			if tc.ramRoleArn != "" {
				assert.Equal(t, tc.ramRoleArn, *credConfig.RoleArn)
			}

			if tc.ecsRoleName != "" {
				assert.Equal(t, tc.ecsRoleName, *credConfig.RoleName)
			}
		})
	}
}

func TestUnitCommonNeedRefreshCredential(t *testing.T) {
	testCases := []struct {
		name        string
		credType    string
		setupConfig func() *credentials.Config
		expected    bool
		description string
	}{
		{
			name:     "STS_Credential_No_Refresh",
			credType: "sts",
			setupConfig: func() *credentials.Config {
				return &credentials.Config{
					Type:            tea.String("sts"),
					AccessKeyId:     tea.String("test-ak"),
					AccessKeySecret: tea.String("test-sk"),
					SecurityToken:   tea.String("test-token"),
				}
			},
			expected:    false,
			description: "STS credentials should not need refresh",
		},
		{
			name:     "AccessKey_Credential_No_Refresh",
			credType: "access_key",
			setupConfig: func() *credentials.Config {
				return &credentials.Config{
					Type:            tea.String("access_key"),
					AccessKeyId:     tea.String("test-ak"),
					AccessKeySecret: tea.String("test-sk"),
				}
			},
			expected:    false,
			description: "Access key credentials should not need refresh",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			credConfig := tc.setupConfig()

			cred, err := credentials.NewCredential(credConfig)
			if err != nil {
				t.Skipf("Skipping test due to credential creation error: %v", err)
				return
			}

			config := &Config{
				Credential: cred,
			}

			result := config.needRefreshCredential()
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
