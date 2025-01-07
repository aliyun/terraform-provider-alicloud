package alicloud

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/aliyun/fc-go-sdk"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Fcv2 Function. >>> Resource test cases, automatically generated.
// Case 3393
func TestAccAliCloudFcv2Function_basic3393(t *testing.T) {
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
			testAccPreCheckWithRegions(t, true, connectivity.FCV2FunctionSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": name,
					"runtime":       "custom.debian10",
					"service_name":  "${alicloud_fc_service.default.name}",
					"handler":       "index.handler",
					"code": []map[string]interface{}{
						{
							"zip_file": "UEsDBAoAAAAAAOt8v1YAAAAAAAAAAAAAAAAFABwAY29kZS9VVAkAA6r5dmTE+XZkdXgLAAEE9gEAAAQUAAAAUEsDBBQAAAAIAOt8v1Y14c8e4QAAAH8BAAANABwAY29kZS9pbmRleC5qc1VUCQADqvl2ZKv5dmR1eAsAAQT2AQAABBQAAAB1kFFOwzAMht9zCr+lnapGPDAQ1XYK3lGamiXCjavGYQO0i3AW7sQVSDU0oU28Wf4/+7Osc0JIMgcnulNmpR4ZMNqeEMQjhBgkWArvOMMzWskzQuVFpvRgjEea2hK+5dg6Hs3ALo8Y5WlAsYHMze36/m7dehmpVhOhLaowlmKBrtfn6CRwBJugR+L999enwsPEs6T2L7iBynEUPEgDzhL11r3UsNnChwIoSWLClnhX6fNUiDtdd0v8y1cxEzWgl+6xUytzNnkbBzpZ8LXc2cD/sgtd+Qcx7Hmm4WS79l0gxx9QSwECHgMKAAAAAADrfL9WAAAAAAAAAAAAAAAABQAYAAAAAAAAABAA7UEAAAAAY29kZS9VVAUAA6r5dmR1eAsAAQT2AQAABBQAAABQSwECHgMUAAAACADrfL9WNeHPHuEAAAB/AQAADQAYAAAAAAABAAAApIE/AAAAY29kZS9pbmRleC5qc1VUBQADqvl2ZHV4CwABBPYBAAAEFAAAAFBLBQYAAAAAAgACAJ4AAABnAQAAAAA=",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": name,
						"runtime":       "custom.debian10",
						"service_name":  CHECKSET,
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
							"zip_file": "UEsDBBQAAAAIAEqkA1VQpdMOPwUAANQOAAAHABwAbWFpbi50ZlVUCQADS2vqYkxr6mJ1eAsAAQT2AQAABBQAAADlV91v2zYQf89fISh72IDqy1YS22gfMqfFVmxLtrVphzYQKIqymUiiQlJ2nMD/+46kKMmOO2BD97IpcCQd746/++RJEs5RznjpPB05Dif3DeUkS2rOVjQjXGiy46CC4oI1mfOqJTiOYA3HxAGKC6ubpgosk9tyrECeskpxRH4UR35oVrZH6rc9sptoBUYStB/r/TAmQiR3ZKOkf3p3/uOJnE7KPzbX02V9evWGnqPb32AjzSwI5kRaZrQer5t5HX9/9/7y9Ql/8zGtKjp++8vHePU251rkWFm6UNDUpWRqD8yRS4KE9CLgOT4M4eLt/fION1fzdZ3S+YcPm3eGdxcBqc6vpz9jdPH6/uHhh8n1Pf/1Qvx+dTq6qOa40Aja7YEZV94SVYvHJWtc8MlRwTAqjNszkqOmkElGBOa0lhowiJw3kjmwIZIkc9IN7M7B1QWgBZGVcNZULh1pA+sONFWoJI69QFMv6SlJL8dwL9lQZFXjBGeUW5FoOvKj04kf+mEQne5wCtgYLw33gBPiHozifeNoteDKvwpsgmnGkxRW7wRIftJZYvMs56xMasal2t6LWqpkLW2HChklGWZFG9WisKnYb3DYBiNMGadyo4WjMDS5+uIrgaG4rA+jCX39F4T/GEgcjw8h6cm7UCT+15BMwkNAOurXx3Gj0goqaillLWZBoApLSL7xu/z3KQu6dhbstaqggCISMsgYFiBrepoIcpzkTYVVxR0dW3LfphJYV6VDMXFNEekn3bt0iYENS1IUzFszXmS25xRskWBW5XRhWI9bl9wSrJ3XqVeMLd23d6W3EwIGIRkn+0Ka6Jv/A4Gtbjms0MAQFjOOytksmkTjs7N4PJlMT8ZnYThTHK2DctxWtaIBfqWAVuDTCtqc6YygK4ditnvoPrFvHPTFRsUtWXDW1AnNhoB3F32x8GnWSWp9wD+4BqKw6KvfQMB2H5qZFtIpcgZihseHO4i+6Hhuej9tla2HI24zAkLeP2pb2wzYR9lniW/vNijDRqwqAXLQqzdyCSrVck4LYllgOXgvVO6ukeqZVYZUkFKUoqDv38FOK68LtilJJT1SLWhFAswy4j/SWisvocPzTSLoo86HaDQxR1JTSVpaI9wFizQZDqesIF37LxE1EBUvHJktfRwq133Jc5LTxYJwcFz3pP3W5otR8fLl68s35kB/cskK0At35nxymRCzy1TVwNwcebOrRhqCe/MCYkELyEtgfXLh/NX3mpOcPsCjG7SPwCeavCVqT2y3EGq9o3Z4G9AvRmTn5FQMAMuz1nTlNWQY1Nn45Gw6jeM4HI3i8WhYZ6BFm1oxSXOKkcLQltxuWjmHeorxdYJ4NdhT+WswVDzfXp3wXtrgOyK90F4mppv6mZGuDuzhuBZoQ3jSjnkQXf3u6pPBLNlmqLypejlmZQ02pgVJbLpBrboVpOetiEbujebJiB0DYf/EIO1UwYiitMH44GUnaUTGWUTi8QkKUe72QkznRyfU579igIdEVZjqYHBLkSCn8bfuN0+AbenjdbbtC+a7dlb9+01Be0D3oi+6zNdvPsTvZlhqAHinIu27zQfz/hed6lm9eeSBap/9BwqvImvv/1V8K8QpgpqBIbJEC+K2Hwf9RwEwzxs49UsVWAkdGpJIs/oq4ZXiVzATQfj7aV3JdLMSIE4JvQUG33gGCx9KFcYgD2u9XqdXfx8Ea5J6OfiXgDvuvOwWzGUzNbNFeszvAYOWErLaQNYWvnIOQfn02fWDlDEJi6j+7N7s60F6nNzRUzVlSviOnik482j7J1BLAQIeAxQAAAAIAEqkA1VQpdMOPwUAANQOAAAHABgAAAAAAAEAAACAgQAAAABtYWluLnRmVVQFAANLa+pidXgLAAEE9gEAAAQUAAAAUEsFBgAAAAABAAEATQAAAIAFAAAAAA==",
						},
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
								"223.5.5.5", "221.10.10.10", "210.22.22.10"},
							"searches": []string{
								"mydomain.com", "mydomain1.com", "mydomain2.com"},
							"dns_options": []map[string]interface{}{
								{
									"name":  name + "_update",
									"value": "1",
								},
								{
									"name":  name + "_update",
									"value": "2",
								},
								{
									"name":  name + "_update",
									"value": "3",
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
					"cpu":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "512",
						"cpu":       "1",
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
						"acs:fc:cn-hangzhou:official:layers/Python310/versions/2", "acs:fc:cn-hangzhou:official:layers/Nodejs18/versions/1", "acs:fc:cn-hangzhou:official:layers/ServerlessDevs/versions/1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"layers.#": "3",
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
					"environment_variables": map[string]string{
						"env_k": "env_v",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_variables.%":     "1",
						"environment_variables.env_k": "env_v",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_variables.%":     "0",
						"environment_variables.env_k": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": map[string]string{
						"LD_LIBRARY_PATH": "/code:/code/lib:/usr/local/lib:/opt/lib:/opt/php8.1/lib:/opt/php8.0/lib:/opt/php7.2/lib",
						"NODE_PATH":       "/opt/nodejs/node_modules",
						"PATH":            "/opt/nodejs16/bin:/usr/local/bin/apache-maven/bin:/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/ruby/bin:/opt/bin:/code:/code/bin",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_variables.%":               "3",
						"environment_variables.LD_LIBRARY_PATH": "/code:/code/lib:/usr/local/lib:/opt/lib:/opt/php8.1/lib:/opt/php8.0/lib:/opt/php7.2/lib",
						"environment_variables.NODE_PATH":       "/opt/nodejs/node_modules",
						"environment_variables.PATH":            "/opt/nodejs16/bin:/usr/local/bin/apache-maven/bin:/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/ruby/bin:/opt/bin:/code:/code/bin",
					}),
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
								"command1", "command2", "command3"},
							"args": []string{
								"arg1", "arg2", "arg3"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"code_checksum": "14474122665904472553",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"code_checksum": "14474122665904472553",
					}),
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
					"handler": "index.handler1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"handler": "index.handler1",
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
					"code_checksum": "7715138221701825864",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"code_checksum": "7715138221701825864",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_dns": []map[string]interface{}{
						{
							"name_servers": []string{
								"223.6.6.6", "221.10.10.1", "210.10.10.1"},
							"searches": []string{
								"newdomain.com", "newdomain2.com", "newdomain3.com"},
							"dns_options": []map[string]interface{}{
								{
									"name":  name + "_update",
									"value": "2",
								},
								{
									"name":  name + "_update",
									"value": "3",
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
						"acs:fc:cn-hangzhou:official:layers/Python310-Package-Collection/versions/2"},
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
					"function_name":          name + "_update",
					"memory_size":            "1024",
					"runtime":                "custom.debian10",
					"description":            "terraform测试case",
					"service_name":           "${alicloud_fc_service.default.name}",
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
							"zip_file": "UEsDBAoAAAAAAEZU5FYAAAAAAAAAAAAAAAAFABwAY29kZS9VVAkAAySFo2Q8haNkdXgLAAEE9gEAAAQUAAAAUEsDBBQAAAAIAEZU5FaBE2T3pAEAANQCAAANABwAY29kZS9pbmRleC5weVVUCQADJIWjZCWFo2R1eAsAAQT2AQAABBQAAABlUsFq20AQvesrpuSwUrElN1AwhpySgqEGl2DowTViLY2sJasddXcUJ/2UfEv+qb/QkS0b0expd96b92YeewPTz1MoqDTusICOq+m8r0SRaVryDJYOB4Gi5bfVap3/XD+uHuAO9mqJ1hIcydvy0y+nougGNgTo9N4icI1gnGGjrfmDHirU3HmEuGZuwyLLarRtKuBr59KCmqykomvQcV4ia2OzL1/nt7N5WnNjE1FuLeogko1cetpHg84VbMiBDrBHS8e/72/SV2I1psUFOcYXThaCyel3k+a7y5LpAXl1qsXJmJEaV1GsrkpCVUkURb18rV1ppQHds/HkJhBYe849hpZcQLHqdQZjsRp4W1UV6VBVuxPH4+8OA+edN//xRsjArcjD0wSeZbsLMzWMTYgHQwBTwVN6GiYcDdexWm42P3J1xQEkV08FhgBFF5iaywRQoy7Rhyux1eH8kEQJAjXItWQgPI+nuthwF2RqdTubwfq7GhY6Z5APeoJvY3XfL+14yq8tqgmoPoGstdo4lewuaqME47P45INcMpjIz3KwHf3P3T9QSwECHgMKAAAAAABGVORWAAAAAAAAAAAAAAAABQAYAAAAAAAAABAA7UEAAAAAY29kZS9VVAUAAySFo2R1eAsAAQT2AQAABBQAAABQSwECHgMUAAAACABGVORWgRNk96QBAADUAgAADQAYAAAAAAABAAAApIE/AAAAY29kZS9pbmRleC5weVVUBQADJIWjZHV4CwABBPYBAAAEFAAAAFBLBQYAAAAAAgACAJ4AAAAqAgAAAAA=",
						},
					},
					"code_checksum": "1750949844529959033",
					"custom_dns": []map[string]interface{}{
						{
							"name_servers": []string{
								"223.5.5.5", "221.10.10.10", "210.22.22.10"},
							"searches": []string{
								"mydomain.com", "mydomain1.com", "mydomain2.com"},
							"dns_options": []map[string]interface{}{
								{
									"name":  name + "_update",
									"value": "1",
								},
								{
									"name":  name + "_update",
									"value": "2",
								},
								{
									"name":  name + "_update",
									"value": "3",
								},
							},
						},
					},
					"disk_size":            "512",
					"instance_concurrency": "10",
					"layers": []string{
						"acs:fc:cn-hangzhou:official:layers/Python310/versions/2", "acs:fc:cn-hangzhou:official:layers/Nodejs18/versions/1", "acs:fc:cn-hangzhou:official:layers/ServerlessDevs/versions/1"},
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
								"command1", "command2", "command3"},
							"args": []string{
								"arg1", "arg2", "arg3"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":          name + "_update",
						"memory_size":            "1024",
						"runtime":                "custom.debian10",
						"description":            "terraform测试case",
						"service_name":           CHECKSET,
						"initializer":            "index.initializer",
						"initialization_timeout": "10",
						"timeout":                "60",
						"handler":                "index.handler",
						"instance_type":          "e1",
						"disk_size":              "512",
						"instance_concurrency":   "10",
						"layers.#":               "3",
						"cpu":                    "1",
						"ca_port":                "9000",
						"code_checksum":          "1750949844529959033",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code"},
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
	"function_arn":           CHECKSET,
}

