package alicloud

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
)

func writeTestConfig(t *testing.T, path string, profiles []map[string]interface{}, current string) {
	t.Helper()
	data, err := json.Marshal(map[string]interface{}{"profiles": profiles, "current": current})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
}

func writeTestProfileFile(t *testing.T, path string, profile map[string]interface{}) {
	t.Helper()
	writeTestConfig(t, path, []map[string]interface{}{profile}, profile["name"].(string))
}

func providerResourceData(t *testing.T, raw map[string]interface{}) (*schema.Provider, *schema.ResourceData) {
	t.Helper()
	provider := Provider()
	return provider, schema.TestResourceDataRaw(t, provider.Schema, raw)
}

func configureProvider(t *testing.T, raw map[string]interface{}) *connectivity.AliyunClient {
	t.Helper()
	provider, resourceData := providerResourceData(t, raw)
	configured, diags := providerConfigure(context.Background(), resourceData, provider)
	if diags.HasError() {
		t.Fatalf("providerConfigure failed: %v", diags)
	}
	client, ok := configured.(*connectivity.AliyunClient)
	if !ok || client == nil {
		t.Fatalf("expected an AliyunClient, got %#v", configured)
	}
	return client
}

func clearCredentialEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"ALICLOUD_ACCESS_KEY", "ALICLOUD_SECRET_KEY", "ALICLOUD_SECURITY_TOKEN",
		"ALIBABA_CLOUD_ACCESS_KEY_ID", "ALIBABA_CLOUD_ACCESS_KEY_SECRET", "ALIBABA_CLOUD_SECURITY_TOKEN",
		"ALIBABACLOUD_ACCESS_KEY_ID", "ALICLOUD_PROFILE", "ALIBABA_CLOUD_PROFILE",
		"ALICLOUD_ECS_ROLE_NAME", "ALICLOUD_SHARED_CREDENTIALS_FILE",
		"ALIBABA_CLOUD_OIDC_TOKEN", "ALIBABA_CLOUD_OIDC_TOKEN_FILE",
		"ALIBABA_CLOUD_OIDC_PROVIDER_ARN", "ALIBABA_CLOUD_ROLE_ARN", "ALIBABA_CLOUD_ROLE_SESSION_NAME",
	} {
		t.Setenv(key, "")
	}
}

func diagnosticsText(diags diag.Diagnostics) string {
	var text strings.Builder
	for _, d := range diags {
		text.WriteString(d.Summary)
		text.WriteString(" ")
		text.WriteString(d.Detail)
	}
	return text.String()
}

func appendMap(base, extra map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{}, len(base)+len(extra))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range extra {
		merged[k] = v
	}
	return merged
}

func TestProviderCredentials_MultipleProfiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	profiles := []map[string]interface{}{
		{"name": "static-ak", "mode": "AK", "access_key_id": "STATIC_AK"},
		{"name": "external", "mode": "External", "process_command": "/bin/true"},
		{"name": "oauth", "mode": "OAuth", "oauth_client_id": "client-id"},
		{"name": "cloudsso", "mode": "CloudSSO", "start_url": "https://sso.aliyun.com"},
	}
	writeTestConfig(t, path, profiles, "external")

	for _, profile := range profiles {
		t.Run(profile["name"].(string), func(t *testing.T) {
			_, resourceData := providerResourceData(t, map[string]interface{}{
				"profile": profile["name"], "shared_credentials_file": path, "region": "cn-beijing",
			})
			got, err := loadProfileConfig(resourceData)
			if err != nil {
				t.Fatalf("loadProfileConfig failed: %v", err)
			}
			if got["name"] != profile["name"] || got["mode"] != profile["mode"] {
				t.Fatalf("loaded wrong profile: %#v", got)
			}
		})
	}
}

