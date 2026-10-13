package alicloud

import (
	"strings"
	"testing"
)

func sweepTestCenName(name string) bool {
	for _, prefix := range []string{"tf-testAcc", "tf_testAcc", "tf_example", "tf-example", "terraform-example"} {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

func sweepCenChildInRegion(region, childRegion string) bool {
	return region != "" && childRegion == region
}

func sweepTestVpcName(name string, isDefault bool) bool {
	if isDefault {
		return false
	}
	for _, prefix := range []string{"tf-testAcc", "tf_testAcc", "tf_test_", "tf-test-", "testAcc", "tf-example", "tf-test"} {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

func TestSweeperScope(t *testing.T) {
	for name, want := range map[string]bool{"tf-testAcc-Cen": true, "tf_example": true, "TF-TESTACC-Cen": true, "production": false, "production-tf-testAcc": false, "": false} {
		if got := sweepTestCenName(name); got != want {
			t.Fatalf("CEN name eligibility %q = %v, want %v", name, got, want)
		}
	}
	for _, tc := range []struct {
		name      string
		isDefault bool
		want      bool
	}{
		{"tf-testAccVpc123", false, true},
		{"TF-TESTACCVpc123", false, true},
		{"tf-example", false, true},
		{"tf-testAccDefault", true, false},
		{"production", false, false},
		{"production-tf-testAcc", false, false},
		{"", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := sweepTestVpcName(tc.name, tc.isDefault); got != tc.want {
				t.Fatalf("name eligibility = %v, want %v", got, tc.want)
			}
		})
	}
	for _, tc := range []struct {
		region, child string
		want          bool
	}{
		{"cn-beijing", "cn-beijing", true},
		{"cn-beijing", "cn-hangzhou", false},
		{"cn-beijing", "", false},
		{"", "", false},
	} {
		if got := sweepCenChildInRegion(tc.region, tc.child); got != tc.want {
			t.Fatalf("region eligibility (%q, %q) = %v, want %v", tc.region, tc.child, got, tc.want)
		}
	}
}
