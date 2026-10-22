package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_router_interface", &resource.Sweeper{
		Name: "alicloud_router_interface",
		F:    testSweepRouterInterfaces,
	})
}

func testSweepRouterInterfaces(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var ris []vpc.RouterInterfaceType
	req := vpc.CreateDescribeRouterInterfacesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouterInterfaces(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Router Interfaces: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeRouterInterfacesResponse)
		if resp == nil || len(resp.RouterInterfaceSet.RouterInterfaceType) < 1 {
			break
		}
		ris = append(ris, resp.RouterInterfaceSet.RouterInterfaceType...)

		if len(resp.RouterInterfaceSet.RouterInterfaceType) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}
	service := VpcService{client}
	for _, v := range ris {
		name := v.Name
		id := v.RouterInterfaceId
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			// If a RI name is not set successfully, it should be fetched by vpc name and deleted.
			if skip {
				if need, err := service.needSweepVpc(v.VpcInstanceId, ""); err == nil {
					skip = !need
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Router Interface: %s (%s)", name, id)
				continue
			}
		}
		log.Printf("[INFO] Deleting Router Interface: %s (%s)", name, id)
		req := vpc.CreateDeleteRouterInterfaceRequest()
		req.RouterInterfaceId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteRouterInterface(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Router Interface (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckRouterInterfaceExists(n string, ri *vpc.RouterInterfaceType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No interface ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}

		response, err := vpcService.DescribeRouterInterface(rs.Primary.ID, client.RegionId)
		if err != nil {
			return fmt.Errorf("Error finding interface %s: %#v", rs.Primary.ID, err)
		}
		ri = &response
		return nil
	}
}

func testAccCheckRouterInterfaceDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_router_interface" {
			continue
		}

		// Try to find the interface
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}

		ri, err := vpcService.DescribeRouterInterface(rs.Primary.ID, client.RegionId)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if ri.RouterInterfaceId == rs.Primary.ID {
			return WrapError(Error("Interface %s still exists.", rs.Primary.ID))
		}
	}
	return nil
}

func TestAccAliCloudVPCRouterInterfaceBasic(t *testing.T) {
	var v vpc.RouterInterfaceType
	resourceId := "alicloud_router_interface.default"
	ra := resourceAttrInit(resourceId, testAccRouterInterfaceCheckMap)

	rand := acctest.RandInt()
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouterInterfaceConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccRouterInterfaceConfig%d", rand),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
			{
				Config: testAccRouterInterfaceConfig_role(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"role":          "InitiatingSide",
						"specification": "Large.1",
					}),
				),
			},
			{
				Config: testAccRouterInterfaceConfig_specification(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"role":          "InitiatingSide",
						"specification": "Large.2",
					}),
				),
			},
			{
				Config: testAccRouterInterfaceConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccRouterInterfaceConfig%d_change", rand),
					}),
				),
			},
			{
				Config: testAccRouterInterfaceConfig_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccRouterInterfaceConfig%d_description", rand),
					}),
				),
			},
			{
				Config: testAccRouterInterfaceConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"name":          fmt.Sprintf("tf-testAccRouterInterfaceConfig%d", rand),
						"description":   fmt.Sprintf("tf-testAccRouterInterfaceConfig%d", rand),
						"role":          "InitiatingSide",
						"specification": "Large.1",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudVPCRouterInterfaceMulti(t *testing.T) {
	var v vpc.RouterInterfaceType
	resourceId := "alicloud_router_interface.default.2"
	ra := resourceAttrInit(resourceId, testAccRouterInterfaceCheckMap)

	rand := acctest.RandInt()
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouterInterfaceConfig_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccRouterInterfaceConfig%d", rand),
					}),
				),
			},
		},
	})

}

