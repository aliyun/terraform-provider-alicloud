package alikafka_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alikafka"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudAliKafkaInstanceAllowedIpAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_instance_allowed_ip_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAliKafkaInstanceAllowedIpAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAliKafkaInstanceAllowedIpAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salikafkainstanceallowedipattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliKafkaInstanceAllowedIpAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_alikafka_instance.default.id}",
					"allowed_type": "vpc",
					"port_range":   "9092/9092",
					"allowed_ip":   "172.168.4.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
						"allowed_type": "vpc",
						"port_range":   "9092/9092",
						"allowed_ip":   "172.168.4.0/24",
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

func TestAccAliCloudAliKafkaInstanceAllowedIpAttachment_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_instance_allowed_ip_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAliKafkaInstanceAllowedIpAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAliKafkaInstanceAllowedIpAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salikafkainstanceallowedipattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliKafkaInstanceAllowedIpAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_alikafka_instance.default.id}",
					"allowed_type": "internet",
					"port_range":   "9093/9093",
					"allowed_ip":   "1.1.1.1/1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
						"allowed_type": "internet",
						"port_range":   "9093/9093",
						"allowed_ip":   "1.1.1.1/1",
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

func TestAccAliCloudAliKafkaInstanceAllowedIpAttachment_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_instance_allowed_ip_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAliKafkaInstanceAllowedIpAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAliKafkaInstanceAllowedIpAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salikafkainstanceallowedipattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliKafkaInstanceAllowedIpAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_alikafka_instance.default.id}",
					"allowed_type": "vpc",
					"port_range":   "9094/9094",
					"allowed_ip":   "172.168.4.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
						"allowed_type": "vpc",
						"port_range":   "9094/9094",
						"allowed_ip":   "172.168.4.0/24",
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

func TestAccAliCloudAliKafkaInstanceAllowedIpAttachment_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_instance_allowed_ip_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAliKafkaInstanceAllowedIpAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAliKafkaInstanceAllowedIpAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salikafkainstanceallowedipattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliKafkaInstanceAllowedIpAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_alikafka_instance.default.id}",
					"allowed_type": "vpc",
					"port_range":   "9095/9095",
					"allowed_ip":   "172.168.4.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
						"allowed_type": "vpc",
						"port_range":   "9095/9095",
						"allowed_ip":   "172.168.4.0/24",
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

var AliCloudAliKafkaInstanceAllowedIpAttachmentMap0 = map[string]string{}

func AliCloudAliKafkaInstanceAllowedIpAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_alikafka_instance" "default" {
	  		name           = var.name
	  		partition_num  = "50"
	  		disk_type      = "1"
	  		disk_size      = "500"
	  		deploy_type    = "4"
	  		eip_max        = "3"
	  		io_max         = "20"
	  		vswitch_id     = data.alicloud_vswitches.default.ids.0
	  		security_group = alicloud_security_group.default.id
		}
`, name)
}
