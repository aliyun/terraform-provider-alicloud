package alicloud

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudKmsKeyDataSource_All(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudKmsKeyDataSourceAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_keys.keys"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.description", "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.status", "Enabled"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.arn"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.creation_date"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.delete_date", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.creator"),
				),
			},
			{
				Config: testAccCheckAlicloudKmsKeyDataSourceAllWithErrorDescription,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_keys.keys"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.#", "0"),
				),
			},
			{
				Config: testAccCheckAlicloudKmsKeyDataSourceAllWithErrorIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_keys.keys"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.#", "0"),
				),
			},
			{
				Config: testAccCheckAlicloudKmsKeyDataSourceAllWithErrorStatus,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_keys.keys"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudKmsKeyDataSource_Description(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	randseed := rand.Int31n(10000)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAlicloudKmsKeyDataSourceDescription, randseed),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_keys.keys"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.description"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.status", "Enabled"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.arn"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.creation_date"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.delete_date", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.creator"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckAlicloudKmsKeyDataSourceDescriptionfake, randseed),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_keys.keys"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudKmsKeyDataSource_Status(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudKmsKeyDataSourceStatus,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_keys.keys"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.description", "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.status", "Enabled"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.arn"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.creation_date"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.delete_date", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.creator"),
				),
			},
			{
				Config: testAccCheckAlicloudKmsKeyDataSourceStatusfake,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_keys.keys"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudKmsKeyDataSource_Ids(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudKmsKeyDataSourceIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_keys.keys"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.description", "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.status", "Enabled"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.arn"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.creation_date"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.delete_date", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_keys.keys", "keys.0.creator"),
				),
			},
			{
				Config: testAccCheckAlicloudKmsKeyDataSourceIdsfake,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_keys.keys"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.#", "0"),
				),
			},
		},
	})
}

var testAccCheckAlicloudKmsKeyDataSourceDescription = `
resource "alicloud_kms_key" "key" {
    description = "tf_testaccCheckAlicloudKmsKeyDataSourceBasic%d"
    deletion_window_in_days = 7
}

data "alicloud_kms_keys" "keys" {
    description_regex = "^${alicloud_kms_key.key.description}"
}
`

var testAccCheckAlicloudKmsKeyDataSourceDescriptionfake = `
resource "alicloud_kms_key" "key" {
    description = "tf_testaccCheckAlicloudKmsKeyDataSourceBasic%d"
    deletion_window_in_days = 7
}

data "alicloud_kms_keys" "keys" {
    description_regex = "^${alicloud_kms_key.key.description}_fake"
}
`

const testAccCheckAlicloudKmsKeyDataSourceStatus = `
resource "alicloud_kms_key" "key" {
    description = "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"
    deletion_window_in_days = 7
}

data "alicloud_kms_keys" "keys" {
    status = "Enabled"
    ids = ["${alicloud_kms_key.key.id}"]
}
`

const testAccCheckAlicloudKmsKeyDataSourceStatusfake = `
resource "alicloud_kms_key" "key" {
    description = "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"
    deletion_window_in_days = 7
}

data "alicloud_kms_keys" "keys" {
    status = "Disabled"
    ids = ["${alicloud_kms_key.key.id}"]
}
`

const testAccCheckAlicloudKmsKeyDataSourceIds = `
resource "alicloud_kms_key" "key" {
    description = "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"
    deletion_window_in_days = 7
}

data "alicloud_kms_keys" "keys" {
    ids = ["${alicloud_kms_key.key.id}"]
}
`

const testAccCheckAlicloudKmsKeyDataSourceIdsfake = `
resource "alicloud_kms_key" "key" {
    description = "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"
    deletion_window_in_days = 7
}

data "alicloud_kms_keys" "keys" {
    ids = ["${alicloud_kms_key.key.id}test"]
}
`

const testAccCheckAlicloudKmsKeyDataSourceAll = `
resource "alicloud_kms_key" "key" {
    description = "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"
    deletion_window_in_days = 7
}

data "alicloud_kms_keys" "keys" {
    description_regex = "^${alicloud_kms_key.key.description}"
    ids = ["${alicloud_kms_key.key.id}"]
    status = "Enabled"
}
`

const testAccCheckAlicloudKmsKeyDataSourceAllWithErrorDescription = `
resource "alicloud_kms_key" "key" {
    description = "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"
    deletion_window_in_days = 7
}

data "alicloud_kms_keys" "keys" {
    description_regex = "^${alicloud_kms_key.key.description}_fake"
    ids = ["${alicloud_kms_key.key.id}"]
    status = "Enabled"
}
`

const testAccCheckAlicloudKmsKeyDataSourceAllWithErrorIds = `
resource "alicloud_kms_key" "key" {
    description = "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"
    deletion_window_in_days = 7
}

data "alicloud_kms_keys" "keys" {
    description_regex = "^${alicloud_kms_key.key.description}"
    ids = ["${alicloud_kms_key.key.id}test"]
    status = "Enabled"
}
`

const testAccCheckAlicloudKmsKeyDataSourceAllWithErrorStatus = `
resource "alicloud_kms_key" "key" {
    description = "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"
    deletion_window_in_days = 7
}

data "alicloud_kms_keys" "keys" {
    description_regex = "^${alicloud_kms_key.key.description}"
    ids = ["${alicloud_kms_key.key.id}"]
    status = "Disabled"
}
`
