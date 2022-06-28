package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_sag_qos", &resource.Sweeper{
		Name: "alicloud_sag_qos",
		F:    testSweepSagQos,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			//"alicloud_cs_kubernetes",
		},
	})
}

func testSweepSagQos(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var qoses []smartag.Qos
	req := smartag.CreateDescribeQosesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeQoses(req)
		})
		if err != nil {
			log.Println(fmt.Errorf("Error retrieving Sag Qoses: %s", err))
			return nil
		}
		resp, _ := raw.(*smartag.DescribeQosesResponse)
		if resp == nil || len(resp.Qoses.Qos) < 1 {
			break
		}
		qoses = append(qoses, resp.Qoses.Qos...)

		if len(resp.Qoses.Qos) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	for _, qos := range qoses {
		name := qos.QosName
		id := qos.QosId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}

		if skip {
			log.Printf("[INFO] Skipping Smartag Qos: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Smartag Qos: %s (%s)", name, id)
		request := smartag.CreateDeleteQosRequest()
		request.QosId = id

		_, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DeleteQos(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Smart qos (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudSagQos_basic(t *testing.T) {
	var qos smartag.Qos
	resourceId := "alicloud_sag_qos.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &qos, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testSagQosName-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagQosDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
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
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": fmt.Sprintf("%s-Update", name),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("%s-Update", name),
					}),
				),
			},
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

func TestAccAlicloudSagQos_multi(t *testing.T) {
	var qos smartag.Qos
	resourceId := "alicloud_sag_qos.default.9"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &qos, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testSagQosName-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagQosDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":  "${var.name}-${count.index}",
					"count": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("%s-9", name),
					}),
				),
			},
		},
	})
}

func resourceSagQosDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
