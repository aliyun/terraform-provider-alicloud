package alicloud

import (
	"log"
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

func TestAccAlicloudCSKubernetes_basic(t *testing.T) {
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
					"name":                        name,
					"master_vswitch_ids":          []string{"${local.vswitch_id}", "${local.vswitch_id}", "${local.vswitch_id}"},
					"worker_vswitch_ids":          []string{"${local.vswitch_id}"},
					"master_instance_types":       []string{"${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}"},
					"worker_instance_types":       []string{"${data.alicloud_instance_types.default1.instance_types.0.id}"},
					"master_disk_category":        "cloud_ssd",
					"worker_disk_size":            "50",
					"password":                    "Yourpassword1234",
					"pod_cidr":                    "10.72.0.0/16",
					"service_cidr":                "172.18.0.0/16",
					"enable_ssh":                  "true",
					"load_balancer_spec":          "slb.s2.small",
					"install_cloud_monitor":       "true",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"deletion_protection":         "false",
					"timezone":                    "Asia/Shanghai",
					"os_type":                     "Linux",
					"platform":                    "CentOS",
					"node_port_range":             "30000-32767",
					"cluster_domain":              "cluster.local",
					"custom_san":                  "www.terraform.io",
					"rds_instances":               []string{"${alicloud_db_instance.default.id}"},
					"taints":                      []map[string]string{{"key": "tf-key1", "value": "tf-value1", "effect": "NoSchedule"}},
					"proxy_mode":                  "ipvs",
					"new_nat_gateway":             "true",
					"worker_instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        name,
						"master_disk_category":        "cloud_ssd",
						"worker_disk_size":            "50",
						"password":                    "Yourpassword1234",
						"pod_cidr":                    "10.72.0.0/16",
						"service_cidr":                "172.18.0.0/16",
						"enable_ssh":                  "true",
						"install_cloud_monitor":       "true",
						"resource_group_id":           CHECKSET,
						"deletion_protection":         "false",
						"timezone":                    "Asia/Shanghai",
						"os_type":                     "Linux",
						"platform":                    "CentOS",
						"node_port_range":             "30000-32767",
						"cluster_domain":              "cluster.local",
						"custom_san":                  "www.terraform.io",
						"rds_instances.#":             "1",
						"taints.#":                    "1",
						"taints.0.key":                "tf-key1",
						"taints.0.value":              "tf-value1",
						"taints.0.effect":             "NoSchedule",
						"proxy_mode":                  "ipvs",
						"new_nat_gateway":             "true",
						"worker_instance_charge_type": "PostPaid",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"name", "new_nat_gateway", "pod_cidr",
					"service_cidr", "enable_ssh", "password", "install_cloud_monitor", "user_ca", "force_update",
					"master_disk_category", "master_disk_size", "master_instance_charge_type", "master_instance_types",
					"node_cidr_mask", "slb_internet_enabled", "vswitch_ids", "worker_disk_category", "worker_disk_size",
					"worker_instance_charge_type", "worker_instance_types", "log_config", "worker_number",
					"worker_data_disk_category", "worker_data_disk_size", "master_vswitch_ids", "worker_vswitch_ids", "exclude_autoscaler_nodes", "cpu_policy", "proxy_mode", "cluster_domain", "custom_san", "node_port_range", "os_type", "platform", "timezone", "runtime", "taints", "rds_instances", "load_balancer_spec"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"install_cloud_monitor": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"install_cloud_monitor": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "tf-dedicated-k8s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-dedicated-k8s",
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
		},
	})
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
	availability_zone          = data.alicloud_zones.default.zones.0.id
	cpu_core_count             = 4
	memory_size                = 8
	kubernetes_node_role       = "Worker"
    instance_type_family       = "ecs.sn1ne"
}

data "alicloud_instance_types" "default1" {
  availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
  instance_type_family       = "ecs.sn1ne"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
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

resource "alicloud_snapshot_policy" "default" {
  name            = "${var.name}"
  repeat_weekdays = ["1", "2", "3"]
  retention_days  = -1
  time_points     = ["1", "22", "23"]
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
	availability_zone          = data.alicloud_zones.default.zones.0.id
	cpu_core_count             = 4
	memory_size                = 8
	kubernetes_node_role       = "Worker"
	system_disk_category 	   = "cloud_essd"
    instance_type_family       = "ecs.sn1ne"
}

data "alicloud_instance_types" "default1" {
  availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count       = 4
  memory_size          = 8
  system_disk_category = "cloud_essd"
  kubernetes_node_role = "Worker"
  instance_type_family       = "ecs.sn1ne"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
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

resource "alicloud_snapshot_policy" "default" {
  name            = "${var.name}"
  repeat_weekdays = ["1", "2", "3"]
  retention_days  = -1
  time_points     = ["1", "22", "23"]
}
	`, name)
}

var csKubernetesBasicMap = map[string]string{
	"name": CHECKSET,
}

func Test_parseRRSAMetadata(t *testing.T) {
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
