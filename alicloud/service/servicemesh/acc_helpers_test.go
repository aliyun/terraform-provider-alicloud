package servicemesh_test

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	roacs "github.com/alibabacloud-go/cs-20151215/v8/client"
	tea "github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The generic ACC engine (sentinels, ResourceCheck/ResourceAttr/ResourceAttrCheck
// and the reflect-based describe-method dispatch, the resourceConfig builder, and
// the ResourceTestAccConfig/DataSourceTestAccConfigFunc factories) is now provided
// by alicloud/acctest/common.go and is called directly by the servicemesh callers.
//
// Only what cannot be reached from the public package is kept here:
//   - defaultRegionToTest / testAccPreCheck / testAccPreCheckWithRegions: the
//     package-level region sync wrappers mirror alicloud/acctest.PreCheck so the
//     local defaultRegionToTest stays in lock-step with ALICLOUD_REGION;
//   - dataSourceAttr + dataSourceTestCheck(WithPreCheck): the data-source runner
//     must keep ProtoV5ProviderFactories, because the public
//     acctest.DataSourceAttr.DataSourceTestCheck(WithPreCheck) hard-codes
//     SDKProviderFactories. buildDataSourceSteps / checkDataSourceAttr stay local
//     for the same reason (the public counterparts are unexported), but the
//     attribute core they build on is delegated to acctest.ResourceAttrInit;
//   - vpcSweepRPC / vpcACKAPI / nativeVpcACKAPI: the real VPC/ACK sweep adapters
//     used by the servicemesh sweeper, copied from the unreachable root test
//     package with identical semantics.

// defaultRegionToTest mirrors the package-level region used by the root test
// helpers so sharedClientForRegion stays in sync with ALICLOUD_REGION.
var defaultRegionToTest = func() string {
	if v := os.Getenv("ALICLOUD_REGION"); v != "" {
		return v
	}
	return "cn-beijing"
}()

// testAccPreCheck mirrors alicloud/acctest.PreCheck and keeps the package-level
// defaultRegionToTest in sync with the resolved region.
func testAccPreCheck(t *testing.T) {
	acctest.PreCheck(t)
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		defaultRegionToTest = "cn-beijing"
	} else {
		defaultRegionToTest = v
	}
}

// testAccPreCheckWithRegions mirrors alicloud/acctest.PreCheckWithRegions.
func testAccPreCheckWithRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	acctest.PreCheckWithRegions(t, supported, regions)
}

// in most cases, the TestCheckFunc list of dataSource test case is repeated，so we make an abstract in
// order to reduce redundant code.
// dataSourceAttr has 3 field ,incloud resourceId  existMapFunc fakeMapFunc, every dataSource test can use only one
type dataSourceAttr struct {
	// IDRefreshName, like "data.alicloud_dns_records.record"
	resourceId string

	// get existMap function
	existMapFunc func(rand int) map[string]string

	// get fakeMap function
	fakeMapFunc func(rand int) map[string]string
}

// get exist and empty resourceAttrMapUpdate function
func (dsa *dataSourceAttr) checkDataSourceAttr(rand int) (exist, empty acctest.ResourceAttrMapUpdate) {
	exist = acctest.ResourceAttrInit(dsa.resourceId, dsa.existMapFunc(rand)).ResourceAttrMapUpdateSet()
	empty = acctest.ResourceAttrInit(dsa.resourceId, dsa.fakeMapFunc(rand)).ResourceAttrMapUpdateSet()
	return
}

// according to configs generate step list and execute the test
// lintignore: AT001
func (dsa *dataSourceAttr) dataSourceTestCheck(t *testing.T, rand int, configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps:                    steps,
	})
}

// according to configs generate step list and execute the test with preCheck
// lintignore: AT001
func (dsa *dataSourceAttr) dataSourceTestCheckWithPreCheck(t *testing.T, rand int, preCheck func(), configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 preCheck,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps:                    steps,
	})
}

