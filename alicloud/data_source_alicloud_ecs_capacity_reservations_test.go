package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudECSCapacityReservationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_capacity_reservation.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_capacity_reservation.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecs_capacity_reservation.default.capacity_reservation_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecs_capacity_reservation.default.capacity_reservation_name}_fake"`,
		}),
	}
	capacityReservationIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"capacity_reservation_ids": `["${alicloud_ecs_capacity_reservation.default.id}"]`,
			"ids":                      `["${alicloud_ecs_capacity_reservation.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"capacity_reservation_ids": `["${alicloud_ecs_capacity_reservation.default.id}_fake"]`,
			"ids":                      `["${alicloud_ecs_capacity_reservation.default.id}"]`,
		}),
	}
	instanceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"instance_type": `"${alicloud_ecs_capacity_reservation.default.instance_type}"`,
			"ids":           `["${alicloud_ecs_capacity_reservation.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"instance_type": `"${alicloud_ecs_capacity_reservation.default.instance_type}_fake"`,
			"ids":           `["${alicloud_ecs_capacity_reservation.default.id}"]`,
		}),
	}
	paymentTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"payment_type": `"${alicloud_ecs_capacity_reservation.default.payment_type}"`,
			"ids":          `["${alicloud_ecs_capacity_reservation.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"payment_type": `"PrePaid"`,
			"ids":          `["${alicloud_ecs_capacity_reservation.default.id}"]`,
		}),
	}
	platformConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_ecs_capacity_reservation.default.id}"]`,
			"platform": `"${alicloud_ecs_capacity_reservation.default.platform}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"platform": `"windows"`,
			"ids":      `["${alicloud_ecs_capacity_reservation.default.id}"]`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_ecs_capacity_reservation.default.resource_group_id}"`,
			"ids":               `["${alicloud_ecs_capacity_reservation.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_ecs_capacity_reservation.default.resource_group_id}_fake"`,
			"ids":               `["${alicloud_ecs_capacity_reservation.default.id}"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"status": `"${alicloud_ecs_capacity_reservation.default.status}"`,
			"ids":    `["${alicloud_ecs_capacity_reservation.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"status": `"Released"`,
			"ids":    `["${alicloud_ecs_capacity_reservation.default.id}"]`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_capacity_reservation.default.id}"]`,
			"tags": `{ 
						"Created" = "tfTestAcc0"
    					"For"     = "Tftestacc 0" 
					}`,
		}),
		fakeConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_capacity_reservation.default.id}"]`,
			"tags": `{ 
						"Created" = "tfTestAcc0-fake"
    					"For"     = "Tftestacc 0-fake" 
					}`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"status": `"${alicloud_ecs_capacity_reservation.default.status}"`,
			"tags": `{ 
						"Created" = "tfTestAcc0"
    					"For"     = "Tftestacc 0" 
					}`,
			"name_regex":               `"${alicloud_ecs_capacity_reservation.default.capacity_reservation_name}"`,
			"capacity_reservation_ids": `["${alicloud_ecs_capacity_reservation.default.id}"]`,
			"instance_type":            `"${alicloud_ecs_capacity_reservation.default.instance_type}"`,
			"resource_group_id":        `"${alicloud_ecs_capacity_reservation.default.resource_group_id}"`,
			"ids":                      `["${alicloud_ecs_capacity_reservation.default.id}"]`,
			"payment_type":             `"${alicloud_ecs_capacity_reservation.default.payment_type}"`,
			"platform":                 `"${alicloud_ecs_capacity_reservation.default.platform}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ecs_capacity_reservation.default.id}_fake"]`,
			"payment_type":      `"PrePaid"`,
			"resource_group_id": `"${alicloud_ecs_capacity_reservation.default.resource_group_id}_fake"`,
			"tags": `{ 
						"Created" = "tfTestAcc0-fake"
    					"For"     = "Tftestacc 0-fake" 
					}`,
			"name_regex":               `"${alicloud_ecs_capacity_reservation.default.capacity_reservation_name}_fake"`,
			"capacity_reservation_ids": `["${alicloud_ecs_capacity_reservation.default.id}_fake"]`,
			"instance_type":            `"${alicloud_ecs_capacity_reservation.default.instance_type}_fake"`,
			"platform":                 `"windows"`,
			"status":                   `"Released"`,
		}),
	}
	var existAlicloudEcsCapacityReservationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"names.#":                                "1",
			"reservations.#":                         "1",
			"reservations.0.id":                      CHECKSET,
			"reservations.0.capacity_reservation_id": CHECKSET,
			"reservations.0.capacity_reservation_name": CHECKSET,
			"reservations.0.description":               CHECKSET,
			"reservations.0.end_time":                  CHECKSET,
			"reservations.0.end_time_type":             "Unlimited",
			"reservations.0.instance_type":             CHECKSET,
			"reservations.0.match_criteria":            "Open",
			"reservations.0.payment_type":              "PostPaid",
			"reservations.0.platform":                  "linux",
			"reservations.0.resource_group_id":         CHECKSET,
			"reservations.0.start_time":                CHECKSET,
			"reservations.0.start_time_type":           CHECKSET,
			"reservations.0.status":                    "Active",
			"reservations.0.tags.%":                    "2",
			"reservations.0.tags.Created":              "tfTestAcc0",
			"reservations.0.tags.For":                  "Tftestacc 0",
			"reservations.0.time_slot":                 "",
			"reservations.0.zone_ids.#":                "1",
			"reservations.0.instance_amount":           "1",
		}
	}
	var fakeAlicloudEcsCapacityReservationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var AlicloudEcsCapacityReservationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_capacity_reservations.default",
		existMapFunc: existAlicloudEcsCapacityReservationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcsCapacityReservationsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	AlicloudEcsCapacityReservationsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, capacityReservationIdsConf, instanceTypeConf, paymentTypeConf, platformConf, resourceGroupIdConf, statusConf, tagsConf, allConf)
}
func testAccCheckAlicloudEcsCapacityReservationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccCapacityReservation-%d"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}


resource "alicloud_ecs_capacity_reservation" "default" {
	description =                var.name
	platform =                   "linux"
	capacity_reservation_name =  var.name
	end_time_type =              "Unlimited"
	resource_group_id =          data.alicloud_resource_manager_resource_groups.default.ids.0
	instance_amount =            1
	instance_type =              "ecs.c5.2xlarge"
	match_criteria =             "Open"
	tags = {
		Created =  "tfTestAcc0"
		For =      "Tftestacc 0"
	}
	zone_ids = [data.alicloud_zones.default.zones[0].id]
}

data "alicloud_ecs_capacity_reservations" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
