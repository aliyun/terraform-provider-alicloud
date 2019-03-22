package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_mongodb_instance", &resource.Sweeper{
		Name: "alicloud_mongodb_instance",
		F:    testSweepMongoDBInstances,
	})
}

func TestAccAlicloudMongoDBInstance_classic(t *testing.T) {
	const res_format = `
resource "alicloud_mongodb_instance" "foo" {
	%s
}
`
	const res_name = "alicloud_mongodb_instance.foo"
	var instance dds.DBInstance
	var args testMongoDBArgs
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.MongoDBClassicNoSupportedRegions)
		},
		IDRefreshName: res_name,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBInstanceDestroy,
		Steps: append(
			[]resource.TestStep{
				//create instance
				{
					Config: testMongoDBCreateConfig(res_format, args.SetItem(
						testDefaultItems(res_name, "engine_version", "3.4"),
						testDefaultItems(res_name, "db_instance_storage", 10),
						testDefaultItems(res_name, "db_instance_class", "dds.mongo.mid"),
					)),
					Check: func() resource.TestCheckFunc {
						check := []resource.TestCheckFunc{testAccCheckMongoDBInstanceExists(res_name, &instance)}
						check = append(check, testCheckMongDBArgs(res_name, args)...)
						check = append(check, testCheckDefualtMongDBArgs(res_name, args)...)
						return resource.ComposeTestCheckFunc(check...)
					}(),
				}}, testAccAlicloudCommonSteps(res_format, res_name, &args, &instance)...),
	})

}

func TestAccAlicloudMongoDBInstance_vpc(t *testing.T) {
	const res_format = `
data "alicloud_zones" "default" {
  available_resource_creation = "${var.creation}"
}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

variable "creation" {
  default = "MongoDB"
}

variable "name" {
  default = "tf-testAccMongoDBInstance_vpc"
}

resource "alicloud_mongodb_instance" "foo" {
  vswitch_id = "${alicloud_vswitch.default.id}"
%s
}
`
	const res_name = "alicloud_mongodb_instance.foo"
	var instance dds.DBInstance
	var args testMongoDBArgs
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: res_name,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBInstanceDestroy,
		Steps: append(
			[]resource.TestStep{
				//create instance
				{
					Config: testMongoDBCreateConfig(res_format, args.SetItem(
						testDefaultItems(res_name, "engine_version", "3.4"),
						testDefaultItems(res_name, "db_instance_storage", 10),
						testDefaultItems(res_name, "db_instance_class", "dds.mongo.mid"),
						testDefaultItems(res_name, "name", "tf-testAccMongoDBInstance"),
					)),
					Check: func() resource.TestCheckFunc {
						check := []resource.TestCheckFunc{testAccCheckMongoDBInstanceExists(res_name, &instance)}
						check = append(check, testCheckMongDBArgs(res_name, args)...)
						check = append(check, testCheckDefualtMongDBArgs(res_name, args)...)
						return resource.ComposeTestCheckFunc(check...)
					}(),
				}}, testAccAlicloudCommonSteps(res_format, res_name, &args, &instance)...),
	})
}

func TestAccAlicloudMongoDBInstance_multiAZ(t *testing.T) {
	const res_format = `
data "alicloud_zones" "default" {
  available_resource_creation = "${var.creation}"
  multi                       = true
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.multi_zone_ids[0]}"
  name              = "${var.name}"
}
variable "creation" {
  default = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multiAZ"
}

resource "alicloud_mongodb_instance" "foo" {
  zone_id    = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_id = "${alicloud_vswitch.default.id}"
  %s
}
`
	const res_name = "alicloud_mongodb_instance.foo"
	var instance dds.DBInstance
	var args testMongoDBArgs
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.MongoDBMultiAzNoSupportedRegions)
		},
		IDRefreshName: res_name,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBInstanceDestroy,
		Steps: append(
			[]resource.TestStep{
				//create instance
				{
					Config: testMongoDBCreateConfig(res_format, args.SetItem(
						testDefaultItems(res_name, "engine_version", "3.4"),
						testDefaultItems(res_name, "db_instance_storage", 10),
						testDefaultItems(res_name, "db_instance_class", "dds.mongo.mid"),
						testDefaultItems(res_name, "name", "tf-testAccMongoDBInstance"),
					)),
					Check: func() resource.TestCheckFunc {
						check := []resource.TestCheckFunc{testAccCheckMongoDBInstanceExists(res_name, &instance),
							testAccCheckMongoDBInstanceMultiIZ(&instance)}
						check = append(check, testCheckMongDBArgs(res_name, args)...)
						check = append(check, testCheckDefualtMongDBArgs(res_name, args)...)
						return resource.ComposeTestCheckFunc(check...)
					}(),
				}},
			testAccAlicloudCommonSteps(res_format, res_name, &args, &instance)...),
	})
}

