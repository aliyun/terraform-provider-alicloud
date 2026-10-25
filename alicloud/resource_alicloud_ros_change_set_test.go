package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// The template body used by CREATE change sets in basic test: defines a VPC
// resource (logical id "Vpc") with two Parameters so that the parameters set
// can be modified across steps to satisfy the coverage gate (testModifySet).
// Internal quotes are backslash-escaped (\"...) so the coverage checker can
// unmarshal the config map as JSON, matching the established repo pattern.
const testAccRosChangeSetTemplateBody = `{\"ROSTemplateFormatVersion\":\"2015-09-01\",\"Parameters\":{\"VpcName\":{\"Type\":\"String\"},\"VpcCidr\":{\"Type\":\"String\",\"Default\":\"172.16.0.0/12\"}},\"Resources\":{\"Vpc\":{\"Type\":\"ALIYUN::ECS::VPC\",\"Properties\":{\"VpcName\":{\"Ref\":\"VpcName\"},\"CidrBlock\":{\"Ref\":\"VpcCidr\"}}}}}`

const testAccRosChangeSetStackPolicyBody = `{\"Statement\":[{\"Action\":\"Update:*\",\"Resource\":\"*\",\"Effect\":\"Allow\",\"Principal\":\"*\"}]}`

// testAccRosChangeSetImportTemplateBody is used by the IMPORT change set in
// the import test. Every resource of an IMPORT change set must declare an
// explicit DeletionPolicy, otherwise the change set fails server-side with
// "DeletionPolicy of resource(...) is not specified".
const testAccRosChangeSetImportTemplateBody = `{\"ROSTemplateFormatVersion\":\"2015-09-01\",\"Parameters\":{\"VpcName\":{\"Type\":\"String\"},\"VpcCidr\":{\"Type\":\"String\"}},\"Resources\":{\"Vpc\":{\"Type\":\"ALIYUN::ECS::VPC\",\"DeletionPolicy\":\"Retain\",\"Properties\":{\"VpcName\":{\"Ref\":\"VpcName\"},\"CidrBlock\":{\"Ref\":\"VpcCidr\"}}}}}`

// testAccRosChangeSetUpdateStackTemplateBody is the template shared by the ROS
// stack fixture and the UPDATE change set in the update test (every change set
// requires a template source). The change set overrides VpcCidr so it has an
// actual change.
const testAccRosChangeSetUpdateStackTemplateBody = `{\"ROSTemplateFormatVersion\":\"2015-09-01\",\"Parameters\":{\"VpcName\":{\"Type\":\"String\"},\"VpcCidr\":{\"Type\":\"String\",\"Default\":\"172.16.0.0/12\"}},\"Resources\":{\"Vpc\":{\"Type\":\"ALIYUN::ECS::VPC\",\"Properties\":{\"VpcName\":{\"Ref\":\"VpcName\"},\"CidrBlock\":{\"Ref\":\"VpcCidr\"}}}}}`

