# lorry
lorry是原KubeBlocks中提供命令执行通道的服务，提供了多种wellknown数据库引擎的action实现，比如apecloud-MySQL、MySQL、Redis、PostgreSQL、MongoDB等等。对于lorry已经支持的数据库引擎，在做addon接入时，可以声明lorry作为engines plugin，快速接入KubeBlocks。

## 声明engines plugin
KubeBlocks本身对engine plugin的形态不做限制，可以以sidecarCar或container中守护进程或其它形态运行。lorry目前默认采用localhost方式与DB进程通信，所以建议使用sidecar方式部署lorry，部署模版：
```
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: mysql
    image: apecloud-registry.cn-zhangjiakou.cr.aliyuncs.com/apecloud/apecloud-mysql-server:8.0.30
    ...
  - name: lorry
    image: apecloud-registry.cn-zhangjiakou.cr.aliyuncs.com/apecloud/lorry:1.0.0
    command:
    - lorry
    - --port
    - "3501"
    - --grpcport
    - "50001"
    env:
    - name: MYSQL_ROOT_USER
      valueFrom:
        secretKeyRef:
          key: username
          name: cluster-mysql-account-root
    - name: MYSQL_ROOT_PASSWORD
      valueFrom:
        secretKeyRef:
          key: password
          name: cluster-mysql-account-root
    - name: KB_BD_TYPE
      value: mysql
```

## Action定义
在KB 1.0版本中，action支持Exec和GRPC两种定义方式，addon Provider在使用lorry时，可以通过以下方式使用lorry提供的action能力，以roleProbe为例：

### Exec
```
  lifecycleActions:
    roleProbe:
      customHandler:
        exec:
          command:
            - /bin/bash
            - -c
            - curl -X GET -H 'Content-Type: application/json' 'http://127.0.0.1:3501/v1.0/checkrole'
```

### GRPC
正在实现中，此处先给出接入的方式（待定，可讨论）
```
  lifecycleActions:
    roleProbe:
      customHandler:
        grpc:
          host: 3501
          port: 50001
```
              
