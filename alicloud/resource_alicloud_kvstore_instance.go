package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudKVStoreInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKVStoreInstanceCreate,
		Read:   resourceAlicloudKVStoreInstanceRead,
		Update: resourceAlicloudKVStoreInstanceUpdate,
		Delete: resourceAlicloudKVStoreInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRKVInstanceName,
			},
			"password": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validateRKVPassword,
			},
			"instance_class": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateInstanceChargeType,
				Optional:     true,
				Default:      PostPaid,
			},
			"period": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: rkvPostPaidDiffSuppressFunc,
			},
			"instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(KVStoreRedis),
				ValidateFunc: validateAllowedStringValue([]string{
					string(KVStoreMemcache),
					string(KVStoreRedis),
				}),
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"engine_version": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      KVStore2Dot8,
				ValidateFunc: validateAllowedStringValue([]string{string(KVStore2Dot8), string(KVStore4Dot0)}),
			},
			"connection_domain": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"backup_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"security_ips": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudKVStoreInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}

	request, err := buildKVStoreCreateRequest(d, meta)
	if err != nil {
		return err
	}

	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.CreateInstance(request)
	})

	if err != nil {
		return fmt.Errorf("Error creating Alicloud db instance: %#v", err)
	}
	resp, _ := raw.(*r_kvstore.CreateInstanceResponse)
	d.SetId(resp.InstanceId)

	// wait instance status change from Creating to Normal
	if err := kvstoreService.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Normal, err)
	}

	return resourceAlicloudKVStoreInstanceUpdate(d, meta)
}

func resourceAlicloudKVStoreInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}
	d.Partial(true)

	if d.HasChange("security_ips") {
		// wait instance status is Normal before modifying
		if err := kvstoreService.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Normal, err)
		}
		request := r_kvstore.CreateModifySecurityIpsRequest()
		request.SecurityIpGroupName = "default"
		request.InstanceId = d.Id()
		if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
			request.SecurityIps = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
		} else {
			return fmt.Errorf("Security ips cannot be empty")
		}
		_, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifySecurityIps(request)
		})
		if err != nil {
			return fmt.Errorf("Create security whitelist ips got an error: %#v", err)
		}
		d.SetPartial("security_ips")
		// wait instance status is Normal after modifying
		if err := kvstoreService.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Normal, err)
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudKVStoreInstanceRead(d, meta)
	}

	if d.HasChange("instance_class") {
		// wait instance status is Normal before modifying
		if err := kvstoreService.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Normal, err)
		}
		request := r_kvstore.CreateModifyInstanceSpecRequest()
		request.InstanceId = d.Id()
		request.InstanceClass = d.Get("instance_class").(string)
		request.EffectiveTime = "Immediately"
		_, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifyInstanceSpec(request)
		})
		if err != nil {
			return err
		}
		// wait instance status is Normal after modifying
		if err := kvstoreService.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Normal, err)
		}
		// There needs more time to sync instance class update
		if err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			instance, err := kvstoreService.DescribeRKVInstanceById(d.Id())
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if instance.InstanceClass != request.InstanceClass {
				return resource.RetryableError(fmt.Errorf("Waitting for instance class is changed timeout. "+
					"Expect instance class %s, got %s.", instance.InstanceClass, request.InstanceClass))
			}
			return nil
		}); err != nil {
			return err
		}

		d.SetPartial("instance_class")
	}

	request := r_kvstore.CreateModifyInstanceAttributeRequest()
	request.InstanceId = d.Id()
	update := false
	if d.HasChange("instance_name") {
		request.InstanceName = d.Get("instance_name").(string)
		update = true
	}

	if d.HasChange("password") {
		request.NewPassword = d.Get("password").(string)
		update = true
	}

	if d.HasChange("instance_charge_type") || d.HasChange("period") {
		prePaidRequest := r_kvstore.CreateTransformToPrePaidRequest()
		prePaidRequest.InstanceId = d.Id()
		prePaidRequest.Period = requests.Integer(strconv.Itoa(d.Get("period").(int)))

		// for now we just support charge change from PostPaid to PrePaid
		configPayType := PayType(d.Get("instance_charge_type").(string))
		if configPayType == PrePaid {
			_, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
				return rkvClient.TransformToPrePaid(prePaidRequest)
			})
			if err != nil {
				return fmt.Errorf("TransformToPrePaid got an error: %#v", err)
			}
			// wait instance status is Normal after modifying
			if err := kvstoreService.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
				return fmt.Errorf("WaitForInstance %s got error: %#v", Normal, err)
			}
			d.SetPartial("instance_charge_type")
			d.SetPartial("period")
		}
	}

	if update {
		// wait instance status is Normal before modifying
		if err := kvstoreService.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Normal, err)
		}
		_, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifyInstanceAttribute(request)
		})
		if err != nil {
			return fmt.Errorf("ModifyRKVInstanceAttribute got an error: %#v", err)
		}
		d.SetPartial("instance_name")
		d.SetPartial("password")
		// wait instance status is Normal after modifying
		if err := kvstoreService.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Normal, err)
		}
	}

	d.Partial(false)
	return resourceAlicloudKVStoreInstanceRead(d, meta)
}

func resourceAlicloudKVStoreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}
	instance, err := kvstoreService.DescribeRKVInstanceById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe rKV InstanceAttribute: %#v", err)
	}
	d.Set("instance_name", instance.InstanceName)
	d.Set("instance_class", instance.InstanceClass)
	d.Set("availability_zone", instance.ZoneId)
	d.Set("instance_charge_type", instance.ChargeType)
	d.Set("instance_type", instance.InstanceType)
	d.Set("vswitch_id", instance.VSwitchId)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("connection_domain", instance.ConnectionDomain)
	d.Set("private_ip", instance.PrivateIp)
	d.Set("security_ips", strings.Split(instance.SecurityIPList, COMMA_SEPARATED))

	return nil
}

func resourceAlicloudKVStoreInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}

	instance, err := kvstoreService.DescribeRKVInstanceById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return fmt.Errorf("Error Describe KVStore InstanceAttribute: %#v", err)
	}
	if PayType(instance.ChargeType) == Prepaid {
		return fmt.Errorf("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically")
	}
	request := r_kvstore.CreateDeleteInstanceRequest()
	request.InstanceId = d.Id()

	return resource.Retry(8*time.Minute, func() *resource.RetryError {
		_, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DeleteInstance(request)
		})

		if err != nil {
			if IsExceptedError(err, InvalidKVStoreInstanceIdNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete KVStore instance timeout and got an error: %#v", err))
		}

		if _, err := kvstoreService.DescribeRKVInstanceById(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error Describe KVStore InstanceAttribute: %#v", err))
		}

		return resource.RetryableError(fmt.Errorf("Delete KVStore instance timeout and got an error: %#v", err))
	})
}

func buildKVStoreCreateRequest(d *schema.ResourceData, meta interface{}) (*r_kvstore.CreateInstanceRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := r_kvstore.CreateCreateInstanceRequest()
	request.InstanceName = Trim(d.Get("instance_name").(string))
	request.RegionId = client.RegionId
	request.InstanceType = Trim(d.Get("instance_type").(string))
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	if request.InstanceType == string(KVStoreMemcache) && request.EngineVersion == string(KVStore4Dot0) {
		return nil, fmt.Errorf("Currently Memcache instance only supports engine version 2.8.")
	}
	request.InstanceClass = Trim(d.Get("instance_class").(string))
	request.ChargeType = Trim(d.Get("instance_charge_type").(string))
	request.Password = Trim(d.Get("password").(string))
	request.BackupId = Trim(d.Get("backup_id").(string))

	if PayType(request.ChargeType) == PrePaid {
		request.Period = strconv.Itoa(d.Get("period").(int))
	}

	if zone, ok := d.GetOk("availability_zone"); ok && Trim(zone.(string)) != "" {
		request.ZoneId = Trim(zone.(string))
	}

	request.NetworkType = strings.ToUpper(string(Classic))
	if vswitchId, ok := d.GetOk("vswitch_id"); ok && vswitchId.(string) != "" {
		request.VSwitchId = vswitchId.(string)
		request.NetworkType = strings.ToUpper(string(Vpc))
		request.PrivateIpAddress = Trim(d.Get("private_ip").(string))

		// check vswitchId in zone
		vsw, err := vpcService.DescribeVswitch(vswitchId.(string))
		if err != nil {
			return nil, fmt.Errorf("DescribeVSwitch got an error: %#v", err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return nil, fmt.Errorf("The specified vswitch %s isn't in the multi zone %s", vsw.VSwitchId, request.ZoneId)
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, fmt.Errorf("The specified vswitch %s isn't in the zone %s", vsw.VSwitchId, request.ZoneId)
		}

		request.VpcId = vsw.VpcId
	}

	request.Token = buildClientToken(fmt.Sprintf("TF-Create%sInstance", request.InstanceType))

	return request, nil
}
