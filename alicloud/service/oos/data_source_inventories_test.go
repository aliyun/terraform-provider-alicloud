package oos_test

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudOosInventories_basic0(t *testing.T) {
	dataSourceId := "data.alicloud_oos_inventories.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOosInventoriesConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceId, "id"),
					resource.TestCheckResourceAttrSet(dataSourceId, "next_token_result"),
					testAccCheckOosInventoriesEntityCount(dataSourceId, 2),
					testAccCheckOosInventoriesEntities(dataSourceId),
				),
			},
			{
				Config: testAccOosInventoriesConfigNextToken(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_oos_inventories.first", "entities.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_oos_inventories.first", "next_token_result"),
					resource.TestCheckResourceAttr(dataSourceId, "entities.#", "1"),
					testAccCheckOosInventoriesNextTokenPage("data.alicloud_oos_inventories.first", dataSourceId),
				),
			},
			{
				Config: testAccOosInventoriesConfigFilteredMiss(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "entities.#", "0"),
				),
			},
			{
				Config: testAccOosInventoriesConfigFilteredMatch(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOosInventoriesEntityCount(dataSourceId, 1),
					testAccCheckOosInventoriesMatchesSeed(dataSourceId, "data.alicloud_oos_inventories.seed"),
				),
			},
			{
				Config: testAccOosInventoriesConfigAggregator(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "aggregator.#", "1"),
					resource.TestCheckResourceAttr(dataSourceId, "aggregator.0", "ACS:Application.Name"),
					resource.TestCheckResourceAttr(dataSourceId, "entities.#", "1"),
					testAccCheckOosInventoriesAggregation(dataSourceId),
				),
			},
			{
				Config: testAccOosInventoriesConfigPinnedRegion(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceId, "region_id"),
					testAccCheckOosInventoriesEntityCount(dataSourceId, 1),
					testAccCheckOosInventoriesEntities(dataSourceId),
				),
			},
			{
				Config:      testAccOosInventoriesConfigInvalidOperator(),
				ExpectError: regexp.MustCompile("invalid filter rule"),
			},
		},
	})
}

// testAccCheckOosInventoriesEntityCount requires the query to return at least min
// entities. The unfiltered query is expected to page through the inventory of the
// managed instances that the acceptance account already owns.
func testAccCheckOosInventoriesEntityCount(name string, min int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		count, err := oosInventoriesEntityCount(s, name)
		if err != nil {
			return err
		}
		if count < min {
			return fmt.Errorf("%s: entities.# is %d, expected at least %d", name, count, min)
		}
		return nil
	}
}

// testAccCheckOosInventoriesEntities verifies that every returned entity is a JSON
// object that carries its identifier and at least one dynamic component key, which
// is the shape a query without aggregator returns.
func testAccCheckOosInventoriesEntities(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		entities, err := oosInventoriesEntities(s, name)
		if err != nil {
			return err
		}
		if len(entities) == 0 {
			return fmt.Errorf("%s: entities is empty", name)
		}
		for i, entity := range entities {
			id, _ := entity["Id"].(string)
			if id == "" {
				return fmt.Errorf("%s: entities.%d has no Id", name, i)
			}
			components := 0
			for key := range entity {
				if key != "Id" {
					components++
				}
			}
			if components == 0 {
				return fmt.Errorf("%s: entities.%d has no component key besides Id", name, i)
			}
		}
		return nil
	}
}

// testAccCheckOosInventoriesNextTokenPage verifies that feeding next_token_result
// back as next_token reads a different entity instead of repeating the first one.
// Entities of the same instance share the Id key, so the whole entity is compared.
func testAccCheckOosInventoriesNextTokenPage(firstName, nextName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		first, err := oosInventoriesRawEntities(s, firstName)
		if err != nil {
			return err
		}
		next, err := oosInventoriesRawEntities(s, nextName)
		if err != nil {
			return err
		}
		if len(first) != 1 {
			return fmt.Errorf("%s: entities.# is %d, expected exactly 1", firstName, len(first))
		}
		if len(next) != 1 {
			return fmt.Errorf("%s: entities.# is %d, expected exactly 1", nextName, len(next))
		}
		if first[0] == next[0] {
			return fmt.Errorf("%s and %s returned the same entity; next_token did not advance the page", firstName, nextName)
		}
		return nil
	}
}

// testAccCheckOosInventoriesMatchesSeed verifies that filtering on the instance ID
// of a known entity returns entities of that instance only, so the filter rule
// changes the result set instead of being ignored by the service.
func testAccCheckOosInventoriesMatchesSeed(name, seedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		seedID, err := oosInventoriesSingleID(s, seedName)
		if err != nil {
			return err
		}
		entities, err := oosInventoriesEntities(s, name)
		if err != nil {
			return err
		}
		if len(entities) == 0 {
			return fmt.Errorf("%s: the filtered query returned no entity for instance %s", name, seedID)
		}
		for i, entity := range entities {
			if id, _ := entity["Id"].(string); id != seedID {
				return fmt.Errorf("%s: entities.%d is %q, expected the filtered instance %q", name, i, id, seedID)
			}
		}
		return nil
	}
}

