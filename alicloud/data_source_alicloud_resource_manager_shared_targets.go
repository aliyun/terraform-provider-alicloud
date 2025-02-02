package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudResourceManagerSharedTargets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudResourceManagerSharedTargetsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_share_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Associated", "Associating", "Disassociated", "Disassociating", "Failed"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"targets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_share_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudResourceManagerSharedTargetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListResourceShareAssociations"
	request := make(map[string]interface{})

	request["MaxResults"] = PageSizeMedium
	request["AssociationType"] = "Target"

	if v, ok := d.GetOk("resource_share_id"); ok {
		request["ResourceShareIds"] = []string{v.(string)}
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

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_resource_manager_shared_targets", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.ResourceShareAssociations", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ResourceShareAssociations", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["EntityId"])]; !ok {
					continue
				}
			}

			if statusOk && status.(string) != "" && status.(string) != item["AssociationStatus"].(string) {
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
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["EntityId"]),
			"target_id":         fmt.Sprint(object["EntityId"]),
			"resource_share_id": object["ResourceShareId"],
			"status":            object["AssociationStatus"],
		}

		ids = append(ids, fmt.Sprint(object["EntityId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("targets", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
