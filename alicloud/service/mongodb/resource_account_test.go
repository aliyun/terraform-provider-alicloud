package mongodb_test

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mongodb"
	sdkSchema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkTerraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudMongodbAccount_basicNormal(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccmongodbaccount%d", rand)
	resourceId := "alicloud_mongodb_account.default"
	testAccConfig := testAccMongodbAccountNormalConfigFunc(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// MongoDB capacity for the acceptance account is provisioned in
			// cn-shanghai; the zone itself is resolved from data sources.
			os.Setenv("ALICLOUD_REGION", "cn-shanghai")
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// The description of a replica set account can only be set through
				// ModifyAccountDescription, so it is applied inside Create.
				Config: testAccConfig(map[string]interface{}{
					"account_name":        "root",
					"account_password":    "YourPassword+12345",
					"character_type":      "normal",
					"db_name":             "admin",
					"account_description": name,
					"instance_id":         "${alicloud_mongodb_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "account_name", "root"),
					resource.TestCheckResourceAttr(resourceId, "character_type", "normal"),
					resource.TestCheckResourceAttr(resourceId, "db_name", "admin"),
					resource.TestCheckResourceAttr(resourceId, "account_description", name),
					resource.TestCheckResourceAttrSet(resourceId, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_name":        "root",
					"account_password":    "YourPassword+12345",
					"character_type":      "normal",
					"db_name":             "admin",
					"account_description": name + "update",
					"instance_id":         "${alicloud_mongodb_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "account_description", name+"update"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_name":        "root",
					"account_password":    "YourPassword+12345678",
					"character_type":      "normal",
					"db_name":             "admin",
					"account_description": name + "update",
					"instance_id":         "${alicloud_mongodb_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "account_description", name+"update"),
					resource.TestCheckResourceAttr(resourceId, "db_name", "admin"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

func TestAccAliCloudMongodbAccount_basicDbAccount(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccmongodbaccountdb%d", rand)
	accountName := fmt.Sprintf("tf_acc_%d", rand)
	importAccountName := fmt.Sprintf("tf_imp_%d", rand)
	resourceId := "alicloud_mongodb_account.default"
	importResourceId := "alicloud_mongodb_account.imported"
	testAccConfig := testAccMongodbAccountDbConfigFunc(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// MongoDB capacity for the acceptance account is provisioned in
			// cn-shanghai; the zone itself is resolved from data sources.
			os.Setenv("ALICLOUD_REGION", "cn-shanghai")
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_name":        accountName,
					"account_password":    "qwer1234!!!",
					"character_type":      "db",
					"db_name":             "admin",
					"account_description": name,
					"instance_id":         "${alicloud_mongodb_sharding_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "account_name", accountName),
					resource.TestCheckResourceAttr(resourceId, "character_type", "db"),
					resource.TestCheckResourceAttr(resourceId, "db_name", "admin"),
					resource.TestCheckResourceAttr(resourceId, "account_description", name),
					resource.TestCheckResourceAttrSet(resourceId, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
				),
			},
			{
				// CreateAccount rejects a database that does not exist yet
				// (InvalidParameter DBName), so the identity scope is exercised
				// with a second account in the existing admin database.
				Config: testAccConfig(map[string]interface{}{
					"account_name":        accountName,
					"account_password":    "qwer1234!!!",
					"character_type":      "db",
					"db_name":             "admin",
					"account_description": name,
					"import_account":      true,
					"import_account_name": importAccountName,
					"instance_id":         "${alicloud_mongodb_sharding_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "db_name", "admin"),
					resource.TestCheckResourceAttr(importResourceId, "account_name", importAccountName),
					resource.TestCheckResourceAttr(importResourceId, "db_name", "admin"),
					resource.TestCheckResourceAttrPair(resourceId, "instance_id", importResourceId, "instance_id"),
				),
			},
			{
				// ResetAccountPassword must carry DBName to address the account.
				Config: testAccConfig(map[string]interface{}{
					"account_name":        accountName,
					"account_password":    "qwer12345!!!",
					"character_type":      "db",
					"db_name":             "admin",
					"account_description": name,
					"import_account":      true,
					"import_account_name": importAccountName,
					"instance_id":         "${alicloud_mongodb_sharding_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "db_name", "admin"),
					resource.TestCheckResourceAttr(importResourceId, "account_name", importAccountName),
				),
			},
			{
				// Import with the database scoped id; the account has no
				// description because shard accounts do not return it.
				ResourceName:            importResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[importResourceId]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", importResourceId)
					}
					return fmt.Sprintf("%s:%s:%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["db_name"], rs.Primary.Attributes["account_name"]), nil
				},
			},
			{
				// CreateAccount documents that shard account privileges are fixed to
				// read-only, so configuring roles is kept as a plan-only check until
				// the product exposes a creatable role scenario.
				Config: testAccConfig(map[string]interface{}{
					"account_name":        accountName,
					"account_password":    "qwer12345!!!",
					"character_type":      "db",
					"db_name":             "admin",
					"account_description": name,
					"import_account":      true,
					"import_account_name": importAccountName,
					"instance_id":         "${alicloud_mongodb_sharding_instance.default.id}",
					"roles": []map[string]interface{}{
						{
							"role":    "read",
							"db_name": "admin",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccMongodbAccountNormalConfigFunc(name string) func(map[string]interface{}) string {
	return func(attrs map[string]interface{}) string {
		return fmt.Sprintf(`
%s

resource "alicloud_mongodb_account" "default" {
  account_name        = "%v"
  account_password    = "%v"
  instance_id         = "%v"
  character_type      = "%v"
  db_name             = "%v"
  account_description = "%v"
}
`, testAccMongodbAccountNormalDependence(name),
			attrs["account_name"], attrs["account_password"], attrs["instance_id"], attrs["character_type"],
			attrs["db_name"], attrs["account_description"])
	}
}

func testAccMongodbAccountDbConfigFunc(name string) func(map[string]interface{}) string {
	return func(attrs map[string]interface{}) string {
		roles, _ := attrs["roles"].([]map[string]interface{})
		var roleBlocks string
		for _, role := range roles {
			roleBlocks += fmt.Sprintf(`
  roles {
    role    = "%v"
    db_name = "%v"
  }`, role["role"], role["db_name"])
		}
		otherAccount := ""
		if v, ok := attrs["import_account"].(bool); ok && v {
			otherAccount = fmt.Sprintf(`
resource "alicloud_mongodb_account" "imported" {
  account_name     = "%v"
  account_password = "%v"
  instance_id      = "%v"
  character_type   = "db"
  db_name          = "%v"
}
`, attrs["import_account_name"], attrs["account_password"], attrs["instance_id"], attrs["db_name"])
		}
		return fmt.Sprintf(`
%s

resource "alicloud_mongodb_account" "default" {
  account_name        = "%v"
  account_password    = "%v"
  instance_id         = "%v"
  character_type      = "%v"
  db_name             = "%v"
  account_description = "%v"%s
}
%s
`, testAccMongodbAccountDbDependence(name),
			attrs["account_name"], attrs["account_password"], attrs["instance_id"], attrs["character_type"],
			attrs["db_name"], attrs["account_description"], roleBlocks, otherAccount)
	}
}

func testAccMongodbAccountNormalDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_mongodb_zones" "default" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

locals {
  vswitch_zone_ids = [for zone in data.alicloud_zones.default.zones : zone.id]
  mongodb_zone_ids = [for zone in data.alicloud_mongodb_zones.default.zones : zone.id if contains(local.vswitch_zone_ids, zone.id)]
  zone_id          = local.mongodb_zone_ids[0]
}

resource "alicloud_vpc" "default" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = local.zone_id
  cidr_block = "10.0.0.0/24"
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.2"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 20
  vswitch_id          = alicloud_vswitch.default.id
  zone_id             = local.zone_id
  name                = var.name
}
`, name)
}

func testAccMongodbAccountDbDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_mongodb_zones" "default" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

locals {
  vswitch_zone_ids = [for zone in data.alicloud_zones.default.zones : zone.id]
  mongodb_zone_ids = [for zone in data.alicloud_mongodb_zones.default.zones : zone.id if contains(local.vswitch_zone_ids, zone.id)]
  zone_id          = local.mongodb_zone_ids[0]
}

resource "alicloud_vpc" "default" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = local.zone_id
  cidr_block = "10.0.0.0/24"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  engine_version = "4.2"
  storage_type   = "cloud_essd1"
  vswitch_id     = alicloud_vswitch.default.id
  vpc_id         = alicloud_vpc.default.id
  storage_engine = "WiredTiger"
  network_type   = "VPC"
  zone_id        = local.zone_id
  name           = var.name

  mongo_list {
    node_class = "mdb.shard.4x.large.d"
  }
  mongo_list {
    node_class = "mdb.shard.4x.large.d"
  }
  shard_list {
    node_class   = "mdb.shard.4x.large.d"
    node_storage = "20"
  }
  shard_list {
    node_class   = "mdb.shard.4x.large.d"
    node_storage = "20"
  }
  config_server_list {
    node_class   = "mdb.shard.2x.xlarge.d"
    node_storage = "80"
  }
}
`, name)
}

func unitMongodbAccountResourceData(t *testing.T, id string, attrs map[string]interface{}, changes map[string]interface{}) *sdkSchema.ResourceData {
	t.Helper()
	state := &sdkTerraform.InstanceState{ID: id, Attributes: map[string]string{}}
	for k, v := range attrs {
		state.Attributes[k] = fmt.Sprint(v)
	}
	diff := &sdkTerraform.InstanceDiff{Attributes: map[string]*sdkTerraform.ResourceAttrDiff{}}
	for k, v := range changes {
		diff.Attributes[k] = &sdkTerraform.ResourceAttrDiff{Old: fmt.Sprint(attrs[k]), New: fmt.Sprint(v)}
	}
	d, err := sdkSchema.InternalMap(mongodb.ResourceAccount().Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("failed to build resource data: %v", err)
	}
	return d
}

func TestUnitMongodbAccountSchemaContract(t *testing.T) {
	r := mongodb.ResourceAccount()
	assert.True(t, r.Schema["db_name"].ForceNew, "db_name is part of the account identity")
	assert.True(t, r.Schema["db_name"].Computed)
	assert.True(t, r.Schema["roles"].Computed, "roles are returned by DescribeAccounts")
	assert.True(t, r.Schema["roles"].ForceNew)
	assert.True(t, r.Schema["character_type"].ForceNew)
}

func TestUnitMongodbAccountDescribeRejectsInvalidId(t *testing.T) {
	service := mongodb.MongodbServiceV2{}
	if _, err := service.DescribeMongodbAccount("dds-1"); err == nil {
		t.Fatal("expected an error for an id without an account name")
	}
	if _, err := service.DescribeMongodbAccount("dds-1:db:acc:extra"); err == nil {
		t.Fatal("expected an error for an id with too many parts")
	}
}

func TestUnitMongodbAccountStateRefreshChecksFieldError(t *testing.T) {
	service := mongodb.MongodbServiceV2{}
	refresh := service.MongodbAccountStateRefreshFuncWithApi("dds-1:admin:acc", "AccountStatus", nil, func(id string) (map[string]interface{}, error) {
		return map[string]interface{}{"DBName": "admin"}, nil
	})
	if _, _, err := refresh(); err == nil {
		t.Fatal("expected an error when the watched field is missing from the response")
	}
}

func TestUnitMongodbAccountReadNotFoundClearsId(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:acc", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc",
	}, nil)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&mongodb.MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *mongodb.MongodbServiceV2, id string) (map[string]interface{}, error) {
		return nil, tferrs.WrapErrorf(tferrs.NotFoundErr("Account", id), tferrs.NotFoundMsg, "mock")
	})
	defer patches.Reset()

	assert.Nil(t, mongodb.ResourceAccount().Read(d, &connectivity.AliyunClient{}))
	assert.Equal(t, "", d.Id())
}

func TestUnitMongodbAccountReadSetsDatabaseScopedId(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:acc1", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc1",
	}, nil)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&mongodb.MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *mongodb.MongodbServiceV2, id string) (map[string]interface{}, error) {
		return map[string]interface{}{
			"DBInstanceId":  "dds-1",
			"AccountName":   "acc1",
			"DBName":        "admin",
			"CharacterType": "db",
			"AccountStatus": "Available",
			"Roles": map[string]interface{}{
				"Role": []interface{}{map[string]interface{}{"DBName": "admin", "Role": "read"}},
			},
		}, nil
	})
	defer patches.Reset()

	assert.Nil(t, mongodb.ResourceAccount().Read(d, &connectivity.AliyunClient{}))
	assert.Equal(t, "dds-1:admin:acc1", d.Id())
	assert.Equal(t, "admin", d.Get("db_name"))
	assert.Equal(t, 1, d.Get("roles").(*sdkSchema.Set).Len())
}

func TestUnitMongodbAccountReadPropagatesError(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:acc", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc",
	}, nil)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&mongodb.MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *mongodb.MongodbServiceV2, id string) (map[string]interface{}, error) {
		return nil, errors.New("Forbidden.RAM")
	})
	defer patches.Reset()

	err := mongodb.ResourceAccount().Read(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected the DescribeAccounts error to be propagated")
	}
	assert.Equal(t, "dds-1:acc", d.Id())
}

