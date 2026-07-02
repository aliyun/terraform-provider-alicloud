# 阿里云错误码官方定义查证(跨 skill 复用 reference)

> 适用场景:给定 `<ProductCode>` + `<ErrorCode>`,找该错误码的官方定义 —— HTTP 状态、原文 message(中/英)、是否建议重试、相邻/相关错误码。用于 aone-triage 工单查证、terraform-pr-review 评审错误码补丁、provider-resource-dev 排查客户报错。

## 一、查证入口(按顺序)

不凭记忆、不看博客定性;只信阿里云官方(help.aliyun.com / alibabacloud.com/help / next.api.aliyun.com / api.aliyun.com/troubleshoot)。

1. **help.aliyun.com 中文错误码表(首选,最新)**
   `https://help.aliyun.com/zh/{product-slug}/api-{product-code}-{yyyy-mm-dd}-errorcodes`
   例:ESA → `https://help.aliyun.com/zh/edge-security-acceleration/esa/api-esa-2024-09-10-errorcodes`

2. **alibabacloud.com 英文镜像**
   `https://www.alibabacloud.com/help/en/{product-slug}/api-{product-code}-{yyyy-mm-dd}-errorcodes`
   同页翻译,取英文 message 原文常用它。

3. **api.aliyun.com 通用诊断入口(无需产品,直接错误码搜)**
   `https://api.aliyun.com/troubleshoot?q={ErrorCode}`
   跨产品聚合,快速判断该错误码是否某产品独有;命中后再回入口 1/2 拿正式定义。

4. **next.api.aliyun.com 产品页 / API 页**
   `https://next.api.aliyun.com/product/{PRODUCT}`
   `https://next.api.aliyun.com/api/{PRODUCT}/{VERSION}/{ApiName}`
   产品页跳具体 API 的错误码块,API 页 URL 里带 `{VERSION}`(=错误码表 URL 里的 `yyyy-mm-dd`)。

5. **WebSearch 兜底(仅当上述均无命中,验证阿里云官方结果时)**
   关键词:`阿里云 {product} {ErrorCode} 错误码` / `alicloud {product} {ErrorCode} error code`。
   命中官方文档链接才用,第三方来源不作定性。

## 二、产品 slug 与 API 版本

**产品 slug 规律**(help.aliyun.com URL 里 `/zh/` 后那一段):

| 产品 | slug |
|---|---|
| ESA(边缘安全加速) | `edge-security-acceleration/esa` |
| OSS | `oss` |
| ECS | `elastic-compute-service` |
| ALB | `server-load-balancer/alb` |
| NLB | `server-load-balancer/nlb` |
| VPC | `virtual-private-cloud` |

**不确定 slug 就从 next.api.aliyun.com 产品页跳**:产品页右上角"帮助文档"链接会带 slug。

**API 版本日期怎么拿**:`next.api.aliyun.com/api/{PRODUCT}/{VERSION}/{ApiName}` 页面 URL 里就带 `{VERSION}`(格式 `yyyy-mm-dd`),把它拼进 errorcodes URL 即可。

## 三、输出模板(每次查证都按这个格式给结论)

```
Product: <name> (<version>)
Code: <ErrorCode>
HTTP: <status>
EN msg: "<英文原文>"
ZH msg: "<中文原文>"
Retry advice (official): <yes/no + 原文引用>
相邻/相关错误码: <其他并列 code + 一句话区别>
Source URLs:
  - <中文错误码表 URL>
  - <英文镜像 URL 或 API 页 URL>
```

## 四、反例 / 不要做的事

- **不臆测机制**:官方文档没写"乐观锁 / 悲观锁 / 队列 mutex / etcd revision 冲突"这类实现机制 → 一律不写;只保守说"服务端返回的可重试错误码"或"并发写冲突/前一请求处理中",跟官方 message 措辞对齐。
- **第三方博客、SO、CSDN、公众号只作旁证**:绝不用它们做官方定性;若官方无定义而只有社区讨论,报"官方未定义,仅有社区经验",不下结论。
- **一码一查,别跨产品套语义**:错误码名字看起来像锁(如 `LockFailed`),不同产品语义可能完全不同 —— 一个产品可能是"上一请求未完成",另一个可能是"分布式锁抢占失败"。永远回到该 product 的 errorcodes 表原文。
- **不省略 HTTP 状态**:HTTP 400/409/429/503 语义差别大,retry 语义与之强相关;缺 HTTP 就回补,不补就标 `HTTP: unknown`,别猜。
- **不合并"相邻错误码"进结论**:例如 `LockFailed` 与 `TooManyRequests` / `Site.ServiceBusy` 虽都触发"重试",但触发路径/接口层/退避策略完全不同,分列列出,一句话说清区别。

## 五、谁应该读这个

- **aone-triage**:工单里出现具体错误码(客户贴报错、SDK/CLI 报错、Terraform apply 报错)且需要判断"这是什么/该怎么回"时。
- **terraform-pr-review**:PR diff 涉及错误码补丁(NeedRetry / retryable errors 白名单 / IsExpectedErrors 等),核对补的 code 是否与官方 message 语义一致、HTTP 是否匹配。
- **provider-resource-dev**:开发中 acc 测/客户报错命中未识别错误码,先查官方定义,再决定加入 retry 列表还是暴露给用户。
