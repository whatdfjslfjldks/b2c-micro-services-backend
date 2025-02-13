# 定义 Go 相关的变量
GO = go
GO_RUN = $(GO) run

# api网关
API_GATE_WAY = ./api-gateway/main.go


# 默认目标：运行程序
all: run

# 运行程序
run:
	$(GO_RUN) $(API_GATE_WAY)

## 清理构建文件
#clean:
#	rm -f $(TARGET)
