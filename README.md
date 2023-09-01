# kratos_consul_bug_report

本工程稳定复现了使用 consul 注册中心时重复打印日志的问题。代码取自 [go-kratos/examples](https://github.com/go-kratos/examples)。

## 工程结构

```sh
.
├── README.md
├── kratos_2_2_1
│   ├── client       # registry/consul/client
│   ├── compose.yaml
│   ├── go.mod       # go 1.16 kratos v2.2.1
│   ├── go.sum
│   ├── helloworld   # helloworld/helloworld
│   ├── main.go      # 复现脚本
│   └── server       # registry/consul/server
└── kratos_2_7_0
    ├── client
    ├── compose.yaml
    ├── go.mod       # go 1.21.0 kratos v2.7.0
    ├── go.sum
    ├── helloworld   # .proto 文件一致，其他文件为重新生成
    ├── main.go      # 复现脚本
    └── server
```

## 运行步骤

```sh
# 在 kratos v2.2.1 环境中执行复现脚本
cd kratos_2_2_1; go run .

# 在 kratos v2.7.0 环境中也执行一遍，对比日志异同。
cd ../kratos_2_7_0; go run .
```

## 运行结果

* v2.2.1 版本。日志 "agent.http: Request cancelled..." 只在 client 端退出时打印一次。

```sh
...
consul  | 2023-09-01T06:04:34.842Z [INFO]  agent: Synced node info
consul  | 2023-09-01T06:04:34.846Z [INFO]  agent: Synced service: service=6c42f6d0-488d-11ee-b985-fa40210d196d
consul  | 2023-09-01T06:04:35.858Z [INFO]  agent: Synced check: check=service:6c42f6d0-488d-11ee-b985-fa40210d196d
consul  | 2023-09-01T06:04:39.743Z [INFO]  agent: Synced check: check=service:6c42f6d0-488d-11ee-b985-fa40210d196d:2
consul  | 2023-09-01T06:04:41.444Z [INFO]  agent: Synced check: check=service:6c42f6d0-488d-11ee-b985-fa40210d196d:1
consul  | 2023-09-01T06:05:03.782Z [INFO]  agent.http: Request cancelled: method=GET url="/v1/health/service/helloworld?index=21&passing=1&wait=55000ms" from=192.168.160.1:56464 error="context canceled"
...
```

* v2.7.0 版本。日志 "agent.http: Request cancelled..." 一直在打印。

```sh
...
consul  | 2023-09-01T06:05:16.235Z [INFO]  agent: Synced node info
consul  | 2023-09-01T06:05:16.238Z [INFO]  agent: Synced service: service=83143d06-488d-11ee-b46c-fa40210d196d
consul  | 2023-09-01T06:05:16.868Z [INFO]  agent: Synced check: check=service:83143d06-488d-11ee-b46c-fa40210d196d:1
consul  | 2023-09-01T06:05:16.871Z [INFO]  agent: Synced check: check=service:83143d06-488d-11ee-b46c-fa40210d196d:2
consul  | 2023-09-01T06:05:16.956Z [ERROR] agent.server.autopilot: Failed to reconcile current state with the desired state
consul  | 2023-09-01T06:05:17.246Z [INFO]  agent: Synced check: check=service:83143d06-488d-11ee-b46c-fa40210d196d
consul  | 2023-09-01T06:05:29.198Z [INFO]  agent.http: Request cancelled: method=GET url="/v1/health/service/helloworld?index=21&passing=1&wait=55000ms" from=192.168.176.1:59556 error="context canceled"
consul  | 2023-09-01T06:05:40.199Z [INFO]  agent.http: Request cancelled: method=GET url="/v1/health/service/helloworld?index=21&passing=1&wait=55000ms" from=192.168.176.1:44994 error="context canceled"
consul  | 2023-09-01T06:05:42.043Z [INFO]  agent.http: Request cancelled: method=GET url="/v1/health/service/helloworld?index=21&passing=1&wait=55000ms" from=192.168.176.1:32978 error="context canceled"
...
```

可以根据脚本描述手动执行命令复现该过程。
