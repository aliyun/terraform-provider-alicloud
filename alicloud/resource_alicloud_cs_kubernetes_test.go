package alicloud

import (
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// # Generate a CA cert pair.
//echo '{"CN":"CA","key":{"algo":"rsa","size":2048}}' | cfssl gencert -initca - | cfssljson -bare ca -
//echo '{"signing":{"default":{"expiry":"438000h","usages":["signing","key encipherment","server auth","client auth"]}}}' > ca-config.json

//# Use CA cert to sign a client cert pair.
//export ADDRESS=
//export NAME=kubernetes-admin
//echo '{"CN":"'$NAME'","O": "system:masters","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME

const caCert = `-----BEGIN CERTIFICATE-----
MIIC6jCCAdKgAwIBAgIUfSCqJB17C26e20n1wp3QZ0ypbsQwDQYJKoZIhvcNAQEL
BQAwDTELMAkGA1UEAxMCQ0EwHhcNMTkwMTA5MDcwNTAwWhcNMjQwMTA4MDcwNTAw
WjANMQswCQYDVQQDEwJDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AMwEC6JM3KxwWcg8l+bpehMBGAE/gUzglYaSwZt0SxlTYyLjx16HIUZcv2JQX5kP
D0jSLqTrd6C2lc3EaLGyU9SLAdjB6uV5jDKeCJhbAXJEtdHiaI5SuJPd1f/RwUym
7aUcG9puLN18203zvfp+Ot4uaoKlUd/sq+VREiojEz5oGbRrrHIVMD4VyqZidmfL
bOG2Zfz3XSKwcJEs2EuI7nlXYLEtWm2YDZQlC2goLfDbj/QkChjgpooyrlo9+Pp3
JTpydgrE3aTecrpVRRzioUKuUJ4RsXHqLfdDFVpN0GB0JdDfGdjGaptfUVD5Dn7I
Zfe69kTXmH9qwEALEmF1pl8CAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgEGMA8GA1Ud
EwEB/wQFMAMBAf8wHQYDVR0OBBYEFM/M1RKZToaXmUK434EthsoNsNunMA0GCSqG
SIb3DQEBCwUAA4IBAQA2cuOyLAb3mkYrfJsv8PHDuZ/c6TUNRDdHUpq6ItQRKuu3
a5fhmAcJD5MZp67n1gVzVZsQ95qrsduwnSCnwDBSZJP21vcqdeIaG+mjlg/zXXnw
b3qCqbtk27Yuypw91+3Jza834vzEAUvHQiWgVOiiHzFO5jImQhAosTMV838ae/kd
ws6mhM65UuWFg5HLbdM2J/zrjWrhuAJZgR1Kx2eReleUyDg97Bc0SPTBth28tGvH
UjY7X0eHM5vuv6NUOyElHVteY8oQ1f0f06K5K4lB7lJ1SB/9PdxYv/AawQwqJIQr
iPn9wR9zlLX6d0Zge293YJ/HGOsm7UzI65DxfDfp
-----END CERTIFICATE-----`

func init() {
	resource.AddTestSweepers("alicloud_cs_kubernetes", &resource.Sweeper{
		Name: "alicloud_cs_kubernetes",
		F:    testSweepCSKubernetes,
	})
}

func testSweepCSKubernetes(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	sweepOtherResourceSuffixes := make([]string, 0)

	raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
		return csClient.DescribeClusters("")
	})
	if err != nil {
		return fmt.Errorf("Error retrieving CS Clusters: %s", err)
	}
	clusters, _ := raw.([]cs.ClusterType)
	sweeped := false

	var vpcIds, vswIds, groupIds, slbIds []string
	for _, v := range clusters {
		name := v.Name
		id := v.ClusterID
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping CS Clusters: %s (%s)", name, id)
				continue
			}
		}
		log.Printf("[INFO] Close CS Clusters: %s (%s) deletion protection", name, id)
		invoker := NewInvoker()

		var requestInfo cs.ModifyClusterArgs
		requestInfo.DeletionProtection = false

		if err := invoker.Run(func() error {
			_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.ModifyCluster(id, &requestInfo)
			})
			return err
		}); err != nil {
			log.Printf("[INFO] Close CS Clusters: %s (%s) deletion protection failed", name, id)
		}

		log.Printf("[INFO] Deleting CS Clusters: %s (%s)", name, id)
		sweepOtherResourceSuffixes = append(sweepOtherResourceSuffixes, id)

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := invoker.Run(func() error {
				_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
					return nil, csClient.DeleteKubernetesCluster(id)
				})
				return err
			}); err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CS Clusters (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
		vpcIds = append(vpcIds, v.VPCID)
		vswIds = append(vswIds, strings.Split(v.VSwitchID, ",")...)
		groupIds = append(groupIds, strings.Split(v.SecurityGroupID, ",")...)
		slbIds = append(slbIds, strings.Split(v.ExternalLoadbalancerID, ",")...)
	}
	if sweeped {
		// Waiting 30 seconds to eusure these swarms have been deleted.
		time.Sleep(30 * time.Second)
	}
	// Currently, the CS will retain some resources after the cluster is deleted.
	slbS := SlbService{client}
	for _, id := range slbIds {
		if err := slbS.sweepSlb(id); err != nil {
			log.Printf("[ERROR] Failed to deleting slb %s: %s", id, WrapError(err))
		}
	}
	ecsS := EcsService{client}
	for _, id := range groupIds {
		if err := ecsS.sweepSecurityGroup(id); err != nil {
			log.Printf("[ERROR] Failed to deleting SG %s: %s", id, WrapError(err))
		}
	}
	vpcS := VpcService{client}
	for _, id := range vswIds {
		if err := vpcS.sweepVSwitch(id); err != nil {
			log.Printf("[ERROR] Failed to deleting VSW %s: %s", id, WrapError(err))
		}
	}
	for _, id := range vpcIds {
		request := vpc.CreateDescribeNatGatewaysRequest()
		request.VpcId = id
		if raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNatGateways(request)
		}); err != nil {
			log.Printf("[ERROR] %#v", WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR))
		} else {
			response, _ := raw.(vpc.DescribeNatGatewaysResponse)
			for _, nat := range response.NatGateways.NatGateway {
				if err := vpcS.sweepNatGateway(nat.NatGatewayId); err != nil {
					log.Printf("[ERROR] Failed to delete nat gateway %s: %s", nat.Name, err)
				}
			}
		}
		if err := vpcS.sweepVpc(id); err != nil {
			log.Printf("[ERROR] Failed to deleting VPC %s: %s", id, WrapError(err))
		}
	}
	// Sweep the log projects which created by K8s Service
	testSweepLogProjectsWithPrefixAndSuffix(region, []string{}, sweepOtherResourceSuffixes)
	return nil
}

