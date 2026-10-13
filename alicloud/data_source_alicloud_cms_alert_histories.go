// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsAlertHistories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsAlertHistoriesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"alert_history_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alert_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"biz_source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"latest_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name_keyword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"label_filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"labels": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"opt": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"histories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_history_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"biz_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"alert_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"latest_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"instance_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"send": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"actions": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"notification": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"contacts": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"custom_webhooks": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"ding_webhooks": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"fs_webhooks": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"groups": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"silence_time": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"slack_webhooks": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"wx_webhooks": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
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

func dataSourceAliCloudCmsAlertHistoriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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

	action := "/alertHistories"
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	query["RegionId"] = StringPointer(client.RegionId)
	query["maxResults"] = StringPointer(strconv.Itoa(PageSizeLarge))
	if v, ok := d.GetOk("alert_history_id"); ok {
		query["alertHistoryId"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("alert_rule_id"); ok {
		query["alertRuleId"] = StringPointer(v.(string))
	}
	if v := d.Get("biz_source"); !IsNil(v) {
		query["bizSource"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}
	if v := d.Get("max_level"); !IsNil(v) {
		query["maxLevel"] = StringPointer(v.(string))
	}
	if v := d.Get("latest_level"); !IsNil(v) {
		query["latestLevel"] = StringPointer(v.(string))
	}
	if v := d.Get("instance_key"); !IsNil(v) {
		query["instanceKey"] = StringPointer(v.(string))
	}
	if v := d.Get("display_name_keyword"); !IsNil(v) {
		query["displayNameKeyWord"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("label_filter"); ok {
		labelFilterList := v.([]interface{})
		if len(labelFilterList) > 0 && labelFilterList[0] != nil {
			labelFilterMap := labelFilterList[0].(map[string]interface{})
			if labels, ok := labelFilterMap["labels"].(map[string]interface{}); ok && len(labels) > 0 {
				labelsJson, err := json.Marshal(labels)
				if err != nil {
					return WrapError(err)
				}
				query["labelFilter.labels"] = StringPointer(string(labelsJson))
			}
			if opt, ok := labelFilterMap["opt"].(string); ok && opt != "" {
				query["labelFilter.opt"] = StringPointer(opt)
			}
		}
	}

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, query)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.data[*]", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["alertHistoryId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	ids := make([]string, 0, len(objects))
	s := make([]map[string]interface{}, 0, len(objects))
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["alert_history_id"] = objectRaw["alertHistoryId"]
		mapping["status"] = objectRaw["status"]
		mapping["biz_source"] = objectRaw["bizSource"]
		mapping["end_time"] = objectRaw["endTime"]
		mapping["alert_rule_id"] = objectRaw["alertRuleId"]
		mapping["rule_display_name"] = objectRaw["ruleDisplayName"]
		mapping["start_time"] = objectRaw["startTime"]
		mapping["labels"] = objectRaw["labels"]
		mapping["count"] = objectRaw["count"]
		mapping["latest_level"] = objectRaw["latestLevel"]
		mapping["annotations"] = objectRaw["annotations"]
		mapping["instance_key"] = objectRaw["instanceKey"]
		mapping["max_level"] = objectRaw["maxLevel"]
		mapping["region_id"] = objectRaw["regionId"]
		mapping["workspace"] = objectRaw["workspace"]

		sendMaps := make([]map[string]interface{}, 0)
		sendMap := make(map[string]interface{})
		sendRaw := make(map[string]interface{})
		if objectRaw["send"] != nil {
			sendRaw = objectRaw["send"].(map[string]interface{})
		}
		if len(sendRaw) > 0 {
			actionMaps := make([]map[string]interface{}, 0)
			actionMap := make(map[string]interface{})
			actionRaw := make(map[string]interface{})
			if sendRaw["action"] != nil {
				actionRaw = sendRaw["action"].(map[string]interface{})
			}
			if len(actionRaw) > 0 {
				actionsRaw := make([]interface{}, 0)
				if actionRaw["actions"] != nil {
					actionsRaw = convertToInterfaceArray(actionRaw["actions"])
				}
				actionMap["actions"] = actionsRaw
				actionMaps = append(actionMaps, actionMap)
			}
			sendMap["action"] = actionMaps

			notificationMaps := make([]map[string]interface{}, 0)
			notificationMap := make(map[string]interface{})
			notificationRaw := make(map[string]interface{})
			if sendRaw["notification"] != nil {
				notificationRaw = sendRaw["notification"].(map[string]interface{})
			}
			if len(notificationRaw) > 0 {
				contactsRaw := make([]interface{}, 0)
				if notificationRaw["contacts"] != nil {
					contactsRaw = convertToInterfaceArray(notificationRaw["contacts"])
				}
				notificationMap["contacts"] = contactsRaw

				customWebhooksRaw := make([]interface{}, 0)
				if notificationRaw["customWebhooks"] != nil {
					customWebhooksRaw = convertToInterfaceArray(notificationRaw["customWebhooks"])
				}
				notificationMap["custom_webhooks"] = customWebhooksRaw

				dingWebhooksRaw := make([]interface{}, 0)
				if notificationRaw["dingWebhooks"] != nil {
					dingWebhooksRaw = convertToInterfaceArray(notificationRaw["dingWebhooks"])
				}
				notificationMap["ding_webhooks"] = dingWebhooksRaw

				fsWebhooksRaw := make([]interface{}, 0)
				if notificationRaw["fsWebhooks"] != nil {
					fsWebhooksRaw = convertToInterfaceArray(notificationRaw["fsWebhooks"])
				}
				notificationMap["fs_webhooks"] = fsWebhooksRaw

				groupsRaw := make([]interface{}, 0)
				if notificationRaw["groups"] != nil {
					groupsRaw = convertToInterfaceArray(notificationRaw["groups"])
				}
				notificationMap["groups"] = groupsRaw

				notificationMap["silence_time"] = notificationRaw["silenceTime"]

				slackWebhooksRaw := make([]interface{}, 0)
				if notificationRaw["slackWebhooks"] != nil {
					slackWebhooksRaw = convertToInterfaceArray(notificationRaw["slackWebhooks"])
				}
				notificationMap["slack_webhooks"] = slackWebhooksRaw

				wxWebhooksRaw := make([]interface{}, 0)
				if notificationRaw["wxWebhooks"] != nil {
					wxWebhooksRaw = convertToInterfaceArray(notificationRaw["wxWebhooks"])
				}
				notificationMap["wx_webhooks"] = wxWebhooksRaw

				notificationMaps = append(notificationMaps, notificationMap)
			}
			sendMap["notification"] = notificationMaps
			sendMaps = append(sendMaps, sendMap)
		}
		mapping["send"] = sendMaps

		ids = append(ids, fmt.Sprint(mapping["alert_history_id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("histories", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return WrapError(err)
		}
	}
	return nil
}
