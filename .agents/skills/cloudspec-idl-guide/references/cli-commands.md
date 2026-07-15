# Aliyun CLI cspec plugin 命令参考

源自 Aliyun CLI cspec plugin 的 AI 工作流，常用命令速查。

## 安装

```bash
aliyun plugin install --names cspec --source-base https://cli.aliyun-inc.com/registry_id/2/env/pre/plugins
```

## 核心命令

| 命令 | 用途 |
|------|------|
| `aliyun cspec build` | 编译和验证 cspec 文件 |
| `aliyun cspec check --name <组件名称>` | 检查指定组件的语法和规范（**推荐单组件校验**） |
| `aliyun cspec format` | 格式化 cspec 文件 |
| `aliyun cspec codefix` | 自动修复规范问题（需配合 -r/-t/-c 等参数） |

## 其他常用命令

| 命令 | 用途                    |
|------|-----------------------|
| `aliyun cspec create service` | 创建服务                  |
| `aliyun cspec create operation` | 创建操作                  |
| `aliyun cspec create resource` | 创建资源                  |
| `aliyun cspec name-gen` | 生成测试用例名称或 API 场景名称    |
| `aliyun cspec version` | 显示cli版本信息  确定是否要更新cli |

## 校验建议

- **推荐**：使用 `aliyun cspec check --name <组件名称>` 进行单组件校验
- **避免**：使用 `aliyun cspec check` 全量校验，建议逐个组件校验，便于定位问题

## 编辑后验证流程

1. 修改 .cspec 文件后，执行 `aliyun cspec build` 验证语法
2. 通过后，执行 `aliyun cspec check --name <组件名>` 进行规范检测
3. 若有规范问题，可尝试 `aliyun cspec codefix` 或手动修复，详见 [codefix-usage.md](codefix-usage.md)
