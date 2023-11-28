# Telegram Notification Bot

> 这是一个简单的 Telegram 通知机器人, 主要用于通过 HTTP 发送 Telegram 通知.

## 一、安装

### 1.1、二进制安装

可直接从 Release 页面下载预编译的二进制文件, 然后直接启动即可:

```sh
./notibot --recipient 123456789 --auth-mode keyauth --access-token wzdnpOjAhqOENtJKvxsSmCIftCYCUBYY

2023-11-28 12:35:55 INFO NotiBot Starting...
2023-11-28 12:35:55 INFO NotiBot Auth Mode: keyauth
2023-11-28 12:35:55 INFO NotiBot API User: noti
2023-11-28 12:35:55 INFO NotiBot API Password: BCgySJqEDcxTrRil
2023-11-28 12:35:55 INFO NotiBot API Access Token: BHrNhJBbOtaVVnqRhWKJzLBsUYyGJhsr
2023-11-28 12:35:55 INFO NotiBot Telegram Recipient: 123456789
```

### 1.2、Docker 安装

您可以直接使用预编译的 Docker 镜像来启动 NotiBot, 以下为 Docker Compose 样例:

```yaml
version: '3.9'
services:
  notibot:
    image: mritd/notibot
    container_name: notibot
    restart: unless-stopped
    environment:
      - TZ=Asia/Shanghai
      - NOTI_AUTH_MODE=keyauth
      - NOTI_AUTH_TOKEN=BHrNhJBbOtaVVnqRhWKJzLBsUYyGJhsr
      - NOTI_TELEGRAM_TOKEN=123456789:AFDKsdfs-kjabnjfbsdfSfsdfsSgut
      - NOTI_TELEGRAM_RECIPIENT=-345647567
    volumes:
      - /etc/timezone:/etc/timezone
```

## 二、配置参数

> NotiBot 接受两种配置方式: 命令行参数/环境变量

### 2.1、命令行参数

```sh
~ ❯❯❯ ./notibot --help
Telegram Notification Bot

Usage:
  notibot [flags]

Flags:
  -t, --access-token string   Server API AccessToken
  -m, --auth-mode string      Server API Auth Mode(basicauth/keyauth/none)
  -a, --bot-api string        Telegram API Address (default "https://api.telegram.org")
  -s, --bot-token string      Telegram Bot Token
  -h, --help                  help for notibot
  -l, --listen string         Server Listen Address
  -p, --password string       Server API Basic Auth Password
  -r, --recipient string      Telegram Message Recipient
  -u, --username string       Server API Basic Auth User
```

- `--listen`: 配置 NotiBot 监听地址, 默认为 `0.0.0.0:8080`
- `--bot-api`: 配置 Telegram API 地址, 默认为 `https://api.telegram.org`(可通过国外 VPS 反向代理然后使用自己的私有地址)
- `--bot-token`: Telegram Bot 的 API Token, 需要自行通过 BotFather 创建机器人获取
- `--recipient`: Telegram 接收此通知的用户/群组 ID, 支持多个 ID 以及混合推送, 多个 ID 请使用英文逗号分割
- `--auth-mode`: NotiBot 的认证模式, 为防止未授权用户滥用通知, NotiBot 默认采用 keyauth 模式认证(可选值: `basicauth`、`keyauth`、`none`)
- `--access-token`: 当使用 `keyauth` 模式时, 通过此参数设置访问 Token, 不写每次随机生成
- `--username`: 当使用 `basicauth` 模式时, 通过此参数设置访问的用户名, 不写默认为 `noti`
- `--password`: 当使用 `basicauth` 模式时, 通过此参数设置访问用户的密码, 不写默认随机生成

### 2.2、环境变量

为了方便容器化使用, NotiBot 也会从环境变量中获取配置, 环境变量与命令行参数一一对应, 对应关系如下:

- `--listen`: `NOTI_LISTEN_ADDR`
- `--bot-api`: `NOTI_TELEGRAM_API`
- `--bot-token`: `NOTI_TELEGRAM_TOKEN`
- `--recipient`: `NOTI_TELEGRAM_RECIPIENT`
- `--auth-mode`: `NOTI_AUTH_MODE`
- `--access-token`: `NOTI_AUTH_TOKEN`
- `--username`: `NOTI_AUTH_USERNAME`
- `--password`: `NOTI_AUTH_PASSWORD`

## 三、如何使用

### 3.1、请求方式

当 NotiBot 启动成功后, 所有推送发送全部采用 HTTP POST FORM 发送, 根据认证模式不同其请求样例如下:

#### Auth Mode - None

```sh
curl -sSL -XPOST https://example.com/message -d "message=测试消息"
```

#### Auth Mode - KeyAuth

```sh
curl -sSL -XPOST -H 'Authorization: Bearer RLGE9ydOlKzdf2ECALDW2cHQwQUBbGOR' https://example.com/message -d "message=测试消息"
```

#### Auth Mode - BasicAuth

```sh
curl -sSL -XPOST -u 'username:password' https://example.com/message -d "message=测试消息"
```

### 3.2、通知类型

#### 文本消息 `/message`

```sh
curl -sSL -XPOST https://example.com/message -d "message=测试消息" -d "markdown=true"
```

可选参数:

- `markdown`: Markdown 解析, 如果为 `true` 则表示 Telegram 应该将消息解析为 Markdown 格式(默认 `fasle`)
- `silent`: 静默发送通知, 为 `true` 时手机客户端会显示消息角标, 但不会发出声音

#### 图片消息 `/image`

```sh
curl -sSL -XPOST https://example.com/image -F 'image=@蔡徐坤.jpg'
```

可选参数:

- `silent`: 静默发送通知, 为 `true` 时手机客户端会显示消息角标, 但不会发出声音

#### 文件消息 `/file`

```sh
curl -sSL -XPOST https://example.com/file -F 'file=@Monaco Nerd Font Complete.ttf'
```

可选参数:

- `silent`: 静默发送通知, 为 `true` 时手机客户端会显示消息角标, 但不会发出声音