func AlicloudFcv2FunctionBasicDependence3393(name string) string {
	dir, err := ioutil.TempDir(os.TempDir(), name)
	if err != nil {
		log.Printf("Failed to create temp directory: %s. Error: %v", dir, err)
		return ""
	}
	result := `'use strict';
	           exports.initializer = (context, callback) => {
		           console.log('hello init');
		           callback(null, 'hello init');
			   }
	           exports.handler = (event, context, callback) => {
		           console.log(event.toString())
		           callback(null, 'hello world');
			   }`
	filePath := filepath.Join(dir, "hello.js")
	err = ioutil.WriteFile(filePath, []byte(result), 0644)
	if err != nil {
		log.Printf("Failed to write file: %s. Error: %v", filePath, err)
		return ""
	}
	// Create the zip file.
	zipped := &bytes.Buffer{}
	err = fc.ZipDir(dir, zipped)
	if err != nil {
		return ""
	}
	zipFilePath := filepath.Join(os.TempDir(), name+".zip")
	err = ioutil.WriteFile(zipFilePath, zipped.Bytes(), 0644)
	if err != nil {
		log.Printf("Failed to write zip file: %s. Error: %v", zipFilePath, err)
		return ""
	}

	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

locals {
  container_command = "[\"python\", \"server.py\"]"
  container_args    = "[\"a1\", \"a2\"]"
}

output "container_command" {
  value = local.container_command
}

output "container_args" {
  value = local.container_args
}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = "tf unit test"
}