func TestAccAlicloudMongoDBInstance_multi_instance(t *testing.T) {
	const res_format = `
data "alicloud_zones" "default" {
  available_resource_creation = "${var.creation}"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
variable "creation" {
  default = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multi_instance"
}

resource "alicloud_mongodb_instance" "foo" {
vswitch_id = "${alicloud_vswitch.default.id}"
%s
}
resource "alicloud_mongodb_instance" "foo2" {
vswitch_id = "${alicloud_vswitch.default.id}"
%s
}
`
	const res_name = "alicloud_mongodb_instance.foo"
	var instance dds.DBInstance
	var args testMongoDBArgs

	const res_name2 = "alicloud_mongodb_instance.foo2"
	var instance2 dds.DBInstance
	var args2 testMongoDBArgs

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: res_name,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBInstanceDestroy,
		Steps: append(
			[]resource.TestStep{
				//create instance
				{
					Config: testMongoDBCreateConfig(res_format,
						args.SetItem(
							testDefaultItems(res_name, "engine_version", "3.4"),
							testDefaultItems(res_name, "db_instance_storage", 10),
							testDefaultItems(res_name, "db_instance_class", "dds.mongo.mid"),
							testDefaultItems(res_name, "name", "tf-testAccMongoDBInstance"),
						),
						args2.SetItem(
							testDefaultItems(res_name2, "engine_version", "3.2"),
							testDefaultItems(res_name2, "db_instance_storage", 20),
							testDefaultItems(res_name2, "db_instance_class", "dds.mongo.mid"),
							testDefaultItems(res_name2, "name", "tf-testAccMongoDBInstance_2"),
						)),
					Check: resource.ComposeTestCheckFunc(func() []resource.TestCheckFunc {
						check1 := []resource.TestCheckFunc{testAccCheckMongoDBInstanceExists(res_name, &instance)}
						check1 = append(check1, testCheckMongDBArgs(res_name, args)...)
						check1 = append(check1, testCheckDefualtMongDBArgs(res_name, args)...)

						check2 := []resource.TestCheckFunc{testAccCheckMongoDBInstanceExists(res_name2, &instance2)}
						check2 = append(check1, testCheckMongDBArgs(res_name2, args2)...)
						check2 = append(check1, testCheckDefualtMongDBArgs(res_name2, args2)...)
						return append(check1, check2...)
					}()...),
				},
				//update name
				{
					Config: testMongoDBCreateConfig(res_format,
						args.SetItem(testDefaultItems(res_name, "name", "tf-testAccMongoDBInstance_test")),
						args2.SetItem(testDefaultItems(res_name2, "name", "tf-testAccMongoDBInstance_test2")),
					),
					Check: resource.ComposeTestCheckFunc(func() []resource.TestCheckFunc {
						check1 := append(testCheckMongDBArgs(res_name, args), testCheckDefualtMongDBArgs(res_name, args)...)
						check2 := append(testCheckMongDBArgs(res_name2, args2), testCheckDefualtMongDBArgs(res_name2, args2)...)
						return append(check1, check2...)
					}()...),
				},
				//update instance_storage
				{
					Config: testMongoDBCreateConfig(res_format,
						args.SetItem(testDefaultItems(res_name, "db_instance_storage", 30)),
						args2.SetItem(testDefaultItems(res_name2, "db_instance_storage", 40)),
					),
					Check: resource.ComposeTestCheckFunc(func() []resource.TestCheckFunc {
						check1 := append(testCheckMongDBArgs(res_name, args), testCheckDefualtMongDBArgs(res_name, args)...)
						check2 := append(testCheckMongDBArgs(res_name2, args2), testCheckDefualtMongDBArgs(res_name2, args2)...)
						return append(check1, check2...)
					}()...),
				},
				//set account_password
				{
					Config: testMongoDBCreateConfig(res_format,
						args.SetItem(testDefaultItems(res_name, "account_password", "1234567@tests321")),
						args2.SetItem(testDefaultItems(res_name2, "account_password", "1234567@tests123")),
					),
					Check: resource.ComposeTestCheckFunc(func() []resource.TestCheckFunc {
						check1 := append(testCheckMongDBArgs(res_name, args), testCheckDefualtMongDBArgs(res_name, args)...)
						check2 := append(testCheckMongDBArgs(res_name2, args2), testCheckDefualtMongDBArgs(res_name2, args2)...)
						return append(check1, check2...)
					}()...),
				},
				//update all together
				{
					Config: testMongoDBCreateConfig(res_format,
						args.SetItem(testDefaultItems(res_name, "account_password", "1234567@tests321"),
							testDefaultItems(res_name, "db_instance_storage", 50),
							testDefaultItems(res_name, "name", "tf-testAccMongoDBInstance")),
						args2.SetItem(testDefaultItems(res_name2, "account_password", "1234567@tests321"),
							testDefaultItems(res_name2, "db_instance_storage", 50),
							testDefaultItems(res_name2, "name", "tf-testAccMongoDBInstance")),
					),
					Check: resource.ComposeTestCheckFunc(func() []resource.TestCheckFunc {
						check1 := append(testCheckMongDBArgs(res_name, args), testCheckDefualtMongDBArgs(res_name, args)...)
						check2 := append(testCheckMongDBArgs(res_name2, args2), testCheckDefualtMongDBArgs(res_name2, args2)...)
						return append(check1, check2...)
					}()...),
				},
			}),
	})
}

