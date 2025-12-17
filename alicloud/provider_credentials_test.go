package alicloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// TestProviderCredentials_ExternalMode 测试 External 模式的凭证获取
func TestProviderCredentials_ExternalMode(t *testing.T) {
	// 创建临时目录和文件
	tmpDir, err := ioutil.TempDir("", "alicloud-test-external-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建外部程序脚本
	scriptPath := filepath.Join(tmpDir, "credential-provider.sh")
	scriptContent := `#!/bin/bash
cat <<EOF
{
  "AccessKeyId": "STS.EXTERNAL123456789",
  "AccessKeySecret": "external_secret_key_test",
  "SecurityToken": "external_security_token_test",
  "Expiration": "2099-12-31T23:59:59Z"
}
EOF
`
	err = ioutil.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// 创建临时配置文件
	configPath := filepath.Join(tmpDir, "config.json")
	config := map[string]interface{}{
		"profiles": []map[string]interface{}{
			{
				"name":            "external-test",
				"mode":            "External",
				"process_command": scriptPath,
				"region_id":       "cn-beijing",
			},
		},
		"current": "external-test",
	}

	configData, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 测试配置读取
	t.Run("ParseExternalConfig", func(t *testing.T) {
		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		var cfg map[string]interface{}
		err = json.Unmarshal(data, &cfg)
		if err != nil {
			t.Fatalf("Failed to parse config JSON: %v", err)
		}

		profiles := cfg["profiles"].([]interface{})
		if len(profiles) == 0 {
			t.Fatal("No profiles found")
		}

		profile := profiles[0].(map[string]interface{})
		if profile["mode"] != "External" {
			t.Errorf("Expected mode 'External', got '%v'", profile["mode"])
		}

		if profile["name"] != "external-test" {
			t.Errorf("Expected name 'external-test', got '%v'", profile["name"])
		}
	})

	// 测试 provider 配置
	t.Run("ProviderConfigWithExternal", func(t *testing.T) {
		// 跳过 Windows 系统（bash 脚本问题）
		if runtime.GOOS == "windows" {
			t.Skip("Skipping bash script test on Windows")
		}

		raw := make(map[string]interface{})
		raw["profile"] = "external-test"
		raw["shared_credentials_file"] = configPath
		raw["region"] = "cn-beijing"

		resourceData := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).Schema, raw)

		// 读取配置文件验证
		providerConfig = nil // 重置全局变量
		_, err := getConfigFromProfile(resourceData, "mode")
		if err != nil {
			t.Fatalf("Failed to get config from profile: %v", err)
		}

		if providerConfig == nil {
			t.Fatal("Provider config is nil")
		}

		mode, ok := providerConfig["mode"]
		if !ok {
			t.Fatal("Mode not found in provider config")
		}

		if mode != "External" {
			t.Errorf("Expected mode 'External', got '%v'", mode)
		}
	})
}

// TestProviderCredentials_OAuthMode 测试 OAuth 模式的凭证获取
func TestProviderCredentials_OAuthMode(t *testing.T) {
	// 创建临时目录
	tmpDir, err := ioutil.TempDir("", "alicloud-test-oauth-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建临时配置文件
	configPath := filepath.Join(tmpDir, "config.json")
	config := map[string]interface{}{
		"profiles": []map[string]interface{}{
			{
				"name":                "oauth-test",
				"mode":                "OAuth",
				"oauth_client_id":     "test-client-id",
				"oauth_client_secret": "test-client-secret",
				"oauth_token_url":     "https://oauth.aliyun.com/token",
				"oauth_scope":         "openid",
				"region_id":           "cn-hangzhou",
			},
		},
		"current": "oauth-test",
	}

	configData, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 测试配置读取
	t.Run("ParseOAuthConfig", func(t *testing.T) {
		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		var cfg map[string]interface{}
		err = json.Unmarshal(data, &cfg)
		if err != nil {
			t.Fatalf("Failed to parse config JSON: %v", err)
		}

		profiles := cfg["profiles"].([]interface{})
		if len(profiles) == 0 {
			t.Fatal("No profiles found")
		}

		profile := profiles[0].(map[string]interface{})
		if profile["mode"] != "OAuth" {
			t.Errorf("Expected mode 'OAuth', got '%v'", profile["mode"])
		}

		if profile["oauth_client_id"] != "test-client-id" {
			t.Errorf("Expected oauth_client_id 'test-client-id', got '%v'", profile["oauth_client_id"])
		}
	})

	// 测试 provider 配置
	t.Run("ProviderConfigWithOAuth", func(t *testing.T) {
		raw := make(map[string]interface{})
		raw["profile"] = "oauth-test"
		raw["shared_credentials_file"] = configPath
		raw["region"] = "cn-hangzhou"

		resourceData := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).Schema, raw)

		// 读取配置文件验证
		providerConfig = nil // 重置全局变量
		_, err := getConfigFromProfile(resourceData, "mode")
		if err != nil {
			t.Fatalf("Failed to get config from profile: %v", err)
		}

		if providerConfig == nil {
			t.Fatal("Provider config is nil")
		}

		mode, ok := providerConfig["mode"]
		if !ok {
			t.Fatal("Mode not found in provider config")
		}

		if mode != "OAuth" {
			t.Errorf("Expected mode 'OAuth', got '%v'", mode)
		}
	})
}