func testAccRouterInterfaceConfigBasic(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

data "alicloud_regions" "default" {
	current = true
}

resource "alicloud_router_interface" "default" {
	opposite_region = "${data.alicloud_regions.default.regions.0.id}"
	router_type = "VRouter"
	router_id = "${alicloud_vpc.default.router_id}"
	role = "AcceptingSide"
    instance_charge_type = "PostPaid"
	name = "${var.name}"
}`, rand)
}

func testAccRouterInterfaceConfig_multi(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING"
}

data "alicloud_regions" "default" {
	current = true
}

resource "alicloud_router_interface" "default" {
	count = 3
	opposite_region = "${data.alicloud_regions.default.regions.0.id}"
	router_type = "VRouter"
	router_id = "${data.alicloud_vpcs.default.vpcs.0.router_id}"
	role = "AcceptingSide"
    instance_charge_type = "PostPaid"
	name = "${var.name}"
}`, rand)
}

func testAccRouterInterfaceConfig_role(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

data "alicloud_regions" "default" {
	current = true
}

resource "alicloud_router_interface" "default" {
	opposite_region = "${data.alicloud_regions.default.regions.0.id}"
	router_type = "VRouter"
	router_id = "${alicloud_vpc.default.router_id}"
	role = "InitiatingSide"
	specification = "Large.1"
    instance_charge_type = "PostPaid"
	name = "${var.name}"
}`, rand)
}

func testAccRouterInterfaceConfig_specification(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

data "alicloud_regions" "default" {
	current = true
}

resource "alicloud_router_interface" "default" {
	opposite_region = "${data.alicloud_regions.default.regions.0.id}"
	router_type = "VRouter"
	router_id = "${alicloud_vpc.default.router_id}"
	role = "InitiatingSide"
    instance_charge_type = "PostPaid"
	name = "${var.name}"
	specification = "Large.2"
}`, rand)
}

func testAccRouterInterfaceConfig_name(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

data "alicloud_regions" "default" {
	current = true
}

resource "alicloud_router_interface" "default" {
	opposite_region = "${data.alicloud_regions.default.regions.0.id}"
	router_type = "VRouter"
	router_id = "${alicloud_vpc.default.router_id}"
	role = "InitiatingSide"
    instance_charge_type = "PostPaid"
	name = "${var.name}_change"
	specification = "Large.2"
}`, rand)
}

func testAccRouterInterfaceConfig_description(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

data "alicloud_regions" "default" {
	current = true
}

resource "alicloud_router_interface" "default" {
	opposite_region = "${data.alicloud_regions.default.regions.0.id}"
	router_type = "VRouter"
	router_id = "${alicloud_vpc.default.router_id}"
	role = "InitiatingSide"
    instance_charge_type = "PostPaid"
	name = "${var.name}_change"
	specification = "Large.2"
	description = "${var.name}_description"
}`, rand)
}

