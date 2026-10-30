package alicloud

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

// Real-time log delivery does not require an enterpriseplan ESA instance: CreateSiteDeliveryTask
// only needs a plan that supports it, and the acceptance account's advancedplan instance reports
// free delivery quota. The instance owns a single site slot, which the pending test site these
// cases share already occupies, so the cases run against that site instead of creating one.
//
// Every case takes its task_name from the per-case random `name`, so a task left behind by an
// interrupted run cannot collide with a later run (RealTimeLog.ProjectExistError) and the cases
// stay independent on the shared site. The site is resolved through
// `plan_subscribe_type = "advancedplan"`, which returns the only site that instance owns. This
// resource has no sweeper, so a task abandoned by a killed run has to be deleted with
// DeleteSiteDeliveryTask before rerunning.

// Every case delivers into a pre-created SLS project/log store that the test account owns. The
// fixture is an environment input, not something the cases invent:
//   - the account is already at its SLS project quota (CreateProject answers 403 ProjectQuotaExceed,
//     "most has 300 project"), so creating a project per case is not possible;
//   - the project name the cases used before is no longer usable: an account-wide
//     CreateProject answers 400 ProjectAlreadyExist for it while the account's own project list
//     does not contain it, so the cases now share an existing test account project that was left
//     behind by a real-time log delivery acceptance run (the tf-testacc... project below);
//   - pointing the cases at an arbitrary project would write logs into an object another team
//     owns, so the shared project only hosts a log store dedicated to these cases.
//
// The values default to that dedicated fixture and can be overridden with environment variables.
// The dependence below proves at plan time that the project and the log store exist in the
// configured region (a missing or misnamed fixture fails the read, it never silently skips or
// passes). The generated SiteDeliveryTask_test uses the same fixture through these constants.
const (
	esaSiteDeliveryTaskSlsProjectEnv  = "ALICLOUD_ESA_SITE_DELIVERY_TASK_SLS_PROJECT"
	esaSiteDeliveryTaskSlsLogStoreEnv = "ALICLOUD_ESA_SITE_DELIVERY_TASK_SLS_LOG_STORE"
	esaSiteDeliveryTaskSlsRegionEnv   = "ALICLOUD_ESA_SITE_DELIVERY_TASK_SLS_REGION"

	esaSiteDeliveryTaskSlsDefaultProject  = "tf-testacccn-beijingcdnrealtimelogdelivery644953"
	esaSiteDeliveryTaskSlsDefaultLogStore = "esa-dt-1029-110681"
	esaSiteDeliveryTaskSlsDefaultRegion   = "cn-beijing"
)

// esaSiteDeliveryTaskSlsFixture is the SLS project, log store and region the delivery cases write
// to. All three belong together: the region must be the region that owns the project, otherwise
// the verification reads and the task's sls_delivery block disagree.
type esaSiteDeliveryTaskSlsFixture struct {
	Project  string
	LogStore string
	Region   string
}

// esaSiteDeliveryTaskSlsFixtureFromEnv resolves the fixture from the environment, falling back to
// the dedicated fixture. Empty environment values are ignored so that an operator only has to
// override the field that changed.
func esaSiteDeliveryTaskSlsFixtureFromEnv() esaSiteDeliveryTaskSlsFixture {
	fixture := esaSiteDeliveryTaskSlsFixture{
		Project:  esaSiteDeliveryTaskSlsDefaultProject,
		LogStore: esaSiteDeliveryTaskSlsDefaultLogStore,
		Region:   esaSiteDeliveryTaskSlsDefaultRegion,
	}
	if v := strings.TrimSpace(os.Getenv(esaSiteDeliveryTaskSlsProjectEnv)); v != "" {
		fixture.Project = v
	}
	if v := strings.TrimSpace(os.Getenv(esaSiteDeliveryTaskSlsLogStoreEnv)); v != "" {
		fixture.LogStore = v
	}
	if v := strings.TrimSpace(os.Getenv(esaSiteDeliveryTaskSlsRegionEnv)); v != "" {
		fixture.Region = v
	}
	return fixture
}

// esaSiteDeliveryTaskBasicDependence renders the shared fixture used by the first-create and
// disappears cases: the account's site, plus read-only lookups that prove the SLS project and log
// store of the fixture exist in the configured region before the task is created.
func esaSiteDeliveryTaskBasicDependence(fixture esaSiteDeliveryTaskSlsFixture) string {
	return fmt.Sprintf(`
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "advancedplan"
}

provider "alicloud" {
  alias  = "slsfixture"
  region = %q
}

data "alicloud_log_projects" "default" {
  provider   = alicloud.slsfixture
  name_regex = "^%s$"
}

data "alicloud_log_stores" "default" {
  provider   = alicloud.slsfixture
  project    = %q
  name_regex = "^%s$"
}
`, fixture.Region, regexp.QuoteMeta(fixture.Project), fixture.Project, regexp.QuoteMeta(fixture.LogStore))
}

// testAccCheckEsaSiteDeliveryTaskSlsFixture proves the fixture reads returned exactly one project
// and one log store. Without it an empty lookup would only surface as an index error elsewhere.
func testAccCheckEsaSiteDeliveryTaskSlsFixture(fixture esaSiteDeliveryTaskSlsFixture) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.alicloud_log_projects.default", "names.#", "1"),
		resource.TestCheckResourceAttr("data.alicloud_log_projects.default", "names.0", fixture.Project),
		resource.TestCheckResourceAttr("data.alicloud_log_stores.default", "names.#", "1"),
		resource.TestCheckResourceAttr("data.alicloud_log_stores.default", "names.0", fixture.LogStore),
	)
}