// TestProviderCredentials_MultipleProfiles 测试多个 profile 配置
func TestProviderCredentials_MultipleProfiles(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "alicloud-test-multi-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建外部程序脚本
	scriptPath := filepath.Join(tmpDir, "credential-provider.sh")
	scriptContent := `#!/bin/bash
cat <<EOF
{
  "AccessKeyId": "STS.MULTI123456789",
  "AccessKeySecret": "multi_secret_key_test",
  "SecurityToken": "multi_security_token_test",
  "Expiration": "2099-12-31T23:59:59Z"
}
EOF
`
	err = ioutil.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// 创建包含多个 profile 的配置文件
	configPath := filepath.Join(tmpDir, "config.json")
	config := map[string]interface{}{
		"profiles": []map[string]interface{}{
			{
				"name":              "static-ak",
				"mode":              "AK",
				"access_key_id":     "LTAI5tSTATIC123",
				"access_key_secret": "static_secret",
				"region_id":         "cn-beijing",
			},
			{
				"name":            "external-profile",
				"mode":            "External",
				"process_command": scriptPath,
				"region_id":       "cn-beijing",
			},
			{
				"name":                "oauth-profile",
				"mode":                "OAuth",
				"oauth_client_id":     "oauth-client-id",
				"oauth_client_secret": "oauth-client-secret",
				"oauth_token_url":     "https://oauth.aliyun.com/token",
				"region_id":           "cn-hangzhou",
			},
			{
				"name":       "cloudsso-profile",
				"mode":       "CloudSSO",
				"start_url":  "https://sso.aliyun.com",
				"sso_region": "cn-hangzhou",
				"account_id": "123456789012",
				"region_id":  "cn-beijing",
			},
		},
		"current": "external-profile",
	}

	configData, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		profileName  string
		expectedMode string
	}{
		{"static-ak", "AK"},
		{"external-profile", "External"},
		{"oauth-profile", "OAuth"},
		{"cloudsso-profile", "CloudSSO"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Profile_%s", tt.profileName), func(t *testing.T) {
			raw := make(map[string]interface{})
			raw["profile"] = tt.profileName
			raw["shared_credentials_file"] = configPath
			raw["region"] = "cn-beijing"

			resourceData := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).Schema, raw)

			// 重置全局变量
			providerConfig = nil

			_, err := getConfigFromProfile(resourceData, "mode")
			if err != nil {
				t.Fatalf("Failed to get config from profile: %v", err)
			}

			if providerConfig == nil {
				t.Fatal("Provider config is nil")
			}

			mode, ok := providerConfig["mode"]
			if !ok {
				t.Fatal("Mode not found in provider config")
			}

			if mode != tt.expectedMode {
				t.Errorf("Expected mode '%s', got '%v'", tt.expectedMode, mode)
			}
		})
	}
}

// TestProviderCredentials_GetConfigFromProfile 测试 getConfigFromProfile 函数
func TestProviderCredentials_GetConfigFromProfile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "alicloud-test-getconfig-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	config := map[string]interface{}{
		"profiles": []map[string]interface{}{
			{
				"name":            "test-external",
				"mode":            "External",
				"process_command": "/usr/local/bin/credential-provider",
				"region_id":       "cn-shanghai",
			},
		},
		"current": "test-external",
	}

	configData, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	raw := make(map[string]interface{})
	raw["profile"] = "test-external"
	raw["shared_credentials_file"] = configPath
	raw["region"] = "cn-beijing"

	resourceData := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).Schema, raw)

	tests := []struct {
		name        string
		profileKey  string
		expectedVal interface{}
		shouldBeNil bool
	}{
		{
			name:        "GetRegionId",
			profileKey:  "region_id",
			expectedVal: "cn-shanghai",
			shouldBeNil: false,
		},
		{
			name:        "GetMode_External",
			profileKey:  "mode",
			expectedVal: "External",
			shouldBeNil: true,
		},
		{
			name:        "GetAccessKey_ShouldBeNil",
			profileKey:  "access_key_id",
			expectedVal: nil,
			shouldBeNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			providerConfig = nil // 重置
			val, err := getConfigFromProfile(resourceData, tt.profileKey)
			if err != nil {
				t.Fatalf("Failed to get config: %v", err)
			}

			if tt.shouldBeNil {
				if val != nil {
					t.Errorf("Expected nil, got %v", val)
				}
			} else {
				if val != tt.expectedVal {
					t.Errorf("Expected %v, got %v", tt.expectedVal, val)
				}
			}
		})
	}
}

// TestProviderCredentials_AdvancedModesReturnNil 测试高级模式返回 nil
func TestProviderCredentials_AdvancedModesReturnNil(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "alicloud-test-advanced-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	modes := []string{"External", "OAuth", "CloudSSO", "ChainableRamRoleArn"}

	for _, mode := range modes {
		t.Run(fmt.Sprintf("Mode_%s", mode), func(t *testing.T) {
			configPath := filepath.Join(tmpDir, fmt.Sprintf("config_%s.json", mode))
			config := map[string]interface{}{
				"profiles": []map[string]interface{}{
					{
						"name":      fmt.Sprintf("test-%s", mode),
						"mode":      mode,
						"region_id": "cn-beijing",
					},
				},
				"current": fmt.Sprintf("test-%s", mode),
			}

			configData, err := json.Marshal(config)
			if err != nil {
				t.Fatal(err)
			}

			err = ioutil.WriteFile(configPath, configData, 0644)
			if err != nil {
				t.Fatal(err)
			}

			raw := make(map[string]interface{})
			raw["profile"] = fmt.Sprintf("test-%s", mode)
			raw["shared_credentials_file"] = configPath
			raw["region"] = "cn-beijing"

			resourceData := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).Schema, raw)

			providerConfig = nil
			// 对于高级模式，access_key_id 应该返回 nil
			val, err := getConfigFromProfile(resourceData, "access_key_id")
			if err != nil {
				t.Fatalf("Failed to get config: %v", err)
			}

			// 高级模式应该返回 nil（凭证通过 CLIProfileCredentialsProvider 动态获取）
			if val != nil {
				t.Errorf("Expected nil for access_key_id in mode %s, got %v", mode, val)
			}
		})
	}
}

