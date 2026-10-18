package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cr_ee_namespace", &resource.Sweeper{
		Name:         "alicloud_cr_ee_namespace",
		Dependencies: []string{"alicloud_cr_ee_repo"},
		F:            testSweepCrEENamespace,
	})
}

// testSweepCrEENamespace cleans up Container Registry Enterprise Edition
// namespaces left on CR EE instances by acceptance tests. Repositories are
// swept first (alicloud_cr_ee_namespace DependsOn alicloud_cr_ee_repo) so a
// namespace with test repos under it can be removed. Namespaces marked
// NODELETING (or whose instance is marked NODELETING) are preserved; in
// non-SWEEP_ALL mode only names with the standard test prefixes are removed.
func testSweepCrEENamespace(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(fmt.Errorf("error getting AliCloud client: %s", err))
	}
	client := rawClient.(*connectivity.AliyunClient)
	crService := &CrService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tftestAcc",
	}

	// Enumerate CR EE instances, then paginate namespaces per instance.
	instPageNo, pageSize := 1, PageSizeLarge
	for {
		instances, err := crService.ListCrEEInstances(instPageNo, pageSize)
		if err != nil {
			log.Printf("[ERROR] Failed to ListCrEEInstances: %s", err)
			return nil
		}
		for _, instance := range instances.Instances {
			nsPageNo := 1
			for {
				resp, err := crService.ListCrEENamespaces(instance.InstanceId, nsPageNo, pageSize)
				if err != nil {
					log.Printf("[ERROR] Failed to ListCrEENamespaces for instance %s: %s", instance.InstanceId, err)
					break
				}
				for _, namespace := range resp.Namespaces {
					name := namespace.NamespaceName
					if protectedSweepVpcName(name) || protectedSweepVpcName(instance.InstanceName) {
						log.Printf("[INFO] Skipping CR EE namespace %s:%s: protected (NODELETING)", instance.InstanceId, name)
						continue
					}
					if !sweepAll() {
						skip := true
						for _, p := range prefixes {
							if strings.HasPrefix(strings.ToLower(name), strings.ToLower(p)) {
								skip = false
								break
							}
						}
						if skip {
							continue
						}
					}
					id := fmt.Sprintf("%s:%s", instance.InstanceId, name)
					log.Printf("[INFO] Deleting CR EE namespace: %s", id)
					if _, derr := crService.DeleteCrEENamespace(id); derr != nil {
						if NotFoundError(derr) {
							continue
						}
						log.Printf("[ERROR] Failed to delete CR EE namespace %s: %s", id, derr)
						continue
					}
				}
				if len(resp.Namespaces) < pageSize {
					break
				}
				nsPageNo++
			}
		}
		if len(instances.Instances) < pageSize {
			break
		}
		instPageNo++
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
