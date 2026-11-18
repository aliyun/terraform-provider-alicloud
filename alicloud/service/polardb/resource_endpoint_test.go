package polardb_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	sdkpolardb "github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var privateConnectionStringRegexp = "^[a-z-A-Z-0-9]+.rwlb.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com"

func TestAccAliCloudPolarDBEndpointConfigUpdate(t *testing.T) {
	var v *sdkpolardb.DBEndpoint
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBendpoint-%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_endpoint.default"
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterEndpoint")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePolarDBEndpointConfigDependence)
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
					"db_cluster_id": "${alicloud_polardb_cluster.cluster.id}",
					"endpoint_type": "Custom",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_type": "Custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_add_new_nodes": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_add_new_nodes": "Enable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_write_mode": "ReadWrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_write_mode": "ReadWrite",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nodes": []string{"${data.alicloud_polardb_clusters.default.clusters.0.db_nodes.0.db_node_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nodes.#": "1",
					}),
				),
			},
			//todo: After resource polardb_node is supported, it is necessary to add a modification check on the “nodes” parameter
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled":     "Enable",
					"net_type":        "Private",
					"ssl_auto_rotate": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":           "Enable",
						"net_type":              "Private",
						"ssl_connection_string": REGEXMATCH + privateConnectionStringRegexp,
						"ssl_auto_rotate":       "Enable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_endpoint_description": "polar_db_endpoint_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_endpoint_description": "polar_db_endpoint_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": "tf-privatetestacc-3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": "tf-privatetestacc-3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3307",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3307",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_config": map[string]string{
						"ConsistLevel": "1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_config.ConsistLevel": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"net_type", "endpoint_config"},
			},
		},
	})
}

func TestAccAliCloudPolarDBEndpointConfigUpdate_SslConnectionStringAndConnectionPrefix(t *testing.T) {
	var v *sdkpolardb.DBEndpoint
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBendpoint-%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_endpoint.default"
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterEndpoint")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePolarDBEndpointConfigDependence)
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
					"db_cluster_id": "${alicloud_polardb_cluster.cluster.id}",
					"endpoint_type": "Custom",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_type": "Custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_add_new_nodes": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_add_new_nodes": "Enable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_write_mode": "ReadWrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_write_mode": "ReadWrite",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nodes": []string{"${data.alicloud_polardb_clusters.default.clusters.0.db_nodes.0.db_node_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nodes.#": "1",
					}),
				),
			},
			//todo: After resource polardb_node is supported, it is necessary to add a modification check on the “nodes” parameter
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": "tf-privatetestacc-4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": "tf-privatetestacc-4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled":     "Enable",
					"net_type":        "Private",
					"ssl_auto_rotate": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":           "Enable",
						"net_type":              "Private",
						"ssl_connection_string": REGEXMATCH + privateConnectionStringRegexp,
						"ssl_auto_rotate":       "Enable",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"net_type", "endpoint_config"},
			},
		},
	})
}

func TestAccAliCloudPolarDBEndpointSccModeCreate(t *testing.T) {
	var v *sdkpolardb.DBEndpoint
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBendpoint-scc-%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_endpoint.default"
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterEndpoint")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePolarDBEndpointConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":   "${alicloud_polardb_cluster.cluster.id}",
					"endpoint_type":   "Custom",
					"read_write_mode": "ReadWrite",
					"endpoint_config": map[string]string{
						"ConsistLevel": "0",
					},
					"scc_mode": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_type":                "Custom",
						"read_write_mode":              "ReadWrite",
						"endpoint_config.ConsistLevel": "0",
						"scc_mode":                     "on",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"net_type", "endpoint_config"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scc_mode": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scc_mode": "on",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"net_type", "endpoint_config"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scc_mode": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scc_mode": "off",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"net_type", "endpoint_config"},
			},
		},
	})
}

func TestUnitPolarDBEndpointSccMode(t *testing.T) {
	sccModeSchema := polardb.ResourceEndpoint().Schema["scc_mode"]
	for _, input := range []string{"on", "off"} {
		if _, errors := sccModeSchema.ValidateFunc(input, "scc_mode"); len(errors) > 0 {
			t.Fatalf("expected %q to be valid, got: %v", input, errors)
		}
	}
	for _, input := range []string{"ON", "OFF", "INVALID"} {
		if _, errors := sccModeSchema.ValidateFunc(input, "scc_mode"); len(errors) == 0 {
			t.Fatalf("expected %q to be rejected", input)
		}
	}
}

