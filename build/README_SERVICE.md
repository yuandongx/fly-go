本目录包含用于将本项目作为系统服务安装的模板与脚本。

- `fly-go.service.template`：systemd 单元模板。将 `@@EXEC_PATH@@` 替换为可执行文件的绝对路径后复制到 `/etc/systemd/system/<name>.service`。
- `install_service.sh`：在 Linux（systemd）系统上安装并启动服务。用法：

  ./install_service.sh [path/to/binary] [service_name]

  - 默认二进制：`/usr/local/bin/fly-go`
  - 默认服务名：`fly-go`

- `uninstall_service.sh`：卸载 Linux systemd 服务并删除二进制（如果在 `/usr/local/bin`）。
- `install_windows_service.ps1`：在 Windows 上使用 `sc.exe` 创建自启动服务（需要以管理员权限运行）。用法示例（管理员 PowerShell）：

  .\install_windows_service.ps1 -BinaryPath "C:\\Program Files\\fly-go\\fly-go.exe" -ServiceName "fly-go"

注意：脚本尽量保持简单、通用；在生产环境请根据需要调整用户、日志、环境变量、安全权限等。