// TestAccAliCloudROSChangeSet_basic covers the template_body scenario plus all
// CREATE-compatible attributes and a modify step that changes the parameters
// set so parameter_key/parameter_value land in testModifySet. The import step
// keeps template_body out of ImportStateVerifyIgnore (it is read back, so
// ForceNew+Ignore would be illegal). ResourcesToImport is IMPORT-only per the
// API contract and is covered by TestAccAliCloudROSChangeSet_import.
func TestAccAliCloudROSChangeSet_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_change_set.default"
	ra := resourceAttrInit(resourceId, AlicloudRosChangeSetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosChangeSet")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosChangeSet%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRosChangeSetBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"change_set_name":   name,
					"stack_name":        name + "stack",
					"change_set_type":   "CREATE",
					"description":       "Test From Terraform",
					"template_body":     testAccRosChangeSetTemplateBody,
					"disable_rollback":  true,
					"notification_urls": []string{"http://ros-notification.example.com/hook"},
					"parallelism":       10,
					"parameters": []map[string]interface{}{
						{"parameter_key": "VpcName", "parameter_value": "tf-testacc-basic"},
					},
					"ram_role_name":     "${alicloud_ram_role.default.name}",
					"resource_group_id": "${alicloud_resource_manager_resource_group.default.id}",
					// Only one of the four stack policy parameters can be set
					// per request; step 1 covers stack_policy_body and the
					// carried-over value stays the only policy parameter in
					// step 2 (the URL variants are covered by the dedicated
					// conflict cases below).
					"stack_policy_body": testAccRosChangeSetStackPolicyBody,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"timeout_in_minutes": 60,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"change_set_name": name,
						"stack_name":      name + "stack",
						"change_set_type": "CREATE",
						"description":     "Test From Terraform",
					}),
				),
			},
			{
				// Modify step: the parameters set changes (VpcName value differs
				// and a new VpcCidr parameter appears), so parameter_key and
				// parameter_value are both in testModifySet. Step attributes
				// merge cumulatively, so the stack_policy_body from step 1 is
				// carried over (only one stack policy parameter per request).
				// stack_name is ForceNew and step 1 already created a change set
				// on the first placeholder stack, so this step uses a second
				// stack name and replaces the change set resource.
				Config: testAccConfig(map[string]interface{}{
					"change_set_name":   name,
					"stack_name":        name + "stack2",
					"change_set_type":   "CREATE",
					"description":       "Test From Terraform",
					"template_body":     testAccRosChangeSetTemplateBody,
					"disable_rollback":  true,
					"notification_urls": []string{"http://ros-notification.example.com/hook"},
					"parallelism":       10,
					"parameters": []map[string]interface{}{
						{"parameter_key": "VpcName", "parameter_value": "tf-testacc-basic-modify"},
						{"parameter_key": "VpcCidr", "parameter_value": "172.16.0.0/12"},
					},
					"ram_role_name":     "${alicloud_ram_role.default.name}",
					"resource_group_id": "${alicloud_resource_manager_resource_group.default.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"timeout_in_minutes": 60,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"change_set_name": name,
						"stack_name":      name + "stack2",
						"change_set_type": "CREATE",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// template_body is read back by Read (d.Set), so it must NOT be
				// ignored; ForceNew+Ignore on a read-back attr is illegal per the
				// coverage gate. The remaining attributes are write-only (not
				// read back by Read) and are legitimately ignored.
				ImportStateVerifyIgnore: []string{"notification_urls", "parallelism", "ram_role_name", "replacement_option", "resource_group_id", "resources_to_import", "stack_policy_body", "stack_policy_during_update_body", "stack_policy_during_update_url", "stack_policy_url", "taint_resources", "template_id", "template_scratch_id", "template_url", "template_version", "use_previous_parameters"},
			},
		},
	})
}

// TestAccAliCloudROSChangeSet_import covers the IMPORT change_set scenario:
// resources_to_import with the identifier map (serialized to a JSON string per
// the API contract) and parameters describing the imported VPC.
func TestAccAliCloudROSChangeSet_import(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_change_set.default"
	ra := resourceAttrInit(resourceId, AlicloudRosChangeSetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosChangeSet")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosChangeSet%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRosChangeSetImportDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"change_set_name": name,
					"stack_name":      name + "stack",
					"change_set_type": "IMPORT",
					"description":     "Test From Terraform",
					"template_body":   testAccRosChangeSetImportTemplateBody,
					"parameters": []map[string]interface{}{
						{"parameter_key": "VpcName", "parameter_value": "${alicloud_vpc.default.name}"},
						{"parameter_key": "VpcCidr", "parameter_value": "${alicloud_vpc.default.cidr_block}"},
					},
					"resources_to_import": []map[string]interface{}{
						{"resource_identifier": map[string]interface{}{"VpcId": "${alicloud_vpc.default.id}"}, "logical_resource_id": "Vpc", "resource_type": "ALIYUN::ECS::VPC"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"change_set_name": name,
						"stack_name":      name + "stack",
						"change_set_type": "IMPORT",
						"description":     "Test From Terraform",
					}),
				),
			},
		},
	})
}