func resourcePolarDBEndpointConfigDependence(name string) string {
	return fmt.Sprintf(`
        variable "name" {
                default = "%s"
        }

		data "alicloud_vpcs" "default" {
			name_regex = "^default-NODELETING$"
		}

		resource "alicloud_vpc" "default" {
			vpc_name = var.name
		}
	
		resource "alicloud_vswitch" "default" {
			zone_id = data.alicloud_polardb_node_classes.default.classes.0.zone_id
			vpc_id = alicloud_vpc.default.id
			cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
		}

		data "alicloud_polardb_node_classes" "default" {
			pay_type   = "PostPaid"
			db_type    = "MySQL"
			db_version = "8.0"
			category   = "Normal"
		}

        resource "alicloud_polardb_cluster" "cluster" {
                db_type = "MySQL"
                db_version = "8.0"
                pay_type = "PostPaid"
                db_node_class = "polar.mysql.x4.medium"
                vswitch_id = alicloud_vswitch.default.id
                description = "${var.name}"
                lifecycle {
                        ignore_changes = [encrypt_new_tables]
                }
        }

		data "alicloud_polardb_clusters" "default" {
		  	ids = [alicloud_polardb_cluster.cluster.id]
		}
        `, name)
}

// Test PolarDB Endpoint DynamoDB type
func TestAccAliCloudPolarDBEndpoint_dynamoDB(t *testing.T) {
	var v *sdkpolardb.DBEndpoint
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBendpoint-dynamodb-%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_endpoint.default"
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterEndpoint")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePolarDBEndpointDynamoDBDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id": "${alicloud_polardb_cluster.default.id}",
					"endpoint_type": "DynamoDB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"endpoint_type": "DynamoDB",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourcePolarDBEndpointDynamoDBDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.7.id
  vswitch_name = var.name
}

resource "alicloud_polardb_cluster" "default" {
  db_type         = "PostgreSQL"
  db_version      = "14"
  db_node_class   = "polar.pg.x4.medium"
  pay_type        = "PostPaid"
  vswitch_id      = alicloud_vswitch.default.id
  enable_dynamodb = true
  lifecycle {
    ignore_changes = [encrypt_new_tables]
  }
}
`, name)
}

type sccModeEndpointFixture struct {
	*httptest.Server
	mu            sync.Mutex
	endpointID    string
	sccMode       string
	readWriteMode string
	created       atomic.Bool
	createForm    map[string][]string
	modifyForm    map[string][]string
	createCalls   int
	modifyCalls   int
}

func newSccModeEndpointFixture(t *testing.T, endpointID, sccMode, readWriteMode string) *sccModeEndpointFixture {
	t.Helper()
	fixture := &sccModeEndpointFixture{
		endpointID:    endpointID,
		sccMode:       sccMode,
		readWriteMode: readWriteMode,
	}
	fixture.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "DescribeDBClusterEndpoints":
			primary := map[string]interface{}{"DBEndpointId": "pe-primary", "EndpointType": "Primary"}
			endpoint := map[string]interface{}{
				"DBEndpointId":    fixture.endpointID,
				"EndpointType":    "Custom",
				"ReadWriteMode":   fixture.readWriteMode,
				"SccMode":         fixture.sccMode,
				"EndpointConfig":  `{"ConsistLevel":"0"}`,
				"Nodes":           "",
				"AutoAddNewNodes": "",
				"AddressItems": []map[string]string{
					{"NetType": "Private", "ConnectionString": fixture.endpointID + ".rwlb.rds.aliyuncs.com", "Port": "3306"},
				},
			}
			var items []map[string]interface{}
			switch requested := r.Form.Get("DBEndpointId"); {
			case requested == fixture.endpointID:
				items = []map[string]interface{}{endpoint}
			case requested != "":
				items = []map[string]interface{}{primary}
			case fixture.created.Load():
				items = []map[string]interface{}{endpoint, primary}
			default:
				items = []map[string]interface{}{primary}
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"Items": items})
		case "DescribeDBClusterSSL":
			_, _ = w.Write([]byte(`{"Items":[]}`))
		case "CreateDBClusterEndpoint":
			fixture.mu.Lock()
			fixture.createCalls++
			fixture.createForm = copySccModeForm(r.Form)
			fixture.mu.Unlock()
			fixture.created.Store(true)
			_ = json.NewEncoder(w).Encode(map[string]string{"DBEndpointId": fixture.endpointID})
		case "ModifyDBClusterEndpoint":
			fixture.mu.Lock()
			fixture.modifyCalls++
			fixture.modifyForm = copySccModeForm(r.Form)
			fixture.mu.Unlock()
			_, _ = w.Write([]byte(`{}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(fixture.Server.Close)
	return fixture
}

func (f *sccModeEndpointFixture) requests() (int, map[string][]string, int, map[string][]string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.createCalls, f.createForm, f.modifyCalls, f.modifyForm
}

func copySccModeForm(form map[string][]string) map[string][]string {
	result := make(map[string][]string, len(form))
	for key, values := range form {
		result[key] = append([]string(nil), values...)
	}
	return result
}

