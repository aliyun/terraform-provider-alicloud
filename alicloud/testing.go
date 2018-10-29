package alicloud

import (
	"log"
	"os"
)

type RegionalFeature string

const (
	SlbSpecification  = RegionalFeature("SLB_SPECIFICATION")
	SlbPayByBandwidth = RegionalFeature("SLB_PAY_BY_BANDWIDTH")
	FunctionCompute   = RegionalFeature("FUNCTION_COMPUTE")
	PrivateZone       = RegionalFeature("PRIVATE_ZONE")
	RdsMultiAZ        = RegionalFeature("RDS_MULTI_AZ")
	ClassicNetwork    = RegionalFeature("CLASSIC_NETWORK")
)

func isRegionSupports(features ...RegionalFeature) bool {
	for _, feature := range features {
		featureSkipped := os.Getenv("ALICLOUD_SKIP_TESTS_FOR_"+string(feature)) == "true"
		if featureSkipped {
			return false
		}
	}
	return true
}

func logTestSkippedBecauseOfUnsupportedRegionalFeatures(testName string, features ...RegionalFeature) {
	log.Printf("[INFO] Test '%v' skipped because the current region doesn't support all the following features: %v\n", testName, features)
}