func TestProviderConfigure_GetConfigFromProfile_AdvancedModes(t *testing.T) {
	for _, mode := range []string{"External", "OAuth", "CloudSSO", "ChainableRamRoleArn"} {
		t.Run(mode, func(t *testing.T) {
			name := "test-" + mode
			path := filepath.Join(t.TempDir(), "config.json")
			writeTestProfileFile(t, path, map[string]interface{}{
				"name": name, "mode": mode, "region_id": "cn-hangzhou",
			})
			_, resourceData := providerResourceData(t, map[string]interface{}{
				"profile": name, "shared_credentials_file": path, "region": "cn-hangzhou",
			})
			profile, err := loadProfileConfig(resourceData)
			if err != nil {
				t.Fatalf("loadProfileConfig failed: %v", err)
			}
			if got, err := getConfigFromProfile(profile, "region_id"); err != nil || got != "cn-hangzhou" {
				t.Fatalf("region_id = %v, err = %v", got, err)
			}
			for _, key := range []string{"mode", "access_key_id", "ram_role_name"} {
				if got, err := getConfigFromProfile(profile, key); err != nil || got != nil {
					t.Fatalf("%s = %v, err = %v; want nil", key, got, err)
				}
			}
		})
	}
}

func TestProviderCredentials_MultiAliasIsolation(t *testing.T) {
	dir := t.TempDir()
	pathA := filepath.Join(dir, "a.json")
	pathB := filepath.Join(dir, "b.json")
	writeTestProfileFile(t, pathA, map[string]interface{}{
		"name": "alias-a", "mode": "AK", "access_key_id": "ALIAS_A_AK",
		"access_key_secret": "ALIAS_A_SK", "region_id": "cn-hangzhou",
	})
	writeTestProfileFile(t, pathB, map[string]interface{}{
		"name": "alias-b", "mode": "StsToken", "access_key_id": "ALIAS_B_AK",
		"access_key_secret": "ALIAS_B_SK", "sts_token": "ALIAS_B_STS", "region_id": "cn-beijing",
	})

	clientA := configureProvider(t, map[string]interface{}{
		"profile": "alias-a", "shared_credentials_file": pathA, "region": "cn-hangzhou",
	})
	clientB := configureProvider(t, map[string]interface{}{
		"profile": "alias-b", "shared_credentials_file": pathB, "region": "cn-beijing",
	})
	if clientA.AccessKey != "ALIAS_A_AK" || clientA.SecretKey != "ALIAS_A_SK" || clientA.RegionId != "cn-hangzhou" {
		t.Fatalf("alias A loaded another profile: %#v", clientA)
	}
	if clientB.AccessKey != "ALIAS_B_AK" || clientB.SecretKey != "ALIAS_B_SK" ||
		clientB.SecurityToken != "ALIAS_B_STS" || clientB.RegionId != "cn-beijing" {
		t.Fatalf("alias B loaded another profile: %#v", clientB)
	}
}

func TestProviderCredentials_ConcurrentAliasesNoRace(t *testing.T) {
	const aliasCount = 8
	type aliasProfile struct {
		resourceData *schema.ResourceData
		accessKey    string
		secretKey    string
		region       string
	}

	provider := Provider()
	profiles := make([]aliasProfile, 0, aliasCount)
	for i := 0; i < aliasCount; i++ {
		accessKey := fmt.Sprintf("CONCURRENT_ALIAS_%d_AK", i)
		secretKey := fmt.Sprintf("CONCURRENT_ALIAS_%d_SK", i)
		region := []string{"cn-hangzhou", "cn-beijing", "cn-shanghai", "cn-shenzhen"}[i%4]
		name := fmt.Sprintf("concurrent-alias-%d", i)
		path := filepath.Join(t.TempDir(), "config.json")
		writeTestProfileFile(t, path, map[string]interface{}{
			"name": name, "mode": "AK", "access_key_id": accessKey,
			"access_key_secret": secretKey, "region_id": region,
		})
		profiles = append(profiles, aliasProfile{
			resourceData: schema.TestResourceDataRaw(t, provider.Schema, map[string]interface{}{
				"profile": name, "shared_credentials_file": path, "region": region,
			}),
			accessKey: accessKey, secretKey: secretKey, region: region,
		})
	}

	var wg sync.WaitGroup
	errCh := make(chan error, aliasCount)
	for _, profile := range profiles {
		wg.Add(1)
		go func(profile aliasProfile) {
			defer wg.Done()
			config, err := loadProfileConfig(profile.resourceData)
			if err != nil {
				errCh <- err
				return
			}
			for key, want := range map[string]interface{}{
				"access_key_id": profile.accessKey, "access_key_secret": profile.secretKey, "region_id": profile.region,
			} {
				if got, err := getConfigFromProfile(config, key); err != nil || got != want {
					errCh <- fmt.Errorf("%s = %v, err = %v; want %v", key, got, err, want)
					return
				}
			}
		}(profile)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Error(err)
	}
}