resource "alicloud_log_store" "default" {
  project          = alicloud_log_project.default.name
  name             = var.name
  retention_period = "3000"
  shard_count      = 1
}
resource "alicloud_fc_service" "default" {
  name        = var.name
  description = "tf unit test"
  log_config {
    project  = alicloud_log_project.default.name
    logstore = alicloud_log_store.default.name
  }
  role       = alicloud_ram_role.default.arn
  depends_on = ["alicloud_ram_role_policy_attachment.default"]
}
resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket = alicloud_oss_bucket.default.id
  key    = "fc/hello.zip"
  source = "%s"
}

resource "alicloud_ram_role" "default" {
  name        = var.name
  document    = <<EOF
  %s
  EOF
  description = "this is a test"
  force       = true
}
resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}
resource "alicloud_ram_role_policy_attachment" "acr" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunContainerRegistryReadOnlyAccess"
  policy_type = "System"
}

`, name, zipFilePath, testFCRoleTemplate)
}

// Case 3270
func TestAccAliCloudFcv2Function_basic3270(t *testing.T) {
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
			testAccPreCheckWithRegions(t, true, connectivity.FCV2FunctionSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": name,
					"runtime":       "nodejs14",
					"service_name":  "${alicloud_fc_service.default.name}",
					"handler":       "index.handler",
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "code-sample-cn-hangzhou",
							"oss_object_name": "quick-start-sample-codes/quick-start-sample-codes-nodejs/RocketMQ-producer-nodejs14-event/code.zip",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": name,
						"runtime":       "nodejs14",
						"service_name":  CHECKSET,
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
							"oss_bucket_name": "function-template-code-cn-hangzhou",
							"oss_object_name": "alimebot-nodejs-code.zip",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": map[string]string{
						"env_k": "env_v",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_variables.%":     "1",
						"environment_variables.env_k": "env_v",
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
					"disk_size": "512",
					"cpu":       "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "512",
						"cpu":       "2",
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
						"acs:fc:cn-hangzhou:official:layers/Nodejs-Puppeteer17x/versions/3"},
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "code-sample-cn-hangzhou",
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
					"service_name":           "${alicloud_fc_service.default.name}",
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
							"oss_bucket_name": "code-sample-cn-hangzhou",
							"oss_object_name": "quick-start-sample-codes/quick-start-sample-codes-nodejs/RocketMQ-producer-nodejs14-event/code.zip",
						},
					},
					"disk_size":            "512",
					"instance_concurrency": "10",
					"layers": []string{
						"acs:fc:cn-hangzhou:official:layers/Nodejs-Puppeteer17x/versions/3"},
					"cpu": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":          name + "_update",
						"memory_size":            "1024",
						"runtime":                "nodejs14",
						"description":            "terraform测试case",
						"service_name":           CHECKSET,
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
				ImportStateVerifyIgnore: []string{"code"},
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
	dir, err := ioutil.TempDir(os.TempDir(), name)
	if err != nil {
		log.Printf("Failed to create temp directory: %s. Error: %v", dir, err)
		return ""
	}
	result := `'use strict';
	           exports.initializer = (context, callback) => {
		           console.log('hello init');
		           callback(null, 'hello init');
			   }
	           exports.handler = (event, context, callback) => {
		           console.log(event.toString())
		           callback(null, 'hello world');
			   }`
	filePath := filepath.Join(dir, "hello.js")
	err = ioutil.WriteFile(filePath, []byte(result), 0644)
	if err != nil {
		log.Printf("Failed to write file: %s. Error: %v", filePath, err)
		return ""
	}
	// Create the zip file.
	zipped := &bytes.Buffer{}
	err = fc.ZipDir(dir, zipped)
	if err != nil {
		return ""
	}
	zipFilePath := filepath.Join(os.TempDir(), name+".zip")
	err = ioutil.WriteFile(zipFilePath, zipped.Bytes(), 0644)
	if err != nil {
		log.Printf("Failed to write zip file: %s. Error: %v", zipFilePath, err)
		return ""
	}

	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

locals {
  container_command = "[\"python\", \"server.py\"]"
  container_args    = "[\"a1\", \"a2\"]"
}

output "container_command" {
  value = local.container_command
}

output "container_args" {
  value = local.container_args
}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = "tf unit test"
}

resource "alicloud_log_store" "default" {
  project          = alicloud_log_project.default.name
  name             = var.name
  retention_period = "3000"
  shard_count      = 1
}
resource "alicloud_fc_service" "default" {
  name        = var.name
  description = "tf unit test"
  log_config {
    project  = alicloud_log_project.default.name
    logstore = alicloud_log_store.default.name
  }
  role       = alicloud_ram_role.default.arn
  depends_on = ["alicloud_ram_role_policy_attachment.default"]
}
resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket = alicloud_oss_bucket.default.id
  key    = "fc/hello.zip"
  source = "%s"
}

resource "alicloud_ram_role" "default" {
  name        = var.name
  document    = <<EOF
  %s
  EOF
  description = "this is a test"
  force       = true
}
resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}
resource "alicloud_ram_role_policy_attachment" "acr" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunContainerRegistryReadOnlyAccess"
  policy_type = "System"
}

`, name, zipFilePath, testFCRoleTemplate)
}

