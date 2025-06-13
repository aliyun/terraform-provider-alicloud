package connectivity

import "testing"

func TestRegionDefinitions(t *testing.T) {
	t.Run("AllExportedRegionsInValidRegions", func(t *testing.T) {
		definedRegions := []Region{
			Hangzhou, Qingdao, Beijing,
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