func TestProviderConfigure(t *testing.T) {
	base := map[string]interface{}{
		"access_key": "test-access-key", "secret_key": "test-secret-key", "region": "cn-hangzhou",
	}
	testCases := []struct {
		name       string
		raw        map[string]interface{}
		wantRegion string
		wantToken  string
	}{
		{name: "static AK", raw: base, wantRegion: "cn-hangzhou"},
		{name: "default region", raw: map[string]interface{}{"access_key": "test-access-key", "secret_key": "test-secret-key"}, wantRegion: "cn-beijing"},
		{name: "static STS", raw: appendMap(base, map[string]interface{}{"security_token": "test-token"}), wantRegion: "cn-hangzhou", wantToken: "test-token"},
		{name: "endpoints", raw: appendMap(base, map[string]interface{}{"endpoints": []interface{}{map[string]interface{}{"ecs": "ecs.custom.aliyuncs.com", "vpc": "vpc.custom.aliyuncs.com"}}})},
		{name: "timeouts", raw: appendMap(base, map[string]interface{}{"client_read_timeout": 60, "client_connect_timeout": 30, "max_retry_timeout": 90})},
		{name: "account ID", raw: appendMap(base, map[string]interface{}{"account_id": "123456789012"})},
		{name: "configuration source", raw: appendMap(base, map[string]interface{}{"configuration_source": "test/1.0.0"})},
		{name: "HTTPS", raw: appendMap(base, map[string]interface{}{"protocol": "HTTPS"})},
		{name: "HTTP", raw: appendMap(base, map[string]interface{}{"protocol": "HTTP"})},
		{name: "skip region validation", raw: appendMap(base, map[string]interface{}{"region": "invalid-region", "skip_region_validation": true})},
		{name: "secure transport", raw: appendMap(base, map[string]interface{}{"secure_transport": "true"})},
		{name: "source IP", raw: appendMap(base, map[string]interface{}{"source_ip": "192.168.1.100"})},
		{name: "sign version", raw: appendMap(base, map[string]interface{}{"sign_version": []interface{}{map[string]interface{}{"oss": "v4", "sls": "v4"}}})},
	}
	for _, region := range []string{"cn-beijing", "cn-shanghai", "cn-shenzhen", "cn-hongkong", "ap-southeast-1", "us-west-1", "eu-central-1"} {
		testCases = append(testCases, struct {
			name       string
			raw        map[string]interface{}
			wantRegion string
			wantToken  string
		}{name: "region " + region, raw: appendMap(base, map[string]interface{}{"region": region}), wantRegion: region})
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv("ALICLOUD_REGION", "")
			t.Setenv("ALIBABA_CLOUD_REGION", "")
			client := configureProvider(t, testCase.raw)
			if testCase.wantRegion != "" && client.RegionId != testCase.wantRegion {
				t.Fatalf("RegionId = %q, want %q", client.RegionId, testCase.wantRegion)
			}
			if client.AccessKey != "test-access-key" || client.SecurityToken != testCase.wantToken {
				t.Fatalf("unexpected credentials: access key %q, token %q", client.AccessKey, client.SecurityToken)
			}
		})
	}
}

func TestProviderConfigureCredentialErrors(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		raw         map[string]interface{}
		wantMessage string
	}{
		{name: "missing credentials", raw: map[string]interface{}{"region": "cn-hangzhou"}, wantMessage: "no valid credential sources"},
		{name: "missing secret key", raw: map[string]interface{}{"access_key": "test-access-key", "region": "cn-hangzhou"}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			clearCredentialEnv(t)
			provider, resourceData := providerResourceData(t, testCase.raw)
			_, diags := providerConfigure(context.Background(), resourceData, provider)
			if !diags.HasError() {
				t.Fatal("expected providerConfigure to fail")
			}
			if testCase.wantMessage != "" && !strings.Contains(diagnosticsText(diags), testCase.wantMessage) {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
		})
	}
}