func TestAccAliCloudCSKubernetes_basic(t *testing.T) {
	var v *cs.KubernetesClusterDetail

	resourceId := "alicloud_cs_kubernetes.default"
	ra := resourceAttrInit(resourceId, csKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKubernetes_basic-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSKubernetesConfigDependence)

	tmpCAFile, err := os.CreateTemp("", "tf-acc-alicloud-cs-kubernetes-userca")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpCAFile.Name())
	err = os.WriteFile(tmpCAFile.Name(), []byte(caCert), 0644)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                           name,
					"version":                        "${data.alicloud_cs_kubernetes_version.kubernetes_versions.metadata.2.version}",
					"master_vswitch_ids":             []string{"${local.vswitch_id}", "${local.vswitch_id}", "${local.vswitch_id}"},
					"master_instance_types":          []string{"${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}"},
					"master_disk_category":           "cloud_essd",
					"master_disk_performance_level":  "PL0",
					"master_disk_size":               "80",
					"master_disk_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.default.id}",
					"key_name":                       "${alicloud_key_pair.default.key_pair_name}",
					"pod_cidr":                       "10.72.0.0/16",
					"service_cidr":                   "172.18.0.0/16",
					"enable_ssh":                     "false",
					"load_balancer_spec":             "slb.s1.small",
					"install_cloud_monitor":          "true",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"security_group_id":              "${alicloud_security_group.default.id}",
					"deletion_protection":            "false",
					"timezone":                       "Asia/Shanghai",
					"os_type":                        "Linux",
					"platform":                       "AliyunLinux3",
					"image_id":                       "aliyun_3_x64_20G_alibase_20240819.vhd",
					"runtime":                        map[string]interface{}{"name": "containerd", "version": "1.6.20"},
					"node_name_mode":                 "customized,aliyun.com-,5,-test",
					"cluster_domain":                 "cluster.local",
					"custom_san":                     "www.terraform.io",
					"rds_instances":                  []string{"${alicloud_db_instance.default.id}"},
					"tags":                           map[string]string{"Platform": "TF"},
					"proxy_mode":                     "ipvs",
					"new_nat_gateway":                "true",
					"slb_internet_enabled":           "true",
					"node_cidr_mask":                 "24",
					"api_audiences":                  []string{"https://kubernetes.default.svc"},
					"service_account_issuer":         "https://kubernetes.default.svc",
					"user_ca":                        tmpCAFile.Name(),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                          name,
						"master_disk_category":          "cloud_essd",
						"master_disk_performance_level": "PL0",
						"master_disk_size":              "80",
						"key_name":                      name,
						"pod_cidr":                      "10.72.0.0/16",
						"service_cidr":                  "172.18.0.0/16",
						"enable_ssh":                    "false",
						"install_cloud_monitor":         "true",
						"resource_group_id":             CHECKSET,
						"security_group_id":             CHECKSET,
						"deletion_protection":           "false",
						"timezone":                      "Asia/Shanghai",
						"os_type":                       "Linux",
						"platform":                      "AliyunLinux3",
						"cluster_domain":                "cluster.local",
						"custom_san":                    "www.terraform.io",
						"rds_instances.#":               "1",
						"tags.%":                        "1",
						"tags.Platform":                 "TF",
						"proxy_mode":                    "ipvs",
						"new_nat_gateway":               "true",
						"nat_gateway_id":                CHECKSET,
						"slb_internet_enabled":          "true",
						"node_cidr_mask":                "24",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"new_nat_gateway", "password", "user_ca", "rds_instances",
					"cluster_ca_cert", "client_key", "client_cert", "kms_encryption_context", "kms_encrypted_password",
					"retain_resources", "name_prefix", "enable_ssh", "timezone", "runtime",
					"api_audiences", "service_account_issuer", "load_balancer_spec", "platform",
				},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Platform": "TF",
						"Env":      "Pre",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":        "2",
						"tags.Platform": "TF",
						"tags.Env":      "Pre",
					}),
				),
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
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "true",
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				// upgrade
				Config: testAccConfig(map[string]interface{}{
					"version": "${data.alicloud_cs_kubernetes_version.kubernetes_versions.metadata.1.version}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": CHECKSET,
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
							"resource_type": "SLS_Data",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "SLS_ControlPlane",
						},
						{
							"delete_mode":   "delete",
							"resource_type": "ALB",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{})),
			},
		},
	})
}

