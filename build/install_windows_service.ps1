<#
.SYNOPSIS
  使用 sc.exe 在 Windows 上创建自启动服务（需要管理员权限）

#>
param(
  [string]$BinaryPath = "C:\\Program Files\\fly-go\\fly-go.exe",
  [string]$ServiceName = "fly-go"
)

if (-not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
  Write-Error "请以管理员身份运行此脚本"
  exit 1
}

Write-Host "Creating service $ServiceName with binary $BinaryPath"

New-Item -ItemType Directory -Force -Path (Split-Path $BinaryPath) | Out-Null
if (-not (Test-Path $BinaryPath)) {
  Write-Warning "二进制文件不存在，脚本不会拷贝；请确保 $BinaryPath 已存在或手动复制。"
}

$exists = sc.exe query $ServiceName 2>&1 | Out-String
if ($exists -notmatch "FAILED 1060") {
  Write-Host "服务 '$ServiceName' 已存在，尝试删除旧服务"
  sc.exe delete $ServiceName | Out-Null
  Start-Sleep -Seconds 1
}

sc.exe create $ServiceName binPath= "`"$BinaryPath`"" start= auto DisplayName= "$ServiceName" | Out-Null
Write-Host "服务已创建。使用 'Start-Service $ServiceName' 启动。"