func TestUnitMongodbAccountUpdateDescriptionDescribeError(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:admin:acc", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc",
		"db_name":      "admin",
	}, map[string]interface{}{
		"account_description": "updated",
	})
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&mongodb.MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *mongodb.MongodbServiceV2, id string) (map[string]interface{}, error) {
		return nil, errors.New("Forbidden.RAM")
	})
	defer patches.Reset()

	if err := mongodb.ResourceAccount().Update(d, &connectivity.AliyunClient{}); err == nil {
		t.Fatal("expected the pre-update DescribeAccounts error to be propagated")
	}
}

func TestUnitMongodbAccountUpdateShardDescriptionRejected(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:admin:acc", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc",
		"db_name":      "admin",
	}, map[string]interface{}{
		"account_description": "updated",
	})
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&mongodb.MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *mongodb.MongodbServiceV2, id string) (map[string]interface{}, error) {
		return map[string]interface{}{"DBInstanceId": "dds-1", "AccountName": "acc", "DBName": "admin", "CharacterType": "db", "AccountStatus": "Available"}, nil
	})
	defer patches.Reset()

	err := mongodb.ResourceAccount().Update(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected a shard account description update to be rejected")
	}
	assert.Contains(t, err.Error(), "replica set")
}

func TestUnitMongodbAccountCreateRejectsRolesForReplicaSet(t *testing.T) {
	d := sdkSchema.TestResourceDataRaw(t, mongodb.ResourceAccount().Schema, map[string]interface{}{
		"instance_id":      "dds-1",
		"account_name":     "root",
		"account_password": "YourPassword+12345",
		"character_type":   "normal",
		"db_name":          "admin",
		"roles": []interface{}{
			map[string]interface{}{"role": "read", "db_name": "admin"},
		},
	})
	err := mongodb.ResourceAccount().Create(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected roles on a replica set account to be rejected")
	}
	assert.Contains(t, err.Error(), "shard accounts")
}

func TestUnitMongodbAccountDeleteNoop(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:acc", map[string]interface{}{}, nil)
	assert.Nil(t, mongodb.ResourceAccount().Delete(d, &connectivity.AliyunClient{}))
}
