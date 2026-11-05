package polardb_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	sdkpolardb "github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var privateConnectionStringRegexp = "^[a-z-A-Z-0-9]+.rwlb.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com"

func TestAccAliCloudPolarDBEndpointConfigUpdate(t *testing.T) {
	var v *sdkpolardb.DBEndpoint
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBendpoint-%s", rand)
	connectionPrefix := fmt.Sprintf("tf-testacc-%s", tfacctest.RandStringFromCharSet(8, tfacctest.CharSetAlphaNum))
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
					"connection_prefix": connectionPrefix,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": connectionPrefix,
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
				Config: testAccConfig(map[string]interface{}{
					"endpoint_config": map[string]string{
						"ConsistLevel": "0",
					},
					"scc_mode": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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
					"db_cluster_id": "${alicloud_polardb_cluster.cluster.id}",
					"endpoint_type": "Custom",
					"endpoint_config": map[string]string{
						"ConsistLevel": "0",
					},
					"scc_mode": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_type":                "Custom",
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

func TestAccAliCloudPolarDBEndpointConfigUpdate_SslConnectionStringAndConnectionPrefix(t *testing.T) {
	var v *sdkpolardb.DBEndpoint
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBendpoint-%s", rand)
	connectionPrefix := fmt.Sprintf("tf-testacc-%s", tfacctest.RandStringFromCharSet(8, tfacctest.CharSetAlphaNum))
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
					"connection_prefix": connectionPrefix,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": connectionPrefix,
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

func TestUnitPolarDBEndpointSccConfig(t *testing.T) {
	sccModeSchema := polardb.ResourceEndpoint().Schema["scc_mode"]
	for _, input := range []string{"on", "off"} {
		if _, errors := sccModeSchema.ValidateFunc(input, "scc_mode"); len(errors) > 0 {
			t.Fatalf("expected %q to be valid, got: %v", input, errors)
		}
	}
	for _, input := range []string{"ON", "OFF", "invalid"} {
		if _, errors := sccModeSchema.ValidateFunc(input, "scc_mode"); len(errors) == 0 {
			t.Fatalf("expected %q to be rejected", input)
		}
	}
}

func TestUnitPolarDBEndpointSccModeWire(t *testing.T) {
	var mu sync.Mutex
	var modifyForm url.Values
	var describeCalls atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch action := r.Form.Get("Action"); action {
		case "ModifyDBClusterEndpoint":
			mu.Lock()
			modifyForm = r.Form
			mu.Unlock()
			_, _ = w.Write([]byte(`{"RequestId":"mock"}`))
		case "DescribeDBClusterEndpoints":
			describeCalls.Add(1)
			_, _ = w.Write([]byte(`{"RequestId":"mock","Items":[{"DBEndpointId":"pe-unit","EndpointType":"Custom","SccMode":"ON","EndpointConfig":"{\"ConsistLevel\":\"0\"}","AutoAddNewNodes":"Disable","ReadWriteMode":"ReadOnly","AddressItems":[{"NetType":"Private","Port":"3306","ConnectionString":"pe-unit.rwlb.rds.aliyuncs.com"}]}]}`))
		case "DescribeDBClusterSSL":
			_, _ = w.Write([]byte(`{"RequestId":"mock","Items":[]}`))
		default:
			t.Errorf("unexpected action: %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
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
	endpoints.Store("polardb", endpoint)

	resourceSchema := polardb.ResourceEndpoint()
	ctx := context.Background()
	state := &terraform.InstanceState{
		ID: "pc-unit:pe-unit",
		Attributes: map[string]string{
			"id":            "pc-unit:pe-unit",
			"db_cluster_id": "pc-unit",
			"endpoint_type": "Custom",
			"scc_mode":      "off",
		},
	}
	configData := terraform.NewResourceConfigRaw(map[string]interface{}{
		"db_cluster_id": "pc-unit",
		"endpoint_type": "Custom",
		"endpoint_config": map[string]interface{}{
			"ConsistLevel": "0",
		},
		"scc_mode": "on",
	})
	diff, err := resourceSchema.Diff(ctx, state, configData, client)
	if err != nil {
		t.Fatalf("diff failed: %s", err)
	}
	newState, diags := resourceSchema.Apply(ctx, state, diff, client)
	if diags.HasError() {
		t.Fatalf("apply failed: %v", diags)
	}

	mu.Lock()
	form := modifyForm
	mu.Unlock()
	if form == nil {
		t.Fatal("ModifyDBClusterEndpoint was not called")
	}
	if got := form.Get("DBClusterId"); got != "pc-unit" {
		t.Errorf("DBClusterId = %q, want %q", got, "pc-unit")
	}
	if got := form.Get("DBEndpointId"); got != "pe-unit" {
		t.Errorf("DBEndpointId = %q, want %q", got, "pe-unit")
	}
	if got := form.Get("SccMode"); got != "on" {
		t.Errorf("SccMode request parameter = %q, want %q", got, "on")
	}
	var endpointConfig map[string]string
	if err := json.Unmarshal([]byte(form.Get("EndpointConfig")), &endpointConfig); err != nil {
		t.Fatalf("decode EndpointConfig request parameter %q: %s", form.Get("EndpointConfig"), err)
	}
	if endpointConfig["ConsistLevel"] != "0" {
		t.Errorf("EndpointConfig.ConsistLevel = %q, want %q", endpointConfig["ConsistLevel"], "0")
	}
	if got := newState.Attributes["scc_mode"]; got != "on" {
		t.Errorf("scc_mode after apply = %q, want %q for Describe value ON", got, "on")
	}
	if describeCalls.Load() == 0 {
		t.Error("DescribeDBClusterEndpoints was not called")
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