func TestProviderConfigureProfileCredentials(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		mode      string
		token     string
		accessKey string
	}{
		{name: "AK", mode: "AK", accessKey: "PROFILE_AK"},
		{name: "StsToken", mode: "StsToken", accessKey: "PROFILE_STS_AK", token: "PROFILE_STS_TOKEN"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			writeTestProfileFile(t, path, map[string]interface{}{
				"name": "test-profile", "mode": testCase.mode, "access_key_id": testCase.accessKey,
				"access_key_secret": "PROFILE_SECRET", "sts_token": testCase.token, "region_id": "cn-beijing",
			})
			client := configureProvider(t, map[string]interface{}{
				"profile": "test-profile", "shared_credentials_file": path, "region": "cn-beijing",
			})
			if client.AccessKey != testCase.accessKey || client.SecretKey != "PROFILE_SECRET" || client.SecurityToken != testCase.token {
				t.Fatalf("unexpected profile credentials: access key %q, secret %q, token %q", client.AccessKey, client.SecretKey, client.SecurityToken)
			}
		})
	}
}

func TestProviderConfigure_ProfileExternal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell script fixture")
	}
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "credential-provider.sh")
	script := `#!/bin/sh
cat <<EOF
{"mode":"External","access_key_id":"STS.ExternalAK","access_key_secret":"ExternalSecret","sts_token":"ExternalToken","expiration":"2099-12-31T23:59:59Z"}
EOF
`
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "config.json")
	writeTestProfileFile(t, path, map[string]interface{}{
		"name": "external", "mode": "External", "process_command": scriptPath, "region_id": "cn-hangzhou",
	})
	configureProvider(t, map[string]interface{}{
		"profile": "external", "shared_credentials_file": path, "region": "cn-hangzhou",
	})
}

func TestProviderConfigure_StaticOverridesProfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	writeTestProfileFile(t, path, map[string]interface{}{
		"name": "priority-test", "mode": "AK", "access_key_id": "PROFILE_AK",
		"access_key_secret": "PROFILE_SECRET", "region_id": "cn-shanghai",
	})
	client := configureProvider(t, map[string]interface{}{
		"access_key": "STATIC_AK", "secret_key": "STATIC_SECRET", "profile": "priority-test",
		"shared_credentials_file": path, "region": "cn-hangzhou",
	})
	if client.AccessKey != "STATIC_AK" || client.SecretKey != "STATIC_SECRET" {
		t.Fatalf("static credentials did not override profile: %q / %q", client.AccessKey, client.SecretKey)
	}
}

func TestProviderConfigure_ProfileNotFound(t *testing.T) {
	clearCredentialEnv(t)
	path := filepath.Join(t.TempDir(), "config.json")
	writeTestProfileFile(t, path, map[string]interface{}{
		"name": "existing", "mode": "AK", "access_key_id": "AK", "access_key_secret": "SK",
	})
	provider, resourceData := providerResourceData(t, map[string]interface{}{
		"profile": "missing", "shared_credentials_file": path, "region": "cn-hangzhou",
	})
	if _, diags := providerConfigure(context.Background(), resourceData, provider); !diags.HasError() {
		t.Fatal("expected a missing-profile error")
	}
}

func TestProviderConfigure_AssumeRoleConfig(t *testing.T) {
	_, resourceData := providerResourceData(t, map[string]interface{}{
		"assume_role": []interface{}{map[string]interface{}{
			"role_arn": "acs:ram::123456789012:role/test-role", "session_name": "test-session", "session_expiration": 3600,
		}},
	})
	assumeRoles := resourceData.Get("assume_role").(*schema.Set).List()
	if len(assumeRoles) != 1 {
		t.Fatalf("got %d assume_role blocks", len(assumeRoles))
	}
	assumeRole := assumeRoles[0].(map[string]interface{})
	if assumeRole["role_arn"] != "acs:ram::123456789012:role/test-role" || assumeRole["session_name"] != "test-session" || assumeRole["session_expiration"] != 3600 {
		t.Fatalf("unexpected assume_role config: %#v", assumeRole)
	}
}

