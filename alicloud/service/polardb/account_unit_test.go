package polardb

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestUnitAlicloudPolarDbAccountPrivilegeValidator(t *testing.T) {
	valid900 := strings.TrimSuffix(strings.Repeat("ReadOnly,", 99), ",") + ",ReadWrite"
	valid901 := strings.Replace(valid900, "ReadOnly", "ReadWrite", 1)
	cases := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"single value", "ReadWrite", false},
		{"multiple values", "ReadWrite,ReadOnly", false},
		{"all values", "ReadWrite,ReadOnly,DMLOnly,DDLOnly,ReadIndex", false},
		{"empty value", "", true},
		{"trailing empty item", "ReadWrite,", true},
		{"empty item in list", "ReadWrite,,ReadOnly", true},
		{"unknown value", "ReadWrite,Unknown", true},
		{"900 characters", valid900, false},
		{"901 characters", valid901, true},
		{"non string", 1, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, errs := validatePolarDbAccountPrivilege(tc.value, "account_privilege")
			if tc.wantErr != (len(errs) > 0) {
				t.Fatalf("value %q: wantErr=%v, got %v", tc.name, tc.wantErr, errs)
			}
		})
	}
}

func TestUnitAlicloudPolarDbAccountCreatePairing(t *testing.T) {
	cases := []struct {
		name         string
		dbName       string
		dbNameKnown  bool
		privilege    string
		privKnown    bool
		privForAllDB string
		privAllKnown bool
		wantErr      bool
	}{
		{"privilege with database", "db", true, "ReadWrite", true, "", true, false},
		{"neither set", "", true, "", true, "", true, false},
		{"empty database with privilege", "", true, "ReadWrite", true, "", true, true},
		{"unknown database defers", "", false, "ReadWrite", true, "", true, false},
		{"unknown privilege defers", "", true, "ReadWrite", false, "", true, false},
		{"priv for all without privilege", "db", true, "", true, "1", true, true},
		{"priv for all with unknown privilege defers", "db", true, "", false, "1", true, false},
		{"priv for all with privilege", "db", true, "ReadWrite", true, "1", true, false},
		{"priv for all disabled", "db", true, "", true, "0", true, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateAccountCreatePairing(tc.dbName, tc.dbNameKnown, tc.privilege, tc.privKnown, tc.privForAllDB, tc.privAllKnown)
			if tc.wantErr != (err != nil) {
				t.Fatalf("%s: wantErr=%v, got %v", tc.name, tc.wantErr, err)
			}
		})
	}
}

func TestUnitAlicloudPolarDbAccountTypeDiffSuppress(t *testing.T) {
	cases := []struct {
		old, new string
		want     bool
	}{
		{"ReadOnly", "Normal", true},
		{"Normal", "ReadOnly", false},
		{"Super", "Normal", false},
		{"ReadOnly", "Super", false},
		{"", "Normal", false},
	}
	for _, tc := range cases {
		if got := polarDbAccountTypeDiffSuppressFunc("", tc.old, tc.new, nil); got != tc.want {
			t.Fatalf("old=%q new=%q: want %v, got %v", tc.old, tc.new, tc.want, got)
		}
	}
}

func TestUnitAlicloudPolarDbAccountPasswordSourceChanged(t *testing.T) {
	cases := []struct {
		name             string
		passwordChanged  bool
		password         string
		kmsChanged       bool
		kmsPassword      string
		woVersionChanged bool
		woPassword       string
		want             bool
	}{
		{"password changed", true, "pw", false, "", false, "", true},
		{"password removed", true, "", false, "", false, "", false},
		{"password type only", false, "pw", false, "", false, "", false},
		{"kms changed", false, "", true, "cipher", false, "", true},
		{"kms removed", false, "", true, "", false, "", false},
		{"write-only version changed", false, "", false, "", true, "wo", true},
		{"write-only version without value", false, "", false, "", true, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := accountPasswordSourceChanged(tc.passwordChanged, tc.password, tc.kmsChanged, tc.kmsPassword, tc.woVersionChanged, tc.woPassword)
			if got != tc.want {
				t.Fatalf("%s: want %v, got %v", tc.name, tc.want, got)
			}
		})
	}
}

