package alicloud

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/fc-go-sdk"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var defaultRegionToTest = os.Getenv("ALICLOUD_REGION")

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"alicloud": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ALICLOUD_ACCESS_KEY"); v == "" {
		t.Fatal("ALICLOUD_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_SECRET_KEY"); v == "" {
		t.Fatal("ALICLOUD_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		log.Println("[INFO] Test: Using cn-beijing as test region")
		os.Setenv("ALICLOUD_REGION", "cn-beijing")
	} else {
		defaultRegionToTest = v
	}
}

// currently not all account site type support create PostPaid resources, PayByBandwidth and other limits.
// The setting of account site type can skip some unsupported cases automatically.

func testAccPreCheckWithAccountSiteType(t *testing.T, account AccountSite) {
	defaultAccount := string(DomesticSite)
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_ACCOUNT_SITE")); v != "" {
		defaultAccount = v
	}
	if defaultAccount != string(account) {
		t.Skipf("Skipping unsupported account type %s-Site. It only supports %s-Site.", defaultAccount, account)
		t.Skipped()
	}
}

func testAccPreCheckPrePaidResources(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ENABLE_CHECKING_PRE_PAID")); v != "true" {
		t.Skip("Skipping testing PrePaid resources, otherwise setting environment parameter 'ENABLE_CHECKING_PRE_PAID'.")
		t.Skipped()
	}
}

func testAccClassicNetworkResources(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ENABLE_CLASSIC_NETWORK")); v != "true" {
		t.Skip("Skipping testing classic resources, otherwise setting environment parameter 'ENABLE_CLASSIC_NETWORK'.")
		t.Skipped()
	}
}

// Skip automatically the testcases which does not support some known regions.
// If supported is true, the regions should a list of supporting the service regions.
// If supported is false, the regions should a list of unsupporting the service regions.
// If the region is unsupported and has backend region, the backend region will instead
func testAccPreCheckWithRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	if v := os.Getenv("ALICLOUD_ACCESS_KEY"); v == "" {
		t.Fatal("ALICLOUD_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_SECRET_KEY"); v == "" {
		t.Fatal("ALICLOUD_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		t.Logf("[WARNING] The region is not set and using cn-beijing as test region")
		os.Setenv("ALICLOUD_REGION", "cn-beijing")
	}
	checkoutSupportedRegions(t, supported, regions)
}

func checkoutSupportedRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	region := os.Getenv("ALICLOUD_REGION")
	find := false
	backupRegion := string(connectivity.APSouthEast1)
	if region == string(connectivity.APSouthEast1) {
		backupRegion = string(connectivity.EUCentral1)
	}

	checkoutRegion := os.Getenv("CHECKOUT_REGION")
	if checkoutRegion == "true" {
		if region == string(connectivity.Hangzhou) {
			region = string(connectivity.EUCentral1)
			os.Setenv("ALICLOUD_REGION", region)
		}
	}
	backupRegionFind := false
	hangzhouRegionFind := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
		if string(r) == backupRegion {
			backupRegionFind = true
		}
		if string(connectivity.Hangzhou) == string(r) {
			hangzhouRegionFind = true
		}
	}

	if (find && !supported) || (!find && supported) {
		if supported {
			if backupRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, backupRegion)
				os.Setenv("ALICLOUD_REGION", backupRegion)
				defaultRegionToTest = backupRegion
				return
			}
			if hangzhouRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, connectivity.Hangzhou)
				os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
				defaultRegionToTest = string(connectivity.Hangzhou)
				return
			}
			t.Skipf("Skipping unsupported region %s. Supported regions: %s.", region, regions)
		} else {
			if !backupRegionFind {
				t.Logf("Skipping unsupported region %s. Unsupported regions: %s. Using %s as this test region", region, regions, backupRegion)
				os.Setenv("ALICLOUD_REGION", backupRegion)
				defaultRegionToTest = backupRegion
				return
			}
			if !hangzhouRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, connectivity.Hangzhou)
				os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
				defaultRegionToTest = string(connectivity.Hangzhou)
				return
			}
			t.Skipf("Skipping unsupported region %s. Unsupported regions: %s.", region, regions)
		}
		t.Skipped()
	}
}