// per schema attribute test config
type dataSourceTestAccConfig struct {
	// be equal to testCase config string,but the result has only one record
	existConfig string

	// if the dataSourceAttr.existMapFunc returned map value not match we want, existChangMap can alter checkMap for existConfig
	existChangMap map[string]string

	// be equal to testCase config string,but the result is empty
	fakeConfig string

	// if the dataSourceAttr.fakeMapFunc returned map value not match we want, fakeChangMap can alter checkMap for fakeConfig
	fakeChangMap map[string]string
}

// build test cases for each attribute
func (conf *dataSourceTestAccConfig) buildDataSourceSteps(t *testing.T, info *dataSourceAttr, rand int) []resource.TestStep {
	testAccCheckExist, testAccCheckEmpty := info.checkDataSourceAttr(rand)
	var steps []resource.TestStep
	if conf.existConfig != "" {
		step := resource.TestStep{
			Config: conf.existConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckExist(conf.existChangMap),
			),
		}
		steps = append(steps, step)
	}
	if conf.fakeConfig != "" {
		step := resource.TestStep{
			Config: conf.fakeConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckEmpty(conf.fakeChangMap),
			),
		}
		steps = append(steps, step)
	}
	return steps
}

// vpcSweepRPC is a package-local copy of the RPC signature originally defined in
// alicloud/resource_alicloud_vpc_test.go. The root test package cannot be
// imported from a service package, so the function type used by the ASM delete
// adapter (nativeVpcASMDelete) is kept here with identical semantics.
type vpcSweepRPC func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error)

// vpcACKAPI and nativeVpcACKAPI are package-local copies of the flat-layout
// originals in alicloud/resource_alicloud_vpc_test.go. The root test package
// cannot be imported from a service package, so the ACK cluster helpers used by
// the servicemesh sweeper's vpcASMAPI are kept here with the same semantics.
type vpcACKAPI struct {
	list              func(int) ([]map[string]interface{}, error)
	detail            func(string) (map[string]interface{}, error)
	disableProtection func(string) error
	remove            func(string) error
	// Optional wait seams keep lifecycle regressions pure and bounded.
	pause           func(time.Duration)
	now             func() time.Time
	initialAttempts int
}

func nativeVpcACKAPI(client *connectivity.AliyunClient, region string) vpcACKAPI {
	var once sync.Once
	var csClient *roacs.Client
	var initErr error
	get := func() (*roacs.Client, error) {
		once.Do(func() { csClient, initErr = client.NewRoaCsClient() })
		return csClient, initErr
	}
	bodyMap := func(body interface{}) (map[string]interface{}, error) {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		var result map[string]interface{}
		err = json.Unmarshal(raw, &result)
		return result, err
	}
	return vpcACKAPI{
		list: func(page int) ([]map[string]interface{}, error) {
			sdk, err := get()
			if err != nil {
				return nil, err
			}
			response, err := sdk.DescribeClustersV1((&roacs.DescribeClustersV1Request{}).SetRegionId(region).SetPageNumber(int64(page)).SetPageSize(int64(helper.PageSizeLarge)))
			if err != nil {
				return nil, err
			}
			if response == nil || response.Body == nil {
				return nil, fmt.Errorf("ACK empty cluster list response")
			}
			var rows []map[string]interface{}
			for _, cluster := range response.Body.Clusters {
				item, err := bodyMap(cluster)
				if err != nil {
					return nil, err
				}
				rows = append(rows, item)
			}
			return rows, nil
		},
		detail: func(id string) (map[string]interface{}, error) {
			sdk, err := get()
			if err != nil {
				return nil, err
			}
			response, err := sdk.DescribeClusterDetail(tea.String(id))
			if err != nil {
				return nil, err
			}
			if response == nil || response.Body == nil {
				return nil, fmt.Errorf("ACK empty cluster detail")
			}
			return bodyMap(response.Body)
		},
		disableProtection: func(id string) error {
			sdk, err := get()
			if err != nil {
				return err
			}
			_, err = sdk.ModifyCluster(tea.String(id), (&roacs.ModifyClusterRequest{}).SetDeletionProtection(false))
			return err
		},
		remove: func(id string) error {
			sdk, err := get()
			if err != nil {
				return err
			}
			_, err = sdk.DeleteCluster(tea.String(id), &roacs.DeleteClusterRequest{})
			return err
		},
	}
}
