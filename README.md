与上游的区别：

- Docker 镜像发布源已迁移到 GHCR，推荐使用 `ghcr.io/tunecc/moments:latest`。
- 仓库内的 `compose.yml` 已默认指向本 fork 镜像。
- 发布帖子支持自定义时间



默认管理员账号：

- 用户名：`admin`
- 密码：`a123456`

首次登录后请立即修改默认密码。

## Docker Compose 部署

推荐使用 Docker Compose 部署，便于持久化数据和后续升级。

### 1. 生成 JWT_KEY

`JWT_KEY` 用于签发登录状态。建议固定配置，否则随机生成后每次容器重启都可能需要重新登录。

```bash
openssl rand -hex 32
```

### 2. 创建 docker-compose.yml

```yaml
services:
  moments:
    image: ghcr.io/tunecc/moments:latest
    container_name: moments
    restart: always
    environment:
      PORT: 3000
      JWT_KEY: your_jwt_key_here
    ports:
      - 3000:3000
    volumes:
      - ./data:/app/data
```

把 `your_jwt_key_here` 替换为上一步生成的值。

### 3. 启动

```bash
docker compose up -d
```

启动后访问：

```text
http://你的服务器IP:3000
```

## Docker CLI 部署

如果不使用 Docker Compose，也可以直接运行容器。

```bash
docker run -d \
  --name moments \
  --restart always \
  -e PORT=3000 \
  -e JWT_KEY=your_jwt_key_here \
  -p 3000:3000 \
  -v /var/moments:/app/data \
  ghcr.io/tunecc/moments:latest
```

停止并删除容器：

```bash
docker stop moments
docker rm moments
```

使用 Docker Compose 升级：

```bash
docker compose pull
docker compose up -d
```

## 常用环境变量

| 变量名 | 说明 | Docker 场景建议值 |
| --- | --- | --- |
| `PORT` | 容器内监听端口 | `3000` |
| `JWT_KEY` | JWT 密钥，建议固定配置 | 使用 `openssl rand -hex 32` 生成 |
| `DB` | SQLite 数据库路径 | 默认 `/app/data/db.sqlite` |
| `UPLOAD_DIR` | 本地上传文件目录 | 默认 `/app/data/upload` |
| `CORS_ORIGIN` | 允许的跨域 Origin，多个值用英文逗号分隔 | 按需配置 |
| `LOG_LEVEL` | 日志级别 | `INFO`，排障时可用 `DEBUG` |
| `ENABLE_SWAGGER` | 是否启用 Swagger 文档 | `false` |
| `ENABLE_SQL_OUTPUT` | 是否输出 SQL 日志 | `false` |

如果使用反向代理并且前后端同域访问，通常不需要配置 `CORS_ORIGIN`。

## 许可

本项目继承上游 Moments 的许可证，详见 [LICENSE.txt](./LICENSE.txt)。
