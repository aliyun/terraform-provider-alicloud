package alicloud

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	sdkendpoints "github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_drds_instance", &resource.Sweeper{
		Name: "alicloud_drds_instance",
		F:    testSweepDRDSInstances,
	})
}

func testSweepDRDSInstances(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DrdsSupportedRegions) {
		log.Printf("[INFO] Skipping DRDS Instance unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := drds.CreateDescribeDrdsInstancesRequest()
	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstances(request)
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving DRDS Instances: %s", WrapError(err))
	}
	response, _ := raw.(*drds.DescribeDrdsInstancesResponse)

	vpcService := VpcService{client}
	for _, v := range response.Instances.Instance {
		name := v.Description
		id := v.DrdsInstanceId
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
			if skip {
				instanceDetailRequest := drds.CreateDescribeDrdsInstanceRequest()
				instanceDetailRequest.DrdsInstanceId = id
				raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
					return drdsClient.DescribeDrdsInstance(instanceDetailRequest)
				})
				if err != nil {
					log.Printf("[ERROR] Error retrieving DRDS Instance: %s. %s", id, WrapError(err))
				}
				instanceDetailResponse, _ := raw.(*drds.DescribeDrdsInstanceResponse)
				for _, vip := range instanceDetailResponse.Data.Vips.Vip {
					if need, err := vpcService.needSweepVpc(vip.VpcId, ""); err == nil {
						skip = !need
						break
					}
				}

			}
			if skip {
				log.Printf("[INFO] Skipping DRDS Instance: %s (%s)", name, id)
				continue
			}
		}
		log.Printf("[INFO] Deleting DRDS Instance: %s (%s)", name, id)
		req := drds.CreateRemoveDrdsInstanceRequest()
		req.DrdsInstanceId = id
		_, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
			return drdsClient.RemoveDrdsInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete DRDS Instance (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccDRDSInstanceRegion(t *testing.T) {
	t.Helper()
	// Keep service availability and provider/data-source regions consistent.
	t.Setenv("ALICLOUD_REGION", "cn-hangzhou")
	t.Setenv("CHECKOUT_REGION", "false")
	originalDefaultRegion := defaultRegionToTest
	t.Cleanup(func() { defaultRegionToTest = originalDefaultRegion })
	defaultRegionToTest = "cn-hangzhou"
}

func TestAccAliCloudDRDSInstance_Vpc(t *testing.T) {
	testAccDRDSInstanceRegion(t)
	t.Logf("DRDS acceptance test region: %s", defaultRegionToTest)
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "alicloud_drds_instance.default"
	ra := resourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sDrdsdatabase-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDRDSInstanceConfigDependence)
	createConfig := testAccConfig(map[string]interface{}{
		"description":          "${var.name}",
		"zone_id":              "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
		"instance_series":      "${var.instance_series}",
		"instance_charge_type": "PostPaid",
		"vswitch_id":           "${data.alicloud_vswitches.default.vswitches.0.id}",
		"specification":        "drds.sn2.4c16g.8C32G",
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
			if region := os.Getenv("ALICLOUD_REGION"); region != "cn-hangzhou" {
				t.Fatalf("ALICLOUD_REGION = %q, want cn-hangzhou", region)
			}
			if defaultRegionToTest != "cn-hangzhou" {
				t.Fatalf("defaultRegionToTest = %q, want cn-hangzhou", defaultRegionToTest)
			}
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":     name,
						"mysql_version":   "5",
						"specification":   "drds.sn2.4c16g.8C32G",
						"instance_series": "drds.sn2.4c16g",
						"vswitch_id":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             createConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudDRDSInstance_Multi(t *testing.T) {
	testAccDRDSInstanceRegion(t)
	t.Logf("DRDS acceptance test region: %s", defaultRegionToTest)
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "alicloud_drds_instance.default.2"
	ra := resourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sDrdsdatabase-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDRDSInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
			testAccPreCheckWithRegions(t, false, connectivity.DrdsClassicNoSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"zone_id":              "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"instance_series":      "${var.instance_series}",
					"instance_charge_type": "PostPaid",
					"specification":        "drds.sn2.4c16g.8C32G",
					"vswitch_id":           "${data.alicloud_vswitches.default.vswitches.0.id}",
					"count":                "3",
					"mysql_version":        "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudDRDSInstance_VpcId(t *testing.T) {
	testAccDRDSInstanceRegion(t)
	t.Logf("DRDS acceptance test region: %s", defaultRegionToTest)
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "alicloud_drds_instance.default"
	ra := resourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sDrdsdatabase-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDRDSInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"zone_id":              "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"instance_series":      "${var.instance_series}",
					"instance_charge_type": "PostPaid",
					"vswitch_id":           "${data.alicloud_vswitches.default.vswitches.0.id}",
					"specification":        "drds.sn2.4c16g.8C32G",
					"vpc_id":               "${data.alicloud_vpcs.default.ids.0}",
					"mysql_version":        "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
						"vpc_id":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccAliCloudDRDSInstance_MySQLVersion(t *testing.T) {
	testAccDRDSInstanceRegion(t)
	t.Logf("DRDS acceptance test region: %s", defaultRegionToTest)
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "alicloud_drds_instance.default"
	ra := resourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sDrdsdatabase-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDRDSInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"zone_id":              "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"instance_series":      "${var.instance_series}",
					"instance_charge_type": "PostPaid",
					"vswitch_id":           "${data.alicloud_vswitches.default.vswitches.0.id}",
					"specification":        "drds.sn2.4c16g.8C32G",
					"mysql_version":        "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   name,
						"mysql_version": "5",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func resourceDRDSInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
	provider "alicloud" {
	  region = "cn-hangzhou"
	}
	variable "name" {
		default = "%s"
	}
	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	
	variable "instance_series" {
		default = "drds.sn2.4c16g"
	}
	
	data "alicloud_vpcs" "default"	{
        name_regex = "default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
	  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
	}
`, name)
}

var drdsInstancebasicMap = map[string]string{
	"description":          CHECKSET,
	"zone_id":              CHECKSET,
	"instance_series":      "drds.sn2.4c16g",
	"instance_charge_type": "PostPaid",
	"specification":        "drds.sn2.4c16g.8C32G",
	"connection_string":    CHECKSET,
	"port":                 CHECKSET,
}

// TestUnitAliCloudDRDSInstanceFlattenVips guards the DRDS instance Read against a
// panic when the DescribeDrdsInstance response returns an empty VIP list. A valid
// instance can transiently report no VIPs, and indexing Vip[0] directly used to
// crash the provider. The empty case must return zero values; the populated case
// must surface paired network IDs and the intranet connection string/port.
func TestUnitAliCloudDRDSInstanceFlattenVips(t *testing.T) {
	// Empty VIP list: must not panic and must yield zero values.
	vpcId, connectionString, port, vswitchId := flattenDrdsInstanceVips([]drds.Vip{})
	if vpcId != "" || connectionString != "" || port != "" || vswitchId != "" {
		t.Fatalf("empty VIP list should yield zero values, got vpcId=%q connectionString=%q port=%q vswitchId=%q", vpcId, connectionString, port, vswitchId)
	}

	// Nil VIP list: same guarantee.
	vpcId, connectionString, port, vswitchId = flattenDrdsInstanceVips(nil)
	if vpcId != "" || connectionString != "" || port != "" || vswitchId != "" {
		t.Fatalf("nil VIP list should yield zero values, got vpcId=%q connectionString=%q port=%q vswitchId=%q", vpcId, connectionString, port, vswitchId)
	}

	// The first VIP has no VSwitch. Both network IDs must come from the next VIP.
	vips := []drds.Vip{
		{Type: "internet", VpcId: "vpc-external", Dns: "public.example.com", Port: "3306"},
		{Type: "intranet", VpcId: "vpc-internal", Dns: "intranet.example.com", Port: "3307", VswitchId: "vsw-intranet"},
	}
	vpcId, connectionString, port, vswitchId = flattenDrdsInstanceVips(vips)
	if vpcId != "vpc-internal" {
		t.Fatalf("vpcId should come from the VIP with the VSwitch, got %q", vpcId)
	}
	if connectionString != "intranet.example.com" {
		t.Fatalf("connectionString should come from the intranet VIP, got %q", connectionString)
	}
	if port != "3307" {
		t.Fatalf("port should come from the intranet VIP, got %q", port)
	}
	if vswitchId != "vsw-intranet" {
		t.Fatalf("vswitchId should be the first non-empty VswitchId, got %q", vswitchId)
	}
}

func TestUnitAliCloudDRDSInstanceRegionIsolation(t *testing.T) {
	t.Setenv("ALICLOUD_REGION", "cn-beijing")
	t.Setenv("CHECKOUT_REGION", "true")
	originalDefaultRegion := defaultRegionToTest
	t.Cleanup(func() { defaultRegionToTest = originalDefaultRegion })
	defaultRegionToTest = "cn-beijing"

	t.Run("hangzhou", func(t *testing.T) {
		testAccDRDSInstanceRegion(t)
		if got := os.Getenv("ALICLOUD_REGION"); got != "cn-hangzhou" {
			t.Fatalf("ALICLOUD_REGION = %q, want cn-hangzhou", got)
		}
		if got := os.Getenv("CHECKOUT_REGION"); got != "false" {
			t.Fatalf("CHECKOUT_REGION = %q, want false", got)
		}
		if defaultRegionToTest != "cn-hangzhou" {
			t.Fatalf("defaultRegionToTest = %q, want cn-hangzhou", defaultRegionToTest)
		}
	})

	if got := os.Getenv("ALICLOUD_REGION"); got != "cn-beijing" {
		t.Errorf("ALICLOUD_REGION after cleanup = %q, want cn-beijing", got)
	}
	if got := os.Getenv("CHECKOUT_REGION"); got != "true" {
		t.Errorf("CHECKOUT_REGION after cleanup = %q, want true", got)
	}
	if defaultRegionToTest != "cn-beijing" {
		t.Errorf("defaultRegionToTest after cleanup = %q, want cn-beijing", defaultRegionToTest)
	}
}

func drdsInstanceReadTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	for _, product := range []string{"drds", "vpc"} {
		product := product
		endpoint := sdkendpoints.GetEndpointFromMap("cn-hangzhou", product)
		t.Cleanup(func() { sdkendpoints.AddEndpointMapping("cn-hangzhou", product, endpoint) })
	}
	for _, variable := range []string{"HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy", "ALL_PROXY", "all_proxy", "TF_ENDPOINT_PATH"} {
		t.Setenv(variable, "")
	}
	t.Setenv("NO_PROXY", "*")
	t.Setenv("no_proxy", "*")
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	var endpoints, signVersion sync.Map
	endpoints.Store("drds", strings.TrimPrefix(server.URL, "http://"))
	endpoints.Store("vpc", strings.TrimPrefix(server.URL, "http://"))
	config := &connectivity.Config{
		Region: connectivity.Hangzhou, RegionId: "cn-hangzhou", Protocol: "HTTP",
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		Endpoints: &endpoints, SignVersion: &signVersion,
		AccountType: "test", SkipRegionValidation: true,
		ClientReadTimeout: 1000, ClientConnectTimeout: 1000,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	return client
}

const drdsInstanceCompleteResponse = `{"Data":{"Status":"RUN","ZoneId":"cn-hangzhou-j","Description":"test instance","InstanceSeries":"drds.sn2.4c16g","InstanceSpec":"8C32G","CommodityCode":"drdsPost","MysqlVersion":5,"Vips":{"Vip":[{"Type":"intranet","VpcId":"vpc-test","VswitchId":"vsw-test","Dns":"db.example.com","Port":"3306"}]}}}`

func TestUnitDrdsInstanceReadConsistency(t *testing.T) {
	for _, tc := range []struct {
		name     string
		response string
		want     map[string]interface{}
	}{
		{name: "non-prefixed specification", response: drdsInstanceCompleteResponse},
		{name: "prefixed specification", response: strings.Replace(drdsInstanceCompleteResponse, `"8C32G"`, `"drds.sn2.4c16g.8C32G"`, 1), want: map[string]interface{}{"specification": "drds.sn2.4c16g.8C32G"}},
		{name: "blank description", response: strings.Replace(drdsInstanceCompleteResponse, `"Description":"test instance"`, `"Description":""`, 1), want: map[string]interface{}{"description": ""}},
		{name: "public VIP precedes network VIP", response: strings.Replace(drdsInstanceCompleteResponse, `"Vip":[`, `"Vip":[{"Type":"internet","Dns":"public.example.com","Port":"3306"},`, 1)},
		{name: "lowercase dns", response: strings.Replace(drdsInstanceCompleteResponse, `"Dns"`, `"dns"`, 1)},
		{name: "numeric port", response: strings.Replace(drdsInstanceCompleteResponse, `"Port":"3306"`, `"Port":3306`, 1)},
		{name: "PrePaid mapping", response: strings.Replace(drdsInstanceCompleteResponse, "drdsPost", "drdsPre", 1), want: map[string]interface{}{"instance_charge_type": "PrePaid"}},
		{name: "no VSwitch", response: strings.Replace(drdsInstanceCompleteResponse, `"VswitchId":"vsw-test",`, "", 1), want: map[string]interface{}{"vpc_id": "", "vswitch_id": ""}},
		{name: "empty VIPs", response: strings.Replace(drdsInstanceCompleteResponse, `[{"Type":"intranet","VpcId":"vpc-test","VswitchId":"vsw-test","Dns":"db.example.com","Port":"3306"}]`, `[]`, 1), want: map[string]interface{}{"vpc_id": "", "vswitch_id": "", "connection_string": "", "port": ""}},
		{name: "empty fields overwrite state", response: `{"Data":{"Status":"DO_CREATE"}}`, want: map[string]interface{}{
			"zone_id": "", "description": "", "specification": "", "instance_series": "", "instance_charge_type": "",
			"vpc_id": "", "vswitch_id": "", "connection_string": "", "port": "", "mysql_version": 0, "status": "DO_CREATE",
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls int32
			client := drdsInstanceReadTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				if r.Form.Get("Action") != "DescribeDrdsInstance" {
					t.Errorf("unexpected action %q", r.Form.Get("Action"))
				}
				atomic.AddInt32(&calls, 1)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, tc.response)
			})
			want := map[string]interface{}{
				"zone_id": "cn-hangzhou-j", "description": "test instance", "specification": "8C32G", "instance_series": "drds.sn2.4c16g",
				"instance_charge_type": "PostPaid", "vpc_id": "vpc-test", "vswitch_id": "vsw-test", "connection_string": "db.example.com", "port": "3306", "mysql_version": 5, "status": "RUN",
			}
			attributes := map[string]string{}
			for key, value := range want {
				attributes[key] = fmt.Sprint(value)
			}
			for key, value := range tc.want {
				want[key] = value
			}
			r := resourceAlicloudDRDSInstance()
			if r.Timeouts.Read != nil {
				t.Fatal("Read must not declare a metadata wait timeout")
			}
			d := r.Data(&terraform.InstanceState{ID: "drds-test", Attributes: attributes})
			err := resourceAliCloudDRDSInstanceRead(d, client)
			if err != nil {
				t.Errorf("single Read failed: %v", err)
			}
			if got := atomic.LoadInt32(&calls); got != 1 {
				t.Errorf("Describe calls = %d, want 1", got)
			}
			for key, value := range want {
				if got := d.Get(key); got != value {
					t.Errorf("%s = %v, want %v", key, got, value)
				}
			}
			if d.Id() != "drds-test" {
				t.Error("successful response cleared ID")
			}
		})
	}
}

func TestUnitDrdsInstanceReadTransientError(t *testing.T) {
	for _, code := range []string{"InternalError", "Throttling", "InvalidDrdsInstanceId.NotFound", "Forbidden", "status5"} {
		t.Run(code, func(t *testing.T) {
			var calls int32
			client := drdsInstanceReadTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if atomic.AddInt32(&calls, 1) == 1 {
					if code == "status5" {
						fmt.Fprint(w, `{"Data":{"Status":"5"}}`)
						return
					}
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, `{"Code":%q,"Message":"test response"}`, code)
					return
				}
				fmt.Fprint(w, strings.Replace(drdsInstanceCompleteResponse, `"8C32G"`, `"drds.sn2.4c16g.8C32G"`, 1))
			})
			attributes := map[string]string{"description": "unchanged", "specification": "old-spec", "instance_series": "old-series", "vpc_id": "old-vpc", "vswitch_id": "old-vsw", "zone_id": "old-zone", "instance_charge_type": "PrePaid", "connection_string": "old.example.com", "port": "1234", "status": "RUN"}
			d := resourceAlicloudDRDSInstance().Data(&terraform.InstanceState{ID: "drds-test", Attributes: attributes})
			err := resourceAliCloudDRDSInstanceRead(d, client)
			if code == "InvalidDrdsInstanceId.NotFound" || code == "status5" {
				if err != nil || d.Id() != "" {
					t.Errorf("NotFound: ID=%q err=%v", d.Id(), err)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), code) || d.Id() != "drds-test" {
					t.Errorf("API error: ID=%q err=%v", d.Id(), err)
				}
				for key, want := range attributes {
					if got := d.Get(key); got != want {
						t.Errorf("%s = %v, want preserved %v", key, got, want)
					}
				}
			}
			if got := atomic.LoadInt32(&calls); got != 1 {
				t.Errorf("Describe calls = %d, want 1", got)
			}
		})
	}
}

func TestUnitDrdsInstanceCreateReadiness(t *testing.T) {
	for _, tc := range []struct {
		name       string
		response   string
		requireVPC bool
		wantCalls  int32
	}{
		{name: "non-prefixed specification", response: drdsInstanceCompleteResponse, requireVPC: true, wantCalls: 1},
		{name: "RUN missing series", response: strings.Replace(drdsInstanceCompleteResponse, `"InstanceSeries":"drds.sn2.4c16g",`, "", 1), requireVPC: true, wantCalls: 2},
		{name: "RUN missing specification", response: strings.Replace(drdsInstanceCompleteResponse, `"InstanceSpec":"8C32G",`, "", 1), requireVPC: true, wantCalls: 2},
		{name: "RUN missing billing", response: strings.Replace(drdsInstanceCompleteResponse, `"CommodityCode":"drdsPost",`, "", 1), requireVPC: true, wantCalls: 2},
		{name: "VPC creation missing VSwitch", response: strings.Replace(drdsInstanceCompleteResponse, `"VswitchId":"vsw-test",`, "", 1), requireVPC: true, wantCalls: 2},
		{name: "non-VPC creation", response: strings.Replace(drdsInstanceCompleteResponse, `"VswitchId":"vsw-test",`, "", 1), wantCalls: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls int32
			client := drdsInstanceReadTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if atomic.AddInt32(&calls, 1) == 1 {
					fmt.Fprint(w, tc.response)
				} else {
					fmt.Fprint(w, drdsInstanceCompleteResponse)
				}
			})
			service := DrdsService{client}
			if err := service.waitDrdsInstanceReady("drds-test", tc.requireVPC, 2*time.Second); err != nil {
				t.Fatal(err)
			}
			if got := atomic.LoadInt32(&calls); got != tc.wantCalls {
				t.Errorf("Describe calls = %d, want %d", got, tc.wantCalls)
			}
		})
	}
}