// Test ESA SiteDeliveryTask. >>> Resource test cases, automatically generated.
// Case resource_SiteDeliveryTask_test
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)
	sls := esaSiteDeliveryTaskSlsFixtureFromEnv()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckEsaSiteDeliveryTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_delivery": []map[string]interface{}{
						{
							"sls_project":   sls.Project,
							"sls_region":    sls.Region,
							"sls_log_store": sls.LogStore,
						},
					},
					"site_id":       "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"data_center":   "global",
					"discard_rate":  "0.0",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "sls",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_esa_sites.default", "sites.#", "1"),
					testAccCheckEsaSiteDeliveryTaskSlsFixture(sls),
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"discard_rate":  "0.1",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
					"field_name":    "OriginIP,ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"discard_rate":  "0.2",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "offline",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "online",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sls_delivery"},
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_testBasicDependence(name string) string {
	return esaSiteDeliveryTaskBasicDependence(esaSiteDeliveryTaskSlsFixtureFromEnv())
}

// Case resource_SiteDeliveryTask_http_test
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_http_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_http_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_http_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckEsaSiteDeliveryTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id": "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"http_delivery": []map[string]interface{}{
						{
							"compress":        "gzip",
							"log_body_suffix": "cdnVersion:1.0",
							"standard_auth_param": []map[string]interface{}{
								{
									"private_key":  "***",
									"url_path":     "v1/log/upload",
									"expired_time": "300",
								},
							},
							"standard_auth_on": "false",
							"log_body_prefix":  "cdnVersion:1.0",
							"header_param": map[string]interface{}{
								"x-auth": "test-header-value",
							},
							"query_param": map[string]interface{}{
								"auth_token": "test-query-value",
							},
							"dest_url":          "http://11.177.129.13:8081",
							"max_batch_size":    "1000",
							"max_retry":         "3",
							"transform_timeout": "10",
							"max_batch_mb":      "5",
						},
					},
					"data_center":   "global",
					"discard_rate":  "0.0",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_esa_sites.default", "sites.#", "1"),
					testAccCheck(map[string]string{}),
					testAccCheckEsaSiteDeliveryTaskHttpSinkParams(resourceId, "x-auth", "test-header-value", "auth_token", "test-query-value"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "offline",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_delivery"},
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_http_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_http_testBasicDependence(name string) string {
	return `
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "advancedplan"
}
`
}

// Case resource_SiteDeliveryTask_oss_test_test
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_oss_test_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_oss_test_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_oss_test_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckEsaSiteDeliveryTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":       "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"data_center":   "global",
					"discard_rate":  "0.0",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"oss_delivery": []map[string]interface{}{
						{
							"bucket_name": "test-log3",
							"region":      "cn-hongkong",
							"prefix_path": "",
							"aliuid":      "1097011697834102",
						},
					},
					"delivery_type": "oss",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_esa_sites.default", "sites.#", "1"),
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "offline",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oss_delivery"},
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_oss_test_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_oss_test_testBasicDependence(name string) string {
	return `
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "advancedplan"
}
`
}

// Case resource_SiteDeliveryTask_kafka_test_test
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_kafka_test_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_kafka_test_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_kafka_test_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckEsaSiteDeliveryTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":       "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"data_center":   "global",
					"discard_rate":  "0.0",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
					"kafka_delivery": []map[string]interface{}{
						{
							"compress":       "lz4",
							"user_name":      "kafkaadmin",
							"machanism_type": "plain",
							"brokers": []string{
								"11.177.129.13:9092",
								"11.177.129.14:9092",
							},
							"balancer":  "kafka.LeastBytes",
							"topic":     "access",
							"user_auth": "true",
							"password":  "kafkaadmin-cdnlog",
						},
					},
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "kafka",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_esa_sites.default", "sites.#", "1"),
					testAccCheck(map[string]string{}),
				),
			},
			{
				// TypeList order coverage (kafka_delivery.brokers): the only change is the broker
				// order, and the post-apply plan must be empty. The provider's Read does not
				// populate kafka_delivery (GetSiteDeliveryTask does return SinkConfig, but Read
				// ignores it), so this step is a regression guard against a future permadiff and
				// does not assert the service's ordering semantics.
				Config: testAccConfig(map[string]interface{}{
					"kafka_delivery": []map[string]interface{}{
						{
							"compress":       "lz4",
							"user_name":      "kafkaadmin",
							"machanism_type": "plain",
							"brokers": []string{
								"11.177.129.14:9092",
								"11.177.129.13:9092",
							},
							"balancer":  "kafka.LeastBytes",
							"topic":     "access",
							"user_auth": "true",
							"password":  "kafkaadmin-cdnlog",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "offline",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kafka_delivery"},
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_kafka_test_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_kafka_test_testBasicDependence(name string) string {
	return `
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "advancedplan"
}
`
}

// Case resource_SiteDeliveryTask_aws3_test_test
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_aws3_test_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_aws3_test_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_aws3_test_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckEsaSiteDeliveryTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"data_center":  "global",
					"discard_rate": "0.0",
					"task_name":    name,
					"s3_delivery": []map[string]interface{}{
						{
							"secret_key":             "qYl*****PJeFh",
							"endpoint":               "https://s3.oss-cn-hangzhou.aliyuncs.com",
							"vertify_type":           "",
							"region":                 "us-east-2",
							"bucket_path":            "openapi-test-esa/test",
							"server_side_encryption": "false",
							"access_key":             "AK****YMC",
							"prefix_path":            "logriver-test/log",
							"s3_cmpt":                "true",
						},
					},
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "aws3",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_esa_sites.default", "sites.#", "1"),
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "offline",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"s3_delivery"},
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_aws3_test_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_aws3_test_testBasicDependence(name string) string {
	return `
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "advancedplan"
}
`
}

// Test ESA SiteDeliveryTask first create with status online. <<< Resource test cases, automatically generated.
// Case resource_SiteDeliveryTask_first_online_test
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_first_online_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_first_online_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)
	sls := esaSiteDeliveryTaskSlsFixtureFromEnv()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_first_online_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckEsaSiteDeliveryTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_delivery": []map[string]interface{}{
						{
							"sls_project":   sls.Project,
							"sls_region":    sls.Region,
							"sls_log_store": sls.LogStore,
						},
					},
					"site_id":       "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"data_center":   "global",
					"discard_rate":  "0.0",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "sls",
					"status":        "online",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_esa_sites.default", "sites.#", "1"),
					testAccCheckEsaSiteDeliveryTaskSlsFixture(sls),
					testAccCheck(map[string]string{
						"status": "online",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "offline",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "offline",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "online",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "online",
					}),
				),
			},
			{
				// UpdateSiteDeliveryTask accepts BusinessType with HTTP 200 but the service keeps the
				// old value, so the apply cannot converge: Terraform reports a non-empty plan right
				// after the apply. Pin that exact diagnostic together with the attribute and both
				// values (state keeps the actual dcdn_log_access_l1, the config asks for
				// dcdn_log_er), so an unrelated apply-time error or another attribute cannot pass.
				Config: testAccConfig(map[string]interface{}{
					"business_type": "dcdn_log_er",
				}),
				ExpectError: regexp.MustCompile(`(?s)After applying this test step, the non-refresh plan was not empty.*~ business_type = "dcdn_log_access_l1" -> "dcdn_log_er"`),
			},
			{
				// Restore the original business_type (state never left the original value), so the
				// apply below is a no-op and the case converges.
				Config: testAccConfig(map[string]interface{}{
					"business_type": "dcdn_log_access_l1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"business_type": "dcdn_log_access_l1",
						"status":        "online",
					}),
				),
			},
			{
				// Same config as the previous step: re-applying must produce no change.
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "online",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sls_delivery"},
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_first_online_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_first_online_testBasicDependence(name string) string {
	return esaSiteDeliveryTaskBasicDependence(esaSiteDeliveryTaskSlsFixtureFromEnv())
}

// Test ESA SiteDeliveryTask first create with status offline. <<< Resource test cases, automatically generated.
// Case resource_SiteDeliveryTask_first_offline_test
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_first_offline_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_first_offline_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)
	sls := esaSiteDeliveryTaskSlsFixtureFromEnv()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_first_offline_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckEsaSiteDeliveryTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_delivery": []map[string]interface{}{
						{
							"sls_project":   sls.Project,
							"sls_region":    sls.Region,
							"sls_log_store": sls.LogStore,
						},
					},
					"site_id":       "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"data_center":   "global",
					"discard_rate":  "0.0",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "sls",
					"status":        "offline",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_esa_sites.default", "sites.#", "1"),
					testAccCheckEsaSiteDeliveryTaskSlsFixture(sls),
					testAccCheck(map[string]string{
						"status": "offline",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "offline",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sls_delivery"},
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_first_offline_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_first_offline_testBasicDependence(name string) string {
	return esaSiteDeliveryTaskBasicDependence(esaSiteDeliveryTaskSlsFixtureFromEnv())
}

// Case resource_SiteDeliveryTask_disappears_test
// 覆盖删除后回读语义：GetSiteDeliveryTask 对已删除任务返回 400 InternalError（非 404），因此
// Read 必须用完整成功的 ListSiteDeliveryTasks 才能判定任务消失并清理 state；同站点其它任务
// （非目标）不得受影响。已删除任务再次 Delete 只返回 400/406 错误，收敛判据在同文件的
// unit test 中覆盖。
func TestAccAliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_disappears_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_delivery_task.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_disappears_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteDeliveryTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sESASiteDeliveryTask%d", defaultRegionToTest, rand)
	bystander := fmt.Sprintf("tf-testAcc%sESASiteDeliveryTask%db", defaultRegionToTest, rand)
	sls := esaSiteDeliveryTaskSlsFixtureFromEnv()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_disappears_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckEsaSiteDeliveryTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_delivery": []map[string]interface{}{
						{
							"sls_project":   sls.Project,
							"sls_region":    sls.Region,
							"sls_log_store": sls.LogStore,
						},
					},
					"site_id":       "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"data_center":   "global",
					"discard_rate":  "0.0",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "sls",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_esa_sites.default", "sites.#", "1"),
					testAccCheckEsaSiteDeliveryTaskSlsFixture(sls),
					testAccCheck(map[string]string{}),
					testAccCheckEsaSiteDeliveryTaskListed(resourceId, true),
				),
			},
			{
				// 任务被带外删除后，Read 只能凭成功 List 证明其消失并清空 state，随后 plan 重新创建。
				Config: testAccConfig(map[string]interface{}{
					"sls_delivery": []map[string]interface{}{
						{
							"sls_project":   sls.Project,
							"sls_region":    sls.Region,
							"sls_log_store": sls.LogStore,
						},
					},
					"site_id":       "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"data_center":   "global",
					"discard_rate":  "0.0",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "sls",
				}),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccEsaSiteDeliveryTaskCreateOutOfBand(t, bystander, sls),
					testAccEsaSiteDeliveryTaskDeleteOutOfBand(resourceId),
					testAccCheckEsaSiteDeliveryTaskListed(resourceId, false),
					testAccCheckEsaSiteDeliveryTaskByNameListed(bystander, true),
				),
			},
			{
				// 目标任务被重新创建；同站点的非目标任务保持存在，最后带外清理旁观任务。
				Config: testAccConfig(map[string]interface{}{
					"sls_delivery": []map[string]interface{}{
						{
							"sls_project":   sls.Project,
							"sls_region":    sls.Region,
							"sls_log_store": sls.LogStore,
						},
					},
					"site_id":       "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"data_center":   "global",
					"discard_rate":  "0.0",
					"task_name":     name,
					"business_type": "dcdn_log_access_l1",
					"field_name":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
					"delivery_type": "sls",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
					testAccCheckEsaSiteDeliveryTaskListed(resourceId, true),
					testAccCheckEsaSiteDeliveryTaskByNameListed(bystander, true),
					testAccEsaSiteDeliveryTaskDeleteByNameOutOfBand(bystander),
					testAccCheckEsaSiteDeliveryTaskByNameListed(bystander, false),
					testAccCheckEsaSiteDeliveryTaskListed(resourceId, true),
				),
			},
		},
	})
}

var AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_disappears_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteDeliveryTaskresource_SiteDeliveryTask_disappears_testBasicDependence(name string) string {
	return esaSiteDeliveryTaskBasicDependence(esaSiteDeliveryTaskSlsFixtureFromEnv())
}

// Test ESA SiteDeliveryTask. <<< Resource test cases, automatically generated.

// testAccEsaSiteDeliveryTaskSinkParam 在 GetSiteDeliveryTask 返回的 SinkConfig 中定位指定参数集合。
// SinkConfig 在 OpenAPI 中为 any 类型，实际可能是对象或 JSON 字符串，两种形态都做兼容。
func testAccEsaSiteDeliveryTaskSinkParam(sinkConfig interface{}, paramName string) map[string]interface{} {
	switch value := sinkConfig.(type) {
	case string:
		var decoded interface{}
		if err := json.Unmarshal([]byte(value), &decoded); err == nil {
			return testAccEsaSiteDeliveryTaskSinkParam(decoded, paramName)
		}
	case []byte:
		return testAccEsaSiteDeliveryTaskSinkParam(string(value), paramName)
	case map[string]interface{}:
		for key, item := range value {
			if strings.EqualFold(key, paramName) {
				if params, ok := item.(map[string]interface{}); ok {
					return params
				}
			}
		}
		for _, item := range value {
			if params := testAccEsaSiteDeliveryTaskSinkParam(item, paramName); params != nil {
				return params
			}
		}
	}
	return nil
}

