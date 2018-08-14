package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunNatGatewayCreate,
		Read:   resourceAliyunNatGatewayRead,
		Update: resourceAliyunNatGatewayUpdate,
		Delete: resourceAliyunNatGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"spec": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'spec' has been deprecated from provider version 1.7.1, and new field 'specification' can replace it.",
			},
			"specification": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNatGatewaySpec,
				Default:      NatGatewaySmallSpec,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"bandwidth_package_ids": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"snat_table_ids": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"forward_table_ids": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"bandwidth_packages": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_count": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"bandwidth": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"zone": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"public_ip_addresses": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				MaxItems: 4,
				Optional: true,
			},
		},
	}
}

func resourceAliyunNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	args := vpc.CreateCreateNatGatewayRequest()
	args.RegionId = string(getRegion(d, meta))
	args.VpcId = string(d.Get("vpc_id").(string))
	args.Spec = string(d.Get("specification").(string))

	bandwidthPackages := []vpc.CreateNatGatewayBandwidthPackage{}
	for _, e := range d.Get("bandwidth_packages").([]interface{}) {
		pack := e.(map[string]interface{})
		bandwidthPackage := vpc.CreateNatGatewayBandwidthPackage{
			IpCount:   strconv.Itoa(pack["ip_count"].(int)),
			Bandwidth: strconv.Itoa(pack["bandwidth"].(int)),
		}
		if pack["zone"].(string) != "" {
			bandwidthPackage.Zone = pack["zone"].(string)
		}
		bandwidthPackages = append(bandwidthPackages, bandwidthPackage)
	}

	args.BandwidthPackage = &bandwidthPackages

	if v, ok := d.GetOk("name"); ok {
		args.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		args.Description = v.(string)
	}

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ar := args
		resp, err := conn.CreateNatGateway(ar)
		if err != nil {
			if IsExceptedError(err, VswitchStatusError) || IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(fmt.Errorf("CreateNatGateway got error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("CreateNatGateway got error: %#v", err))
		}
		d.SetId(resp.NatGatewayId)
		return nil
	}); err != nil {
		return err
	}

	return resourceAliyunNatGatewayRead(d, meta)
}

func resourceAliyunNatGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	natGateway, err := client.DescribeNatGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", natGateway.Name)
	d.Set("specification", natGateway.Spec)
	d.Set("bandwidth_package_ids", strings.Join(natGateway.BandwidthPackageIds.BandwidthPackageId, ","))
	d.Set("snat_table_ids", strings.Join(natGateway.SnatTableIds.SnatTableId, ","))
	d.Set("forward_table_ids", strings.Join(natGateway.ForwardTableIds.ForwardTableId, ","))
	d.Set("description", natGateway.Description)
	d.Set("vpc_id", natGateway.VpcId)

	bindWidthPackages, err := flattenBandWidthPackages(natGateway.BandwidthPackageIds.BandwidthPackageId, meta, d)
	if err != nil {
		log.Printf("[ERROR] bindWidthPackages flattenBandWidthPackages failed. natgateway id is %#v", d.Id())
	} else {
		d.Set("bandwidth_packages", bindWidthPackages)
	}

	return nil
}

func resourceAliyunNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	conn := client.vpcconn

	natGateway, err := client.DescribeNatGateway(d.Id())
	if err != nil {
		return err
	}

	d.Partial(true)
	attributeUpdate := false
	args := vpc.CreateModifyNatGatewayAttributeRequest()
	args.RegionId = natGateway.RegionId
	args.NatGatewayId = natGateway.NatGatewayId

	if d.HasChange("name") {
		d.SetPartial("name")
		var name string
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		} else {
			return fmt.Errorf("cann't change name to empty string")
		}
		args.Name = name

		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		var description string
		if v, ok := d.GetOk("description"); ok {
			description = v.(string)
		} else {
			return fmt.Errorf("can to change description to empty string")
		}

		args.Description = description

		attributeUpdate = true
	}

	if attributeUpdate {
		if _, err := conn.ModifyNatGatewayAttribute(args); err != nil {
			return err
		}
	}

	if d.HasChange("specification") {
		d.SetPartial("specification")
		request := vpc.CreateModifyNatGatewaySpecRequest()
		request.RegionId = natGateway.RegionId
		request.NatGatewayId = natGateway.NatGatewayId
		request.Spec = d.Get("specification").(string)

		if _, err := conn.ModifyNatGatewaySpec(request); err != nil {
			return fmt.Errorf("ModifyNatGatewaySpec got an error: %#v with args: %#v", err, *args)
		}

	}
	d.Partial(false)

	return resourceAliyunNatGatewayRead(d, meta)
}

func resourceAliyunNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	conn := client.vpcconn

	packRequest := vpc.CreateDescribeBandwidthPackagesRequest()
	packRequest.RegionId = string(getRegion(d, meta))
	packRequest.NatGatewayId = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {

		resp, err := conn.DescribeBandwidthPackages(packRequest)
		if err != nil {
			log.Printf("[ERROR] Describe bandwidth package is failed, natGateway Id: %s", d.Id())
			return resource.NonRetryableError(err)
		}

		retry := false
		if resp != nil && len(resp.BandwidthPackages.BandwidthPackage) > 0 {
			for _, pack := range resp.BandwidthPackages.BandwidthPackage {
				request := vpc.CreateDeleteBandwidthPackageRequest()
				request.RegionId = string(getRegion(d, meta))
				request.BandwidthPackageId = pack.BandwidthPackageId
				if _, err := conn.DeleteBandwidthPackage(request); err != nil {
					if IsExceptedError(err, NatGatewayInvalidRegionId) {
						log.Printf("[ERROR] Delete bandwidth package is failed, bandwidthPackageId: %#v", pack.BandwidthPackageId)
						return resource.NonRetryableError(err)
					}
					retry = true
				}
			}
		}

		if retry {
			return resource.RetryableError(fmt.Errorf("Delete bandwidth package timeout and got an error: %#v.", err))
		}

		args := vpc.CreateDeleteNatGatewayRequest()
		args.RegionId = string(getRegion(d, meta))
		args.NatGatewayId = d.Id()

		if _, err := conn.DeleteNatGateway(args); err != nil {
			if IsExceptedError(err, DependencyViolationBandwidthPackages) {
				return resource.RetryableError(fmt.Errorf("Delete nat gateway timeout and got an error: %#v.", err))
			}
			if IsExceptedError(err, InvalidNatGatewayIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		nat, err := client.DescribeNatGateway(d.Id())

		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			log.Printf("[ERROR] Describe NatGateways failed.")
			return resource.NonRetryableError(err)
		} else if nat.NatGatewayId != d.Id() {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Delete nat gateway timeout and got an error: %#v.", err))
	})
}

func flattenBandWidthPackages(bandWidthPackageIds []string, meta interface{}, d *schema.ResourceData) ([]map[string]interface{}, error) {
	packageLen := len(bandWidthPackageIds)
	result := make([]map[string]interface{}, 0, packageLen)
	for i := packageLen - 1; i >= 0; i-- {
		packageId := bandWidthPackageIds[i]
		packages, err := getPackages(packageId, meta, d)
		if err != nil {
			return result, err
		}
		ipAddress := flattenPackPublicIp(packages.PublicIpAddresses.PublicIpAddresse)
		ipCont, ipContErr := strconv.Atoi(packages.IpCount)
		bandWidth, bandWidthErr := strconv.Atoi(packages.Bandwidth)
		if ipContErr != nil {
			return result, ipContErr
		}
		if bandWidthErr != nil {
			return result, bandWidthErr
		}
		l := map[string]interface{}{
			"ip_count":            ipCont,
			"bandwidth":           bandWidth,
			"zone":                packages.ZoneId,
			"public_ip_addresses": ipAddress,
		}
		result = append(result, l)
	}
	return result, nil
}
func getPackages(packageId string, meta interface{}, d *schema.ResourceData) (pack vpc.BandwidthPackage, err error) {
	req := vpc.CreateDescribeBandwidthPackagesRequest()
	req.NatGatewayId = d.Id()
	req.BandwidthPackageId = packageId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		packages, err := meta.(*AliyunClient).vpcconn.DescribeBandwidthPackages(req)
		if err != nil {
			return err
		}
		if packages == nil || len(packages.BandwidthPackages.BandwidthPackage) < 1 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("Bandwidth Package", packageId))
		}
		pack = packages.BandwidthPackages.BandwidthPackage[0]
		return nil
	})
	return
}
func flattenPackPublicIp(publicIpAddressList []vpc.PublicIpAddresse) string {
	var result []string
	for _, publicIpAddresses := range publicIpAddressList {
		ipAddress := publicIpAddresses.IpAddress
		result = append(result, ipAddress)
	}
	return strings.Join(result, ",")
}
