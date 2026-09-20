package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAlicloudSLBAclsDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_slb_acl.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_slb_acl.default.name}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_slb_acl.default.name}"`,
			"tags":       `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_slb_acl.default.name}"`,
			"tags":       `{Created = "TF1"}`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_slb_acl.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_slb_acl.default.id}_fake"]`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_slb_acl.default.id}"]`,
			"resource_group_id": `""`,
		}),
		fakeConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_slb_acl.default.id}_fake"]`,
			"resource_group_id": fmt.Sprintf(`"%s_fake"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_slb_acl.default.id}"]`,
			"name_regex": `"${alicloud_slb_acl.default.name}"`,
			// The resource route tables do not support resource_group_id, so it was set empty.
			"resource_group_id": `""`,
		}),
		fakeConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_slb_acl.default.id}_fake"]`,
			"name_regex":        `"${alicloud_slb_acl.default.name}"`,
			"resource_group_id": `""`,
		}),
	}

	var existSLBAclsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":                     "1",
			"ids.#":                      "1",
			"names.#":                    "1",
			"acls.0.id":                  CHECKSET,
			"acls.0.resource_group_id":   CHECKSET,
			"acls.0.name":                fmt.Sprintf("tf-testAccSlbAclDataSourceBisic-%d", rand),
			"acls.0.ip_version":          "ipv4",
			"acls.0.entry_list.#":        "2",
			"acls.0.related_listeners.#": "0",
			"acls.0.tags.%":              "2",
			"acls.0.tags.Created":        "TF",
			"acls.0.tags.For":            "acceptance test",
		}
	}

	var fakeSLBAclsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":  "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var slbaclsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_acls.default",
		existMapFunc: existSLBAclsMapFunc,
		fakeMapFunc:  fakeSLBAclsMapFunc,
	}

	slbaclsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, tagsConf, idsConf, resourceGroupIdConf, allConf)
}

func testAccCheckAlicloudSlbAclsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSlbAclDataSourceBisic-%d"
}
variable "ip_version" {
	default = "ipv4"
}

resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
    entry = "10.10.10.0/24"
    comment = "first"
  }
  entry_list {
      entry = "168.10.10.0/24"
      comment = "second"
  }
   tags = {
      Created = "TF"
       For     = "acceptance test"
    }
}


data "alicloud_slb_acls" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

// TestAccAlicloudSLBAclsDataSource_multipleAcls exercises the multi-ACL loop in
// slbAclsDescriptionAttributes where the ListTagResources error-swallowing bug
// previously left the data source returning success with a half-complete state.
// Both ACLs carry tags; the data source must surface both ACLs, both id/name
// lists, and per-ACL tags (d.Set("acls", s)/d.Set("ids", ids)/d.Set("names", names)).
func TestAccAlicloudSLBAclsDataSource_multipleAcls(t *testing.T) {
	rand := acctest.RandInt()
	secondTagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAclsMultipleConfig(rand),
	}
	allAttrs := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.alicloud_slb_acls.default", "acls.#", "2"),
		resource.TestCheckResourceAttr("data.alicloud_slb_acls.default", "ids.#", "2"),
		resource.TestCheckResourceAttr("data.alicloud_slb_acls.default", "names.#", "2"),
		// Both ACLs carry two tags; the per-index ordering is not guaranteed but
		// each returned ACL must expose its two tags rather than an empty set.
		resource.TestCheckResourceAttrSet("data.alicloud_slb_acls.default", "acls.0.id"),
		resource.TestCheckResourceAttrSet("data.alicloud_slb_acls.default", "acls.1.id"),
		resource.TestCheckResourceAttr("data.alicloud_slb_acls.default", "acls.0.tags.%", "2"),
		resource.TestCheckResourceAttr("data.alicloud_slb_acls.default", "acls.1.tags.%", "2"),
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: secondTagsConf.existConfig,
				Check:  allAttrs,
			},
		},
	})
}

func testAccCheckAlicloudSlbAclsMultipleConfig(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSlbAclDataSourceMulti-%d"
}
variable "ip_version" {
  default = "ipv4"
}

resource "alicloud_slb_acl" "first" {
  name = "${var.name}-1"
  ip_version = "${var.ip_version}"
  entry_list {
    entry = "10.10.10.0/24"
    comment = "first"
  }
  tags = {
    Created = "TF"
    Order   = "first"
  }
}

resource "alicloud_slb_acl" "second" {
  name = "${var.name}-2"
  ip_version = "${var.ip_version}"
  entry_list {
    entry = "168.10.10.0/24"
    comment = "second"
  }
  tags = {
    Created = "TF"
    Order   = "second"
  }
}