// testAccEsaSiteDeliveryTaskSinkParamValue 读取参数集合中指定 key 的保存值，只接受
// {"StaticValue":"..."} 形态：扁平字符串视为未按契约保存，避免修复前的扁平序列化也能通过。
func testAccEsaSiteDeliveryTaskSinkParamValue(params map[string]interface{}, key string) (string, bool) {
	for paramKey, paramValue := range params {
		if !strings.EqualFold(paramKey, key) {
			continue
		}
		switch value := paramValue.(type) {
		case map[string]interface{}:
			for staticKey, staticValueItem := range value {
				if strings.EqualFold(staticKey, "StaticValue") {
					text, ok := staticValueItem.(string)
					return text, ok
				}
			}
			return "", false
		case string:
			// 服务端可能把整个参数对象编码成 JSON 字符串；要求解码后仍是 StaticValue 形态。
			var decoded map[string]interface{}
			if err := json.Unmarshal([]byte(value), &decoded); err == nil {
				if staticItem, found := decoded["StaticValue"]; found {
					if text, ok := staticItem.(string); ok {
						return text, true
					}
				}
			}
			return "", false
		}
		return "", false
	}
	return "", false
}

// testAccCheckEsaSiteDeliveryTaskHttpSinkParams 通过只读查询接口断言 HTTP 投递的
// HeaderParam/QueryParam 值已被服务端实际保存，而不是仅证明 Create 被接受。取值助手只接受
// {"StaticValue":"..."} 形态（对象或等价的 JSON 字符串），因此该断言同时证明服务端保存的是
// 修复后的嵌套形态，扁平字符串不会通过。
func testAccCheckEsaSiteDeliveryTaskHttpSinkParams(resourceId, headerKey, headerValue, queryKey, queryValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		esaServiceV2 := EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
		object, err := esaServiceV2.DescribeEsaSiteDeliveryTask(rs.Primary.ID)
		if err != nil {
			return err
		}
		sinkConfig, ok := object["SinkConfig"]
		if !ok {
			return fmt.Errorf("GetSiteDeliveryTask response has no SinkConfig: %v", object)
		}
		headerParams := testAccEsaSiteDeliveryTaskSinkParam(sinkConfig, "HeaderParam")
		if savedValue, found := testAccEsaSiteDeliveryTaskSinkParamValue(headerParams, headerKey); !found || savedValue != headerValue {
			return fmt.Errorf("HeaderParam %q is not persisted as expected: found=%t, savedValue=%q, want %q, SinkConfig=%v", headerKey, found, savedValue, headerValue, sinkConfig)
		}
		queryParams := testAccEsaSiteDeliveryTaskSinkParam(sinkConfig, "QueryParam")
		if savedValue, found := testAccEsaSiteDeliveryTaskSinkParamValue(queryParams, queryKey); !found || savedValue != queryValue {
			return fmt.Errorf("QueryParam %q is not persisted as expected: found=%t, savedValue=%q, want %q, SinkConfig=%v", queryKey, found, savedValue, queryValue, sinkConfig)
		}
		return nil
	}
}

// testAccCheckEsaSiteDeliveryTaskDestroy 用完整成功的 ListSiteDeliveryTasks 证明任务已删除。
// GetSiteDeliveryTask 对不存在的任务返回 400 InternalError 而非 404，不能用于销毁判定；List
// 失败时删除状态即为未证实，直接报错而不当作已销毁。
func testAccCheckEsaSiteDeliveryTaskDestroy(s *terraform.State) error {
	esaServiceV2 := EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_esa_site_delivery_task" {
			continue
		}
		parts := strings.SplitN(rs.Primary.ID, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("unexpected id format %q", rs.Primary.ID)
		}
		exists, err := esaServiceV2.EsaSiteDeliveryTaskExists(parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("checking destroy of %s: ListSiteDeliveryTasks failed, absence is unproven: %v", rs.Primary.ID, err)
		}
		if exists {
			return fmt.Errorf("resource %s still exists in ListSiteDeliveryTasks", rs.Primary.ID)
		}
	}
	return nil
}

func testAccEsaSiteDeliveryTaskSiteId(s *terraform.State, resourceId string) (string, string, error) {
	rs, ok := s.RootModule().Resources[resourceId]
	if !ok {
		return "", "", fmt.Errorf("resource %s not found in state", resourceId)
	}
	parts := strings.SplitN(rs.Primary.ID, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected id format %q", rs.Primary.ID)
	}
	return parts[0], parts[1], nil
}

func testAccCheckEsaSiteDeliveryTaskListed(resourceId string, want bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		siteId, taskName, err := testAccEsaSiteDeliveryTaskSiteId(s, resourceId)
		if err != nil {
			return err
		}
		esaServiceV2 := EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
		exists, err := esaServiceV2.EsaSiteDeliveryTaskExists(siteId, taskName)
		if err != nil {
			return err
		}
		if exists != want {
			return fmt.Errorf("task %s:%s listed=%t, want %t", siteId, taskName, exists, want)
		}
		return nil
	}
}

