package cs_test

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cs"
	aliyungocs "github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Configurable TypeList attributes on alicloud_cs_kubernetes are covered by
// convergent reorder tests in TestAccAliCloudCSKubernetes_collectionOrder below.

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

var csTestEndpoints sync.Map
var csTestSignVersion sync.Map

// csTestClientForRegion rebuilds the package-level sweeper client on the
// exported connectivity API: the flat-suite sharedClientForRegion is not
// reachable from this package.
func csTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	var accessKey, secretKey string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}
	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}
	securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &csTestEndpoints,
		SignVersion: &csTestSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential

	return conf.Client()
}

func init() {
	resource.AddTestSweepers("alicloud_cs_kubernetes", &resource.Sweeper{
		Name: "alicloud_cs_kubernetes",
		F:    testSweepCSKubernetes,
	})
}

func testSweepCSKubernetes(region string) error {
	rawClient, err := csTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf-example",
		"tf_example",
		"terraform-example",
	}
	sweepOtherResourceSuffixes := make([]string, 0)

	raw, err := client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
		return csClient.DescribeClusters("")
	})
	if err != nil {
		return fmt.Errorf("Error retrieving CS Clusters: %s", err)
	}
	clusters, _ := raw.([]aliyungocs.ClusterType)
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
		invoker := helper.NewInvoker()

		var requestInfo aliyungocs.ModifyClusterArgs
		requestInfo.DeletionProtection = false

		if err := invoker.Run(func() error {
			_, err := client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
				return nil, csClient.ModifyCluster(id, &requestInfo)
			})
			return err
		}); err != nil {
			log.Printf("[INFO] Close CS Clusters: %s (%s) deletion protection failed", name, id)
		}

		log.Printf("[INFO] Deleting CS Clusters: %s (%s)", name, id)
		sweepOtherResourceSuffixes = append(sweepOtherResourceSuffixes, id)

		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			if err := invoker.Run(func() error {
				_, err := client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
					return nil, csClient.DeleteKubernetesCluster(id)
				})
				return err
			}); err != nil {
				return retry.RetryableError(err)
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
	// The flat-layout sweep helpers (SlbService.sweepSlb, EcsService.sweepSecurityGroup,
	// VpcService.sweepVpc/sweepVSwitch/sweepNatGateway, testSweepLogProjectsWithPrefixAndSuffix)
	// live in the root alicloud test package and are unreachable from the split
	// cs_test package. They are reproduced below as free functions on the exported
	// connectivity client (nlb/dts split-package precedent), preserving the original
	// sweep semantics verbatim.
	for _, id := range slbIds {
		if err := sweepSlb(client, id); err != nil {
			log.Printf("[ERROR] Failed to deleting slb %s: %s", id, errs.WrapError(err))
		}
	}
	for _, id := range groupIds {
		if err := sweepSecurityGroup(client, id); err != nil {
			log.Printf("[ERROR] Failed to deleting SG %s: %s", id, errs.WrapError(err))
		}
	}
	for _, id := range vswIds {
		if err := sweepVSwitch(client, id); err != nil {
			log.Printf("[ERROR] Failed to deleting VSW %s: %s", id, errs.WrapError(err))
		}
	}
	for _, id := range vpcIds {
		request := vpc.CreateDescribeNatGatewaysRequest()
		request.VpcId = id
		if raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNatGateways(request)
		}); err != nil {
			log.Printf("[ERROR] %#v", errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR))
		} else {
			response, _ := raw.(vpc.DescribeNatGatewaysResponse)
			for _, nat := range response.NatGateways.NatGateway {
				if err := sweepNatGateway(client, nat.NatGatewayId); err != nil {
					log.Printf("[ERROR] Failed to delete nat gateway %s: %s", nat.Name, err)
				}
			}
		}
		if err := sweepVpc(client, id); err != nil {
			log.Printf("[ERROR] Failed to deleting VPC %s: %s", id, errs.WrapError(err))
		}
	}
	// Sweep the log projects which created by K8s Service
	sweepLogProjectsWithPrefixAndSuffix(client, []string{}, sweepOtherResourceSuffixes)
	return nil
}

// sweepAll mirrors the root helper of the same name so the migrated sweeper
// keeps its original gate without referencing the flat layout (nlb/dts
// split-package precedent).
func sweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}

