package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudCSKMSEncryption_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kms_encryption.default"
	ra := resourceAttrInit(resourceId, csKMSEncryptionBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCsKMSEncryption")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cs-kms-encryption-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSKMSEncryptionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":         "${alicloud_cs_managed_kubernetes.CreateCluster.id}",
					"disable_encryption": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":         CHECKSET,
						"disable_encryption": "true",
					}),
				),
			},
			{
				PreConfig: func() { time.Sleep(5 * time.Minute) },
				Config: testAccConfig(map[string]interface{}{
					"disable_encryption": "false",
					"kms_key_id":         "${data.alicloud_kms_keys.kms_keys_ds.keys.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disable_encryption": "false",
						"kms_key_id":         CHECKSET,
					}),
				),
			},
			{
				PreConfig: func() { time.Sleep(5 * time.Minute) },
				Config: testAccConfig(map[string]interface{}{
					"kms_key_id": "${data.alicloud_kms_keys.kms_keys_ds.keys.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_key_id": CHECKSET,
					}),
					func(s *terraform.State) error {
						time.Sleep(5 * time.Minute)
						return nil
					},
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

func TestAccAliCloudCSKMSEncryption_Create(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cs_kms_encryption.create"
	ra := resourceAttrInit(resourceId, csKMSEncryptionCreateMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCsKMSEncryption")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cs-kms-encryption-create-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSKMSEncryptionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":         "${alicloud_cs_managed_kubernetes.CreateCluster.id}",
					"disable_encryption": "false",
					"kms_key_id":         "${data.alicloud_kms_keys.kms_keys_ds.keys.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":         CHECKSET,
						"disable_encryption": "false",
						"kms_key_id":         CHECKSET,
					}),
					func(s *terraform.State) error {
						time.Sleep(5 * time.Minute)
						return nil
					},
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

func resourceCSKMSEncryptionConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "vpc_cidr" {
  default = "10.0.0.0/8"
}

variable "vswitch_cidrs" {
  type    = list(string)
  default = ["10.1.0.0/16", "10.2.0.0/16"]
}

variable "pod_cidr" {
  default = "172.16.0.0/16"
}

variable "service_cidr" {
  default = "192.168.0.0/16"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {}

data "alicloud_kms_keys" "kms_keys_ds" {
  filters = "[{\"Key\":\"KeyState\",\"Values\":[\"Enabled\"]},{\"Key\":\"KeySpec\",\"Values\":[\"Aliyun_AES_256\"]},{\"Key\":\"KeyUsage\",\"Values\":[\"ENCRYPT/DECRYPT\"]},{\"Key\":\"CreatorType\",\"Values\":[\"User\"]}]"
}

resource "alicloud_vpc" "CreateVPC" {
  cidr_block = var.vpc_cidr
}

resource "alicloud_vswitch" "CreateVSwitch" {
  count      = length(var.vswitch_cidrs)
  vpc_id     = alicloud_vpc.CreateVPC.id
  cidr_block = element(var.vswitch_cidrs, count.index)
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index].zone_id
}

resource "alicloud_cs_managed_kubernetes" "CreateCluster" {
  name                         = var.name
  cluster_spec                 = "ack.standard"
  profile                      = "Default"
  vswitch_ids                  = split(",", join(",", alicloud_vswitch.CreateVSwitch.*.id))
  pod_cidr                     = var.pod_cidr
  service_cidr                 = var.service_cidr
  is_enterprise_security_group = true
  ip_stack                     = "ipv4"
  proxy_mode                   = "ipvs"
  deletion_protection          = false

  addons {
    name = "gatekeeper"
  }
  addons {
    name = "loongcollector"
  }
  addons {
    name = "policy-template-controller"
  }

  operation_policy {
    cluster_auto_upgrade {
      enabled = false
    }
  }
  maintenance_window {
    enable = false
  }
}
`, name)
}

var csKMSEncryptionBasicMap = map[string]string{
	"cluster_id":         CHECKSET,
	"disable_encryption": CHECKSET,
}

var csKMSEncryptionCreateMap = map[string]string{
	"cluster_id":         CHECKSET,
	"disable_encryption": CHECKSET,
	"kms_key_id":         CHECKSET,
}
