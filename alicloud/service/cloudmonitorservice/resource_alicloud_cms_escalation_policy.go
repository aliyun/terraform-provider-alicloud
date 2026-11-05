// Package cloudmonitorservice. This file is hand-written from the Cms 2024-03-30 CloudSpec EscalationPolicy definition.
package cloudmonitorservice

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCmsEscalationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsEscalationPolicyCreate,
		Read:   resourceAliCloudCmsEscalationPolicyRead,
		Update: resourceAliCloudCmsEscalationPolicyUpdate,
		Delete: resourceAliCloudCmsEscalationPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"escalation_stage_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"cycle_notify_interval": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"cycle_notify_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"trigger_delay": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"notify_channels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"receivers": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"enabled_sub_channels": {
										Type:             schema.TypeSet,
										Optional:         true,
										Elem:             &schema.Schema{Type: schema.TypeString},
										DiffSuppressFunc: suppressCmsEscalationEnabledSubChannelsDiff,
									},
								},
							},
						},
						"target_incident_state": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"effect_time_range": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_zone": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"day_in_week": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
									"start_time_in_minute": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"end_time_in_minute": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// suppressCmsEscalationEnabledSubChannelsDiff suppresses spurious diffs on the
// enabled_sub_channels TypeSet count key that arise when the API returns a
// nil/empty set during refresh. The SDK represents a config-declared empty set
// (e.g. []string{}) with a "0" count key, while a state that omits the field
// surfaces an empty count key; the IDRefreshName path then plans a no-op
// change between "" and "0". This only suppresses that empty/zero count-key
// drift and never masks real element additions/removals.
func suppressCmsEscalationEnabledSubChannelsDiff(k, old, new string, d *schema.ResourceData) bool {
	if strings.HasSuffix(k, ".#") {
		if (old == "" && new == "0") || (old == "0" && new == "") {
			return true
		}
	}
	return false
}

