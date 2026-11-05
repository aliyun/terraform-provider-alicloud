package drds_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	drdssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/drds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_drds_instance", &resource.Sweeper{
		Name: "alicloud_drds_instance",
		F:    testSweepDRDSInstances,
	})
}

func testSweepDRDSInstances(region string) error {
	if sweep.SweepPreCheckWithRegions(region, true, connectivity.DrdsSupportedRegions) {
		log.Printf("[INFO] Skipping DRDS Instance unsupported region: %s", region)
		return nil
	}
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := drdssdk.CreateDescribeDrdsInstancesRequest()
	raw, err := client.WithDrdsClient(func(drdsClient *drdssdk.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstances(request)
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving DRDS Instances: %s", errs.WrapError(err))
	}
	response, _ := raw.(*drdssdk.DescribeDrdsInstancesResponse)

	for _, v := range response.Instances.Instance {
		name := v.Description
		id := v.DrdsInstanceId
		skip := true
		if !sweep.SweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
			if skip {
				instanceDetailRequest := drdssdk.CreateDescribeDrdsInstanceRequest()
				instanceDetailRequest.DrdsInstanceId = id
				raw, err := client.WithDrdsClient(func(drdsClient *drdssdk.Client) (interface{}, error) {
					return drdsClient.DescribeDrdsInstance(instanceDetailRequest)
				})
				if err != nil {
					log.Printf("[ERROR] Error retrieving DRDS Instance: %s. %s", id, errs.WrapError(err))
				}
				instanceDetailResponse, _ := raw.(*drdssdk.DescribeDrdsInstanceResponse)
				for _, vip := range instanceDetailResponse.Data.Vips.Vip {
					if need, err := needSweepVpc(client, vip.VpcId); err == nil {
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
		req := drdssdk.CreateRemoveDrdsInstanceRequest()
		req.DrdsInstanceId = id
		_, err := client.WithDrdsClient(func(drdsClient *drdssdk.Client) (interface{}, error) {
			return drdsClient.RemoveDrdsInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete DRDS Instance (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

// needSweepVpc replicates the VpcId-non-empty path of the root-package
// VpcService.needSweepVpc together with the DescribeVpc it delegates to (see
// alicloud/service_alicloud_common_test.go and service_alicloud_vpc.go), so the
// DRDS sweeper keeps its secondary VPC-name recheck after migrating to a leaf
// package, without importing the product-owned VPC service package or the root
// package. It calls DescribeVpcs directly through the shared connectivity client.
//
// The baseline DescribeVpc guards the returned VPC by identity: the first
// element's VpcId must equal the requested id, otherwise the entry is treated as
// missing. This identity check (service_alicloud_vpc.go L77), addDebug, the hit
// Debug log and the WrapErrorf error context are all preserved here. The error
// contract is preserved too: a missing VPC (NotFound error codes, an empty list
// or a VpcId identity mismatch) returns (false, nil) so the caller's `err == nil`
// guard treats it as "no sweep needed"; any other error is returned with its
// original context so the caller skips acting on it, matching the original
// behavior. Only the VpcId path is needed because the DRDS sweeper always passes
// an empty vswitchId.
func needSweepVpc(client *connectivity.AliyunClient, vpcId string) (bool, error) {
	if vpcId == "" {
		return false, nil
	}
	action := "DescribeVpcs"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"VpcId":    vpcId,
	}
	response, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Forbidden.VpcNotFound", "InvalidVpcID.NotFound"}) {
			return false, nil
		}
		return false, errs.WrapErrorf(err, errs.DefaultErrorMsg, vpcId, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Vpcs.Vpc", response)
	if err != nil {
		return false, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, vpcId, "$.Vpcs.Vpc", response)
	}
	if len(v.([]interface{})) < 1 {
		return false, nil
	}
	object := v.([]interface{})[0].(map[string]interface{})
	// Preserve the baseline DescribeVpc identity guard (service_alicloud_vpc.go
	// L77): a returned entry whose VpcId differs from the requested one is
	// treated as missing, so the sweeper does not act on a stale/foreign VPC.
	if object["VpcId"].(string) != vpcId {
		return false, nil
	}
	name := strings.ToLower(object["VpcName"].(string))
	if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
		log.Printf("[DEBUG] Need to sweep the VPC (%v (%v)).", object["VpcId"], object["VpcName"])
		return true, nil
	}
	return false, nil
}

// drdsNeedSweepVpcTestClient stands up a credential-free AliyunClient whose VPC
// RPC transport is pointed at an in-process httptest server, mirroring the
// rds/agentloop unit-test construction. This keeps needSweepVpc exercisable as
// a plain `go test` without gomonkey binary patching or real cloud credentials.
func drdsNeedSweepVpcTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("vpc", endpoint)
	return client
}

// TestUnitDRDSNeedSweepVpc locks the selection boundaries of the migrated
// needSweepVpc helper against an in-process VPC mock, so the保真 of the sweep
// recheck stays guarded after the DRDS package migration. It covers the four
// boundaries required by the round-3 review plus the restored VpcId identity
// guard: a tf-testacc VPC name triggers a sweep (and the DEBUG log); a non-test
// name does not; a NotFound error code is swallowed to (false, nil); any other
// error is surfaced with the DescribeVpcs action context; and a returned VPC
// whose VpcId differs from the request is rejected by the identity guard even
// when its name would match.
func TestUnitDRDSNeedSweepVpc(t *testing.T) {
	const vpcID = "vpc-test-f5"

	writeVpc := func(w http.ResponseWriter, vpcID, vpcName string) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Vpcs": map[string]interface{}{
				"Vpc": []map[string]interface{}{
					{"VpcId": vpcID, "VpcName": vpcName},
				},
			},
		})
	}
	writeErr := func(w http.ResponseWriter, code, message string) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Code": code, "Message": message,
		})
	}

	cases := []struct {
		name    string
		handler http.HandlerFunc
		want    bool
		wantErr bool
		wantLog bool
	}{
		{
			name: "prefix_hit_sweeps",
			handler: func(w http.ResponseWriter, r *http.Request) {
				writeVpc(w, vpcID, "tf-testacc-drds-instance")
			},
			want: true, wantErr: false, wantLog: true,
		},
		{
			name: "non_hit_no_sweep",
			handler: func(w http.ResponseWriter, r *http.Request) {
				writeVpc(w, vpcID, "prod-drds-instance")
			},
			want: false, wantErr: false, wantLog: false,
		},
		{
			name: "notfound_swallowed",
			handler: func(w http.ResponseWriter, r *http.Request) {
				writeErr(w, "Forbidden.VpcNotFound", "the specified VPC is not found")
			},
			want: false, wantErr: false, wantLog: false,
		},
		{
			name: "other_error_surfaced",
			handler: func(w http.ResponseWriter, r *http.Request) {
				writeErr(w, "Internal.Error", "boom")
			},
			want: false, wantErr: true, wantLog: false,
		},
		{
			name: "identity_mismatch_no_sweep",
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Name would match, but the returned VpcId is foreign: the
				// identity guard must reject it instead of sweeping.
				writeVpc(w, "vpc-someone-else", "tf-testacc-stale")
			},
			want: false, wantErr: false, wantLog: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := drdsNeedSweepVpcTestClient(t, tc.handler)

			var buf bytes.Buffer
			orig := log.Writer()
			log.SetOutput(&buf)
			defer log.SetOutput(orig)

			got, err := needSweepVpc(client, vpcID)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got nil")
				}
				if !strings.Contains(err.Error(), "DescribeVpcs") {
					t.Fatalf("error should carry the DescribeVpcs action context, got %q", err.Error())
				}
			} else if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got != tc.want {
				t.Fatalf("needSweepVpc = %v, want %v", got, tc.want)
			}
			hasLog := strings.Contains(buf.String(), "Need to sweep the VPC")
			if hasLog != tc.wantLog {
				t.Fatalf("debug log hit = %v, want %v (buf=%q)", hasLog, tc.wantLog, buf.String())
			}
		})
	}
}

