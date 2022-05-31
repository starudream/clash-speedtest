# Clash Speedtest

![Golang](https://img.shields.io/github/workflow/status/starudream/clash-speedtest/Golang/master?label=Golang&style=for-the-badge)
![Docker](https://img.shields.io/github/workflow/status/starudream/clash-speedtest/Docker/master?label=Docker&style=for-the-badge)
![Release](https://img.shields.io/github/v/release/starudream/clash-speedtest?include_prereleases&style=for-the-badge)
![License](https://img.shields.io/badge/License-Apache%20License%202.0-blue?style=for-the-badge)

`Clash` 节点测速

## Usage

```
Usage:
  clash-speedtest [flags]

Flags:
      --debug              (env: SCS_DEBUG) show debug information
      --exclude strings    (env: SCS_EXCLUDE) filter nodes by exclude
  -h, --help               help for clash-speedtest
      --include strings    (env: SCS_INCLUDE) filter nodes by include
      --output string      (env: SCS_OUTPUT) output directory
      --proxy string       (env: SCS_PROXY) clash http proxy url (default "http://127.0.0.1:7890")
      --retry int          (env: SCS_RETRY) retry times when failed (default 3)
      --secret string      (env: SCS_SECRET) clash external controller secret
      --timeout duration   (env: SCS_TIMEOUT) timeout for http request (default 5s)
      --url string         (env: SCS_URL) clash external controller url (default "http://127.0.0.1:9090")
  -v, --version            version for clash-speedtest
```

### Docker

![Version](https://img.shields.io/docker/v/starudream/clash-speedtest?style=for-the-badge)
![Size](https://img.shields.io/docker/image-size/starudream/clash-speedtest/latest?style=for-the-badge)
![Pull](https://img.shields.io/docker/pulls/starudream/clash-speedtest?style=for-the-badge)

```bash
docker pull starudream/clash-speedtest
```

```bash
docker run --rm \
    --name clash-speedtest \
    -e SCS_DEBUG=true \
    -e SCS_URL=http://host.docker.internal:9090 \
    -e SCS_PROXY=http://host.docker.internal:7890 \
    starudream/clash-speedtest
```

## Screenshot

![screenshot](./docs/screenshot.png)

## License

[Apache License 2.0](./LICENSE)