func testAccCheckEsaSiteDeliveryTaskByNameListed(taskName string, want bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceId := "alicloud_esa_site_delivery_task.default"
		siteId, _, err := testAccEsaSiteDeliveryTaskSiteId(s, resourceId)
		if err != nil {
			return err
		}
		esaServiceV2 := EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
		exists, err := esaServiceV2.EsaSiteDeliveryTaskExists(siteId, taskName)
		if err != nil {
			return err
		}
		if exists != want {
			return fmt.Errorf("task %s:%s listed=%t, want %t", siteId, taskName, exists, want)
		}
		return nil
	}
}

// testAccEsaSiteDeliveryTaskCreateOutOfBand creates the bystander task outside Terraform state and
// registers its cleanup immediately, before any later check in the case can fail. The cleanup is
// bound to the site and task captured here, so it keeps working after the state read below is no
// longer available (for example once the target task has been deleted out of band and removed from
// state). The out-of-band task delivers into the same verified SLS fixture as the resource under
// test.
func testAccEsaSiteDeliveryTaskCreateOutOfBand(t *testing.T, taskName string, sls esaSiteDeliveryTaskSlsFixture) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		siteId, _, err := testAccEsaSiteDeliveryTaskSiteId(s, "alicloud_esa_site_delivery_task.default")
		if err != nil {
			return err
		}
		sinkConfig, err := json.Marshal(map[string]string{
			"SLSProject":  sls.Project,
			"SLSRegion":   sls.Region,
			"SLSLogStore": sls.LogStore,
		})
		if err != nil {
			return fmt.Errorf("encoding SlsDelivery for the out-of-band task failed: %v", err)
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		request := map[string]interface{}{
			"SiteId":       siteId,
			"TaskName":     taskName,
			"DataCenter":   "global",
			"DeliveryType": "sls",
			"BusinessType": "dcdn_log_access_l1",
			"DiscardRate":  0.0,
			"FieldName":    "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID",
			"SlsDelivery":  string(sinkConfig),
		}
		if _, err := client.RpcPost("ESA", "2024-09-10", "CreateSiteDeliveryTask", map[string]interface{}{}, request, true); err != nil {
			return fmt.Errorf("creating bystander task %s failed: %v", taskName, err)
		}
		testAccEsaSiteDeliveryTaskRegisterOutOfBandCleanup(t, siteId, taskName)
		return nil
	}
}

// esaSiteDeliveryTaskCleanupRegistrar is the subset of *testing.T the out-of-band cleanup needs, so
// the registration and the failure-path behaviour can be unit tested without a live client.
type esaSiteDeliveryTaskCleanupRegistrar interface {
	Helper()
	Cleanup(func())
	Errorf(format string, args ...interface{})
}

// testAccEsaSiteDeliveryTaskRegisterOutOfBandCleanup is the acceptance-test wiring of
// registerEsaSiteDeliveryTaskOutOfBandCleanup with the live ESA client.
func testAccEsaSiteDeliveryTaskRegisterOutOfBandCleanup(t *testing.T, siteId, taskName string) {
	esaServiceV2 := EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	registerEsaSiteDeliveryTaskOutOfBandCleanup(t, siteId, taskName, esaServiceV2.EsaSiteDeliveryTaskExists, func() error {
		return testAccEsaSiteDeliveryTaskDelete(siteId, taskName)
	})
}

// registerEsaSiteDeliveryTaskOutOfBandCleanup registers a t.Cleanup that removes a task created
// outside Terraform state. Cleanup registered through t.Cleanup still runs when a later check,
// apply or destroy fails, so an interrupted case cannot leave the bystander on the shared site.
// Removal is idempotent: a fully successful ListSiteDeliveryTasks that no longer lists the exact
// task is the only proof of absence and skips the delete, and a delete error is only tolerated when
// a later successful List proves the task gone. Any other outcome is reported as a test error
// instead of being swallowed, and the exact site/task captured at creation is used, so other tasks
// on the site are never probed for deletion.
func registerEsaSiteDeliveryTaskOutOfBandCleanup(t esaSiteDeliveryTaskCleanupRegistrar, siteId, taskName string, exists func(siteId, taskName string) (bool, error), remove func() error) {
	t.Helper()
	t.Cleanup(func() {
		if err := esaSiteDeliveryTaskDeleteWithApis(siteId+":"+taskName, siteId, taskName, exists, remove); err != nil {
			t.Errorf("out-of-band task %s:%s could not be removed: %v", siteId, taskName, err)
		}
	})
}

// esaSiteDeliveryTaskCleanupRecorder stands in for *testing.T in the unit tests: it records the
// registered cleanup so the tests can invoke it the way t.Cleanup does after a failed case.
type esaSiteDeliveryTaskCleanupRecorder struct {
	cleanups []func()
	errors   []string
}

func (r *esaSiteDeliveryTaskCleanupRecorder) Helper() {}

func (r *esaSiteDeliveryTaskCleanupRecorder) Cleanup(fn func()) {
	r.cleanups = append(r.cleanups, fn)
}

func (r *esaSiteDeliveryTaskCleanupRecorder) Errorf(format string, args ...interface{}) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}

func (r *esaSiteDeliveryTaskCleanupRecorder) runCleanups() {
	for _, fn := range r.cleanups {
		fn()
	}
}

func testAccEsaSiteDeliveryTaskDeleteOutOfBand(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		siteId, taskName, err := testAccEsaSiteDeliveryTaskSiteId(s, resourceId)
		if err != nil {
			return err
		}
		return testAccEsaSiteDeliveryTaskDelete(siteId, taskName)
	}
}

func testAccEsaSiteDeliveryTaskDeleteByNameOutOfBand(taskName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		siteId, _, err := testAccEsaSiteDeliveryTaskSiteId(s, "alicloud_esa_site_delivery_task.default")
		if err != nil {
			return err
		}
		return testAccEsaSiteDeliveryTaskDelete(siteId, taskName)
	}
}

