# Lorry
Lorry是原KubeBlocks中提供命令执行通道的服务，提供了多种wellknown数据库引擎的action实现，比如apecloud-MySQL、MySQL、Redis、PostgreSQL、MongoDB等等。对于Lorry已经支持的数据库引擎，在做addon接入时，可以声明Lorry作为engines plugin，快速接入KubeBlocks。

Lorry本身提供两中运行模式：守护进程模式和临时任务模式，可以根据业务场景，自行选择适合的运行方式。
## 守护进程模式
该模式下lorry已常驻进程存在，以APIServer方式对外提供服务，可以理解为engines plugin的一种实现形态。
KubeBlocks本身对engines plugin的形态不做限制，可以以sidecarCar、container守护进程或其它形态运行。Lorry目前默认采用localhost方式与DB进程通信，所以该模式下，建议使用sidecar方式部署Lorry，部署模版：
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
    - service
    env:
    - name: KB_POD_NAME
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.name
    - name: KB_POD_UID
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.uid
    - name: KB_NAMESPACE
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.namespace
    - name: KB_SERVICE_USER
      valueFrom:
        secretKeyRef:
          key: username
          name: cluster-mysql-account-root
    - name: KB_SERVICE_PASSWORD
      valueFrom:
        secretKeyRef:
          key: password
          name: cluster-mysql-account-root
    - name: KB_CLUSTER_NAME
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.labels['app.kubernetes.io/instance']
    - name: KB_COMP_NAME
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.labels['apps.kubeblocks.io/component-name']
    - name: KB_ENGINE_TYPE
      value: mysql
```

当使用lorry，且以守护进程运行时，action定义可以使用调用lorry api的方式实现：
```
  lifecycleActions:
    roleProbe:
      exec:
        command:
          - /bin/bash
          - -c
          - curl -X GET -H 'Content-Type: application/json' 'http://127.0.0.1:3501/v1.0/getrole'
```

## 临时任务
该模式下，lorry执行完成相应任务后会立即退出，适合作为一个工具使用，可以action中直接调用。对于已支持的action列表及使用可参考docs目录下的文档。

```
  lifecycleActions:
    roleProbe:
      exec:
        image: registry.cn-hangzhou.aliyuncs.com/xuriwuyun/lorry:latest
        command:
          - lorry
          - getrole
```
              