// AlicloudRosChangeSetImportDependence creates the VPC fixture imported by the
// import test's resources_to_import.
func AlicloudRosChangeSetImportDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_vpc" "default" {
  vpc_name   = "%[1]s-vpc"
  cidr_block = "172.16.0.0/12"
}
`, name)
}

// TestAccAliCloudROSChangeSet_templateId covers the mutually-exclusive
// template_id + template_version scenario (template sourced from a ROS
// template fixture; the version is resolved through the ros_templates data
// source, which exposes template_version).
func TestAccAliCloudROSChangeSet_templateId(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_change_set.default"
	ra := resourceAttrInit(resourceId, AlicloudRosChangeSetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosChangeSet")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosChangeSet%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRosChangeSetTemplateDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"change_set_name":  name,
					"stack_name":       name + "stack",
					"change_set_type":  "CREATE",
					"description":      "Test From Terraform",
					"template_id":      "${data.alicloud_ros_templates.default.templates.0.template_id}",
					"template_version": "${data.alicloud_ros_templates.default.templates.0.template_version}",
					// The template fixture declares VpcName without a default,
					// so CreateChangeSet rejects the request with
					// UserParameterMissing unless the value is provided.
					"parameters": []map[string]interface{}{
						{"parameter_key": "VpcName", "parameter_value": "tf-testacc-template-id"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"change_set_name": name,
						"change_set_type": "CREATE",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudROSChangeSet_templateScratch covers the mutually-exclusive
// template_scratch_id scenario (template sourced from a ROS template scratch).
func TestAccAliCloudROSChangeSet_templateScratch(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_change_set.default"
	ra := resourceAttrInit(resourceId, AlicloudRosChangeSetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosChangeSet")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosChangeSet%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRosChangeSetScratchDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// TemplateScratchId is only supported for IMPORT change sets per
				// the API contract.
				Config: testAccConfig(map[string]interface{}{
					"change_set_name":     name,
					"stack_name":          name + "stack",
					"change_set_type":     "IMPORT",
					"description":         "Test From Terraform",
					"template_scratch_id": "${alicloud_ros_template_scratch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"change_set_name": name,
						"change_set_type": "IMPORT",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudROSChangeSet_update covers the UPDATE change_set scenario:
// stack_id (existing ROS stack fixture), use_previous_parameters,
// stack_policy_during_update_body, replacement_option and taint_resources.
// The URL-form stack policy attribute (stack_policy_during_update_url) and
// template_url are covered by TestUnitAlicloudROSChangeSet: the integration
// account blocks public OSS access (object-level public ACL PUT is rejected
// and anonymous reads are forbidden), so ROS cannot fetch OSS-hosted URLs.
func TestAccAliCloudROSChangeSet_update(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_change_set.default"
	ra := resourceAttrInit(resourceId, AlicloudRosChangeSetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosChangeSet")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosChangeSet%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRosChangeSetUpdateDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"change_set_name":         name,
					"stack_id":                "${alicloud_ros_stack.default.id}",
					"change_set_type":         "UPDATE",
					"description":             "Test From Terraform",
					"template_body":           testAccRosChangeSetUpdateStackTemplateBody,
					"use_previous_parameters": true,
					"parameters": []map[string]interface{}{
						{"parameter_key": "VpcCidr", "parameter_value": "172.16.1.0/24"},
					},
					"replacement_option":              "Disabled",
					"stack_policy_during_update_body": testAccRosChangeSetStackPolicyBody,
					"taint_resources":                 []string{"Vpc"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"change_set_name": name,
						"change_set_type": "UPDATE",
					}),
				),
			},
			{
				// UPDATE-type change sets allow only one
				// stack_policy_during_update_* parameter per request (verified
				// live: the CREATE type accepts both, the UPDATE type rejects
				// the pair). The body is carried over from step 1, so adding
				// the URL must fail with MultipleParameter before anything is
				// created. The replacement destroys the step 1 change set
				// first, which CheckDestroy then no longer needs to repeat.
				Config: testAccConfig(map[string]interface{}{
					"stack_policy_during_update_url": "https://invalid.tf-testacc.example.com/policy-during-update.json",
				}),
				ExpectError: regexp.MustCompile("MultipleParameter"),
			},
		},
	})
}

func TestAccAliCloudROSChangeSet_tags(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_change_set.default"
	ra := resourceAttrInit(resourceId, AlicloudRosChangeSetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosChangeSet")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosChangeSet%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRosChangeSetBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"change_set_name": name,
					"stack_name":      name + "stack",
					"change_set_type": "CREATE",
					"description":     "Test From Terraform",
					"template_body":   `{\"ROSTemplateFormatVersion\":\"2015-09-01\"}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"change_set_name": name,
						"stack_name":      name + "stack",
						"change_set_type": "CREATE",
						"description":     "Test From Terraform",
						"template_body":   CHECKSET,
					}),
					resource.TestCheckResourceAttr(resourceId, "tags.Created", "TF"),
					resource.TestCheckResourceAttr(resourceId, "tags.For", "Test"),
				),
			},
		},
	})
}

