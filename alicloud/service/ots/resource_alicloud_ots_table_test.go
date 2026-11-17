// collection-order: bypass defined_column -- pre-existing attribute outside this change
package ots_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
	"github.com/hashicorp/go-cty/cty"
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
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"primary_key.#":                "2",
		"primary_key.0.name":           "pk1",
		"primary_key.0.type":           "Integer",
		"primary_key.0.auto_increment": "false",
		"primary_key.1.name":           "pk2",
		"primary_key.1.type":           "Integer",
		"primary_key.1.auto_increment": "true",
		"time_to_live":                 "-1",
		"max_version":                  "1",
	})
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
					"table_name":    "tf_${var.name}",
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
					testAccCheck(nil),
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

func TestAccAliCloudOtsTable_withLocalTxn(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_table.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"primary_key.#":      "1",
		"primary_key.0.name": "pk1",
		"primary_key.0.type": "Integer",
		"time_to_live":       "-1",
		"max_version":        "1",
		"enable_local_txn":   "true",
	})
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
					"table_name":    "tf_${var.name}",
					"primary_key": []map[string]interface{}{
						{
							"name": "pk1",
							"type": "Integer",
						},
					},
					"time_to_live":     "-1",
					"max_version":      "1",
					"enable_local_txn": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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
					"time_to_live": "86401",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_to_live":     "86401",
						"enable_local_txn": "true",
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
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"primary_key.#":      "1",
		"primary_key.0.name": "pk1",
		"primary_key.0.type": "Integer",
		"time_to_live":       "-1",
		"max_version":        "1",
		"allow_update":       "false",
	})
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
					"table_name":    "tf_${var.name}",
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
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudOtsTable_primaryKeyOrder(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_table.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"primary_key.#":      "2",
		"primary_key.0.name": "pk1",
		"primary_key.0.type": "Integer",
		"primary_key.1.name": "pk2",
		"primary_key.1.type": "String",
		"time_to_live":       "-1",
		"max_version":        "1",
	})
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
					"table_name":    "tf_${var.name}",
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Re-apply the imported configuration so the reordered apply below
				// stays consecutive; an import step between two applies resets the
				// collection-order snapshot.
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "tf_${var.name}",
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
				Config: testAccConfig(map[string]interface{}{
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"primary_key.0.name": "pk2",
						"primary_key.0.type": "String",
						"primary_key.1.name": "pk1",
						"primary_key.1.type": "Integer",
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

func TestAccAliCloudOtsTable_autoIncrementValidation(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alicloud_ots_table.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
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
					"table_name":    "tf_${var.name}",
					"primary_key": []map[string]interface{}{
						{
							"name": "pk1",
							"type": "String",
						},
						{
							"name":           "pk2",
							"type":           "String",
							"auto_increment": "true",
						},
					},
					"time_to_live": "-1",
					"max_version":  "1",
				}),
				ExpectError: regexp.MustCompile("`auto_increment` can only be enabled for a primary key of type"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"primary_key": []map[string]interface{}{
						{
							"name":           "pk1",
							"type":           "Integer",
							"auto_increment": "true",
						},
					},
				}),
				ExpectError: regexp.MustCompile("`auto_increment` cannot be enabled for the first primary key"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
						{
							"name":           "pk3",
							"type":           "Integer",
							"auto_increment": "true",
						},
					},
				}),
				ExpectError: regexp.MustCompile("only one primary key can enable `auto_increment`"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
					"enable_local_txn": "true",
				}),
				ExpectError: regexp.MustCompile("`enable_local_txn` cannot be enabled when a primary key enables `auto_increment`"),
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

