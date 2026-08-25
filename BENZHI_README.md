# 门店巡店整改系统

标准构建、运行和测试命令见下方。

```bash
go build ./...
go run ./cmd/server
go test -count=1 ./...
cd web && npm install && npm run build
```

架构构建：`./build_benzhi_docker.sh inspection linux/arm64` 和 `./build_benzhi_docker.sh inspection linux/amd64`。