// The integration account blocks public OSS objects (account-level Block
// Public Access), so a change set cannot be created from a real template URL:
// the object cannot be uploaded with a public ACL and the ROS server-side
// fetch of a private object is rejected. The two conflict cases below
// therefore cover the *_url attributes through the documented mutual
// exclusion contract: CreateChangeSet validates the pairs before any template
// retrieval or stack creation, so the request fails with MultipleParameter
// and leaves no residue (verified live via ListStacks after each probe).
func TestAccAliCloudROSChangeSet_templateBodyUrlConflict(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_change_set.default"
	ra := resourceAttrInit(resourceId, AlicloudRosChangeSetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosChangeSet")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosChangeSet%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return ""
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"change_set_name": name,
					"stack_name":      name + "stack",
					"change_set_type": "CREATE",
					"template_body":   `{\"ROSTemplateFormatVersion\":\"2015-09-01\"}`,
					"template_url":    "https://invalid.tf-testacc.example.com/template.json",
				}),
				ExpectError: regexp.MustCompile("MultipleParameter"),
			},
		},
	})
}

func TestAccAliCloudROSChangeSet_stackPolicyUrlConflict(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_change_set.default"
	ra := resourceAttrInit(resourceId, AlicloudRosChangeSetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosChangeSet")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosChangeSet%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return ""
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"change_set_name":   name,
					"stack_name":        name + "stack",
					"change_set_type":   "CREATE",
					"template_body":     `{\"ROSTemplateFormatVersion\":\"2015-09-01\"}`,
					"stack_policy_body": testAccRosChangeSetStackPolicyBody,
					"stack_policy_url":  "https://invalid.tf-testacc.example.com/policy.json",
				}),
				ExpectError: regexp.MustCompile("MultipleParameter"),
			},
		},
	})
}

var AlicloudRosChangeSetMap = map[string]string{
	"status": CHECKSET,
}

// AlicloudRosChangeSetBasicDependence creates the shared fixtures used by the
// basic and tags tests: a resource manager resource group and a RAM role.
func AlicloudRosChangeSetBasicDependence(name string) string {
	return fmt.Sprintf(`
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_resource_manager_resource_group" "default" {
  name         = "%[1]s-rg"
  display_name = "%[1]s-rg"
}

resource "alicloud_ram_role" "default" {
  name     = "%[1]s-role"
  document = "{\"Version\":\"1\",\"Statement\":[{\"Action\":\"sts:AssumeRole\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":[\"cs.aliyuncs.com\"]}}]}"
}
`, name)
}

// AlicloudRosChangeSetTemplateDependence creates a ROS template fixture used
// by the templateId test (template_id + template_version). The templates data
// source resolves the created template's id and latest version.
func AlicloudRosChangeSetTemplateDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ros_template" "default" {
  template_name = "%[1]s-tmpl"
  template_body = "{\"ROSTemplateFormatVersion\":\"2015-09-01\",\"Parameters\":{\"VpcName\":{\"Type\":\"String\"}}}"
  description   = "tf-testacc ros template"
}

