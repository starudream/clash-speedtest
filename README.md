# Clash-SpeedTest

![Go](https://img.shields.io/github/workflow/status/starudream/clash-speedtest/Go/master?style=for-the-badge)
![Docker](https://img.shields.io/github/workflow/status/starudream/clash-speedtest/Docker/master?style=for-the-badge)
![License](https://img.shields.io/badge/License-Apache%20License%202.0-blue?style=for-the-badge)

## Usage

![Version](https://img.shields.io/docker/v/starudream/clash-speedtest?style=for-the-badge)
![Size](https://img.shields.io/docker/image-size/starudream/clash-speedtest/latest?style=for-the-badge)
![Pull](https://img.shields.io/docker/pulls/starudream/clash-speedtest?style=for-the-badge)

```bash
docker pull starudream/clash
docker pull starudream/clash-speedtest
```

```yaml
version: "3.8"

services:

  clash:
    image: starudream/clash
    volumes:
      - /opt/docker/clash/config.yaml:/root/.config/clash/config.yaml

  clash-speedtest:
    image: starudream/clash-speedtest
    depends_on:
      - clash
    command: /clash-speedtest -url http://clash:9090 -proxy http://clash:7890
```

```bash
docker-compose up
```

## License

[Apache License 2.0](./LICENSE)
