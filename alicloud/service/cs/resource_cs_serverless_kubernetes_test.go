package cs_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cs"
	aliyungocs "github.com/denverdino/aliyungo/cs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func getTimezone(region string) string {
	var timeZoneMap = map[string]string{
		"eu-central-1": "Europe/London",
		"cn-hangzhou":  "Asia/Shanghai",
		"cn-shanghai":  "Asia/Shanghai",
		"cn-beijing":   "Asia/Shanghai",
	}
	timeZone := "Asia/Shanghai"
	if v, ok := timeZoneMap[region]; ok {
		timeZone = v
	}
	return timeZone
}

func TestAccAliCloudCSServerlessKubernetes_basic(t *testing.T) {
	var v *aliyungocs.ServerlessClusterResponse
	resourceId := "alicloud_cs_serverless_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, csServerlessKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	clusterCaCertFile, clientCertFile, clientKeyFile, err := CreateTempFiles()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(clientCertFile.Name())
	defer os.Remove(clientKeyFile.Name())
	defer os.Remove(clusterCaCertFile.Name())

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccserverlesskubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSServerlessKubernetesConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                           name,
					"version":                        "${data.alicloud_cs_kubernetes_version.version-126.metadata.0.version}",
					"vpc_id":                         "${alicloud_vpc.default.id}",
					"vswitch_ids":                    []string{"${alicloud_vswitch.default.id}"},
					"security_group_id":              "${alicloud_security_group.default.id}",
					"new_nat_gateway":                "true",
					"deletion_protection":            "false",
					"enable_rrsa":                    "false",
					"endpoint_public_access_enabled": "true",
					"load_balancer_spec":             "slb.s2.small",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"tags": map[string]string{
						"Platform": "TF",
					},
					"service_cidr":     "10.0.1.0/24",
					"private_zone":     "true",
					"logging_type":     "SLS",
					"sls_project_name": name,
					"time_zone":        getTimezone(os.Getenv("ALICLOUD_REGION")),
					"cluster_spec":     "ack.standard",
					"custom_san":       "www.terraform.io,1.1.1.1",
					"addons": []map[string]string{
						{
							"name":     "managed-arms-prometheus",
							"config":   "",
							"version":  "",
							"disabled": "false",
						},
					},
					"maintenance_window": []map[string]string{{"enable": "true", "maintenance_time": "2024-10-15T12:31:00.000+08:00", "duration": "3h", "weekly_period": "Thursday"}},
					"operation_policy": []map[string]interface{}{
						{
							"cluster_auto_upgrade": []map[string]interface{}{
								{
									"enabled": "true",
									"channel": "patch",
								},
							},
						},
					},
					"cluster_ca_cert": clusterCaCertFile.Name(),
					"client_key":      clientKeyFile.Name(),
					"client_cert":     clientCertFile.Name(),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                      name,
						"version":                                   acctest.CHECKSET,
						"deletion_protection":                       "false",
						"enable_rrsa":                               "false",
						"new_nat_gateway":                           "true",
						"endpoint_public_access_enabled":            "true",
						"resource_group_id":                         acctest.CHECKSET,
						"security_group_id":                         acctest.CHECKSET,
						"vswitch_ids.#":                             "1",
						"time_zone":                                 getTimezone(os.Getenv("ALICLOUD_REGION")),
						"cluster_spec":                              "ack.standard",
						"custom_san":                                "www.terraform.io,1.1.1.1",
						"maintenance_window.#":                      "1",
						"maintenance_window.0.enable":               "true",
						"maintenance_window.0.maintenance_time":     "2024-10-15T12:31:00.000+08:00",
						"maintenance_window.0.duration":             "3h",
						"maintenance_window.0.weekly_period":        "Thursday",
						"operation_policy.#":                        "1",
						"operation_policy.0.cluster_auto_upgrade.#": "1",
						"operation_policy.0.cluster_auto_upgrade.0.enabled": "true",
						"operation_policy.0.cluster_auto_upgrade.0.channel": "patch",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"load_balancer_spec", "new_nat_gateway", "private_zone", "sls_project_name", "service_discovery_types", "logging_type", "time_zone", "addons", "cluster_ca_cert", "client_key", "client_cert"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_spec": "ack.pro.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_spec": "ack.pro.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "${data.alicloud_cs_kubernetes_version.version-128.metadata.0.version}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_rrsa": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_rrsa": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_san": "www.terraform.io,terraform.test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_san": "www.terraform.io,terraform.test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Platform": "TF",
						"Env":      "Pre",
					},
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":              "2",
						"tags.Platform":       "TF",
						"tags.Env":            "Pre",
						"deletion_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintenance_window": []map[string]string{{"enable": "false", "maintenance_time": "2024-10-15T11:31:00.000+08:00", "duration": "5h", "weekly_period": "Monday,Thursday"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintenance_window.#":                  "1",
						"maintenance_window.0.enable":           "false",
						"maintenance_window.0.maintenance_time": "2024-10-15T11:31:00.000+08:00",
						"maintenance_window.0.duration":         "5h",
						"maintenance_window.0.weekly_period":    "Monday,Thursday",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"operation_policy": []map[string]interface{}{
						{
							"cluster_auto_upgrade": []map[string]interface{}{
								{
									"enabled": "false",
									"channel": "rapid",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operation_policy.#":                                "1",
						"operation_policy.0.cluster_auto_upgrade.#":         "1",
						"operation_policy.0.cluster_auto_upgrade.0.enabled": "false",
						"operation_policy.0.cluster_auto_upgrade.0.channel": "rapid",
					}),
				),
			},
			{
				// retain_resources is a delete-time, write-only input: it is only
				// sent in the Delete request and is never read back into state, so
				// it is set here purely to exercise the config path without a
				// read-back assertion. It is reset to empty in the next step so the
				// final destroy is unaffected.
				Config: testAccConfig(map[string]interface{}{
					"retain_resources": []string{"${alicloud_vswitch.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retain_resources": []string{},
					"delete_options": []map[string]interface{}{
						{
							"delete_mode":   "delete",
							"resource_type": "SLB",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "ALB",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "PrivateZone",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "SLS_Data",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "SLS_ControlPlane",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{})),
			},
		},
	})
}

func TestAccAliCloudCSServerlessKubernetes_encryption(t *testing.T) {
	var v *aliyungocs.ServerlessClusterResponse
	resourceId := "alicloud_cs_serverless_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, csServerlessKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccserverlesskubernetes-encryption-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSServerlessKubernetesConfigDependenceEncryption)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                           name,
					"version":                        "${data.alicloud_cs_kubernetes_version.version-126.metadata.0.version}",
					"vpc_id":                         "${alicloud_vpc.default.id}",
					"vswitch_ids":                    []string{"${alicloud_vswitch.default.id}"},
					"security_group_id":              "${alicloud_security_group.default.id}",
					"new_nat_gateway":                "true",
					"deletion_protection":            "false",
					"endpoint_public_access_enabled": "true",
					"load_balancer_spec":             "slb.s2.small",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"service_cidr":                   "10.0.1.0/24",
					"private_zone":                   "true",
					"cluster_spec":                   "ack.pro.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"version":                        acctest.CHECKSET,
						"deletion_protection":            "false",
						"endpoint_public_access_enabled": "true",
						"cluster_spec":                   "ack.pro.small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encryption_provider_key": "${data.alicloud_kms_keys.default.keys[0].key_id}",
					"disable_encryption":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encryption_provider_key": acctest.CHECKSET,
						"disable_encryption":      "false",
					}),
				),
			},
			{
				PreConfig: func() { time.Sleep(5 * time.Minute) },
				Config: testAccConfig(map[string]interface{}{
					"encryption_provider_key": "",
					"disable_encryption":      "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disable_encryption": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"load_balancer_spec", "new_nat_gateway", "private_zone", "sls_project_name", "service_discovery_types", "logging_type", "time_zone", "addons", "cluster_ca_cert", "client_key", "client_cert"},
			},
		},
	})
}