func checkoutAccount(t *testing.T, SLAVE bool) {
	if SLAVE {
		if os.Getenv("ALICLOUD_ACCESS_KEY_SLAVE") == "" || os.Getenv("ALICLOUD_SECRET_KEY_SLAVE") == "" {
			t.Logf("\nALICLOUD_ACCESS_KEY_SLAVE or ALICLOUD_SECRET_KEY_SLAVE is empty and please add them.")
		} else {
			os.Setenv("ALICLOUD_ACCESS_KEY_MASTER", os.Getenv("ALICLOUD_ACCESS_KEY"))
			os.Setenv("ALICLOUD_SECRET_KEY_MASTER", os.Getenv("ALICLOUD_SECRET_KEY"))
			os.Setenv("ALICLOUD_ACCESS_KEY", os.Getenv("ALICLOUD_ACCESS_KEY_SLAVE"))
			os.Setenv("ALICLOUD_SECRET_KEY", os.Getenv("ALICLOUD_SECRET_KEY_SLAVE"))
			t.Logf("%s is using the slave account", t.Name())
		}
	} else {
		if os.Getenv("ALICLOUD_ACCESS_KEY_MASTER") == "" || os.Getenv("ALICLOUD_SECRET_KEY_MASTER") == "" {
			t.Logf("\nALICLOUD_ACCESS_KEY_MASTER or ALICLOUD_SECRET_KEY_MASTER is empty and please add them.")
		} else {
			os.Setenv("ALICLOUD_ACCESS_KEY", os.Getenv("ALICLOUD_ACCESS_KEY_MASTER"))
			os.Setenv("ALICLOUD_SECRET_KEY", os.Getenv("ALICLOUD_SECRET_KEY_MASTER"))
		}
	}
}

// Skip automatically the sweep testcases which does not support some known regions.
// If supported is true, the regions should a list of supporting the service regions.
// If supported is false, the regions should a list of unsupporting the service regions.
func testSweepPreCheckWithRegions(region string, supported bool, regions []connectivity.Region) bool {
	find := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
	}
	return (find && !supported) || (!find && supported)
}

func testAccCheckAlicloudDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("data source ID not set")
		}
		return nil
	}
}

func testAccPreCheckWithMultipleAccount(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_ACCESS_KEY_2")); v == "" {
		t.Skipf("Skipping unsupported test with multiple account")
		t.Skipped()
	}
}

func testAccPreCheckOSSForImageImport(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_OSS_BUCKET_FOR_IMAGE")); v == "" {
		t.Skipf("Skipping tests without OSS_Bucket set.")
		t.Skipped()
	}
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_OSS_OBJECT_FOR_IMAGE")); v == "" {
		t.Skipf("Skipping OSS_Object does not exist.")
		t.Skipped()
	}
}

func testAccPreCheckKMSForKeyIdImport(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_KMS_KEY_ID")); v == "" {
		t.Skipf("Skipping tests without KEY_ID set.")
		t.Skipped()
	}
}

func testAccPreCheckWithCmsContactGroupSetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_CMS_CONTACT_GROUP")); v == "" {
		t.Skipf("Skipping the test case with no cms contact group setting")
		t.Skipped()
	}
}

func testAccPreCheckWithSmartAccessGatewaySetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("SAG_INSTANCE_ID")); v == "" {
		t.Skipf("Skipping the test case with no sag instance id setting")
		t.Skipped()
	}
}

func testAccPreCheckWithResourceManagerFloderIdSetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_RESOURCE_MANAGER_FOLDER_ID1")); v == "" {
		t.Skip("Skipping the test case with no sag folder id setting")
		t.Skipped()
	}
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_RESOURCE_MANAGER_FOLDER_ID2")); v == "" {
		t.Skip("Skipping the test case with no sag folder id setting")
		t.Skipped()
	}
}

func testAccPreCheckWithEnvVariable(t *testing.T, envVariableName string) {
	if v := strings.TrimSpace(os.Getenv(envVariableName)); v == "" {
		t.Skipf("Skipping the test case with no env variable %s", envVariableName)
		t.Skipped()
	}
}

func testAccPreCheckWithSlbInstanceSetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_SLB_INSTANCE_ID")); v == "" {
		t.Skipf("Skipping the test case with no slb instance id setting")
		t.Skipped()
	}
}

func testAccPreCheckWithSmartAccessGatewayAppSetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("SAG_APP_INSTANCE_ID")); v == "" {
		t.Skipf("Skipping the test case with no sag app instance id setting")
		t.Skipped()
	}
}

func testAccPreCheckWithTime(t *testing.T, days []int) {
	skipped := true
	for _, d := range days {
		if time.Now().Day() == d {
			skipped = false
			break
		}
	}
	if skipped {
		t.Skipf("Skipping the test case when not in specified days %#v of every month", days)
		t.Skipped()
	}
}

func testAccPreCheckWithAlikafkaAclEnable(t *testing.T) {
	aclEnable := os.Getenv("ALICLOUD_ALIKAFKA_ACL_ENABLE")

	if aclEnable != "true" {
		t.Skipf("Skipping the test case because the acl is not enabled.")
		t.Skipped()
	}
}

func testAccPreCheckEnterpriseAccountEnabled(t *testing.T) {
	enable := os.Getenv("ENTERPRISE_ACCOUNT_ENABLED")
	if enable != "true" {
		t.Skipf("Skipping the test case because the enterprise account is not enabled.")
		t.Skipped()
	}
}