// sweepSlb mirrors (SlbService).sweepSlb from the flat layout's
// service_alicloud_common_test.go on the exported connectivity client.
func sweepSlb(client *connectivity.AliyunClient, id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Set SLB DeleteProtection to off before deleting %s ...", id)
	action := "SetLoadBalancerDeleteProtection"
	request := map[string]interface{}{
		"RegionId":         client.RegionId,
		"LoadBalancerId":   id,
		"DeleteProtection": "off",
	}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 10*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] Set SLB %s DeleteProtection to off failed.", id)
	}
	log.Printf("[DEBUG] Deleting SLB %s ...", id)
	delRequest := map[string]interface{}{
		"RegionId":       client.RegionId,
		"LoadBalancerId": id,
	}
	action = "DeleteLoadBalancer"
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("Slb", "2014-05-15", action, nil, delRequest, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return errs.WrapError(err)
}

// sweepSecurityGroup mirrors (EcsService).sweepSecurityGroup from the flat
// layout's service_alicloud_common_test.go on the exported connectivity client.
func sweepSecurityGroup(client *connectivity.AliyunClient, id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Deleting Security Group %s ...", id)
	request := ecs.CreateDeleteSecurityGroupRequest()
	request.SecurityGroupId = id
	_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteSecurityGroup(request)
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return errs.WrapError(err)
}

// sweepVpc mirrors (VpcService).sweepVpc from the flat layout's
// service_alicloud_common_test.go on the exported connectivity client.
func sweepVpc(client *connectivity.AliyunClient, id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Deleting Vpc %s ...", id)
	request := vpc.CreateDeleteVpcRequest()
	request.VpcId = id
	_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteVpc(request)
	})
	return errs.WrapError(err)
}

// sweepVSwitch mirrors (VpcService).sweepVSwitch from the flat layout's
// service_alicloud_common_test.go on the exported connectivity client.
func sweepVSwitch(client *connectivity.AliyunClient, id string) error {
	var err error
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Deleting Vswitch %s ...", id)
	action := "DeleteVSwitch"
	request := map[string]interface{}{
		"VSwitchId": id,
	}
	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(10*time.Second, func() *retry.RetryError {
		_, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return errs.WrapError(err)
}

// sweepNatGateway mirrors (VpcService).sweepNatGateway from the flat layout's
// service_alicloud_common_test.go on the exported connectivity client.
func sweepNatGateway(client *connectivity.AliyunClient, id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[INFO] Deleting Nat Gateway %s ...", id)
	request := vpc.CreateDeleteNatGatewayRequest()
	request.NatGatewayId = id
	request.Force = requests.NewBoolean(true)
	_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteNatGateway(request)
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return errs.WrapError(err)
}

// sweepLogProjectsWithPrefixAndSuffix mirrors
// testSweepLogProjectsWithPrefixAndSuffix from the flat layout's
// resource_alicloud_log_project_test.go on the exported connectivity client;
// the original rebuilt its own client via sharedClientForRegion, but the
// sweeper already holds an equivalent client, so it is passed in directly.
func sweepLogProjectsWithPrefixAndSuffix(client *connectivity.AliyunClient, prefixes, suffixes []string) error {
	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		return slsClient.ListProject()
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving Log Projects: %s", errs.WrapError(err))
	}
	names, _ := raw.([]string)

	for _, v := range names {
		name := v
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				for _, suffix := range suffixes {
					if strings.HasSuffix(strings.ToLower(name), strings.ToLower(suffix)) {
						skip = false
						break
					}
				}
			}
			// Sweep the project which from the k8s cluster
			if skip && strings.HasPrefix(name, "k8s-log-") {
				k8sId := strings.TrimPrefix(name, "k8s-log-")
				raw, err := client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
					return csClient.DescribeCluster(k8sId)
				})
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
						skip = false
					} else {
						log.Printf("[ERROR] DescribeCluster got an error: %#v", err)
					}
				} else {
					cluster, _ := raw.(aliyungocs.ClusterType)
					if strings.HasPrefix(strings.ToLower(cluster.Name), "tf-testacc") || strings.HasPrefix(strings.ToLower(cluster.Name), "tf_testacc") {
						skip = false
					}
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Log Project: %s", name)
				continue
			}
		}
		log.Printf("[INFO] Deleting Log Project: %s", name)
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteProject(name)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Log Project (%s): %s", name, err)
		}
	}
	return nil
}