func TestAccAliCloudCSKubernetes_prepaid(t *testing.T) {
	var v *cs.KubernetesClusterDetail

	resourceId := "alicloud_cs_kubernetes.default"
	ra := resourceAttrInit(resourceId, csKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKubernetes_prepaid-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSKubernetesConfigDependence_essd)

	clusterCaCertFile, clientCertFile, clientKeyFile, err := CreateTempFiles()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(clientCertFile.Name())
	defer os.Remove(clientKeyFile.Name())
	defer os.Remove(clusterCaCertFile.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name_prefix":                  "tf-testAccKubernetes_prepaid",
					"master_vswitch_ids":           []string{"${local.vswitch_id}", "${local.vswitch_id}", "${local.vswitch_id}"},
					"master_instance_types":        []string{"${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}"},
					"master_disk_category":         "cloud_ssd",
					"master_auto_renew":            "true",
					"master_auto_renew_period":     "1",
					"master_instance_charge_type":  "PrePaid",
					"master_period":                "1",
					"master_period_unit":           "Month",
					"pod_vswitch_ids":              []string{"${local.vswitch_id}"},
					"password":                     "Yourpassword1234",
					"service_cidr":                 "172.18.0.0/16",
					"enable_ssh":                   "false",
					"load_balancer_spec":           "slb.s1.small",
					"install_cloud_monitor":        "true",
					"resource_group_id":            "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"deletion_protection":          "false",
					"timezone":                     "Asia/Shanghai",
					"os_type":                      "Linux",
					"platform":                     "AliyunLinux3",
					"cluster_domain":               "cluster.local",
					"custom_san":                   "www.terraform.io",
					"proxy_mode":                   "ipvs",
					"new_nat_gateway":              "true",
					"is_enterprise_security_group": "true",
					"addons":                       []map[string]string{{"name": "terway-eniip", "config": "", "version": "", "disabled": "false"}},
					"cluster_ca_cert":              clusterCaCertFile.Name(),
					"client_key":                   clientKeyFile.Name(),
					"client_cert":                  clientCertFile.Name(),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                         CHECKSET,
						"master_disk_category":         "cloud_ssd",
						"master_auto_renew":            "true",
						"master_auto_renew_period":     "1",
						"master_instance_charge_type":  "PrePaid",
						"master_period":                "1",
						"master_period_unit":           "Month",
						"password":                     "Yourpassword1234",
						"service_cidr":                 "172.18.0.0/16",
						"enable_ssh":                   "false",
						"install_cloud_monitor":        "true",
						"resource_group_id":            CHECKSET,
						"deletion_protection":          "false",
						"timezone":                     "Asia/Shanghai",
						"os_type":                      "Linux",
						"platform":                     "AliyunLinux3",
						"cluster_domain":               "cluster.local",
						"custom_san":                   "www.terraform.io",
						"proxy_mode":                   "ipvs",
						"new_nat_gateway":              "true",
						"nat_gateway_id":               CHECKSET,
						"is_enterprise_security_group": "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"new_nat_gateway", "password", "user_ca", "runtime",
					"rds_instances", "cluster_ca_cert", "client_key", "client_cert", "kms_encryption_context",
					"kms_encrypted_password", "retain_resources", "name_prefix", "enable_ssh", "timezone", "addons",
					"load_balancer_spec", "pod_vswitch_ids", "slb_internet_enabled", "platform",
				},
			},
		},
	})
}

