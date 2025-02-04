package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudAmqpStaticAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAmqpStaticAccountsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"accounts": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"access_key": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instance_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"master_uid": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"password": {
							Computed:  true,
							Sensitive: true,
							Type:      schema.TypeString,
						},
						"user_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAmqpStaticAccountsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}
	action := "ListAccounts"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("amqp-open", "2019-12-12", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_amqp_static_accounts", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data[*]", response)
	}
	result, _ := resp.(map[string]interface{})
	for _, v := range result {
		item := v.([]interface{})
		for _, i := range item {
			detail := i.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(detail["cInstanceId"], ":", detail["accessKey"])]; !ok {
					continue
				}
			}
			objects = append(objects, i)
		}

	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(object["cInstanceId"], ":", object["accessKey"]),
			"access_key":  object["accessKey"],
			"create_time": object["createTimestamp"],
			"instance_id": object["cInstanceId"],
			"master_uid":  object["masterUid"],
			"password":    object["password"],
			"user_name":   object["userName"],
		}

		ids = append(ids, fmt.Sprint(object["cInstanceId"], ":", object["accessKey"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("accounts", s); err != nil {
		return WrapError(err)
	}
	return nil
}
