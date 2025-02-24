# SMPP Application REST Server

[English](#english) | [中文](#chinese)

<a name="english"></a>
## English

### Overview
SMPP Application REST Server is a high-performance SMPP client that supports transmitter, receiver, and transceiver modes. It provides both REST API and Web GUI interfaces for configuration management and operation control.

### Features
- Multiple SMPP connection types support (transmitter/receiver/transceiver)
- Configurable message content (random/pre-defined/mixed mode)
- Flexible destination address generation
- Real-time metrics monitoring
- Web-based configuration management
- REST API for automation

### Quick Start
1. Build the application:
```bash
./build.sh
```

2. Configure the application:
- Copy `config/smpp-app.yaml.example` to `config/smpp-app.yaml`
- Modify the configuration according to your needs

3. Start the server:
```bash
./rest-server -c config/smpp-app.yaml
```

### Web Interface
Access the Web GUI at `http://<server-address>:8081`

The Web interface provides:
- Configuration viewing and editing
- Start/Stop message sending
- Real-time operation feedback

### REST API Endpoints
1. Start Message Loop
```
POST /api/startloop
```

2. Stop Message Loop
```
POST /api/stoploop
```

3. Get Configuration
```
GET /api/config
```

4. Update Configuration
```
POST /api/config
Content-Type: application/yaml
```

### Configuration
Key configuration items in `smpp-app.yaml`:
```yaml
service:
  smpp:
    - server:
        addr: "smpp-server-address"
        port: 5588
        user: "username"
        password: "password"
      client:
        bind-type: "transmitter"  # transmitter/receiver/transceiver
        conn-num: 1
      message:
        send:
          content-mode: "mixed"   # random/pre-defined/mixed
          text-file: "data/text.txt"
          url-file: "data/url.txt"
          tps: 100
  rest:
    addr: "0.0.0.0"
    port: 8081
  log:
    level: "debug"
```

### Metrics
The application provides real-time metrics:
- ao: Number of messages sent
- ao failure: Number of failed messages
- at: Number of messages received
- at failure: Number of failed receives

Metrics are updated every 5 seconds.

---

<a name="chinese"></a>
## 中文

### 概述
SMPP应用REST服务器是一个高性能的SMPP客户端，支持发送器、接收器和收发器模式。它提供REST API和Web图形界面用于配置管理和操作控制。

### 功能特点
- 支持多种SMPP连接类型（发送器/接收器/收发器）
- 可配置的消息内容（随机/预定义/混合模式）
- 灵活的目标地址生成
- 实时指标监控
- 基于Web的配置管理
- 支持自动化的REST API

### 快速开始
1. 构建应用：
```bash
./build.sh
```

2. 配置应用：
- 将 `config/smpp-app.yaml.example` 复制为 `config/smpp-app.yaml`
- 根据需要修改配置

3. 启动服务器：
```bash
./rest-server -c config/smpp-app.yaml
```

### Web界面
访问Web界面：`http://<服务器地址>:8081`

Web界面提供：
- 配置查看和编辑
- 启动/停止消息发送
- 实时操作反馈

### REST API接口
1. 启动消息循环
```
POST /api/startloop
```

2. 停止消息循环
```
POST /api/stoploop
```

3. 获取配置
```
GET /api/config
```

4. 更新配置
```
POST /api/config
Content-Type: application/yaml
```

### 配置说明
`smpp-app.yaml` 中的主要配置项：
```yaml
service:
  smpp:
    - server:
        addr: "smpp服务器地址"
        port: 5588
        user: "用户名"
        password: "密码"
      client:
        bind-type: "transmitter"  # transmitter(发送器)/receiver(接收器)/transceiver(收发器)
        conn-num: 1
      message:
        send:
          content-mode: "mixed"   # random(随机)/pre-defined(预定义)/mixed(混合)
          text-file: "data/text.txt"
          url-file: "data/url.txt"
          tps: 100
  rest:
    addr: "0.0.0.0"
    port: 8081
  log:
    level: "debug"
```

### 监控指标
应用提供实时监控指标：
- ao：已发送的消息数量
- ao failure：发送失败的消息数量
- at：已接收的消息数量
- at failure：接收失败的消息数量

指标每5秒更新一次。
