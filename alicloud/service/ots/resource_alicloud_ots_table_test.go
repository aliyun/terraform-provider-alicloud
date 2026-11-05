// collection-order: bypass defined_column -- pre-existing attribute outside this change
// collection-order: bypass primary_key -- pre-existing attribute outside this change
package ots_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
	"github.com/golang/protobuf/proto"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// collection-order: bypass defined_column -- pre-existing attribute outside this change
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

func TestAccAliCloudOtsTable_withAutoIncrementPrimaryKey(t *testing.T) {
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
						{
							"name":           "pk2",
							"type":           "Integer",
							"auto_increment": "true",
						},
					},
					"time_to_live": "-1",
					"max_version":  "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"primary_key.#":                 "2",
						"primary_key.0.name":            "pk1",
						"primary_key.0.type":            "Integer",
						"primary_key.0.auto_increment":  "false",
						"primary_key.1.name":            "pk2",
						"primary_key.1.type":            "Integer",
						"primary_key.1.auto_increment":  "true",
						"time_to_live":                  "-1",
						"max_version":                   "1",
						"deviation_cell_version_in_sec": "86400",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_local_txn"},
			},
		},
	})
}

func TestAccAliCloudOtsTable_withLocalTxn(t *testing.T) {
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
					"enable_local_txn": "true",
					"time_to_live":     "-1",
					"max_version":      "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_local_txn":              "true",
						"primary_key.#":                 "1",
						"primary_key.0.name":            "pk1",
						"primary_key.0.type":            "Integer",
						"primary_key.0.auto_increment":  "false",
						"time_to_live":                  "-1",
						"max_version":                   "1",
						"deviation_cell_version_in_sec": "86400",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_local_txn"},
			},
		},
	})
}

func TestAccAliCloudOtsTable_withAllowUpdate(t *testing.T) {
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
					"allow_update": "true",
					"time_to_live": "-1",
					"max_version":  "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_update":                  "true",
						"primary_key.#":                 "1",
						"primary_key.0.name":            "pk1",
						"primary_key.0.type":            "Integer",
						"primary_key.0.auto_increment":  "false",
						"time_to_live":                  "-1",
						"max_version":                   "1",
						"deviation_cell_version_in_sec": "86400",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_local_txn"},
			},
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
					"allow_update": "false",
					"time_to_live": "-1",
					"max_version":  "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_update":                  "false",
						"primary_key.#":                 "1",
						"primary_key.0.name":            "pk1",
						"primary_key.0.type":            "Integer",
						"primary_key.0.auto_increment":  "false",
						"time_to_live":                  "-1",
						"max_version":                   "1",
						"deviation_cell_version_in_sec": "86400",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_local_txn"},
			},
		},
	})
}

// TestAccAliCloudOtsTable_primaryKeyOrder exercises the convergent reorder coverage
// required for the configurable multi-member primary_key TypeList: apply two distinct
// primary keys in one order, then apply the same members reordered. primary_key is
// ForceNew, so the second step replaces the table; the OTS CreateTable API echoes the
// submitted primary-key order (the first key is the partition key), so the post-apply
// plan converges to empty. See scripts/collection-order/README.md.
func TestAccAliCloudOtsTable_primaryKeyOrder(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_table.default"
	ra := acctest.ResourceAttrInit(resourceId, otsTableBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

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
						{
							"name": "pk2",
							"type": "String",
						},
					},
					"time_to_live": "-1",
					"max_version":  "1",
				}),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_local_txn"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${var.name}",
					"primary_key": []map[string]interface{}{
						{
							"name": "pk2",
							"type": "String",
						},
						{
							"name": "pk1",
							"type": "Integer",
						},
					},
					"time_to_live": "-1",
					"max_version":  "1",
				}),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_local_txn"},
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

