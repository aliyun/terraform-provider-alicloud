package config

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
	ChengDu     = Region("cn-chengdu")
	HeYuan      = Region("cn-heyuan")
	WuLanChaBu  = Region("cn-wulanchabu")
	GuangZhou   = Region("cn-guangzhou")
	NanJing     = Region("cn-nanjing")
	FuZhou      = Region("cn-fuzhou")
	ZhongWei    = Region("cn-zhongwei")

	APSouthEast1 = Region("ap-southeast-1")
	APNorthEast1 = Region("ap-northeast-1")
	APNorthEast2 = Region("ap-northeast-2")
	APSouthEast2 = Region("ap-southeast-2")
	APSouthEast3 = Region("ap-southeast-3")
	APSouthEast5 = Region("ap-southeast-5")
	APSouthEast6 = Region("ap-southeast-6")
	APSouthEast7 = Region("ap-southeast-7")
	APSouthEast8 = Region("ap-southeast-8")

	APSouth1 = Region("ap-south-1")

	USWest1      = Region("us-west-1")
	USEast1      = Region("us-east-1")
	USSouthEast1 = Region("us-southeast-1")

	NASouth1 = Region("na-south-1")

	MEEast1    = Region("me-east-1")
	MECentral1 = Region("me-central-1")

	EUCentral1 = Region("eu-central-1")
	EUWest1    = Region("eu-west-1")

	RusWest1 = Region("rus-west-1")

	HangzhouFinance     = Region("cn-hangzhou-finance")
	HangzhouFinanceOSS  = Region("cn-hzfinance")
	HangzhouFinanceOSS1 = Region("cn-hzjbp")

	BeijingFinance1   = Region("cn-beijing-finance-1")
	BeijingFinancePub = Region("cn-beijing-finance-1-pub")

	ShanghaiFinance     = Region("cn-shanghai-finance-1")
	ShanghaiFinance1Pub = Region("cn-shanghai-finance-1-pub")
	ShenZhenFinance1    = Region("cn-shenzhen-finance-1")
	ShenzhenFinance2    = Region("cn-szfinance")
	ShenzhenFinance     = Region("cn-shenzhen-finance")

	CnNorth2Gov1 = Region("cn-north-2-gov-1")
)

var ValidRegions = []Region{
	Hangzhou, Qingdao, Beijing, Shenzhen, Hongkong, Shanghai, Zhangjiakou, Huhehaote, ChengDu, HeYuan, WuLanChaBu, GuangZhou, NanJing, FuZhou, ZhongWei,
	USWest1, USEast1, USSouthEast1,
	NASouth1,
	APNorthEast1, APNorthEast2, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5, APSouthEast6, APSouthEast7, APSouthEast8,
	APSouth1,
	MEEast1, MECentral1,
	EUCentral1, EUWest1, RusWest1,
	CnNorth2Gov1,
	ShenZhenFinance1, ShenzhenFinance2, ShenzhenFinance,
	ShanghaiFinance1Pub, ShanghaiFinance,
	HangzhouFinance, HangzhouFinanceOSS, HangzhouFinanceOSS1,
	BeijingFinance1, BeijingFinancePub,
}