func TestAccAliCloudDRDSInstance_Vpc(t *testing.T) {
	var v *drdssdk.DescribeDrdsInstanceResponse

	resourceId := "alicloud_drds_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return drds.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sDrdsdatabase-%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceDRDSInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
		},
		// module name
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"zone_id":              "${alicloud_vswitch.default.zone_id}",
					"instance_series":      "${var.instance_series}",
					"instance_charge_type": "PostPaid",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"specification":        "drds.sn2.4c16g.8c32g",
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
	var v *drdssdk.DescribeDrdsInstanceResponse

	resourceId := "alicloud_drds_instance.default.2"
	ra := acctest.ResourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return drds.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sDrdsdatabase-%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceDRDSInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"zone_id":              "${alicloud_vswitch.default.zone_id}",
					"instance_series":      "${var.instance_series}",
					"instance_charge_type": "PostPaid",
					"specification":        "drds.sn2.4c16g.8c32g",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"count":                "3",
					"mysql_version":        "5",
					"vpc_id":               "${alicloud_vswitch.default.vpc_id}",
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
	var v *drdssdk.DescribeDrdsInstanceResponse

	resourceId := "alicloud_drds_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return drds.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sDrdsdatabase-%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceDRDSInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"zone_id":              "${alicloud_vswitch.default.zone_id}",
					"instance_series":      "${var.instance_series}",
					"instance_charge_type": "PostPaid",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"specification":        "drds.sn2.4c16g.8c32g",
					"vpc_id":               "${alicloud_vswitch.default.vpc_id}",
					"mysql_version":        "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
						"vpc_id":      acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
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

