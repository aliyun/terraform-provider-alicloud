# jarvis

可被 Claude 完全接替日常工作的自包含起点母版。

## 起步

1. `git clone` 本仓库
2. `cp bootstrap/.env.example bootstrap/.env` 并填写 `GH_TOKEN` 和阿里云密钥（a1 用容器内凭证）
3. `bash bootstrap/install.sh` 安装依赖
4. `bash bootstrap/verify.sh` 验证环境 —— 全绿才干活，任一 FAIL 整体退非零
5. 读 `CLAUDE.md` 开局

## 结构

- **CLAUDE.md** 自举入口
- **autonomy.md** 决策权，硬门 = 正式发布
- **loops/** 工作流
- **skills/** vendored 技能
- **bootstrap/** 装配与验证
- **runs/** 审计
- **escalation/** 停队

## 自治边界

全链跑到预发/CR，正式发布前停；低置信 → 起草不发出入 `escalation/`。