data "alicloud_slb_acls" "default" {
  ids = ["${alicloud_slb_acl.first.id}", "${alicloud_slb_acl.second.id}"]
}
`, rand)
}

// TestUnitAlicloudSlbAcls_listTagsResourcesError verifies the fix for the bug
// where the alicloud_slb_acls data source swallowed the ListTagResources error
// (returned nil) and reported the data source as successful with a
// half-complete state (missing tags/ids/names/acls). With the fix, the error
// must be propagated. Mocks the SLB client via gomonkey so no real cloud calls
// are made; skipped when the shared test client cannot be constructed.
func TestUnitAlicloudSlbAcls_listTagsResourcesError(t *testing.T) {
	p := Provider().(*schema.Provider).DataSourcesMap
	d, _ := schema.InternalMap(p["alicloud_slb_acls"].Schema).Data(nil, nil)

	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)

	// Single ACL whose tags query fails: the data source must surface the error
	// instead of returning success with an incomplete state.
	t.Run("SingleAclTagsErrorPropagated", func(t *testing.T) {
		callCount := 0
		patchWithSlbClient := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "WithSlbClient",
			func(_ *connectivity.AliyunClient, do func(*slb.Client) (interface{}, error)) (interface{}, error) {
				callCount++
				if callCount == 1 {
					resp := slb.CreateDescribeAccessControlListsResponse()
					resp.Acls = slb.Acls{Acl: []slb.Acl{{AclId: "acl-test-1", AclName: "tf-test-acl-1"}}}
					return resp, nil
				}
				attrResp := slb.CreateDescribeAccessControlListAttributeResponse()
				attrResp.AclId = "acl-test-1"
				attrResp.AclName = "tf-test-acl-1"
				attrResp.AddressIPVersion = "ipv4"
				attrResp.ResourceGroupId = "rg-test"
				return attrResp, nil
			})
		defer patchWithSlbClient.Reset()

		patchListTags := gomonkey.ApplyMethod(reflect.TypeOf(&SlbService{}), "ListTagResources",
			func(_ *SlbService, id string, resourceType string) (interface{}, error) {
				return nil, fmt.Errorf("ListTagResources failed for acl %s", id)
			})
		defer patchListTags.Reset()

		err := dataSourceAlicloudSlbAclsRead(d, rawClient)
		assert.NotNil(t, err, "data source must propagate ListTagResources error instead of swallowing it and returning success")
	})

	// Multiple ACLs where only the second ACL's tags query fails: the whole
	// data source must return an error, not silently succeed with the first
	// ACL's state (the original bug's most damaging symptom).
	t.Run("MultipleAclsPartialTagsErrorPropagated", func(t *testing.T) {
		callCount := 0
		patchWithSlbClient := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "WithSlbClient",
			func(_ *connectivity.AliyunClient, do func(*slb.Client) (interface{}, error)) (interface{}, error) {
				callCount++
				if callCount == 1 {
					resp := slb.CreateDescribeAccessControlListsResponse()
					resp.Acls = slb.Acls{Acl: []slb.Acl{
						{AclId: "acl-test-1", AclName: "tf-test-acl-1"},
						{AclId: "acl-test-2", AclName: "tf-test-acl-2"},
					}}
					return resp, nil
				}
				attrResp := slb.CreateDescribeAccessControlListAttributeResponse()
				attrResp.AclId = "acl-test-1"
				attrResp.AclName = "tf-test-acl-1"
				attrResp.AddressIPVersion = "ipv4"
				attrResp.ResourceGroupId = "rg-test"
				return attrResp, nil
			})
		defer patchWithSlbClient.Reset()

		// First ACL's tags query succeeds, second ACL's fails: the data source
		// must error out rather than returning a half-complete success.
		tagsCallCount := 0
		patchListTags := gomonkey.ApplyMethod(reflect.TypeOf(&SlbService{}), "ListTagResources",
			func(_ *SlbService, id string, resourceType string) (interface{}, error) {
				tagsCallCount++
				if tagsCallCount == 1 {
					return nil, nil
				}
				return nil, fmt.Errorf("ListTagResources failed for acl %s", id)
			})
		defer patchListTags.Reset()

		err := dataSourceAlicloudSlbAclsRead(d, rawClient)
		assert.NotNil(t, err, "multi-ACL partial tags failure must surface an error, not silent success with half-complete state")
		assert.GreaterOrEqual(t, tagsCallCount, 2, "ListTagResources must be invoked for each ACL until the failure surfaces")
	})
}
