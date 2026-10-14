package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceAlicloudSecurityGroupSnapshots is handwritten against the ACube pre model
// contract SecurityGroupSnapshot (resourceTypeVersion 1.0.0.34.pre, category=readonly,
// resourceTypeName ACS::Ecs::SecurityGroupSnapshot). The resource is read-only and
// system produced, so it is exposed as a Terraform data source rather than a lifecycle
// resource. Its only operation is gets -> DescribeSecurityGroupSnapshots (Ecs, version
// 2016-03-14, runtimeType=pop), rootMapping $.SecurityGroupSnapshots[*]. The model is the
// source of truth for request/response field names and paths; this file does not modify it.
func dataSourceAlicloudSecurityGroupSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSecurityGroupSnapshotsRead,
		Schema: map[string]*schema.Schema{
			// Model gets attributeMappings: $.SecurityGroupSnapshotId ->
			// SecurityGroupSnapshotIds (query, repeatList, requestRequired=false).
			"security_group_snapshot_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Model gets attributeMappings: $.RegionId -> RegionId (query,
			// requestRequired=true); defaulted from the provider client region.
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Filter is the model list private property (operatePrivateProperties.list).
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Model publicProperties: SecurityGroupSnapshotId, RegionId. Only the
			// model-confirmed fields are mapped; additional Inner API response fields are
			// not mapped pending Inner API schema verification.
			"snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSecurityGroupSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeSecurityGroupSnapshots"
	var query map[string]interface{}
	request := make(map[string]interface{})
	var response map[string]interface{}
	var err error

	// Model gets attributeMappings: RegionId -> RegionId (query, required).
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("region_id"); ok && v.(string) != "" {
		request["RegionId"] = v
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("security_group_snapshot_ids"); ok {
		ids := make([]interface{}, 0)
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
			ids = append(ids, vv.(string))
		}
		if len(ids) > 0 {
			// Model maps $.SecurityGroupSnapshotId -> SecurityGroupSnapshotIds (query,
			// repeatList); sent as a JSON string list consistent with the sibling
			// SecurityGroupSnapshotPolicy Inner API binding.
			request["SecurityGroupSnapshotIds"] = convertListToJsonString(ids)
		}
	}

	// Filter is the model list private property (operatePrivateProperties.list).
	if v, ok := d.GetOk("filter"); ok {
		request["Filter"] = v
	}

	// The model gets operation has no pagination configuration; DescribeSecurityGroupSnapshots
	// returns the full matching list in a single response, so no paging loop is used.
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2016-03-14", action, query, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_security_group_snapshots", action, AlibabaCloudSdkGoERROR)
	}

	// Model rootMapping: $.SecurityGroupSnapshots[*]. POP list responses sometimes nest the
	// array under an inner key, so fall back to that shape before treating the result as empty.
	resp, jsonErr := jsonpath.Get("$.SecurityGroupSnapshots[*]", response)
	if jsonErr != nil {
		resp, jsonErr = jsonpath.Get("$.SecurityGroupSnapshots.SecurityGroupSnapshot[*]", response)
	}
	if jsonErr != nil {
		return WrapErrorf(jsonErr, FailedGetAttributeMsg, action, "$.SecurityGroupSnapshots[*]", response)
	}
	result, _ := resp.([]interface{})

	var objects []map[string]interface{}
	for _, v := range result {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		// Defensive local filter in case the server ignores SecurityGroupSnapshotIds.
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["SecurityGroupSnapshotId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		id := fmt.Sprint(object["SecurityGroupSnapshotId"])
		mapping := map[string]interface{}{
			"id":                         id,
			"security_group_snapshot_id": id,
			"region_id":                  object["RegionId"],
		}
		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("snapshots", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
