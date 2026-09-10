package connectivity

import "testing"

func TestUnitCommonRegionDefinitions(t *testing.T) {
	t.Run("AllExportedRegionsInValidRegions", func(t *testing.T) {
		definedRegions := []Region{
			Hangzhou, Qingdao, Beijing, ZhongWei, APSouthEast8, NASouth1,
		}

		validMap := make(map[Region]bool)
		for _, r := range ValidRegions {
			validMap[r] = true
		}

		for _, dr := range definedRegions {
			if !validMap[dr] {
				t.Errorf("Region %s not found in ValidRegions", dr)
			}
		}
	})

	t.Run("ServiceSpecificRegionsValidation", func(t *testing.T) {
		testCases := []struct {
			serviceRegions []Region
		}{
			{EcsClassicSupportedRegions},
			{DdosBgpSupportRegions},
		}

		validMap := make(map[Region]bool)
		for _, r := range ValidRegions {
			validMap[r] = true
		}

		for _, tc := range testCases {
			for _, sr := range tc.serviceRegions {
				if !validMap[sr] {
					t.Errorf("Service region %s not in ValidRegions", sr)
				}
			}
		}
	})
}

func TestUnitVpnGatewayVpnAttachmentSupportRegions(t *testing.T) {
	if len(VpnGatewayVpnAttachmentSupportRegions) != 1 {
		t.Fatalf("expected exactly one supported region, got %d", len(VpnGatewayVpnAttachmentSupportRegions))
	}

	if got := VpnGatewayVpnAttachmentSupportRegions[0]; got != EUCentral1 {
		t.Fatalf("expected supported region %q, got %q", EUCentral1, got)
	}
}
