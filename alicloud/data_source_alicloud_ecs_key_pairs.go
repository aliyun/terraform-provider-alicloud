package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudEcsKeyPairs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEcsKeyPairsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"finger_print": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchemaForceNew(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"pairs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_pair_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finger_print": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Resource{Schema: outputInstancesSchema()},
						},
					},
				},
			},
			"key_pairs": {
				Type:       schema.TypeList,
				Computed:   true,
				Deprecated: "Field `key_pairs` has been deprecated from provider version 1.121.0. New field `pairs` instead.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_pair_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finger_print": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Resource{Schema: outputInstancesSchema()},
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEcsKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeKeyPairs"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	if v, ok := d.GetOk("finger_print"); ok {
		request["KeyPairFingerPrint"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	var objects []map[string]interface{}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var keyPairNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		keyPairNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_key_pairs", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.KeyPairs.KeyPair", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.KeyPairs.KeyPair", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["KeyPairName"])]; !ok {
					continue
				}
			}

			if keyPairNameRegex != nil {
				if !keyPairNameRegex.MatchString(fmt.Sprint(item["KeyPairName"])) {
					continue
				}
			}

			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["KeyPairName"]),
			"key_pair_name":     fmt.Sprint(object["KeyPairName"]),
			"key_name":          fmt.Sprint(object["KeyPairName"]),
			"finger_print":      object["KeyPairFingerPrint"],
			"resource_group_id": object["ResourceGroupId"],
			"tags":              tagsToMap(object["Tags"].(map[string]interface{})["Tag"]),
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["KeyPairName"])

		instancesMaps, err := getInstances(client, fmt.Sprint(mapping["id"]))
		if err != nil {
			return WrapError(err)
		}

		mapping["instances"] = instancesMaps

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("pairs", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("key_pairs", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

func getInstances(client *connectivity.AliyunClient, id string) (instances []map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeInstances"

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"KeyPairName": id,
		"MaxResults":  PageSizeLarge,
	}

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return instances, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Instances.Instance", response)
		if err != nil {
			return instances, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instances.Instance", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return instances, nil
		}

		for _, v := range resp.([]interface{}) {
			instanceArg := v.(map[string]interface{})
			instanceMap := map[string]interface{}{}

			if fmt.Sprint(instanceArg["KeyPairName"]) == id {

				if instanceId, ok := instanceArg["InstanceId"]; ok {
					instanceMap["instance_id"] = instanceId
				}

				if instanceName, ok := instanceArg["InstanceName"]; ok {
					instanceMap["instance_name"] = instanceName
				}

				if description, ok := instanceArg["Description"]; ok {
					instanceMap["description"] = description
				}

				if imageId, ok := instanceArg["ImageId"]; ok {
					instanceMap["image_id"] = imageId
				}

				if regionId, ok := instanceArg["RegionId"]; ok {
					instanceMap["region_id"] = regionId
				}

				if zoneId, ok := instanceArg["ZoneId"]; ok {
					instanceMap["availability_zone"] = zoneId
				}

				if instanceType, ok := instanceArg["InstanceType"]; ok {
					instanceMap["instance_type"] = instanceType
				}

				if vpcAttributes, ok := instanceArg["VpcAttributes"]; ok {
					if vswitchId, ok := vpcAttributes.(map[string]interface{})["VSwitchId"]; ok {
						instanceMap["vswitch_id"] = vswitchId
					}
				}

				publicIpAddress := instanceArg["PublicIpAddress"].(map[string]interface{})["IpAddress"].([]interface{})
				EipAddress := instanceArg["EipAddress"].(map[string]interface{})["IpAddress"]

				if len(publicIpAddress) > 0 {
					instanceMap["public_ip"] = publicIpAddress[0]
				} else {
					instanceMap["public_ip"] = EipAddress
				}

				innerIpAddress := instanceArg["InnerIpAddress"].(map[string]interface{})["IpAddress"].([]interface{})
				vpcPrivateIpAddress := instanceArg["VpcAttributes"].(map[string]interface{})["PrivateIpAddress"].(map[string]interface{})["IpAddress"].([]interface{})

				if len(innerIpAddress) > 0 {
					instanceMap["private_ip"] = innerIpAddress[0]
				} else if len(vpcPrivateIpAddress) > 0 {
					instanceMap["private_ip"] = vpcPrivateIpAddress[0]
				}

				if keyPairName, ok := instanceArg["KeyPairName"]; ok {
					instanceMap["key_name"] = keyPairName
				}

				if status, ok := instanceArg["Status"]; ok {
					instanceMap["status"] = status
				}

				instances = append(instances, instanceMap)
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	return instances, nil
}
