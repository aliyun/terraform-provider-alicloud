package config

import (
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/stretchr/testify/assert"
)

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

			err := config.validateRegion(ValidRegions)
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

			err := config.loadAndValidate(ValidRegions)
			if tc.expectError {
				assert.Error(t, err, tc.description)
			} else {
				assert.NoError(t, err, tc.description)
			}
		})
	}
}

func TestUnitCommonGetCredentialConfig(t *testing.T) {
	testCases := []struct {
		name          string
		accessKey     string
		secretKey     string
		securityToken string
		ramRoleArn    string
		ecsRoleName   string
		stsSupported  bool
		expectedType  string
		description   string
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

			result := config.NeedRefreshCredential()
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSetAuthByAssumeRole_EmptyAccessKey(t *testing.T) {
	config := &Config{
		AccessKey:  "",
		RamRoleArn: "acs:ram::123456789012:role/testrole",
	}
	err := config.setAuthByAssumeRole()
	assert.NoError(t, err, "Empty AccessKey should cause early return with no error")
	assert.Nil(t, config.Credential, "Credential should remain nil when AccessKey is empty")
}

func TestUnitCommonSetAuthByAssumeRole_EmptyRamRoleArn(t *testing.T) {
	config := &Config{
		AccessKey:  "test-ak",
		SecretKey:  "test-sk",
		RamRoleArn: "",
	}
	err := config.setAuthByAssumeRole()
	assert.NoError(t, err, "Empty RamRoleArn should cause early return with no error")
	assert.Nil(t, config.Credential, "Credential should remain nil when RamRoleArn is empty")
}

func TestUnitCommonSetAuthByAssumeRole_BothEmpty(t *testing.T) {
	config := &Config{}
	err := config.setAuthByAssumeRole()
	assert.NoError(t, err, "Both empty should cause early return with no error")
}

func TestUnitCommonSetAuthByAssumeRole_WithStsEndpoint(t *testing.T) {
	// With StsEndpoint set, but invalid credentials → GetCredential will fail with network error.
	// We verify the error message wraps "refresh Ram Role Arn credential failed".
	config := &Config{
		AccessKey:          "fake-ak",
		SecretKey:          "fake-sk",
		RamRoleArn:         "acs:ram::123456789012:role/testrole",
		RamRoleSessionName: "test-session",
		StsEndpoint:        "sts.cn-hangzhou.aliyuncs.com",
	}
	err := config.setAuthByAssumeRole()
	if err != nil {
		assert.Contains(t, err.Error(), "refresh Ram Role Arn credential failed",
			"Error should mention credential refresh failure")
	}
	// Either succeeds (unlikely without real creds) or fails with the expected error message.
}

func TestUnitCommonSetAuthByAssumeRole_WithTimeouts(t *testing.T) {
	// Ensure timeout config fields are used without panics.
	config := &Config{
		AccessKey:            "fake-ak",
		SecretKey:            "fake-sk",
		RamRoleArn:           "acs:ram::123456789012:role/testrole",
		RamRoleSessionName:   "test-session",
		ClientConnectTimeout: 3000,
		ClientReadTimeout:    5000,
	}
	err := config.setAuthByAssumeRole()
	if err != nil {
		assert.Contains(t, err.Error(), "refresh Ram Role Arn credential failed")
	}
}

func TestUnitCommonSetAuthByAssumeRole_WithSecurityToken(t *testing.T) {
	// SecurityToken path: config.SetSecurityToken is called when c.SecurityToken != "".
	config := &Config{
		AccessKey:          "fake-ak",
		SecretKey:          "fake-sk",
		SecurityToken:      "fake-sts-token",
		RamRoleArn:         "acs:ram::123456789012:role/testrole",
		RamRoleSessionName: "test-session",
	}
	err := config.setAuthByAssumeRole()
	if err != nil {
		assert.Contains(t, err.Error(), "refresh Ram Role Arn credential failed")
	}
}

func TestUnitCommonSetAuthCredentialByEcsRoleName_WithAccessKey(t *testing.T) {
	config := &Config{
		AccessKey:   "existing-ak",
		EcsRoleName: "test-role",
	}
	err := config.setAuthCredentialByEcsRoleName()
	assert.NoError(t, err, "Non-empty AccessKey should cause early return with no error")
	assert.Nil(t, config.Credential, "Credential should remain nil when AccessKey is set")
}

func TestUnitCommonSetAuthCredentialByEcsRoleName_EmptyRoleName(t *testing.T) {
	config := &Config{
		AccessKey:   "",
		EcsRoleName: "",
	}
	err := config.setAuthCredentialByEcsRoleName()
	assert.NoError(t, err, "Empty EcsRoleName should cause early return with no error")
	assert.Nil(t, config.Credential, "Credential should remain nil when EcsRoleName is empty")
}

func TestUnitCommonSetAuthCredentialByEcsRoleName_BothEmpty(t *testing.T) {
	config := &Config{}
	err := config.setAuthCredentialByEcsRoleName()
	assert.NoError(t, err, "Both empty should cause early return with no error")
}

func TestUnitCommonSetAuthCredentialByEcsRoleName_WithRoleName(t *testing.T) {
	// EcsRoleName set, AccessKey empty → attempts to call ECS metadata service.
	// In a non-ECS environment this will fail; verify error message is correct.
	config := &Config{
		AccessKey:   "",
		EcsRoleName: "test-ecs-role",
	}
	err := config.setAuthCredentialByEcsRoleName()
	if err != nil {
		assert.Contains(t, err.Error(), "refresh Ecs Ram Role credential failed",
			"Error should mention ECS credential refresh failure")
	}
	// If running on an actual ECS instance with the role, it should succeed and set Credential.
}
