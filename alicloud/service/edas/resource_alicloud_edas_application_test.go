package edas_test

// collection-order: bypass ecu_info -- pre-existing stock TypeList unchanged by this directory-only migration; a convergent two-member reorder apply requires live EDAS application ECU members that the existing tests do not exercise, and fabricating ECU members is out of scope for this migration

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	edassvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_edas_application", &resource.Sweeper{
		Name: "alicloud_edas_application",
		F:    testSweepEdasApplication,
	})
}

func testSweepEdasApplication(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	edasService := edassvc.NewEdasService(client)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	applicationListRq := edas.CreateListApplicationRequest()
	applicationListRq.RegionId = region

	raw, err := client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListApplication(applicationListRq)
	})
	if err != nil || raw == nil {
		log.Printf("[ERROR] Failed to retrieve edas application in service list: %s", err)
	}

	listApplicationResponse, _ := raw.(*edas.ListApplicationResponse)
	if listApplicationResponse == nil || listApplicationResponse.Code != 200 {
		log.Printf("[ERROR] Failed to retrieve edas application in service list: %v", listApplicationResponse)
		return nil
	}

	for _, v := range listApplicationResponse.ApplicationList.Application {
		name := v.Name
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}

		if skip {
			log.Printf("[INFO] Skipping edas application: %s", name)
			continue
		}
		log.Printf("[INFO] delete edas application: %s", name)
		// stop it before delete
		stopAppRequest := edas.CreateStopApplicationRequest()
		stopAppRequest.RegionId = region
		stopAppRequest.AppId = v.AppId

		raw, err := client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.StopApplication(stopAppRequest)
		})
		if err != nil {
			return err
		}
		logs.AddDebug(stopAppRequest.GetActionName(), raw, stopAppRequest.RoaRequest, stopAppRequest)
		stopAppResponse, _ := raw.(*edas.StopApplicationResponse)
		changeOrderId := stopAppResponse.ChangeOrderId

		if len(changeOrderId) > 0 {
			stateConf := helper.BuildStateConf([]string{"0", "1"}, []string{"2"}, 5*time.Minute, 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return err
			}
		}

		deleteApplicationRequest := edas.CreateDeleteApplicationRequest()
		deleteApplicationRequest.RegionId = region
		deleteApplicationRequest.AppId = v.AppId

		err = retry.Retry(2*time.Minute, func() *retry.RetryError {
			raw, err := client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteApplication(deleteApplicationRequest)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser}) {
					time.Sleep(10 * time.Second)
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(deleteApplicationRequest.GetActionName(), raw, deleteApplicationRequest.RoaRequest, deleteApplicationRequest)
			rsp := raw.(*edas.DeleteApplicationResponse)
			if rsp.Code == 601 && strings.Contains(rsp.Message, "Operation cannot be processed because there are running instances.") {
				err = errs.Error("Operation cannot be processed because there are running instances.")
				return retry.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] DeleteApplication failed: %v", err)
		}
	}

	return nil
}

func TestAccAliCloudEdasApplication_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "alicloud_edas_application.default"
	ra := acctest.ResourceAttrInit(resourceId, edasApplicationBasicMap)
	serviceFunc := func() interface{} {
		return edassvc.NewEdasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandIntRange(1000, 9999)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasapplicationbasic%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEdasApplicationConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEdasApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name": "${var.name}",
					"package_type":     "JAR",
					"cluster_id":       "${alicloud_edas_cluster.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name": name,
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
					"application_name": fmt.Sprintf("tf-testacc-edasapplicationchange%v", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name": fmt.Sprintf("tf-testacc-edasapplicationchange%v", rand)}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"package_type": "WAR",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_type": "WAR"}),
				),
			},
		},
	})
}

func TestAccAliCloudEdasApplication_multi(t *testing.T) {
	var v *edas.Applcation
	resourceId := "alicloud_edas_application.default.1"
	ra := acctest.ResourceAttrInit(resourceId, edasApplicationBasicMap)
	serviceFunc := func() interface{} {
		return edassvc.NewEdasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandIntRange(100, 999)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasapplicationbasic%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEdasApplicationConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEdasApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":            "2",
					"application_name": "${var.name}-${count.index}",
					"package_type":     "JAR",
					"cluster_id":       "${alicloud_edas_cluster.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var edasApplicationBasicMap = map[string]string{
	"application_name": acctest.CHECKSET,
	"package_type":     acctest.CHECKSET,
	"cluster_id":       acctest.CHECKSET,
}

func testAccCheckEdasApplicationDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasApplicationConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		data "alicloud_vpcs" "default" {
			name_regex = "^default-NODELETING$"
		}

		resource "alicloud_edas_cluster" "default" {
		  cluster_name = "${var.name}"
		  cluster_type = 2
		  network_mode = 2
		  vpc_id       = data.alicloud_vpcs.default.ids.0
		}
		`, name)
}