// Case 3395
func TestAccAliCloudFcv2Function_basic3395(t *testing.T) {
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
			testAccPreCheckWithRegions(t, true, connectivity.FCV2FunctionSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": name,
					"runtime":       "custom-container",
					"service_name":  "${alicloud_fc_service.default.name}",
					"handler":       "index.handler",
					"custom_container_config": []map[string]interface{}{
						{
							"args":              "[\\\"--debug\\\"]",
							"command":           "[\\\"python\\\",\\\"app.py\\\"]",
							"image":             "registry-vpc.cn-hangzhou.aliyuncs.com/fc-stable-diffusion/stable-diffusion:v1",
							"acceleration_type": "Default",
							"web_server_mode":   "true",
						},
					},
					"instance_type":   "fc.gpu.tesla.1",
					"gpu_memory_size": "2048",
					"memory_size":     "1280",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":   name,
						"runtime":         "custom-container",
						"service_name":    CHECKSET,
						"handler":         "index.handler",
						"instance_type":   "fc.gpu.tesla.1",
						"gpu_memory_size": "2048",
						"memory_size":     "1280",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gpu_memory_size": "4096",
					"memory_size":     "2560",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gpu_memory_size": "4096",
						"memory_size":     "2560",
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
					"environment_variables": map[string]string{
						"env_k": "env_v",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_variables.%":     "1",
						"environment_variables.env_k": "env_v",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_variables.%":     "0",
						"environment_variables.env_k": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_variables": map[string]string{
						"LD_LIBRARY_PATH": "/code:/code/lib:/usr/local/lib:/opt/lib:/opt/php8.1/lib:/opt/php8.0/lib:/opt/php7.2/lib",
						"NODE_PATH":       "/opt/nodejs/node_modules",
						"PATH":            "/opt/nodejs16/bin:/usr/local/bin/apache-maven/bin:/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/ruby/bin:/opt/bin:/code:/code/bin",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_variables.%":               "3",
						"environment_variables.LD_LIBRARY_PATH": "/code:/code/lib:/usr/local/lib:/opt/lib:/opt/php8.1/lib:/opt/php8.0/lib:/opt/php7.2/lib",
						"environment_variables.NODE_PATH":       "/opt/nodejs/node_modules",
						"environment_variables.PATH":            "/opt/nodejs16/bin:/usr/local/bin/apache-maven/bin:/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/ruby/bin:/opt/bin:/code:/code/bin",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size": "512",
					"cpu":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "512",
						"cpu":       "1",
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
							"image":             "registry-vpc.cn-hangzhou.aliyuncs.com/fc-stable-diffusion/stable-diffusion:v1",
							"acceleration_type": "None",
							"web_server_mode":   "false",
						},
					},
					"instance_concurrency": "1",
					"ca_port":              "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_concurrency": "1",
						"ca_port":              "0",
					}),
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
							"image":             "registry.cn-hangzhou.aliyuncs.com/serverlessdevshanxie/sd-auto-nas:v3",
							"acceleration_type": "None",
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
					"service_name":           "${alicloud_fc_service.default.name}",
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
							"image":             "registry-vpc.cn-hangzhou.aliyuncs.com/fc-stable-diffusion/stable-diffusion:v1",
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
						"service_name":           CHECKSET,
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
}