func TestAccAliCloudDRDSInstance_MySQLVersion(t *testing.T) {
	var v *drdssdk.DescribeDrdsInstanceResponse

	resourceId := "alicloud_drds_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return drds.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sDrdsdatabase-%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceDRDSInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"zone_id":              "${alicloud_vswitch.default.zone_id}",
					"instance_series":      "${var.instance_series}",
					"instance_charge_type": "PostPaid",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"specification":        "drds.sn2.4c16g.8c32g",
					"vpc_id":               "${alicloud_vswitch.default.vpc_id}",
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
	variable "name" {
		default = "%s"
	}

	variable "instance_series" {
		default = "drds.sn2.4c16g"
	}

	# DRDS 1.0 (classic, Type=1/V1) instance creation is only supported in
	# specific zones of cn-beijing. An existing running DRDS 1.0 instance
	# (drds.sn2.4c16g, drdsPost, V1) in this account confirms cn-beijing-h is a
	# valid DRDS 1.0 zone; the first VSwitch-capable zone (cn-beijing-f) is a
	# DRDS 2.0/polardbx zone and returns InternalError on classic CreateDrdsInstance.
	# ACC region is fixed to cn-beijing, so the zone is pinned to cn-beijing-h.
	resource "alicloud_vpc" "default" {
		name       = "tf-testAcc-vpc-${var.name}"
		cidr_block = "172.16.0.0/16"
	}
	resource "alicloud_vswitch" "default" {
		name       = "tf-testAcc-vsw-${var.name}"
		vpc_id     = alicloud_vpc.default.id
		cidr_block = "172.16.0.0/24"
		zone_id    = "cn-beijing-h"
	}
`, name)
}

var drdsInstancebasicMap = map[string]string{
	"description":          acctest.CHECKSET,
	"zone_id":              acctest.CHECKSET,
	"instance_series":      "drds.sn2.4c16g",
	"instance_charge_type": "PostPaid",
	"specification":        "drds.sn2.4c16g.8c32g",
	"connection_string":    acctest.CHECKSET,
	"port":                 acctest.CHECKSET,
}