func TestAccAliCloudCSKubernetes_basic(t *testing.T) {
	var v *aliyungocs.KubernetesClusterDetail

	resourceId := "alicloud_cs_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, csKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKubernetes_basic-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSKubernetesConfigDependence)

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
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                           name,
					"version":                        "${data.alicloud_cs_kubernetes_version.kubernetes_versions.metadata.2.version}",
					"master_vswitch_ids":             []string{"${local.vswitch_id}", "${local.vswitch_id}", "${local.vswitch_id}"},
					"master_instance_types":          []string{"${var.instance_type}", "${var.instance_type}", "${var.instance_type}"},
					"master_disk_category":           "cloud_essd",
					"master_disk_size":               "80",
					"master_disk_snapshot_policy_id": "${alicloud_ecs_auto_snapshot_policy.default.id}",
					"key_name":                       "${alicloud_key_pair.default.key_pair_name}",
					"pod_cidr":                       "10.72.0.0/16",
					"service_cidr":                   "172.18.0.0/16",
					"enable_ssh":                     "false",
					"install_cloud_monitor":          "true",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"security_group_id":              "${alicloud_security_group.default.id}",
					"deletion_protection":            "false",
					"os_type":                        "Linux",
					"platform":                       "AliyunLinux4ContainerOptimized",
					"image_id":                       "aliyun_4_x64_20G_container_optimized_alibase_20260402.vhd",
					"runtime":                        []map[string]interface{}{{"name": "containerd", "version": "1.6.20"}},
					"node_name_mode":                 "customized,aliyun.com-,5,-test",
					"cluster_domain":                 "cluster.local",
					"custom_san":                     "www.terraform.io",
					"rds_instances":                  []string{"${alicloud_db_instance.default.id}"},
					"tags":                           map[string]string{"Platform": "TF"},
					"proxy_mode":                     "ipvs",
					"new_nat_gateway":                "true",
					"slb_internet_enabled":           "false",
					"node_cidr_mask":                 "24",
					"api_audiences":                  []string{"https://kubernetes.default.svc"},
					"service_account_issuer":         "https://kubernetes.default.svc",
					"user_ca":                        tmpCAFile.Name(),
					"skip_set_certificate_authority": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"master_disk_category":           "cloud_essd",
						"master_disk_size":               "80",
						"key_name":                       name,
						"pod_cidr":                       "10.72.0.0/16",
						"service_cidr":                   "172.18.0.0/16",
						"enable_ssh":                     "false",
						"install_cloud_monitor":          "true",
						"resource_group_id":              acctest.CHECKSET,
						"security_group_id":              acctest.CHECKSET,
						"deletion_protection":            "false",
						"os_type":                        "Linux",
						"platform":                       "AliyunLinux4ContainerOptimized",
						"cluster_domain":                 "cluster.local",
						"custom_san":                     "www.terraform.io",
						"rds_instances.#":                "1",
						"tags.%":                         "1",
						"tags.Platform":                  "TF",
						"proxy_mode":                     "ipvs",
						"new_nat_gateway":                "true",
						"nat_gateway_id":                 acctest.CHECKSET,
						"slb_internet_enabled":           "false",
						"node_cidr_mask":                 "24",
						"skip_set_certificate_authority": "false",
					}),
				),
			},
			{
				// KMS-encrypted login credentials, master disk performance
				// tiers and the delete-time retain list are node-provisioning
				// concerns this suite cannot exercise. key_name from the
				// previous step conflicts with kms_encrypted_password, so the
				// plan stops at validation before the placeholders reach the
				// API; the next config step drops them again.
				Config: testAccConfig(map[string]interface{}{
					"kms_encrypted_password": "kms-tfacc-placeholder",
					// The test config renderer emits map keys verbatim, and `for` is a
					// reserved word in HCL object constructors, so a placeholder key
					// that is a plain identifier keeps the step parseable.
					"kms_encryption_context":        map[string]string{"kms": "tfacc"},
					"master_disk_performance_level": "PL1",
					"retain_resources":              []string{"ecs"},
					"runtime":                       []map[string]interface{}{{"name": "docker", "version": "19.03.15"}},
				}),
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "new_nat_gateway", "password", "user_ca", "rds_instances", "cluster_ca_cert", "client_key", "client_cert", "kms_encryption_context", "kms_encrypted_password", "retain_resources", "name_prefix", "enable_ssh", "runtime", "api_audiences", "service_account_issuer", "load_balancer_spec", "platform", "timezone"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					// The KMS validation step stops at plan time, so its
					// placeholders were never applied; drop them and restore the
					// creation-time runtime before the remaining updates.
					"kms_encrypted_password":        acctest.REMOVEKEY,
					"kms_encryption_context":        acctest.REMOVEKEY,
					"master_disk_performance_level": acctest.REMOVEKEY,
					"retain_resources":              acctest.REMOVEKEY,
					"runtime":                       []map[string]interface{}{{"name": "containerd", "version": "1.6.20"}},
					"addons": []map[string]interface{}{
						{
							"config":   "",
							"disabled": "false",
							"name":     "terway-eniip",
							"version":  "",
						},
					},
					"tags": map[string]string{
						"Platform": "TF",
						"Env":      "Pre",
					},
					"name": name,
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
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				// upgrade
				Config: testAccConfig(map[string]interface{}{
					"version": "1.36.2-aliyun.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": "1.36.2-aliyun.1",
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

// Note that you may need to set slb_internet_enabled=false to run the test
func TestAccAliCloudCSKubernetes_prepaid(t *testing.T) {
	var v *aliyungocs.KubernetesClusterDetail

	resourceId := "alicloud_cs_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, csKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKubernetes_prepaid-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSKubernetesConfigDependence_essd)

	clusterCaCertFile, clientCertFile, clientKeyFile, err := CreateTempFiles()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(clientCertFile.Name())
	defer os.Remove(clientKeyFile.Name())
	defer os.Remove(clusterCaCertFile.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name_prefix":                    "tf-testAccKubernetes_prepaid",
					"master_vswitch_ids":             []string{"${local.vswitch_id}", "${local.vswitch_id}", "${local.vswitch_id}"},
					"master_instance_types":          []string{"${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}"},
					"master_disk_category":           "cloud_essd",
					"master_auto_renew":              "true",
					"master_auto_renew_period":       "1",
					"master_instance_charge_type":    "PrePaid",
					"master_period":                  "1",
					"master_period_unit":             "Month",
					"pod_vswitch_ids":                []string{"${local.vswitch_id}"},
					"password":                       "Yourpassword1234",
					"service_cidr":                   "172.18.0.0/16",
					"enable_ssh":                     "false",
					"install_cloud_monitor":          "true",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"deletion_protection":            "false",
					"os_type":                        "Linux",
					"platform":                       "AliyunLinux4ContainerOptimized",
					"cluster_domain":                 "cluster.local",
					"custom_san":                     "www.terraform.io",
					"proxy_mode":                     "ipvs",
					"new_nat_gateway":                "true",
					"is_enterprise_security_group":   "false",
					"addons":                         []map[string]string{{"name": "terway-eniip", "config": "", "version": "", "disabled": "false"}},
					"cluster_ca_cert":                clusterCaCertFile.Name(),
					"client_key":                     clientKeyFile.Name(),
					"client_cert":                    clientCertFile.Name(),
					"skip_set_certificate_authority": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           acctest.CHECKSET,
						"master_disk_category":           "cloud_essd",
						"master_auto_renew":              "true",
						"master_auto_renew_period":       "1",
						"master_instance_charge_type":    "PrePaid",
						"master_period":                  "1",
						"master_period_unit":             "Month",
						"password":                       "Yourpassword1234",
						"service_cidr":                   "172.18.0.0/16",
						"enable_ssh":                     "false",
						"install_cloud_monitor":          "true",
						"resource_group_id":              acctest.CHECKSET,
						"deletion_protection":            "false",
						"os_type":                        "Linux",
						"platform":                       "AliyunLinux4ContainerOptimized",
						"cluster_domain":                 "cluster.local",
						"custom_san":                     "www.terraform.io",
						"proxy_mode":                     "ipvs",
						"new_nat_gateway":                "true",
						"nat_gateway_id":                 acctest.CHECKSET,
						"skip_set_certificate_authority": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_set_certificate_authority", "is_enterprise_security_group", "new_nat_gateway", "password", "user_ca", "runtime", "rds_instances", "cluster_ca_cert", "client_key", "client_cert", "kms_encryption_context", "kms_encrypted_password", "retain_resources", "name_prefix", "enable_ssh", "addons", "load_balancer_spec", "pod_vswitch_ids", "slb_internet_enabled", "platform", "timezone"},
			},
		},
	})
}

