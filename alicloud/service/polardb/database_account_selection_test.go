package polardb

import (
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func databaseAccountResponse(records ...map[string]interface{}) map[string]interface{} {
	accounts := make([]interface{}, 0, len(records))
	for _, record := range records {
		accounts = append(accounts, record)
	}
	return map[string]interface{}{
		"Accounts": map[string]interface{}{"Account": accounts},
	}
}

// TestUnitAliCloudPolarDbDatabase_accountSelection covers the multi-account
// DescribeDatabases responses consumed by the resource read path. The configured
// account must be selected by name regardless of list order, and an unknown or
// empty account name must never fall back to the first entry: that entry may
// belong to another account, and binding its privilege to the ForceNew
// account_privilege would plan a spurious database replacement.
func TestUnitAliCloudPolarDbDatabase_accountSelection(t *testing.T) {
	accountA := map[string]interface{}{"AccountName": "alice", "AccountPrivilege": "ReadOnly", "PrivilegeStatus": "Empowered"}
	accountB := map[string]interface{}{"AccountName": "bob", "AccountPrivilege": "ReadWrite", "PrivilegeStatus": "Empowered"}

	tests := []struct {
		name          string
		response      map[string]interface{}
		accountName   string
		want          map[string]interface{}
		wantFound     bool
		wantAuthority bool
	}{
		{
			name:          "configured account listed after another account",
			response:      databaseAccountResponse(accountB, accountA),
			accountName:   "alice",
			want:          accountA,
			wantFound:     true,
			wantAuthority: true,
		},
		{
			name:          "configured account listed first",
			response:      databaseAccountResponse(accountA, accountB),
			accountName:   "alice",
			want:          accountA,
			wantFound:     true,
			wantAuthority: true,
		},
		{
			name:          "configured account missing from a populated response",
			response:      databaseAccountResponse(accountB),
			accountName:   "alice",
			wantFound:     false,
			wantAuthority: true,
		},
		{
			name:          "legal empty account list is authoritative",
			response:      databaseAccountResponse(),
			accountName:   "alice",
			wantFound:     false,
			wantAuthority: true,
		},
		{
			name:          "empty accounts container is authoritative",
			response:      map[string]interface{}{"Accounts": map[string]interface{}{}},
			accountName:   "alice",
			wantFound:     false,
			wantAuthority: true,
		},
		{
			name:          "no configured account does not fall back to the first entry",
			response:      databaseAccountResponse(accountB, accountA),
			accountName:   "",
			wantFound:     false,
			wantAuthority: true,
		},
		{
			name:          "missing accounts block is not authoritative",
			response:      map[string]interface{}{},
			accountName:   "alice",
			wantFound:     false,
			wantAuthority: false,
		},
		{
			name:          "malformed account block is not authoritative",
			response:      map[string]interface{}{"Accounts": map[string]interface{}{"Account": "unexpected"}},
			accountName:   "alice",
			wantFound:     false,
			wantAuthority: false,
		},
		{
			name:          "malformed entries make the block non-authoritative",
			response:      map[string]interface{}{"Accounts": map[string]interface{}{"Account": []interface{}{"unexpected", accountA}}},
			accountName:   "alice",
			wantFound:     false,
			wantAuthority: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			records, authoritative := databaseAccounts(test.response)
			if authoritative != test.wantAuthority {
				t.Fatalf("expected authoritative=%v, got %v", test.wantAuthority, authoritative)
			}
			got, found := selectDatabaseAccount(records, test.accountName)
			if found != test.wantFound {
				t.Fatalf("expected found=%v, got %v (account %v)", test.wantFound, found, got)
			}
			if !test.wantFound {
				if got != nil {
					t.Fatalf("expected no account, got %v", got)
				}
				return
			}
			if got == nil {
				t.Fatalf("expected account %v, got nil", test.want)
			}
			if gotName, wantName := got["AccountName"], test.want["AccountName"]; gotName != wantName {
				t.Fatalf("expected account %v, got %v", wantName, gotName)
			}
			if gotPrivilege, wantPrivilege := got["AccountPrivilege"], test.want["AccountPrivilege"]; gotPrivilege != wantPrivilege {
				t.Fatalf("expected privilege %v, got %v", wantPrivilege, gotPrivilege)
			}
		})
	}
}

