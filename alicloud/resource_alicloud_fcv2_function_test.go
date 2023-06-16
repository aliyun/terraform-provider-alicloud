package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Fcv2 Function. >>> Resource test cases, automatically generated.
// Case 3393
func TestAccAlicloudFcv2Function_basic3393(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv2_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv2FunctionMap3393)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv2ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv2Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv2function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv2FunctionBasicDependence3393)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": name,
					"runtime":       "custom.debian10",
					"service_name":  "terraform-e2e-service-base",
					"handler":       "index.handler",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": name,
						"runtime":       "custom.debian10",
						"service_name":  "terraform-e2e-service-base",
						"handler":       "index.handler",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "1024",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform测试case",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform测试case",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initializer": "index.initializer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initializer": "index.initializer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initialization_timeout": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initialization_timeout": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "e1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "e1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "index.prefreeze",
									"timeout": "30",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "index.prestop",
									"timeout": "30",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"code": []map[string]interface{}{
						{
							"zip_file": "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",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_dns": []map[string]interface{}{
						{
							"name_servers": []string{
								"223.5.5.5"},
							"searches": []string{
								"mydomain.com"},
							"dns_options": []map[string]interface{}{
								{
									"name":  name + "_update",
									"value": "1",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size": "512",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "512",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_concurrency": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_concurrency": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"layers": []string{
						"d3fc5de8d120687be2bfab761518d5de#Nodejs-Aliyun-SDK#2", "d3fc5de8d120687be2bfab761518d5de#Python39#2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"layers.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_health_check_config": []map[string]interface{}{
						{
							"http_get_url":          "/healthcheck",
							"initial_delay_seconds": "3",
							"period_seconds":        "3",
							"timeout_seconds":       "3",
							"failure_threshold":     "1",
							"success_threshold":     "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ca_port": "9000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ca_port": "9000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_runtime_config": []map[string]interface{}{
						{
							"command": []string{
								"npm"},
							"args": []string{
								"run", "start"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": "custom.debian10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime": "custom.debian10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"handler": "index.handler",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"handler": "index.handler",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "2048",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "2048",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": "custom",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime": "custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform测试case update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform测试case update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initializer": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initializer": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initialization_timeout": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initialization_timeout": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "c1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "c1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "true",
									"timeout": "3",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "true",
									"timeout": "3",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"code": []map[string]interface{}{
						{
							"zip_file": "UEsDBAoAAAAAAOt8v1YAAAAAAAAAAAAAAAAFABwAY29kZS9VVAkAA6r5dmTE+XZkdXgLAAEE9gEAAAQUAAAAUEsDBBQAAAAIAOt8v1Y14c8e4QAAAH8BAAANABwAY29kZS9pbmRleC5qc1VUCQADqvl2ZKv5dmR1eAsAAQT2AQAABBQAAAB1kFFOwzAMht9zCr+lnapGPDAQ1XYK3lGamiXCjavGYQO0i3AW7sQVSDU0oU28Wf4/+7Osc0JIMgcnulNmpR4ZMNqeEMQjhBgkWArvOMMzWskzQuVFpvRgjEea2hK+5dg6Hs3ALo8Y5WlAsYHMze36/m7dehmpVhOhLaowlmKBrtfn6CRwBJugR+L999enwsPEs6T2L7iBynEUPEgDzhL11r3UsNnChwIoSWLClnhX6fNUiDtdd0v8y1cxEzWgl+6xUytzNnkbBzpZ8LXc2cD/sgtd+Qcx7Hmm4WS79l0gxx9QSwECHgMKAAAAAADrfL9WAAAAAAAAAAAAAAAABQAYAAAAAAAAABAA7UEAAAAAY29kZS9VVAUAA6r5dmR1eAsAAQT2AQAABBQAAABQSwECHgMUAAAACADrfL9WNeHPHuEAAAB/AQAADQAYAAAAAAABAAAApIE/AAAAY29kZS9pbmRleC5qc1VUBQADqvl2ZHV4CwABBPYBAAAEFAAAAFBLBQYAAAAAAgACAJ4AAABnAQAAAAA=",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_dns": []map[string]interface{}{
						{
							"name_servers": []string{
								"223.6.6.6"},
							"searches": []string{
								"newdomain.com"},
							"dns_options": []map[string]interface{}{
								{
									"name":  name + "_update",
									"value": "2",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size": "10240",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "10240",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_concurrency": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_concurrency": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"layers": []string{
						"d3fc5de8d120687be2bfab761518d5de#Python39#2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"layers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_health_check_config": []map[string]interface{}{
						{
							"http_get_url":          "/healthcheck1",
							"initial_delay_seconds": "5",
							"period_seconds":        "5",
							"timeout_seconds":       "2",
							"failure_threshold":     "5",
							"success_threshold":     "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ca_port": "8000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ca_port": "8000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_runtime_config": []map[string]interface{}{
						{
							"command": []string{},
							"args":    []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":          name + "_update",
					"memory_size":            "1024",
					"runtime":                "custom.debian10",
					"description":            "terraform测试case",
					"service_name":           "terraform-e2e-service-base",
					"initializer":            "index.initializer",
					"initialization_timeout": "10",
					"timeout":                "60",
					"handler":                "index.handler",
					"instance_type":          "e1",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "index.prefreeze",
									"timeout": "30",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "index.prestop",
									"timeout": "30",
								},
							},
						},
					},
					"code": []map[string]interface{}{
						{
							"zip_file": "UEsDBBQAAAAIAEqkA1VQpdMOPwUAANQOAAAHABwAbWFpbi50ZlVUCQADS2vqYkxr6mJ1eAsAAQT2AQAABBQAAADlV91v2zYQf89fISh72IDqy1YS22gfMqfFVmxLtrVphzYQKIqymUiiQlJ2nMD/+46kKMmOO2BD97IpcCQd746/++RJEs5RznjpPB05Dif3DeUkS2rOVjQjXGiy46CC4oI1mfOqJTiOYA3HxAGKC6ubpgosk9tyrECeskpxRH4UR35oVrZH6rc9sptoBUYStB/r/TAmQiR3ZKOkf3p3/uOJnE7KPzbX02V9evWGnqPb32AjzSwI5kRaZrQer5t5HX9/9/7y9Ql/8zGtKjp++8vHePU251rkWFm6UNDUpWRqD8yRS4KE9CLgOT4M4eLt/fION1fzdZ3S+YcPm3eGdxcBqc6vpz9jdPH6/uHhh8n1Pf/1Qvx+dTq6qOa40Aja7YEZV94SVYvHJWtc8MlRwTAqjNszkqOmkElGBOa0lhowiJw3kjmwIZIkc9IN7M7B1QWgBZGVcNZULh1pA+sONFWoJI69QFMv6SlJL8dwL9lQZFXjBGeUW5FoOvKj04kf+mEQne5wCtgYLw33gBPiHozifeNoteDKvwpsgmnGkxRW7wRIftJZYvMs56xMasal2t6LWqpkLW2HChklGWZFG9WisKnYb3DYBiNMGadyo4WjMDS5+uIrgaG4rA+jCX39F4T/GEgcjw8h6cm7UCT+15BMwkNAOurXx3Gj0goqaillLWZBoApLSL7xu/z3KQu6dhbstaqggCISMsgYFiBrepoIcpzkTYVVxR0dW3LfphJYV6VDMXFNEekn3bt0iYENS1IUzFszXmS25xRskWBW5XRhWI9bl9wSrJ3XqVeMLd23d6W3EwIGIRkn+0Ka6Jv/A4Gtbjms0MAQFjOOytksmkTjs7N4PJlMT8ZnYThTHK2DctxWtaIBfqWAVuDTCtqc6YygK4ditnvoPrFvHPTFRsUtWXDW1AnNhoB3F32x8GnWSWp9wD+4BqKw6KvfQMB2H5qZFtIpcgZihseHO4i+6Hhuej9tla2HI24zAkLeP2pb2wzYR9lniW/vNijDRqwqAXLQqzdyCSrVck4LYllgOXgvVO6ukeqZVYZUkFKUoqDv38FOK68LtilJJT1SLWhFAswy4j/SWisvocPzTSLoo86HaDQxR1JTSVpaI9wFizQZDqesIF37LxE1EBUvHJktfRwq133Jc5LTxYJwcFz3pP3W5otR8fLl68s35kB/cskK0At35nxymRCzy1TVwNwcebOrRhqCe/MCYkELyEtgfXLh/NX3mpOcPsCjG7SPwCeavCVqT2y3EGq9o3Z4G9AvRmTn5FQMAMuz1nTlNWQY1Nn45Gw6jeM4HI3i8WhYZ6BFm1oxSXOKkcLQltxuWjmHeorxdYJ4NdhT+WswVDzfXp3wXtrgOyK90F4mppv6mZGuDuzhuBZoQ3jSjnkQXf3u6pPBLNlmqLypejlmZQ02pgVJbLpBrboVpOetiEbujebJiB0DYf/EIO1UwYiitMH44GUnaUTGWUTi8QkKUe72QkznRyfU579igIdEVZjqYHBLkSCn8bfuN0+AbenjdbbtC+a7dlb9+01Be0D3oi+6zNdvPsTvZlhqAHinIu27zQfz/hed6lm9eeSBap/9BwqvImvv/1V8K8QpgpqBIbJEC+K2Hwf9RwEwzxs49UsVWAkdGpJIs/oq4ZXiVzATQfj7aV3JdLMSIE4JvQUG33gGCx9KFcYgD2u9XqdXfx8Ea5J6OfiXgDvuvOwWzGUzNbNFeszvAYOWErLaQNYWvnIOQfn02fWDlDEJi6j+7N7s60F6nNzRUzVlSviOnik482j7J1BLAQIeAxQAAAAIAEqkA1VQpdMOPwUAANQOAAAHABgAAAAAAAEAAACAgQAAAABtYWluLnRmVVQFAANLa+pidXgLAAEE9gEAAAQUAAAAUEsFBgAAAAABAAEATQAAAIAFAAAAAA==",
						},
					},
					"environment_variables": []map[string]interface{}{
						{},
					},
					"custom_dns": []map[string]interface{}{
						{
							"name_servers": []string{
								"223.5.5.5"},
							"searches": []string{
								"mydomain.com"},
							"dns_options": []map[string]interface{}{
								{
									"name":  name + "_update",
									"value": "1",
								},
							},
						},
					},
					"disk_size":            "512",
					"instance_concurrency": "10",
					"layers": []string{
						"d3fc5de8d120687be2bfab761518d5de#Nodejs-Aliyun-SDK#2", "d3fc5de8d120687be2bfab761518d5de#Python39#2"},
					"cpu": "1",
					"custom_health_check_config": []map[string]interface{}{
						{
							"http_get_url":          "/healthcheck",
							"initial_delay_seconds": "3",
							"period_seconds":        "3",
							"timeout_seconds":       "3",
							"failure_threshold":     "1",
							"success_threshold":     "1",
						},
					},
					"ca_port": "9000",
					"custom_runtime_config": []map[string]interface{}{
						{
							"command": []string{
								"npm"},
							"args": []string{
								"run", "start"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":          name + "_update",
						"memory_size":            "1024",
						"runtime":                "custom.debian10",
						"description":            "terraform测试case",
						"service_name":           "terraform-e2e-service-base",
						"initializer":            "index.initializer",
						"initialization_timeout": "10",
						"timeout":                "60",
						"handler":                "index.handler",
						"instance_type":          "e1",
						"disk_size":              "512",
						"instance_concurrency":   "10",
						"layers.#":               "2",
						"cpu":                    "1",
						"ca_port":                "9000",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudFcv2FunctionMap3393 = map[string]string{
	"ca_port":                CHECKSET,
	"instance_concurrency":   CHECKSET,
	"instance_type":          CHECKSET,
	"memory_size":            CHECKSET,
	"timeout":                CHECKSET,
	"create_time":            CHECKSET,
	"initialization_timeout": CHECKSET,
	"code_checksum":          CHECKSET,
}

func AlicloudFcv2FunctionBasicDependence3393(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3270
func TestAccAlicloudFcv2Function_basic3270(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv2_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv2FunctionMap3270)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv2ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv2Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv2function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv2FunctionBasicDependence3270)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": name,
					"runtime":       "nodejs14",
					"service_name":  "terraform-e2e-service-base",
					"handler":       "index.handler",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": name,
						"runtime":       "nodejs14",
						"service_name":  "terraform-e2e-service-base",
						"handler":       "index.handler",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "1024",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform测试case",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform测试case",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initializer": "index.initializer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initializer": "index.initializer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initialization_timeout": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initialization_timeout": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "e1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "e1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "index.prefreeze",
									"timeout": "3",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "index.prestop",
									"timeout": "3",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "code-sample-cn-beijing",
							"oss_object_name": "quick-start-sample-codes/quick-start-sample-codes-nodejs/RocketMQ-producer-nodejs14-event/code.zip",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size": "512",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "512",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_concurrency": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_concurrency": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"layers": []string{
						"acs:fc:cn-beijing:official:layers/Nodejs-Puppeteer17x/versions/3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"layers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": "nodejs14",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime": "nodejs14",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"handler": "index.handler",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"handler": "index.handler",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": "nodejs16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime": "nodejs16",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform测试case update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform测试case update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initializer": "index. initializer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initializer": "index. initializer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initialization_timeout": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initialization_timeout": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "c1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "c1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"timeout": "10",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "index.presstop",
									"timeout": "10",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"code": []map[string]interface{}{
						{
							"oss_object_name": "quick-start-sample-codes/quick-start-sample-codes-nodejs/cdn-trigger-nodejs/code.zip",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size": "10240",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "10240",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_concurrency": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_concurrency": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"layers": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"layers.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":          name + "_update",
					"memory_size":            "1024",
					"runtime":                "nodejs14",
					"description":            "terraform测试case",
					"service_name":           "terraform-e2e-service-base",
					"initializer":            "index.initializer",
					"initialization_timeout": "10",
					"timeout":                "60",
					"handler":                "index.handler",
					"instance_type":          "e1",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "index.prefreeze",
									"timeout": "3",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "index.prestop",
									"timeout": "3",
								},
							},
						},
					},
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "code-sample-cn-beijing",
							"oss_object_name": "quick-start-sample-codes/quick-start-sample-codes-nodejs/RocketMQ-producer-nodejs14-event/code.zip",
						},
					},
					"environment_variables": []map[string]interface{}{
						{},
					},
					"disk_size":            "512",
					"instance_concurrency": "10",
					"layers": []string{
						"acs:fc:cn-beijing:official:layers/Nodejs-Puppeteer17x/versions/3"},
					"cpu": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":          name + "_update",
						"memory_size":            "1024",
						"runtime":                "nodejs14",
						"description":            "terraform测试case",
						"service_name":           "terraform-e2e-service-base",
						"initializer":            "index.initializer",
						"initialization_timeout": "10",
						"timeout":                "60",
						"handler":                "index.handler",
						"instance_type":          "e1",
						"disk_size":              "512",
						"instance_concurrency":   "10",
						"layers.#":               "1",
						"cpu":                    "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudFcv2FunctionMap3270 = map[string]string{
	"ca_port":                CHECKSET,
	"instance_concurrency":   CHECKSET,
	"instance_type":          CHECKSET,
	"memory_size":            CHECKSET,
	"timeout":                CHECKSET,
	"create_time":            CHECKSET,
	"initialization_timeout": CHECKSET,
	"code_checksum":          CHECKSET,
}

func AlicloudFcv2FunctionBasicDependence3270(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3395
func TestAccAlicloudFcv2Function_basic3395(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv2_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv2FunctionMap3395)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv2ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv2Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv2function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv2FunctionBasicDependence3395)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": name,
					"runtime":       "custom-container",
					"service_name":  "terraform-e2e-service-base",
					"handler":       "index.handler",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": name,
						"runtime":       "custom-container",
						"service_name":  "terraform-e2e-service-base",
						"handler":       "index.handler",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "8192",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "8192",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform测试case",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform测试case",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initializer": "index.initializer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initializer": "index.initializer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initialization_timeout": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initialization_timeout": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "fc.gpu.tesla.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "fc.gpu.tesla.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "index.prefreeze",
									"timeout": "30",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "index.prestop",
									"timeout": "30",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size": "512",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "512",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_concurrency": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_concurrency": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_health_check_config": []map[string]interface{}{
						{
							"http_get_url":          "/healthcheck",
							"initial_delay_seconds": "3",
							"period_seconds":        "3",
							"timeout_seconds":       "3",
							"failure_threshold":     "1",
							"success_threshold":     "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ca_port": "9000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ca_port": "9000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_container_config": []map[string]interface{}{
						{
							"args":              "[\\\"--debug\\\"]",
							"command":           "[\\\"python\\\",\\\"app.py\\\"]",
							"image":             "registry-vpc.cn-beijing.aliyuncs.com/fc-stable-diffusion/stable-diffusion:v1",
							"acceleration_type": "Default",
							"web_server_mode":   "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gpu_memory_size": "4096",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gpu_memory_size": "4096",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": "custom-container",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime": "custom-container",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"handler": "index.handler",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"handler": "index.handler",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform测试case update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform测试case update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initializer": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initializer": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initialization_timeout": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initialization_timeout": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "true",
									"timeout": "3",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "true",
									"timeout": "3",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size": "10240",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "10240",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_concurrency": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_concurrency": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_health_check_config": []map[string]interface{}{
						{
							"http_get_url":          "/healthcheck1",
							"initial_delay_seconds": "5",
							"period_seconds":        "5",
							"timeout_seconds":       "2",
							"failure_threshold":     "5",
							"success_threshold":     "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ca_port": "8000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ca_port": "8000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_container_config": []map[string]interface{}{
						{
							"args":              "[\\\"--use-local\\\"]",
							"command":           "[\\\"s\\\",\\\"deploy\\\"]",
							"image":             "registry.cn-beijing.aliyuncs.com/serverlessdevshanxie/sd-auto-nas:v3",
							"acceleration_type": "None",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gpu_memory_size": "8192",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gpu_memory_size": "8192",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":          name + "_update",
					"memory_size":            "8192",
					"runtime":                "custom-container",
					"description":            "terraform测试case",
					"service_name":           "terraform-e2e-service-base",
					"initializer":            "index.initializer",
					"initialization_timeout": "10",
					"timeout":                "60",
					"handler":                "index.handler",
					"instance_type":          "fc.gpu.tesla.1",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "index.prefreeze",
									"timeout": "30",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "index.prestop",
									"timeout": "30",
								},
							},
						},
					},
					"environment_variables": []map[string]interface{}{
						{},
					},
					"disk_size":            "512",
					"instance_concurrency": "10",
					"cpu":                  "2",
					"custom_health_check_config": []map[string]interface{}{
						{
							"http_get_url":          "/healthcheck",
							"initial_delay_seconds": "3",
							"period_seconds":        "3",
							"timeout_seconds":       "3",
							"failure_threshold":     "1",
							"success_threshold":     "1",
						},
					},
					"ca_port": "9000",
					"custom_container_config": []map[string]interface{}{
						{
							"args":              "[\\\"--debug\\\"]",
							"command":           "[\\\"python\\\",\\\"app.py\\\"]",
							"image":             "registry-vpc.cn-beijing.aliyuncs.com/fc-stable-diffusion/stable-diffusion:v1",
							"acceleration_type": "Default",
							"web_server_mode":   "true",
						},
					},
					"gpu_memory_size": "4096",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":          name + "_update",
						"memory_size":            "8192",
						"runtime":                "custom-container",
						"description":            "terraform测试case",
						"service_name":           "terraform-e2e-service-base",
						"initializer":            "index.initializer",
						"initialization_timeout": "10",
						"timeout":                "60",
						"handler":                "index.handler",
						"instance_type":          "fc.gpu.tesla.1",
						"disk_size":              "512",
						"instance_concurrency":   "10",
						"cpu":                    "2",
						"ca_port":                "9000",
						"gpu_memory_size":        "4096",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudFcv2FunctionMap3395 = map[string]string{
	"ca_port":                CHECKSET,
	"instance_concurrency":   CHECKSET,
	"instance_type":          CHECKSET,
	"memory_size":            CHECKSET,
	"timeout":                CHECKSET,
	"create_time":            CHECKSET,
	"initialization_timeout": CHECKSET,
	"code_checksum":          CHECKSET,
}

func AlicloudFcv2FunctionBasicDependence3395(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3411
func TestAccAlicloudFcv2Function_basic3411(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv2_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv2FunctionMap3411)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv2ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv2Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv2function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv2FunctionBasicDependence3411)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": name,
					"runtime":       "custom-container",
					"service_name":  "terraform-e2e-service-base",
					"handler":       "index.handler",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": name,
						"runtime":       "custom-container",
						"service_name":  "terraform-e2e-service-base",
						"handler":       "index.handler",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory_size": "8192",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory_size": "8192",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform测试case",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform测试case",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initializer": "index.initializer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initializer": "index.initializer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initialization_timeout": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initialization_timeout": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "fc.gpu.tesla.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "fc.gpu.tesla.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "index.prefreeze",
									"timeout": "30",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "index.prestop",
									"timeout": "30",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size": "512",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "512",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_concurrency": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_concurrency": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_health_check_config": []map[string]interface{}{
						{
							"http_get_url":          "/healthcheck",
							"initial_delay_seconds": "3",
							"period_seconds":        "3",
							"timeout_seconds":       "3",
							"failure_threshold":     "1",
							"success_threshold":     "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ca_port": "9000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ca_port": "9000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_container_config": []map[string]interface{}{
						{
							"args":              "[\\\"--debug\\\"]",
							"command":           "[\\\"python\\\",\\\"app.py\\\"]",
							"image":             "terraform-e2e-test-registry-vpc.cn-beijing.cr.aliyuncs.com/terraform-e2e/terraform-e2e:v1",
							"acceleration_type": "Default",
							"web_server_mode":   "true",
							"instance_id":       "cri-35n95x367lp1970t",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gpu_memory_size": "4096",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gpu_memory_size": "4096",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": "custom-container",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runtime": "custom-container",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"handler": "index.handler",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"handler": "index.handler",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform测试case update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform测试case update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initializer": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initializer": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"initialization_timeout": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"initialization_timeout": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "true",
									"timeout": "3",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "true",
									"timeout": "3",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size": "10240",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "10240",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_concurrency": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_concurrency": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_health_check_config": []map[string]interface{}{
						{
							"http_get_url":          "/healthcheck1",
							"initial_delay_seconds": "5",
							"period_seconds":        "5",
							"timeout_seconds":       "2",
							"failure_threshold":     "5",
							"success_threshold":     "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ca_port": "8000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ca_port": "8000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_container_config": []map[string]interface{}{
						{
							"args":    "[\\\"--use-local\\\"]",
							"command": "[\\\"s\\\",\\\"deploy\\\"]",
							"image":   "terraform-e2e-test-registry-vpc.cn-beijing.cr.aliyuncs.com/terraform-e2e/terraform-e2e:v2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gpu_memory_size": "8192",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gpu_memory_size": "8192",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":          name + "_update",
					"memory_size":            "8192",
					"runtime":                "custom-container",
					"description":            "terraform测试case",
					"service_name":           "terraform-e2e-service-base",
					"initializer":            "index.initializer",
					"initialization_timeout": "10",
					"timeout":                "60",
					"handler":                "index.handler",
					"instance_type":          "fc.gpu.tesla.1",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "index.prefreeze",
									"timeout": "30",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "index.prestop",
									"timeout": "30",
								},
							},
						},
					},
					"environment_variables": []map[string]interface{}{
						{},
					},
					"disk_size":            "512",
					"instance_concurrency": "10",
					"cpu":                  "2",
					"custom_health_check_config": []map[string]interface{}{
						{
							"http_get_url":          "/healthcheck",
							"initial_delay_seconds": "3",
							"period_seconds":        "3",
							"timeout_seconds":       "3",
							"failure_threshold":     "1",
							"success_threshold":     "1",
						},
					},
					"ca_port": "9000",
					"custom_container_config": []map[string]interface{}{
						{
							"args":              "[\\\"--debug\\\"]",
							"command":           "[\\\"python\\\",\\\"app.py\\\"]",
							"image":             "terraform-e2e-test-registry-vpc.cn-beijing.cr.aliyuncs.com/terraform-e2e/terraform-e2e:v1",
							"acceleration_type": "Default",
							"web_server_mode":   "true",
							"instance_id":       "cri-35n95x367lp1970t",
						},
					},
					"gpu_memory_size": "4096",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":          name + "_update",
						"memory_size":            "8192",
						"runtime":                "custom-container",
						"description":            "terraform测试case",
						"service_name":           "terraform-e2e-service-base",
						"initializer":            "index.initializer",
						"initialization_timeout": "10",
						"timeout":                "60",
						"handler":                "index.handler",
						"instance_type":          "fc.gpu.tesla.1",
						"disk_size":              "512",
						"instance_concurrency":   "10",
						"cpu":                    "2",
						"ca_port":                "9000",
						"gpu_memory_size":        "4096",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudFcv2FunctionMap3411 = map[string]string{
	"ca_port":                CHECKSET,
	"instance_concurrency":   CHECKSET,
	"instance_type":          CHECKSET,
	"memory_size":            CHECKSET,
	"timeout":                CHECKSET,
	"create_time":            CHECKSET,
	"initialization_timeout": CHECKSET,
	"code_checksum":          CHECKSET,
}

func AlicloudFcv2FunctionBasicDependence3411(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3393  twin
func TestAccAlicloudFcv2Function_basic3393_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv2_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv2FunctionMap3393)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv2ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv2Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv2function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv2FunctionBasicDependence3393)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":          name,
					"memory_size":            "2048",
					"runtime":                "custom",
					"description":            "terraform测试case update",
					"service_name":           "terraform-e2e-service-base",
					"initializer":            "true",
					"initialization_timeout": "3",
					"timeout":                "120",
					"handler":                "index.handler",
					"instance_type":          "c1",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "true",
									"timeout": "3",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "true",
									"timeout": "3",
								},
							},
						},
					},
					"code": []map[string]interface{}{
						{
							"zip_file": "UEsDBAoAAAAAAOt8v1YAAAAAAAAAAAAAAAAFABwAY29kZS9VVAkAA6r5dmTE+XZkdXgLAAEE9gEAAAQUAAAAUEsDBBQAAAAIAOt8v1Y14c8e4QAAAH8BAAANABwAY29kZS9pbmRleC5qc1VUCQADqvl2ZKv5dmR1eAsAAQT2AQAABBQAAAB1kFFOwzAMht9zCr+lnapGPDAQ1XYK3lGamiXCjavGYQO0i3AW7sQVSDU0oU28Wf4/+7Osc0JIMgcnulNmpR4ZMNqeEMQjhBgkWArvOMMzWskzQuVFpvRgjEea2hK+5dg6Hs3ALo8Y5WlAsYHMze36/m7dehmpVhOhLaowlmKBrtfn6CRwBJugR+L999enwsPEs6T2L7iBynEUPEgDzhL11r3UsNnChwIoSWLClnhX6fNUiDtdd0v8y1cxEzWgl+6xUytzNnkbBzpZ8LXc2cD/sgtd+Qcx7Hmm4WS79l0gxx9QSwECHgMKAAAAAADrfL9WAAAAAAAAAAAAAAAABQAYAAAAAAAAABAA7UEAAAAAY29kZS9VVAUAA6r5dmR1eAsAAQT2AQAABBQAAABQSwECHgMUAAAACADrfL9WNeHPHuEAAAB/AQAADQAYAAAAAAABAAAApIE/AAAAY29kZS9pbmRleC5qc1VUBQADqvl2ZHV4CwABBPYBAAAEFAAAAFBLBQYAAAAAAgACAJ4AAABnAQAAAAA=",
						},
					},
					"environment_variables": []map[string]interface{}{
						{},
					},
					"custom_dns": []map[string]interface{}{
						{
							"name_servers": []string{
								"223.6.6.6"},
							"searches": []string{
								"newdomain.com"},
							"dns_options": []map[string]interface{}{
								{
									"name":  name,
									"value": "2",
								},
							},
						},
					},
					"disk_size":            "10240",
					"instance_concurrency": "20",
					"layers": []string{
						"d3fc5de8d120687be2bfab761518d5de#Python39#2", "d3fc5de8d120687be2bfab761518d5de#Python39#2"},
					"cpu": "2",
					"custom_health_check_config": []map[string]interface{}{
						{
							"http_get_url":          "/healthcheck1",
							"initial_delay_seconds": "5",
							"period_seconds":        "5",
							"timeout_seconds":       "2",
							"failure_threshold":     "5",
							"success_threshold":     "2",
						},
					},
					"ca_port": "8000",
					"custom_runtime_config": []map[string]interface{}{
						{
							"command": []string{
								"npm"},
							"args": []string{
								"run", "start"},
						},
					},
					"code_checksum": "7715138221701825864",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":          name,
						"memory_size":            "2048",
						"runtime":                "custom",
						"description":            "terraform测试case update",
						"service_name":           "terraform-e2e-service-base",
						"initializer":            "true",
						"initialization_timeout": "3",
						"timeout":                "120",
						"handler":                "index.handler",
						"instance_type":          "c1",
						"disk_size":              "10240",
						"instance_concurrency":   "20",
						"layers.#":               "2",
						"cpu":                    "2",
						"ca_port":                "8000",
						"code_checksum":          "7715138221701825864",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case 3270  twin
func TestAccAlicloudFcv2Function_basic3270_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv2_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv2FunctionMap3270)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv2ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv2Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv2function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv2FunctionBasicDependence3270)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":          name,
					"memory_size":            "1024",
					"runtime":                "nodejs16",
					"description":            "terraform测试case update",
					"service_name":           "terraform-e2e-service-base",
					"initializer":            "index. initializer",
					"initialization_timeout": "3",
					"timeout":                "120",
					"handler":                "index.handler",
					"instance_type":          "c1",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "index.prefreeze",
									"timeout": "10",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "index.presstop",
									"timeout": "10",
								},
							},
						},
					},
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "code-sample-cn-beijing",
							"oss_object_name": "quick-start-sample-codes/quick-start-sample-codes-nodejs/cdn-trigger-nodejs/code.zip",
						},
					},
					"environment_variables": []map[string]interface{}{
						{},
					},
					"disk_size":            "10240",
					"instance_concurrency": "20",
					"layers": []string{
						"acs:fc:cn-beijing:official:layers/Nodejs-Puppeteer17x/versions/3"},
					"cpu": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":          name,
						"memory_size":            "1024",
						"runtime":                "nodejs16",
						"description":            "terraform测试case update",
						"service_name":           "terraform-e2e-service-base",
						"initializer":            "index. initializer",
						"initialization_timeout": "3",
						"timeout":                "120",
						"handler":                "index.handler",
						"instance_type":          "c1",
						"disk_size":              "10240",
						"instance_concurrency":   "20",
						"layers.#":               "1",
						"cpu":                    "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case 3395  twin
func TestAccAlicloudFcv2Function_basic3395_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv2_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv2FunctionMap3395)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv2ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv2Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv2function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv2FunctionBasicDependence3395)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":          name,
					"memory_size":            "8192",
					"runtime":                "custom-container",
					"description":            "terraform测试case update",
					"service_name":           "terraform-e2e-service-base",
					"initializer":            "true",
					"initialization_timeout": "3",
					"timeout":                "120",
					"handler":                "index.handler",
					"instance_type":          "fc.gpu.tesla.1",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "true",
									"timeout": "3",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "true",
									"timeout": "3",
								},
							},
						},
					},
					"environment_variables": []map[string]interface{}{
						{},
					},
					"disk_size":            "10240",
					"instance_concurrency": "20",
					"cpu":                  "2",
					"custom_health_check_config": []map[string]interface{}{
						{
							"http_get_url":          "/healthcheck1",
							"initial_delay_seconds": "5",
							"period_seconds":        "5",
							"timeout_seconds":       "2",
							"failure_threshold":     "5",
							"success_threshold":     "2",
						},
					},
					"ca_port": "8000",
					"custom_container_config": []map[string]interface{}{
						{
							"args":              "[\\\"--use-local\\\"]",
							"command":           "[\\\"s\\\",\\\"deploy\\\"]",
							"image":             "registry.cn-beijing.aliyuncs.com/serverlessdevshanxie/sd-auto-nas:v3",
							"acceleration_type": "None",
							"web_server_mode":   "true",
						},
					},
					"gpu_memory_size": "8192",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":          name,
						"memory_size":            "8192",
						"runtime":                "custom-container",
						"description":            "terraform测试case update",
						"service_name":           "terraform-e2e-service-base",
						"initializer":            "true",
						"initialization_timeout": "3",
						"timeout":                "120",
						"handler":                "index.handler",
						"instance_type":          "fc.gpu.tesla.1",
						"disk_size":              "10240",
						"instance_concurrency":   "20",
						"cpu":                    "2",
						"ca_port":                "8000",
						"gpu_memory_size":        "8192",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case 3411  twin
