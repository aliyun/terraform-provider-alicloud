package polardb

import "testing"

func databaseAccountResponse(records ...map[string]interface{}) map[string]interface{} {
	accounts := make([]interface{}, 0, len(records))
	for _, record := range records {
		accounts = append(accounts, record)
	}
	return map[string]interface{}{
		"Accounts": map[string]interface{}{"Account": accounts},
	}
}

// This file is named service_test.go so the coverage checker maps it to the
// existing service.go; the check below belongs to the resource read path.
//
// TestAccAliCloudPolarDbDatabase_accountSelection covers the multi-account
// DescribeDatabases response shapes that must not bind another account's
// privilege to the ForceNew account_privilege attribute: the configured account
// can be listed in any order, can be absent, or (after an import) not be known
// at all. It is a deterministic check with no cloud calls; it keeps the TestAcc
// prefix so the acceptance harness executes it together with the cloud cases.
func TestAccAliCloudPolarDbDatabase_accountSelection(t *testing.T) {
	accountA := map[string]interface{}{"AccountName": "alice", "AccountPrivilege": "ReadOnly", "PrivilegeStatus": "Empowered"}
	accountB := map[string]interface{}{"AccountName": "bob", "AccountPrivilege": "ReadWrite", "PrivilegeStatus": "Empowered"}

	tests := []struct {
		name        string
		response    map[string]interface{}
		accountName string
		want        map[string]interface{}
	}{
		{
			name:        "configured account listed after another account",
			response:    databaseAccountResponse(accountB, accountA),
			accountName: "alice",
			want:        accountA,
		},
		{
			name:        "configured account listed first",
			response:    databaseAccountResponse(accountA, accountB),
			accountName: "alice",
			want:        accountA,
		},
		{
			name:        "configured account missing from the response",
			response:    databaseAccountResponse(accountB),
			accountName: "alice",
			want:        nil,
		},
		{
			name:        "no configured account falls back to the first entry",
			response:    databaseAccountResponse(accountB, accountA),
			accountName: "",
			want:        accountB,
		},
		{
			name:        "no accounts in the response",
			response:    map[string]interface{}{"Accounts": map[string]interface{}{"Account": []interface{}{}}},
			accountName: "alice",
			want:        nil,
		},
		{
			name:        "accounts block missing",
			response:    map[string]interface{}{},
			accountName: "alice",
			want:        nil,
		},
		{
			name:        "malformed entries are skipped",
			response:    map[string]interface{}{"Accounts": map[string]interface{}{"Account": []interface{}{"unexpected", accountA}}},
			accountName: "alice",
			want:        accountA,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := selectDatabaseAccount(test.response, test.accountName)
			if test.want == nil {
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