// TestUnitAliCloudPolarDbDatabase_refreshAccountState covers the state effect of
// the resource read path, not only the account selection helper. A legal empty
// account list proves that the last grant was revoked and must clear the stale
// ForceNew privilege; a missing or malformed block cannot prove that and keeps
// the recorded value.
func TestUnitAliCloudPolarDbDatabase_refreshAccountState(t *testing.T) {
	accountA := map[string]interface{}{"AccountName": "alice", "AccountPrivilege": "ReadOnly"}
	accountB := map[string]interface{}{"AccountName": "bob", "AccountPrivilege": "ReadWrite"}

	tests := []struct {
		name            string
		response        map[string]interface{}
		wantAccountName string
		wantPrivilege   string
	}{
		{
			name:            "legal empty list clears the stale privilege",
			response:        databaseAccountResponse(),
			wantAccountName: "alice",
			wantPrivilege:   "",
		},
		{
			name:            "empty accounts container clears the stale privilege",
			response:        map[string]interface{}{"Accounts": map[string]interface{}{}},
			wantAccountName: "alice",
			wantPrivilege:   "",
		},
		{
			name:            "missing accounts block keeps the recorded privilege",
			response:        map[string]interface{}{},
			wantAccountName: "alice",
			wantPrivilege:   "ReadOnly",
		},
		{
			name:            "malformed account block keeps the recorded privilege",
			response:        map[string]interface{}{"Accounts": map[string]interface{}{"Account": "unexpected"}},
			wantAccountName: "alice",
			wantPrivilege:   "ReadOnly",
		},
		{
			name:            "other accounts only clear the stale privilege",
			response:        databaseAccountResponse(accountB),
			wantAccountName: "alice",
			wantPrivilege:   "",
		},
		{
			name:            "matched account refreshes name and privilege",
			response:        databaseAccountResponse(accountB, accountA),
			wantAccountName: "alice",
			wantPrivilege:   "ReadOnly",
		},
		{
			name:            "matched account without privilege keeps the configured value",
			response:        databaseAccountResponse(map[string]interface{}{"AccountName": "alice", "PrivilegeStatus": "Empowered"}),
			wantAccountName: "alice",
			wantPrivilege:   "ReadOnly",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, ResourceDatabase().Schema, map[string]interface{}{
				"account_name":      "alice",
				"account_privilege": "ReadOnly",
			})
			refreshDatabaseAccountState(d, test.response)
			if got := d.Get("account_name"); got != test.wantAccountName {
				t.Fatalf("expected account_name %q, got %v", test.wantAccountName, got)
			}
			if got := d.Get("account_privilege"); got != test.wantPrivilege {
				t.Fatalf("expected account_privilege %q, got %v", test.wantPrivilege, got)
			}
		})
	}
}

// TestUnitAliCloudPolarDbDatabase_importKeepsAccountStateEmpty proves that an
// imported database carries no account identity and therefore no account
// selection is attempted from the list position.
func TestUnitAliCloudPolarDbDatabase_importKeepsAccountStateEmpty(t *testing.T) {
	accountB := map[string]interface{}{"AccountName": "bob", "AccountPrivilege": "ReadWrite"}
	d := schema.TestResourceDataRaw(t, ResourceDatabase().Schema, map[string]interface{}{})
	refreshDatabaseAccountState(d, databaseAccountResponse(accountB))
	if got := d.Get("account_name"); got != "" {
		t.Fatalf("expected account_name to stay empty, got %v", got)
	}
	if got := d.Get("account_privilege"); got != "" {
		t.Fatalf("expected account_privilege to stay empty, got %v", got)
	}
}

// TestUnitAliCloudPolarDbDatabase_createRequiresAccountName covers the
// apply-time guard: an account_name that is unknown at plan time can resolve to
// an empty string when Create runs. The guard must reject the request before any
// client call, so an empty client is never dereferenced here.
func TestUnitAliCloudPolarDbDatabase_createRequiresAccountName(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceDatabase().Schema, map[string]interface{}{
		"db_cluster_id":     "pc-database-guard-only",
		"db_name":           "tf_testacc_guard",
		"account_privilege": "ReadOnly",
	})

	err := resourceDatabaseCreate(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected the create guard to reject account_privilege without account_name")
	}
	if !strings.Contains(err.Error(), `"account_privilege" requires "account_name"`) {
		t.Fatalf("expected the account_name guard error, got %v", err)
	}
	if d.Id() != "" {
		t.Fatalf("expected no resource id after the rejected create, got %q", d.Id())
	}
}