func AlicloudFcv2FunctionBasicDependence3395(name string) string {
	dir, err := ioutil.TempDir(os.TempDir(), name)
	if err != nil {
		log.Printf("Failed to create temp directory: %s. Error: %v", dir, err)
		return ""
	}
	result := `'use strict';
	           exports.initializer = (context, callback) => {
		           console.log('hello init');
		           callback(null, 'hello init');
			   }
	           exports.handler = (event, context, callback) => {
		           console.log(event.toString())
		           callback(null, 'hello world');
			   }`
	filePath := filepath.Join(dir, "hello.js")
	err = ioutil.WriteFile(filePath, []byte(result), 0644)
	if err != nil {
		log.Printf("Failed to write file: %s. Error: %v", filePath, err)
		return ""
	}
	// Create the zip file.
	zipped := &bytes.Buffer{}
	err = fc.ZipDir(dir, zipped)
	if err != nil {
		return ""
	}
	zipFilePath := filepath.Join(os.TempDir(), name+".zip")
	err = ioutil.WriteFile(zipFilePath, zipped.Bytes(), 0644)
	if err != nil {
		log.Printf("Failed to write zip file: %s. Error: %v", zipFilePath, err)
		return ""
	}

	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

locals {
  container_command = "[\"python\", \"server.py\"]"
  container_args    = "[\"a1\", \"a2\"]"
}

output "container_command" {
  value = local.container_command
}

output "container_args" {
  value = local.container_args
}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = "tf unit test"
}

