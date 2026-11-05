package gwlb_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/gwlb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Gwlb LoadBalancer. >>> Resource test cases, automatically generated.
// Case LoadBalancer Test 7994
func TestAccAliCloudGwlbLoadBalancer_basic7994(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gwlb_load_balancer.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGwlbLoadBalancerMap7994)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gwlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgwlbloadbalancer%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGwlbLoadBalancerBasicDependence7994)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":             "${alicloud_vpc.defaulti9Axhl.id}",
					"load_balancer_name": name,
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default9NaKmL.id}",
							"zone_id":    "${var.zone_id1}",
						},
					},
					"address_ip_version": "Ipv4",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"dry_run":            "false",
					"traffic_mode":       "ByPass",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":             acctest.CHECKSET,
						"load_balancer_name": name,
						"zone_mappings.#":    "1",
						"address_ip_version": "Ipv4",
						"resource_group_id":  acctest.CHECKSET,
						"dry_run":            "false",
						"traffic_mode":       "ByPass",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaultH4pKT4.id}",
							"zone_id":    "${var.zone_id2}",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"traffic_mode":      "LoadBalance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
						"zone_mappings.#":    "1",
						"resource_group_id":  acctest.CHECKSET,
						"traffic_mode":       "LoadBalance",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"traffic_mode":      "ByPass",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
						"traffic_mode":      "ByPass",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_mode": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mode": "ByPass",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudGwlbLoadBalancerMap7994 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudGwlbLoadBalancerBasicDependence7994(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

variable "zone_id2" {
  default = "cn-wulanchabu-c"
}

variable "zone_id1" {
  default = "cn-wulanchabu-b"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaulti9Axhl" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default9NaKmL" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%%s1", var.name)
}

resource "alicloud_vswitch" "defaultH4pKT4" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id2
  cidr_block   = "10.0.1.0/24"
  vswitch_name = format("%%s2", var.name)
}


`, name)
}

// Case LoadBalancer Test_依赖资源 8565
func TestAccAliCloudGwlbLoadBalancer_basic8565(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gwlb_load_balancer.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGwlbLoadBalancerMap8565)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gwlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgwlbloadbalancer%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGwlbLoadBalancerBasicDependence8565)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":             "${alicloud_vpc.defaulti9Axhl.id}",
					"load_balancer_name": name,
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default9NaKmL.id}",
							"zone_id":    "${var.zone_id1}",
						},
					},
					"address_ip_version": "Ipv4",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"dry_run":            "false",
					"traffic_mode":       "ByPass",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":             acctest.CHECKSET,
						"load_balancer_name": name,
						"zone_mappings.#":    "1",
						"address_ip_version": "Ipv4",
						"resource_group_id":  acctest.CHECKSET,
						"dry_run":            "false",
						"traffic_mode":       "ByPass",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaultH4pKT4.id}",
							"zone_id":    "${var.zone_id2}",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"traffic_mode":      "LoadBalance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
						"zone_mappings.#":    "1",
						"resource_group_id":  acctest.CHECKSET,
						"traffic_mode":       "LoadBalance",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"traffic_mode":      "ByPass",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
						"traffic_mode":      "ByPass",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_mode": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mode": "ByPass",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudGwlbLoadBalancerMap8565 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudGwlbLoadBalancerBasicDependence8565(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

variable "zone_id2" {
  default = "cn-wulanchabu-c"
}

variable "zone_id1" {
  default = "cn-wulanchabu-b"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaulti9Axhl" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default9NaKmL" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%%s1", var.name)
}

resource "alicloud_vswitch" "defaultH4pKT4" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id2
  cidr_block   = "10.0.1.0/24"
  vswitch_name = format("%%s2", var.name)
}


`, name)
}

// Test Gwlb LoadBalancer. <<< Resource test cases, automatically generated.

// gwlbLoadBalancerUnitClient serves handler from a local HTTP endpoint and returns a client
// bound to it, so unit tests can drive the resource entry points without any cloud access.
func gwlbLoadBalancerUnitClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{AccessKey: "test-key", SecretKey: "test-secret", Credential: credential, RegionId: "cn-wulanchabu", AccountType: "test", Protocol: "http", Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("gwlb", endpoint)
	return client
}

// TestUnitAliCloudGwlbLoadBalancerRenameKeepsStatusOnlyWait drives the update entry with a
// state written before traffic_mode existed plus a rename-only change. It asserts that the
// request carries the new name and no TrafficMode, and that the update converges through the
// status-only wait instead of a traffic-mode target.
func TestUnitAliCloudGwlbLoadBalancerRenameKeepsStatusOnlyWait(t *testing.T) {
	r := gwlb.ResourceLoadBalancer()

	var mu sync.Mutex
	updateParams := url.Values{}
	statusPolls := 0
	client := gwlbLoadBalancerUnitClient(t, func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		params := req.URL.Query()
		if form, err := url.ParseQuery(string(body)); err == nil {
			for key, values := range form {
				params[key] = append(params[key], values...)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		switch params.Get("Action") {
		case "UpdateLoadBalancerAttribute":
			mu.Lock()
			updateParams = params
			mu.Unlock()
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"RequestId": "unit-update"})
		case "GetLoadBalancerAttribute":
			mu.Lock()
			statusPolls++
			mu.Unlock()
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"LoadBalancerId":     "lb-unit",
				"LoadBalancerStatus": "Active",
				"TrafficMode":        "LoadBalance",
			})
		default:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"RequestId": "unit-other"})
		}
	})

	state := &terraform.InstanceState{ID: "lb-unit", Attributes: map[string]string{
		"id":                 "lb-unit",
		"load_balancer_name": "before",
		"vpc_id":             "vpc-unit",
		"address_ip_version": "Ipv4",
	}}
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"load_balancer_name": "after",
		"vpc_id":             "vpc-unit",
		"address_ip_version": "Ipv4",
	})
	diff, err := r.Diff(context.Background(), state, config, client)
	if err != nil {
		t.Fatalf("diff: %s", err)
	}
	d, err := schema.InternalMap(r.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("resource data: %s", err)
	}

	done := make(chan error, 1)
	go func() { done <- r.Update(d, client) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("update: %s", err)
		}
	case <-time.After(2 * time.Minute):
		t.Fatal("update did not converge: a rename-only update must use the status-only wait")
	}

	mu.Lock()
	defer mu.Unlock()
	if values, ok := updateParams["TrafficMode"]; ok {
		t.Fatalf("rename-only update must not send TrafficMode, got %v", values)
	}
	if got := updateParams.Get("LoadBalancerName"); got != "after" {
		t.Fatalf("unexpected LoadBalancerName %q", got)
	}
	if statusPolls == 0 {
		t.Fatal("expected at least one status poll from the wait")
	}
}
