package alicloud

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// TestGetConfigFromProfile_AbnormalStructure verifies that a syntactically
// valid but structurally abnormal Alibaba Cloud CLI profile file
// (~/.aliyun/config.json) — produced by concurrent writes, truncated writes,
// manual edits or CLI version incompatibility — returns a clear error from
// getConfigFromProfile during provider configure instead of panicking or
// silently degrading to an empty configuration.
//
// It also regression-checks the normal authentication paths (AK, StsToken,
// EcsRamRole, RamRoleArn and the advanced modes) and the historical fallback
// for a missing profile file.
func TestGetConfigFromProfile_AbnormalStructure(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "alicloud-test-profile-panic-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Reuse a single non-existent path for the "file missing" case.
	missingPath := filepath.Join(tmpDir, "does-not-exist.json")

	tests := []struct {
		name        string
		profile     string
		content     string // raw profile file content; ignored when fileMissing is true
		fileMissing bool
		profileKey  string
		wantErr     bool
		wantVal     interface{} // asserted (with comma-ok) when !wantErr and checkVal
		checkVal    bool
	}{
		// Abnormal profile files: every case must return an error and must not
		// panic. These correspond to the acceptance cases from the route plan.
		{name: "empty_file", profile: "p", content: "", profileKey: "access_key_id", wantErr: true},
		{name: "truncated_json", profile: "p", content: `{"profiles":`, profileKey: "access_key_id", wantErr: true},
		{name: "empty_object", profile: "p", content: `{}`, profileKey: "access_key_id", wantErr: true},
		{name: "profiles_is_object", profile: "p", content: `{"profiles":{}}`, profileKey: "access_key_id", wantErr: true},
		{name: "profiles_null_element", profile: "p", content: `{"profiles":[null]}`, profileKey: "access_key_id", wantErr: true},
		{name: "profiles_number_element", profile: "p", content: `{"profiles":[123]}`, profileKey: "access_key_id", wantErr: true},
		{name: "mode_type_error", profile: "p", content: `{"profiles":[{"name":"p","mode":123}]}`, profileKey: "access_key_id", wantErr: true},
		{name: "missing_mode_field", profile: "p", content: `{"profiles":[{"name":"p","region_id":"cn-hangzhou"}]}`, profileKey: "access_key_id", wantErr: true},
		{name: "profile_not_found", profile: "missing", content: `{"profiles":[{"name":"p","mode":"AK"}]}`, profileKey: "access_key_id", wantErr: true},

		// Regression: normal authentication paths must keep working and must
		// not return an error.
		{name: "valid_AK", profile: "p", content: `{"profiles":[{"name":"p","mode":"AK","access_key_id":"AK123","access_key_secret":"SK123","region_id":"cn-hangzhou"}],"current":"p"}`, profileKey: "access_key_id", wantErr: false, wantVal: "AK123", checkVal: true},
		{name: "valid_StsToken", profile: "p", content: `{"profiles":[{"name":"p","mode":"StsToken","sts_token":"STS123","region_id":"cn-hangzhou"}],"current":"p"}`, profileKey: "sts_token", wantErr: false, wantVal: "STS123", checkVal: true},
		{name: "valid_EcsRamRole", profile: "p", content: `{"profiles":[{"name":"p","mode":"EcsRamRole","ram_role_name":"role-x","region_id":"cn-hangzhou"}],"current":"p"}`, profileKey: "ram_role_name", wantErr: false, wantVal: "role-x", checkVal: true},
		{name: "valid_RamRoleArn", profile: "p", content: `{"profiles":[{"name":"p","mode":"RamRoleArn","ram_role_arn":"acs:ram::1:role/x","ram_session_name":"sess","region_id":"cn-hangzhou"}],"current":"p"}`, profileKey: "ram_role_arn", wantErr: false, wantVal: "acs:ram::1:role/x", checkVal: true},
		{name: "valid_advanced_mode_returns_nil", profile: "p", content: `{"profiles":[{"name":"p","mode":"ChainableRamRoleArn","region_id":"cn-hangzhou"}],"current":"p"}`, profileKey: "access_key_id", wantErr: false, wantVal: nil, checkVal: true},
		{name: "region_id_special_case", profile: "p", content: `{"profiles":[{"name":"p","mode":"ChainableRamRoleArn","region_id":"cn-hangzhou"}],"current":"p"}`, profileKey: "region_id", wantErr: false, wantVal: "cn-hangzhou", checkVal: true},

		// Regression: a missing profile file is not a structural defect; the
		// provider must keep the historical fallback (no error, nil value).
		{name: "file_not_exist", profile: "p", fileMissing: true, profileKey: "access_key_id", wantErr: false, wantVal: nil, checkVal: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := missingPath
			if !tt.fileMissing {
				configPath = filepath.Join(tmpDir, tt.name+".json")
				if err := ioutil.WriteFile(configPath, []byte(tt.content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			raw := map[string]interface{}{
				"profile":                 tt.profile,
				"shared_credentials_file": configPath,
				"region":                  "cn-beijing",
			}
			resourceData := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).Schema, raw)

			// Reset the package-level cache so each case loads its own file.
			providerConfig = nil

			// Recover from any panic so a regression is reported as a test
			// failure rather than crashing the test binary.
			var got interface{}
			var err error
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("getConfigFromProfile panicked for %q: %v", tt.name, r)
					}
				}()
				got, err = getConfigFromProfile(resourceData, tt.profileKey)
			}()

			if tt.wantErr && err == nil {
				t.Fatalf("expected an error for %q, got nil (value=%v)", tt.name, got)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.name, err)
			}
			if tt.checkVal && !tt.wantErr {
				if tt.wantVal == nil {
					if got != nil {
						t.Errorf("expected nil value for %q, got %v", tt.name, got)
					}
					return
				}
				s, ok := got.(string)
				if !ok {
					t.Fatalf("expected string value %q for %q, got %T (%v)", tt.wantVal, tt.name, got, got)
				}
				if s != tt.wantVal.(string) {
					t.Errorf("expected %q for %q, got %q", tt.wantVal, tt.name, s)
				}
			}
		})
	}
}
