package sweep

import (
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestSweepAll(t *testing.T) {
	for _, value := range []string{"", "true", "false", "TRUE", "True", " true ", "1"} {
		t.Run("value="+value, func(t *testing.T) {
			t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", value)
			if got, want := SweepAll(), value == "true"; got != want {
				t.Fatalf("got %t, want %t", got, want)
			}
		})
	}
}

func TestSweepOlderThan(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  time.Duration
	}{
		{"", 0}, {" \t", 0}, {"invalid", 0}, {"12", 0}, {"0", 0},
		{"12h", 12 * time.Hour}, {" 1h30m ", 90 * time.Minute},
		{"-1h", -time.Hour}, {"1.5s", 1500 * time.Millisecond},
	} {
		t.Run("value="+tc.value, func(t *testing.T) {
			t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", tc.value)
			if got := SweepOlderThan(); got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestParseSweepCreationTime(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  time.Time
	}{
		{"2021-04-18T15:02:37Z", time.Date(2021, 4, 18, 15, 2, 37, 0, time.UTC)},
		{"2021-04-18T23:02:37+08:00", time.Date(2021, 4, 18, 15, 2, 37, 0, time.UTC)},
		{"2021-04-18T15:02:37.123Z", time.Date(2021, 4, 18, 15, 2, 37, 123000000, time.UTC)},
		{"2021-04-18T15:02Z", time.Date(2021, 4, 18, 15, 2, 0, 0, time.UTC)},
		{"2021-04-18T23:02+08:00", time.Date(2021, 4, 18, 15, 2, 0, 0, time.UTC)},
		{"2021-04-18 15:02:37Z", time.Date(2021, 4, 18, 15, 2, 37, 0, time.UTC)},
		{"2021-04-18 23:02:37+08:00", time.Date(2021, 4, 18, 15, 2, 37, 0, time.UTC)},
		{"2021-04-18 15:02:37", time.Date(2021, 4, 18, 15, 2, 37, 0, time.UTC)},
		{"", time.Time{}}, {"invalid", time.Time{}}, {"2021-04-18", time.Time{}},
		{" 2021-04-18T15:02:37Z ", time.Time{}}, {"1618758157", time.Time{}},
	} {
		t.Run("value="+tc.value, func(t *testing.T) {
			got, ok := ParseSweepCreationTime(tc.value)
			if ok != !tc.want.IsZero() || !got.Equal(tc.want) {
				t.Fatalf("got %s, %t; want %s, %t", got, ok, tc.want, !tc.want.IsZero())
			}
		})
	}
}

func TestSweepRecentCheck(t *testing.T) {
	old := time.Now().Add(-48 * time.Hour).UTC().Format(time.RFC3339)
	recent := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
	future := time.Now().Add(48 * time.Hour).UTC().Format(time.RFC3339)
	for _, tc := range []struct {
		name, threshold string
		item            map[string]interface{}
		skip            bool
		reason          string
	}{
		{"disabled", "", nil, false, ""},
		{"zero", "0", nil, false, ""},
		{"negative", "-1h", nil, false, ""},
		{"invalid_threshold", "bad", nil, false, ""},
		{"missing", "12h", nil, true, "no creation time field found; ALICLOUD_SWEEP_OLDER_THAN=12h0m0s set, protecting"},
		{"empty", "12h", map[string]interface{}{"CreationTime": ""}, true, "no creation time field found"},
		{"nil", "12h", map[string]interface{}{"CreationTime": nil}, true, "no creation time field found"},
		{"invalid", "12h", map[string]interface{}{"CreationTime": "bad"}, true, "cannot parse creation time \"bad\"; ALICLOUD_SWEEP_OLDER_THAN=12h0m0s set, protecting"},
		{"numeric", "12h", map[string]interface{}{"CreationTime": 123}, true, "cannot parse creation time \"123\""},
		{"old", "12h", map[string]interface{}{"CreationTime": old}, false, ""},
		{"recent", "12h", map[string]interface{}{"CreationTime": recent}, true, "newer than ALICLOUD_SWEEP_OLDER_THAN=12h0m0s"},
		{"future", "12h", map[string]interface{}{"CreationTime": future}, true, "newer than ALICLOUD_SWEEP_OLDER_THAN=12h0m0s"},
		{"first_field_wins", "12h", map[string]interface{}{"CreationTime": recent, "CreateTime": old}, true, "newer than"},
		{"invalid_first_wins", "12h", map[string]interface{}{"CreationTime": "bad", "CreateTime": old}, true, "cannot parse creation time"},
		{"skip_empty_field", "12h", map[string]interface{}{"CreationTime": "", "CreateTime": old}, false, ""},
		{"skip_nil_field", "12h", map[string]interface{}{"CreationTime": nil, "CreateTime": old}, false, ""},
		{"second_field_precedence", "12h", map[string]interface{}{"CreateTime": recent, "creation_time": old}, true, "newer than"},
		{"third_field_precedence", "12h", map[string]interface{}{"creation_time": recent, "create_time": old}, true, "newer than"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", tc.threshold)
			skip, reason := SweepRecentCheck(tc.item)
			if skip != tc.skip || (tc.reason == "" && reason != "") || !strings.Contains(reason, tc.reason) {
				t.Fatalf("got skip=%t, reason=%q; want skip=%t, reason containing %q", skip, reason, tc.skip, tc.reason)
			}
		})
	}
	for _, key := range []string{"CreationTime", "CreateTime", "creation_time", "create_time"} {
		t.Run(key, func(t *testing.T) {
			t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "12h")
			if skip, reason := SweepRecentCheck(map[string]interface{}{key: old}); skip || reason != "" {
				t.Fatalf("old resource unexpectedly protected: %s", reason)
			}
		})
	}
}

func TestSweepPreCheckWithRegions(t *testing.T) {
	for _, tc := range []struct {
		name      string
		supported bool
		regions   []connectivity.Region
		skip      bool
	}{
		{"allowed", true, []connectivity.Region{connectivity.Hangzhou}, false},
		{"not_allowed", true, []connectivity.Region{connectivity.Beijing}, true},
		{"empty_allowlist", true, nil, true},
		{"excluded", false, []connectivity.Region{connectivity.Hangzhou}, true},
		{"not_excluded", false, []connectivity.Region{connectivity.Beijing}, false},
		{"empty_exclusions", false, nil, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := SweepPreCheckWithRegions("cn-hangzhou", tc.supported, tc.regions); got != tc.skip {
				t.Fatalf("got skip=%t, want %t", got, tc.skip)
			}
		})
	}
}
