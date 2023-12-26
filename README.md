# Clash Speedtest

![golang](https://img.shields.io/github/actions/workflow/status/starudream/clash-speedtest/golang.yml?style=for-the-badge&logo=github&label=golang)
![release](https://img.shields.io/github/v/release/starudream/clash-speedtest?style=for-the-badge)
![license](https://img.shields.io/github/license/starudream/clash-speedtest?style=for-the-badge)

## Usage

```
> clash-speedtest -h
Usage:
  clash-speedtest [flags]

Flags:
      --clash-addr string     clash external controller address (default "http://127.0.0.1:9090")
      --clash-proxy string    clash proxy url, http or socks5
      --clash-secret string   clash external controller secret
  -y, --confirm               confirm to start speedtest
  -d, --download string       download type: cloudflare, speedtest, fast (default "cloudflare")
  -e, --excludes strings      exclude proxy names, after filter by include
  -h, --help                  help for clash-speedtest
  -i, --includes strings      include proxy names
  -o, --output string         output file path (default "output")
  -r, --size uint16           download size for each thread (default 20)
  -t, --threads uint16        download threads for each type, each thread will took 2MB traffic (default 3)
  -v, --version               version for clash-speedtest
```

- `--clash-addr`
  > 默认使用 http://127.0.0.1:9090

- `--clash-proxy`
  > 通过 `clash-addr` 的 `ip` + 接口 `/configs` 获取 `Port(http)` 或 `SocksPort(socks5)` 或 `MixedPort(http)`  
  > 如果配置了 `authentication` 需要手动指定 `--clash-secret http://user:pass@127.0.0.1:7890`

- `--download` `-d`
  > 选择测试下载速度的方式，目前支持 `cloudflare` `speedtest` `fast`  
  > - `cloudflare` https://speed.cloudflare.com
  > - `speedtest` https://www.speedtest.net
  > - `fast` https://fast.com

- `--size` `-r`
  > 指定下载方式每个线程的下载大小
  > - `cloudflare` `fast` 单位为 MB
  > - `speedtest` 参考 https://github.com/starudream/clash-speedtest/blob/v3/api/speedtest/download.go

- `--threads` `-t`
  > 指定下载方式的线程数，会同时下载多个文件，每个线程会消耗大约 2MB 流量

- `includes` `-i`
  > 包含的代理名称，满足其一即包含  
  > 如果有多个，可以使用多个 `-i` 参数，或者使用 `,` 分隔，如果有 `|`，需要使用 `\|` 转义或使用引号包裹

- `--excludes` `-e`
  > 排除的代理名称，满足其一即排除

- `output` `-o`
  > 输出文件路径，默认为当前二进制文件同目录下 `output` 文件夹

### Example

```shell
clash-speedtest -i 'IEPL|2' -e Premium
```

![test](./docs/test.png)

- `CONN` 连接时间
- `DOWN` 下载速度 `total_size / response_time`

![TraceInfo](./docs/vearne.cc-39953-trace-info.png)

## Thanks

- [cloudflare speed](https://speed.cloudflare.com)
- [speedtest.net](https://www.speedtest.net)
- [fast.com](https://fast.com)
- [go-resty trace info](https://vearne.cc/archives/39953)

## [License](./LICENSE)
