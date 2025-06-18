# go-manus
<mark style="background - color: yellow;">中文文档-------[README-cn.md](https://github.com/xdl2003/go-agent/blob/main/README-cn.md)--------</mark>

A multi - agent self - planning system written in Go, adapted from [OpenManus](https://github.com/mannaandpoem/OpenManus/). It incorporates MCP (Model Context Protocol) configuration capabilities similar to Claude, allowing users to freely integrate third - party MCPs.

## Features
1. Multi - agent self - planning
2. MCP configuration capabilities
3. Easy integration of third - party MCPs

## Configuration
The `config.yaml` file is used to configure the system. Here is a basic guide on how to fill it:

### `PrimaryConfig`
- `ModelSource`: Specify the source of the model, e.g., `doubao`.
- `ModelName`: Name of the model, e.g., `chatglm_pro`.
- `ApiKey`: API key for accessing the model.

### `ExecutorConfig`
Similar to `PrimaryConfig`, used for the executor's model configuration.

### `AllMcpConfig`
- `type`: Protocol type, usually `stdio`.
- `command`: Command to start the MCP server, e.g., `npx`.
- `env`: Environment variables, such as API keys.
- `args`: Command - line arguments for starting the MCP server.

**Important**: The Agent requires `webSearch` and `FileSystem` MCP plugins to function properly. We strongly recommend using the following plugins:
- [FileSystem MCP](https://github.com/modelcontextprotocol/servers/tree/main/src/filesystem)
- [webSearch MCP](https://github.com/exa - labs/exa - mcp - server)

### Example Configuration
Refer to `config/config - example.yaml` for an example configuration.
