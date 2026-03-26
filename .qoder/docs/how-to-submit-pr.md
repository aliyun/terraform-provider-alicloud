# 如何向开源项目提交 PR（Pull Request）

> 以本项目 `aliyun/terraform-provider-alicloud` 为例，通俗讲解完整流程。

---

## 一、搞清楚三个概念

| 概念 | 通俗理解 |
|------|---------|
| **Fork** | 把别人的仓库"复制"一份到你自己的 GitHub 账号下 |
| **Branch（分支）** | 在你的副本上开辟一条独立的修改通道，不影响主线 |
| **Pull Request（PR）** | 向原作者说"我改了这些，请你合并到官方仓库" |

---

## 二、完整步骤

### 第 1 步：Fork 官方仓库

1. 打开官方仓库页面：https://github.com/aliyun/terraform-provider-alicloud
2. 点击右上角 **Fork** 按钮
3. 完成后你的账号下会出现：`https://github.com/你的用户名/terraform-provider-alicloud`

---

### 第 2 步：克隆到本地（首次）

```bash
git clone https://github.com/你的用户名/terraform-provider-alicloud.git
cd terraform-provider-alicloud
```

> 如果已经 clone 了官方仓库，跳到第 3 步添加远程地址即可。

---

### 第 3 步：添加远程地址（已有本地仓库时）

```bash
# 查看当前远程地址
git remote -v

# 添加你的 fork 仓库为新的远程（命名为 fork）
git remote add fork https://github.com/你的用户名/terraform-provider-alicloud.git
```

---

### 第 4 步：创建功能分支

```bash
# 从最新的 master 拉取
git checkout master
git pull origin master

# 创建新分支（命名规范：类型/简短描述）
git checkout -b feat/kvstore-add-storage-type
```

**分支命名建议：**

| 类型 | 前缀 | 示例 |
|------|------|------|
| 新功能 | `feat/` | `feat/kvstore-add-storage-type` |
| Bug 修复 | `fix/` | `fix/redis-engine-version-nil` |
| 文档 | `docs/` | `docs/update-kvstore-readme` |

---

### 第 5 步：修改代码

在本地编辑器中完成代码修改，例如本次修改了：

```
alicloud/resource_alicloud_kvstore_instance.go
```

涉及三处改动：
1. Schema 中新增字段定义
2. Create 函数中透传参数
3. Read 函数中回填字段值

---

### 第 6 步：提交代码

```bash
# 将修改的文件加入暂存区
git add alicloud/resource_alicloud_kvstore_instance.go

# 提交，message 格式：类型(范围): 简短描述
git commit -m "feat(kvstore): add storage_type support for Redis 7.0 cloud-disk instances"
```

**commit message 规范：**

```
feat(kvstore):     新增功能
fix(redis):        修复 bug
docs(kvstore):     文档更新
refactor(vpc):     代码重构
```

---

### 第 7 步：推送到你的 Fork 仓库

```bash
git push fork feat/kvstore-add-storage-type
```

> 首次推送会弹出浏览器要求登录 GitHub 授权，完成授权后自动继续。

---

### 第 8 步：在 GitHub 创建 PR

打开以下链接（替换为你的用户名和分支名）：

```
https://github.com/aliyun/terraform-provider-alicloud/compare/master...你的用户名:terraform-provider-alicloud:你的分支名
```

本次示例：
```
https://github.com/aliyun/terraform-provider-alicloud/compare/master...cxx-12333:terraform-provider-alicloud:feat/kvstore-add-storage-type
```

填写 PR 信息：
- **Title**：和 commit message 保持一致
- **Description**：说明问题背景、改了什么、怎么用

点击 **Create pull request** 完成提交。

---

## 三、PR Description 模板

```markdown
## Summary
（一句话说明这个 PR 做了什么）

## Problem
（描述当前存在的问题，附上报错信息）

## Root Cause
（根本原因分析）

## Changes
（列出改动点）
- Schema: 新增 xxx 字段
- Create: 透传 xxx 参数
- Read: 回填 xxx 字段

## Usage
（给出使用示例代码）
```hcl
resource "alicloud_xxx" "example" {
  new_field = "value"
}
```
```

---

## 四、常见问题

### Q：push 时报错 `fatal: not a git repository`
**原因**：当前终端工作目录不在仓库根目录下。  
**解决**：使用 `git -C "仓库绝对路径" push ...` 或先 `cd` 到仓库目录。

### Q：push 后在 GitHub 找不到分支
**原因**：push 的是 `origin`（官方）而不是 `fork`（你的副本）。  
**解决**：确认命令是 `git push fork 分支名` 而不是 `git push origin 分支名`。

### Q：PR 提交后被要求修改
**方法**：在本地同一分支继续修改，再次 `git add` → `git commit` → `git push fork 分支名`，PR 会自动更新。

---

## 五、本次 PR 完整记录

| 步骤 | 命令 / 操作 |
|------|------------|
| 添加远程 | `git remote add fork https://github.com/cxx-12333/terraform-provider-alicloud.git` |
| 创建分支 | `git checkout -b feat/kvstore-add-storage-type` |
| 提交代码 | `git commit -m "feat(kvstore): add storage_type support for Redis 7.0 cloud-disk instances"` |
| 推送分支 | `git push fork feat/kvstore-add-storage-type` |
| 创建 PR | https://github.com/aliyun/terraform-provider-alicloud/compare/master...cxx-12333:terraform-provider-alicloud:feat/kvstore-add-storage-type |
