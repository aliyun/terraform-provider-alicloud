// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudApigSecrets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudApigSecretRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"gateway_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"AI", "API"}, false),
			},
			"name_like": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"ENABLE", "DISABLE", "DELETED"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"secrets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reference_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kms_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kms_instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kms_key_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kms_secret_arn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudApigSecretRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/secrets")
	request := make(map[string]*string)
	request["pageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
	request["pageNumber"] = StringPointer("1")

	if v, ok := d.GetOk("gateway_type"); ok {
		request["gatewayType"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("name_like"); ok {
		request["nameLike"] = StringPointer(v.(string))
	}

	status, statusOk := d.GetOk("status")

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

	var secretNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}

		secretNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("APIG", "2024-03-27", action, request, nil, nil)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_apig_secrets", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.data.items", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.data.items", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["secretId"])]; !ok {
					continue
				}
			}

			if secretNameRegex != nil {
				if !secretNameRegex.MatchString(fmt.Sprint(item["name"])) {
					continue
				}
			}

			if statusOk && status.(string) != "" && status.(string) != item["status"].(string) {
				continue
			}

			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}

		pageNumber, err := strconv.Atoi(*request["pageNumber"])
		if err != nil {
			return WrapError(err)
		}

		request["pageNumber"] = StringPointer(strconv.Itoa(pageNumber + 1))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":               fmt.Sprint(object["secretId"]),
			"secret_id":        fmt.Sprint(object["secretId"]),
			"gateway_type":     object["gatewayType"],
			"name":             object["name"],
			"description":      object["description"],
			"secret_source":    object["secretSource"],
			"reference_count":  object["referenceCount"],
			"status":           object["status"],
			"create_timestamp": fmt.Sprint(object["createTimestamp"]),
			"update_timestamp": fmt.Sprint(object["updateTimestamp"]),
		}

		if v, ok := object["kmsConfig"]; ok && v != nil {
			kmsConfigMaps := make([]map[string]interface{}, 0)
			kmsConfigMap := make(map[string]interface{})
			kmsConfig := v.(map[string]interface{})

			if kmsInstanceId, ok := kmsConfig["kmsInstanceId"]; ok {
				kmsConfigMap["kms_instance_id"] = kmsInstanceId
			}

			if kmsKeyId, ok := kmsConfig["kmsKeyId"]; ok {
				kmsConfigMap["kms_key_id"] = kmsKeyId
			}

			if kmsSecretArn, ok := kmsConfig["kmsSecretArn"]; ok {
				kmsConfigMap["kms_secret_arn"] = kmsSecretArn
			}

			if versionId, ok := kmsConfig["versionId"]; ok {
				kmsConfigMap["version_id"] = versionId
			}

			kmsConfigMaps = append(kmsConfigMaps, kmsConfigMap)

			mapping["kms_config"] = kmsConfigMaps
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("secrets", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
