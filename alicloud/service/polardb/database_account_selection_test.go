package polardb

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
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

func databaseAccountResponseForEngine(engine string, records ...map[string]interface{}) map[string]interface{} {
	response := databaseAccountResponse(records...)
	response["Engine"] = engine
	return response
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
		config          map[string]interface{}
		response        map[string]interface{}
		wantAccountName string
		wantPrivilege   string
	}{
		{
			name:            "legal empty list clears the stale identity and privilege",
			response:        databaseAccountResponse(),
			wantAccountName: "",
			wantPrivilege:   "",
		},
		{
			name:            "empty accounts container clears the stale identity and privilege",
			response:        map[string]interface{}{"Accounts": map[string]interface{}{}},
			wantAccountName: "",
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
			name:            "other accounts only clear the stale identity and privilege",
			response:        databaseAccountResponse(accountB),
			wantAccountName: "",
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
		{
			// Legacy PostgreSQL configuration that only manages account_name: the
			// owner is transferred to another account outside Terraform, so the
			// managed binding must become visible even though the privilege field
			// was never configured.
			name:            "postgresql owner transfer clears the managed account name",
			config:          map[string]interface{}{"account_name": "alice"},
			response:        databaseAccountResponseForEngine("PostgreSQL", accountB),
			wantAccountName: "",
			wantPrivilege:   "",
		},
		{
			name:            "oracle owner transfer clears the managed account name",
			config:          map[string]interface{}{"account_name": "alice"},
			response:        databaseAccountResponseForEngine("Oracle", accountB),
			wantAccountName: "",
			wantPrivilege:   "",
		},
		{
			name:            "postgresql empty owner list clears the managed account name",
			config:          map[string]interface{}{"account_name": "alice"},
			response:        databaseAccountResponseForEngine("PostgreSQL"),
			wantAccountName: "",
			wantPrivilege:   "",
		},
		{
			name:            "postgresql response containing the configured owner keeps it",
			config:          map[string]interface{}{"account_name": "alice"},
			response:        databaseAccountResponseForEngine("PostgreSQL", map[string]interface{}{"AccountName": "alice", "PrivilegeStatus": "Empowered"}),
			wantAccountName: "alice",
			wantPrivilege:   "",
		},
		{
			// The response is authoritative and no longer contains the managed
			// account: clearing the name exposes the lost binding even when the
			// privilege field was never configured (the legacy configuration).
			name:            "mysql clears the account name when the managed grant disappears",
			config:          map[string]interface{}{"account_name": "alice"},
			response:        databaseAccountResponseForEngine("MySQL", accountB),
			wantAccountName: "",
			wantPrivilege:   "",
		},
		{
			name:            "non-owner engine clears the account name on an empty list",
			config:          map[string]interface{}{"account_name": "alice"},
			response:        databaseAccountResponseForEngine("MySQL"),
			wantAccountName: "",
			wantPrivilege:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := test.config
			if config == nil {
				config = map[string]interface{}{
					"account_name":      "alice",
					"account_privilege": "ReadOnly",
				}
			}
			d := schema.TestResourceDataRaw(t, ResourceDatabase().Schema, config)
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

// TestUnitAliCloudPolarDbDatabase_createKeepsAccountStateDuringRead covers the
// create-time read: a create can observe an eventually consistent account list,
// so the configured identity and privilege must survive until a later refresh
// can prove the removal.
func TestUnitAliCloudPolarDbDatabase_createKeepsAccountStateDuringRead(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceDatabase().Schema, map[string]interface{}{
		"account_name":      "alice",
		"account_privilege": "ReadOnly",
	})
	d.MarkNewResource()
	refreshDatabaseAccountState(d, databaseAccountResponseForEngine("PostgreSQL"))
	if got := d.Get("account_name"); got != "alice" {
		t.Fatalf("expected account_name to keep the configured value during create, got %v", got)
	}
	if got := d.Get("account_privilege"); got != "ReadOnly" {
		t.Fatalf("expected account_privilege to keep the configured value during create, got %v", got)
	}
}

// TestUnitAliCloudPolarDbDatabase_accountNameChange covers the guard that
// rejects an account_name change on engines whose update does not re-grant the
// account. Accepting such a change would send no request, and the read would
// bind the new account's privilege to the ForceNew account_privilege after the
// plan, so every change on these engines is rejected.
func TestUnitAliCloudPolarDbDatabase_accountNameChange(t *testing.T) {
	tests := []struct {
		name     string
		response map[string]interface{}
		wantErr  string
	}{
		{name: "postgresql is supported", response: map[string]interface{}{"Engine": "PostgreSQL"}},
		{name: "oracle is supported", response: map[string]interface{}{"Engine": "Oracle"}},
		{
			name:     "mysql is rejected",
			response: databaseAccountResponseForEngine("MySQL", map[string]interface{}{"AccountName": "alice", "AccountPrivilege": "ReadOnly"}),
			wantErr:  `"account_name" can only be changed`,
		},
		{
			name:     "mysql with a malformed account block is rejected",
			response: map[string]interface{}{"Engine": "MySQL", "Accounts": map[string]interface{}{"Account": "unexpected"}},
			wantErr:  `"account_name" can only be changed`,
		},
		{
			name:     "missing engine is rejected",
			response: map[string]interface{}{},
			wantErr:  `the database engine is not reported`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateAccountNameChange(test.response)
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("expected the engine to be accepted, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected an error containing %q, got nil", test.wantErr)
			}
			if !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("expected an error containing %q, got %v", test.wantErr, err)
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

// TestUnitAliCloudPolarDbDatabase_waitAccountEmpowered covers the post-grant
// wait: it must only return once the account is listed with
// PrivilegeStatus=Empowered, keep waiting while the status is transitional or
// the account list is not yet visible, and surface other errors immediately.
func TestUnitAliCloudPolarDbDatabase_waitAccountEmpowered(t *testing.T) {
	empowered := func(status string) map[string]interface{} {
		return databaseAccountResponseForEngine("PostgreSQL", map[string]interface{}{
			"AccountName":     "alice",
			"PrivilegeStatus": status,
		})
	}

	t.Run("waits through a missing list and Empowering", func(t *testing.T) {
		responses := []map[string]interface{}{
			databaseAccountResponseForEngine("PostgreSQL"),
			empowered("Empowering"),
			empowered("Empowered"),
		}
		calls := 0
		describe := func(string) (map[string]interface{}, error) {
			if calls >= len(responses) {
				return empowered("Empowered"), nil
			}
			response := responses[calls]
			calls++
			return response, nil
		}
		if err := waitDatabaseAccountEmpowered(describe, "pc-1:db", "alice", 60*time.Second); err != nil {
			t.Fatalf("expected the wait to finish once the account is empowered, got %v", err)
		}
		if calls != len(responses) {
			t.Fatalf("expected %d describe calls, got %d", len(responses), calls)
		}
	})

	t.Run("times out while the account stays Empowering", func(t *testing.T) {
		describe := func(string) (map[string]interface{}, error) {
			return empowered("Empowering"), nil
		}
		if err := waitDatabaseAccountEmpowered(describe, "pc-1:db", "alice", 100*time.Millisecond); err == nil {
			t.Fatal("expected a timeout while the account stays Empowering")
		}
	})

	t.Run("propagates other errors", func(t *testing.T) {
		describe := func(string) (map[string]interface{}, error) {
			return nil, fmt.Errorf("unexpected describe failure")
		}
		err := waitDatabaseAccountEmpowered(describe, "pc-1:db", "alice", 5*time.Second)
		if err == nil || !strings.Contains(err.Error(), "unexpected describe failure") {
			t.Fatalf("expected the describe error to propagate, got %v", err)
		}
	})

	t.Run("retries a malformed account list", func(t *testing.T) {
		responses := []map[string]interface{}{
			{"Engine": "PostgreSQL", "Accounts": map[string]interface{}{"Account": "unexpected"}},
			empowered("Empowered"),
		}
		calls := 0
		describe := func(string) (map[string]interface{}, error) {
			index := calls
			if calls < len(responses)-1 {
				calls++
			}
			return responses[index], nil
		}
		if err := waitDatabaseAccountEmpowered(describe, "pc-1:db", "alice", 60*time.Second); err != nil {
			t.Fatalf("expected the wait to retry the malformed account list, got %v", err)
		}
	})

	t.Run("retries a missing database", func(t *testing.T) {
		responses := []map[string]interface{}{nil, empowered("Empowered")}
		calls := 0
		describe := func(string) (map[string]interface{}, error) {
			index := calls
			if calls < len(responses)-1 {
				calls++
			}
			if responses[index] == nil {
				return nil, errs.GetNotFoundErrorFromString("the database does not exist")
			}
			return responses[index], nil
		}
		if err := waitDatabaseAccountEmpowered(describe, "pc-1:db", "alice", 60*time.Second); err != nil {
			t.Fatalf("expected the wait to retry the missing database, got %v", err)
		}
	})

	t.Run("retries while the privilege is removed", func(t *testing.T) {
		responses := []map[string]interface{}{empowered("Removing"), empowered("Empowered")}
		calls := 0
		describe := func(string) (map[string]interface{}, error) {
			index := calls
			if calls < len(responses)-1 {
				calls++
			}
			return responses[index], nil
		}
		if err := waitDatabaseAccountEmpowered(describe, "pc-1:db", "alice", 60*time.Second); err != nil {
			t.Fatalf("expected the wait to retry the removing privilege, got %v", err)
		}
	})

	t.Run("empty account name is a no-op", func(t *testing.T) {
		describe := func(string) (map[string]interface{}, error) {
			t.Fatal("must not describe the database for an empty account name")
			return nil, nil
		}
		if err := waitDatabaseAccountEmpowered(describe, "pc-1:db", "  ", time.Second); err != nil {
			t.Fatalf("expected no wait for an empty account name, got %v", err)
		}
	})
}
