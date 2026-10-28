// collection-order: bypass defined_column -- pre-existing attribute outside this change
// collection-order: bypass primary_key -- pre-existing attribute outside this change
package ots_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
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