// TestUnitOtsTableCustomizeDiffRejectsInvalidConfigs drives the plan-time
// CustomizeDiff through Resource.Diff so the four rejection rules are exercised
// without any cloud call: Diff never reaches Create, and a nil meta would panic
// if the callback touched the API client.
func TestUnitOtsTableCustomizeDiffRejectsInvalidConfigs(t *testing.T) {
	primaryKey := func(name, pkType string, autoIncrement bool) map[string]interface{} {
		return map[string]interface{}{
			"name":           name,
			"type":           pkType,
			"auto_increment": autoIncrement,
		}
	}
	tableConfig := func(keys []map[string]interface{}, localTxn bool) map[string]interface{} {
		return map[string]interface{}{
			"instance_name":    "tf-test-instance",
			"table_name":       "tf_test_table",
			"time_to_live":     -1,
			"max_version":      1,
			"primary_key":      keys,
			"enable_local_txn": localTxn,
		}
	}

	cases := []struct {
		name    string
		config  map[string]interface{}
		wantErr string
	}{
		{
			name:    "auto_increment primary key must be Integer",
			config:  tableConfig([]map[string]interface{}{primaryKey("pk1", "Integer", false), primaryKey("pk2", "String", true)}, false),
			wantErr: `the primary_key with auto_increment set to true must be of type Integer, got "String"`,
		},
		{
			name:    "auto_increment primary key must not be the partition key",
			config:  tableConfig([]map[string]interface{}{primaryKey("pk1", "Integer", true), primaryKey("pk2", "Integer", false)}, false),
			wantErr: "the auto_increment primary key must not be the partition key (the first primary_key)",
		},
		{
			name:    "at most one auto_increment primary key",
			config:  tableConfig([]map[string]interface{}{primaryKey("pk1", "Integer", false), primaryKey("pk2", "Integer", true), primaryKey("pk3", "Integer", true)}, false),
			wantErr: "at most one primary_key can be set to auto_increment, got 2",
		},
		{
			name:    "auto_increment conflicts with local transaction",
			config:  tableConfig([]map[string]interface{}{primaryKey("pk1", "Integer", false), primaryKey("pk2", "Integer", true)}, true),
			wantErr: "enable_local_txn cannot be set to true when any primary_key has auto_increment set to true",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := ots.ResourceOtsTable()
			if _, err := r.Diff(context.Background(), nil, sdktf.NewResourceConfigRaw(tc.config), nil); err == nil {
				t.Fatalf("expected the configuration to be rejected, got no error")
			} else if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("unexpected rejection reason: %v, want it to contain %q", err, tc.wantErr)
			}
		})
	}

	t.Run("valid auto_increment configuration is accepted", func(t *testing.T) {
		r := ots.ResourceOtsTable()
		config := tableConfig([]map[string]interface{}{primaryKey("pk1", "Integer", false), primaryKey("pk2", "Integer", true)}, false)
		if _, err := r.Diff(context.Background(), nil, sdktf.NewResourceConfigRaw(config), nil); err != nil {
			t.Fatalf("expected the configuration to be accepted, got error: %v", err)
		}
	})

	t.Run("valid local transaction configuration is accepted", func(t *testing.T) {
		r := ots.ResourceOtsTable()
		config := tableConfig([]map[string]interface{}{primaryKey("pk1", "Integer", false)}, true)
		if _, err := r.Diff(context.Background(), nil, sdktf.NewResourceConfigRaw(config), nil); err != nil {
			t.Fatalf("expected the configuration to be accepted, got error: %v", err)
		}
	})
}

type otsTableWireFixture struct {
	createRequests [][]byte
	describe       *otsprotocol.DescribeTableResponse
}