func buildEscalationPolicyBody(d *schema.ResourceData) map[string]interface{} {
	request := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		request["name"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOkExists("enable"); ok {
		request["enable"] = v
	}

	if v, ok := d.GetOk("escalation_stage_list"); ok || d.HasChange("escalation_stage_list") {
		stages := make([]interface{}, 0)
		stageList := v.([]interface{})
		for _, stageRaw := range stageList {
			stageMap := make(map[string]interface{})
			stage, ok := stageRaw.(map[string]interface{})
			if !ok || stage == nil {
				stages = append(stages, stageMap)
				continue
			}
			stageMap["index"] = stage["index"]
			stageMap["cycleNotifyInterval"] = stage["cycle_notify_interval"]
			stageMap["cycleNotifyCount"] = stage["cycle_notify_count"]
			stageMap["triggerDelay"] = stage["trigger_delay"]
			stageMap["targetIncidentState"] = stage["target_incident_state"]

			notifyChannels := make([]interface{}, 0)
			if channelsRaw, ok := stage["notify_channels"].([]interface{}); ok {
				for _, channelRaw := range channelsRaw {
					channelMap := make(map[string]interface{})
					channel, ok := channelRaw.(map[string]interface{})
					if !ok || channel == nil {
						notifyChannels = append(notifyChannels, channelMap)
						continue
					}
					channelMap["channelType"] = channel["channel_type"]
					receivers := make([]interface{}, 0)
					if r, ok := channel["receivers"].([]interface{}); ok {
						receivers = r
					}
					channelMap["receivers"] = receivers
					if enabledSubChannelsSet, ok := channel["enabled_sub_channels"].(*schema.Set); ok {
						channelMap["enabledSubChannels"] = enabledSubChannelsSet.List()
					} else {
						channelMap["enabledSubChannels"] = channel["enabled_sub_channels"]
					}
					notifyChannels = append(notifyChannels, channelMap)
				}
			}
			stageMap["notifyChannels"] = notifyChannels

			effectTimeRangeMap := make(map[string]interface{})
			effectList, ok := stage["effect_time_range"].([]interface{})
			if ok && len(effectList) > 0 {
				if effect, ok := effectList[0].(map[string]interface{}); ok && effect != nil {
					effectTimeRangeMap["timeZone"] = effect["time_zone"]
					effectTimeRangeMap["startTimeInMinute"] = effect["start_time_in_minute"]
					effectTimeRangeMap["endTimeInMinute"] = effect["end_time_in_minute"]
					dayInWeek := make([]interface{}, 0)
					if d, ok := effect["day_in_week"].([]interface{}); ok {
						dayInWeek = d
					}
					effectTimeRangeMap["dayInWeek"] = dayInWeek
				}
			}
			if len(effectTimeRangeMap) > 0 {
				stageMap["effectTimeRange"] = effectTimeRangeMap
			}

			stages = append(stages, stageMap)
		}
		request["escalationStageList"] = stages
	}

	return request
}

func resourceAliCloudCmsEscalationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "/escalationPolicies"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var body map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query["workspace"] = helper.StringPointer(d.Get("workspace").(string))
	query["regionId"] = helper.StringPointer(client.RegionId)

	request = buildEscalationPolicyBody(d)
	body = request

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cms_escalation_policy", action, errs.AlibabaCloudSdkGoERROR)
	}

	uuidVar, _ := jsonpath.Get("$.uuid", response)
	if uuidVar == nil || fmt.Sprint(uuidVar) == "" {
		return errs.WrapErrorf(fmt.Errorf("uuid not found in Create response: %v", response), errs.DefaultErrorMsg, "alicloud_cms_escalation_policy", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%v:%v", uuidVar, d.Get("workspace").(string)))

	return resourceAliCloudCmsEscalationPolicyRead(d, meta)
}

func resourceAliCloudCmsEscalationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsEscalationPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_escalation_policy DescribeCmsEscalationPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("enable", objectRaw["enable"])
	d.Set("name", objectRaw["name"])
	d.Set("update_time", objectRaw["updateTime"])
	d.Set("uuid", objectRaw["uuid"])
	d.Set("workspace", objectRaw["workspace"])

	escalationStageListMaps := make([]map[string]interface{}, 0)
	if stageListRaw, ok := objectRaw["escalationStageList"]; ok && stageListRaw != nil {
		for _, stageChildRaw := range convert.ConvertToInterfaceArray(stageListRaw) {
			stageMap := make(map[string]interface{})
			stageChild, ok := stageChildRaw.(map[string]interface{})
			if !ok || stageChild == nil {
				escalationStageListMaps = append(escalationStageListMaps, stageMap)
				continue
			}
			stageMap["index"] = stageChild["index"]
			stageMap["cycle_notify_interval"] = stageChild["cycleNotifyInterval"]
			stageMap["cycle_notify_count"] = stageChild["cycleNotifyCount"]
			stageMap["trigger_delay"] = stageChild["triggerDelay"]
			stageMap["target_incident_state"] = stageChild["targetIncidentState"]

			notifyChannelsMaps := make([]map[string]interface{}, 0)
			if channelsRaw, ok := stageChild["notifyChannels"]; ok && channelsRaw != nil {
				for _, channelChildRaw := range convert.ConvertToInterfaceArray(channelsRaw) {
					channelMap := make(map[string]interface{})
					channelChild, ok := channelChildRaw.(map[string]interface{})
					if !ok || channelChild == nil {
						notifyChannelsMaps = append(notifyChannelsMaps, channelMap)
						continue
					}
					channelMap["channel_type"] = channelChild["channelType"]

					receiversRaw := make([]interface{}, 0)
					if channelChild["receivers"] != nil {
						receiversRaw = convert.ConvertToInterfaceArray(channelChild["receivers"])
					}
					channelMap["receivers"] = receiversRaw

					enabledSubChannelsRaw := make([]interface{}, 0)
					if channelChild["enabledSubChannels"] != nil {
						enabledSubChannelsRaw = convert.ConvertToInterfaceArray(channelChild["enabledSubChannels"])
					}
					// Only populate enabled_sub_channels when the API returns a
					// non-empty set. When the API returns nil/empty, omit the
					// field so the framework preserves the config's empty
					// representation instead of writing an explicit zero-count
					// set that diverges from a config-declared empty list during
					// refresh.
					if len(enabledSubChannelsRaw) > 0 {
						channelMap["enabled_sub_channels"] = enabledSubChannelsRaw
					}
					notifyChannelsMaps = append(notifyChannelsMaps, channelMap)
				}
			}
			stageMap["notify_channels"] = notifyChannelsMaps

			effectTimeRangeMaps := make([]map[string]interface{}, 0)
			if effectRaw, ok := stageChild["effectTimeRange"]; ok && effectRaw != nil {
				effectTimeRangeMap := make(map[string]interface{})
				if effectChild, ok := effectRaw.(map[string]interface{}); ok && effectChild != nil {
					effectTimeRangeMap["time_zone"] = effectChild["timeZone"]
					effectTimeRangeMap["start_time_in_minute"] = effectChild["startTimeInMinute"]
					effectTimeRangeMap["end_time_in_minute"] = effectChild["endTimeInMinute"]
					dayInWeekRaw := make([]interface{}, 0)
					if effectChild["dayInWeek"] != nil {
						dayInWeekRaw = convert.ConvertToInterfaceArray(effectChild["dayInWeek"])
					}
					effectTimeRangeMap["day_in_week"] = dayInWeekRaw
					effectTimeRangeMaps = append(effectTimeRangeMaps, effectTimeRangeMap)
				}
			}
			stageMap["effect_time_range"] = effectTimeRangeMaps

			escalationStageListMaps = append(escalationStageListMaps, stageMap)
		}
	}
	if err := d.Set("escalation_stage_list", escalationStageListMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudCmsEscalationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length 2, got %d", d.Id(), len(parts)))
	}
	uuid := parts[0]
	workspace := parts[1]

	action := fmt.Sprintf("/escalationPolicies/%s", uuid)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var body map[string]interface{}
	var err error
	query["workspace"] = helper.StringPointer(workspace)
	query["regionId"] = helper.StringPointer(client.RegionId)

	request = buildEscalationPolicyBody(d)
	body = request

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudCmsEscalationPolicyRead(d, meta)
}

func resourceAliCloudCmsEscalationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length 2, got %d", d.Id(), len(parts)))
	}
	uuid := parts[0]
	workspace := parts[1]

	action := fmt.Sprintf("/escalationPolicies/%s", uuid)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["workspace"] = helper.StringPointer(workspace)
	query["regionId"] = helper.StringPointer(client.RegionId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
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
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