// The create request is asserted from the bytes the real Tablestore client
// sends to an isolated endpoint, driven through the public resource apply path.
func TestUnitTableCreateRequestOptions(t *testing.T) {
	baseConfig := func() map[string]interface{} {
		return map[string]interface{}{
			"instance_name": "test-instance",
			"table_name":    "test_table",
			"primary_key": []interface{}{
				map[string]interface{}{
					"name": "pk1",
					"type": "Integer",
				},
			},
			"time_to_live": -1,
			"max_version":  1,
		}
	}

	t.Run("LocalTransactionTrue", func(t *testing.T) {
		config := baseConfig()
		config["enable_local_txn"] = true
		request, err := createTableThroughFixture(t, config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if request.EnableLocalTxn == nil || !*request.EnableLocalTxn {
			t.Fatalf("expected EnableLocalTxn=true in the serialized CreateTable request, got %v", request.EnableLocalTxn)
		}
		for _, pk := range request.TableMeta.PrimaryKey {
			if pk.Option != nil && *pk.Option == otsprotocol.PrimaryKeyOption_AUTO_INCREMENT {
				t.Fatalf("expected primary key %s not to be auto-incremented", pk.GetName())
			}
		}
	})

	t.Run("LocalTransactionFalse", func(t *testing.T) {
		config := baseConfig()
		config["enable_local_txn"] = false
		request, err := createTableThroughFixture(t, config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if request.EnableLocalTxn == nil || *request.EnableLocalTxn {
			t.Fatalf("expected an explicit EnableLocalTxn=false in the serialized CreateTable request, got %v", request.EnableLocalTxn)
		}
	})

	t.Run("LocalTransactionOmitted", func(t *testing.T) {
		request, err := createTableThroughFixture(t, baseConfig())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if request.EnableLocalTxn != nil {
			t.Fatalf("expected EnableLocalTxn to stay unset when the argument is omitted, got %v", *request.EnableLocalTxn)
		}
	})

	t.Run("AutoIncrementPrimaryKey", func(t *testing.T) {
		config := baseConfig()
		config["primary_key"] = []interface{}{
			map[string]interface{}{
				"name": "pk1",
				"type": "Integer",
			},
			map[string]interface{}{
				"name":           "pk2",
				"type":           "Integer",
				"auto_increment": true,
			},
		}
		request, err := createTableThroughFixture(t, config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if request.EnableLocalTxn != nil {
			t.Fatalf("expected EnableLocalTxn to stay unset when the argument is omitted, got %v", *request.EnableLocalTxn)
		}
		if len(request.TableMeta.PrimaryKey) != 2 {
			t.Fatalf("expected 2 primary keys, got %d", len(request.TableMeta.PrimaryKey))
		}
		if pk := request.TableMeta.PrimaryKey[0]; pk.Option != nil && *pk.Option == otsprotocol.PrimaryKeyOption_AUTO_INCREMENT {
			t.Fatalf("expected the first primary key not to be auto-incremented")
		}
		if pk := request.TableMeta.PrimaryKey[1]; pk.Option == nil || *pk.Option != otsprotocol.PrimaryKeyOption_AUTO_INCREMENT {
			t.Fatalf("expected the second primary key to be auto-incremented")
		}
	})
}

// createTableThroughFixture runs the public resource create path against a local
// endpoint and returns the decoded CreateTable request the SDK actually sent.
func createTableThroughFixture(t *testing.T, config map[string]interface{}) (*otsprotocol.CreateTableRequest, error) {
	t.Helper()
	t.Setenv("NO_PROXY", "127.0.0.1,localhost")

	var mutex sync.Mutex
	var createTableBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/openapi/getinstance":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"InstanceName":"test-instance","InstanceStatus":"Normal"}`)
		case "/CreateTable":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			mutex.Lock()
			createTableBody = body
			mutex.Unlock()
			w.WriteHeader(http.StatusOK)
		case "/DescribeTable":
			body, err := describeTableFixtureResponse(config).XXX_Marshal(nil, false)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/x-protobuf")
			_, _ = w.Write(body)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-access-key").
		SetAccessKeySecret("test-secret-key"))
	if err != nil {
		t.Fatalf("failed to build the test credential: %v", err)
	}
	endpoints := &sync.Map{}
	endpoints.Store("tablestore", strings.TrimPrefix(server.URL, "http://"))
	client, err := (&connectivity.Config{
		AccessKey:            "test-access-key",
		SecretKey:            "test-secret-key",
		Region:               connectivity.Region("cn-beijing"),
		RegionId:             "cn-beijing",
		AccountType:          "Domestic",
		Protocol:             "HTTP",
		ClientReadTimeout:    5000,
		ClientConnectTimeout: 5000,
		SkipRegionValidation: true,
		OtsEndpoint:          server.URL,
		Endpoints:            endpoints,
		Credential:           credential,
	}).Client()
	if err != nil {
		t.Fatalf("failed to build the test client: %v", err)
	}

	resourceOtsTable := ots.ResourceOtsTable()
	resourceConfig := sdkterraform.NewResourceConfigRaw(config)
	diff, err := resourceOtsTable.Diff(context.Background(), nil, resourceConfig, client)
	if err != nil {
		return nil, err
	}
	if diff == nil {
		return nil, fmt.Errorf("no diff produced for the test configuration")
	}
	// Resource.Diff does not carry the raw configuration into the diff, so the
	// create path would read an empty raw config without this.
	diff.RawConfig = rawConfigValue(t, resourceOtsTable.CoreConfigSchema().ImpliedType(), config)
	if _, diags := resourceOtsTable.Apply(context.Background(), nil, diff, client); diags.HasError() {
		return nil, fmt.Errorf("failed to apply the test configuration: %v", diags)
	}

	mutex.Lock()
	body := createTableBody
	mutex.Unlock()
	if len(body) == 0 {
		return nil, fmt.Errorf("no CreateTable request reached the endpoint")
	}
	request := new(otsprotocol.CreateTableRequest)
	if err := request.XXX_Unmarshal(body); err != nil {
		return nil, err
	}
	return request, nil
}

// rawConfigValue builds the raw configuration a real plan produces: every schema
// attribute is present and an omitted argument is a typed null, not a missing key.
func rawConfigValue(t *testing.T, impliedType cty.Type, config map[string]interface{}) cty.Value {
	t.Helper()

	return rawConfigTypedValue(t, impliedType, config)
}

func rawConfigTypedValue(t *testing.T, typ cty.Type, value interface{}) cty.Value {
	t.Helper()

	switch {
	case typ.IsObjectType():
		attributes := make(map[string]cty.Value, len(typ.AttributeTypes()))
		configured, ok := value.(map[string]interface{})
		if !ok {
			t.Fatalf("expected an object value for %s, got %T", typ.FriendlyName(), value)
		}
		for name, attributeType := range typ.AttributeTypes() {
			attributeValue, present := configured[name]
			if !present {
				attributes[name] = cty.NullVal(attributeType)
				continue
			}
			attributes[name] = rawConfigTypedValue(t, attributeType, attributeValue)
		}
		return cty.ObjectVal(attributes)
	case typ.IsListType() || typ.IsSetType():
		items, ok := value.([]interface{})
		if !ok {
			t.Fatalf("expected a collection value for %s, got %T", typ.FriendlyName(), value)
		}
		elementType := typ.ElementType()
		elements := make([]cty.Value, 0, len(items))
		for _, item := range items {
			elements = append(elements, rawConfigTypedValue(t, elementType, item))
		}
		switch {
		case typ.IsSetType() && len(elements) == 0:
			return cty.SetValEmpty(elementType)
		case typ.IsSetType():
			return cty.SetVal(elements)
		case len(elements) == 0:
			return cty.ListValEmpty(elementType)
		default:
			return cty.ListVal(elements)
		}
	case typ.Equals(cty.String):
		stringValue, ok := value.(string)
		if !ok {
			t.Fatalf("expected a string value for %s, got %T", typ.FriendlyName(), value)
		}
		return cty.StringVal(stringValue)
	case typ.Equals(cty.Bool):
		boolValue, ok := value.(bool)
		if !ok {
			t.Fatalf("expected a bool value for %s, got %T", typ.FriendlyName(), value)
		}
		return cty.BoolVal(boolValue)
	case typ.Equals(cty.Number):
		intValue, ok := value.(int)
		if !ok {
			t.Fatalf("expected an int value for %s, got %T", typ.FriendlyName(), value)
		}
		return cty.NumberIntVal(int64(intValue))
	default:
		t.Fatalf("unsupported raw config type %s", typ.FriendlyName())
		return cty.NilVal
	}
}

func describeTableFixtureResponse(config map[string]interface{}) *otsprotocol.DescribeTableResponse {
	tableName := config["table_name"].(string)
	timeToLive := int32(config["time_to_live"].(int))
	maxVersion := int32(config["max_version"].(int))
	allowUpdate := true
	deviation := int64(0)

	primaryKeys := make([]*otsprotocol.PrimaryKeySchema, 0)
	for _, item := range config["primary_key"].([]interface{}) {
		member := item.(map[string]interface{})
		keyType := otsprotocol.PrimaryKeyType_INTEGER
		if member["type"] == "String" {
			keyType = otsprotocol.PrimaryKeyType_STRING
		} else if member["type"] == "Binary" {
			keyType = otsprotocol.PrimaryKeyType_BINARY
		}
		key := &otsprotocol.PrimaryKeySchema{
			Name: primaryKeyStringPointer(member["name"].(string)),
			Type: &keyType,
		}
		if autoIncrement, ok := member["auto_increment"].(bool); ok && autoIncrement {
			option := otsprotocol.PrimaryKeyOption_AUTO_INCREMENT
			key.Option = &option
		}
		primaryKeys = append(primaryKeys, key)
	}

	return &otsprotocol.DescribeTableResponse{
		TableMeta: &otsprotocol.TableMeta{
			TableName:  primaryKeyStringPointer(tableName),
			PrimaryKey: primaryKeys,
		},
		ReservedThroughputDetails: &otsprotocol.ReservedThroughputDetails{
			CapacityUnit: &otsprotocol.CapacityUnit{
				Read:  primaryKeyInt32Pointer(0),
				Write: primaryKeyInt32Pointer(0),
			},
			LastIncreaseTime: primaryKeyInt64Pointer(0),
		},
		TableOptions: &otsprotocol.TableOptions{
			TimeToLive:                &timeToLive,
			MaxVersions:               &maxVersion,
			DeviationCellVersionInSec: &deviation,
			AllowUpdate:               &allowUpdate,
		},
		TableStatus: tableStatusPointer(otsprotocol.TableStatus_ACTIVE),
	}
}

func tableStatusPointer(value otsprotocol.TableStatus) *otsprotocol.TableStatus {
	return &value
}

func primaryKeyStringPointer(value string) *string {
	return &value
}

func primaryKeyInt32Pointer(value int32) *int32 {
	return &value
}

func primaryKeyInt64Pointer(value int64) *int64 {
	return &value
}