func newSccModeTestClient(t *testing.T, serverURL string) *connectivity.AliyunClient {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoint := strings.TrimPrefix(serverURL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := connectivity.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("polardb", endpoint)
	return client
}

func assertSccModeParam(t *testing.T, action string, form map[string][]string, want string, configured bool) {
	t.Helper()
	got, present := form["SccMode"]
	if !configured {
		if present {
			t.Fatalf("%s sent SccMode %v without configuration", action, got)
		}
		return
	}
	if !present {
		t.Fatalf("%s did not send SccMode, want %q", action, want)
	}
	if len(got) != 1 || got[0] != want {
		t.Fatalf("%s SccMode = %v, want %q", action, got, want)
	}
}

func TestUnitPolarDBEndpointSccModeCreateRequest(t *testing.T) {
	res := polardb.ResourceEndpoint()
	for _, tc := range []struct {
		name       string
		sccMode    string
		configured bool
	}{
		{name: "configured_on", sccMode: "on", configured: true},
		{name: "configured_off", sccMode: "off", configured: true},
		{name: "unconfigured"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newSccModeEndpointFixture(t, "pe-new", tc.sccMode, "ReadWrite")
			client := newSccModeTestClient(t, fixture.URL)
			raw := map[string]interface{}{
				"db_cluster_id":   "pc-unit-test",
				"endpoint_type":   "Custom",
				"read_write_mode": "ReadWrite",
				"endpoint_config": map[string]interface{}{"ConsistLevel": "0"},
			}
			if tc.configured {
				raw["scc_mode"] = tc.sccMode
			}
			d := schema.TestResourceDataRaw(t, res.Schema, raw)
			d.MarkNewResource()
			if err := res.Create(d, client); err != nil {
				t.Fatalf("create failed: %s", err)
			}
			createCalls, createForm, modifyCalls, modifyForm := fixture.requests()
			if createCalls < 1 {
				t.Fatal("CreateDBClusterEndpoint was not called")
			}
			assertSccModeParam(t, "create", createForm, tc.sccMode, tc.configured)
			if tc.configured && modifyCalls < 1 {
				t.Fatal("ModifyDBClusterEndpoint was not called after the configured create")
			}
			if modifyCalls > 0 {
				assertSccModeParam(t, "modify", modifyForm, tc.sccMode, tc.configured)
			}
		})
	}
}

func TestUnitPolarDBEndpointSccModeUpdateRequest(t *testing.T) {
	for _, tc := range []struct {
		name          string
		priorSccMode  string
		nextSccMode   string
		priorRW       string
		nextRW        string
		serverSccMode string
		serverRW      string
		wantSccMode   string
	}{
		{name: "update_off", priorSccMode: "on", nextSccMode: "off", priorRW: "ReadWrite", nextRW: "ReadWrite", serverSccMode: "off", serverRW: "ReadWrite", wantSccMode: "off"},
		{name: "unchanged_scc_mode", priorSccMode: "on", nextSccMode: "on", priorRW: "ReadWrite", nextRW: "ReadOnly", serverSccMode: "on", serverRW: "ReadOnly"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newSccModeEndpointFixture(t, "pe-existing", tc.serverSccMode, tc.serverRW)
			fixture.created.Store(true)
			client := newSccModeTestClient(t, fixture.URL)
			res := polardb.ResourceEndpoint()
			prior := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
				"db_cluster_id":   "pc-unit-test",
				"endpoint_type":   "Custom",
				"read_write_mode": tc.priorRW,
				"scc_mode":        tc.priorSccMode,
				"endpoint_config": map[string]interface{}{"ConsistLevel": "0"},
			})
			prior.SetId("pc-unit-test:pe-existing")
			state := prior.State()
			if state == nil {
				t.Fatal("expected an initial state for the update scenario")
			}
			config := terraform.NewResourceConfigRaw(map[string]interface{}{
				"db_cluster_id":   "pc-unit-test",
				"endpoint_type":   "Custom",
				"read_write_mode": tc.nextRW,
				"scc_mode":        tc.nextSccMode,
				"endpoint_config": map[string]interface{}{"ConsistLevel": "0"},
			})
			diff, err := res.Diff(context.Background(), state, config, nil)
			if err != nil {
				t.Fatalf("diff failed: %s", err)
			}
			if diff == nil {
				t.Fatal("expected a diff for the update scenario")
			}
			if _, diags := res.Apply(context.Background(), state, diff, client); diags.HasError() {
				t.Fatalf("apply failed: %v", diags)
			}
			createCalls, _, modifyCalls, modifyForm := fixture.requests()
			if createCalls != 0 {
				t.Fatalf("CreateDBClusterEndpoint was called %d times during the update", createCalls)
			}
			if modifyCalls < 1 {
				t.Fatal("ModifyDBClusterEndpoint was not called")
			}
			got, present := modifyForm["SccMode"]
			if tc.wantSccMode == "" {
				if present {
					t.Fatalf("modify sent SccMode %v although it did not change", got)
				}
			} else if !present || len(got) != 1 || got[0] != tc.wantSccMode {
				t.Fatalf("modify SccMode = %v (present=%v), want %q", got, present, tc.wantSccMode)
			}
			if tc.name == "unchanged_scc_mode" {
				if rw := modifyForm["ReadWriteMode"]; len(rw) != 1 || rw[0] != tc.nextRW {
					t.Fatalf("modify ReadWriteMode = %v, want %q", rw, tc.nextRW)
				}
			}
		})
	}
}
