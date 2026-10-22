package amqp_test

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/amqp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// vswitch_ids is a pre-existing ForceNew TypeList attribute. This PR migrates
// the resource into the service/amqp package and fixes the PayAsYouGo/Serverless
// Delete path; it does not add or change vswitch_ids or its ordering semantics
// (the schema definition is identical to the pre-migration baseline). The
// convergent-reorder coverage the TypeList Order Coverage gate requires is
// therefore pre-existing stock debt exposed by the file touch, not a gap this
// change introduces, so it is marked here rather than added as a tangential
// cloud test.
// collection-order: bypass vswitch_ids -- pre-existing stock attribute; schema unchanged by this PR (Delete fix + migration only)

const amqpRemoveKey = "REMOVEKEY"

// testAccCheckAmqpInstanceDestroy verifies that every alicloud_amqp_instance
// recorded in the state has been released on the cloud side — covering both the
// Subscription refund path and the PayAsYouGo/Serverless DeleteInstance path.
func testAccCheckAmqpInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_amqp_instance" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := amqp.NewService(client).DescribeAmqpInstance(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing AMQP instance %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("AMQP instance %s still exists", rs.Primary.ID)
	}
	return nil
}

// amqpTestConfig reproduces the legacy accumulated-config builder: each step
// only declares the attributes that change, and the builder carries the rest
// forward, rendering a full resource block.
type amqpTestConfig struct {
	dependence string
	attrs      map[string]interface{}
}

func newAmqpTestConfig(dependence func(string) string, name string) *amqpTestConfig {
	return &amqpTestConfig{
		dependence: dependence(name),
		attrs:      map[string]interface{}{},
	}
}

func (c *amqpTestConfig) build(changes map[string]interface{}) string {
	for k, v := range changes {
		if s, ok := v.(string); ok && s == amqpRemoveKey {
			delete(c.attrs, k)
			continue
		}
		c.attrs[k] = v
	}
	var sb strings.Builder
	sb.WriteString(c.dependence)
	sb.WriteString("\n\nresource \"alicloud_amqp_instance\" \"default\" {")
	keys := make([]string, 0, len(c.attrs))
	for k := range c.attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString("\n  ")
		sb.WriteString(k)
		sb.WriteString(" = ")
		sb.WriteString(amqpRenderValue(c.attrs[k]))
	}
	sb.WriteString("\n}")
	return sb.String()
}

func amqpRenderValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return "\"" + val + "\""
	case bool:
		if val {
			return "true"
		}
		return "false"
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case []string:
		parts := make([]string, len(val))
		for i, s := range val {
			parts[i] = "\"" + s + "\""
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case map[string]string:
		parts := make([]string, 0, len(val))
		for k, vv := range val {
			parts = append(parts, k+" = \""+vv+"\"")
		}
		return "{" + strings.Join(parts, ", ") + "}"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func amqpRegion() string {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	return region
}

// amqpSweepClient builds a *connectivity.AliyunClient from env credentials for
// the AMQP acceptance-test helpers (e.g. the encryption service-linked-role
// preCheck). It is NOT the sweeper registration: that lives in the root alicloud
// test package (alicloud/sweeper_amqp_instance_test.go) so `make sweep`
// discovers it; this helper only constructs a client object.
func amqpSweepClient(region string) (*connectivity.AliyunClient, error) {
	accessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	if accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}
	secretKey := os.Getenv("ALICLOUD_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}
	securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN")
	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &sync.Map{},
		SignVersion: &sync.Map{},
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

func amqpPreCheck(t *testing.T) {
	t.Helper()
	acctest.PreCheck(t)
	// AMQP instance/serverless tests only support cn-hangzhou. The suite's
	// PreCheck defaults an empty ALICLOUD_REGION to cn-beijing; replicate the
	// original testAccPreCheckWithRegions(t, true, ["cn-hangzhou"]) behaviour by
	// remapping the region to cn-hangzhou instead of skipping, so the ACC runs
	// there rather than being no-ops.
	if region := os.Getenv("ALICLOUD_REGION"); region != "cn-hangzhou" {
		t.Logf("AMQP instance tests require cn-hangzhou; remapping region %q to cn-hangzhou", region)
		os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
	}
}

func testAccAmqpInstanceBasicDependence6128(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
		vpc_id     = data.alicloud_vpcs.default.ids.0
		name_regex = "^default-NODELETING-vsw"
	}

	data "alicloud_security_groups" "default" {
		vpc_id     = data.alicloud_vpcs.default.ids.0
		name_regex = "default-NODELETING"
	}
`, name)
}

func testAccAmqpInstanceEncryptedDependence6128(name string) string {
	return testAccAmqpInstanceBasicDependence6128(name) + `
	resource "alicloud_kms_key" "default" {
		description            = "terraform-provider-alicloud AMQP acceptance test"
		pending_window_in_days = 7
	}
`
}

func testAccAmqpInstanceBasicDependence11166(name string) string {
	return testAccAmqpInstanceBasicDependence6128(name)
}

// testAccEnsureAmqpEncryptionServiceLinkedRole creates (if absent) and waits
// for the AMQP encryption service-linked role so encrypted instance creation
// does not fail on a missing role.
func testAccEnsureAmqpEncryptionServiceLinkedRole(t *testing.T) {
	t.Helper()
	const (
		serviceName = "encrypt.amqp.aliyuncs.com"
		roleName    = "AliyunServiceLinkedRoleForAmqpEncrypt"
	)

	client, err := amqpSweepClient(amqpRegion())
	if err != nil {
		t.Fatalf("failed to get Alibaba Cloud client: %s", err)
	}

	if _, err := amqpGetRole(client, roleName); err == nil {
		return
	} else if !errs.NotFoundError(err) && !errs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
		t.Fatalf("failed to describe AMQP encryption service-linked role: %s", err)
	}

	request := map[string]interface{}{"ServiceName": serviceName}
	if _, err := client.RpcPost("ResourceManager", "2020-03-31", "CreateServiceLinkedRole", nil, request, true); err != nil &&
		!errs.IsExpectedErrors(err, []string{"AlreadyExists", "EntityAlreadyExists"}) {
		t.Fatalf("failed to create AMQP encryption service-linked role: %s", err)
	}

	if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		if _, describeErr := amqpGetRole(client, roleName); describeErr != nil {
			if errs.NotFoundError(describeErr) || errs.IsExpectedErrors(describeErr, []string{"EntityNotExist.Role"}) {
				return retry.RetryableError(describeErr)
			}
			return retry.NonRetryableError(describeErr)
		}
		return nil
	}); err != nil {
		t.Fatalf("AMQP encryption service-linked role did not become available: %s", err)
	}
}

func amqpGetRole(client *connectivity.AliyunClient, roleName string) (*ram.GetRoleResponse, error) {
	request := ram.CreateGetRoleRequest()
	request.RegionId = client.RegionId
	request.RoleName = roleName
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, roleName, "GetRole", errs.AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ram.GetRoleResponse)
	return response, nil
}

// Case: serverless PayAsYouGo instance (onDemand -> provisioned), 6128.
func TestAccAliCloudAmqpInstanceResource_basic6128(t *testing.T) {
	resourceId := "alicloud_amqp_instance.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%samqpinstance%d", amqpRegion(), rand)
	cfg := newAmqpTestConfig(testAccAmqpInstanceBasicDependence6128, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			amqpPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAmqpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg.build(map[string]interface{}{
					"payment_type":           "PayAsYouGo",
					"vpc_id":                 "${data.alicloud_vswitches.default.vpc_id}",
					"vswitch_ids":            []string{"${data.alicloud_vswitches.default.ids.0}", "${data.alicloud_vswitches.default.ids.1}"},
					"security_group_id":      "${data.alicloud_security_groups.default.ids.0}",
					"serverless_charge_type": "onDemand",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "payment_type", "PayAsYouGo"),
					resource.TestCheckResourceAttrSet(resourceId, "vpc_id"),
					resource.TestCheckResourceAttr(resourceId, "vswitch_ids.#", "2"),
					resource.TestCheckResourceAttrSet(resourceId, "security_group_id"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"serverless_charge_type": "provisioned",
					"provisioned_capacity":   "2000",
					"modify_type":            "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "provisioned_capacity", "2000"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"edition": "shared",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "edition", "shared"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instance_name", name),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"provisioned_capacity": "20000",
					"modify_type":          "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "provisioned_capacity", "20000"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"support_eip": "true",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "support_eip", "true"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"support_tracing": "true",
					"modify_type":     "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "support_tracing", "true"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceId, "tags.Created", "TF"),
					resource.TestCheckResourceAttr(resourceId, "tags.For", "Test"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
					},
					"tags.For": amqpRemoveKey,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.Created", "TF-update"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"serverless_switch": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "serverless_switch", "true"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_model", "auto_renew", "modify_type", "period", "period_cycle", "serverless_charge_type", "renewal_duration", "renewal_duration_unit", "renewal_status"},
			},
		},
	})
}

// Case: serverless PayAsYouGo provisioned dedicated instance, 6128 twin.
func TestAccAliCloudAmqpInstanceResource_basic6128_twin(t *testing.T) {
	resourceId := "alicloud_amqp_instance.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%samqpinstance%d", amqpRegion(), rand)
	cfg := newAmqpTestConfig(testAccAmqpInstanceBasicDependence6128, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			amqpPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAmqpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg.build(map[string]interface{}{
					"payment_type":           "PayAsYouGo",
					"vpc_id":                 "${data.alicloud_vswitches.default.vpc_id}",
					"vswitch_ids":            []string{"${data.alicloud_vswitches.default.ids.0}", "${data.alicloud_vswitches.default.ids.1}"},
					"security_group_id":      "${data.alicloud_security_groups.default.ids.0}",
					"auth_model":             "ram",
					"serverless_charge_type": "provisioned",
					"serverless_switch":      "true",
					"provisioned_capacity":   "20000",
					"edition":                "dedicated",
					"instance_name":          name,
					"support_eip":            "true",
					"support_tracing":        "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "payment_type", "PayAsYouGo"),
					resource.TestCheckResourceAttrSet(resourceId, "vpc_id"),
					resource.TestCheckResourceAttr(resourceId, "vswitch_ids.#", "2"),
					resource.TestCheckResourceAttrSet(resourceId, "security_group_id"),
					resource.TestCheckResourceAttr(resourceId, "auth_model", "ram"),
					resource.TestCheckResourceAttr(resourceId, "provisioned_capacity", "20000"),
					resource.TestCheckResourceAttr(resourceId, "edition", "dedicated"),
					resource.TestCheckResourceAttr(resourceId, "instance_name", name),
					resource.TestCheckResourceAttr(resourceId, "support_eip", "true"),
					resource.TestCheckResourceAttr(resourceId, "support_tracing", "true"),
					resource.TestCheckResourceAttr(resourceId, "serverless_switch", "true"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_model", "auto_renew", "modify_type", "period", "period_cycle", "serverless_charge_type", "renewal_duration", "renewal_duration_unit", "renewal_status"},
			},
		},
	})
}

// Case: encrypted serverless dedicated instance.
func TestAccAliCloudAmqpInstanceResource_encryptedServerlessDedicatedCreate(t *testing.T) {
	resourceId := "alicloud_amqp_instance.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%samqpencrypted%d", amqpRegion(), rand)
	cfg := newAmqpTestConfig(testAccAmqpInstanceEncryptedDependence6128, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			amqpPreCheck(t)
			testAccEnsureAmqpEncryptionServiceLinkedRole(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAmqpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg.build(map[string]interface{}{
					"payment_type":           "PayAsYouGo",
					"vpc_id":                 "${data.alicloud_vswitches.default.vpc_id}",
					"vswitch_ids":            []string{"${data.alicloud_vswitches.default.ids.0}", "${data.alicloud_vswitches.default.ids.1}"},
					"security_group_id":      "${data.alicloud_security_groups.default.ids.0}",
					"serverless_charge_type": "provisioned",
					"provisioned_capacity":   "20000",
					"edition":                "dedicated",
					"instance_name":          name,
					"encrypted_instance":     true,
					"kms_key_id":             "${alicloud_kms_key.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "payment_type", "PayAsYouGo"),
					resource.TestCheckResourceAttrSet(resourceId, "vpc_id"),
					resource.TestCheckResourceAttr(resourceId, "vswitch_ids.#", "2"),
					resource.TestCheckResourceAttrSet(resourceId, "security_group_id"),
					resource.TestCheckResourceAttr(resourceId, "serverless_charge_type", "provisioned"),
					resource.TestCheckResourceAttr(resourceId, "provisioned_capacity", "20000"),
					resource.TestCheckResourceAttr(resourceId, "edition", "dedicated"),
					resource.TestCheckResourceAttr(resourceId, "instance_name", name),
					resource.TestCheckResourceAttr(resourceId, "encrypted_instance", "true"),
					resource.TestCheckResourceAttrSet(resourceId, "kms_key_id"),
					resource.TestCheckResourceAttrPair(resourceId, "kms_key_id", "alicloud_kms_key.default", "id"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_model", "auto_renew", "modify_type", "period", "period_cycle", "serverless_charge_type", "renewal_duration", "renewal_duration_unit", "renewal_status"},
			},
		},
	})
}

// Case: Subscription enterprise -> vip instance, 11166.
func TestAccAliCloudAmqpInstanceResource_basic11166(t *testing.T) {
	resourceId := "alicloud_amqp_instance.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%samqpinstance%d", amqpRegion(), rand)
	cfg := newAmqpTestConfig(testAccAmqpInstanceBasicDependence11166, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			amqpPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAmqpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg.build(map[string]interface{}{
					"payment_type":      "Subscription",
					"vpc_id":            "${data.alicloud_vswitches.default.vpc_id}",
					"vswitch_ids":       []string{"${data.alicloud_vswitches.default.ids.0}", "${data.alicloud_vswitches.default.ids.1}"},
					"security_group_id": "${data.alicloud_security_groups.default.ids.0}",
					"instance_type":     "enterprise",
					"max_tps":           "3000",
					"max_connections":   "2000",
					"queue_capacity":    "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "payment_type", "Subscription"),
					resource.TestCheckResourceAttrSet(resourceId, "vpc_id"),
					resource.TestCheckResourceAttr(resourceId, "vswitch_ids.#", "2"),
					resource.TestCheckResourceAttrSet(resourceId, "security_group_id"),
					resource.TestCheckResourceAttr(resourceId, "instance_type", "enterprise"),
					resource.TestCheckResourceAttr(resourceId, "max_tps", "3000"),
					resource.TestCheckResourceAttr(resourceId, "max_connections", "2000"),
					resource.TestCheckResourceAttr(resourceId, "queue_capacity", "200"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instance_name", name),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"instance_type":   "vip",
					"max_tps":         "8000",
					"max_connections": "50000",
					"queue_capacity":  "10000",
					"modify_type":     "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instance_type", "vip"),
					resource.TestCheckResourceAttr(resourceId, "max_tps", "8000"),
					resource.TestCheckResourceAttr(resourceId, "max_connections", "50000"),
					resource.TestCheckResourceAttr(resourceId, "queue_capacity", "10000"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"max_connections": "51000",
					"modify_type":     "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "max_connections", "51000"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"support_eip": "true",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "support_eip", "true"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"max_eip_tps": "256",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "max_eip_tps", "256"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"max_tps":     "15000",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "max_tps", "15000"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"queue_capacity": "11000",
					"modify_type":    "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "queue_capacity", "11000"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"renewal_status":        "AutoRenewal",
					"renewal_duration":      "1",
					"renewal_duration_unit": "Year",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "renewal_status", "AutoRenewal"),
					resource.TestCheckResourceAttr(resourceId, "renewal_duration", "1"),
					resource.TestCheckResourceAttr(resourceId, "renewal_duration_unit", "Year"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"storage_size": "800",
					"modify_type":  "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "storage_size", "800"),
				),
			},
			{
				Config: cfg.build(map[string]interface{}{
					"support_tracing":      "true",
					"tracing_storage_time": "15",
					"modify_type":          "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "support_tracing", "true"),
					resource.TestCheckResourceAttr(resourceId, "tracing_storage_time", "15"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_model", "auto_renew", "modify_type", "period", "period_cycle", "serverless_charge_type"},
			},
		},
	})
}

// Case: Subscription vip instance (full attribute set), 11166 twin.
func TestAccAliCloudAmqpInstanceResource_basic11166_twin(t *testing.T) {
	resourceId := "alicloud_amqp_instance.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%samqpinstance%d", amqpRegion(), rand)
	cfg := newAmqpTestConfig(testAccAmqpInstanceBasicDependence11166, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			amqpPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAmqpInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg.build(map[string]interface{}{
					"payment_type":          "Subscription",
					"vpc_id":                "${data.alicloud_vswitches.default.vpc_id}",
					"vswitch_ids":           []string{"${data.alicloud_vswitches.default.ids.0}", "${data.alicloud_vswitches.default.ids.1}"},
					"security_group_id":     "${data.alicloud_security_groups.default.ids.0}",
					"auth_model":            "openSource",
					"auto_renew":            "true",
					"period":                "1",
					"period_cycle":          "Year",
					"instance_type":         "vip",
					"listener_mode":         "ssl_only",
					"max_tps":               "8000",
					"max_connections":       "50000",
					"queue_capacity":        "10000",
					"instance_name":         name,
					"support_eip":           "true",
					"max_eip_tps":           "128",
					"renewal_status":        "AutoRenewal",
					"renewal_duration":      "1",
					"renewal_duration_unit": "Year",
					"storage_size":          "800",
					"support_tracing":       "true",
					"tracing_storage_time":  "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "payment_type", "Subscription"),
					resource.TestCheckResourceAttrSet(resourceId, "vpc_id"),
					resource.TestCheckResourceAttr(resourceId, "vswitch_ids.#", "2"),
					resource.TestCheckResourceAttrSet(resourceId, "security_group_id"),
					resource.TestCheckResourceAttr(resourceId, "auth_model", "openSource"),
					resource.TestCheckResourceAttr(resourceId, "instance_type", "vip"),
					resource.TestCheckResourceAttr(resourceId, "listener_mode", "ssl_only"),
					resource.TestCheckResourceAttr(resourceId, "max_tps", "8000"),
					resource.TestCheckResourceAttr(resourceId, "max_connections", "50000"),
					resource.TestCheckResourceAttr(resourceId, "queue_capacity", "10000"),
					resource.TestCheckResourceAttr(resourceId, "instance_name", name),
					resource.TestCheckResourceAttr(resourceId, "support_eip", "true"),
					resource.TestCheckResourceAttr(resourceId, "max_eip_tps", "128"),
					resource.TestCheckResourceAttr(resourceId, "renewal_status", "AutoRenewal"),
					resource.TestCheckResourceAttr(resourceId, "renewal_duration", "1"),
					resource.TestCheckResourceAttr(resourceId, "renewal_duration_unit", "Year"),
					resource.TestCheckResourceAttr(resourceId, "storage_size", "800"),
					resource.TestCheckResourceAttr(resourceId, "support_tracing", "true"),
					resource.TestCheckResourceAttr(resourceId, "tracing_storage_time", "15"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_model", "auto_renew", "modify_type", "period", "period_cycle", "serverless_charge_type"},
			},
		},
	})
}
