package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDataWorksNodeOnBaseline() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksNodeOnBaselineRead,
		Schema: map[string]*schema.Schema{
			"baseline_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudDataWorksNodeOnBaselineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "GetNodeOnBaseline"
	request := map[string]interface{}{
		"BaselineId": d.Get("baseline_id").(string),
	}

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_node_on_baseline", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
	}
	result, ok := resp.([]interface{})
	if !ok || len(result) == 0 {
		d.SetId("")
		return WrapErrorf(NotFoundErr("DataWorks NodeOnBaseline", d.Get("baseline_id").(string)), NotFoundWithResponse, response)
	}

	var (
		nodeIDFilter    = d.Get("node_id").(string)
		ownerFilter     = d.Get("owner").(string)
		projectIDFilter = d.Get("project_id").(string)
		nameRegexStr    = d.Get("name_regex").(string)
		nameRegex       *regexp.Regexp
	)
	if nameRegexStr != "" {
		if r, e := regexp.Compile(nameRegexStr); e == nil {
			nameRegex = r
		}
	}

	var matched map[string]interface{}
	for _, v := range result {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if nodeIDFilter != "" && fmt.Sprint(item["NodeId"]) != nodeIDFilter {
			continue
		}
		if ownerFilter != "" && fmt.Sprint(item["Owner"]) != ownerFilter {
			continue
		}
		if projectIDFilter != "" && fmt.Sprint(item["ProjectId"]) != projectIDFilter {
			continue
		}
		if nameRegex != nil {
			nodeName := fmt.Sprint(item["NodeName"])
			if !nameRegex.MatchString(nodeName) {
				continue
			}
		}
		matched = item
		break
	}

	if matched == nil {
		d.SetId("")
		return WrapErrorf(NotFoundErr("DataWorks NodeOnBaseline", d.Get("baseline_id").(string)), NotFoundWithResponse, response)
	}

	nodeID := fmt.Sprint(matched["NodeId"])
	d.SetId(fmt.Sprintf("%s:%s", d.Get("baseline_id").(string), nodeID))

	if err := d.Set("node_id", nodeID); err != nil {
		return WrapError(err)
	}
	if err := d.Set("node_name", matched["NodeName"]); err != nil {
		return WrapError(err)
	}
	if err := d.Set("owner", matched["Owner"]); err != nil {
		return WrapError(err)
	}
	if err := d.Set("project_id", matched["ProjectId"]); err != nil {
		return WrapError(err)
	}
	if err := d.Set("region_id", client.RegionId); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		s := map[string]interface{}{
			"id":          d.Id(),
			"baseline_id": d.Get("baseline_id").(string),
			"node_id":     nodeID,
			"node_name":   matched["NodeName"],
			"owner":       matched["Owner"],
			"project_id":  matched["ProjectId"],
			"region_id":   client.RegionId,
		}
		writeToFile(output.(string), s)
	}

	return nil
}
