// collection-order: bypass defined_column -- pre-existing attribute outside this change
// collection-order: bypass primary_key -- pre-existing attribute outside this change
package ots_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
	"testing"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudOtsTable_basic(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_table.default"
	ra := acctest.ResourceAttrInit(resourceId, otsTableBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTableConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${var.name}",
					"primary_key": []map[string]interface{}{
						{
							"name": "pk1",
							"type": "Integer",
						},
					},
					"time_to_live": "-1",
					"max_version":  "1",
					"allow_update": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-" + name,
						"table_name":    name,
						"allow_update":  "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"defined_column": []map[string]interface{}{
						{
							"name": "col1",
							"type": "Integer",
						},
						{
							"name": "col2",
							"type": "String",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"defined_column.#":      "2",
						"defined_column.0.name": "col1",
						"defined_column.0.type": "Integer",
						"defined_column.1.name": "col2",
						"defined_column.1.type": "String",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"deviation_cell_version_in_sec": "86401",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deviation_cell_version_in_sec": "86401",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_to_live": "86401",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_to_live": "86401",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_version": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_version": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_update": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_update": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_to_live":                  "-1",
					"max_version":                   "1",
					"deviation_cell_version_in_sec": "86400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_to_live":                  "-1",
						"max_version":                   "1",
						"deviation_cell_version_in_sec": "86400",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOtsTable_stream(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_table.default"
	ra := acctest.ResourceAttrInit(resourceId, otsTableBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTableConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${var.name}",
					"primary_key": []map[string]interface{}{
						{
							"name": "pk1",
							"type": "Integer",
						},
					},
					"time_to_live": "-1",
					"max_version":  "1",
					"stream_specification": []map[string]interface{}{
						{
							"enable_stream":   true,
							"expiration_time": 24,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stream_specification.#":                 "1",
						"stream_specification.0.enable_stream":   "true",
						"stream_specification.0.expiration_time": "24",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stream_specification": []map[string]interface{}{
						{
							"enable_stream":   true,
							"expiration_time": 168,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stream_specification.0.expiration_time": "168",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stream_specification": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stream_specification.0.enable_stream":   acctest.REMOVEKEY,
						"stream_specification.0.expiration_time": acctest.REMOVEKEY,
						"stream_specification.#":                 "0",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stream_specification": []map[string]interface{}{
						{
							"enable_stream":   true,
							"expiration_time": 12,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stream_specification.#":                 "1",
						"stream_specification.0.enable_stream":   "true",
						"stream_specification.0.expiration_time": "12",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stream_specification": []map[string]interface{}{
						{
							"enable_stream": true,
						},
					},
				}),
				ExpectError: regexp.MustCompile("must be set to a value between 1 and 168 when `enable_stream` is true"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stream_specification": []map[string]interface{}{
						{
							"enable_stream":   true,
							"expiration_time": 169,
						},
					},
				}),
				ExpectError: regexp.MustCompile(`expected stream_specification\.0\.expiration_time to be in the range \(1 - 168\), got 169`),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stream_specification": []map[string]interface{}{
						{
							"enable_stream":   true,
							"expiration_time": 0,
						},
					},
				}),
				ExpectError: regexp.MustCompile(`expected stream_specification\.0\.expiration_time to be in the range \(1 - 168\), got 0`),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stream_specification": []map[string]interface{}{
						{
							"enable_stream": false,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stream_specification.0.enable_stream":   acctest.REMOVEKEY,
						"stream_specification.0.expiration_time": acctest.REMOVEKEY,
						"stream_specification.#":                 "0",
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

func TestAccAliCloudOtsTable_highPerformance(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_table.default"
	ra := acctest.ResourceAttrInit(resourceId, otsTableBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTableConfigDependenceHighperformance)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${var.name}",
					"primary_key": []map[string]interface{}{
						{
							"name": "pk1",
							"type": "Integer",
						},
					},
					"defined_column": []map[string]interface{}{
						{
							"name": "col1",
							"type": "Integer",
						},
					},
					"time_to_live": "-1",
					"max_version":  "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-" + name,
						"table_name":    name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deviation_cell_version_in_sec": "86401",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deviation_cell_version_in_sec": "86401",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_to_live": "86401",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_to_live": "86401",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_version": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_version": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_to_live":                  "-1",
					"max_version":                   "1",
					"deviation_cell_version_in_sec": "86400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_to_live":                  "-1",
						"max_version":                   "1",
						"deviation_cell_version_in_sec": "86400",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOtsTable_multi(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_table.default.4"
	ra := acctest.ResourceAttrInit(resourceId, otsTableBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTableConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${var.name}${count.index}",
					"primary_key": []map[string]interface{}{
						{
							"name": "pk1",
							"type": "Integer",
						},
					},
					"defined_column": []map[string]interface{}{
						{
							"name": "col1",
							"type": "Integer",
						},
					},
					"time_to_live": "-1",
					"max_version":  "1",
					"count":        "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAliCloudOtsTable_withTableEncryption(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_table.default"
	ra := acctest.ResourceAttrInit(resourceId, otsTableBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTableConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${var.name}",
					"primary_key": []map[string]interface{}{
						{
							"name": "pk1",
							"type": "Integer",
						},
					},
					"defined_column": []map[string]interface{}{
						{
							"name": "col1",
							"type": "Integer",
						},
					},
					"time_to_live": "-1",
					"max_version":  "1",
					"enable_sse":   "true",
					"sse_key_type": "SSE_KMS_SERVICE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-" + name,
						"table_name":    name,
						"time_to_live":  "-1",
						"max_version":   "1",
						"enable_sse":    "true",
						"sse_key_type":  "SSE_KMS_SERVICE",
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

func TestAccAliCloudOtsTable_withTableByOKEncryption(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_table.default"
	ra := acctest.ResourceAttrInit(resourceId, otsTableBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTableConfigDependenceSseByOK)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${var.name}",
					"primary_key": []map[string]interface{}{
						{
							"name": "pk1",
							"type": "Integer",
						},
					},
					"defined_column": []map[string]interface{}{
						{
							"name": "col1",
							"type": "Integer",
						},
					},
					"time_to_live": "-1",
					"max_version":  "1",
					"enable_sse":   "true",
					"sse_key_type": "SSE_BYOK",
					"sse_key_id":   "${alicloud_kms_key.default.id}",
					"sse_role_arn": "${alicloud_ram_role.default.arn}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-" + name,
						"table_name":    name,
						"time_to_live":  "-1",
						"max_version":   "1",
						"enable_sse":    "true",
						"sse_key_type":  "SSE_BYOK",
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

func resourceOtsTableConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	resource "alicloud_ots_instance" "default" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "%s"
	  tags = {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}
	`, name, string(ots.OtsCapacity))
}

func resourceOtsTableConfigDependenceSseByOK(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "default" {
		name        = "tf-${var.name}"
		document    = <<EOF
		  {
			"Statement": [
			  {
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": {
				  "Service": [
					  "ots.aliyuncs.com"
	  
				  ]
				}
			  }
			],
			"Version": "1"
		  }
		  EOF
		description = "this is TfOtsAccessKMS role."
	  }
	  resource "alicloud_ram_policy" "default" {
		policy_name     = "tf-${var.name}"
		policy_document = <<EOF
		{
		  "Statement": [
			{
			  "Effect": "Allow",
			  "Action": [
					  "kms:Decrypt",
					  "kms:GenerateDataKey"
				  ],
				  "Resource": [
					  "*"
				  ]
			}
		  ],
			"Version": "1"
		}
		EOF
		description     = "this is a policy"
	  }
	  
	  resource "alicloud_ram_role_policy_attachment" "default" {
		policy_name = alicloud_ram_policy.default.policy_name
		policy_type = alicloud_ram_policy.default.type
		role_name   = alicloud_ram_role.default.name
	  }
	  
	  resource "alicloud_kms_key" "default" {
		description            = "Hello KMS"
		pending_window_in_days = "7"
		status                 = "Enabled"
	  }


	variable "name" {
	  default = "%s"
	}
	resource "alicloud_ots_instance" "default" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "%s"
	  tags = {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}
	`, name, string(ots.OtsCapacity))
}

func resourceOtsTableConfigDependenceHighperformance(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	resource "alicloud_ots_instance" "default" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "%s"
	  tags = {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}
	`, name, string(ots.OtsHighPerformance))
}

var otsTableBasicMap = map[string]string{
	"primary_key.#":                 "1",
	"primary_key.0.name":            "pk1",
	"primary_key.0.type":            "Integer",
	"time_to_live":                  "-1",
	"max_version":                   "1",
	"deviation_cell_version_in_sec": "86400",
}

func testAccCheckOtsTableExist(n string, table *tablestore.DescribeTableResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found OTS table: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no OTS table ID is set")
		}

		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		otsService := ots.NewService(client)

		response, err := otsService.DescribeOtsTable(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error finding OTS table %s: %#v", rs.Primary.ID, err)
		}

		table = response
		return nil
	}
}

func testAccCheckOtsTableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ots_table" {
			continue
		}

		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		otsService := ots.NewService(client)

		if _, err := otsService.DescribeOtsTable(rs.Primary.ID); err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("error! Ots table still exists")
	}

	return nil
}

const (
	otsUnitInstanceName = "terraform-ex"
	otsUnitTableName    = "terraform_example"
	otsUnitResourceId   = "terraform-ex:terraform_example"
)

type otsUnitCall struct {
	path string
	body []byte
}

type otsUnitStream struct {
	enabled    bool
	expiration int32
}

type otsUnitFailure struct {
	path    string
	code    string
	message string
	times   int
}

// otsUnitGateway replays the OTS REST and Tablestore endpoints so the resource entry points can be
// exercised without a cloud account and every serialized request can be asserted.
type otsUnitGateway struct {
	mu       sync.Mutex
	stream   otsUnitStream
	calls    []otsUnitCall
	failures []*otsUnitFailure
}

func (g *otsUnitGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	g.mu.Lock()
	g.calls = append(g.calls, otsUnitCall{path: r.URL.Path, body: body})
	stream := g.stream
	failure := g.takeFailureLocked(r.URL.Path)
	g.mu.Unlock()

	if failure != nil {
		writeOtsUnitError(w, failure.code, failure.message)
		return
	}

	switch r.URL.Path {
	case "/v2/openapi/getinstance":
		payload, err := json.Marshal(map[string]interface{}{
			"InstanceName": r.URL.Query().Get("InstanceName"),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	case "/CreateTable":
		w.WriteHeader(http.StatusOK)
	case "/DescribeTable":
		request := new(otsprotocol.DescribeTableRequest)
		if err := proto.Unmarshal(body, request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		writeOtsUnitProto(w, otsUnitDescribeResponse(request.GetTableName(), stream))
	case "/UpdateTable":
		request := new(otsprotocol.UpdateTableRequest)
		if err := proto.Unmarshal(body, request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if spec := request.GetStreamSpec(); spec != nil {
			stream.enabled = spec.GetEnableStream()
			if stream.enabled {
				stream.expiration = spec.GetExpirationTime()
			}
			g.mu.Lock()
			g.stream = stream
			g.mu.Unlock()
		}
		writeOtsUnitProto(w, otsUnitUpdateResponse(stream))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (g *otsUnitGateway) takeFailureLocked(path string) *otsUnitFailure {
	for _, failure := range g.failures {
		if failure.path == path && failure.times > 0 {
			failure.times--
			return failure
		}
	}
	return nil
}

func (g *otsUnitGateway) bodies(path string) [][]byte {
	g.mu.Lock()
	defer g.mu.Unlock()
	var result [][]byte
	for _, call := range g.calls {
		if call.path == path {
			result = append(result, call.body)
		}
	}
	return result
}

func (g *otsUnitGateway) count(path string) int {
	return len(g.bodies(path))
}

func (g *otsUnitGateway) paths() []string {
	g.mu.Lock()
	defer g.mu.Unlock()
	paths := make([]string, 0, len(g.calls))
	for _, call := range g.calls {
		paths = append(paths, call.path)
	}
	return paths
}

func writeOtsUnitError(w http.ResponseWriter, code, message string) {
	payload, err := proto.Marshal(&otsprotocol.Error{Code: proto.String(code), Message: proto.String(message)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write(payload)
}

func writeOtsUnitProto(w http.ResponseWriter, message proto.Message) {
	payload, err := proto.Marshal(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}

func otsUnitTableOptions() *otsprotocol.TableOptions {
	return &otsprotocol.TableOptions{
		TimeToLive:                proto.Int32(-1),
		MaxVersions:               proto.Int32(1),
		DeviationCellVersionInSec: proto.Int64(86400),
	}
}

func otsUnitReservedThroughput() *otsprotocol.ReservedThroughputDetails {
	return &otsprotocol.ReservedThroughputDetails{
		CapacityUnit:     &otsprotocol.CapacityUnit{Read: proto.Int32(0), Write: proto.Int32(0)},
		LastIncreaseTime: proto.Int64(0),
	}
}

func otsUnitStreamDetails(stream otsUnitStream) *otsprotocol.StreamDetails {
	details := &otsprotocol.StreamDetails{EnableStream: proto.Bool(stream.enabled)}
	if stream.enabled {
		details.StreamId = proto.String("unit-stream-id")
		details.ExpirationTime = proto.Int32(stream.expiration)
		details.LastEnableTime = proto.Int64(1700000000)
	}
	return details
}

func otsUnitDescribeResponse(tableName string, stream otsUnitStream) *otsprotocol.DescribeTableResponse {
	primaryKeyType := otsprotocol.PrimaryKeyType_INTEGER
	return &otsprotocol.DescribeTableResponse{
		TableMeta: &otsprotocol.TableMeta{
			TableName:  proto.String(tableName),
			PrimaryKey: []*otsprotocol.PrimaryKeySchema{{Name: proto.String("pk1"), Type: &primaryKeyType}},
		},
		ReservedThroughputDetails: otsUnitReservedThroughput(),
		TableOptions:              otsUnitTableOptions(),
		TableStatus:               otsprotocol.TableStatus_ACTIVE.Enum(),
		StreamDetails:             otsUnitStreamDetails(stream),
	}
}

func otsUnitUpdateResponse(stream otsUnitStream) *otsprotocol.UpdateTableResponse {
	return &otsprotocol.UpdateTableResponse{
		ReservedThroughputDetails: otsUnitReservedThroughput(),
		TableOptions:              otsUnitTableOptions(),
		StreamDetails:             otsUnitStreamDetails(stream),
	}
}

func newOtsUnitClient(t *testing.T, server *httptest.Server) *connectivity.AliyunClient {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("unit-access-key").
		SetAccessKeySecret("unit-secret-key"))
	if err != nil {
		t.Fatalf("create unit credential: %s", err)
	}
	endpoints := &sync.Map{}
	// The ROA client appends the path to this value directly, so it must stay a bare host.
	endpoints.Store("tablestore", server.Listener.Addr().String())
	config := &connectivity.Config{
		Region:               connectivity.Hangzhou,
		RegionId:             "cn-hangzhou",
		AccessKey:            "unit-access-key",
		SecretKey:            "unit-secret-key",
		AccountType:          "Domestic",
		Protocol:             "HTTP",
		Endpoints:            endpoints,
		OtsEndpoint:          server.URL,
		SkipRegionValidation: true,
		Credential:           credential,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("create unit client: %s", err)
	}
	return client
}

func otsUnitRaw(streamSpec []interface{}) map[string]interface{} {
	raw := map[string]interface{}{
		"instance_name": otsUnitInstanceName,
		"table_name":    otsUnitTableName,
		"primary_key": []interface{}{
			map[string]interface{}{"name": "pk1", "type": "Integer"},
		},
		"time_to_live": -1,
		"max_version":  1,
	}
	if streamSpec != nil {
		raw["stream_specification"] = streamSpec
	}
	return raw
}

func otsUnitStreamConfig(enableStream bool, expirationTime int) []interface{} {
	spec := map[string]interface{}{"enable_stream": enableStream}
	if enableStream {
		spec["expiration_time"] = expirationTime
	}
	return []interface{}{spec}
}

func otsUnitStateStream(t *testing.T, d *schema.ResourceData) otsUnitStream {
	t.Helper()
	specs, _ := d.Get("stream_specification").([]interface{})
	if len(specs) == 0 || specs[0] == nil {
		return otsUnitStream{}
	}
	spec := specs[0].(map[string]interface{})
	stream := otsUnitStream{}
	stream.enabled, _ = spec["enable_stream"].(bool)
	expirationTime, _ := spec["expiration_time"].(int)
	stream.expiration = int32(expirationTime)
	return stream
}

func TestUnitAliCloudOtsTableStreamCreate(t *testing.T) {
	t.Run("enabled stream is sent and response is read back", func(t *testing.T) {
		gateway := &otsUnitGateway{stream: otsUnitStream{enabled: true, expiration: 48}}
		server := httptest.NewServer(gateway)
		defer server.Close()

		d := schema.TestResourceDataRaw(t, ots.ResourceOtsTable().Schema, otsUnitRaw(otsUnitStreamConfig(true, 24)))
		if err := ots.ResourceOtsTable().Create(d, newOtsUnitClient(t, server)); err != nil {
			t.Fatalf("Create with stream: %s (requests: %v)", err, gateway.paths())
		}

		payload := otsUnitCreateTablePayload(t, gateway)
		if payload.GetStreamSpec() == nil {
			t.Fatalf("CreateTable payload has no stream_spec")
		}
		if !payload.GetStreamSpec().GetEnableStream() {
			t.Fatalf("CreateTable payload enable_stream = false, want true")
		}
		if got := payload.GetStreamSpec().GetExpirationTime(); got != 24 {
			t.Fatalf("CreateTable payload expiration_time = %d, want 24", got)
		}

		stream := otsUnitStateStream(t, d)
		if !stream.enabled {
			t.Fatalf("state enable_stream = false, want true")
		}
		if stream.expiration != 48 {
			t.Fatalf("state expiration_time = %d, want the DescribeTable value 48", stream.expiration)
		}
	})

	t.Run("absent block sends no stream specification", func(t *testing.T) {
		gateway := &otsUnitGateway{}
		server := httptest.NewServer(gateway)
		defer server.Close()

		d := schema.TestResourceDataRaw(t, ots.ResourceOtsTable().Schema, otsUnitRaw(nil))
		if err := ots.ResourceOtsTable().Create(d, newOtsUnitClient(t, server)); err != nil {
			t.Fatalf("Create without stream: %s (requests: %v)", err, gateway.paths())
		}

		payload := otsUnitCreateTablePayload(t, gateway)
		if payload.GetStreamSpec() != nil {
			t.Fatalf("CreateTable payload stream_spec = %#v, want omitted", payload.GetStreamSpec())
		}
	})

	t.Run("explicit false block sends no stream specification", func(t *testing.T) {
		gateway := &otsUnitGateway{}
		server := httptest.NewServer(gateway)
		defer server.Close()

		d := schema.TestResourceDataRaw(t, ots.ResourceOtsTable().Schema, otsUnitRaw(otsUnitStreamConfig(false, 0)))
		if err := ots.ResourceOtsTable().Create(d, newOtsUnitClient(t, server)); err != nil {
			t.Fatalf("Create with disabled stream: %s (requests: %v)", err, gateway.paths())
		}

		payload := otsUnitCreateTablePayload(t, gateway)
		if payload.GetStreamSpec() != nil {
			t.Fatalf("CreateTable payload stream_spec = %#v, want omitted", payload.GetStreamSpec())
		}
	})
}

func otsUnitCreateTablePayload(t *testing.T, gateway *otsUnitGateway) *otsprotocol.CreateTableRequest {
	t.Helper()
	bodies := gateway.bodies("/CreateTable")
	if len(bodies) != 1 {
		t.Fatalf("CreateTable calls = %d, want exactly 1 (requests: %v)", len(bodies), gateway.paths())
	}
	payload := new(otsprotocol.CreateTableRequest)
	if err := proto.Unmarshal(bodies[0], payload); err != nil {
		t.Fatalf("decode CreateTable payload: %s", err)
	}
	return payload
}

func otsUnitUpdateTablePayloads(t *testing.T, gateway *otsUnitGateway) []*otsprotocol.UpdateTableRequest {
	t.Helper()
	bodies := gateway.bodies("/UpdateTable")
	if len(bodies) == 0 {
		t.Fatalf("UpdateTable was not called (requests: %v)", gateway.paths())
	}
	payloads := make([]*otsprotocol.UpdateTableRequest, 0, len(bodies))
	for i, body := range bodies {
		payload := new(otsprotocol.UpdateTableRequest)
		if err := proto.Unmarshal(body, payload); err != nil {
			t.Fatalf("decode UpdateTable payload %d: %s", i, err)
		}
		payloads = append(payloads, payload)
	}
	return payloads
}

func TestUnitAliCloudOtsTableStreamUpdate(t *testing.T) {
	unitResource := ots.ResourceOtsTable()
	schemaMap := schema.InternalMap(unitResource.Schema)

	update := func(t *testing.T, gateway *otsUnitGateway, server *httptest.Server, oldStream, newSpec []interface{}) (*schema.ResourceData, error) {
		t.Helper()
		oldData := schema.TestResourceDataRaw(t, unitResource.Schema, otsUnitRaw(oldStream))
		oldData.SetId(otsUnitResourceId)
		oldState := oldData.State()

		diff, err := schemaMap.Diff(context.Background(), oldState, sdkterraform.NewResourceConfigRaw(otsUnitRaw(newSpec)), unitResource.CustomizeDiff, nil, false)
		if err != nil {
			t.Fatalf("plan diff: %s", err)
		}
		d, err := schemaMap.Data(oldState, diff)
		if err != nil {
			t.Fatalf("build resource data: %s", err)
		}
		d.SetId(otsUnitResourceId)
		return d, unitResource.Update(d, newOtsUnitClient(t, server))
	}

	t.Run("retention update sends the new value only", func(t *testing.T) {
		gateway := &otsUnitGateway{stream: otsUnitStream{enabled: true, expiration: 24}}
		server := httptest.NewServer(gateway)
		defer server.Close()

		d, err := update(t, gateway, server, otsUnitStreamConfig(true, 24), otsUnitStreamConfig(true, 168))
		if err != nil {
			t.Fatalf("Update: %s (requests: %v)", err, gateway.paths())
		}

		payloads := otsUnitUpdateTablePayloads(t, gateway)
		if len(payloads) != 1 {
			t.Fatalf("UpdateTable calls = %d, want exactly 1", len(payloads))
		}
		for i, payload := range payloads {
			if payload.GetTableOptions() != nil {
				t.Fatalf("UpdateTable call %d must not send table_options", i)
			}
		}
		payload := payloads[len(payloads)-1]
		if payload.GetStreamSpec() == nil {
			t.Fatalf("UpdateTable payload has no stream_spec")
		}
		if !payload.GetStreamSpec().GetEnableStream() {
			t.Fatalf("UpdateTable payload enable_stream = false, want true")
		}
		if got := payload.GetStreamSpec().GetExpirationTime(); got != 168 {
			t.Fatalf("UpdateTable payload expiration_time = %d, want 168", got)
		}

		stream := otsUnitStateStream(t, d)
		if !stream.enabled || stream.expiration != 168 {
			t.Fatalf("state stream = %#v, want enabled with 168", stream)
		}
	})

	t.Run("transient quota exhaustion is retried", func(t *testing.T) {
		gateway := &otsUnitGateway{
			stream: otsUnitStream{enabled: true, expiration: 24},
			failures: []*otsUnitFailure{{
				path:    "/UpdateTable",
				code:    "OTSQuotaExhausted",
				message: "Table operations are too frequent.",
				times:   1,
			}},
		}
		server := httptest.NewServer(gateway)
		defer server.Close()

		d, err := update(t, gateway, server, otsUnitStreamConfig(true, 24), otsUnitStreamConfig(true, 168))
		if err != nil {
			t.Fatalf("Update with a transient quota error: %s (requests: %v)", err, gateway.paths())
		}
		if calls := gateway.count("/UpdateTable"); calls != 2 {
			t.Fatalf("UpdateTable calls = %d, want 2 so the transient error is retried", calls)
		}
		if stream := otsUnitStateStream(t, d); !stream.enabled || stream.expiration != 168 {
			t.Fatalf("state stream = %#v, want enabled with 168", stream)
		}
	})

	t.Run("permanent quota exhaustion fails fast", func(t *testing.T) {
		gateway := &otsUnitGateway{
			stream: otsUnitStream{enabled: true, expiration: 24},
			failures: []*otsUnitFailure{{
				path:    "/UpdateTable",
				code:    "OTSQuotaExhausted",
				message: "Number of tables exceeds the upper limit.",
				times:   1,
			}},
		}
		server := httptest.NewServer(gateway)
		defer server.Close()

		if _, err := update(t, gateway, server, otsUnitStreamConfig(true, 24), otsUnitStreamConfig(true, 168)); err == nil {
			t.Fatalf("Update with a permanent quota error succeeded, want failure")
		}
		if calls := gateway.count("/UpdateTable"); calls != 1 {
			t.Fatalf("UpdateTable calls = %d, want 1 so the permanent error is not retried", calls)
		}
	})

	t.Run("removing the block disables stream", func(t *testing.T) {
		gateway := &otsUnitGateway{stream: otsUnitStream{enabled: true, expiration: 12}}
		server := httptest.NewServer(gateway)
		defer server.Close()

		d, err := update(t, gateway, server, otsUnitStreamConfig(true, 12), nil)
		if err != nil {
			t.Fatalf("Update: %s (requests: %v)", err, gateway.paths())
		}

		payloads := otsUnitUpdateTablePayloads(t, gateway)
		for i, payload := range payloads {
			if payload.GetTableOptions() != nil {
				t.Fatalf("UpdateTable call %d must not send table_options", i)
			}
			if payload.GetStreamSpec() == nil {
				t.Fatalf("UpdateTable call %d has no stream_spec", i)
			}
			if payload.GetStreamSpec().GetEnableStream() {
				t.Fatalf("UpdateTable call %d enable_stream = true, want false", i)
			}
			if payload.GetStreamSpec().ExpirationTime != nil {
				t.Fatalf("UpdateTable call %d expiration_time = %d, want omitted", i, payload.GetStreamSpec().GetExpirationTime())
			}
		}

		if stream := otsUnitStateStream(t, d); stream.enabled {
			t.Fatalf("state stream = %#v, want disabled", stream)
		}
	})
}

func TestUnitAliCloudOtsTableStreamCustomizeDiff(t *testing.T) {
	unitResource := ots.ResourceOtsTable()
	schemaMap := schema.InternalMap(unitResource.Schema)
	streamType := cty.Object(map[string]cty.Type{"enable_stream": cty.Bool, "expiration_time": cty.Number})

	cases := []struct {
		name      string
		block     cty.Value
		spec      []interface{}
		wantError bool
	}{
		{
			name:  "known valid retention passes",
			block: cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"enable_stream": cty.True, "expiration_time": cty.NumberIntVal(24)})}),
			spec:  otsUnitStreamConfig(true, 24),
		},
		{
			name:      "known zero retention is rejected",
			block:     cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"enable_stream": cty.True, "expiration_time": cty.NumberIntVal(0)})}),
			spec:      otsUnitStreamConfig(true, 0),
			wantError: true,
		},
		{
			name:      "known retention above the limit is rejected",
			block:     cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"enable_stream": cty.True, "expiration_time": cty.NumberIntVal(169)})}),
			spec:      otsUnitStreamConfig(true, 169),
			wantError: true,
		},
		{
			name:  "unknown retention is deferred",
			block: cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"enable_stream": cty.True, "expiration_time": cty.UnknownVal(cty.Number)})}),
			spec:  otsUnitStreamConfig(true, 0),
		},
		{
			name:  "absent block passes",
			block: cty.NullVal(cty.List(streamType)),
			spec:  nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			state := &sdkterraform.InstanceState{
				ID:        otsUnitResourceId,
				RawConfig: cty.ObjectVal(map[string]cty.Value{"stream_specification": tc.block}),
			}
			_, err := schemaMap.Diff(context.Background(), state, sdkterraform.NewResourceConfigRaw(otsUnitRaw(tc.spec)), unitResource.CustomizeDiff, nil, false)
			if tc.wantError && err == nil {
				t.Fatalf("expected a plan-time rejection")
			}
			if !tc.wantError && err != nil {
				t.Fatalf("unexpected plan-time rejection: %s", err)
			}
		})
	}
}
