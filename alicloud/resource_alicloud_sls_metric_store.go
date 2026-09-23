package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlsMetricStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsMetricStoreCreate,
		Read:   resourceAliCloudSlsMetricStoreRead,
		Update: resourceAliCloudSlsMetricStoreUpdate,
		Delete: resourceAliCloudSlsMetricStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metric_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"shard_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"auto_split": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"max_split_shard_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      64,
				ValidateFunc: IntBetween(0, 256),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOkExists("auto_split"); ok && !v.(bool) {
						return true
					}
					return false
				},
			},
			"append_meta": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"hot_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"infrequent_access_ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntAtLeast(61),
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						return true
					}
					return old != "" && new != "" && old == new
				},
			},
			"encrypt_conf": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"encrypt_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"user_cmk_info": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cmk_key_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"arn": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildSlsMetricStore(d *schema.ResourceData) *sls.MetricStore {
	metricStore := &sls.MetricStore{
		Name:          d.Get("metric_store_name").(string),
		TTL:           d.Get("ttl").(int),
		ShardCount:    d.Get("shard_count").(int),
		AutoSplit:     d.Get("auto_split").(bool),
		MaxSplitShard: d.Get("max_split_shard_count").(int),
		AppendMeta:    d.Get("append_meta").(bool),
		Mode:          d.Get("mode").(string),
	}
	if hotTTL, ok := d.GetOk("hot_ttl"); ok {
		metricStore.HotTTL = int32(hotTTL.(int))
	}
	if infrequentAccessTTL, ok := d.GetOk("infrequent_access_ttl"); ok {
		ttl := int32(infrequentAccessTTL.(int))
		metricStore.InfrequentAccessTTL = &ttl
	}
	if encrypt := buildSlsMetricStoreEncrypt(d); encrypt != nil {
		metricStore.EncryptConf = encrypt
	}
	return metricStore
}

func buildSlsMetricStoreEncrypt(d *schema.ResourceData) *sls.MetricStoreEncryptConf {
	if field, ok := d.GetOk("encrypt_conf"); ok {
		list, ok := field.([]interface{})
		if !ok || len(list) == 0 {
			return nil
		}
		values, ok := list[0].(map[string]interface{})
		if !ok {
			return nil
		}
		encryptConf := &sls.MetricStoreEncryptConf{}
		if v, ok := values["enable"].(bool); ok {
			encryptConf.Enable = v
		}
		if v, ok := values["encrypt_type"].(string); ok && v != "" {
			encryptConf.EncryptType = v
		}
		cmkInfoList, _ := values["user_cmk_info"].([]interface{})
		if len(cmkInfoList) > 0 {
			if cmkInfo, ok := cmkInfoList[0].(map[string]interface{}); ok {
				userCmkInfo := &sls.MetricStoreEncryptUserCmkConf{}
				if v, ok := cmkInfo["cmk_key_id"].(string); ok {
					userCmkInfo.CmkKeyId = v
				}
				if v, ok := cmkInfo["arn"].(string); ok {
					userCmkInfo.Arn = v
				}
				if v, ok := cmkInfo["region_id"].(string); ok {
					userCmkInfo.RegionId = v
				}
				encryptConf.UserCmkInfo = userCmkInfo
			}
		}
		return encryptConf
	}
	return nil
}

func resourceAliCloudSlsMetricStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	projectName := d.Get("project_name").(string)
	metricStoreName := d.Get("metric_store_name").(string)
	metricStore := buildSlsMetricStore(d)

	var requestinfo *sls.Client
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.CreateMetricStoreV2(projectName, metricStore)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateMetricStoreV2", raw, requestinfo, map[string]interface{}{
			"project":     projectName,
			"metricStore": metricStore,
		})
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_metric_store", "CreateMetricStoreV2", AliyunLogGoSdkERROR)
	}

	d.SetId(fmt.Sprintf("%s:%s", projectName, metricStoreName))
	return resourceAliCloudSlsMetricStoreRead(d, meta)
}

func resourceAliCloudSlsMetricStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsLogStore(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[WARN] Resource alicloud_sls_metric_store %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("metric_store_name", objectRaw["logstoreName"])
	if objectRaw["ttl"] != nil {
		d.Set("ttl", objectRaw["ttl"])
	}
	if objectRaw["shardCount"] != nil {
		d.Set("shard_count", objectRaw["shardCount"])
	}
	if objectRaw["autoSplit"] != nil {
		d.Set("auto_split", objectRaw["autoSplit"])
	}
	if objectRaw["maxSplitShard"] != nil {
		d.Set("max_split_shard_count", objectRaw["maxSplitShard"])
	}
	if objectRaw["appendMeta"] != nil {
		d.Set("append_meta", objectRaw["appendMeta"])
	}
	if objectRaw["hot_ttl"] != nil {
		d.Set("hot_ttl", objectRaw["hot_ttl"])
	}
	if objectRaw["infrequentAccessTTL"] != nil {
		d.Set("infrequent_access_ttl", objectRaw["infrequentAccessTTL"])
	}
	if objectRaw["mode"] != nil {
		d.Set("mode", objectRaw["mode"])
	}
	if objectRaw["createTime"] != nil {
		d.Set("create_time", objectRaw["createTime"])
	}
	if objectRaw["lastModifyTime"] != nil {
		d.Set("last_modify_time", objectRaw["lastModifyTime"])
	}

	encryptConfMaps := make([]map[string]interface{}, 0)
	if encryptConfRaw, ok := objectRaw["encrypt_conf"].(map[string]interface{}); ok && len(encryptConfRaw) > 0 {
		encryptConfMap := map[string]interface{}{
			"enable":       encryptConfRaw["enable"],
			"encrypt_type": encryptConfRaw["encrypt_type"],
		}
		userCmkInfoMaps := make([]map[string]interface{}, 0)
		if userCmkInfoRaw, ok := encryptConfRaw["user_cmk_info"].(map[string]interface{}); ok && len(userCmkInfoRaw) > 0 {
			userCmkInfoMaps = append(userCmkInfoMaps, map[string]interface{}{
				"cmk_key_id": userCmkInfoRaw["cmk_key_id"],
				"arn":        userCmkInfoRaw["arn"],
				"region_id":  userCmkInfoRaw["region_id"],
			})
		}
		encryptConfMap["user_cmk_info"] = userCmkInfoMaps
		encryptConfMaps = append(encryptConfMaps, encryptConfMap)
	}
	if objectRaw["encrypt_conf"] != nil {
		if err := d.Set("encrypt_conf", encryptConfMaps); err != nil {
			return err
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("project_name", parts[0])

	return nil
}

func resourceAliCloudSlsMetricStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	projectName := parts[0]

	metricStore := buildSlsMetricStore(d)
	var requestinfo *sls.Client
	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.UpdateMetricStoreV2(projectName, metricStore)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("UpdateMetricStoreV2", raw, requestinfo, map[string]interface{}{
			"project":     projectName,
			"metricStore": metricStore,
		})
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateMetricStoreV2", AliyunLogGoSdkERROR)
	}

	return resourceAliCloudSlsMetricStoreRead(d, meta)
}

func resourceAliCloudSlsMetricStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	projectName := parts[0]
	metricStoreName := parts[1]

	action := fmt.Sprintf("/logstores/%s", metricStoreName)
	var request map[string]interface{}
	request = make(map[string]interface{})
	hostMap := map[string]*string{
		"project": StringPointer(projectName),
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err := client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteMetricStore", action), make(map[string]*string), nil, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteMetricStore", AlibabaCloudSdkGoERROR)
	}

	return nil
}
