package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_cr_ee_namespace", &resource.Sweeper{
		Name: "alicloud_cr_ee_namespace",
		F:    testSweepCrEENamespace,
	})
}

var testaccCrEEInstanceId string

func setTestaccCrEEInstanceId(t *testing.T) {
	if testaccCrEEInstanceId != "" {
		return
	}

	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping cr ee test case with err: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	crService := &CrService{client}
	resp, err := crService.ListCrEEInstances(1, 10)
	if err != nil {
		t.Skipf("Skipping cr ee test case with err: %s", err)
	}
	if len(resp.Instances) == 0 {
		t.Skipf("Skipping cr ee test case without default instances")
	}
	testaccCrEEInstanceId = resp.Instances[0].InstanceId
}

func testSweepCrEENamespace(region string) error {
	if testaccCrEEInstanceId == "" {
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(fmt.Errorf("error getting Alicloud client: %s", err))
	}
	client := rawClient.(*connectivity.AliyunClient)
	crService := &CrService{client}

	pageNo := 1
	pageSize := 50
	var namespaces []cr_ee.NamespacesItem
	for {
		resp, err := crService.ListCrEENamespaces(testaccCrEEInstanceId, pageNo, pageSize)
		if err != nil {
			return WrapError(err)
		}
		namespaces = append(namespaces, resp.Namespaces...)
		if len(resp.Namespaces) < pageSize {
			break
		}
		pageNo++
	}

	testPrefix := "tf-testacc"
	for _, namespace := range namespaces {
		if strings.HasPrefix(namespace.NamespaceName, testPrefix) {
			crService.DeleteCrEENamespace(namespace.InstanceId, namespace.NamespaceName)
		}
	}
	return nil
}

func TestAccAlicloudCrEENamespace_Basic(t *testing.T) {
	setTestaccCrEEInstanceId(t)
	var v *cr_ee.GetNamespaceResponse
	resourceId := "alicloud_cr_ee_namespace.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEENamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ns-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEENamespaceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithCrEE(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        testaccCrEEInstanceId,
					"name":               name,
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        testaccCrEEInstanceId,
						"name":               name,
						"auto_create":        "false",
						"default_visibility": "PUBLIC",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
					"default_visibility": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_visibility": "PRIVATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               name,
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"auto_create":        "false",
						"default_visibility": "PUBLIC",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCrEENamespace_Multi(t *testing.T) {
	setTestaccCrEEInstanceId(t)
	var v *cr_ee.GetNamespaceResponse
	resourceId := "alicloud_cr_ee_namespace.default.4"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEENamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ns-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEENamespaceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithCrEE(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        testaccCrEEInstanceId,
					"name":               name + "${count.index}",
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
					"count":              "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name + fmt.Sprint(4),
						"auto_create":        "false",
						"default_visibility": "PUBLIC",
					}),
				),
			},
		},
	})
}

func resourceCrEENamespaceConfigDependence(name string) string {
	return ""
}

func testAccPreCheckWithCrEE(t *testing.T) {
	testAccPreCheck(t)
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping cr ee test case with err: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	crService := &CrService{client}
	resp, err := crService.ListCrEEInstances(1, 10)
	if err != nil {
		//Maybe crEE has not opened int the region
		t.Skipf("Skipping cr ee test case with err: %s", err)
	}
	if len(resp.Instances) == 0 {
		t.Skipf("Skipping cr ee test case without default instances")
	}
}