// newOtsTableWireClient points both the OTS REST API and the TableStore
// protobuf endpoint at a local server, so a full resource Create runs without
// touching the cloud and every serialized request can be inspected.
func newOtsTableWireClient(t *testing.T) (*connectivity.AliyunClient, *otsTableWireFixture) {
	t.Helper()
	fixture := &otsTableWireFixture{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v2/openapi/getinstance":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"body":{"InstanceName":"tf-test-instance","Status":"Normal","InstanceType":"Capacity"}}`)
		case "/CreateTable":
			raw, err := io.ReadAll(req.Body)
			if err != nil {
				t.Errorf("read CreateTable request: %v", err)
			}
			fixture.createRequests = append(fixture.createRequests, raw)
			w.WriteHeader(http.StatusOK)
		case "/DescribeTable":
			raw, err := proto.Marshal(fixture.describe)
			if err != nil {
				t.Errorf("marshal DescribeTable response: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(raw)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	endpoints := &sync.Map{}
	signVersions := &sync.Map{}
	endpoints.Store("tablestore", strings.TrimPrefix(server.URL, "http://"))
	config := &connectivity.Config{
		RegionId:             "cn-hangzhou",
		AccessKey:            "test-access-key",
		SecretKey:            "test-secret-key",
		Protocol:             "HTTP",
		AccountType:          "test",
		Endpoints:            endpoints,
		SignVersion:          signVersions,
		OtsEndpoint:          server.URL,
		SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("build connectivity client: %v", err)
	}
	return client, fixture
}

func otsTableDescribeResponse(secondAutoIncrement bool) *otsprotocol.DescribeTableResponse {
	timeToLive := int32(-1)
	maxVersion := int32(1)
	deviation := int64(86400)
	allowUpdate := true
	readCapacity := int32(0)
	writeCapacity := int32(0)
	status := otsprotocol.TableStatus_ACTIVE
	pkType := otsprotocol.PrimaryKeyType_INTEGER
	second := &otsprotocol.PrimaryKeySchema{Name: proto.String("pk2"), Type: &pkType}
	if secondAutoIncrement {
		option := otsprotocol.PrimaryKeyOption_AUTO_INCREMENT
		second.Option = &option
	}
	return &otsprotocol.DescribeTableResponse{
		TableMeta: &otsprotocol.TableMeta{
			TableName: proto.String("tf_test_table"),
			PrimaryKey: []*otsprotocol.PrimaryKeySchema{
				{Name: proto.String("pk1"), Type: &pkType},
				second,
			},
		},
		ReservedThroughputDetails: &otsprotocol.ReservedThroughputDetails{
			CapacityUnit:     &otsprotocol.CapacityUnit{Read: &readCapacity, Write: &writeCapacity},
			LastIncreaseTime: proto.Int64(0),
		},
		TableOptions: &otsprotocol.TableOptions{
			TimeToLive:                &timeToLive,
			MaxVersions:               &maxVersion,
			DeviationCellVersionInSec: &deviation,
			AllowUpdate:               &allowUpdate,
		},
		TableStatus: &status,
	}
}

// TestUnitOtsTableCreateAndDescribeWire captures the serialized CreateTable
// request and the DescribeTable response the Read path maps back, covering the
// two fields' plan/apply wire semantics without a cloud call.
func TestUnitOtsTableCreateAndDescribeWire(t *testing.T) {
	primaryKey := func(name string, autoIncrement bool) map[string]interface{} {
		return map[string]interface{}{"name": name, "type": "Integer", "auto_increment": autoIncrement}
	}
	tableConfig := func(localTxn bool, keys ...map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"instance_name":    "tf-test-instance",
			"table_name":       "tf_test_table",
			"time_to_live":     -1,
			"max_version":      1,
			"primary_key":      keys,
			"enable_local_txn": localTxn,
		}
	}

	t.Run("auto_increment is sent on CreateTable and read back from DescribeTable", func(t *testing.T) {
		client, fixture := newOtsTableWireClient(t)
		fixture.describe = otsTableDescribeResponse(true)

		r := ots.ResourceOtsTable()
		d := schema.TestResourceDataRaw(t, r.Schema, tableConfig(false, primaryKey("pk1", false), primaryKey("pk2", true)))
		if err := r.Create(d, client); err != nil {
			t.Fatalf("create table: %v", err)
		}
		if len(fixture.createRequests) != 1 {
			t.Fatalf("expected 1 CreateTable request, got %d", len(fixture.createRequests))
		}
		request := &otsprotocol.CreateTableRequest{}
		if err := proto.Unmarshal(fixture.createRequests[0], request); err != nil {
			t.Fatalf("decode CreateTable request: %v", err)
		}
		if len(request.TableMeta.PrimaryKey) != 2 {
			t.Fatalf("expected 2 primary keys in the request, got %d", len(request.TableMeta.PrimaryKey))
		}
		if request.TableMeta.PrimaryKey[0].Option != nil {
			t.Fatalf("partition key must not carry an option, got %v", *request.TableMeta.PrimaryKey[0].Option)
		}
		if request.TableMeta.PrimaryKey[1].Option == nil || *request.TableMeta.PrimaryKey[1].Option != otsprotocol.PrimaryKeyOption_AUTO_INCREMENT {
			t.Fatalf("second primary key option = %v, want AUTO_INCREMENT", request.TableMeta.PrimaryKey[1].Option)
		}
		if request.EnableLocalTxn == nil || *request.EnableLocalTxn {
			t.Fatalf("EnableLocalTxn = %v, want false", request.EnableLocalTxn)
		}

		keys := d.Get("primary_key").([]interface{})
		if len(keys) != 2 {
			t.Fatalf("expected 2 primary keys in state, got %d", len(keys))
		}
		if got := keys[0].(map[string]interface{})["auto_increment"]; got != false {
			t.Fatalf("pk1 auto_increment = %v, want false", got)
		}
		if got := keys[1].(map[string]interface{})["auto_increment"]; got != true {
			t.Fatalf("pk2 auto_increment = %v, want true", got)
		}
	})

	t.Run("enable_local_txn is sent on CreateTable", func(t *testing.T) {
		client, fixture := newOtsTableWireClient(t)
		fixture.describe = otsTableDescribeResponse(false)

		r := ots.ResourceOtsTable()
		d := schema.TestResourceDataRaw(t, r.Schema, tableConfig(true, primaryKey("pk1", false)))
		if err := r.Create(d, client); err != nil {
			t.Fatalf("create table: %v", err)
		}
		if len(fixture.createRequests) != 1 {
			t.Fatalf("expected 1 CreateTable request, got %d", len(fixture.createRequests))
		}
		request := &otsprotocol.CreateTableRequest{}
		if err := proto.Unmarshal(fixture.createRequests[0], request); err != nil {
			t.Fatalf("decode CreateTable request: %v", err)
		}
		if request.EnableLocalTxn == nil || !*request.EnableLocalTxn {
			t.Fatalf("EnableLocalTxn = %v, want true", request.EnableLocalTxn)
		}
	})
}
