package sweep

import (
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestSharedEndpointCacheForRegion(t *testing.T) {
	region := t.Name()
	otherRegion := region + "-other"
	t.Cleanup(func() {
		regionalEndpoints.Delete(region)
		regionalEndpoints.Delete(otherRegion)
	})
	cache := SharedEndpointCacheForRegion(region)
	if cache != SharedEndpointCacheForRegion(region) {
		t.Fatal("same region must reuse its cache")
	}
	otherCache := SharedEndpointCacheForRegion(otherRegion)
	if cache == otherCache {
		t.Fatal("different regions must not share a cache")
	}
	cache.Store("vpc", "vpc.example.invalid")
	if _, ok := otherCache.Load("vpc"); ok {
		t.Fatal("endpoint leaked across regions")
	}

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if SharedEndpointCacheForRegion(region) != cache {
				t.Error("concurrent lookup returned a different cache")
			}
		}()
	}
	wg.Wait()
}

func TestSharedClientForRegionMissingCredentials(t *testing.T) {
	for _, tc := range []struct {
		name, accessKey, secretKey, want string
	}{
		{"both_missing", "", "", "empty ALICLOUD_ACCESS_KEY"},
		{"access_key_missing", "", "test-secret", "empty ALICLOUD_ACCESS_KEY"},
		{"secret_key_missing", "test-access", "", "empty ALICLOUD_SECRET_KEY"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ALICLOUD_ACCESS_KEY", tc.accessKey)
			t.Setenv("ALICLOUD_SECRET_KEY", tc.secretKey)
			client, err := SharedClientForRegion("cn-hangzhou")
			if client != nil || err == nil || err.Error() != tc.want {
				t.Fatalf("got client type %T and error %v; want nil and %q", client, err, tc.want)
			}
		})
	}
}

func setClientTestEnv(t *testing.T) {
	t.Helper()
	t.Setenv("ALICLOUD_ACCESS_KEY", "test-access")
	t.Setenv("ALICLOUD_SECRET_KEY", "test-secret")
	t.Setenv("ALICLOUD_SECURITY_TOKEN", "")
	// Providing both account fields avoids caller-identity API requests.
	t.Setenv("ALICLOUD_ACCOUNT_ID", "1234567890123456")
	t.Setenv("ALICLOUD_ACCOUNT_TYPE", "Domestic")
	t.Setenv("ALIBABA_CLOUD_ACCOUNT_ID", "")
	t.Setenv("ALIBABA_CLOUD_ACCOUNT_TYPE", "")
	t.Setenv("TF_ENDPOINT_PATH", "")
}