func TestAccAliCloudCSKubernetes_collectionOrder(t *testing.T) {
	var v *aliyungocs.KubernetesClusterDetail

	resourceId := "alicloud_cs_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"name":                          acctest.CHECKSET,
		"version":                       "1.32.7-aliyun.1",
		"timezone":                      "Asia/Shanghai",
		"master_disk_performance_level": "PL1",
		"image_id":                      acctest.CHECKSET,
		"master_vswitch_ids.#":          "3",
		"master_instance_types.#":       "3",
		"api_audiences.#":               "2",
	})

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKubernetes_collectionOrder-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSKubernetesConfigDependenceOrder)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                          name,
					"version":                       "1.32.7-aliyun.1",
					"master_vswitch_ids":            []string{"${local.vswitch_id_0}", "${local.vswitch_id_1}", "${local.vswitch_id_2}"},
					"master_instance_types":         []string{"${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.1.id}", "${data.alicloud_instance_types.default.instance_types.2.id}"},
					"master_disk_category":          "cloud_essd",
					"master_disk_size":              "80",
					"master_disk_performance_level": "PL1",
					"service_cidr":                  "172.18.0.0/16",
					"pod_vswitch_ids":               []string{"${local.vswitch_id_0}", "${local.vswitch_id_1}"},
					"proxy_mode":                    "ipvs",
					"addons": []map[string]string{
						{"name": "terway-eniip", "config": "", "version": "", "disabled": "false"},
						{"name": "nginx-ingress-controller", "config": "{\\\"IngressSlbNetworkType\\\":\\\"internet\\\"}", "version": "", "disabled": "false"},
					},
					"new_nat_gateway":        "false",
					"security_group_id":      "${alicloud_security_group.group.id}",
					"timezone":               "Asia/Shanghai",
					"image_id":               "${data.alicloud_images.default.images.0.id}",
					"api_audiences":          []string{"https://kubernetes.default.svc", "https://kubernetes.default.svc.cluster.local"},
					"service_account_issuer": "https://kubernetes.default.svc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pod_vswitch_ids.#": "2",
						"addons.#":          "2",
					}),
					resource.TestCheckResourceAttr(resourceId, "addons.0.name", "terway-eniip"),
					resource.TestCheckResourceAttr(resourceId, "addons.1.name", "nginx-ingress-controller"),
					resource.TestCheckResourceAttr(resourceId, "addons.1.config", `{"IngressSlbNetworkType":"internet"}`),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_vswitch_ids": []string{"${local.vswitch_id_1}", "${local.vswitch_id_2}", "${local.vswitch_id_0}"},
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
			{
				Config: testAccConfig(map[string]interface{}{
					"master_instance_types": []string{"${data.alicloud_instance_types.default.instance_types.1.id}", "${data.alicloud_instance_types.default.instance_types.2.id}", "${data.alicloud_instance_types.default.instance_types.0.id}"},
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
			{
				Config: testAccConfig(map[string]interface{}{
					"api_audiences": []string{"https://kubernetes.default.svc.cluster.local", "https://kubernetes.default.svc"},
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
			// RDS, addons and runtime assertions observe retained state, not cloud updates; stock Read does not independently refresh these fields.
			{
				Config: testAccConfig(map[string]interface{}{
					"rds_instances": []string{"${alicloud_db_instance.default.id}", "${alicloud_db_instance.extra.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rds_instances.#": "2",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "rds_instances.0", "alicloud_db_instance.default", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "rds_instances.1", "alicloud_db_instance.extra", "id"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rds_instances": []string{"${alicloud_db_instance.extra.id}", "${alicloud_db_instance.default.id}"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
					resource.TestCheckResourceAttrPair(resourceId, "rds_instances.0", "alicloud_db_instance.extra", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "rds_instances.1", "alicloud_db_instance.default", "id"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retain_resources": []string{"${alicloud_security_group.group.id}", "${alicloud_vpc.vpc.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retain_resources.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retain_resources": []string{"${alicloud_vpc.vpc.0.id}", "${alicloud_security_group.group.id}"},
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
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_options": []map[string]string{
						{
							"resource_type": "SLB",
							"delete_mode":   "delete",
						},
						{
							"resource_type": "ALB",
							"delete_mode":   "delete",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_options.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_options": []map[string]string{
						{
							"resource_type": "ALB",
							"delete_mode":   "delete",
						},
						{
							"resource_type": "SLB",
							"delete_mode":   "delete",
						},
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
			{
				Config: testAccConfig(map[string]interface{}{
					"addons": []map[string]string{
						{
							"name":     "terway-eniip",
							"config":   "",
							"version":  "",
							"disabled": "false",
						},
						{
							"name":     "nginx-ingress-controller",
							"config":   "{\\\"IngressSlbNetworkType\\\":\\\"internet\\\"}",
							"version":  "",
							"disabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addons.#": "2",
					}),
					resource.TestCheckResourceAttr(resourceId, "addons.0.name", "terway-eniip"),
					resource.TestCheckResourceAttr(resourceId, "addons.1.name", "nginx-ingress-controller"),
					resource.TestCheckResourceAttr(resourceId, "addons.1.config", `{"IngressSlbNetworkType":"internet"}`),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addons": []map[string]string{
						{
							"name":     "nginx-ingress-controller",
							"config":   "{\\\"IngressSlbNetworkType\\\":\\\"internet\\\"}",
							"version":  "",
							"disabled": "false",
						},
						{
							"name":     "terway-eniip",
							"config":   "",
							"version":  "",
							"disabled": "false",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
					resource.TestCheckResourceAttr(resourceId, "addons.0.name", "nginx-ingress-controller"),
					resource.TestCheckResourceAttr(resourceId, "addons.0.config", `{"IngressSlbNetworkType":"internet"}`),
					resource.TestCheckResourceAttr(resourceId, "addons.1.name", "terway-eniip"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addons": []map[string]string{
						{
							"name":     "terway-eniip",
							"config":   "",
							"version":  "",
							"disabled": "false",
						},
						{
							"name":     "logtail-ds",
							"config":   "",
							"version":  "1.0.0",
							"disabled": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addons.#": "2",
					}),
					resource.TestCheckResourceAttr(resourceId, "addons.0.config", ""),
					resource.TestCheckResourceAttr(resourceId, "addons.1.name", "logtail-ds"),
					resource.TestCheckResourceAttr(resourceId, "addons.1.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceId, "addons.1.disabled", "true"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": []map[string]interface{}{
						{"name": "docker", "version": "19.03.15"},
						{"name": "containerd", "version": "1.6.28"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
					resource.TestCheckResourceAttr(resourceId, "runtime.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "runtime.0.name", "docker"),
					resource.TestCheckResourceAttr(resourceId, "runtime.0.version", "19.03.15"),
					resource.TestCheckResourceAttr(resourceId, "runtime.1.name", "containerd"),
					resource.TestCheckResourceAttr(resourceId, "runtime.1.version", "1.6.28"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": []map[string]interface{}{
						{"name": "containerd", "version": "1.6.28"},
						{"name": "docker", "version": "19.03.15"},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
					resource.TestCheckResourceAttr(resourceId, "runtime.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "runtime.0.name", "containerd"),
					resource.TestCheckResourceAttr(resourceId, "runtime.0.version", "1.6.28"),
					resource.TestCheckResourceAttr(resourceId, "runtime.1.name", "docker"),
					resource.TestCheckResourceAttr(resourceId, "runtime.1.version", "19.03.15"),
				),
			},
			// pod_vswitch_ids: verify reorder convergence
			{
				Config: testAccConfig(map[string]interface{}{
					"pod_vswitch_ids": []string{"${local.vswitch_id_1}", "${local.vswitch_id_0}"},
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

func TestAccAliCloudCSKubernetes_kms(t *testing.T) {
	var v *aliyungocs.KubernetesClusterDetail

	resourceId := "alicloud_cs_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"kms_encrypted_password": acctest.CHECKSET,
	})

	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKubernetes_kms-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSKubernetesConfigDependenceKms)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                   name,
					"version":                "1.32.7-aliyun.1",
					"master_vswitch_ids":     []string{"${local.vswitch_id}", "${local.vswitch_id}", "${local.vswitch_id}"},
					"master_instance_types":  []string{"${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}"},
					"master_disk_category":   "cloud_essd",
					"master_disk_size":       "80",
					"service_cidr":           "172.18.0.0/16",
					"pod_cidr":               "10.72.0.0/16",
					"new_nat_gateway":        "false",
					"security_group_id":      "${alicloud_security_group.default.id}",
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{"name": "kms-ctx"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

func resourceCSKubernetesConfigDependenceOrder(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_disk_category     = "cloud_essd"
}

data "alicloud_db_zones" "rds" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  db_instance_storage_type = "cloud_essd"
  category                 = "Basic"
}

data "alicloud_db_instance_classes" "rds" {
  zone_id                  = data.alicloud_db_zones.rds.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  db_instance_storage_type = "cloud_essd"
  category                 = "Basic"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 4
  memory_size          = 8
  system_disk_category = "cloud_essd"
}

data "alicloud_images" "default" {
  name_regex = "^aliyun_3_x64_20G_alibase*"
  owners     = "system"
}

resource "alicloud_vpc" "vpc" {
  count      = 1
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitches" {
  count        = 3
  vpc_id       = alicloud_vpc.vpc.0.id
  cidr_block   = format("192.168.%%d.0/24", count.index + 1)
  zone_id      = data.alicloud_zones.default.zones[count.index].id
  vswitch_name = var.name
}

resource "alicloud_vswitch" "rds" {
  count        = 1
  vpc_id       = alicloud_vpc.vpc.0.id
  cidr_block   = "192.168.30.0/24"
  zone_id      = data.alicloud_db_zones.rds.zones.0.id
  vswitch_name = var.name
}

locals {
  vswitch_id_0       = alicloud_vswitch.vswitches.0.id
  vswitch_id_1       = alicloud_vswitch.vswitches.1.id
  vswitch_id_2       = alicloud_vswitch.vswitches.2.id
  rds_instance_class = [for c in data.alicloud_db_instance_classes.rds.instance_classes : c.instance_class if startswith(c.instance_class, "mysql")][0]
}

resource "alicloud_security_group" "group" {
  security_group_name = var.name
  vpc_id              = alicloud_vpc.vpc.0.id
}

resource "alicloud_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "8.0"
  instance_type        = local.rds_instance_class
  instance_storage     = "30"
  instance_charge_type = "Postpaid"
  instance_name        = "tf-testacckubernetes"
  vswitch_id           = alicloud_vswitch.rds.0.id
  monitoring_period    = "60"
}

resource "alicloud_db_instance" "extra" {
  engine               = "MySQL"
  engine_version       = "8.0"
  instance_type        = local.rds_instance_class
  instance_storage     = "30"
  instance_charge_type = "Postpaid"
  instance_name        = "tf-testacckubernetes-extra"
  vswitch_id           = alicloud_vswitch.rds.0.id
  monitoring_period    = "60"
}
`, name)
}

func resourceCSKubernetesConfigDependenceKms(name string) string {
	return fmt.Sprintf(`%s

data "alicloud_instance_types" "default" {
  availability_zone    = local.zone_id
  cpu_core_count       = 4
  memory_size          = 8
  system_disk_category = "cloud_essd"
}

resource "alicloud_kms_key" "default" {
  description            = var.name
  status                 = "Enabled"
  pending_window_in_days = 7
}

resource "alicloud_kms_ciphertext" "default" {
  key_id    = alicloud_kms_key.default.id
  plaintext = "Yourpassword1234"
  encryption_context = {
    "name" = "kms-ctx"
  }
}
`, resourceCSKubernetesConfigDependence(name))
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

variable "instance_type" {
  default = "ecs.c6.xlarge"
}

data "alicloud_regions" "current" {
  current = true
}

data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
  available_disk_category     = "cloud_essd"
  available_instance_type     = var.instance_type
}

locals {
  # The filtered zone set fluctuates between refreshes, and the runner region varies between
  # environments. Both zones below are verified to offer ecs.c6.xlarge, cloud_essd and the
  # required MySQL 8.0 HighAvailability cloud_essd classes (mysql.n2.small), so pin a literal
  # zone per region to keep the fixture stable across test steps.
  zone_by_region = {
    "cn-hangzhou" = "cn-hangzhou-f"
    "cn-beijing"  = "cn-beijing-f"
  }
  zone_id = contains(keys(local.zone_by_region), data.alicloud_regions.current.regions.0.id) ? local.zone_by_region[data.alicloud_regions.current.regions.0.id] : sort(data.alicloud_zones.default.zones[*].id)[0]
  # The cluster and node-pool steps below all place their vswitches in this zone.
  vswitch_id = alicloud_vswitch.vswitches.id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_cs_kubernetes_version" "kubernetes_versions" {
  cluster_type = "Kubernetes"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitches" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = local.zone_id
  vswitch_name = var.name
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = local.zone_id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

locals {
  rds_instance_class = [for c in data.alicloud_db_instance_classes.default.instance_classes : c if startswith(c.instance_class, "mysql.n2.small")].0
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = local.rds_instance_class.instance_class
  instance_storage         = local.rds_instance_class.storage_range.0.min
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "Postpaid"
  instance_name            = var.name
  vswitch_id               = alicloud_vswitch.vswitches.id
  monitoring_period        = "60"
}

resource "alicloud_key_pair" "default" {
  key_pair_name = var.name
  timeouts {
    delete = "15m"
  }
}

resource "alicloud_security_group" "default" {
  security_group_name = var.name
  vpc_id              = alicloud_vpc.vpc.id
}

resource "alicloud_ecs_auto_snapshot_policy" "default" {
  auto_snapshot_policy_name = var.name
  repeat_weekdays           = ["1", "2", "3"]
  retention_days            = -1
  time_points               = ["1", "22", "23"]
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  cluster_id                    = alicloud_cs_kubernetes.default.id
  node_pool_name                = var.name
  vswitch_ids                   = [local.vswitch_id]
  instance_types                = [var.instance_type]
  password                      = "Test12345"
  system_disk_size              = 50
  system_disk_category          = "cloud_essd"
  instance_charge_type          = "PostPaid"
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
data "alicloud_instance_types" "default" {
  cpu_core_count       = 4
  memory_size          = 8
  system_disk_category = "cloud_essd"
  instance_charge_type = "PrePaid"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "test" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "test" {
  vpc_id       = alicloud_vpc.test.id
  cidr_block   = cidrsubnet(alicloud_vpc.test.cidr_block, 8, 8)
  zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  vswitch_name = var.name
}

locals {
  vswitch_id = alicloud_vswitch.test.id
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  cluster_id                    = alicloud_cs_kubernetes.default.id
  node_pool_name                = var.name
  security_group_ids            = [alicloud_cs_kubernetes.default.security_group_id]
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
	"name":              acctest.CHECKSET,
	"security_group_id": acctest.CHECKSET,
	"image_id":          acctest.CHECKSET,
	"version":           acctest.CHECKSET,
	"platform":          acctest.CHECKSET,
	//"certificate_authority.0.cluster_cert": CHECKSET,
	//"certificate_authority.0.client_cert":  CHECKSET,
	//"certificate_authority.0.client_key":   CHECKSET,
	//"connections.0.api_server_internet": CHECKSET,
	"connections.0.api_server_intranet": acctest.CHECKSET,
	//"connections.0.master_public_ip":    CHECKSET,
	"connections.0.service_domain": acctest.CHECKSET,
	"worker_ram_role_name":         acctest.CHECKSET,
	"node_name_mode":               acctest.CHECKSET,
	"master_nodes.#":               "3",
	"vpc_id":                       acctest.CHECKSET,
	"resource_group_id":            acctest.CHECKSET,
	//"slb_internet":                      CHECKSET,
	"slb_intranet": acctest.CHECKSET,
	"slb_id":       acctest.CHECKSET,
}

func TestAccAliCloudCSKubernetes_StateMigrationV0ToV1(t *testing.T) {
	if os.Getenv("ALICLOUD_STATE_MIGRATION_V0_V1") == "" {
		t.Skip("ALICLOUD_STATE_MIGRATION_V0_V1 not set. Test for V0 -> V1 state migration")
	}
	var v map[string]interface{}
	resourceId := "alicloud_cs_kubernetes.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})
	serviceFunc := func() interface{} {
		return cs.NewCsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-statemig-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"alicloud": {
						Source:            "aliyun/alicloud",
						VersionConstraint: "1.282.0",
					},
				},
				Config: testAccCSKubernetesStateMigrationConfigV0(name),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
				Config:                   testAccCSKubernetesStateMigrationConfigV1(name),
			},
		},
	})
}

func testAccCSKubernetesStateMigrationDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_disk_category     = "cloud_essd"
  available_instance_type     = "ecs.c6.xlarge"
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_key_pair" "default" {
  key_pair_name = var.name
}
`, name)
}

func testAccCSKubernetesStateMigrationConfigV0(name string) string {
	return testAccCSKubernetesStateMigrationDependence(name) + `
resource "alicloud_cs_kubernetes" "default" {
  name                           = var.name
  master_vswitch_ids             = [alicloud_vswitch.default.id, alicloud_vswitch.default.id, alicloud_vswitch.default.id]
  master_instance_types          = ["ecs.c6.xlarge", "ecs.c6.xlarge", "ecs.c6.xlarge"]
  master_disk_category           = "cloud_essd"
  master_disk_size               = 40
  key_name                       = alicloud_key_pair.default.key_pair_name
  pod_cidr                       = "10.75.0.0/16"
  service_cidr                   = "172.25.0.0/16"
  new_nat_gateway                = true
  slb_internet_enabled           = true
  skip_set_certificate_authority = false
  proxy_mode                     = "ipvs"
  runtime = {
    name    = "containerd"
    version = "1.6.20"
  }
}

output "runtime" {
  value = alicloud_cs_kubernetes.default.runtime["name"]
}

output "alicloud_cs_kubernetes" {
  value = alicloud_cs_kubernetes.default.certificate_authority["cluster_cert"]
}

output "connections" {
  value = alicloud_cs_kubernetes.default.connections["api_server_internet"]
}
`
}

func testAccCSKubernetesStateMigrationConfigV1(name string) string {
	return testAccCSKubernetesStateMigrationDependence(name) + `
resource "alicloud_cs_kubernetes" "default" {
  name                           = var.name
  master_vswitch_ids             = [alicloud_vswitch.default.id, alicloud_vswitch.default.id, alicloud_vswitch.default.id]
  master_instance_types          = ["ecs.c6.xlarge", "ecs.c6.xlarge", "ecs.c6.xlarge"]
  master_disk_category           = "cloud_essd"
  master_disk_size               = 40
  key_name                       = alicloud_key_pair.default.key_pair_name
  pod_cidr                       = "10.75.0.0/16"
  service_cidr                   = "172.25.0.0/16"
  new_nat_gateway                = true
  slb_internet_enabled           = true
  skip_set_certificate_authority = false
  proxy_mode                     = "ipvs"
  runtime {
    name    = "containerd"
    version = "1.6.20"
  }
}

output "runtime" {
  value = alicloud_cs_kubernetes.default.runtime.0.name
}

output "alicloud_cs_kubernetes" {
  value = alicloud_cs_kubernetes.default.certificate_authority.0.cluster_cert
}

output "connections" {
  value = alicloud_cs_kubernetes.default.connections.0.api_server_internet
}
`
}