// testAccCheckOosInventoriesAggregation verifies that a query with aggregator
// returns the aggregated content and that the aggregated counts are readable.
func testAccCheckOosInventoriesAggregation(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		entities, err := oosInventoriesEntities(s, name)
		if err != nil {
			return err
		}
		if len(entities) != 1 {
			return fmt.Errorf("%s: entities.# is %d, expected the single aggregated entity", name, len(entities))
		}
		data, ok := entities[0]["Data"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("%s: entities.0 has no Data object", name)
		}
		components := 0
		for component, raw := range data {
			body, ok := raw.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%s: entities.0.Data.%s is not an object", name, component)
			}
			content, ok := body["Content"].([]interface{})
			if !ok || len(content) == 0 {
				return fmt.Errorf("%s: entities.0.Data.%s has no Content entries", name, component)
			}
			components++
			for i, rawEntry := range content {
				entry, ok := rawEntry.(map[string]interface{})
				if !ok {
					return fmt.Errorf("%s: entities.0.Data.%s.Content.%d is not an object", name, component, i)
				}
				count, err := oosInventoriesInt(entry["Count"])
				if err != nil {
					return fmt.Errorf("%s: entities.0.Data.%s.Content.%d: %v", name, component, i, err)
				}
				if count <= 0 {
					return fmt.Errorf("%s: entities.0.Data.%s.Content.%d has Count %d", name, component, i, count)
				}
			}
		}
		if components == 0 {
			return fmt.Errorf("%s: entities.0.Data has no component", name)
		}
		return nil
	}
}

func oosInventoriesEntityCount(s *terraform.State, name string) (int, error) {
	entities, err := oosInventoriesEntities(s, name)
	if err != nil {
		return 0, err
	}
	return len(entities), nil
}

func oosInventoriesSingleID(s *terraform.State, name string) (string, error) {
	entities, err := oosInventoriesEntities(s, name)
	if err != nil {
		return "", err
	}
	if len(entities) != 1 {
		return "", fmt.Errorf("%s: entities.# is %d, expected exactly 1", name, len(entities))
	}
	id, _ := entities[0]["Id"].(string)
	if id == "" {
		return "", fmt.Errorf("%s: entities.0 has no Id", name)
	}
	return id, nil
}

func oosInventoriesEntities(s *terraform.State, name string) ([]map[string]interface{}, error) {
	rawEntities, err := oosInventoriesRawEntities(s, name)
	if err != nil {
		return nil, err
	}
	entities := make([]map[string]interface{}, 0, len(rawEntities))
	for i, value := range rawEntities {
		var entity map[string]interface{}
		if err := json.Unmarshal([]byte(value), &entity); err != nil {
			return nil, fmt.Errorf("%s: entities.%d is not a JSON object: %v", name, i, err)
		}
		if len(entity) == 0 {
			return nil, fmt.Errorf("%s: entities.%d is an empty JSON object", name, i)
		}
		entities = append(entities, entity)
	}
	return entities, nil
}

func oosInventoriesRawEntities(s *terraform.State, name string) ([]string, error) {
	rs, ok := s.RootModule().Resources[name]
	if !ok {
		return nil, fmt.Errorf("data source %s is not found in state", name)
	}
	rawCount, ok := rs.Primary.Attributes["entities.#"]
	if !ok {
		return nil, fmt.Errorf("%s: entities.# is missing from state", name)
	}
	count, err := strconv.Atoi(rawCount)
	if err != nil {
		return nil, fmt.Errorf("%s: entities.# is not a number: %v", name, err)
	}
	entities := make([]string, 0, count)
	for i := 0; i < count; i++ {
		entities = append(entities, rs.Primary.Attributes[fmt.Sprintf("entities.%d", i)])
	}
	return entities, nil
}

func oosInventoriesInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case float64:
		return int(v), nil
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("Count %q is not an integer", v)
		}
		return parsed, nil
	case nil:
		return 0, fmt.Errorf("Count is missing")
	default:
		return 0, fmt.Errorf("Count %v is not an integer", v)
	}
}

func testAccOosInventoriesConfigBasic() string {
	return `
data "alicloud_oos_inventories" "default" {
  max_results = 100
}
`
}

func testAccOosInventoriesConfigNextToken() string {
	return `
data "alicloud_oos_inventories" "first" {
  max_results = 1
}

data "alicloud_oos_inventories" "default" {
  max_results = 1
  next_token  = data.alicloud_oos_inventories.first.next_token_result
}
`
}

func testAccOosInventoriesConfigPinnedRegion() string {
	return `
data "alicloud_regions" "default" {
  current = true
}

data "alicloud_oos_inventories" "default" {
  region_id   = data.alicloud_regions.default.regions.0.id
  max_results = 1
}
`
}

func testAccOosInventoriesConfigFilteredMiss() string {
	return `
data "alicloud_oos_inventories" "default" {
  filter = [
    {
      name     = "ACS:InstanceInformation.InstanceId"
      operator = "Equal"
      value    = ["i-not-exist000000000000"]
    }
  ]
}
`
}

func testAccOosInventoriesConfigFilteredMatch() string {
	return `
data "alicloud_oos_inventories" "seed" {
  max_results = 1
}

data "alicloud_oos_inventories" "default" {
  filter = [
    {
      name     = "ACS:InstanceInformation.InstanceId"
      operator = "Equal"
      value    = [jsondecode(data.alicloud_oos_inventories.seed.entities[0]).Id]
    }
  ]
}
`
}

func testAccOosInventoriesConfigAggregator() string {
	return `
data "alicloud_oos_inventories" "default" {
  aggregator = ["ACS:Application.Name"]
}
`
}

func testAccOosInventoriesConfigInvalidOperator() string {
	return `
data "alicloud_oos_inventories" "default" {
  filter = [
    {
      name     = "ACS:InstanceInformation.InstanceId"
      operator = "Equals"
      value    = ["i-not-exist000000000000"]
    }
  ]
}
`
}