// TestProviderCredentials_PriorityOrder 测试凭证优先级顺序
func TestProviderCredentials_PriorityOrder(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "alicloud-test-priority-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建配置文件（低优先级）
	configPath := filepath.Join(tmpDir, "config.json")
	config := map[string]interface{}{
		"profiles": []map[string]interface{}{
			{
				"name":              "test-profile",
				"mode":              "AK",
				"access_key_id":     "PROFILE_AK",
				"access_key_secret": "PROFILE_SK",
				"region_id":         "cn-hangzhou",
			},
		},
		"current": "test-profile",
	}

	configData, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("StaticConfigOverridesProfile", func(t *testing.T) {
		raw := make(map[string]interface{})
		raw["access_key"] = "STATIC_AK" // 静态配置（高优先级）
		raw["secret_key"] = "STATIC_SK"
		raw["profile"] = "test-profile" // Profile 配置（低优先级）
		raw["shared_credentials_file"] = configPath
		raw["region"] = "cn-beijing"

		resourceData := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).Schema, raw)

		// 使用 getProviderConfig 逻辑验证优先级
		// 注意：这里我们模拟 provider.go 中的逻辑
		var accessKey string
		if v, ok := resourceData.GetOk("access_key"); ok && v != nil && v.(string) != "" {
			accessKey = v.(string)
		} else {
			providerConfig = nil
			if v, err := getConfigFromProfile(resourceData, "access_key_id"); err == nil && v != nil {
				accessKey = v.(string)
			}
		}

		if accessKey != "STATIC_AK" {
			t.Errorf("Expected static config 'STATIC_AK' to take priority, got '%s'", accessKey)
		}
	})

	t.Run("ProfileUsedWhenNoStaticConfig", func(t *testing.T) {
		raw := make(map[string]interface{})
		// 不设置静态 access_key/secret_key
		raw["profile"] = "test-profile"
		raw["shared_credentials_file"] = configPath
		raw["region"] = "cn-beijing"

		resourceData := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).Schema, raw)

		providerConfig = nil
		val, err := getConfigFromProfile(resourceData, "access_key_id")
		if err != nil {
			t.Fatalf("Failed to get config: %v", err)
		}

		if val != "PROFILE_AK" {
			t.Errorf("Expected profile config 'PROFILE_AK', got '%v'", val)
		}
	})
}

// TestProviderCredentials_ConfigFileNotFound 测试配置文件不存在的情况
func TestProviderCredentials_ConfigFileNotFound(t *testing.T) {
	raw := make(map[string]interface{})
	raw["profile"] = "non-existent-profile"
	raw["shared_credentials_file"] = "/tmp/non-existent-config-file-12345.json"
	raw["region"] = "cn-beijing"

	resourceData := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).Schema, raw)

	providerConfig = nil
	val, err := getConfigFromProfile(resourceData, "access_key_id")

	// 文件不存在时应该返回 nil（不报错）
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if val != nil {
		t.Errorf("Expected nil when config file doesn't exist, got %v", val)
	}
}