func TestProviderConfigure_CredentialsURIConsumed(t *testing.T) {
	t.Run("consumed", func(t *testing.T) {
		clearCredentialEnv(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"Code":"Success","AccessKeyId":"STS.ConfigureAccessKey","AccessKeySecret":"configureSecret","SecurityToken":"configureToken","Expiration":"2099-12-31T23:59:59Z"}`)
		}))
		defer server.Close()
		client := configureProvider(t, map[string]interface{}{"credentials_uri": server.URL, "region": "cn-hangzhou"})
		if client.AccessKey != "STS.ConfigureAccessKey" || client.SecurityToken != "configureToken" {
			t.Fatalf("sidecar credentials not consumed: %q / %q", client.AccessKey, client.SecurityToken)
		}
	})

	t.Run("invalid response", func(t *testing.T) {
		clearCredentialEnv(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}))
		defer server.Close()
		provider, resourceData := providerResourceData(t, map[string]interface{}{"credentials_uri": server.URL, "region": "cn-hangzhou"})
		_, diags := providerConfigure(context.Background(), resourceData, provider)
		if !diags.HasError() {
			t.Fatal("expected invalid credentials_uri response to fail")
		}
		msg := diagnosticsText(diags)
		if !strings.Contains(msg, "failed") && !strings.Contains(msg, "no valid credential") {
			t.Fatalf("expected a credentials-related error, got: %v", diags)
		}
	})
}

func TestProviderCredentials_CredentialsURIRefresh(t *testing.T) {
	var hits int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&hits, 1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"Code":"Success","AccessKeyId":"STS.RefreshAccessKey%d","AccessKeySecret":"refreshSecret%d","SecurityToken":"refreshToken%d","Expiration":%q}`,
			n, n, n, time.Now().Add(time.Second).UTC().Format(time.RFC3339))
	}))
	defer server.Close()

	credential := newCredentialsURICredential(server.URL, time.Now().Add(time.Hour).UTC().Format(time.RFC3339), 0, 0)
	if credential == nil {
		t.Fatal("expected a refreshable credential")
	}
	first, err := credential.GetCredential()
	if err != nil {
		t.Fatal(err)
	}
	second, err := credential.GetCredential()
	if err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&hits); got != 2 {
		t.Fatalf("credentials URI requested %d times, want 2", got)
	}
	if *first.SecurityToken == *second.SecurityToken || *second.Type != "credentials_uri" {
		t.Fatalf("credential was not refreshed correctly: %#v / %#v", first, second)
	}
}

func TestProviderCredentials_CredentialsURIStaticFallback(t *testing.T) {
	for _, expiration := range []string{"", "not-a-timestamp"} {
		if credential := newCredentialsURICredential("http://127.0.0.1:1/credentials", expiration, 0, 0); credential != nil {
			t.Fatalf("expected nil credential for expiration %q", expiration)
		}
	}
}

func TestDefaultConfigPath(t *testing.T) {
	origDisableCache := homedir.DisableCache
	homedir.DisableCache = true
	t.Cleanup(func() { homedir.DisableCache = origDisableCache; homedir.Reset() })

	for _, testCase := range []struct {
		name string
		home string
	}{{name: "explicit home", home: t.TempDir()}, {name: "empty home", home: ""}} {
		t.Run(testCase.name, func(t *testing.T) {
			homedir.Reset()
			t.Setenv("HOME", testCase.home)
			t.Setenv("USERPROFILE", testCase.home)
			got, err := defaultConfigPath()
			if strings.HasPrefix(got, "/.aliyun") {
				t.Fatalf("defaultConfigPath regressed to disk root: %q", got)
			}
			if err != nil {
				if got != "" {
					t.Fatalf("path %q returned with error %v", got, err)
				}
				return
			}
			if !strings.HasSuffix(got, filepath.Join(".aliyun", "config.json")) {
				t.Fatalf("unexpected default config path %q", got)
			}
			if testCase.home != "" && got != filepath.Join(testCase.home, ".aliyun", "config.json") {
				t.Fatalf("defaultConfigPath = %q", got)
			}
		})
	}
}

