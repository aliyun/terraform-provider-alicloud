package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cr_ee_namespace", &resource.Sweeper{
		Name: "alicloud_cr_ee_namespace",
		F:    testSweepCrEENamespace,
	})
}

func testSweepCrEENamespace(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(fmt.Errorf("error getting AliCloud client: %s", err))
	}
	client := rawClient.(*connectivity.AliyunClient)
	crService := &CrService{client}

	pageNo := 1
	pageSize := 50
	var namespaces []cr_ee.NamespacesItem
	for {
		instances, err := crService.ListCrEEInstances(pageNo, pageSize)
		if err != nil {
			return WrapError(err)
		}
		for _, instance := range instances.Instances {
			pageNo := 1
			resp, err := crService.ListCrEENamespaces(instance.InstanceId, pageNo, pageSize)
			if err != nil {
				return WrapError(err)
			}
			namespaces = append(namespaces, resp.Namespaces...)
			if len(resp.Namespaces) < pageSize {
				break
			}
			pageNo++
		}
		if len(instances.Instances) < pageSize {
			break
		}
		pageNo++
	}

	testPrefix := "tf-testacc"
	for _, namespace := range namespaces {
		if strings.HasPrefix(namespace.NamespaceName, testPrefix) {
			crService.DeleteCrEENamespace(fmt.Sprint(namespace.InstanceId, ":", namespace.NamespaceName))
		}
	}
	return nil
}

func TestAccAliCloudCREENamespace_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_namespace.default"
	ra := resourceAttrInit(resourceId, AliCloudCREENamespaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEENamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-namespace-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCREENamespaceBasicDependence0)
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
					"instance_id": "${data.alicloud_cr_ee_instances.default.ids.0}",
					"name":        name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_create": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_create": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_visibility": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_visibility": "PUBLIC",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudCREENamespace_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_namespace.default"
	ra := resourceAttrInit(resourceId, AliCloudCREENamespaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEENamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-namespace-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCREENamespaceBasicDependence0)
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
					"instance_id":        "${data.alicloud_cr_ee_instances.default.ids.0}",
					"name":               name,
					"auto_create":        "true",
					"default_visibility": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"name":               name,
						"auto_create":        "true",
						"default_visibility": "PUBLIC",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudCREENamespace_Multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_namespace.default.5"
	ra := resourceAttrInit(resourceId, AliCloudCREENamespaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEENamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-namespace-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCREENamespaceBasicDependence0)
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
					"count":              "6",
					"instance_id":        "${data.alicloud_cr_ee_instances.default.ids.0}",
					"name":               name + "-${count.index}",
					"auto_create":        "false",
					"default_visibility": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"name":               name + fmt.Sprint(-5),
						"auto_create":        "false",
						"default_visibility": "PRIVATE",
					}),
				),
			},
		},
	})
}

var AliCloudCREENamespaceMap0 = map[string]string{
	"default_visibility": CHECKSET,
}

func AliCloudCREENamespaceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_cr_ee_instances" "default" {
	}
`, name)
}
