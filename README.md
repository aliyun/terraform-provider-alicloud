# jarvis

可被 Claude 完全接替日常工作的自包含起点母版。

## 起步（人类只需 3 步）

1. `git clone` 本仓库
2. `cp bootstrap/.env.example bootstrap/.env` 填好 `GH_TOKEN`/阿里云密钥（a1 用容器内凭证）
3. 进目录启动 `claude`

之后全是和 Claude 对话。`CLAUDE.md` 会让 Claude 自己跑 `install.sh`/`verify.sh`、全绿后开工——你不用手敲。环境缺啥它会告诉你。

## 结构

- **CLAUDE.md** 自举入口
- **autonomy.md** 决策权，硬门 = 正式发布
- **loops/** 工作流
- **.claude/skills/** vendored 技能（标准路径，clone 即被加载）
- **bootstrap/** 装配与验证
- **runs/** 审计
- **escalation/** 停队

## 自治边界

全链跑到预发/CR，正式发布前停；低置信 → 起草不发出入 `escalation/`。
