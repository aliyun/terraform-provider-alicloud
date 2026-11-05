package edas_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	edassvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/edas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudEdasApplicationPackageAttachment_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "alicloud_edas_application_deployment.default"
	ra := acctest.ResourceAttrInit(resourceId, edasAPAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return edassvc.NewEdasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc-edasdeploymentbasic%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEdasAPAttachmentDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testEdasCheckDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				// Step 1: create only the ECS / EDAS cluster / instance_cluster_attachment
				// prerequisites (no application referencing ecu_map yet). The attachment
				// Read returns nil with an empty ecu_map (see resource_alicloud_edas_instance_
				// cluster_attachment.go Read); the ECU registration is async. This step polls
				// ListClusterMembers until the ECU for the test's own instance is registered
				// (EcsId exact match + non-empty EcuId), so step 2's refresh re-reads a
				// populated ecu_map before the application + deployment apply.
				Config: resourceEdasAPAttachmentPrerequisites(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdasECUReady(),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":   "${alicloud_edas_application.default.id}",
					"group_id": "all",
					"war_url":  "http://edas-sz.oss-cn-shenzhen.aliyuncs.com/prod/demo/SPRING_CLOUD_CONSUMER.jar",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testEdasCheckDeploymentDestroy(s *terraform.State) error {
	return nil
}

var edasAPAttachmentBasicMap = map[string]string{
	"app_id":   acctest.CHECKSET,
	"group_id": acctest.CHECKSET,
	"war_url":  acctest.CHECKSET,
}

func resourceEdasAPAttachmentPrerequisites(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		data "alicloud_instance_types" "default" {
		  instance_type_family = "ecs.s6"
		  cpu_core_count       = 1
		  memory_size          = 2
		  system_disk_category = "cloud_essd"
		}

		data "alicloud_images" "default" {
		  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64"
		  most_recent = true
		  owners      = "system"
		}

		resource "alicloud_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		  name       = "${var.name}"
		}

		resource "alicloud_vswitch" "default" {
		  vpc_id            = alicloud_vpc.default.id
		  cidr_block        = "172.16.0.0/24"
		  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
		  name              = "${var.name}"
		}

		resource "alicloud_security_group" "default" {
		  name   = "${var.name}"
		  vpc_id = alicloud_vpc.default.id
		}

		resource "alicloud_security_group_rule" "default" {
			type = "ingress"
			ip_protocol = "tcp"
			nic_type = "intranet"
			policy = "accept"
			port_range = "22/22"
			priority = 1
			security_group_id = "${alicloud_security_group.default.id}"
			cidr_ip = "172.16.0.0/24"
		}

		resource "alicloud_instance" "default" {
		  instance_type              = "${data.alicloud_instance_types.default.instance_types.0.id}"
		  system_disk_category       = "cloud_essd"
		  image_id                   = data.alicloud_images.default.images.0.id
		  instance_name              = "${var.name}"
		  vswitch_id                 = alicloud_vswitch.default.id
		  security_groups            = ["${alicloud_security_group.default.id}"]
		  internet_max_bandwidth_out = 10
		}

		resource "alicloud_edas_cluster" "default" {
		  cluster_name = "${var.name}"
		  cluster_type = 2
		  network_mode = 2
		  vpc_id       = alicloud_vpc.default.id
		}

		resource "alicloud_edas_instance_cluster_attachment" "default" {
		  cluster_id = "${alicloud_edas_cluster.default.id}"
		  instance_ids = ["${alicloud_instance.default.id}"]
		}
		`, name)
}

func resourceEdasAPAttachmentDependence(name string) string {
	return resourceEdasAPAttachmentPrerequisites(name) + `
		resource "alicloud_edas_application" "default" {
		  application_name = "${var.name}"
		  cluster_id = "${alicloud_edas_cluster.default.id}"
		  package_type = "JAR"
		  ecu_info = ["${alicloud_edas_instance_cluster_attachment.default.ecu_map[alicloud_instance.default.id]}"]
		}
		`
}

// testAccCheckEdasECUReady polls ListClusterMembers (read-only, bounded) until the
// ECU for the test's own alicloud_instance.default is registered in
// alicloud_edas_cluster.default: the member's EcsId must exactly match the
// instance ID and its EcuId must be non-empty. This gates step 2's application
// apply (which indexes ecu_map) on real ECU readiness instead of a fixed sleep.
func testAccCheckEdasECUReady() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		instanceRs, ok := s.RootModule().Resources["alicloud_instance.default"]
		if !ok {
			return fmt.Errorf("alicloud_instance.default not found in state")
		}
		ecsId := instanceRs.Primary.ID
		clusterRs, ok := s.RootModule().Resources["alicloud_edas_cluster.default"]
		if !ok {
			return fmt.Errorf("alicloud_edas_cluster.default not found in state")
		}
		clusterId := clusterRs.Primary.ID

		wait := helper.IncrementalWait(5*time.Second, 10*time.Second)
		return retry.Retry(15*time.Minute, func() *retry.RetryError {
			request := edas.CreateListClusterMembersRequest()
			request.RegionId = client.RegionId
			request.ClusterId = clusterId
			raw, err := client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.ListClusterMembers(request)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			response := raw.(*edas.ListClusterMembersResponse)
			for _, member := range response.ClusterMemberPage.ClusterMemberList.ClusterMember {
				if member.EcsId == ecsId && member.EcuId != "" {
					return nil
				}
			}
			wait()
			return retry.RetryableError(fmt.Errorf("ECU for instance %s not yet registered in cluster %s", ecsId, clusterId))
		})
	}
}
