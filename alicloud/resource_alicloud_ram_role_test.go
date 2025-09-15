package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ram_role", &resource.Sweeper{
		Name: "alicloud_ram_role",
		F:    testSweepRamRoles,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_fc_service",
		},
	})
}

func testSweepRamRoles(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := ram.CreateListRolesRequest()
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListRoles(request)
	})
	if err != nil {
		return WrapError(err)
	}
	resp, _ := raw.(*ram.ListRolesResponse)
	if len(resp.Roles.Role) < 1 {
		return nil
	}

	sweeped := false

	for _, v := range resp.Roles.Role {
		name := v.RoleName
		id := v.RoleId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ram Role: %s (%s)", name, id)
			continue
		}
		sweeped = true

		log.Printf("[INFO] Detaching Ram Role: %s (%s) policies.", name, id)
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			request := ram.CreateListPoliciesForRoleRequest()
			request.RoleName = name
			return ramClient.ListPoliciesForRole(request)
		})
		resp, _ := raw.(*ram.ListPoliciesForRoleResponse)
		if err != nil {
			log.Printf("[ERROR] Failed to list Ram Role (%s (%s)) policies: %s", name, id, err)
		} else if len(resp.Policies.Policy) > 0 {
			for _, v := range resp.Policies.Policy {
				_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					request := ram.CreateDetachPolicyFromRoleRequest()
					request.PolicyName = v.PolicyName
					request.RoleName = name
					request.PolicyType = v.PolicyType
					return ramClient.DetachPolicyFromRole(request)
				})
				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					log.Printf("[ERROR] Failed detach Policy %s: %#v", v.PolicyName, err)
				}
			}
		}

		log.Printf("[INFO] Deleting Ram Role: %s (%s)", name, id)
		request := ram.CreateDeleteRoleRequest()
		request.RoleName = name

		_, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteRole(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Ram Role (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

// Test Ram Role. >>> Resource test cases, automatically generated.
// Case  Role测试增加tag测试 5886
func TestAccAliCloudRamRole_basic5886(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_role.default"
	ra := resourceAttrInit(resourceId, AliCloudRamRoleMap5886)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamRoleBasicDependence5886)
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
					"role_name":                   name,
					"assume_role_policy_document": `{\"Statement\": [{\"Action\": \"sts:AssumeRole\",\"Effect\": \"Allow\",\"Principal\": {\"Service\": [\"ecs.aliyuncs.com\"]}}],\"Version\": \"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":                   name,
						"assume_role_policy_document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"assume_role_policy_document": `{\"Statement\": [{\"Action\": [\"sts:AssumeRole\",\"sts:SetSourceIdentity\"],\"Effect\": \"Allow\",\"Principal\": {\"Service\": [\"apigateway.aliyuncs.com\",\"ecs.aliyuncs.com\"]}}],\"Version\": \"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"assume_role_policy_document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_session_duration": "6000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_session_duration": "6000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudRamRole_basic5886_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_role.default"
	ra := resourceAttrInit(resourceId, AliCloudRamRoleMap5886)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamRoleBasicDependence5886)
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
					"role_name":                   name,
					"assume_role_policy_document": `{\"Statement\": [{\"Action\": [\"sts:AssumeRole\",\"sts:SetSourceIdentity\"],\"Effect\": \"Allow\",\"Principal\": {\"Service\": [\"apigateway.aliyuncs.com\",\"ecs.aliyuncs.com\"]}}],\"Version\": \"1\"}`,
					"description":                 name,
					"max_session_duration":        "6000",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"force": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":                   name,
						"assume_role_policy_document": CHECKSET,
						"description":                 name,
						"max_session_duration":        "6000",
						"tags.%":                      "2",
						"tags.Created":                "TF",
						"tags.For":                    "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

// Case  Role适配废弃字段name, document
func TestAccAliCloudRamRole_basic5888(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_role.default"
	ra := resourceAttrInit(resourceId, AliCloudRamRoleMap5886)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamRoleBasicDependence5886)
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
					"name":     name,
					"document": `{\"Statement\": [{\"Action\": \"sts:AssumeRole\",\"Effect\": \"Allow\",\"Principal\": {\"Service\": [\"ecs.aliyuncs.com\"]}}],\"Version\": \"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":     name,
						"document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"document": `{\"Statement\": [{\"Action\": \"sts:AssumeRole\",\"Effect\": \"Allow\",\"Principal\": {\"Service\": [\"apigateway.aliyuncs.com\",\"ecs.aliyuncs.com\"]}}],\"Version\": \"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_session_duration": "6000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_session_duration": "6000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudRamRole_basic5888_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_role.default"
	ra := resourceAttrInit(resourceId, AliCloudRamRoleMap5886)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamRoleBasicDependence5886)
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
					"name":                 name,
					"document":             `{\"Statement\": [{\"Action\": \"sts:AssumeRole\",\"Effect\": \"Allow\",\"Principal\": {\"Service\": [\"ecs.aliyuncs.com\"]}}],\"Version\": \"1\"}`,
					"description":          name,
					"max_session_duration": "6000",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"force": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"document":             CHECKSET,
						"description":          name,
						"max_session_duration": "6000",
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

// Case  Role适配废弃字段name, version, services
func TestAccAliCloudRamRole_basic5890(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_role.default"
	ra := resourceAttrInit(resourceId, AliCloudRamRoleMap5886)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamRoleBasicDependence5886)
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
					"name":     name,
					"version":  "1",
					"services": []string{"ecs.aliyuncs.com"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":       name,
						"version":    "1",
						"services.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"services": []string{"ecs.aliyuncs.com", "apigateway.aliyuncs.com"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"services.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_session_duration": "6000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_session_duration": "6000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudRamRole_basic5890_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_role.default"
	ra := resourceAttrInit(resourceId, AliCloudRamRoleMap5886)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamRoleBasicDependence5886)
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
					"name":                 name,
					"version":              "1",
					"services":             []string{"ecs.aliyuncs.com"},
					"description":          name,
					"max_session_duration": "6000",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"force": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"version":              "1",
						"services.#":           "1",
						"description":          name,
						"max_session_duration": "6000",
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudRamRole_multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_role.default.5"
	ra := resourceAttrInit(resourceId, AliCloudRamRoleMap5886)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamRoleBasicDependence5886)
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
					"count":                       "6",
					"role_name":                   name + "-${count.index}",
					"assume_role_policy_document": `{\"Statement\": [{\"Action\": \"sts:AssumeRole\",\"Effect\": \"Allow\",\"Principal\": {\"Service\": [\"ecs.aliyuncs.com\"]}}],\"Version\": \"1\"}`,
					"description":                 name + "-${count.index}",
					"max_session_duration":        "6000",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"force": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":                   name + fmt.Sprint(-5),
						"assume_role_policy_document": CHECKSET,
						"description":                 name + fmt.Sprint(-5),
						"max_session_duration":        "6000",
						"tags.%":                      "2",
						"tags.Created":                "TF",
						"tags.For":                    "Test",
					}),
				),
			},
		},
	})
}

var AliCloudRamRoleMap5886 = map[string]string{
	"arn":         CHECKSET,
	"create_time": CHECKSET,
	"role_id":     CHECKSET,
}

func AliCloudRamRoleBasicDependence5886(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
