package oos

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceSecretParameters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudOosSecretParametersRead,
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
			"secret_parameter_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tags.TagsSchemaForceNew(),
			"sort_field": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Name", "CreatedDate"}, false),
			},
			"sort_order": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Ascending", "Descending"}, false),
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"with_decryption": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_parameter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"share_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_parameter_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"constraints": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudOosSecretParametersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListSecretParameters"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["MaxResults"] = helper.PageSizeLarge

	if v, ok := d.GetOk("secret_parameter_name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("sort_field"); ok {
		request["SortField"] = v
	}

	if v, ok := d.GetOk("sort_order"); ok {
		request["SortOrder"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsJson, err := convert.ConvertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return errs.WrapError(err)
		}

		request["Tags"] = tagsJson
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

	var secretParameterNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		secretParameterNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_oos_secret_parameters", action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Parameters", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Parameters", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Name"])]; !ok {
					continue
				}
			}

			if secretParameterNameRegex != nil && !secretParameterNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}

			objects = append(objects, item)
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		parameter_versionVal, err := convert.FormatInt(object["ParameterVersion"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"id":                    fmt.Sprint(object["Name"]),
			"secret_parameter_id":   object["Id"],
			"type":                  object["Type"],
			"parameter_version":     parameter_versionVal,
			"share_type":            object["ShareType"],
			"key_id":                object["KeyId"],
			"resource_group_id":     object["ResourceGroupId"],
			"secret_parameter_name": object["Name"],
			"description":           object["Description"],
			"tags":                  tags.TagsToMap(object["Tags"]),
			"created_by":            object["CreatedBy"],
			"create_time":           object["CreatedDate"],
			"updated_by":            object["UpdatedBy"],
			"updated_date":          object["UpdatedDate"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["Name"])
		oosService := OosService{client}

		getResp, err := oosService.DescribeOosSecretParameterForDataSource(id, d)
		if err != nil {
			return errs.WrapError(err)
		}

		mapping["constraints"] = getResp["Constraints"]
		mapping["value"] = getResp["Value"]

		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("parameters", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
