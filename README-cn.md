# go-manus

基于 [OpenManus](https://github.com/mannaandpoem/OpenManus/) 改造的Go语言多智能体自规划系统，引入了类似Claude的MCP（模型上下文协议）配置能力，可自由接入第三方MCP。

## 功能特性
1. 多智能体自规划
2. MCP配置能力
3. 便捷接入第三方MCP

## 配置说明
`config.yaml` 文件用于配置系统，以下是基本填写说明：

### `PrimaryConfig`
- `ModelSource`: 指定模型来源，例如 `doubao`。
- `ModelName`: 模型名称，例如 `chatglm_pro`。
- `ApiKey`: 访问模型所需的API密钥。

### `ExecutorConfig`
与 `PrimaryConfig` 类似，用于执行器的模型配置。

### `AllMcpConfig`
- `type`: 协议类型，通常为 `stdio`。
- `command`: 启动MCP服务的命令，例如 `npx`。
- `env`: 环境变量，如API密钥。
- `args`: 启动MCP服务的命令行参数。

**重要提示**：Agent正常工作需要 `webSearch` 和 `FileSystem` MCP插件，强烈推荐使用以下插件：
- [FileSystem MCP](https://github.com/modelcontextprotocol/servers/tree/main/src/filesystem)
- [webSearch MCP](https://github.com/exa-labs/exa-mcp-server)

### 配置示例
示例配置请参考 `config/config-example.yaml`。

## 英文文档
英文文档请查看 [README.md](README.md)