const testAccAlicloudMongoDBInstance_import_config = `
data "alicloud_zones" "default" {
  available_resource_creation = "${var.creation}"
}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

variable "creation" {
  default = "MongoDB"
}

variable "name" {
  default = "tf-testAccMongoDBInstance_vpc"
}

resource "alicloud_mongodb_instance" "foo" {
  engine_version      = "3.4"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  name                = "${var.name}"
  vswitch_id          = "${alicloud_vswitch.default.id}"
}
`

type testResourceArg_security_ips struct {
	ips []string
}

type testMongoDBArgItem struct {
	key        string
	value      interface{}
	check_func func() []resource.TestCheckFunc
}

type testMongoDBArgs []testMongoDBArgItem

func testMongoDBCreateConfig(format string, args ...interface{}) string {
	str := fmt.Sprintf(format, args...)
	return str
}

func testDefaultItems(res_name, key string, value interface{}) testMongoDBArgItem {
	return testMongoDBArgItem{key, value,
		func() []resource.TestCheckFunc {
			return []resource.TestCheckFunc{resource.TestCheckResourceAttr(res_name, key, fmt.Sprint(value))}
		}}
}

func (args *testMongoDBArgs) SetItem(items ...testMongoDBArgItem) *testMongoDBArgs {
	for _, setitem := range items {
		need_append := true
		for i, item := range *args {
			if item.key == setitem.key {
				(*args)[i] = setitem
				need_append = false
				break
			}
		}
		if need_append {
			*args = append(*args, setitem)
		}
	}
	return args
}

func (args testMongoDBArgs) String() string {
	str := "\n"
	for _, item := range args {
		switch item.value.(type) {
		case string:
			str += fmt.Sprintf("%s = %q \n", item.key, item.value)
		case testResourceArg_security_ips:
			temp := ""
			for _, ip := range item.value.(testResourceArg_security_ips).ips {
				temp += fmt.Sprintf("%q,", ip)
			}
			temp = temp[:len(temp)-1]
			str += fmt.Sprintf("%s = [%s] \n", item.key, temp)
		default:
			str += fmt.Sprintf("%s = %v \n", item.key, item.value)
		}
	}
	return str
}

