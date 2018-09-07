package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
				Required:     true,
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
				ForceNew:     true,
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
				Default:  "Redis",
				ValidateFunc: validateAllowedStringValue([]string{
					"Memcache",
					"Redis",
				}),
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"engine_version": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "2.8",
				ValidateFunc: validateAllowedStringValue([]string{
					"2.8",
					"4.0",
				}),
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
	client := meta.(*AliyunClient)
	conn := client.rkvconn

	request, err := buildKVStoreCreateRequest(d, meta)
	if err != nil {
		return err
	}

	resp, err := conn.CreateInstance(request)

	if err != nil {
		return fmt.Errorf("Error creating Alicloud db instance: %#v", err)
	}

	d.SetId(resp.InstanceId)

	// wait instance status change from Creating to Normal
	if err := client.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
	}

	return resourceAlicloudKVStoreInstanceUpdate(d, meta)
}

func resourceAlicloudKVStoreInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rkvconn
	d.Partial(true)

	if d.HasChange("security_ips") {
		request := r_kvstore.CreateModifySecurityIpsRequest()
		request.SecurityIpGroupName = "default"
		request.InstanceId = d.Id()
		if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
			request.SecurityIps = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
		} else {
			return fmt.Errorf("Security ips cannot be empty")
		}
		// wait instance status is Normal before modifying
		if err := client.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}
		if _, err := conn.ModifySecurityIps(request); err != nil {
			return fmt.Errorf("Create security whitelist ips got an error: %#v", err)
		}
		d.SetPartial("security_ips")
		// wait instance status is Normal after modifying
		if err := client.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudKVStoreInstanceRead(d, meta)
	}

	if d.HasChange("instance_class") {
		request := r_kvstore.CreateModifyInstanceSpecRequest()
		request.InstanceId = d.Id()
		request.InstanceClass = d.Get("instance_class").(string)
		request.EffectiveTime = "Immediately"
		if _, err := conn.ModifyInstanceSpec(request); err != nil {
			return err
		}
		// wait instance status is Normal after modifying
		if err := client.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}

		d.SetPartial("instance_class")
	}

	request := r_kvstore.CreateModifyInstanceAttributeRequest()
	request.InstanceId = d.Id()
	update := false
	if d.HasChange("instance_name") {
		request.InstanceName = d.Get("instance_name").(string)
		update = true

		d.SetPartial("instance_name")
	}

	if d.HasChange("password") {
		request.NewPassword = d.Get("password").(string)
		update = true
		d.SetPartial("password")
	}

	if update {
		// wait instance status is Normal before modifying
		if err := client.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}
		if _, err := conn.ModifyInstanceAttribute(request); err != nil {
			return fmt.Errorf("ModifyRKVInstanceDescription got an error: %#v", err)
		}
		// wait instance status is Normal after modifying
		if err := client.WaitForRKVInstance(d.Id(), Normal, DefaultLongTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}
	}

	d.Partial(false)
	return resourceAlicloudKVStoreInstanceRead(d, meta)
}

func resourceAlicloudKVStoreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	instance, err := client.DescribeRKVInstanceById(d.Id())
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
	client := meta.(*AliyunClient)

	instance, err := client.DescribeRKVInstanceById(d.Id())
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

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.rkvconn.DeleteInstance(request)

		if err != nil {
			if IsExceptedError(err, InvalidKVStoreInstanceIdNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete KVStore instance timeout and got an error: %#v", err))
		}

		if _, err := client.DescribeRKVInstanceById(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error Describe KVStore InstanceAttribute: %#v", err))
		}

		return resource.RetryableError(fmt.Errorf("Delete KVStore instance timeout and got an error: %#v", err))
	})
}

func buildKVStoreCreateRequest(d *schema.ResourceData, meta interface{}) (*r_kvstore.CreateInstanceRequest, error) {
	client := meta.(*AliyunClient)
	request := r_kvstore.CreateCreateInstanceRequest()
	request.InstanceName = Trim(d.Get("instance_name").(string))
	request.RegionId = getRegionId(d, meta)
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.InstanceClass = Trim(d.Get("instance_class").(string))
	request.ChargeType = Trim(d.Get("instance_charge_type").(string))
	request.Password = Trim(d.Get("password").(string))
	request.BackupId = Trim(d.Get("backup_id").(string))

	if PayType(request.ChargeType) == PrePaid {
		request.Period = d.Get("Period").(string)
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
		vsw, err := client.DescribeVswitch(vswitchId.(string))
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

	request.Token = buildClientToken("TF-CreateKVStoreInstance")

	return request, nil
}
