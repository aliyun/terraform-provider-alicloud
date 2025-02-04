package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudHbrHanaBackupClients() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHbrHanaBackupClientsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vault_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"REGISTERED", "ACTIVATED", "DEACTIVATED", "INSTALLING", "INSTALL_FAILED", "NOT_INSTALLED", "UPGRADING", "UPGRADE_FAILED", "UNINSTALLING", "UNINSTALL_FAILED", "STOPPED", "UNKNOWN"}, false),
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PageSizeLarge,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hana_backup_clients": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_setting": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"use_https": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_message": {
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

func dataSourceAlicloudHbrHanaBackupClientsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeClients"
	request := make(map[string]interface{})
	setPagingRequest(d, request, PageSizeLarge)
	request["VaultId"] = d.Get("vault_id")
	request["ClientType"] = "ECS_AGENT"
	request["SourceType"] = "HANA"

	if v, ok := d.GetOk("client_id"); ok {
		request["ClientId"] = v
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request["ClusterId"] = v
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
			response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbr_hana_backup_clients", action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		resp, err := jsonpath.Get("$.Clients.Client", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Clients.Client", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v", item["VaultId"], item["ClientId"])]; !ok {
					continue
				}
			}

			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}

			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":             fmt.Sprintf("%v:%v", object["VaultId"], object["ClientId"]),
			"vault_id":       fmt.Sprint(object["VaultId"]),
			"client_id":      fmt.Sprint(object["ClientId"]),
			"client_name":    object["ClientName"],
			"client_type":    object["ClientType"],
			"client_version": object["ClientVersion"],
			"max_version":    object["MaxVersion"],
			"cluster_id":     object["ClusterId"],
			"instance_id":    object["InstanceId"],
			"instance_name":  object["InstanceName"],
			"alert_setting":  object["AlertSetting"],
			"use_https":      object["UseHttps"],
			"network_type":   object["NetworkType"],
			"status_message": object["StatusMessage"],
			"status":         object["Status"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("hana_backup_clients", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
