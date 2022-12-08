package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEdasK8sSlbAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasK8sSlbAttachmentCreate,
		Read:   resourceAlicloudEdasK8sSlbAttachmentRead,
		Update: resourceAlicloudEdasK8sSlbAttachmentUpdate,
		Delete: resourceAlicloudEdasK8sSlbAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"slb_configs": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"internet", "intranet"}, false),
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"rr", "wrr"}, false),
						},
						"specification": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"slb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "if field 'slb_id' is empty, EDAS will purchase a new slb for this config",
						},
						"port_mappings": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"loadbalancer_protocol": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
									},
									"service_port": {
										Type:     schema.TypeSet,
										Required: true,
										MaxItems: 1,
										MinItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(1, 65535),
												},
												"protocol": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
												},
												"target_port": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(1, 65535),
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

func resourceAlicloudEdasK8sSlbAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)

	if v, ok := d.GetOk("slb_configs"); ok {
		slbConfigs := v.(*schema.Set).List()
		if len(slbConfigs) > 0 {
			for _, c := range slbConfigs {
				config := c.(map[string]interface{})
				wait := incrementalWait(3*time.Second, 10*time.Second)
				err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
					if err := edasService.BindK8sSlb(appId, &config, d.Timeout(schema.TimeoutCreate)); err != nil {
						if err.Retryable {
							wait()
						}
						return err
					}
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, appId, "alicloud_edas_k8s_slb_attachment", AlibabaCloudSdkGoERROR)
				}
			}
		}
	}
	d.SetId(appId)
	return resourceAlicloudEdasK8sSlbAttachmentRead(d, meta)
}

func resourceAlicloudEdasK8sSlbAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}
	slbConfigs, err := edasService.DescribeEdasK8sSlbAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_edas_k8s_slb_attachment edasService.DescribeEdasK8sSlbAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "alicloud_edas_k8s_slb_attachment", AlibabaCloudSdkGoERROR)
	}
	d.Set("app_id", d.Id())
	if err := d.Set("slb_configs", slbConfigs); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "alicloud_edas_k8s_slb_attachment", AlibabaCloudSdkGoERROR)
	}

	return nil
}

func parseSlbConfig(info *map[string]interface{}) (config *map[string]interface{}, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("parse slb config err, %v", p)
		}
	}()
	pms := (*info)["portMappings"].([]interface{})
	var mappings []map[string]interface{}
	for _, opm := range pms {
		pm := opm.(map[string]interface{})
		osp := pm["servicePort"].(map[string]interface{})
		sp := []map[string]interface{}{
			{
				"port":        int64(osp["port"].(float64)),
				"protocol":    osp["protocol"],
				"target_port": int64(osp["targetPort"].(float64)),
			},
		}
		mapping := map[string]interface{}{
			"cert_id":               pm["certId"],
			"loadbalancer_protocol": pm["loadBalancerProtocol"],
			"service_port":          sp,
		}
		mappings = append(mappings, mapping)
	}
	return &map[string]interface{}{
		"name":          (*info)["name"],
		"type":          (*info)["addressType"],
		"slb_id":        (*info)["slbId"],
		"scheduler":     (*info)["scheduler"],
		"specification": (*info)["specification"],
		"port_mappings": mappings,
	}, nil
}

func jsonEmpty(s string) bool {
	return s == "" || s == "{}" || s == "[]"
}

func filterSlbInfo(slbInfo string) (*[]map[string]interface{}, error) {
	var slbInfos []map[string]interface{}
	var filteredSlbInfo []map[string]interface{}
	if slbInfo == "" {
		return &filteredSlbInfo, nil
	}
	err := json.Unmarshal([]byte(slbInfo), &slbInfos)
	if err != nil {
		return nil, WrapErrorf(err, "unmarshal slb info failed, value: %s", slbInfo)
	}
	for i := range slbInfos {
		slb := &slbInfos[i]
		if v, ok := (*slb)["addressType"]; ok && (v.(string) == "internet" || v.(string) == "intranet") {
			filteredSlbInfo = append(filteredSlbInfo, *slb)
		}
	}
	return &filteredSlbInfo, nil
}

func resourceAlicloudEdasK8sSlbAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	d.Partial(true)
	if d.HasChange("slb_configs") {
		o, n := d.GetChange("slb_configs")
		oldConfigs := o.(*schema.Set).List()
		newConfigs := n.(*schema.Set).List()
		wait := incrementalWait(3*time.Second, 10*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			if err := edasService.UpdateK8sAppSlbInfos(d.Id(), &oldConfigs, &newConfigs, d.Timeout(schema.TimeoutUpdate)); err != nil {
				if err.Retryable {
					wait()
				}
				return err
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "alicloud_edas_k8s_slb_attachment", AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("slb_configs")
	}
	d.Partial(false)
	return resourceAlicloudEdasK8sSlbAttachmentRead(d, meta)
}

func resourceAlicloudEdasK8sSlbAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}
	if v, ok := d.GetOk("slb_configs"); ok {
		slbConfigs := v.(*schema.Set).List()
		if len(slbConfigs) > 0 {
			for _, c := range slbConfigs {
				config := c.(map[string]interface{})
				slbType := config["type"].(string)
				slbName := config["name"].(string)
				wait := incrementalWait(3*time.Second, 10*time.Second)
				err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
					if err := edasService.UnbindK8sSlb(d.Id(), slbType, slbName, d.Timeout(schema.TimeoutDelete)); err != nil {
						if err.Retryable {
							wait()
						}
						return err
					}
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), "alicloud_edas_k8s_slb_attachment", AlibabaCloudSdkGoERROR)
				}
			}
		}
	}
	return nil
}