func TestSharedClientForRegion(t *testing.T) {
	for _, tc := range []struct {
		name, accountID, fallbackID, accountType, fallbackType, token string
		wantID, wantType                                              string
	}{
		{"primary", "1234567890123456", "", "Domestic", "", "", "1234567890123456", "Domestic"},
		{"fallback", "", "2345678901234567", "", "International", "", "2345678901234567", "International"},
		{"primary_precedence", "1234567890123456", "2345678901234567", "Domestic", "International", "", "1234567890123456", "Domestic"},
		{"sts", "1234567890123456", "", "Domestic", "", "test-token", "1234567890123456", "Domestic"},
		{"padded_sts", "1234567890123456", "", "Domestic", "", "  test-token  ", "1234567890123456", "Domestic"},
		{"whitespace_token", "1234567890123456", "", "Domestic", "", "   ", "1234567890123456", "Domestic"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setClientTestEnv(t)
			t.Setenv("ALICLOUD_ACCOUNT_ID", tc.accountID)
			t.Setenv("ALIBABA_CLOUD_ACCOUNT_ID", tc.fallbackID)
			t.Setenv("ALICLOUD_ACCOUNT_TYPE", tc.accountType)
			t.Setenv("ALIBABA_CLOUD_ACCOUNT_TYPE", tc.fallbackType)
			t.Setenv("ALICLOUD_SECURITY_TOKEN", tc.token)
			raw, err := SharedClientForRegion("cn-hangzhou")
			if err != nil {
				t.Fatal(err)
			}
			client, ok := raw.(*connectivity.AliyunClient)
			if !ok || client == nil {
				t.Fatalf("unexpected client type %T", raw)
			}
			if client.RegionId != "cn-hangzhou" || client.Region != connectivity.Hangzhou {
				t.Fatalf("unexpected client region: %s / %s", client.RegionId, client.Region)
			}
			if client.AccessKey != "test-access" || client.SecretKey != "test-secret" || client.SecurityToken != tc.token {
				t.Fatal("client credentials differ from the supplied test environment")
			}
			// Inspect configuration without adding an accessor or calling cloud APIs.
			config := reflect.ValueOf(client).Elem().FieldByName("config").Elem()
			for field, want := range map[string]string{
				"AccountId": tc.wantID, "AccountType": tc.wantType, "Protocol": "HTTPS",
			} {
				if got := config.FieldByName(field).String(); got != want {
					t.Errorf("%s = %q, want %q", field, got, want)
				}
			}
			if config.FieldByName("Endpoints").Pointer() != reflect.ValueOf(SharedEndpointCacheForRegion("cn-hangzhou")).Pointer() {
				t.Error("client does not use its regional endpoint cache")
			}
			if config.FieldByName("SignVersion").Pointer() != reflect.ValueOf(&signVersion).Pointer() {
				t.Error("client does not share the signing-version cache")
			}
		})
	}

	t.Run("invalid_region", func(t *testing.T) {
		setClientTestEnv(t)
		raw, err := SharedClientForRegion("invalid-region")
		if raw != nil || err == nil || !strings.Contains(err.Error(), "Invalid Alibaba Cloud region: invalid-region") {
			t.Fatalf("unexpected client type %T and error %v", raw, err)
		}
	})
}

func TestSharedClientForRegionWithBackendRegions(t *testing.T) {
	for _, tc := range []struct {
		name      string
		supported bool
		regions   []connectivity.Region
		want      string
	}{
		{"supported_region", true, []connectivity.Region{connectivity.Hangzhou}, "cn-hangzhou"},
		{"prefer_backup", true, []connectivity.Region{connectivity.Beijing, connectivity.APSouthEast1}, "ap-southeast-1"},
		{"first_supported", true, []connectivity.Region{connectivity.Beijing}, "cn-beijing"},
		{"empty_supported", true, nil, "cn-hangzhou"},
		{"not_excluded", false, []connectivity.Region{connectivity.Beijing}, "cn-hangzhou"},
		{"excluded_region", false, []connectivity.Region{connectivity.Hangzhou}, "ap-southeast-1"},
		{"excluded_backup_seen", false, []connectivity.Region{connectivity.APSouthEast1, connectivity.Hangzhou}, "cn-hangzhou"},
		// Preserve the legacy early break: later entries are not examined.
		{"excluded_backup_after_match", false, []connectivity.Region{connectivity.Hangzhou, connectivity.APSouthEast1}, "ap-southeast-1"},
		{"empty_excluded", false, nil, "cn-hangzhou"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setClientTestEnv(t)
			raw, err := SharedClientForRegionWithBackendRegions("cn-hangzhou", tc.supported, tc.regions)
			if err != nil {
				t.Fatal(err)
			}
			client, ok := raw.(*connectivity.AliyunClient)
			if !ok || client == nil {
				t.Fatalf("unexpected client type %T", raw)
			}
			if client.RegionId != tc.want || string(client.Region) != tc.want {
				t.Fatalf("got region %s / %s, want %s", client.RegionId, client.Region, tc.want)
			}
		})
	}
	t.Run("client_error", func(t *testing.T) {
		t.Setenv("ALICLOUD_ACCESS_KEY", "")
		raw, err := SharedClientForRegionWithBackendRegions("cn-hangzhou", true, []connectivity.Region{connectivity.Beijing})
		if raw != nil || err == nil || err.Error() != "empty ALICLOUD_ACCESS_KEY" {
			t.Fatalf("unexpected client type %T and error %v", raw, err)
		}
	})
}