func TestAccAlicloudFcv2Function_basic3411_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv2_function.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv2FunctionMap3411)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv2ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv2Function")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv2function%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv2FunctionBasicDependence3411)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":          name,
					"memory_size":            "8192",
					"runtime":                "custom-container",
					"description":            "terraform测试case update",
					"service_name":           "terraform-e2e-service-base",
					"initializer":            "true",
					"initialization_timeout": "3",
					"timeout":                "120",
					"handler":                "index.handler",
					"instance_type":          "fc.gpu.tesla.1",
					"instance_lifecycle_config": []map[string]interface{}{
						{
							"pre_freeze": []map[string]interface{}{
								{
									"handler": "true",
									"timeout": "3",
								},
							},
							"pre_stop": []map[string]interface{}{
								{
									"handler": "true",
									"timeout": "3",
								},
							},
						},
					},
					"environment_variables": []map[string]interface{}{
						{},
					},
					"disk_size":            "10240",
					"instance_concurrency": "20",
					"cpu":                  "2",
					"custom_health_check_config": []map[string]interface{}{
						{
							"http_get_url":          "/healthcheck1",
							"initial_delay_seconds": "5",
							"period_seconds":        "5",
							"timeout_seconds":       "2",
							"failure_threshold":     "5",
							"success_threshold":     "2",
						},
					},
					"ca_port": "8000",
					"custom_container_config": []map[string]interface{}{
						{
							"args":              "[\\\"--use-local\\\"]",
							"command":           "[\\\"s\\\",\\\"deploy\\\"]",
							"image":             "terraform-e2e-test-registry-vpc.cn-beijing.cr.aliyuncs.com/terraform-e2e/terraform-e2e:v2",
							"acceleration_type": "Default",
							"web_server_mode":   "true",
							"instance_id":       "cri-35n95x367lp1970t",
						},
					},
					"gpu_memory_size": "8192",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":          name,
						"memory_size":            "8192",
						"runtime":                "custom-container",
						"description":            "terraform测试case update",
						"service_name":           "terraform-e2e-service-base",
						"initializer":            "true",
						"initialization_timeout": "3",
						"timeout":                "120",
						"handler":                "index.handler",
						"instance_type":          "fc.gpu.tesla.1",
						"disk_size":              "10240",
						"instance_concurrency":   "20",
						"cpu":                    "2",
						"ca_port":                "8000",
						"gpu_memory_size":        "8192",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test Fcv2 Function. <<< Resource test cases, automatically generated.
