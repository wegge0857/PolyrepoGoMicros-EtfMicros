# Developed by @Viggo Van
[Email](mailto:wayne3van@gmail.com)

### 执行命令
```bash
go get github.com/wegge0857/PolyrepoGoMicros-ApiLink
go mod tidy
```

### 添加分布式事务管理器
```bash
go get github.com/dtm-labs/client
```
### 如barrier表在当前微服务数据库
### 在对应的微服务data层
### import "github.com/dtm-labs/client/dtmcli"
### dtmcli.SetBarrierTableName("barrier")