func testAccPreCheckWithNoDefaultVpc(t *testing.T) {
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	client := rawClient.(*connectivity.AliyunClient)
	request := vpc.CreateDescribeVpcsRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.PageNumber = requests.NewInteger(1)
	request.IsDefault = requests.NewBoolean(true)

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeVpcs(request)
	})
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	response, _ := raw.(*vpc.DescribeVpcsResponse)

	if len(response.Vpcs.Vpc) < 1 {
		request.IsDefault = requests.NewBoolean(false)
		request.VpcName = "default"
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpcs(request)
		})
		if err != nil {
			t.Skipf("Skipping the test case with err: %s", err)
			t.Skipped()
		}
		response2, _ := raw.(*vpc.DescribeVpcsResponse)
		if len(response2.Vpcs.Vpc) < 1 {
			t.Skipf("Skipping the test case with there is no default vpc")
			t.Skipped()
		}
	}
}

func testAccPreCheckWithNoDefaultVswitch(t *testing.T) {
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	client := rawClient.(*connectivity.AliyunClient)
	request := vpc.CreateDescribeVSwitchesRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.PageNumber = requests.NewInteger(1)
	request.IsDefault = requests.NewBoolean(true)

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeVSwitches(request)
	})
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	response, _ := raw.(*vpc.DescribeVSwitchesResponse)

	if len(response.VSwitches.VSwitch) < 1 {
		t.Skipf("Skipping the test case with there is no default vswitche")
		t.Skipped()
	}
}

func testAccPreCheckWithResourceManagerAccountsSetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_RESOURCE_MANAGER_ACCOUNT_ID")); v == "" {
		t.Skip("Skipping the test case with no sag account id setting")
		t.Skipped()
	}
}

func testAccPreCheckWithResourceManagerHandshakesSetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("INVITED_ALICLOUD_ACCOUNT_ID")); v == "" {
		t.Skipf("Skipping the test case with there is no \"INVITED_ALICLOUD_ACCOUNT_ID\" setting")
		t.Skipped()
	}
}

var providerCommon = `
provider "alicloud" {
	assume_role {}
}
`

func TestAccAlicloudProviderEcs(t *testing.T) {
	var v ecs.Instance

	resourceId := "alicloud_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAcc%sEcsInstanceConfigVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return providerCommon + resourceInstanceVpcConfigDependence(name)
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":        "${data.alicloud_images.default.images.0.id}",
					"security_groups": []string{"${alicloud_security_group.default.0.id}"},
					"instance_type":   "${data.alicloud_instance_types.default.instance_types.0.id}",

					"availability_zone":             "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}",
					"system_disk_category":          "cloud_efficiency",
					"instance_name":                 "${var.name}",
					"key_name":                      "${alicloud_key_pair.default.key_name}",
					"spot_strategy":                 "NoSpot",
					"spot_price_limit":              "0",
					"security_enhancement_strategy": "Active",
					"user_data":                     "I_am_user_data",

					"vswitch_id": "${alicloud_vswitch.default.id}",
					"role_name":  "${alicloud_ram_role.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"key_name":      name,
						"role_name":     name,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudProviderFC(t *testing.T) {
	var v *fc.GetFunctionOutput
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%salicloudfcfunction-%d", defaultRegionToTest, rand)
	var basicMap = map[string]string{
		"service":     CHECKSET,
		"name":        name,
		"runtime":     "python2.7",
		"description": "tf",
		"handler":     "hello.handler",
		"oss_bucket":  CHECKSET,
		"oss_key":     CHECKSET,
	}
	resourceId := "alicloud_fc_function.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return providerCommon + resourceFCFunctionConfigDependence(name)
	})

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service":     "${alicloud_fc_service.default.name}",
					"name":        "${var.name}",
					"runtime":     "python2.7",
					"description": "tf",
					"handler":     "hello.handler",
					"oss_bucket":  "${alicloud_oss_bucket.default.id}",
					"oss_key":     "${alicloud_oss_bucket_object.default.key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudProviderOss(t *testing.T) {
	var v oss.GetBucketInfoResult

	resourceId := "alicloud_oss_bucket.default"
	ra := resourceAttrInit(resourceId, ossBucketBasicMap)

	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sbucket-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return providerCommon + resourceOssBucketConfigDependence(name)
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudProviderLog(t *testing.T) {
	var v *sls.LogProject
	resourceId := "alicloud_log_project.default"
	ra := resourceAttrInit(resourceId, logProjectMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%slogproject-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return providerCommon + resourceLogProjectConfigDependence(name)
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudProviderDatahub(t *testing.T) {
	var v *datahub.GetProjectResult

	resourceId := "alicloud_datahub_project.default"
	ra := resourceAttrInit(resourceId, datahubProjectBasicMap)

	serviceFunc := func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testaccdatahubproject%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return providerCommon + resourceDatahubProjectConfigDependence(name)
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
		},
	})
}
