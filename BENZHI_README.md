# 门店巡店整改系统

基于 Go 实现的命令行项目，一款业务数据管理工具，围绕 Store、Inspection、Finding 等业务对象完成创建、校验、更新、查询与结果记录。

```bash
go build ./...
go run ./cmd/server
go test -count=1 ./...
cd web && npm install && npm run build
```

架构构建：`./build_benzhi_docker.sh inspection linux/arm64` 和 `./build_benzhi_docker.sh inspection linux/amd64`。
