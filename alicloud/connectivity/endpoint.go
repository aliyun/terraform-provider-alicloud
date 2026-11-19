package connectivity

import "github.com/aliyun/terraform-provider-alicloud/alicloud/client/config"

type ServiceCode = config.ServiceCode

const (
	MaxcomputeCode      = config.MaxcomputeCode
	CmsCode             = config.CmsCode
	RKvstoreCode        = config.RKvstoreCode
	OnsCode             = config.OnsCode
	DcdnCode            = config.DcdnCode
	MseCode             = config.MseCode
	ActiontrailCode     = config.ActiontrailCode
	OosCode             = config.OosCode
	EcsCode             = config.EcsCode
	NasCode             = config.NasCode
	EciCode             = config.EciCode
	DdoscooCode         = config.DdoscooCode
	BssopenapiCode      = config.BssopenapiCode
	AlidnsCode          = config.AlidnsCode
	ResourcemanagerCode = config.ResourcemanagerCode
	WafOpenapiCode      = config.WafOpenapiCode
	DmsEnterpriseCode   = config.DmsEnterpriseCode
	DnsCode             = config.DnsCode
	KmsCode             = config.KmsCode
	CbnCode             = config.CbnCode
	ECSCode             = config.ECSCode
	ESSCode             = config.ESSCode
	RAMCode             = config.RAMCode
	VPCCode             = config.VPCCode
	SLBCode             = config.SLBCode
	RDSCode             = config.RDSCode
	OSSCode             = config.OSSCode
	ONSCode             = config.ONSCode
	ALIKAFKACode        = config.ALIKAFKACode
	CONTAINCode         = config.CONTAINCode
	CRCode              = config.CRCode
	CDNCode             = config.CDNCode
	CMSCode             = config.CMSCode
	KMSCode             = config.KMSCode
	OTSCode             = config.OTSCode
	DNSCode             = config.DNSCode
	PVTZCode            = config.PVTZCode
	LOGCode             = config.LOGCode
	FCCode              = config.FCCode
	DDSCode             = config.DDSCode
	GPDBCode            = config.GPDBCode
	STSCode             = config.STSCode
	KVSTORECode         = config.KVSTORECode
	POLARDBCode         = config.POLARDBCode
	DATAHUBCode         = config.DATAHUBCode
	MNSCode             = config.MNSCode
	CLOUDAPICode        = config.CLOUDAPICode
	DRDSCode            = config.DRDSCode
	LOCATIONCode        = config.LOCATIONCode
	ELASTICSEARCHCode   = config.ELASTICSEARCHCode
	BSSOPENAPICode      = config.BSSOPENAPICode
	DDOSCOOCode         = config.DDOSCOOCode
	DDOSBGPCode         = config.DDOSBGPCode
	SAGCode             = config.SAGCode
	EMRCode             = config.EMRCode
	CasCode             = config.CasCode
	YUNDUNDBAUDITCode   = config.YUNDUNDBAUDITCode
	MARKETCode          = config.MARKETCode
	HBASECode           = config.HBASECode
	ADBCode             = config.ADBCode
	MAXCOMPUTECode      = config.MAXCOMPUTECode
	EDASCode            = config.EDASCode
	CassandraCode       = config.CassandraCode
)

type Endpoints = config.Endpoints
type Endpoint = config.Endpoint
type RegionIds = config.RegionIds
type Products = config.Products
type Product = config.Product

type ApiVersion = config.ApiVersion

const (
	ApiVersion20140526 = config.ApiVersion20140526
	ApiVersion20160815 = config.ApiVersion20160815
	ApiVersion20140515 = config.ApiVersion20140515
)

func LoadRegionalEndpoint(region string, serviceCode string) string {
	return config.LoadRegionalEndpoint(region, serviceCode)
}

const (
	OpenApiGatewayService          = config.OpenApiGatewayService
	OpenOtsService                 = config.OpenOtsService
	OpenOssService                 = config.OpenOssService
	OpenNasService                 = config.OpenNasService
	OpenCdnService                 = config.OpenCdnService
	OpenKmsService                 = config.OpenKmsService
	OpenSaeService                 = config.OpenSaeService
	OpenCmsService                 = config.OpenCmsService
	OpenDatahubService             = config.OpenDatahubService
	OpenOnsService                 = config.OpenOnsService
	OpenDcdnService                = config.OpenDcdnService
	OpenFcService                  = config.OpenFcService
	OpenAckService                 = config.OpenAckService
	OpenPrivateLinkService         = config.OpenPrivateLinkService
	OpenBrainIndustrialService     = config.OpenBrainIndustrialService
	OpenIotService                 = config.OpenIotService
	OpenVsService                  = config.OpenVsService
	OpenCrService                  = config.OpenCrService
	OpenMaxcomputeService          = config.OpenMaxcomputeService
	OpenCloudStorageGatewayService = config.OpenCloudStorageGatewayService
	DataWorksService               = config.DataWorksService
	OpenHbrService                 = config.OpenHbrService
)

const (
	BssOpenAPIEndpointDomestic                = config.BssOpenAPIEndpointDomestic
	BssOpenAPIEndpointInternational           = config.BssOpenAPIEndpointInternational
	EcdOpenAPIEndpointUser                    = config.EcdOpenAPIEndpointUser
	CloudFirewallOpenAPIEndpointControlPolicy = config.CloudFirewallOpenAPIEndpointControlPolicy
	SaSOpenAPIEndpointInternational           = config.SaSOpenAPIEndpointInternational
)