func CreateTempFiles() (*os.File, *os.File, *os.File, error) {
	clientCertFile, err := os.CreateTemp("", "client-cert")
	if err != nil {
		return nil, nil, nil, err
	}
	clientKeyFile, err := os.CreateTemp("", "client-key")
	if err != nil {
		return nil, nil, nil, err
	}
	clusterCaCertFile, err := os.CreateTemp("", "cluster-ca-cert")
	if err != nil {
		return nil, nil, nil, err
	}
	return clusterCaCertFile, clientCertFile, clientKeyFile, nil
}

func resourceCSKubernetesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Master"
  system_disk_category = "cloud_essd"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_cs_kubernetes_version" "kubernetes_versions" {
  cluster_type       = "Kubernetes"
}


data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "192.168.0.0/16"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  instance_charge_type = "Postpaid"
  instance_name        = "tf-testacckubernetes"
  vswitch_id           = local.vswitch_id
  monitoring_period    = "60"
}

resource "alicloud_key_pair" "default" {
  key_pair_name = var.name
  timeouts {
    delete = "15m"
  }
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_ecs_auto_snapshot_policy" "default" {
  name            = var.name
  repeat_weekdays = ["1", "2", "3"]
  retention_days  = -1
  time_points     = ["1", "22", "23"]
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  cluster_id                    = alicloud_cs_kubernetes.default.id
  name                          = var.name
  vswitch_ids                   = [local.vswitch_id]
  instance_types                = [data.alicloud_instance_types.default.instance_types.0.id]
  password                      = "Test12345"
  system_disk_size              = 50
  system_disk_category          = "cloud_essd"
  system_disk_performance_level = "PL0"
  desired_size                  = 2
}
`, name)
}

func resourceCSKubernetesConfigDependence_essd(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 4
  memory_size          = 8
  system_disk_category = "cloud_essd"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "192.168.0.0/16"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  cluster_id                    = alicloud_cs_kubernetes.default.id
  name                          = var.name
  vswitch_ids                   = [local.vswitch_id]
  instance_types                = [data.alicloud_instance_types.default.instance_types.0.id]
  password                      = "Test12345"
  system_disk_size              = 50
  system_disk_category          = "cloud_essd"
  system_disk_performance_level = "PL0"
  desired_size                  = 2
}
`, name)
}