func testCheckMongDBArgs(res_name string, args testMongoDBArgs) []resource.TestCheckFunc {
	test_funcs := []resource.TestCheckFunc{}
	for _, item := range args {
		test_funcs = append(test_funcs, (item.check_func())...)
	}
	return test_funcs
}

func testCheckDefualtMongDBArgs(res_name string, args testMongoDBArgs) []resource.TestCheckFunc {
	default_check_funcs := map[string]resource.TestCheckFunc{
		"name":                 resource.TestCheckResourceAttr(res_name, "name", ""),
		"replication_factor":   resource.TestCheckResourceAttr(res_name, "replication_factor", "3"),
		"storage_engine":       resource.TestCheckResourceAttr(res_name, "storage_engine", "WiredTiger"),
		"instance_charge_type": resource.TestCheckResourceAttr(res_name, "instance_charge_type", "PostPaid"),
		"security_ip_list":     resource.TestCheckResourceAttrSet(res_name, "security_ip_list.#"),
	}

	test_funcs := []resource.TestCheckFunc{}
	for _, item := range args {
		test_funcs = append(test_funcs, (item.check_func())...)
		delete(default_check_funcs, item.key)
	}

	for _, item := range default_check_funcs {
		test_funcs = append(test_funcs, item)
	}

	return test_funcs
}

func testAccAlicloudCommonSteps(res_format, res_name string, args *testMongoDBArgs, instance *dds.DBInstance) []resource.TestStep {
	Steps := []resource.TestStep{
		//update name
		{
			Config: testMongoDBCreateConfig(res_format,
				args.SetItem(testDefaultItems(res_name, "name", "tf-testAccMongoDBInstance_test"))),
			Check: func() resource.TestCheckFunc {
				check := append(testCheckMongDBArgs(res_name, *args), testCheckDefualtMongDBArgs(res_name, *args)...)
				return resource.ComposeTestCheckFunc(check...)
			}(),
		},
		//Configuration Upgrade
		{
			Config: testMongoDBCreateConfig(res_format,
				args.SetItem(
					testDefaultItems(res_name, "db_instance_storage", 30),
					testDefaultItems(res_name, "db_instance_class", "dds.mongo.standard"),
				)),
			Check: func() resource.TestCheckFunc {
				check := append(testCheckMongDBArgs(res_name, *args), testCheckDefualtMongDBArgs(res_name, *args)...)
				return resource.ComposeTestCheckFunc(check...)
			}(),
		},
		//set-update account_password
		{
			Config: testMongoDBCreateConfig(res_format,
				args.SetItem(testDefaultItems(res_name, "account_password", "1234567@tests"))),
			Check: func() resource.TestCheckFunc {
				check := append(testCheckMongDBArgs(res_name, *args), testCheckDefualtMongDBArgs(res_name, *args)...)
				return resource.ComposeTestCheckFunc(check...)
			}(),
		},
		//set-update security_ip_list
		{
			Config: testMongoDBCreateConfig(res_format,
				args.SetItem(
					testMongoDBArgItem{"security_ip_list", testResourceArg_security_ips{[]string{"10.168.1.12"}}, func() []resource.TestCheckFunc {
						var ips []map[string]interface{}
						return []resource.TestCheckFunc{
							testAccCheckMongoDBSecurityIpExists(res_name, &ips),
							testAccCheckMongoDBInstanceKeyValueInMaps(&ips, "security_ip_list", "10.168.1.12")}
					}})),
			Check: func() resource.TestCheckFunc {
				check := append(testCheckMongDBArgs(res_name, *args), testCheckDefualtMongDBArgs(res_name, *args)...)
				return resource.ComposeTestCheckFunc(check...)
			}(),
		},
		{
			Config: testMongoDBCreateConfig(res_format,
				args.SetItem(
					testMongoDBArgItem{"security_ip_list", testResourceArg_security_ips{[]string{"10.168.1.12", "100.69.7.112"}}, func() []resource.TestCheckFunc {
						var ips []map[string]interface{}
						return []resource.TestCheckFunc{
							testAccCheckMongoDBSecurityIpExists(res_name, &ips),
							testAccCheckMongoDBInstanceKeyValueInMaps(&ips, "security_ip_list", "10.168.1.12,100.69.7.112")}
					}})),
			Check: func() resource.TestCheckFunc {
				check := append(testCheckMongDBArgs(res_name, *args), testCheckDefualtMongDBArgs(res_name, *args)...)
				return resource.ComposeTestCheckFunc(check...)
			}(),
		},
		//all together update
		{
			Config: testMongoDBCreateConfig(res_format,
				args.SetItem(
					testDefaultItems(res_name, "account_password", "1234567@tests"),
					testDefaultItems(res_name, "name", "tf-testAccMongoDBInstanceClassic"),
					testMongoDBArgItem{"security_ip_list", testResourceArg_security_ips{[]string{"10.168.1.12"}}, func() []resource.TestCheckFunc {
						var ips []map[string]interface{}
						ret := []resource.TestCheckFunc{
							testAccCheckMongoDBSecurityIpExists(res_name, &ips),
							testAccCheckMongoDBInstanceKeyValueInMaps(&ips, "security_ip_list", "10.168.1.12")}
						return ret
					}})),
			Check: func() resource.TestCheckFunc {
				check := append(testCheckMongDBArgs(res_name, *args), testCheckDefualtMongDBArgs(res_name, *args)...)
				return resource.ComposeTestCheckFunc(check...)
			}(),
		},
	}
	return Steps
}

func testSweepMongoDBInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []dds.DBInstance
	request := dds.CreateDescribeDBInstancesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DescribeDBInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "testSweepMongoDBInstances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*dds.DescribeDBInstancesResponse)
		addDebug(request.GetActionName(), response)

		if response == nil || len(response.DBInstances.DBInstance) < 1 {
			break
		}
		insts = append(insts, response.DBInstances.DBInstance...)

		if len(response.DBInstances.DBInstance) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range insts {
		name := v.DBInstanceDescription
		id := v.DBInstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			continue
		}

		sweeped = true

		request := dds.CreateDeleteDBInstanceRequest()
		request.DBInstanceId = id
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DeleteDBInstance(request)
		})

		if err != nil {
			log.Printf("[error] %v %v", id, request.GetActionName())
		}
		addDebug(request.GetActionName(), raw)
	}
	if sweeped {
		// Waiting 30 seconds to eusure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func testAccCheckMongoDBInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_mongodb_instance" {
			continue
		}
		_, err := ddsService.DescribeMongoDBInstance(rs.Primary.ID)
		if err != nil {
			if ddsService.NotFoundMongoDBInstance(err) {
				continue
			}
			return WrapError(err)
		}
		return err
	}
	return nil
}

func testAccCheckMongoDBInstanceExists(n string, d *dds.DBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No MongoDB Instance ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		service := MongoDBService{client}
		attr, err := service.DescribeMongoDBInstance(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		*d = *attr
		return nil
	}
}

func testAccCheckMongoDBInstanceMultiIZ(i *dds.DBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !strings.Contains(i.ZoneId, MULTI_IZ_SYMBOL) {
			return WrapError(fmt.Errorf("Current region does not support multiIZ."))
		}
		return nil
	}
}

func testAccCheckMongoDBSecurityIpExists(n string, ips *[]map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No DB Instance ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ddsService := MongoDBService{client}
		respone, err := ddsService.DescribeMongoDBSecurityIps(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}
		if len(respone) < 1 {
			return WrapError(fmt.Errorf("DB security ip not found"))
		}
		result := make([]map[string]interface{}, 0, len(respone))
		for _, i := range respone {
			l := map[string]interface{}{
				"security_ip_list": i.SecurityIpList,
			}
			result = append(result, l)
		}
		*ips = result
		return nil
	}
}

func testAccCheckMongoDBInstanceKeyValueInMaps(ps *[]map[string]interface{}, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, policy := range *ps {
			if policy[key].(string) != value {
				return WrapError(fmt.Errorf("MongoDB attribute '%s' expected %#v, got %#v", key, value, policy[key]))
			}
		}
		return nil
	}
}