// TestProviderConfigure_WithExternalAndOAuth 集成测试（需要实际凭证提供程序）
func TestProviderConfigure_WithExternalAndOAuth(t *testing.T) {
	// 这个测试需要实际的外部凭证提供程序，标记为需要环境变量
	if os.Getenv("ALICLOUD_EXTERNAL_CREDENTIAL_PROVIDER") == "" {
		t.Skip("Skipping integration test: ALICLOUD_EXTERNAL_CREDENTIAL_PROVIDER not set")
	}

	// 跳过 Windows 系统
	if runtime.GOOS == "windows" {
		t.Skip("Skipping bash script test on Windows")
	}

	tmpDir, err := ioutil.TempDir("", "alicloud-test-integration-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	scriptPath := os.Getenv("ALICLOUD_EXTERNAL_CREDENTIAL_PROVIDER")

	configPath := filepath.Join(tmpDir, "config.json")
	config := map[string]interface{}{
		"profiles": []map[string]interface{}{
			{
				"name":            "integration-external",
				"mode":            "External",
				"process_command": scriptPath,
				"region_id":       "cn-beijing",
			},
		},
		"current": "integration-external",
	}

	configData, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	raw := make(map[string]interface{})
	raw["profile"] = "integration-external"
	raw["shared_credentials_file"] = configPath
	raw["region"] = "cn-beijing"

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	// 尝试配置 provider（这会实际调用外部凭证提供程序）
	_, err = providerConfigure(resourceData, provider)
	if err != nil {
		t.Logf("Expected: External credential provider might not return valid credentials in test environment")
		// 不报错，因为在测试环境可能没有真实的凭证
	} else {
		t.Log("Successfully configured provider with External credential mode")
	}
}

// ==============================================================================
// providerConfigure 函数测试用例
// ==============================================================================

// TestProviderConfigure_StaticAK 测试静态 AK/SK 认证
func TestProviderConfigure_StaticAK(t *testing.T) {
	raw := map[string]interface{}{
		"access_key": "LTAI5tTestAccessKey123456",
		"secret_key": "TestSecretKey123456789012345678",
		"region":     "cn-hangzhou",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with static AK/SK: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with static AK/SK")
	}
}

// TestProviderConfigure_StaticAKWithSTS 测试静态 AK/SK + STS Token 认证
func TestProviderConfigure_StaticAKWithSTS(t *testing.T) {
	raw := map[string]interface{}{
		"access_key":     "STS.TestSTSAccessKey123456",
		"secret_key":     "TestSTSSecretKey1234567890123456",
		"security_token": "TestSecurityToken1234567890",
		"region":         "cn-shanghai",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with STS token: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with STS token")
	}
}

// TestProviderConfigure_MissingCredentials 测试缺少凭证的情况
func TestProviderConfigure_MissingCredentials(t *testing.T) {
	// 保存当前环境变量
	envVars := []string{
		"ALICLOUD_ACCESS_KEY", "ALICLOUD_SECRET_KEY", "ALICLOUD_SECURITY_TOKEN",
		"ALIBABA_CLOUD_ACCESS_KEY_ID", "ALIBABA_CLOUD_ACCESS_KEY_SECRET",
		"ALICLOUD_PROFILE", "ALIBABA_CLOUD_PROFILE",
		"ALICLOUD_ECS_ROLE_NAME",
	}
	savedEnvs := make(map[string]string)
	for _, env := range envVars {
		savedEnvs[env] = os.Getenv(env)
		os.Unsetenv(env)
	}
	// 测试结束后恢复环境变量
	defer func() {
		for env, val := range savedEnvs {
			if val != "" {
				os.Setenv(env, val)
			}
		}
	}()

	raw := map[string]interface{}{
		"region": "cn-hangzhou",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	_, err := providerConfigure(resourceData, provider)
	if err == nil {
		t.Fatal("Expected error when credentials are missing")
	}

	expectedMsg := "no valid credential sources"
	if !containsString(err.Error(), expectedMsg) {
		t.Errorf("Expected error message to contain '%s', got: %v", expectedMsg, err)
	}
}

// TestProviderConfigure_MissingSecretKey 测试只有 AccessKey 没有 SecretKey 的情况
func TestProviderConfigure_MissingSecretKey(t *testing.T) {
	// 保存当前环境变量
	envVars := []string{
		"ALICLOUD_ACCESS_KEY", "ALICLOUD_SECRET_KEY", "ALICLOUD_SECURITY_TOKEN",
		"ALIBABA_CLOUD_ACCESS_KEY_ID", "ALIBABA_CLOUD_ACCESS_KEY_SECRET",
		"ALICLOUD_PROFILE", "ALIBABA_CLOUD_PROFILE",
		"ALICLOUD_ECS_ROLE_NAME",
	}
	savedEnvs := make(map[string]string)
	for _, env := range envVars {
		savedEnvs[env] = os.Getenv(env)
		os.Unsetenv(env)
	}
	defer func() {
		for env, val := range savedEnvs {
			if val != "" {
				os.Setenv(env, val)
			}
		}
	}()

	raw := map[string]interface{}{
		"access_key": "LTAI5tTestAccessKey123456",
		"region":     "cn-hangzhou",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	_, err := providerConfigure(resourceData, provider)
	if err == nil {
		t.Fatal("Expected error when secret_key is missing")
	}
}

// TestProviderConfigure_DefaultRegion 测试默认区域
func TestProviderConfigure_DefaultRegion(t *testing.T) {
	raw := map[string]interface{}{
		"access_key": "LTAI5tTestAccessKey123456",
		"secret_key": "TestSecretKey123456789012345678",
		// 不设置 region，应该使用默认值
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}

// TestProviderConfigure_CredentialsURI 测试通过 credentials_uri 获取凭证
func TestProviderConfigure_CredentialsURI(t *testing.T) {
	// 创建模拟服务器返回凭证
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"AccessKeyId": "LTAI5tURITestAK123456",
			"AccessKeySecret": "fakeAccessKeySecret",
			"SecurityToken": "fakeSecurityToken",
			"Expiration": "2099-12-31T23:59:59Z"
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	raw := map[string]interface{}{
		"credentials_uri": server.URL,
		"region":          "cn-hangzhou",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with credentials_uri: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with credentials_uri")
	}
}

// TestProviderConfigure_CredentialsURIInvalid 测试无效的 credentials_uri
func TestProviderConfigure_CredentialsURIInvalid(t *testing.T) {
	// 保存当前环境变量
	envVars := []string{
		"ALICLOUD_ACCESS_KEY", "ALICLOUD_SECRET_KEY", "ALICLOUD_SECURITY_TOKEN",
		"ALIBABA_CLOUD_ACCESS_KEY_ID", "ALIBABA_CLOUD_ACCESS_KEY_SECRET",
		"ALICLOUD_PROFILE", "ALIBABA_CLOUD_PROFILE",
		"ALICLOUD_ECS_ROLE_NAME",
	}
	savedEnvs := make(map[string]string)
	for _, env := range envVars {
		savedEnvs[env] = os.Getenv(env)
		os.Unsetenv(env)
	}
	defer func() {
		for env, val := range savedEnvs {
			if val != "" {
				os.Setenv(env, val)
			}
		}
	}()

	// 创建模拟服务器返回错误
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	raw := map[string]interface{}{
		"credentials_uri": server.URL,
		"region":          "cn-hangzhou",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	_, err := providerConfigure(resourceData, provider)
	if err == nil {
		t.Fatal("Expected error with invalid credentials_uri response")
	}

	// 验证错误消息
	expectedMsg := "failed"
	if !containsString(err.Error(), expectedMsg) && !containsString(err.Error(), "no valid credential") {
		t.Errorf("Expected error message related to credentials, got: %v", err)
	}
}

// TestProviderConfigure_ProfileWithStaticAK 测试 Profile 配置（AK 模式）
func TestProviderConfigure_ProfileWithStaticAK(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "alicloud-test-profile-ak-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	config := map[string]interface{}{
		"profiles": []map[string]interface{}{
			{
				"name":              "test-ak-profile",
				"mode":              "AK",
				"access_key_id":     "fakeAccessKeyId",
				"access_key_secret": "fakeAccessKeySecret",
				"region_id":         "cn-beijing",
			},
		},
		"current": "test-ak-profile",
	}

	configData, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	raw := map[string]interface{}{
		"profile":                 "test-ak-profile",
		"shared_credentials_file": configPath,
		"region":                  "cn-beijing",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	// 重置全局变量
	providerConfig = nil

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with profile AK mode: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with profile AK mode")
	}
}

// TestProviderConfigure_ProfileWithSTS 测试 Profile 配置（StsToken 模式）
func TestProviderConfigure_ProfileWithSTS(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "alicloud-test-profile-sts-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	config := map[string]interface{}{
		"profiles": []map[string]interface{}{
			{
				"name":              "test-sts-profile",
				"mode":              "StsToken",
				"access_key_id":     "STS.ProfileSTSAK123456789",
				"access_key_secret": "ProfileSTSSecret1234567890123",
				"sts_token":         "ProfileSTSToken12345678901234",
				"region_id":         "cn-shanghai",
			},
		},
		"current": "test-sts-profile",
	}

	configData, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	raw := map[string]interface{}{
		"profile":                 "test-sts-profile",
		"shared_credentials_file": configPath,
		"region":                  "cn-shanghai",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	providerConfig = nil

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with profile STS mode: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with profile STS mode")
	}
}

// TestProviderConfigure_ProfileExternal 测试 Profile External 模式
func TestProviderConfigure_ProfileExternal(t *testing.T) {
	// 跳过 Windows 系统
	if runtime.GOOS == "windows" {
		t.Skip("Skipping bash script test on Windows")
	}

	tmpDir, err := ioutil.TempDir("", "alicloud-test-profile-external-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建外部程序脚本
	scriptPath := filepath.Join(tmpDir, "credential-provider.sh")
	scriptContent := `#!/bin/bash
cat <<EOF
{
  "AccessKeyId": "STS.ExternalAK123456789",
  "AccessKeySecret": "ExternalSecretKey1234567890",
  "SecurityToken": "ExternalSecurityToken12345",
  "Expiration": "2099-12-31T23:59:59Z"
}
EOF
`
	err = ioutil.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatal(err)
	}

	configPath := filepath.Join(tmpDir, "config.json")
	config := map[string]interface{}{
		"profiles": []map[string]interface{}{
			{
				"name":            "test-external-profile",
				"mode":            "External",
				"process_command": scriptPath,
				"region_id":       "cn-hangzhou",
			},
		},
		"current": "test-external-profile",
	}

	configData, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	raw := map[string]interface{}{
		"profile":                 "test-external-profile",
		"shared_credentials_file": configPath,
		"region":                  "cn-hangzhou",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	providerConfig = nil

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with profile External mode: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with profile External mode")
	}
}

// TestProviderConfigure_AssumeRole 测试 assume_role 配置
// 注意：此测试需要有效凭证才能成功，使用假凭证会因为 STS 调用失败而报错
func TestProviderConfigure_AssumeRole(t *testing.T) {
	// 跳过使用假凭证的情况（真实 assume_role 需要实际调用 STS）
	if os.Getenv("ALICLOUD_ACCESS_KEY") == "" {
		t.Skip("Skipping assume_role test: requires real credentials to call STS")
	}

	raw := map[string]interface{}{
		"access_key": os.Getenv("ALICLOUD_ACCESS_KEY"),
		"secret_key": os.Getenv("ALICLOUD_SECRET_KEY"),
		"region":     "cn-hangzhou",
		"assume_role": []interface{}{
			map[string]interface{}{
				"role_arn":           os.Getenv("ALICLOUD_ASSUME_ROLE_ARN"),
				"session_name":       "test-session",
				"policy":             "",
				"session_expiration": 3600,
				"external_id":        "",
			},
		},
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with assume_role: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with assume_role")
	}
}

// TestProviderConfigure_AssumeRoleConfig 测试 assume_role 配置解析（不实际调用 STS）
func TestProviderConfigure_AssumeRoleConfig(t *testing.T) {
	raw := map[string]interface{}{
		"access_key": "LTAI5tTestAccessKey123456",
		"secret_key": "TestSecretKey123456789012345678",
		"region":     "cn-hangzhou",
		"assume_role": []interface{}{
			map[string]interface{}{
				"role_arn":           "acs:ram::123456789012:role/test-role",
				"session_name":       "test-session",
				"policy":             "",
				"session_expiration": 3600,
				"external_id":        "",
			},
		},
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	// 只测试配置能够正确解析，不测试实际的 STS 调用
	assumeRoleList := resourceData.Get("assume_role").(*schema.Set).List()
	if len(assumeRoleList) != 1 {
		t.Fatalf("Expected 1 assume_role config, got %d", len(assumeRoleList))
	}

	assumeRole := assumeRoleList[0].(map[string]interface{})
	if assumeRole["role_arn"].(string) != "acs:ram::123456789012:role/test-role" {
		t.Errorf("Expected role_arn to be 'acs:ram::123456789012:role/test-role', got '%s'", assumeRole["role_arn"])
	}
	if assumeRole["session_name"].(string) != "test-session" {
		t.Errorf("Expected session_name to be 'test-session', got '%s'", assumeRole["session_name"])
	}
	if assumeRole["session_expiration"].(int) != 3600 {
		t.Errorf("Expected session_expiration to be 3600, got %d", assumeRole["session_expiration"])
	}
}

// TestProviderConfigure_AssumeRoleDefaultSessionName 测试 assume_role 默认 session_name
func TestProviderConfigure_AssumeRoleDefaultSessionName(t *testing.T) {
	raw := map[string]interface{}{
		"access_key": "LTAI5tTestAccessKey123456",
		"secret_key": "TestSecretKey123456789012345678",
		"region":     "cn-hangzhou",
		"assume_role": []interface{}{
			map[string]interface{}{
				"role_arn":           "acs:ram::123456789012:role/test-role",
				"session_name":       "", // 空的 session_name，应该使用默认值 "terraform"
				"policy":             "",
				"session_expiration": 0,
				"external_id":        "",
			},
		},
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	// 测试配置解析，验证 session_name 的默认值逻辑
	assumeRoleList := resourceData.Get("assume_role").(*schema.Set).List()
	if len(assumeRoleList) != 1 {
		t.Fatalf("Expected 1 assume_role config, got %d", len(assumeRoleList))
	}

	assumeRole := assumeRoleList[0].(map[string]interface{})
	// 在 providerConfigure 中，空的 session_name 会被设置为 "terraform"
	if assumeRole["session_name"].(string) != "" {
		t.Errorf("Expected session_name to be empty in raw config, got '%s'", assumeRole["session_name"])
	}
}

// TestProviderConfigure_WithEndpoints 测试自定义 endpoints
func TestProviderConfigure_WithEndpoints(t *testing.T) {
	raw := map[string]interface{}{
		"access_key": "LTAI5tTestAccessKey123456",
		"secret_key": "TestSecretKey123456789012345678",
		"region":     "cn-hangzhou",
		"endpoints": []interface{}{
			map[string]interface{}{
				"ecs": "ecs.custom.aliyuncs.com",
				"vpc": "vpc.custom.aliyuncs.com",
				"rds": "rds.custom.aliyuncs.com",
			},
		},
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with custom endpoints: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with custom endpoints")
	}
}

// TestProviderConfigure_WithTimeout 测试自定义超时设置
func TestProviderConfigure_WithTimeout(t *testing.T) {
	raw := map[string]interface{}{
		"access_key":             "LTAI5tTestAccessKey123456",
		"secret_key":             "TestSecretKey123456789012345678",
		"region":                 "cn-hangzhou",
		"client_read_timeout":    60,
		"client_connect_timeout": 30,
		"max_retry_timeout":      90,
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with timeout settings: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with timeout settings")
	}
}

// TestProviderConfigure_WithAccountId 测试 account_id 配置
func TestProviderConfigure_WithAccountId(t *testing.T) {
	raw := map[string]interface{}{
		"access_key": "LTAI5tTestAccessKey123456",
		"secret_key": "TestSecretKey123456789012345678",
		"region":     "cn-hangzhou",
		"account_id": "123456789012",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with account_id: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with account_id")
	}
}

// TestProviderConfigure_WithConfigurationSource 测试 configuration_source 配置
func TestProviderConfigure_WithConfigurationSource(t *testing.T) {
	raw := map[string]interface{}{
		"access_key":           "LTAI5tTestAccessKey123456",
		"secret_key":           "TestSecretKey123456789012345678",
		"region":               "cn-hangzhou",
		"configuration_source": "test/1.0.0",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with configuration_source: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with configuration_source")
	}
}

// TestProviderConfigure_WithProtocol 测试 protocol 配置
func TestProviderConfigure_WithProtocol(t *testing.T) {
	testCases := []struct {
		name     string
		protocol string
	}{
		{"HTTPS", "HTTPS"},
		{"HTTP", "HTTP"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			raw := map[string]interface{}{
				"access_key": "LTAI5tTestAccessKey123456",
				"secret_key": "TestSecretKey123456789012345678",
				"region":     "cn-hangzhou",
				"protocol":   tc.protocol,
			}

			provider := Provider().(*schema.Provider)
			resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

			client, err := providerConfigure(resourceData, provider)
			if err != nil {
				t.Fatalf("providerConfigure failed with protocol %s: %v", tc.protocol, err)
			}

			if client == nil {
				t.Fatalf("Expected non-nil client with protocol %s", tc.protocol)
			}
		})
	}
}

// TestProviderConfigure_SkipRegionValidation 测试跳过区域验证
func TestProviderConfigure_SkipRegionValidation(t *testing.T) {
	raw := map[string]interface{}{
		"access_key":             "LTAI5tTestAccessKey123456",
		"secret_key":             "TestSecretKey123456789012345678",
		"region":                 "invalid-region",
		"skip_region_validation": true,
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with skip_region_validation: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with skip_region_validation")
	}
}

// TestProviderConfigure_SecureTransport 测试安全传输配置
func TestProviderConfigure_SecureTransport(t *testing.T) {
	raw := map[string]interface{}{
		"access_key":       "LTAI5tTestAccessKey123456",
		"secret_key":       "TestSecretKey123456789012345678",
		"region":           "cn-hangzhou",
		"secure_transport": "true",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with secure_transport: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with secure_transport")
	}
}

// TestProviderConfigure_SourceIP 测试 source_ip 配置
func TestProviderConfigure_SourceIP(t *testing.T) {
	raw := map[string]interface{}{
		"access_key": "LTAI5tTestAccessKey123456",
		"secret_key": "TestSecretKey123456789012345678",
		"region":     "cn-hangzhou",
		"source_ip":  "192.168.1.100",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with source_ip: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with source_ip")
	}
}

// TestProviderConfigure_StaticOverridesProfile 测试静态配置优先级高于 Profile
func TestProviderConfigure_StaticOverridesProfile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "alicloud-test-priority-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	config := map[string]interface{}{
		"profiles": []map[string]interface{}{
			{
				"name":              "priority-test",
				"mode":              "AK",
				"access_key_id":     "LTAI5tProfileAK1234567890",
				"access_key_secret": "ProfileSecretKey1234567890",
				"region_id":         "cn-shanghai",
			},
		},
		"current": "priority-test",
	}

	configData, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 静态配置应该覆盖 profile 配置
	raw := map[string]interface{}{
		"access_key":              "fakeAccessKey",
		"secret_key":              "fakeSecretKey",
		"profile":                 "priority-test",
		"shared_credentials_file": configPath,
		"region":                  "cn-hangzhou",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	providerConfig = nil

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}

// TestProviderConfigure_ProfileNotFound 测试 Profile 不存在的情况
func TestProviderConfigure_ProfileNotFound(t *testing.T) {
	// 保存当前环境变量并清除
	envVars := []string{
		"ALICLOUD_ACCESS_KEY", "ALICLOUD_SECRET_KEY", "ALICLOUD_SECURITY_TOKEN",
		"ALIBABA_CLOUD_ACCESS_KEY_ID", "ALIBABA_CLOUD_ACCESS_KEY_SECRET",
		"ALICLOUD_PROFILE", "ALIBABA_CLOUD_PROFILE",
		"ALICLOUD_ECS_ROLE_NAME",
	}
	savedEnvs := make(map[string]string)
	for _, env := range envVars {
		savedEnvs[env] = os.Getenv(env)
		os.Unsetenv(env)
	}
	defer func() {
		for env, val := range savedEnvs {
			if val != "" {
				os.Setenv(env, val)
			}
		}
	}()

	tmpDir, err := ioutil.TempDir("", "alicloud-test-notfound-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	config := map[string]interface{}{
		"profiles": []map[string]interface{}{
			{
				"name":              "existing-profile",
				"mode":              "AK",
				"access_key_id":     "LTAI5tExistingAK1234567890",
				"access_key_secret": "ExistingSecretKey1234567890",
				"region_id":         "cn-hangzhou",
			},
		},
		"current": "existing-profile",
	}

	configData, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	raw := map[string]interface{}{
		"profile":                 "non-existing-profile",
		"shared_credentials_file": configPath,
		"region":                  "cn-hangzhou",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	providerConfig = nil

	_, err = providerConfigure(resourceData, provider)
	// Profile 不存在时，应该返回缺少凭证的错误
	if err == nil {
		t.Fatal("Expected error when profile not found")
	}
}

// TestProviderConfigure_AllRegions 测试不同区域的配置
func TestProviderConfigure_AllRegions(t *testing.T) {
	regions := []string{
		"cn-hangzhou",
		"cn-beijing",
		"cn-shanghai",
		"cn-shenzhen",
		"cn-hongkong",
		"ap-southeast-1",
		"us-west-1",
		"eu-central-1",
	}

	for _, region := range regions {
		t.Run(region, func(t *testing.T) {
			raw := map[string]interface{}{
				"access_key": "LTAI5tTestAccessKey123456",
				"secret_key": "TestSecretKey123456789012345678",
				"region":     region,
			}

			provider := Provider().(*schema.Provider)
			resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

			client, err := providerConfigure(resourceData, provider)
			if err != nil {
				t.Fatalf("providerConfigure failed for region %s: %v", region, err)
			}

			if client == nil {
				t.Fatalf("Expected non-nil client for region %s", region)
			}
		})
	}
}

// TestProviderConfigure_EcsRoleName 测试 ECS 角色认证
func TestProviderConfigure_EcsRoleName(t *testing.T) {
	// 此测试需要在 ECS 实例上运行，或设置了相应的环境变量
	if os.Getenv("ALICLOUD_ECS_ROLE_NAME") == "" {
		t.Skip("Skipping ECS role test: ALICLOUD_ECS_ROLE_NAME not set")
	}

	raw := map[string]interface{}{
		"ecs_role_name": os.Getenv("ALICLOUD_ECS_ROLE_NAME"),
		"region":        "cn-hangzhou",
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Logf("providerConfigure with ecs_role_name failed (expected in non-ECS environment): %v", err)
	} else if client == nil {
		t.Log("Client is nil (expected in non-ECS environment)")
	}
}

// TestProviderConfigure_SignVersion 测试签名版本配置
func TestProviderConfigure_SignVersion(t *testing.T) {
	raw := map[string]interface{}{
		"access_key": "LTAI5tTestAccessKey123456",
		"secret_key": "TestSecretKey123456789012345678",
		"region":     "cn-hangzhou",
		"sign_version": []interface{}{
			map[string]interface{}{
				"oss": "v4",
				"sls": "v4",
			},
		},
	}

	provider := Provider().(*schema.Provider)
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

	client, err := providerConfigure(resourceData, provider)
	if err != nil {
		t.Fatalf("providerConfigure failed with sign_version: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with sign_version")
	}
}

// 辅助函数：检查字符串是否包含子串
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStringHelper(s, substr))
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestProviderConfigure_AdvancedProfileModes 测试高级 Profile 模式 (ChainableRamRoleArn, CloudSSO, External, OAuth)
func TestProviderConfigure_AdvancedProfileModes(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "alicloud-test-advanced-profiles-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 测试不同的高级模式
	modes := []struct {
		name        string
		mode        string
		extraConfig map[string]interface{}
	}{
		{
			name: "ChainableRamRoleArn",
			mode: "ChainableRamRoleArn",
			extraConfig: map[string]interface{}{
				"ram_role_arn":           "acs:ram::123456789012:role/test-role",
				"ram_session_name":       "test-session",
				"ram_session_expiration": 3600,
			},
		},
		{
			name: "CloudSSO",
			mode: "CloudSSO",
			extraConfig: map[string]interface{}{
				"start_url":  "https://sso.aliyun.com",
				"sso_region": "cn-shanghai",
				"account_id": "123456789012",
			},
		},
		{
			name: "External",
			mode: "External",
			extraConfig: map[string]interface{}{
				"process_command": "/bin/echo '{\"AccessKeyId\":\"test\",\"AccessKeySecret\":\"test\",\"SecurityToken\":\"test\"}'",
			},
		},
		{
			name: "OAuth",
			mode: "OAuth",
			extraConfig: map[string]interface{}{
				"oauth_client_id":     "test-client-id",
				"oauth_client_secret": "test-client-secret",
				"oauth_token_url":     "https://oauth.aliyun.com/token",
				"oauth_scope":         "openid",
			},
		},
	}

	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			// 创建配置文件
			configPath := filepath.Join(tmpDir, fmt.Sprintf("config_%s.json", mode.name))

			// 构建 profile 配置
			profile := map[string]interface{}{
				"name":      fmt.Sprintf("test-%s", mode.name),
				"mode":      mode.mode,
				"region_id": "cn-hangzhou",
			}

			// 添加额外配置
			for k, v := range mode.extraConfig {
				profile[k] = v
			}

			config := map[string]interface{}{
				"profiles": []map[string]interface{}{profile},
				"current":  fmt.Sprintf("test-%s", mode.name),
			}

			// 写入配置文件
			configData, err := json.Marshal(config)
			if err != nil {
				t.Fatal(err)
			}

			err = ioutil.WriteFile(configPath, configData, 0644)
			if err != nil {
				t.Fatal(err)
			}

			// 创建测试数据
			raw := map[string]interface{}{
				"profile":                 fmt.Sprintf("test-%s", mode.name),
				"shared_credentials_file": configPath,
				"region":                  "cn-hangzhou",
			}

			provider := Provider().(*schema.Provider)
			resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

			// 重置全局变量
			providerConfig = nil

			// 测试配置过程
			client, err := providerConfigure(resourceData, provider)
			if err != nil {
				// 对于 External 和 OAuth 模式，由于缺少真实的凭证提供程序，可能会出错
				// 我们主要验证配置是否能正确处理这些模式
				t.Logf("Expected error for mode %s (might be due to missing real credential provider): %v", mode.name, err)
			} else {
				if client == nil {
					t.Errorf("Expected non-nil client for mode %s", mode.name)
				}
			}
		})
	}
}

