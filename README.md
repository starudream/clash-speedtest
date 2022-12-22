# Clash Speedtest

![Golang](https://img.shields.io/github/actions/workflow/status/starudream/clash-speedtest/golang.yml?label=golang&style=for-the-badge)
![Docker](https://img.shields.io/github/actions/workflow/status/starudream/clash-speedtest/docker.yml?label=docker&style=for-the-badge)
![Release](https://img.shields.io/github/v/release/starudream/clash-speedtest?include_prereleases&sort=semver&style=for-the-badge)
![License](https://img.shields.io/badge/license-Apache--2.0-blue?style=for-the-badge)

`Clash` 节点测速

## Configure

| Variable | Type         | Default               | Description                      |
|----------|--------------|-----------------------|----------------------------------|
| DEBUG    | BOOL         | FALSE                 | show debug information           |
| URL      | STRING       | http://127.0.0.1:9090 | clash external controller url    |
| SECRET   | STRING       | -                     | clash external controller secret |
| PROXY    | STRING       | http://127.0.0.1:7890 | configuration file path          |
| TIMEOUT  | STRING       | 5s                    | timeout for http request         |
| INCLUDE  | STRING ARRAY | -                     | filter nodes by include          |
| EXCLUDE  | STRING ARRAY | -                     | filter nodes by exclude          |
| RETRY    | INT          | 3                     | retry times when failed          |
| OUTPUT   | STRING       | -                     | output directory                 |

### Docker

![Version](https://img.shields.io/docker/v/starudream/clash-speedtest?sort=semver&style=for-the-badge)
![Size](https://img.shields.io/docker/image-size/starudream/clash-speedtest?sort=semver&style=for-the-badge)
![Pull](https://img.shields.io/docker/pulls/starudream/clash-speedtest?style=for-the-badge)

```bash
docker pull starudream/clash-speedtest
```

```bash
docker run --rm \
    --name clash-speedtest \
    -e DEBUG=true \
    -e URL=http://host.docker.internal:9090 \
    -e PROXY=http://host.docker.internal:7890 \
    starudream/clash-speedtest
```

## Screenshot

![screenshot](./docs/screenshot.png)

## License

[Apache License 2.0](./LICENSE)