func TestGetConfigFromProfile_EmptySharedCredentialsFile(t *testing.T) {
	origDisableCache := homedir.DisableCache
	homedir.DisableCache = true
	t.Cleanup(func() { homedir.DisableCache = origDisableCache; homedir.Reset() })
	home := t.TempDir()
	homedir.Reset()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	for _, testCase := range []struct {
		name string
		path string
	}{
		{name: "default path"},
		{name: "explicit missing path", path: filepath.Join(home, "missing.json")},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			_, resourceData := providerResourceData(t, map[string]interface{}{
				"profile": "missing", "shared_credentials_file": testCase.path, "region": "cn-beijing",
			})
			profile, err := loadProfileConfig(resourceData)
			if err != nil {
				t.Fatal(err)
			}
			if got, err := getConfigFromProfile(profile, "access_key_id"); err != nil || got != nil {
				t.Fatalf("access_key_id = %v, err = %v; want nil", got, err)
			}
		})
	}
}

func TestUnitGetConfigFromProfile_InvalidJSON(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		content     string
		wantMessage string
	}{
		{name: "truncated JSON", content: `{"profiles":[{"name":"p1","mode":"AK"`, wantMessage: "invalid JSON"},
		{name: "empty file", content: ``, wantMessage: "invalid JSON"},
		{name: "missing profiles", content: `{"current":"p1"}`},
		{name: "profiles not array", content: `{"profiles":"bad","current":"p1"}`},
		{name: "profile not object", content: `{"profiles":[42],"current":"p1"}`},
		{name: "mode not string", content: `{"profiles":[{"name":"p1","mode":123}],"current":"p1"}`},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			if err := os.WriteFile(path, []byte(testCase.content), 0644); err != nil {
				t.Fatal(err)
			}
			_, resourceData := providerResourceData(t, map[string]interface{}{
				"profile": "p1", "shared_credentials_file": path, "region": "cn-beijing",
			})
			profile, err := loadProfileConfig(resourceData)
			if err == nil {
				t.Fatalf("expected malformed config to fail, got %#v", profile)
			}
			if !strings.Contains(err.Error(), path) || !strings.Contains(err.Error(), "p1") || !strings.Contains(err.Error(), testCase.wantMessage) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestUnitGetConfigFromProfile_RetryAfterParseFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"profiles":[{"name":"retry"`), 0644); err != nil {
		t.Fatal(err)
	}
	_, resourceData := providerResourceData(t, map[string]interface{}{
		"profile": "retry", "shared_credentials_file": path, "region": "cn-beijing",
	})
	if _, err := loadProfileConfig(resourceData); err == nil {
		t.Fatal("expected initial parse failure")
	}
	writeTestProfileFile(t, path, map[string]interface{}{
		"name": "retry", "mode": "AK", "access_key_id": "RETRY_AK", "access_key_secret": "RETRY_SK",
	})
	profile, err := loadProfileConfig(resourceData)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := getConfigFromProfile(profile, "access_key_id"); err != nil || got != "RETRY_AK" {
		t.Fatalf("access_key_id = %v, err = %v", got, err)
	}
}

func TestUnitGetConfigFromProfile_ProfileModeKeySemantics(t *testing.T) {
	configs := map[string]map[string]interface{}{
		"AK":         {"name": "sem-ak", "mode": "AK", "access_key_id": "AK_ID", "access_key_secret": "AK_SECRET", "sts_token": "STS_TOK", "ram_role_name": "ROLE_NAME"},
		"EcsRamRole": {"name": "sem-ecs", "mode": "EcsRamRole", "ram_role_name": "ECS_ROLE"},
		"StsToken":   {"name": "sem-sts", "mode": "StsToken", "access_key_id": "STS_AK_ID", "access_key_secret": "STS_SK", "sts_token": "STS_TOK"},
		"RamRoleArn": {"name": "sem-rra", "mode": "RamRoleArn", "access_key_id": "RRA_AK_ID", "access_key_secret": "RRA_SK", "ram_role_arn": "acs:ram::1234567890123456:role/tf", "ram_session_name": "tf-session", "expired_seconds": 900},
	}
	expectations := map[string]map[string]interface{}{
		"AK":         {"access_key_id": "AK_ID", "sts_token": "", "ram_role_name": ""},
		"EcsRamRole": {"access_key_id": "", "ram_role_name": "ECS_ROLE", "sts_token": ""},
		"StsToken":   {"sts_token": "STS_TOK", "access_key_id": "STS_AK_ID", "ram_role_name": ""},
		"RamRoleArn": {"ram_role_arn": "acs:ram::1234567890123456:role/tf", "ram_session_name": "tf-session", "expired_seconds": float64(900), "sts_token": ""},
	}
	for mode, profile := range configs {
		t.Run(mode, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			writeTestProfileFile(t, path, profile)
			_, resourceData := providerResourceData(t, map[string]interface{}{
				"profile": profile["name"], "shared_credentials_file": path, "region": "cn-beijing",
			})
			loaded, err := loadProfileConfig(resourceData)
			if err != nil {
				t.Fatal(err)
			}
			for key, want := range expectations[mode] {
				if got, err := getConfigFromProfile(loaded, key); err != nil || got != want {
					t.Errorf("%s = %v, err = %v; want %v", key, got, err, want)
				}
			}
		})
	}
}

func TestUnitProviderConfigure_InvalidProfileJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"profiles":[{"name":"broken"`), 0644); err != nil {
		t.Fatal(err)
	}
	provider, resourceData := providerResourceData(t, map[string]interface{}{
		"profile": "broken", "shared_credentials_file": path, "region": "cn-hangzhou",
	})
	_, diags := providerConfigure(context.Background(), resourceData, provider)
	message := diagnosticsText(diags)
	if !diags.HasError() || !strings.Contains(message, path) || !strings.Contains(message, "broken") || !strings.Contains(message, "invalid JSON") {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if strings.Contains(message, "no valid credential sources") {
		t.Fatalf("parse error was masked: %v", diags)
	}
}

func TestUnitProviderConfigure_ProfileAccessKeyNotString(t *testing.T) {
	for _, testCase := range []struct {
		name          string
		profile       map[string]interface{}
		wantAttribute string
		wantMessage   string
	}{
		{
			name: "access key not string",
			profile: map[string]interface{}{
				"name": "typed-ak", "mode": "AK", "access_key_id": 123, "access_key_secret": "SK",
			},
			wantAttribute: "access_key_id",
			wantMessage:   "must be a string",
		},
		{
			name: "expired seconds not number",
			profile: map[string]interface{}{
				"name": "typed-exp", "mode": "RamRoleArn", "access_key_id": "AK", "access_key_secret": "SK",
				"ram_role_arn": "acs:ram::1234567890123456:role/tf", "ram_session_name": "tf-session", "expired_seconds": "900",
			},
			wantAttribute: "expired_seconds",
			wantMessage:   "must be a number",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			clearCredentialEnv(t)
			path := filepath.Join(t.TempDir(), "config.json")
			writeTestProfileFile(t, path, testCase.profile)
			provider, resourceData := providerResourceData(t, map[string]interface{}{
				"profile": testCase.profile["name"], "shared_credentials_file": path, "region": "cn-hangzhou",
			})
			_, diags := providerConfigure(context.Background(), resourceData, provider)
			message := diagnosticsText(diags)
			if !diags.HasError() || !strings.Contains(message, testCase.profile["name"].(string)) ||
				!strings.Contains(message, testCase.wantAttribute) || !strings.Contains(message, testCase.wantMessage) {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
		})
	}
}

func TestProviderConfigure_AssumeRoleWithOIDCConfigParse(t *testing.T) {
	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "oidc-token")
	emptyFile := filepath.Join(dir, "empty-token")
	missingFile := filepath.Join(dir, "missing-token")
	if err := os.WriteFile(tokenFile, []byte("jwt-token-file"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(emptyFile, nil, 0644); err != nil {
		t.Fatal(err)
	}
	base := map[string]interface{}{
		"oidc_provider_arn": "acs:ram::123456789012:oidc-provider/test-idp",
		"role_arn":          "acs:ram::123456789012:role/test-role",
	}

	cfg, err := getAssumeRoleWithOIDCConfig(appendMap(base, map[string]interface{}{"oidc_token_file": tokenFile}))
	if err != nil || cfg.OIDCToken != "" || cfg.OIDCTokenFile != tokenFile {
		t.Fatalf("file form: cfg = %#v, err = %v", cfg, err)
	}
	cfg, err = getAssumeRoleWithOIDCConfig(appendMap(base, map[string]interface{}{
		"oidc_token": "jwt-token-inline", "oidc_token_file": missingFile,
	}))
	if err != nil || cfg.OIDCToken != "jwt-token-inline" || cfg.OIDCTokenFile != missingFile {
		t.Fatalf("inline precedence: cfg = %#v, err = %v", cfg, err)
	}
	for _, testCase := range []struct {
		name        string
		config      map[string]interface{}
		wantMessage string
	}{
		{name: "missing token", config: base, wantMessage: "must be specified"},
		{name: "unreadable file", config: appendMap(base, map[string]interface{}{"oidc_token_file": missingFile}), wantMessage: "reading oidc_token_file failed"},
		{name: "empty file", config: appendMap(base, map[string]interface{}{"oidc_token_file": emptyFile}), wantMessage: "content can not be empty"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			if _, err := getAssumeRoleWithOIDCConfig(testCase.config); err == nil || !strings.Contains(err.Error(), testCase.wantMessage) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestProviderConfigure_AssumeRoleWithOIDCTokenFileConsumed(t *testing.T) {
	clearCredentialEnv(t)
	tokenFile := filepath.Join(t.TempDir(), "oidc-token")
	if err := os.WriteFile(tokenFile, []byte("jwt-token-1"), 0644); err != nil {
		t.Fatal(err)
	}

	var hits int32
	var mu sync.Mutex
	var durations []string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		n := atomic.AddInt32(&hits, 1)
		mu.Lock()
		durations = append(durations, r.FormValue("DurationSeconds"))
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"Credentials":{"AccessKeyId":"STS.OIDCAccessKey%d","AccessKeySecret":"oidcSecret%d","SecurityToken":"oidcToken%d","Expiration":%q}}`,
			n, n, n, time.Now().Add(time.Hour).UTC().Format(time.RFC3339))
	}))
	defer server.Close()

	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		t.Fatalf("unexpected default transport %T", http.DefaultTransport)
	}
	trusted := transport.Clone()
	pool := x509.NewCertPool()
	pool.AddCert(server.Certificate())
	trusted.TLSClientConfig = &tls.Config{RootCAs: pool}
	saved := http.DefaultTransport
	http.DefaultTransport = trusted
	defer func() { http.DefaultTransport = saved }()

	client := configureProvider(t, map[string]interface{}{
		"region": "cn-hangzhou",
		"endpoints": []interface{}{map[string]interface{}{
			"sts": strings.TrimPrefix(server.URL, "https://"),
		}},
		"assume_role_with_oidc": []interface{}{map[string]interface{}{
			"role_arn": "acs:ram::123456789012:role/test-role", "oidc_provider_arn": "acs:ram::123456789012:oidc-provider/test-idp",
			"role_session_name": "test-session", "oidc_token_file": tokenFile, "session_expiration": 1800,
		}},
	})
	if client.AccessKey != "STS.OIDCAccessKey1" || client.SecurityToken != "oidcToken1" {
		t.Fatalf("OIDC credentials not consumed: %q / %q", client.AccessKey, client.SecurityToken)
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Fatalf("STS endpoint hit %d times, want 1", got)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(durations) != 1 || durations[0] != "1800" {
		t.Fatalf("DurationSeconds = %v, want [1800]", durations)
	}
}