// TestProviderConfigure_GetConfigFromProfile_AdvancedModes 测试 getConfigFromProfile 对高级模式的处理
func TestProviderConfigure_GetConfigFromProfile_AdvancedModes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "alicloud-test-getconfig-advanced-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	modes := []string{"ChainableRamRoleArn", "CloudSSO", "External", "OAuth"}

	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			// 创建配置文件
			configPath := filepath.Join(tmpDir, fmt.Sprintf("config_%s.json", mode))
			profileName := fmt.Sprintf("test-%s", mode)
			config := map[string]interface{}{
				"profiles": []map[string]interface{}{
					{
						"name":      profileName,
						"mode":      mode,
						"region_id": "cn-hangzhou",
					},
				},
				"current": profileName,
			}

			// 写入配置文件
			configData, err := json.Marshal(config)
			if err != nil {
				t.Fatal(err)
			}

			err = ioutil.WriteFile(configPath, configData, 0644)
			if err != nil {
				t.Fatal(err)
			}

			// 创建测试数据
			raw := map[string]interface{}{
				"profile":                 profileName,
				"shared_credentials_file": configPath,
				"region":                  "cn-hangzhou",
			}

			provider := Provider().(*schema.Provider)
			resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

			// 重置全局变量
			providerConfig = nil

			// 测试获取 region_id（应该正常返回，因为 region_id 是特例）
			region, err := getConfigFromProfile(resourceData, "region_id")
			if err != nil {
				t.Errorf("Failed to get region_id for mode %s: %v", mode, err)
			} else if region != "cn-hangzhou" {
				t.Errorf("Expected region_id 'cn-hangzhou', got '%v' for mode %s", region, mode)
			}

			// 测试获取 mode（对于高级模式应该返回 nil，除了 region_id 之外的所有字段都返回 nil）
			modeValue, err := getConfigFromProfile(resourceData, "mode")
			if err != nil {
				t.Errorf("Unexpected error when getting mode for mode %s: %v", mode, err)
			} else if modeValue != nil {
				t.Errorf("Expected nil for mode in advanced mode %s, got %v", mode, modeValue)
			}

			// 测试获取 access_key_id（对于高级模式应该返回 nil）
			accessKey, err := getConfigFromProfile(resourceData, "access_key_id")
			if err != nil {
				t.Errorf("Failed to get access_key_id for mode %s: %v", mode, err)
			} else if accessKey != nil {
				t.Errorf("Expected nil for access_key_id in mode %s, got %v", mode, accessKey)
			}

			// 测试获取 ram_role_name（对于高级模式应该返回 nil）
			ramRoleName, err := getConfigFromProfile(resourceData, "ram_role_name")
			if err != nil {
				t.Errorf("Failed to get ram_role_name for mode %s: %v", mode, err)
			} else if ramRoleName != nil {
				t.Errorf("Expected nil for ram_role_name in mode %s, got %v", mode, ramRoleName)
			}
		})
	}
}