resource "alicloud_log_store" "default" {
  project          = alicloud_log_project.default.name
  name             = var.name
  retention_period = "3000"
  shard_count      = 1
}
resource "alicloud_fc_service" "default" {
  name        = var.name
  description = "tf unit test"
  log_config {
    project  = alicloud_log_project.default.name
    logstore = alicloud_log_store.default.name
  }
  role       = alicloud_ram_role.default.arn
  depends_on = ["alicloud_ram_role_policy_attachment.default"]
}
resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket = alicloud_oss_bucket.default.id
  key    = "fc/hello.zip"
  source = "%s"
}

resource "alicloud_ram_role" "default" {
  name        = var.name
  document    = <<EOF
  %s
  EOF
  description = "this is a test"
  force       = true
}
resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}
resource "alicloud_ram_role_policy_attachment" "acr" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunContainerRegistryReadOnlyAccess"
  policy_type = "System"
}

`, name, zipFilePath, testFCRoleTemplate)
}

// Case 3393  twin
func TestAccAliCloudFcv2Function_basic3393_twin(t *testing.T) {
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
			testAccPreCheckWithRegions(t, true, connectivity.FCV2FunctionSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":          name,
					"memory_size":            "2048",
					"runtime":                "custom.debian10",
					"description":            "terraform测试case update",
					"service_name":           "${alicloud_fc_service.default.name}",
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
					"environment_variables": map[string]string{
						"env_k": "env_v",
					},
					"custom_dns": []map[string]interface{}{
						{
							"name_servers": []string{
								"223.6.6.6", "221.10.10.1", "210.10.10.1"},
							"searches": []string{
								"newdomain.com", "newdomain2.com", "newdomain3.com"},
							"dns_options": []map[string]interface{}{
								{
									"name":  name,
									"value": "2",
								},
								{
									"name":  name,
									"value": "3",
								},
								{
									"name":  name,
									"value": "3",
								},
							},
						},
					},
					"disk_size":            "10240",
					"instance_concurrency": "20",
					"layers": []string{
						"acs:fc:cn-hangzhou:official:layers/ServerlessDevs/versions/1", "acs:fc:cn-hangzhou:official:layers/Nodejs18/versions/1"},
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
								"command1", "command2", "command3"},
							"args": []string{
								"arg1", "arg2", "arg3"},
						},
					},
					"code_checksum": "7715138221701825864",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":               name,
						"memory_size":                 "2048",
						"runtime":                     "custom.debian10",
						"description":                 "terraform测试case update",
						"service_name":                CHECKSET,
						"initializer":                 "true",
						"initialization_timeout":      "3",
						"timeout":                     "120",
						"handler":                     "index.handler",
						"instance_type":               "c1",
						"disk_size":                   "10240",
						"instance_concurrency":        "20",
						"layers.#":                    "2",
						"cpu":                         "2",
						"ca_port":                     "8000",
						"code_checksum":               "7715138221701825864",
						"environment_variables.%":     "1",
						"environment_variables.env_k": "env_v",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code"},
			},
		},
	})
}

// Case 3270  twin
func TestAccAliCloudFcv2Function_basic3270_twin(t *testing.T) {
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
			testAccPreCheckWithRegions(t, true, connectivity.FCV2FunctionSupportRegions)
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
					"service_name":           "${alicloud_fc_service.default.name}",
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
							"oss_bucket_name": "code-sample-cn-hangzhou",
							"oss_object_name": "quick-start-sample-codes/quick-start-sample-codes-nodejs/cdn-trigger-nodejs/code.zip",
						},
					},
					"environment_variables": map[string]string{
						"env_k": "env_v",
					},
					"disk_size":            "10240",
					"instance_concurrency": "20",
					"layers": []string{
						"acs:fc:cn-hangzhou:official:layers/Nodejs-Puppeteer17x/versions/3"},
					"cpu": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":               name,
						"memory_size":                 "1024",
						"runtime":                     "nodejs16",
						"description":                 "terraform测试case update",
						"service_name":                CHECKSET,
						"initializer":                 "index. initializer",
						"initialization_timeout":      "3",
						"timeout":                     "120",
						"handler":                     "index.handler",
						"instance_type":               "c1",
						"disk_size":                   "10240",
						"instance_concurrency":        "20",
						"layers.#":                    "1",
						"cpu":                         "1",
						"environment_variables.%":     "1",
						"environment_variables.env_k": "env_v",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code"},
			},
		},
	})
}

// Case 3395  twin
func TestAccAliCloudFcv2Function_basic3395_twin(t *testing.T) {
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
			testAccPreCheckWithRegions(t, true, connectivity.FCV2FunctionSupportRegions)
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
					"service_name":           "${alicloud_fc_service.default.name}",
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
					"environment_variables": map[string]string{
						"env_k": "env_v",
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
							"image":             "registry.cn-hangzhou.aliyuncs.com/serverlessdevshanxie/sd-auto-nas:v3",
							"acceleration_type": "None",
							"web_server_mode":   "true",
						},
					},
					"gpu_memory_size": "8192",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":               name,
						"memory_size":                 "8192",
						"runtime":                     "custom-container",
						"description":                 "terraform测试case update",
						"service_name":                CHECKSET,
						"initializer":                 "true",
						"initialization_timeout":      "3",
						"timeout":                     "120",
						"handler":                     "index.handler",
						"instance_type":               "fc.gpu.tesla.1",
						"disk_size":                   "10240",
						"instance_concurrency":        "20",
						"cpu":                         "2",
						"ca_port":                     "8000",
						"gpu_memory_size":             "8192",
						"environment_variables.%":     "1",
						"environment_variables.env_k": "env_v",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code"},
			},
		},
	})
}

// Test Fcv2 Function. <<< Resource test cases, automatically generated.
