package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_cr_ee_instance.default.id}",
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        "${alicloud_cr_ee_instance.default.id}",
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":              "6",
					"instance_id":        "${alicloud_cr_ee_instance.default.id}",
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
	// instance_name has stricter constraints than other CR names (length cap
	// observed at 30 chars), so derive a sanitized name from `name`.
	instanceName := name
	if len(instanceName) > 30 {
		instanceName = instanceName[:30]
	}
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	# Provision a dedicated EE instance instead of relying on
	# data.alicloud_cr_ee_instances {} — that data source returns an empty list
	# in regions where no instance happens to exist at test time (frequent in
	# cn-beijing when concurrent ACC runs delete each other's instances).
	# Creating one inline keeps the test self-contained and deterministic.
	resource "alicloud_cr_ee_instance" "default" {
		payment_type   = "Subscription"
		period         = 1
		renew_period   = 1
		renewal_status = "AutoRenewal"
		instance_type  = "Advanced"
		instance_name  = "%s"
	}
`, name, instanceName)
}
