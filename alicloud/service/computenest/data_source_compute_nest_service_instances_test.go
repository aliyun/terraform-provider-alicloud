package computenest_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudComputeNestServiceInstancesDataSource exercises every
// Optional argument of the alicloud_compute_nest_service_instances data source
// (ids, name_regex, status, filter, filter.name, filter.value, tags,
// output_file) through inline standard-template TestSteps
// (resource.Test + resource.ComposeTestCheckFunc with
// Config: testAccConfig(map[string]interface{}{...})). The standard-template
// Config blocks are what the TestingCoverageRate checker parses, so each
// Optional argument is recorded as covered; the Exist/Fake pairs reuse the
// same DataSourceAttr assertions the previous DataSourceTestCheckWithPreCheck
// wrapper produced, preserving the runtime behaviour and assertions while
// making attribute coverage measurable.
func TestAccAliCloudComputeNestServiceInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_compute_nest_service_instances.default"
	name := fmt.Sprintf("tf-testacc-ComputeNestServiceInstance%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceServiceInstancesConfig)

	var existServiceInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"names.#":                "1",
			"service_instances.#":    "1",
			"service_instances.0.id": acctest.CHECKSET,
			"service_instances.0.service_instance_id":                         acctest.CHECKSET,
			"service_instances.0.service_instance_name":                       acctest.CHECKSET,
			"service_instances.0.enable_instance_ops":                         acctest.CHECKSET,
			"service_instances.0.operation_start_time":                        acctest.CHECKSET,
			"service_instances.0.operation_end_time":                          acctest.CHECKSET,
			"service_instances.0.resources":                                   acctest.CHECKSET,
			"service_instances.0.source":                                      acctest.CHECKSET,
			"service_instances.0.service.#":                                   "1",
			"service_instances.0.service.0.service_id":                        acctest.CHECKSET,
			"service_instances.0.service.0.service_type":                      acctest.CHECKSET,
			"service_instances.0.service.0.deploy_type":                       acctest.CHECKSET,
			"service_instances.0.service.0.supplier_name":                     acctest.CHECKSET,
			"service_instances.0.service.0.supplier_url":                      acctest.CHECKSET,
			"service_instances.0.service.0.publish_time":                      acctest.CHECKSET,
			"service_instances.0.service.0.version":                           "1",
			"service_instances.0.service.0.version_name":                      acctest.CHECKSET,
			"service_instances.0.service.0.service_infos.#":                   "1",
			"service_instances.0.service.0.service_infos.0.name":              acctest.CHECKSET,
			"service_instances.0.service.0.service_infos.0.short_description": acctest.CHECKSET,
			"service_instances.0.service.0.service_infos.0.image":             acctest.CHECKSET,
			"service_instances.0.service.0.service_infos.0.locale":            acctest.CHECKSET,
			"service_instances.0.service.0.status":                            acctest.CHECKSET,
			"service_instances.0.tags.%":                                      "2",
			"service_instances.0.tags.Created":                                "TF",
			"service_instances.0.tags.For":                                    "ServiceInstance",
			"service_instances.0.status":                                      "Deployed",
		}
	}
	var fakeServiceInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "0",
			"names.#":             "0",
			"service_instances.#": "0",
		}
	}
	// checkDataSourceAttr is an unexported method of acctest.DataSourceAttr and
	// cannot be called from this computenest_test package, so build the exist
	// and empty ResourceAttrMapUpdate check funcs directly with the exported
	// acctest.ResourceAttrInit(...).ResourceAttrMapUpdateSet() chain (the same
	// chain checkDataSourceAttr uses internally), keeping the ExistMapFunc /
	// FakeMapFunc assertions while staying in the standard-template structure
	// the TestingCoverageRate checker parses.
	testAccCheckExist := acctest.ResourceAttrInit(resourceId, existServiceInstancesDataSourceNameMapFunc(rand)).ResourceAttrMapUpdateSet()
	testAccCheckEmpty := acctest.ResourceAttrInit(resourceId, fakeServiceInstancesDataSourceNameMapFunc(rand)).ResourceAttrMapUpdateSet()
	preCheck := func() {
		acctest.PreCheck(t)
		// ComputeNest is Hangzhou-only (connectivity.ComputeNestSupportRegions);
		// the flat suite gate checkoutSupportedRegions converged to cn-hangzhou,
		// so pin it here to keep config and execution region consistent.
		os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 preCheck,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// ids: Exist (matches the created instance) -> Fake (no match).
			{
				Config: testAccConfig(map[string]interface{}{
					"ids": []string{"${alicloud_compute_nest_service_instance.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExist(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ids": []string{"${alicloud_compute_nest_service_instance.default.id}_fake"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmpty(nil),
				),
			},
			// name_regex: Exist (real name) -> Fake (name with _fake suffix).
			{
				Config: testAccConfig(map[string]interface{}{
					"name_regex": "${alicloud_compute_nest_service_instance.default.service_instance_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExist(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name_regex": "${alicloud_compute_nest_service_instance.default.service_instance_name}_fake",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmpty(nil),
				),
			},
			// status: Exist (Deployed, narrowed by name_regex) -> Fake (DeletedFailed).
			{
				Config: testAccConfig(map[string]interface{}{
					"status":     "Deployed",
					"name_regex": "${alicloud_compute_nest_service_instance.default.service_instance_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExist(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DeletedFailed",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmpty(nil),
				),
			},
			// filter: Exist (ServiceInstanceId == created id) -> Fake (no match).
			{
				Config: testAccConfig(map[string]interface{}{
					"filter": []map[string]interface{}{
						{
							"name":  "ServiceInstanceId",
							"value": []string{"${alicloud_compute_nest_service_instance.default.id}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExist(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"filter": []map[string]interface{}{
						{
							"name":  "ServiceInstanceId",
							"value": []string{"${alicloud_compute_nest_service_instance.default.id}_fake"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmpty(nil),
				),
			},
			// tags: Exist (Created=TF/For=ServiceInstance) -> Fake (For=ServiceInstance_Update).
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServiceInstance",
					},
					"name_regex": "${alicloud_compute_nest_service_instance.default.service_instance_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExist(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServiceInstance_Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmpty(nil),
				),
			},
			// all filters together, plus output_file (covers the last Optional
			// argument the previous wrapper-template test never set).
			{
				Config: testAccConfig(map[string]interface{}{
					"ids":        []string{"${alicloud_compute_nest_service_instance.default.id}"},
					"name_regex": "${alicloud_compute_nest_service_instance.default.service_instance_name}",
					"status":     "Deployed",
					"filter": []map[string]interface{}{
						{
							"name":  "ServiceInstanceId",
							"value": []string{"${alicloud_compute_nest_service_instance.default.id}"},
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServiceInstance",
					},
					"output_file": "/tmp/computenest_service_instances_data_source.json",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExist(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ids":        []string{"${alicloud_compute_nest_service_instance.default.id}_fake"},
					"name_regex": "${alicloud_compute_nest_service_instance.default.service_instance_name}_fake",
					"status":     "DeletedFailed",
					"filter": []map[string]interface{}{
						{
							"name":  "ServiceInstanceId",
							"value": []string{"${alicloud_compute_nest_service_instance.default.id}_fake"},
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServiceInstance_Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmpty(nil),
				),
			},
		},
	})
}

func dataSourceServiceInstancesConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
  		available_instance_type     = "ecs.g6.large"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.g6"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.0.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		instance_name              = var.name
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_zones.default.zones.0.id
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = alicloud_vswitch.default.id
	}

	resource "alicloud_compute_nest_service_instance" "default" {
  		service_id            = "service-dd475e6e468348799f0f"
  		service_version       = "1"
  		service_instance_name = var.name
  		resource_group_id     = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  		payment_type          = "Permanent"
  		operation_metadata {
    		operation_start_time = "1681281179000"
    		operation_end_time   = "1681367579000"
    		resources            = "{\"Type\":\"ResourceIds\",\"ResourceIds\":{\"ALIYUN::ECS::INSTANCE\":[\"${alicloud_instance.default.id}\"]},\"RegionId\":\"cn-hangzhou\"}"
  		}
  		tags = {
    		Created = "TF"
    		For     = "ServiceInstance"
  		}
	}
`, name)
}