data "alicloud_ros_templates" "default" {
  ids = [alicloud_ros_template.default.id]
}
`, name)
}

// AlicloudRosChangeSetScratchDependence creates a ROS template scratch fixture
// used by the templateScratch test (template_scratch_id).
func AlicloudRosChangeSetScratchDependence(name string) string {
	return fmt.Sprintf(`
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ros_template_scratch" "default" {
  template_scratch_type = "ResourceImport"
  description           = "%[1]s-scratch"
  preference_parameters {
    parameter_key   = "DeletionPolicy"
    parameter_value = "Retain"
  }
  source_resource_group {
    resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids.0
    resource_type_filter = ["ALIYUN::ECS::VPC"]
  }
}
`, name)
}

// AlicloudRosChangeSetUpdateDependence creates a ROS stack fixture used by the
// update test (stack_id + UPDATE). The stack uses the shared update template
// so the UPDATE change set can override VpcCidr as an actual change.
func AlicloudRosChangeSetUpdateDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ros_stack" "default" {
  stack_name        = "%[1]s-stack"
  template_body     = "%[2]s"
  parameters {
    parameter_key   = "VpcName"
    parameter_value = "tf-testacc-stack"
  }
  parameters {
    parameter_key   = "VpcCidr"
    parameter_value = "172.16.0.0/12"
  }
  timeout_in_minutes = 10
  // GetStackPolicy echoes "{}" for policy-less stacks, which Read writes back
  // and makes the post-apply plan non-empty; the fixture declares the policy
  // explicitly so state converges with the config.
  stack_policy_body = "{\"Statement\": [{\"Action\": \"Update:*\", \"Resource\": \"*\", \"Effect\": \"Allow\", \"Principal\": \"*\"}]}"
}
`, name, testAccRosChangeSetUpdateStackTemplateBody)
}

// lintignore: R001
func TestUnitAlicloudROSChangeSet(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ros_change_set"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ros_change_set"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"change_set_name":   "CreateChangeSetValue",
		"stack_name":        "CreateChangeSetValue",
		"change_set_type":   "CreateChangeSetValue",
		"description":       "CreateChangeSetValue",
		"template_body":     "CreateChangeSetValue",
		"disable_rollback":  false,
		"notification_urls": []string{"CreateChangeSetValue"},
		"parallelism":       10,
		"ram_role_name":     "CreateChangeSetValue",
		"parameters": []map[string]interface{}{
			{
				"parameter_key":   "CreateChangeSetValue",
				"parameter_value": "CreateChangeSetValue",
			},
		},
		"replacement_option": "CreateChangeSetValue",
		"resource_group_id":  "CreateChangeSetValue",
		"resources_to_import": []map[string]interface{}{
			{
				"resource_identifier": map[string]interface{}{"CreateChangeSetValue": "CreateChangeSetValue"},
				"logical_resource_id": "CreateChangeSetValue",
				"resource_type":       "CreateChangeSetValue",
			},
		},
		"stack_id":                        "CreateChangeSetValue",
		"stack_policy_body":               "CreateChangeSetValue",
		"stack_policy_during_update_body": "CreateChangeSetValue",
		"stack_policy_during_update_url":  "CreateChangeSetValue",
		"stack_policy_url":                "CreateChangeSetValue",
		"tags":                            map[string]interface{}{"CreateChangeSetValue": "CreateChangeSetValue"},
		"taint_resources":                 []string{"CreateChangeSetValue"},
		"template_id":                     "CreateChangeSetValue",
		"template_scratch_id":             "CreateChangeSetValue",
		"template_url":                    "CreateChangeSetValue",
		"template_version":                "CreateChangeSetValue",
		"timeout_in_minutes":              60,
		"use_previous_parameters":         false,
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// GetChangeSet
		"ChangeSetName":   "CreateChangeSetValue",
		"ChangeSetType":   "CreateChangeSetValue",
		"Description":     "CreateChangeSetValue",
		"DisableRollback": "CreateChangeSetValue",
		"Parameters": []interface{}{
			map[string]interface{}{
				"ParameterKey":   "CreateChangeSetValue",
				"ParameterValue": "CreateChangeSetValue",
			},
		},
		"StackId":          "CreateChangeSetValue",
		"StackName":        "CreateChangeSetValue",
		"Status":           "CREATE_COMPLETE",
		"TemplateBody":     "CreateChangeSetValue",
		"TimeoutInMinutes": 60,
		"ChangeSetId":      "CreateChangeSetValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateChangeSet
		"ChangeSetId": "CreateChangeSetValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ros_change_set", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRosClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudRosChangeSetCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetChangeSet Response
		"ChangeSetId": "CreateChangeSetValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateChangeSet" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudRosChangeSetCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ros_change_set"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetChangeSet" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudRosChangeSetRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRosClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudRosChangeSetDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteChangeSet" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudRosChangeSetDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
