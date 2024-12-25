package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA Record. >>> Resource test cases, automatically generated.
// Case resource_Record_srv_test
func TestAccAliCloudESARecordresource_Record_srv_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_record.default"
	ra := resourceAttrInit(resourceId, AliCloudESARecordresource_Record_srv_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRecord")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESARecord%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARecordresource_Record_srv_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"record_name": "_udp._sip.idltestrecord.com",
					"comment":     "This is a remark",
					"site_id":     "${alicloud_esa_site.resource_Site_srv_test.id}",
					"record_type": "SRV",
					"data": []map[string]interface{}{
						{
							"priority": "1",
							"port":     "80",
							"value":    "www.eerrraaa.com",
							"weight":   "1",
						},
					},
					"ttl": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "test_record_comment",
					"data": []map[string]interface{}{
						{
							"priority": "2",
							"port":     "8080",
							"value":    "www.qwer.com",
							"weight":   "2",
						},
					},
					"ttl": "86400",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESARecordresource_Record_srv_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESARecordresource_Record_srv_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_srv_test" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_srv_test" {
  site_name   = "idltestrecord.com"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_srv_test.id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_Record_smimea_test
func TestAccAliCloudESARecordresource_Record_smimea_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_record.default"
	ra := resourceAttrInit(resourceId, AliCloudESARecordresource_Record_smimea_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRecord")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESARecord%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARecordresource_Record_smimea_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"record_name": "www.idltestrecord.com",
					"comment":     "This is a remark",
					"site_id":     "${alicloud_esa_site.resource_Site_smimea_test.id}",
					"record_type": "SMIMEA",
					"data": []map[string]interface{}{
						{
							"usage":         "1",
							"matching_type": "1",
							"selector":      "1",
							"certificate":   "7777276264696475536f6d313237",
						},
					},
					"ttl": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "test_record_comment",
					"data": []map[string]interface{}{
						{
							"usage":         "3",
							"matching_type": "3",
							"selector":      "3",
							"certificate":   "7737246264656475536f6d617256",
						},
					},
					"ttl": "86400",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESARecordresource_Record_smimea_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESARecordresource_Record_smimea_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_smimea_test" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_smimea_test" {
  site_name   = "idltestrecord.com"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_smimea_test.id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_Record_test_cname
func TestAccAliCloudESARecordresource_Record_test_cname(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_record.default"
	ra := resourceAttrInit(resourceId, AliCloudESARecordresource_Record_test_cnameMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRecord")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESARecord%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARecordresource_Record_test_cnameBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"record_name": "www.idltestrecord.com",
					"comment":     "This is a remark",
					"proxied":     "true",
					"site_id":     "${alicloud_esa_site.resource_Site_test.id}",
					"record_type": "CNAME",
					"source_type": "S3",
					"data": []map[string]interface{}{
						{
							"value": "www.idltestr.com",
						},
					},
					"biz_name":    "api",
					"host_policy": "follow_hostname",
					"ttl":         "100",
					"auth_conf": []map[string]interface{}{
						{
							"secret_key": "hijklmnhijklmnhijklmnhijklmn",
							"version":    "v4",
							"region":     "us-east-1",
							"auth_type":  "private",
							"access_key": "abcdefgabcdefgabcdefgabcdefg",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"biz_name": "web",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"value": "www.pangleitestupdate.com",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ttl": "3600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "DNS记录测试",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type": "OSS",
					"auth_conf": []map[string]interface{}{
						{
							"secret_key": "secretkey1234567890abcdefghijklmn",
							"version":    "v2",
							"region":     "us-east-2",
							"auth_type":  "public",
							"access_key": "AccessKey1234567890abcdefghijklmn",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_policy": "follow_origin_domain",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment":     "test_record_comment",
					"proxied":     "true",
					"source_type": "S3",
					"biz_name":    "api",
					"data": []map[string]interface{}{
						{
							"value": "www.plexample.com",
						},
					},
					"host_policy": "follow_origin_domain",
					"ttl":         "86400",
					"auth_conf": []map[string]interface{}{
						{
							"secret_key": "secretkey0987654321fedcbafedcba",
							"version":    "v2",
							"region":     "us-gov-west-1",
							"auth_type":  "private",
							"access_key": "AccessKey0987654321fedcbafedcba",
						},
					},
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
				ImportStateVerifyIgnore: []string{"auth_conf.0.secret_key", "auth_conf.0.access_key"},
			},
		},
	})
}

var AliCloudESARecordresource_Record_test_cnameMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESARecordresource_Record_test_cnameBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_test" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_test" {
  site_name   = "idltestrecord.com"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_test.id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_Record_cert_test
func TestAccAliCloudESARecordresource_Record_cert_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_record.default"
	ra := resourceAttrInit(resourceId, AliCloudESARecordresource_Record_cert_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRecord")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESARecord%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARecordresource_Record_cert_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"record_name": "www.idltestrecord.com",
					"comment":     "This is a remark",
					"site_id":     "${alicloud_esa_site.resource_Site_cert_test.id}",
					"record_type": "CERT",
					"data": []map[string]interface{}{
						{
							"type":        "111",
							"key_tag":     "11",
							"algorithm":   "11",
							"certificate": "eGVzdGFsbWNkbg==",
						},
					},
					"ttl": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "test_record_comment",
					"data": []map[string]interface{}{
						{
							"type":        "222",
							"key_tag":     "22",
							"algorithm":   "22",
							"certificate": "bGVzdGGsbWNkbg==",
						},
					},
					"ttl": "86400",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESARecordresource_Record_cert_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESARecordresource_Record_cert_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_cert_test" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_cert_test" {
  site_name   = "idltestrecord.com"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_cert_test.id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_Record_sshfp_test
func TestAccAliCloudESARecordresource_Record_sshfp_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_record.default"
	ra := resourceAttrInit(resourceId, AliCloudESARecordresource_Record_sshfp_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRecord")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESARecord%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARecordresource_Record_sshfp_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"record_name": "www.idltestrecord.com",
					"comment":     "This is a remark",
					"site_id":     "${alicloud_esa_site.resource_Site_sshfp_test.id}",
					"record_type": "SSHFP",
					"data": []map[string]interface{}{
						{
							"type":        "1",
							"fingerprint": "6262626475636f6d",
							"algorithm":   "1",
						},
					},
					"ttl": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "test_record_comment",
					"data": []map[string]interface{}{
						{
							"type":        "3",
							"fingerprint": "6464646475636f6d",
							"algorithm":   "3",
						},
					},
					"ttl": "86400",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESARecordresource_Record_sshfp_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESARecordresource_Record_sshfp_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_sshfp_test" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_sshfp_test" {
  site_name   = "idltestrecord.com"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_sshfp_test.id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_Record_caa_test
func TestAccAliCloudESARecordresource_Record_caa_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_record.default"
	ra := resourceAttrInit(resourceId, AliCloudESARecordresource_Record_caa_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRecord")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESARecord%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARecordresource_Record_caa_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"record_name": "www.idltestrecord.com",
					"comment":     "This is a remark",
					"site_id":     "${alicloud_esa_site.resource_Site_caa_test.id}",
					"record_type": "CAA",
					"data": []map[string]interface{}{
						{
							"value": "www.eerrraaa.com",
							"tag":   "issue",
							"flag":  "1",
						},
					},
					"ttl": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "test_record_comment",
					"data": []map[string]interface{}{
						{
							"value": "www.qwer.com",
							"tag":   "issuewild",
							"flag":  "1",
						},
					},
					"ttl": "86400",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESARecordresource_Record_caa_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESARecordresource_Record_caa_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_caa_test" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_caa_test" {
  site_name   = "idltestrecord.com"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_caa_test.id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA Record. <<< Resource test cases, automatically generated.