var csKubernetesBasicMap = map[string]string{
	"name":                               CHECKSET,
	"is_enterprise_security_group":       CHECKSET,
	"security_group_id":                  CHECKSET,
	"image_id":                           CHECKSET,
	"version":                            CHECKSET,
	"platform":                           CHECKSET,
	"certificate_authority.cluster_cert": CHECKSET,
	"certificate_authority.client_cert":  CHECKSET,
	"certificate_authority.client_key":   CHECKSET,
	"connections.api_server_internet":    CHECKSET,
	"connections.api_server_intranet":    CHECKSET,
	"connections.master_public_ip":       CHECKSET,
	"connections.service_domain":         CHECKSET,
	"worker_ram_role_name":               CHECKSET,
	"node_name_mode":                     CHECKSET,
	"master_nodes.#":                     "3",
	"vpc_id":                             CHECKSET,
	"resource_group_id":                  CHECKSET,
	"slb_internet":                       CHECKSET,
	"slb_intranet":                       CHECKSET,
	"slb_id":                             CHECKSET,
}

func TestUnit_parseRRSAMetadata(t *testing.T) {
	type args struct {
		meta string
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		{
			name: "empty error",
			args: args{
				meta: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid error",
			args: args{
				meta: `
{
	"RRSAConfig": 1234
}`,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid enabled",
			args: args{
				meta: `{
	"RRSAConfig": {
		"enabled": true,
		"issuer": "https://example.com,https://kubernetes.default.svc,kubernetes.default.svc",
		"oidc_name": "ack-rrsa-c12345",
		"oidc_arn": "acs:ram::12345:oidc-provider/ack-rrsa-c12345"
	}
}`,
			},
			want: []map[string]interface{}{{
				"enabled":                true,
				"rrsa_oidc_issuer_url":   "https://example.com",
				"ram_oidc_provider_name": "ack-rrsa-c12345",
				"ram_oidc_provider_arn":  "acs:ram::12345:oidc-provider/ack-rrsa-c12345",
			}},
			wantErr: false,
		},
		{
			name: "valid not enabled",
			args: args{
				meta: `{
	"RRSAConfig": {
		"enabled": false,
		"issuer": "",
		"oidc_name": "",
		"oidc_arn": ""
	}
}`,
			},
			want: []map[string]interface{}{{
				"enabled":                false,
				"rrsa_oidc_issuer_url":   "",
				"ram_oidc_provider_name": "",
				"ram_oidc_provider_arn":  "",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := flattenRRSAMetadata(tt.args.meta)
			if tt.wantErr && err == nil {
				t.Errorf("flattenRRSAMetadata(%v) want error got nil", tt.args.meta)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("flattenRRSAMetadata(%v) want %v got %v", tt.args.meta, tt.want, got)
			}
		})
	}
}