func TestAccAliCloudCSServerlessKubernetesAuto(t *testing.T) {
	var v *aliyungocs.ServerlessClusterResponse
	resourceId := "alicloud_cs_serverless_kubernetes.auto"
	ra := acctest.ResourceAttrInit(resourceId, csServerlessKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	namePrefix := "tf-testaccserverlesskubernetes"
	name := fmt.Sprintf("%s-%d", namePrefix, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSServerlessKubernetesConfigDependenceAuto)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name_prefix":                    namePrefix,
					"zone_id":                        "${data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id}",
					"cluster_spec":                   "ack.pro.small",
					"new_nat_gateway":                "true",
					"deletion_protection":            "false",
					"enable_rrsa":                    "true",
					"endpoint_public_access_enabled": "true",
					"load_balancer_spec":             "slb.s1.small",
					"service_discovery_types":        []string{"PrivateZone"},
					"tags": map[string]string{
						"Platform": "TF",
					},
					"service_cidr": "10.0.1.0/24",
					"time_zone":    getTimezone(os.Getenv("ALICLOUD_REGION")),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           acctest.CHECKSET,
						"name_prefix":                    namePrefix,
						"cluster_spec":                   "ack.pro.small",
						"vpc_id":                         acctest.CHECKSET,
						"deletion_protection":            "false",
						"enable_rrsa":                    "true",
						"new_nat_gateway":                "true",
						"endpoint_public_access_enabled": "true",
						"service_discovery_types.#":      "1",
						"resource_group_id":              acctest.CHECKSET,
						"vswitch_ids.#":                  "1",
						"time_zone":                      getTimezone(os.Getenv("ALICLOUD_REGION")),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"load_balancer_spec", "new_nat_gateway", "private_zone", "sls_project_name", "service_discovery_types", "logging_type", "addons", "zone_id", "name_prefix"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_zone": "Asia/Tokyo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_zone": "Asia/Tokyo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Platform": "TF",
						"Env":      "Pre",
					},
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":              "2",
						"tags.Platform":       "TF",
						"tags.Env":            "Pre",
						"deletion_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_options": []map[string]interface{}{
						{
							"delete_mode":   "delete",
							"resource_type": "SLB",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "ALB",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "PrivateZone",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "SLS_Data",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "SLS_ControlPlane",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{})),
			},
		},
	})
}

func resourceCSServerlessKubernetesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_cs_kubernetes_version" "version-126" {
  cluster_type       = "Kubernetes"
  kubernetes_version = "1.26"
  profile            = "Serverless"
}

data "alicloud_cs_kubernetes_version" "version-128" {
  cluster_type       = "Kubernetes"
  kubernetes_version = "1.28"
  profile            = "Serverless"
}

resource "alicloud_vpc" "default" {
	vpc_name   = var.name
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id            = alicloud_vpc.default.id
	cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
	zone_id           = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
	vswitch_name      = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}
`, name)
}

func resourceCSServerlessKubernetesConfigDependenceEncryption(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_cs_kubernetes_version" "version-126" {
  cluster_type       = "Kubernetes"
  kubernetes_version = "1.26"
  profile            = "Serverless"
}

resource "alicloud_vpc" "default" {
	vpc_name   = var.name
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id            = alicloud_vpc.default.id
	cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
	zone_id           = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
	vswitch_name      = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

data "alicloud_kms_keys" "default" {
  filters = "[{\"Key\":\"KeyState\",\"Values\":[\"Enabled\"]},{\"Key\":\"KeySpec\",\"Values\":[\"Aliyun_AES_256\"]},{\"Key\":\"KeyUsage\",\"Values\":[\"ENCRYPT/DECRYPT\"]},{\"Key\":\"CreatorType\",\"Values\":[\"User\"]}]"
}
`, name)
}

func resourceCSServerlessKubernetesConfigDependenceAuto(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}
`, name)
}

var csServerlessKubernetesBasicMap = map[string]string{
	"new_nat_gateway":                "true",
	"deletion_protection":            "false",
	"endpoint_public_access_enabled": "true",
}

func resourceCSServerlessKubernetesConfigDependenceOrder(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_cs_kubernetes_version" "version-126" {
  cluster_type       = "Kubernetes"
  kubernetes_version = "1.26"
  profile            = "Serverless"
}

resource "alicloud_vpc" "default" {
	vpc_name   = var.name
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id            = alicloud_vpc.default.id
	cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
	zone_id           = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
	vswitch_name      = var.name
}

resource "alicloud_vswitch" "secondary" {
	vpc_id            = alicloud_vpc.default.id
	cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 9)
	zone_id           = data.alicloud_enhanced_nat_available_zones.enhanced.zones.1.zone_id
	vswitch_name      = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}
`, name)
}

func TestAccAliCloudCSServerlessKubernetes_collectionOrder(t *testing.T) {
	var v *aliyungocs.ServerlessClusterResponse
	resourceId := "alicloud_cs_serverless_kubernetes.collectionorder"
	ra := acctest.ResourceAttrInit(resourceId, csServerlessKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccserverlesskubernetes-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSServerlessKubernetesConfigDependenceOrder)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                           name,
					"version":                        "${data.alicloud_cs_kubernetes_version.version-126.metadata.0.version}",
					"vpc_id":                         "${alicloud_vpc.default.id}",
					"vswitch_ids":                    []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.secondary.id}"},
					"security_group_id":              "${alicloud_security_group.default.id}",
					"new_nat_gateway":                "true",
					"deletion_protection":            "false",
					"endpoint_public_access_enabled": "true",
					"service_cidr":                   "10.0.1.0/24",
					"cluster_spec":                   "ack.standard",
					"service_discovery_types":        []string{"PrivateZone", "CoreDNS"},
					"addons": []map[string]string{
						{"name": "managed-arms-prometheus", "disabled": "false"},
						{"name": "metrics-server", "disabled": "false"},
					},
					"retain_resources": []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.secondary.id}"},
					"delete_options": []map[string]interface{}{
						{"delete_mode": "delete", "resource_type": "SLB"},
						{"delete_mode": "delete", "resource_type": "ALB"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#":             "2",
						"service_discovery_types.#": "2",
					}),
				),
			},
			// service_discovery_types: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"service_discovery_types": []string{"CoreDNS", "PrivateZone"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			// addons: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"addons": []map[string]string{
						{"name": "metrics-server", "disabled": "false"},
						{"name": "managed-arms-prometheus", "disabled": "false"},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			// retain_resources: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"retain_resources": []string{"${alicloud_vswitch.secondary.id}", "${alicloud_vswitch.default.id}"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			// delete_options: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_options": []map[string]interface{}{
						{"delete_mode": "delete", "resource_type": "ALB"},
						{"delete_mode": "delete", "resource_type": "SLB"},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			// vswitch_ids: verify reorder convergence (ForceNew, applied last)
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{"${alicloud_vswitch.secondary.id}", "${alicloud_vswitch.default.id}"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}