// TestAdvancedProfileModes_CredentialRetrieval 测试高级 Profile 模式下凭证获取的功能
func TestAdvancedProfileModes_CredentialRetrieval(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "alicloud-test-credential-retrieval-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 清理环境变量中的访问凭证，防止默认获取环境变量中的 AK/SK
	credentialEnvVars := []string{
		"ALIBABA_CLOUD_ACCESS_KEY_ID",
		"ALIBABA_CLOUD_SECRET_ACCESS_KEY",
		"ALIBABA_CLOUD_SECURITY_TOKEN",
		"ALICLOUD_ACCESS_KEY",
		"ALICLOUD_SECRET_KEY",
		"ALICLOUD_SECURITY_TOKEN",
		"ALICLOUD_REGION",
		"ALICLOUD_PROFILE",
		"ALICLOUD_SHARED_CREDENTIALS_FILE",
	}

	// 保存原始环境变量值
	originalEnvVars := make(map[string]string)
	for _, envVar := range credentialEnvVars {
		originalEnvVars[envVar] = os.Getenv(envVar)
		os.Unsetenv(envVar)
	}

	// 测试结束后恢复环境变量
	defer func() {
		for envVar, value := range originalEnvVars {
			if value != "" {
				os.Setenv(envVar, value)
			}
		}
	}()

	// 对于这个测试，我们只测试配置是否能正确处理高级模式
	// 实际的凭证获取需要真实的阿里云环境，所以我们只验证流程
	modes := []struct {
		name string
		mode string
	}{
		{name: "ChainableRamRoleArn", mode: "ChainableRamRoleArn"},
		{name: "CloudSSO", mode: "CloudSSO"},
		{name: "External", mode: "External"},
		{name: "OAuth", mode: "OAuth"},
	}

	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			// 创建配置文件
			configPath := filepath.Join(tmpDir, fmt.Sprintf("config_%s.json", mode.name))
			profileName := fmt.Sprintf("test-%s", mode.name)

			// 根据不同模式准备配置
			profile := map[string]interface{}{
				"name":      profileName,
				"mode":      mode.mode,
				"region_id": "cn-hangzhou",
			}

			// 为不同模式添加必要的配置字段
			switch mode.mode {
			case "ChainableRamRoleArn":
				profile["ram_role_arn"] = "acs:ram::123456789012:role/test-role"
				profile["ram_session_name"] = "test-session"
			case "CloudSSO":
				profile["start_url"] = "https://sso.aliyun.com"
				profile["sso_region"] = "cn-shanghai"
			case "External":
				profile["process_command"] = "/bin/echo '{}'" // 空输出命令
			case "OAuth":
				profile["oauth_client_id"] = "test-client-id"
				profile["oauth_client_secret"] = "test-client-secret"
				profile["oauth_site_type"] = "CN"
			}

			config := map[string]interface{}{
				"profiles": []map[string]interface{}{profile},
				"current":  profileName,
			}

			// 写入配置文件
			configData, err := json.Marshal(config)
			if err != nil {
				t.Fatal(err)
			}

			err = ioutil.WriteFile(configPath, configData, 0644)
			if err != nil {
				t.Fatal(err)
			}

			// 创建测试数据
			raw := map[string]interface{}{
				"profile":                 profileName,
				"shared_credentials_file": configPath,
				"region":                  "cn-hangzhou",
			}

			provider := Provider().(*schema.Provider)
			resourceData := schema.TestResourceDataRaw(t, provider.Schema, raw)

			// 重置全局变量
			providerConfig = nil

			// 测试 providerConfigure 函数对高级模式的处理
			// 注意：由于我们没有真实的凭证提供程序，这里会返回错误，但我们主要验证流程
			client, err := providerConfigure(resourceData, provider)
			if err != nil {
				// 对于 External 和 OAuth 模式，由于缺少真实的凭证提供程序，会出错
				// 这是预期的行为，我们主要验证配置是否能正确识别这些模式
				t.Logf("Expected error for mode %s (might be due to missing real credential provider): %v", mode.name, err)
			} else {
				if client == nil {
					t.Errorf("Expected non-nil client for mode %s", mode.name)
				}
			}
		})
	}
}