func testAccEsaSiteDeliveryTaskDelete(siteId, taskName string) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"SiteId":   siteId,
		"TaskName": taskName,
	}
	if _, err := client.RpcPost("ESA", "2024-09-10", "DeleteSiteDeliveryTask", map[string]interface{}{}, request, true); err != nil {
		return fmt.Errorf("out-of-band delete of %s:%s failed: %v", siteId, taskName, err)
	}
	return nil
}

// TestUnitAliCloudEsaSiteDeliveryTask 覆盖 List 驱动的存在性判定与删除收敛（不再用
// GetSiteDeliveryTask 做存在性探测）：分页、目标在后续页、List 失败不得当作不存在、
// 多任务站点保留非目标，以及删除错误只在 List 前后两次都成功且证明任务消失时收敛。
func TestUnitAliCloudEsaSiteDeliveryTask(t *testing.T) {
	const siteId = "185447364102380"
	const taskName = "task"

	service := &EsaServiceV2{}
	taskPage := func(names ...string) map[string]interface{} {
		tasks := make([]interface{}, 0, len(names))
		for _, name := range names {
			tasks = append(tasks, map[string]interface{}{"TaskName": name})
		}
		return map[string]interface{}{"Tasks": tasks, "TotalCount": len(tasks)}
	}

	t.Run("ListPaginationFindsTargetOnLaterPage", func(t *testing.T) {
		var pages []int
		fetch := func(gotSiteId string, page int) (map[string]interface{}, error) {
			assert.Equal(t, siteId, gotSiteId)
			pages = append(pages, page)
			switch page {
			case 1:
				names := make([]string, 0, PageSizeLarge)
				for i := 0; i < PageSizeLarge; i++ {
					names = append(names, fmt.Sprintf("other-%d", i))
				}
				return taskPage(names...), nil
			case 2:
				return taskPage(taskName), nil
			}
			return taskPage(), nil
		}

		tasks, err := service.listEsaSiteDeliveryTasks(siteId, fetch)
		assert.NoError(t, err)
		assert.Len(t, tasks, PageSizeLarge+1)
		assert.Equal(t, []int{1, 2}, pages)
		assert.True(t, esaSiteDeliveryTaskListed(tasks, taskName))
	})

	t.Run("ListPaginationReportsAbsenceAfterFullTraversal", func(t *testing.T) {
		var pages []int
		fetch := func(gotSiteId string, page int) (map[string]interface{}, error) {
			assert.Equal(t, siteId, gotSiteId)
			pages = append(pages, page)
			if page == 1 {
				names := make([]string, 0, PageSizeLarge)
				for i := 0; i < PageSizeLarge; i++ {
					names = append(names, fmt.Sprintf("other-%d", i))
				}
				return taskPage(names...), nil
			}
			return taskPage("other-on-last-page"), nil
		}

		tasks, err := service.listEsaSiteDeliveryTasks(siteId, fetch)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2}, pages, "every page must be requested until a short page ends the walk")
		assert.Len(t, tasks, PageSizeLarge+1)
		assert.False(t, esaSiteDeliveryTaskListed(tasks, taskName), "a completed full traversal without the task proves absence")
	})

	t.Run("ListFailureIsNotAbsence", func(t *testing.T) {
		fetch := func(string, int) (map[string]interface{}, error) {
			return nil, fmt.Errorf("code: 400, list denied")
		}
		tasks, err := service.listEsaSiteDeliveryTasks(siteId, fetch)
		assert.Error(t, err)
		assert.Nil(t, tasks)
		assert.False(t, esaSiteDeliveryTaskListed(tasks, taskName))
	})

	t.Run("ListMalformedPageIsError", func(t *testing.T) {
		cases := []struct {
			name     string
			response map[string]interface{}
		}{
			{"missing Tasks field", map[string]interface{}{"TotalCount": 0}},
			{"Tasks is not an array", map[string]interface{}{"Tasks": "not-a-list"}},
			{"Tasks contains a non-object element", map[string]interface{}{"Tasks": []interface{}{"not-an-object"}}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				fetch := func(string, int) (map[string]interface{}, error) { return tc.response, nil }
				tasks, err := service.listEsaSiteDeliveryTasks(siteId, fetch)
				assert.Error(t, err)
				assert.Nil(t, tasks)
			})
		}
	})

	t.Run("ListLaterPageParseFailureIsError", func(t *testing.T) {
		fetch := func(_ string, page int) (map[string]interface{}, error) {
			if page == 1 {
				names := make([]string, 0, PageSizeLarge)
				for i := 0; i < PageSizeLarge; i++ {
					names = append(names, fmt.Sprintf("other-%d", i))
				}
				return taskPage(names...), nil
			}
			return map[string]interface{}{"TotalCount": 0}, nil
		}
		tasks, err := service.listEsaSiteDeliveryTasks(siteId, fetch)
		assert.Error(t, err)
		assert.Nil(t, tasks)
	})

	t.Run("ListedMatchesExactTaskAndKeepsOtherTasks", func(t *testing.T) {
		fetch := func(string, int) (map[string]interface{}, error) {
			return taskPage("other-1", "other-2", taskName), nil
		}
		tasks, err := service.listEsaSiteDeliveryTasks(siteId, fetch)
		assert.NoError(t, err)
		assert.True(t, esaSiteDeliveryTaskListed(tasks, taskName))
		assert.True(t, esaSiteDeliveryTaskListed(tasks, "other-1"))
		assert.False(t, esaSiteDeliveryTaskListed(tasks, "absent"))
	})

	t.Run("DeleteConvergesOnlyWhenListProvesAbsence", func(t *testing.T) {
		listErr := fmt.Errorf("code: 400, list denied")
		cases := []struct {
			name     string
			after    bool
			afterErr error
			want     bool
		}{
			{"the task is gone after the failed delete", false, nil, true},
			{"the task is still listed", true, nil, false},
			{"the post-delete list failed", false, listErr, false},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				assert.Equal(t, tc.want, esaSiteDeliveryTaskDeleteConverged(tc.after, tc.afterErr))
			})
		}
	})

	t.Run("OfflineTaskStillCountsAsExisting", func(t *testing.T) {
		fetch := func(string, int) (map[string]interface{}, error) {
			return map[string]interface{}{
				"Tasks": []interface{}{
					map[string]interface{}{"TaskName": taskName, "Status": "offline"},
				},
				"TotalCount": 1,
			}, nil
		}
		tasks, err := service.listEsaSiteDeliveryTasks(siteId, fetch)
		assert.NoError(t, err)
		assert.True(t, esaSiteDeliveryTaskListed(tasks, taskName))
	})

	t.Run("ListRetryGateRetriesInternalErrorOnly", func(t *testing.T) {
		sdkError := func(code, message string) error {
			return &tea.SDKError{
				Code:       String(code),
				Data:       String(code),
				Message:    String(message),
				StatusCode: tea.Int(400),
			}
		}
		assert.True(t, esaSiteDeliveryTaskListRetryable(sdkError("InternalError", "code: 400, An internal exception occurred, please try again later.")))
		assert.False(t, esaSiteDeliveryTaskListRetryable(sdkError("InternalError", "code: 400, Specified parameter SiteId is not valid.")))
		assert.False(t, esaSiteDeliveryTaskListRetryable(sdkError("InvalidParameter", "code: 400, Specified parameter SiteId is not valid.")))
		assert.False(t, esaSiteDeliveryTaskListRetryable(nil))
	})

	t.Run("WriteRetryGateRetriesInternalErrorOnly", func(t *testing.T) {
		sdkError := func(code, message string) error {
			return &tea.SDKError{
				Code:       String(code),
				Data:       String(code),
				Message:    String(message),
				StatusCode: tea.Int(400),
			}
		}
		assert.True(t, esaSiteDeliveryTaskWriteRetryable(sdkError("InternalError", "code: 400, An internal exception occurred, please try again later.")))
		assert.False(t, esaSiteDeliveryTaskWriteRetryable(sdkError("InternalError", "code: 400, Specified parameter TaskName is not valid.")))
		assert.False(t, esaSiteDeliveryTaskWriteRetryable(sdkError("InvalidParameter", "code: 400, Specified parameter is not valid.")))
		assert.False(t, esaSiteDeliveryTaskWriteRetryable(nil))
	})

	t.Run("HttpHeaderParamSerializesStaticValue", func(t *testing.T) {
		got := esaSiteDeliveryTaskHttpParamValues(map[string]interface{}{
			"x-auth": "test-header-value",
		})
		raw, err := json.Marshal(got)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"x-auth":{"StaticValue":"test-header-value"}}`, string(raw))
	})

	t.Run("HttpQueryParamSerializesStaticValue", func(t *testing.T) {
		got := esaSiteDeliveryTaskHttpParamValues(map[string]interface{}{
			"auth_token": "test-query-value",
		})
		raw, err := json.Marshal(got)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"auth_token":{"StaticValue":"test-query-value"}}`, string(raw))
	})

	t.Run("DeleteSequence", func(t *testing.T) {
		repeatDelete406 := fmt.Errorf("code: 406, get project error")
		listErr := fmt.Errorf("code: 400, list denied")

		t.Run("RepeatDelete406ConvergesWhenListsProveAbsence", func(t *testing.T) {
			probes := 0
			exists := func(gotSiteId, gotTaskName string) (bool, error) {
				assert.Equal(t, siteId, gotSiteId)
				assert.Equal(t, taskName, gotTaskName)
				probes++
				return probes == 1, nil
			}
			err := esaSiteDeliveryTaskDeleteWithApis("id", siteId, taskName, exists, func() error {
				return repeatDelete406
			})
			assert.NoError(t, err)
			assert.Equal(t, 2, probes)
		})

		t.Run("RepeatDelete406AndPostListFailureAreBothKept", func(t *testing.T) {
			probes := 0
			exists := func(string, string) (bool, error) {
				probes++
				if probes == 1 {
					return true, nil
				}
				return false, listErr
			}
			err := esaSiteDeliveryTaskDeleteWithApis("id", siteId, taskName, exists, func() error {
				return repeatDelete406
			})
			assert.ErrorIs(t, err, repeatDelete406, "the delete error must stay observable")
			assert.Contains(t, err.Error(), "list denied", "the post-delete List error must not be dropped")
		})

		t.Run("DeleteErrorIsKeptWhenTaskIsStillListed", func(t *testing.T) {
			exists := func(string, string) (bool, error) { return true, nil }
			deleteErr := fmt.Errorf("code: 500, backend failure")
			err := esaSiteDeliveryTaskDeleteWithApis("id", siteId, taskName, exists, func() error {
				return deleteErr
			})
			assert.ErrorIs(t, err, deleteErr)
		})

		t.Run("AbsentTaskSkipsDelete", func(t *testing.T) {
			called := false
			err := esaSiteDeliveryTaskDeleteWithApis("id", siteId, taskName,
				func(string, string) (bool, error) { return false, nil },
				func() error { called = true; return nil })
			assert.NoError(t, err)
			assert.False(t, called)
		})

		t.Run("PreListFailureIsReturnedAndSkipsDelete", func(t *testing.T) {
			removed := false
			err := esaSiteDeliveryTaskDeleteWithApis("id", siteId, taskName,
				func(string, string) (bool, error) { return false, listErr },
				func() error { removed = true; return nil })
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "list denied")
			assert.False(t, removed, "a failed pre-delete List must not delete, even when remove would succeed")
		})

		// resourceAliCloudEsaSiteDeliveryTaskDelete wraps whatever the helper returns; both failures
		// have to survive that wrap so the operator can tell a failed delete from a failed proof.
		t.Run("ResourceDeleteWrapKeepsBothErrors", func(t *testing.T) {
			probes := 0
			err := esaSiteDeliveryTaskDeleteWithApis("id", siteId, taskName, func(string, string) (bool, error) {
				probes++
				if probes == 1 {
					return true, nil
				}
				return false, listErr
			}, func() error {
				return repeatDelete406
			})
			wrapped := WrapErrorf(err, DefaultErrorMsg, "id", "DeleteSiteDeliveryTask", AlibabaCloudSdkGoERROR)
			assert.Contains(t, wrapped.Error(), "get project error", "the delete failure must stay visible after the resource wrap")
			assert.Contains(t, wrapped.Error(), "list denied", "the post-delete List failure must stay visible after the resource wrap")
		})
	})

	// The out-of-band bystander is created after the apply and never enters Terraform state, so the
	// failure half of the disappears case is what matters: the cleanup must be registered before a
	// later check can fail, remove only the exact task it created, and report real failures instead
	// of passing silently.
	t.Run("OutOfBandBystanderCleanup", func(t *testing.T) {
		const bystander = "bystander-task"
		const otherTask = "someone-elses-task"

		t.Run("RunsWhenTheCaseFailsMidwayAndLeavesOtherTasksAlone", func(t *testing.T) {
			registrar := &esaSiteDeliveryTaskCleanupRecorder{}
			probed := make([]string, 0, 1)
			removed := 0
			registerEsaSiteDeliveryTaskOutOfBandCleanup(registrar, siteId, bystander,
				func(gotSiteId, gotTaskName string) (bool, error) {
					assert.Equal(t, siteId, gotSiteId)
					probed = append(probed, gotTaskName)
					return gotTaskName == bystander, nil
				},
				func() error {
					removed++
					return nil
				})
			assert.Len(t, registrar.cleanups, 1, "cleanup must be registered before a later check can fail")

			// t.Cleanup runs after a failed case as well; invoking the registered func is that path.
			registrar.runCleanups()
			assert.Equal(t, []string{bystander}, probed)
			assert.Equal(t, 1, removed)
			assert.NotContains(t, probed, otherTask)
			assert.Empty(t, registrar.errors)
		})

		t.Run("ConvergesWhenALaterListProvesTheTaskGone", func(t *testing.T) {
			registrar := &esaSiteDeliveryTaskCleanupRecorder{}
			probes := 0
			registerEsaSiteDeliveryTaskOutOfBandCleanup(registrar, siteId, bystander,
				func(string, string) (bool, error) {
					probes++
					return probes == 1, nil
				},
				func() error { return fmt.Errorf("code: 406, get project error") })
			registrar.runCleanups()
			assert.Equal(t, 2, probes)
			assert.Empty(t, registrar.errors)
		})

		t.Run("ReportsDeleteFailureWhenTheTaskIsStillListed", func(t *testing.T) {
			registrar := &esaSiteDeliveryTaskCleanupRecorder{}
			registerEsaSiteDeliveryTaskOutOfBandCleanup(registrar, siteId, bystander,
				func(string, string) (bool, error) { return true, nil },
				func() error { return fmt.Errorf("code: 500, backend failure") })
			registrar.runCleanups()
			assert.Len(t, registrar.errors, 1)
			assert.Contains(t, registrar.errors[0], "backend failure")
		})

		t.Run("ReportsWhenAbsenceCannotBeProven", func(t *testing.T) {
			registrar := &esaSiteDeliveryTaskCleanupRecorder{}
			registerEsaSiteDeliveryTaskOutOfBandCleanup(registrar, siteId, bystander,
				func(string, string) (bool, error) { return false, fmt.Errorf("code: 400, list denied") },
				func() error { return fmt.Errorf("code: 406, get project error") })
			registrar.runCleanups()
			assert.Len(t, registrar.errors, 1)
		})

		t.Run("ReportsThePreListErrorEvenWhenRemoveSucceeds", func(t *testing.T) {
			registrar := &esaSiteDeliveryTaskCleanupRecorder{}
			removed := false
			registerEsaSiteDeliveryTaskOutOfBandCleanup(registrar, siteId, bystander,
				func(string, string) (bool, error) { return false, fmt.Errorf("code: 400, list denied") },
				func() error { removed = true; return nil })
			registrar.runCleanups()
			assert.False(t, removed, "the cleanup must not delete what it could not List")
			assert.Len(t, registrar.errors, 1)
			assert.Contains(t, registrar.errors[0], "list denied")
		})

		t.Run("ReportsDeleteAndPostListFailures", func(t *testing.T) {
			registrar := &esaSiteDeliveryTaskCleanupRecorder{}
			probes := 0
			registerEsaSiteDeliveryTaskOutOfBandCleanup(registrar, siteId, bystander,
				func(string, string) (bool, error) {
					probes++
					if probes == 1 {
						return true, nil
					}
					return false, fmt.Errorf("code: 400, list denied")
				},
				func() error { return fmt.Errorf("code: 500, backend failure") })
			registrar.runCleanups()
			assert.Len(t, registrar.errors, 1)
			assert.Contains(t, registrar.errors[0], "backend failure")
			assert.Contains(t, registrar.errors[0], "list denied", "the cleanup must report why absence could not be proven")
		})
	})

	// The SLS fixture is a controlled input: the dedicated fixture is the fallback and each field
	// can be overridden independently by the environment.
	t.Run("SlsFixtureOverrideOrDefault", func(t *testing.T) {
		t.Setenv(esaSiteDeliveryTaskSlsProjectEnv, "")
		t.Setenv(esaSiteDeliveryTaskSlsLogStoreEnv, "")
		t.Setenv(esaSiteDeliveryTaskSlsRegionEnv, "")
		assert.Equal(t, esaSiteDeliveryTaskSlsFixture{
			Project:  esaSiteDeliveryTaskSlsDefaultProject,
			LogStore: esaSiteDeliveryTaskSlsDefaultLogStore,
			Region:   esaSiteDeliveryTaskSlsDefaultRegion,
		}, esaSiteDeliveryTaskSlsFixtureFromEnv())

		t.Setenv(esaSiteDeliveryTaskSlsProjectEnv, "sls-project-from-env")
		assert.Equal(t, "sls-project-from-env", esaSiteDeliveryTaskSlsFixtureFromEnv().Project)
		assert.Equal(t, esaSiteDeliveryTaskSlsDefaultLogStore, esaSiteDeliveryTaskSlsFixtureFromEnv().LogStore)

		t.Setenv(esaSiteDeliveryTaskSlsRegionEnv, "cn-shanghai")
		assert.Equal(t, "cn-shanghai", esaSiteDeliveryTaskSlsFixtureFromEnv().Region)
	})
}