func TestUnitAlicloudPolarDbAccountPageSearch(t *testing.T) {
	pageWithCount := func(count interface{}, records ...interface{}) map[string]interface{} {
		return map[string]interface{}{"Accounts": records, "PageRecordCount": count}
	}
	page := func(records ...interface{}) map[string]interface{} {
		return pageWithCount(json.Number(strconv.Itoa(len(records))), records...)
	}
	pageWithoutCount := func(records ...interface{}) map[string]interface{} {
		return map[string]interface{}{"Accounts": records}
	}
	account := func(name string) map[string]interface{} {
		return map[string]interface{}{"AccountName": name}
	}

	t.Run("found on second page", func(t *testing.T) {
		pages := []map[string]interface{}{
			page(account("other")),
			page(account("target")),
		}
		calls := 0
		fetch := func(pageNumber int) (map[string]interface{}, error) {
			calls++
			return pages[pageNumber-1], nil
		}
		found, _, ok, err := findPolarDbAccountByPages("target", 1, fetch)
		if err != nil || !ok || found["AccountName"] != "target" {
			t.Fatalf("want target on page 2, got found=%v ok=%v err=%v", found, ok, err)
		}
		if calls != 2 {
			t.Fatalf("want 2 fetches, got %d", calls)
		}
	})

	t.Run("not found after last page", func(t *testing.T) {
		pages := []map[string]interface{}{
			page(account("other")),
			page(),
		}
		calls := 0
		fetch := func(pageNumber int) (map[string]interface{}, error) {
			calls++
			return pages[pageNumber-1], nil
		}
		_, _, ok, err := findPolarDbAccountByPages("target", 1, fetch)
		if err != nil || ok {
			t.Fatalf("want not found, got ok=%v err=%v", ok, err)
		}
		if calls != 2 {
			t.Fatalf("want 2 fetches, got %d", calls)
		}
	})

	t.Run("short page stops pagination", func(t *testing.T) {
		calls := 0
		fetch := func(pageNumber int) (map[string]interface{}, error) {
			calls++
			return page(account("other")), nil
		}
		_, _, ok, err := findPolarDbAccountByPages("target", 5, fetch)
		if err != nil || ok {
			t.Fatalf("want not found, got ok=%v err=%v", ok, err)
		}
		if calls != 1 {
			t.Fatalf("want 1 fetch for a short page, got %d", calls)
		}
	})

	t.Run("fetch error propagates", func(t *testing.T) {
		wantErr := errors.New("boom")
		fetch := func(pageNumber int) (map[string]interface{}, error) {
			return nil, wantErr
		}
		_, _, ok, err := findPolarDbAccountByPages("target", 5, fetch)
		if !errors.Is(err, wantErr) || ok {
			t.Fatalf("want propagated error, got ok=%v err=%v", ok, err)
		}
	})

	t.Run("page limit returns an error instead of not found", func(t *testing.T) {
		calls := 0
		fetch := func(pageNumber int) (map[string]interface{}, error) {
			calls++
			return pageWithoutCount(account("other")), nil
		}
		_, _, ok, err := findPolarDbAccountByPages("target", 1, fetch)
		if !errors.Is(err, errPolarDbAccountPageLimitExceeded) || ok {
			t.Fatalf("want page limit error, got ok=%v err=%v", ok, err)
		}
		if calls != polarDbAccountMaxPages {
			t.Fatalf("want fallback pagination to stop at the page cap, got %d fetches", calls)
		}
	})

	t.Run("json number zero count stops without extra fetch", func(t *testing.T) {
		calls := 0
		fetch := func(pageNumber int) (map[string]interface{}, error) {
			calls++
			return pageWithCount(json.Number("0")), nil
		}
		_, _, ok, err := findPolarDbAccountByPages("target", 1, fetch)
		if err != nil || ok {
			t.Fatalf("want not found, got ok=%v err=%v", ok, err)
		}
		if calls != 1 {
			t.Fatalf("want 1 fetch for a zero count page, got %d", calls)
		}
	})

	t.Run("missing count falls back to record list", func(t *testing.T) {
		pages := []map[string]interface{}{
			pageWithoutCount(account("other")),
			page(account("target")),
		}
		calls := 0
		fetch := func(pageNumber int) (map[string]interface{}, error) {
			calls++
			return pages[pageNumber-1], nil
		}
		found, _, ok, err := findPolarDbAccountByPages("target", 1, fetch)
		if err != nil || !ok || found["AccountName"] != "target" {
			t.Fatalf("want target on page 2, got found=%v ok=%v err=%v", found, ok, err)
		}
		if calls != 2 {
			t.Fatalf("want 2 fetches, got %d", calls)
		}
	})

	t.Run("malformed count falls back to record list", func(t *testing.T) {
		pages := []map[string]interface{}{
			pageWithCount("not-a-number", account("other")),
			page(account("target")),
		}
		calls := 0
		fetch := func(pageNumber int) (map[string]interface{}, error) {
			calls++
			return pages[pageNumber-1], nil
		}
		found, _, ok, err := findPolarDbAccountByPages("target", 1, fetch)
		if err != nil || !ok || found["AccountName"] != "target" {
			t.Fatalf("want target on page 2, got found=%v ok=%v err=%v", found, ok, err)
		}
		if calls != 2 {
			t.Fatalf("want 2 fetches, got %d", calls)
		}
	})

	t.Run("json number count larger than records uses record fallback", func(t *testing.T) {
		calls := 0
		fetch := func(pageNumber int) (map[string]interface{}, error) {
			calls++
			return pageWithCount(json.Number("5"), account("other")), nil
		}
		_, _, ok, err := findPolarDbAccountByPages("target", 5, fetch)
		if err != nil || ok {
			t.Fatalf("want not found, got ok=%v err=%v", ok, err)
		}
		if calls != 1 {
			t.Fatalf("want 1 fetch when fallback record list is short, got %d", calls)
		}
	})
}