func testAccRouterInterfaceConfig_all(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

data "alicloud_regions" "default" {
	current = true
}

resource "alicloud_router_interface" "default" {
	opposite_region = "${data.alicloud_regions.default.regions.0.id}"
	router_type = "VRouter"
	router_id = "${alicloud_vpc.default.router_id}"
	role = "InitiatingSide"
    instance_charge_type = "PostPaid"
	name = "${var.name}"
	specification = "Large.1"
	description = "${var.name}"
}`, rand)
}

var testAccRouterInterfaceCheckMap = map[string]string{
	"opposite_region":        CHECKSET,
	"router_type":            "VRouter",
	"router_id":              CHECKSET,
	"role":                   "AcceptingSide",
	"description":            "",
	"health_check_source_ip": "",
	"health_check_target_ip": "",
	"instance_charge_type":   string(PostPaid),
}

func TestAccAliCloudRouterInterface_crossRegionDelete(t *testing.T) {
	resourceId := "alicloud_router_interface.default"
	name := fmt.Sprintf("tf-testAccRouterInterface%d", acctest.RandIntRange(10000, 99999))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, routerInterfaceCrossRegionDeleteDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		ProviderFactories: testAccProviderFactoriesAlternate(),
		CheckDestroy:      testAccCheckCrossRegionRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"opposite_region":      "cn-beijing",
					"router_type":          "VRouter",
					"router_id":            "${alicloud_vpc.local.router_id}",
					"role":                 "InitiatingSide",
					"instance_charge_type": "PostPaid",
					"specification":        "Mini.2",
					"name":                 name,
					"description":          "cross region deletion test",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "opposite_region", "cn-beijing"),
					resource.TestCheckResourceAttr("alicloud_router_interface.remote", "opposite_region", "cn-hangzhou"),
					testAccCheckCrossRegionRouterInterfacesActive,
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name + "-updated",
					"description": "updated cross region deletion test",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "name", name+"-updated"),
					testAccCheckCrossRegionRouterInterfacesActive,
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

func routerInterfaceCrossRegionDeleteDependence(name string) string {
	return fmt.Sprintf(`
provider "alicloud" {
 region = "cn-hangzhou"
}
provider "alicloudalt" {
 region = "cn-beijing"
}
resource "alicloud_vpc" "local" {
 vpc_name = "%[1]s-local"
 cidr_block = "172.16.0.0/16"
}
resource "alicloud_vpc" "remote" {
 provider = alicloudalt
 vpc_name = "%[1]s-remote"
 cidr_block = "172.17.0.0/16"
}
resource "alicloud_router_interface" "remote" {
 provider = alicloudalt
 router_id = alicloud_vpc.remote.router_id
 router_type = "VRouter"
 role = "AcceptingSide"
 specification = "Negative"
 instance_charge_type = "PostPaid"
 opposite_region = "cn-hangzhou"
 name = "%[1]s-remote"
}
resource "alicloud_router_interface_connection" "local" {
 interface_id = alicloud_router_interface.default.id
 opposite_interface_id = alicloud_router_interface.remote.id
 depends_on = [alicloud_router_interface_connection.remote]
}
resource "alicloud_router_interface_connection" "remote" {
 provider = alicloudalt
 interface_id = alicloud_router_interface.remote.id
 opposite_interface_id = alicloud_router_interface.default.id
}
`, name)
}

// Verify both regions explicitly; querying the default region for an aliased
// resource would report NotFound even if the remote resource leaked.
func testAccCheckCrossRegionRouterInterfaceDestroy(s *terraform.State) error {
	for address, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_router_interface" && rs.Type != "alicloud_express_connect_router_interface" && rs.Type != "alicloud_vpc" {
			continue
		}
		region := "cn-hangzhou"
		if strings.HasSuffix(address, ".remote") {
			region = "cn-beijing"
		}
		rawClient, err := sharedClientForRegion(region)
		if err != nil {
			return err
		}
		client := rawClient.(*connectivity.AliyunClient)
		switch rs.Type {
		case "alicloud_router_interface":
			service := VpcService{client}
			_, err = service.DescribeRouterInterface(rs.Primary.ID, region)
		case "alicloud_express_connect_router_interface":
			service := ExpressConnectServiceV2{client}
			_, err = service.DescribeExpressConnectRouterInterface(rs.Primary.ID)
		case "alicloud_vpc":
			service := VpcServiceV2{client}
			_, err = service.DescribeVpcVpc(rs.Primary.ID)
		}
		if err == nil {
			return fmt.Errorf("resource %s still exists in %s", address, region)
		}
		if !NotFoundError(err) {
			return err
		}
	}
	return nil
}

func testAccCheckCrossRegionRouterInterfacesActive(s *terraform.State) error {
	for address, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_router_interface" && rs.Type != "alicloud_express_connect_router_interface" {
			continue
		}
		region := "cn-hangzhou"
		if strings.HasSuffix(address, ".remote") {
			region = "cn-beijing"
		}
		rawClient, err := sharedClientForRegion(region)
		if err != nil {
			return err
		}
		service := VpcService{rawClient.(*connectivity.AliyunClient)}
		if err := service.WaitForRouterInterface(rs.Primary.ID, region, Active, DefaultTimeout); err != nil {
			return err
		}
	}
	return nil
}

func routerInterfaceDeleteTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("RegionId") != "cn-hangzhou" {
			t.Errorf("unexpected region: %q", r.Form.Get("RegionId"))
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{AccessKey: "test-key", SecretKey: "test-secret", Credential: credential, RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http", Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("vpc", endpoint)
	return client
}

func routerInterfaceDeleteTestResponse(t *testing.T, w http.ResponseWriter, response interface{}) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(response); err != nil {
		t.Error(err)
	}
}

func routerInterfaceDeleteTestError(t *testing.T, w http.ResponseWriter, code string) {
	t.Helper()
	if strings.HasSuffix(code, ".NotFound") {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	routerInterfaceDeleteTestResponse(t, w, map[string]interface{}{"Code": code, "Message": "test operation failed"})
}

func TestUnitRouterInterfaceDeleteError(t *testing.T) {
	for _, code := range []string{"Forbidden", "InvalidInstanceId.NotFound"} {
		t.Run(code, func(t *testing.T) {
			describes := 0
			client := routerInterfaceDeleteTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Form.Get("Action") {
				case "DescribeRouterInterfaces":
					describes++
					interfaces := []interface{}{}
					if describes == 1 {
						interfaces = append(interfaces, map[string]interface{}{"RouterInterfaceId": "ri-test", "Status": "Inactive"})
					}
					routerInterfaceDeleteTestResponse(t, w, map[string]interface{}{"RouterInterfaceSet": map[string]interface{}{"RouterInterfaceType": interfaces}})
				case "DeleteRouterInterface":
					if r.Form.Get("RouterInterfaceId") != "ri-test" {
						t.Error("wrong router interface")
					}
					routerInterfaceDeleteTestError(t, w, code)
				default:
					t.Errorf("unexpected action %q", r.Form.Get("Action"))
				}
			})
			d := resourceAlicloudRouterInterface().Data(nil)
			d.SetId("ri-test")
			err := resourceAlicloudRouterInterfaceDelete(d, client)
			if code == "Forbidden" && (err == nil || !strings.Contains(err.Error(), code)) {
				t.Errorf("expected original delete error %s, got %v", code, err)
			}
			if code == "InvalidInstanceId.NotFound" && err != nil {
				t.Errorf("already deleted must succeed: %v", err)
			}
			if describes != 1 {
				t.Errorf("delete error must not fall through to the waiter: got %d reads", describes)
			}
		})
	}
}

func TestUnitRouterInterfaceRegionalEndpoints(t *testing.T) {
	regions := []string{"unit-test-router-interface-a", "unit-test-router-interface-b"}
	caches := []*sync.Map{sharedEndpointCacheForRegion(regions[0]), sharedEndpointCacheForRegion(regions[1])}
	if caches[0] == caches[1] {
		t.Fatal("different regions share an endpoint cache")
	}
	for i, region := range regions {
		if sharedEndpointCacheForRegion(region) != caches[i] {
			t.Fatal("same region does not reuse its endpoint cache")
		}
		t.Cleanup(func() { regionalEndpoints.Delete(region) })
	}
	calls := []int32{0, 0}
	servers := make([]*httptest.Server, len(regions))
	for i, region := range regions {
		i, region := i, region
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				t.Error(err)
			}
			if r.Form.Get("RegionId") != region || r.Form.Get("Action") != "DescribeRouterInterfaces" {
				t.Errorf("region %s received action %s for region %s", region, r.Form.Get("Action"), r.Form.Get("RegionId"))
			}
			atomic.AddInt32(&calls[i], 1)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"TotalCount":1,"RouterInterfaceSet":{"RouterInterfaceType":[{"RouterInterfaceId":"ri-%s","Status":"Active"}]}}`, region)
		}))
		t.Cleanup(servers[i].Close)
		caches[i].Store("vpc", strings.TrimPrefix(servers[i].URL, "http://"))
	}
	t.Setenv("NO_PROXY", "*")
	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	// Rebuild A after B as well: a new regional client must not inherit B's endpoint.
	for _, i := range []int{0, 1, 0, 1} {
		region := regions[i]
		config := &connectivity.Config{Region: connectivity.Region(region), RegionId: region,
			AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
			AccountType: "test", Protocol: "http", SkipRegionValidation: true,
			Endpoints: sharedEndpointCacheForRegion(region), SignVersion: new(sync.Map)}
		client, err := config.Client()
		if err != nil {
			t.Fatal(err)
		}
		service := VpcService{client}
		result, err := service.DescribeRouterInterface("ri-"+region, region)
		if err != nil {
			t.Fatal(err)
		}
		if result.Status != "Active" {
			t.Fatalf("region %s returned status %s", region, result.Status)
		}
	}
	if atomic.LoadInt32(&calls[0]) != 2 || atomic.LoadInt32(&calls[1]) != 2 {
		t.Fatalf("requests reached the wrong regional servers: %v", calls)
	}
}
