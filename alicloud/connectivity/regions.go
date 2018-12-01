package connectivity

// Region represents ECS region
type Region string

// Constants of region definition
const (
	Hangzhou    = Region("cn-hangzhou")
	Qingdao     = Region("cn-qingdao")
	Beijing     = Region("cn-beijing")
	Hongkong    = Region("cn-hongkong")
	Shenzhen    = Region("cn-shenzhen")
	Shanghai    = Region("cn-shanghai")
	Zhangjiakou = Region("cn-zhangjiakou")
	Huhehaote   = Region("cn-huhehaote")

	APSouthEast1 = Region("ap-southeast-1")
	APNorthEast1 = Region("ap-northeast-1")
	APSouthEast2 = Region("ap-southeast-2")
	APSouthEast3 = Region("ap-southeast-3")
	APSouthEast5 = Region("ap-southeast-5")

	APSouth1 = Region("ap-south-1")

	USWest1 = Region("us-west-1")
	USEast1 = Region("us-east-1")

	MEEast1 = Region("me-east-1")

	EUCentral1 = Region("eu-central-1")
	EUWest1    = Region("eu-west-1")

	ShenZhenFinance = Region("cn-shenzhen-finance-1")
	ShanghaiFinance = Region("cn-shanghai-finance-1")
)

var ValidRegions = []Region{
	Hangzhou, Qingdao, Beijing, Shenzhen, Hongkong, Shanghai, Zhangjiakou, Huhehaote,
	USWest1, USEast1,
	APNorthEast1, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5,
	APSouth1,
	MEEast1,
	EUCentral1, EUWest1,
	ShenZhenFinance, ShanghaiFinance,
}

var EcsClassicSupportedRegions = []Region{Shenzhen, Shanghai, Beijing, Qingdao, Hangzhou, Hongkong, USWest1, APSouthEast1}
var SlbGuaranteedSupportedRegions = []Region{Qingdao, Beijing, Hangzhou, Shanghai, Shenzhen, Zhangjiakou, Huhehaote, APSouthEast1, USEast1}
var DrdsSupportedRegions = []Region{Beijing, Shenzhen, Hangzhou, Qingdao, Hongkong}
var DrdsClassicNoSupportedRegions = []Region{Hongkong}
var FcNoSupportedRegions = []Region{Zhangjiakou, Huhehaote, APSouthEast3, APSouthEast5, EUWest1, USEast1, MEEast1}
var DatahubSupportedRegions = []Region{Beijing, Hangzhou, Shanghai, Shenzhen, APSouthEast1}
var RdsMultiAzNoSupportedRegions = []Region{Qingdao, APNorthEast1, APSouthEast5, MEEast1}
var RouteTableNoSupportedRegions = []Region{Beijing, Hangzhou, Shenzhen}
var ApiGatewayNoSupportedRegions = []Region{Zhangjiakou, Huhehaote, USEast1, USWest1, EUWest1, MEEast1}
