package slb_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	slbsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/slb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudSLBMasterSlaveServerGroup_vpc(t *testing.T) {
	var v *slbsdk.DescribeMasterSlaveServerGroupAttributeResponse
	resourceId := "alicloud_slb_master_slave_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, testAccSlbMasterSlaveServerGroupCheckMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return slb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()

	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbMasterSlaveServerGroupVpc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceMasterSlaveServerGroupConfigDependence)
	initialConfig := testAccConfig(map[string]interface{}{
		"load_balancer_id":             "${alicloud_slb_load_balancer.default.id}",
		"delete_protection_validation": false,
		"name":                         "${var.name}",
		"servers":                      masterSlaveServerGroupAccServers("tf-testAccMasterDesc", "tf-testAccSlaveDesc", true),
	})
	changedConfig := testAccConfig(map[string]interface{}{
		"servers": masterSlaveServerGroupAccServers("tf-testAccMasterDescUpdated", "tf-testAccSlaveDescUpdated", true),
	})
	removedConfig := masterSlaveServerGroupConfigWithoutDescriptions(testAccConfig)
	var initialID, changedID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      name,
						"servers.#": "2",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "servers.*", map[string]string{
						"server_type": "Master",
						"description": "tf-testAccMasterDesc",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "servers.*", map[string]string{
						"server_type": "Slave",
						"description": "tf-testAccSlaveDesc",
					}),
					captureMasterSlaveServerGroupID(resourceId, &initialID),
					checkMasterSlaveServerGroupDescriptionRead(resourceId, map[string]string{
						"Master": "tf-testAccMasterDesc",
						"Slave":  "tf-testAccSlaveDesc",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionDestroyBeforeCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionNoop),
					},
				},
				Config: changedConfig,
				Check: resource.ComposeTestCheckFunc(
					checkMasterSlaveServerGroupReplaced(resourceId, &initialID),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "servers.*", map[string]string{
						"server_type": "Master",
						"description": "tf-testAccMasterDescUpdated",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "servers.*", map[string]string{
						"server_type": "Slave",
						"description": "tf-testAccSlaveDescUpdated",
					}),
					checkMasterSlaveServerGroupDescriptionRead(resourceId, map[string]string{
						"Master": "tf-testAccMasterDescUpdated",
						"Slave":  "tf-testAccSlaveDescUpdated",
					}),
					captureMasterSlaveServerGroupID(resourceId, &changedID),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionDestroyBeforeCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionNoop),
					},
				},
				Config: removedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"servers.#": "2"}),
					checkMasterSlaveServerGroupReplaced(resourceId, &changedID),
					checkMasterSlaveServerGroupDescriptionRead(resourceId, map[string]string{
						"Master": "",
						"Slave":  "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

func TestAccAliCloudSLBMasterSlaveServerGroup_multi_vpc(t *testing.T) {
	var v *slbsdk.DescribeMasterSlaveServerGroupAttributeResponse
	resourceId := "alicloud_slb_master_slave_server_group.default.1"
	planResourceId := "alicloud_slb_master_slave_server_group.default[1]"
	ra := acctest.ResourceAttrInit(resourceId, testAccSlbMasterSlaveServerGroupCheckMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return slb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()

	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbMasterSlaveServerGroupVpc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceMasterSlaveServerGroupConfigDependence)
	initialConfig := testAccConfig(map[string]interface{}{
		"count":                        "10",
		"load_balancer_id":             "${alicloud_slb_load_balancer.default.id}",
		"delete_protection_validation": false,
		"name":                         "${var.name}",
		"servers":                      masterSlaveServerGroupAccServers("tf-testAccMasterDesc", "tf-testAccSlaveDesc", true),
	})
	changedConfig := testAccConfig(map[string]interface{}{
		"servers": masterSlaveServerGroupAccServers("tf-testAccMasterDescUpdated", "tf-testAccSlaveDescUpdated", true),
	})
	removedConfig := masterSlaveServerGroupConfigWithoutDescriptions(testAccConfig)
	var initialID, changedID string
	importStep := func() resource.TestStep {
		return resource.TestStep{
			ResourceName:            "alicloud_slb_master_slave_server_group.default[1]",
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			ImportStateIdFunc: func(s *terraform.State) (string, error) {
				rs, ok := s.RootModule().Resources["alicloud_slb_master_slave_server_group.default.1"]
				if !ok {
					return "", fmt.Errorf("master slave server group default.1 not found in state")
				}
				return rs.Primary.ID, nil
			},
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      name,
						"servers.#": "2",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "servers.*", map[string]string{
						"server_type": "Master",
						"description": "tf-testAccMasterDesc",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "servers.*", map[string]string{
						"server_type": "Slave",
						"description": "tf-testAccSlaveDesc",
					}),
					captureMasterSlaveServerGroupID(resourceId, &initialID),
					checkMasterSlaveServerGroupDescriptionRead(resourceId, map[string]string{
						"Master": "tf-testAccMasterDesc",
						"Slave":  "tf-testAccSlaveDesc",
					}),
				),
			},
			importStep(),
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(planResourceId, plancheck.ResourceActionDestroyBeforeCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(planResourceId, plancheck.ResourceActionNoop),
					},
				},
				Config: changedConfig,
				Check: resource.ComposeTestCheckFunc(
					checkMasterSlaveServerGroupReplaced(resourceId, &initialID),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "servers.*", map[string]string{
						"server_type": "Master",
						"description": "tf-testAccMasterDescUpdated",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "servers.*", map[string]string{
						"server_type": "Slave",
						"description": "tf-testAccSlaveDescUpdated",
					}),
					checkMasterSlaveServerGroupDescriptionRead(resourceId, map[string]string{
						"Master": "tf-testAccMasterDescUpdated",
						"Slave":  "tf-testAccSlaveDescUpdated",
					}),
					captureMasterSlaveServerGroupID(resourceId, &changedID),
				),
			},
			importStep(),
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(planResourceId, plancheck.ResourceActionDestroyBeforeCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(planResourceId, plancheck.ResourceActionNoop),
					},
				},
				Config: removedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"servers.#": "2"}),
					checkMasterSlaveServerGroupReplaced(resourceId, &changedID),
					checkMasterSlaveServerGroupDescriptionRead(resourceId, map[string]string{
						"Master": "",
						"Slave":  "",
					}),
				),
			},
			importStep(),
		},
	})
}

func resourceMasterSlaveServerGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_zones" "default" {
    available_disk_category = "cloud_essd"
    available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
	system_disk_category = "cloud_essd"
	eni_amount        = 2
}
data "alicloud_images" "default" {
    instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
    name_regex = "^ubuntu"
    most_recent = true
    owners = "system"
}
resource "alicloud_vpc" "default" {
    vpc_name = "${var.name}"
    cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
    vpc_id = "${alicloud_vpc.default.id}"
    cidr_block = "172.16.0.0/16"
    availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
    name = "${var.name}"
}
resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}
resource "alicloud_ecs_network_interface" "default" {
    count = 1
    name = "${var.name}"
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
}

resource "alicloud_ecs_network_interface_attachment" "default" {
    count = 1
    instance_id = "${alicloud_instance.default.0.id}"
    network_interface_id = "${element(alicloud_ecs_network_interface.default.*.id, count.index)}"
}
resource "alicloud_instance" "default" {
    image_id = "${data.alicloud_images.default.images.0.id}"
    instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
    instance_name = "${var.name}"
    count = "2"
    security_groups = "${alicloud_security_group.default.*.id}"
    internet_charge_type = "PayByTraffic"
    internet_max_bandwidth_out = "10"
    availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
    instance_charge_type = "PostPaid"
    system_disk_category = "cloud_essd"
    vswitch_id = "${alicloud_vswitch.default.id}"
}
resource "alicloud_slb_load_balancer" "default" {
    load_balancer_name = "${var.name}"
    vswitch_id = "${alicloud_vswitch.default.id}"
    load_balancer_spec  = "slb.s2.small"
}
`, name)
}

func masterSlaveServerGroupAccServers(masterDescription, slaveDescription string, includeDescriptions bool) []map[string]interface{} {
	servers := []map[string]interface{}{
		{
			"server_id":   "${alicloud_instance.default.0.id}",
			"port":        "100",
			"weight":      "100",
			"type":        "ecs",
			"server_type": "Master",
		},
		{
			"server_id":   "${alicloud_instance.default.1.id}",
			"port":        "100",
			"weight":      "100",
			"type":        "ecs",
			"server_type": "Slave",
		},
	}
	if includeDescriptions {
		servers[0]["description"] = masterDescription
		servers[1]["description"] = slaveDescription
	}
	return servers
}

func masterSlaveServerGroupConfigWithoutDescriptions(config acctest.ResourceTestAccConfigFunc) string {
	config(map[string]interface{}{"servers": acctest.REMOVEKEY})
	return config(map[string]interface{}{
		"servers": masterSlaveServerGroupAccServers("", "", false),
	})
}

func captureMasterSlaveServerGroupID(resourceId string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok || rs.Primary == nil || rs.Primary.ID == "" {
			return fmt.Errorf("master slave server group %s has no ID", resourceId)
		}
		*id = rs.Primary.ID
		return nil
	}
}

func checkMasterSlaveServerGroupReplaced(resourceId string, previousID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if previousID == nil || *previousID == "" {
			return fmt.Errorf("master slave server group %s has no previous ID", resourceId)
		}
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok || rs.Primary == nil || rs.Primary.ID == "" {
			return fmt.Errorf("master slave server group %s has no current ID", resourceId)
		}
		if rs.Primary.ID == *previousID {
			return fmt.Errorf("master slave server group %s was not replaced: ID remains %s", resourceId, *previousID)
		}
		return nil
	}
}

func checkMasterSlaveServerGroupDescriptionRead(resourceId string, expected map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok || rs.Primary == nil || rs.Primary.ID == "" {
			return fmt.Errorf("master slave server group %s has no ID for Describe", resourceId)
		}
		service := slb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
		response, err := service.DescribeSlbMasterSlaveServerGroup(rs.Primary.ID)
		if err != nil {
			return err
		}
		if response == nil {
			return fmt.Errorf("master slave server group %s returned an empty Describe response", rs.Primary.ID)
		}
		remaining := make(map[string]string, len(expected))
		for role, description := range expected {
			remaining[role] = description
		}
		for _, server := range response.MasterSlaveBackendServers.MasterSlaveBackendServer {
			want, ok := remaining[server.ServerType]
			if !ok {
				return fmt.Errorf("Describe returned unexpected or duplicate server role %q", server.ServerType)
			}
			if server.Description != want {
				return fmt.Errorf("Describe server %s role %s returned description %q, want %q", server.ServerId, server.ServerType, server.Description, want)
			}
			delete(remaining, server.ServerType)
		}
		if len(remaining) != 0 {
			return fmt.Errorf("Describe response is missing server roles: %v", remaining)
		}
		return nil
	}
}

var testAccSlbMasterSlaveServerGroupCheckMap = map[string]string{
	"name":      CHECKSET,
	"servers.#": "2",
}

func TestUnitAlicloudSLBMasterSlaveServerGroupDescriptionConfigRemoval(t *testing.T) {
	config := acctest.ResourceTestAccConfig("alicloud_slb_master_slave_server_group.default", "test", func(string) string {
		return ""
	})
	withDescriptions := config(map[string]interface{}{
		"servers": masterSlaveServerGroupAccServers("master-a", "slave-a", true),
	})
	if !strings.Contains(withDescriptions, "description") {
		t.Fatalf("initial configuration does not contain descriptions: %s", withDescriptions)
	}
	withoutDescriptions := masterSlaveServerGroupConfigWithoutDescriptions(config)
	if strings.Contains(withoutDescriptions, "description") {
		t.Fatalf("removed configuration still contains a description: %s", withoutDescriptions)
	}
}

func TestUnitAlicloudSLBMasterSlaveServerGroupDescriptionValidation(t *testing.T) {
	cases := []struct {
		name        string
		description string
		set         bool
		wantErr     bool
		wantReason  string
	}{
		{name: "unset", set: false, wantErr: false},
		{name: "explicit empty", description: "", set: true, wantErr: true, wantReason: "in the range (1 - 80)"},
		{name: "blank", description: "   ", set: true, wantErr: true, wantReason: "must contain only Chinese characters"},
		{name: "leading and trailing spaces", description: " master-server ", set: true, wantErr: true, wantReason: "must contain only Chinese characters"},
		{name: "valid ascii", description: "master-server_1/2.3", set: true, wantErr: false},
		{name: "valid chinese", description: "主服务器-1", set: true, wantErr: false},
		{name: "80 characters", description: strings.Repeat("a", 80), set: true, wantErr: false},
		{name: "81 characters", description: strings.Repeat("a", 81), set: true, wantErr: true, wantReason: "in the range (1 - 80)"},
		{name: "illegal character", description: "master#server", set: true, wantErr: true, wantReason: "must contain only Chinese characters"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			servers := masterSlaveServerGroupUnitTestServers()
			if c.set {
				servers[0].(map[string]interface{})["description"] = c.description
			}
			diags := slb.ResourceMasterSlaveServerGroup().Validate(sdkterraform.NewResourceConfigRaw(map[string]interface{}{
				"load_balancer_id": "lb-test",
				"name":             "group-test",
				"servers":          servers,
			}))
			var messages []string
			for _, d := range diags {
				messages = append(messages, d.Summary+" "+d.Detail)
			}
			report := strings.Join(messages, " | ")

			if gotErr := diags.HasError(); gotErr != c.wantErr {
				t.Fatalf("description %q: want validation error %v, got %v: %s", c.description, c.wantErr, gotErr, report)
			}
			if !c.wantErr {
				if strings.Contains(report, "description") {
					t.Fatalf("unexpected description diagnostic: %s", report)
				}
				return
			}
			if !strings.Contains(report, "description") {
				t.Fatalf("validation error does not point at the description attribute: %s", report)
			}
			if c.wantReason != "" && !strings.Contains(report, c.wantReason) {
				t.Fatalf("validation error %s does not contain the expected reason %q", report, c.wantReason)
			}
		})
	}
}

func TestUnitAlicloudSLBMasterSlaveServerGroupDescriptionRequest(t *testing.T) {
	cases := []struct {
		name        string
		description string
		set         bool
		wantSubstr  string
		wantAbsent  string
	}{
		{
			name:       "unset description is omitted from the request",
			set:        false,
			wantAbsent: "Description",
		},
		{
			name:        "valid description is sent verbatim",
			description: "master-server_1/2.3",
			set:         true,
			wantSubstr:  "'Description':'master-server_1/2.3'",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, payload := newMasterSlaveServerGroupUnitTestClient(t)

			servers := masterSlaveServerGroupUnitTestServers()
			if c.set {
				servers[0].(map[string]interface{})["description"] = c.description
			}
			resourceSchema := slb.ResourceMasterSlaveServerGroup()
			d := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{
				"load_balancer_id":             "lb-test",
				"name":                         "group-test",
				"delete_protection_validation": false,
				"servers":                      servers,
			})
			if err := resourceSchema.Create(d, client); err != nil {
				t.Fatalf("create master slave server group: %s", err)
			}

			seen, request := payload()
			if !seen || request == "" {
				t.Fatalf("the fixture did not capture a CreateMasterSlaveServerGroup request (seen=%v, payload=%q)", seen, request)
			}
			if !strings.Contains(request, "'IsBackup':'0'") {
				t.Fatalf("captured MasterSlaveBackendServers payload is incomplete: %q", request)
			}
			if c.wantSubstr != "" && !strings.Contains(request, c.wantSubstr) {
				t.Fatalf("MasterSlaveBackendServers %q does not contain %q", request, c.wantSubstr)
			}
			if c.wantAbsent != "" && strings.Contains(request, c.wantAbsent) {
				t.Fatalf("MasterSlaveBackendServers %q unexpectedly contains %q", request, c.wantAbsent)
			}
		})
	}
}

func masterSlaveServerGroupUnitTestServers() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"server_id":   "i-master",
			"port":        80,
			"weight":      100,
			"type":        "ecs",
			"server_type": "Master",
		},
		map[string]interface{}{
			"server_id":   "i-slave",
			"port":        80,
			"weight":      100,
			"type":        "ecs",
			"server_type": "Slave",
		},
	}
}

func newMasterSlaveServerGroupUnitTestClient(t *testing.T) (*connectivity.AliyunClient, func() (bool, string)) {
	t.Helper()

	var (
		mu         sync.Mutex
		createSeen bool
		createForm string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request form: %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "CreateMasterSlaveServerGroup":
			mu.Lock()
			createSeen = true
			createForm = r.Form.Get("MasterSlaveBackendServers")
			mu.Unlock()
			_, _ = w.Write([]byte(`{"MasterSlaveServerGroupId":"mvsg-test"}`))
		case "DescribeMasterSlaveServerGroupAttribute":
			_, _ = w.Write([]byte(`{"MasterSlaveServerGroupId":"mvsg-test","MasterSlaveServerGroupName":"group-test","LoadBalancerId":"lb-test","MasterSlaveBackendServers":{"MasterSlaveBackendServer":[{"ServerId":"i-master","Port":80,"Weight":100,"Type":"ecs","ServerType":"Master"},{"ServerId":"i-slave","Port":80,"Weight":100,"Type":"ecs","ServerType":"Slave"}]}}`))
		default:
			_, _ = w.Write([]byte(`{}`))
		}
	}))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints.Store("slb", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := connectivity.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}

	payload := func() (bool, string) {
		mu.Lock()
		defer mu.Unlock()
		return createSeen, createForm
	}
	return client, payload
}
