package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudEdasK8sApplicationPackageAttachment_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "alicloud_edas_k8s_application_package_attachment.default"
	ra := resourceAttrInit(resourceId, edasK8sAPAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc-edask8sdep%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sAPAttachmentDependence)
	region := os.Getenv("ALICLOUD_REGION")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testEdasCheckK8sDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":    "${alicloud_edas_k8s_application.default.id}",
					"replicas":  "1",
					"image_url": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/edas-demo-image/consumer:1.0", region),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testEdasCheckK8sDeploymentDestroy(s *terraform.State) error {
	return nil
}

var edasK8sAPAttachmentBasicMap = map[string]string{
	"app_id":    CHECKSET,
	"replicas":  CHECKSET,
	"image_url": CHECKSET,
}

func resourceEdasK8sAPAttachmentDependence(name string) string {
	region := os.Getenv("ALICLOUD_REGION")
	img := fmt.Sprintf("registry-vpc.%s.aliyuncs.com/edas-demo-image/consumer:1.0", region)
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}
		
		data "alicloud_zones" default {
		  available_resource_creation = "VSwitch"
		}
		
		data "alicloud_instance_types" "default" {
		  availability_zone = data.alicloud_zones.default.zones.0.id
		  cpu_core_count = 2
		  memory_size = 4
		  kubernetes_node_role = "Worker"
		}
		
		resource "alicloud_vpc" "default" {
		  name = var.name
		  cidr_block = "10.1.0.0/21"
		}
		
		resource "alicloud_vswitch" "default" {
		  name = var.name
		  vpc_id = alicloud_vpc.default.id
		  cidr_block = "10.1.1.0/24"
		  availability_zone = data.alicloud_zones.default.zones.0.id
		}
		
		resource "alicloud_cs_managed_kubernetes" "default" {
		  worker_instance_types = [data.alicloud_instance_types.default.instance_types.0.id]
		  name = var.name
		  worker_vswitch_ids = [alicloud_vswitch.default.id]
		  worker_number = "2"
		  password =                    "Test12345"
		  pod_cidr =                   "172.20.0.0/16"
		  service_cidr =               "172.21.0.0/20"
		  worker_disk_size =            "50"
		  worker_disk_category =         "cloud_ssd"
		  worker_data_disk_size =       "20"
		  worker_data_disk_category =   "cloud_ssd"
		  worker_instance_charge_type = "PostPaid"
		  slb_internet_enabled =        "true"
		}
		
		resource "alicloud_edas_k8s_cluster" "default" {
		  cs_cluster_id = alicloud_cs_managed_kubernetes.default.id
		}
		
		resource "alicloud_edas_k8s_application" "default" {
		  application_name = var.name
		  cluster_id = alicloud_edas_k8s_cluster.default.id
		  replicas = "1"
		  package_type = "Image"
		  image_url = "%v"
		}

		`, name, img)
}
