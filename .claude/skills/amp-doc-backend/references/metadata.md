# API / Struct 元数据参考

受控 wrapper 的 `doc create` 接收两份独立 JSON 文本：`--doc-meta` 与 `--meta`。对于 API/Struct，它们会被后端原样分别传给 V1 `DraftDocMetaSaveDTO.docMeta` 与 `DraftDocMetaSaveDTO.apiMeta`，因此不能合并成一个 JSON。

## 响应与保存映射

| 文档类型 | `get` 响应 | `create` 参数 |
| --- | --- | --- |
| API | `docMetaJsonString` = `ApiDocMetaDTO`；`metaJsonString` = `AlibabaCloudOpenApi` | `--doc-meta` = ApiDocMetaDTO JSON；`--meta` = AlibabaCloudOpenApi JSON |
| Struct | `docMetaJsonString` = `StructDocMetaDTO`；`metaJsonString` = `Components` | `--doc-meta` = StructDocMetaDTO JSON；`--meta` = Components JSON |

`DocOpenApiDTO` 外层还有 `alibabaCloudOpenApiJsonString`，但保存 API/Struct 时以 `metaJsonString` 为准。绝不把 JSON 字符串再包成 `{ "metaJsonString": "..." }`。

## docmeta 最小初始化

只有已发布文档的其中一份元数据缺失时才使用模板；先保留存在的一份，再从当前实体定义构造缺失的另一份。

```json
// API: --doc-meta
{
  "title": "<API 名称>",
  "summary": "<API 的简短中文用途>",
  "paramExtraInfo": {}
}
```

```json
// Struct: --doc-meta
{
  "paramExtraInfo": {}
}
```

`BaseDocMetaDTO` 可包含 `cloudType`、`branchId`、`commitId`、`branchName` 和 `paramExtraInfo`。没有可靠来源时不要初始化前四项。API 的 `ApiDocMetaDTO` 还可包含 `runtimeType`、`docModelType`、`title`、`summary`、请求/响应补充说明、示例、SDK 与权限信息；仅在有可靠内容时填写。Struct 的 `StructDocMetaDTO` 除基础字段外没有专有字段。

`paramExtraInfo` 使用 JSONPath 到文档值的映射。例如：

```json
{
  "$.parameters[0].schema.description": "实例 ID。",
  "$.parameters[0].schema.example": "i-xxxxxxxx",
  "$.responses.200.schema.properties.InstanceId.description": "实例 ID。"
}
```

API 文档还可维护 `paramTagPosition`，按 `参数名-in` 记录该参数的文档 JSONPath 列表；已有值保留。新建时不能可靠定位路径就省略，不能猜测。

## apimeta 字段白名单

允许修改的 schema 字段仅为：

```text
description
title
example
enumValueTitles
```

其中 `example` 是完成“每个出入参有示例”的必要文档字段。它不改变运行时接口语义。`enumValueTitles` 用于已有枚举值的中文展示；默认不得替换 `enum` 原值。
枚举值集合的新增、删除或替换不属于文案变更，必须停止本流程并转分支 E。

### API 的可编辑位置

```text
paths.<path>.<operation>.parameters[*].schema
paths.<path>.<operation>.responses.<status>.schema.properties.*
```

沿 `properties` 和 `items` 递归。保留 `parameters`、`responses`、响应状态码及 schema 的其他全部内容。

### Struct 的可编辑位置

```text
components.schemas.<schema>.
components.schemas.<schema>.properties.*
components.schemas.<schema>.items
```

沿 `properties` 和 `items` 递归，不能更改 `components.schemas` 的结构或 `$ref` 图。

## 提交前差异检查

将原始 apimeta 与编辑后 apimeta 做递归结构比较：每一处新增、删除或值变化必须落在白名单字段。发现非白名单差异（包括 `required`、`type`、`format`、`default`、`pattern`、`minimum`、`maximum`、`$ref`、`properties` 或响应码）立即停止，不保存。docmeta 可更新其文案和 `paramExtraInfo`，但不得凭空设置版本、分支、权限或 SDK 配置。
