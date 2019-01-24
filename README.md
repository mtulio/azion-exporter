# azion-exporter

Azion metrics exporter for Prometheus.

The [Azion](https://www.azion.com.br) Exporter [consumes data from Azion Analytics API](https://www.azion.com.br/developers/api-v2/) using internal GO library exposing it for [Prometheus](https://prometheus.io/).

## BUILD

`make build`

The binary will be created on `./bin` dir.

## ARGUMENTS

### REQUIRED

`-azion.email` : Azion's Account Email

`-azion.token` : Azion's Account Password

### OPTIONAL

None

## USAGE

Show Azion metrics from Analytics:

```bash
./bin/azion-exporter -azion.email=my@email.com -azion.password=myPass
```

> Sample output for `$ curl localhost:9801/metrics`:

```log
# HELP azion_cd_data_transferred_mb Azion Analytics Content Delivery Data Transferred in MB
# TYPE azion_cd_data_transferred_mb gauge
azion_cd_data_transferred_mb{type="missed"} 944.02355
azion_cd_data_transferred_mb{type="saved"} 1958.173775
# HELP azion_cd_requests_count Azion Analytics Content Delivery Requests Count
# TYPE azion_cd_requests_count gauge
azion_cd_requests_count{type="missed"} 21497.0
azion_cd_requests_count{type="saved"} 66719.0
# HELP azion_cd_status_code_total Azion Analytics Content Delivery Status Code 5xx Total
# TYPE azion_cd_status_code_total gauge
azion_cd_status_code_total{code="200"} 75114.0
azion_cd_status_code_total{code="204"} 0.0
azion_cd_status_code_total{code="206"} 85.0
azion_cd_status_code_total{code="2xx"} 8.0
azion_cd_status_code_total{code="301"} 1475.0
azion_cd_status_code_total{code="302"} 4.0
azion_cd_status_code_total{code="304"} 9426.0
azion_cd_status_code_total{code="3xx"} 1.0
azion_cd_status_code_total{code="400"} 27.0
azion_cd_status_code_total{code="403"} 145.0
azion_cd_status_code_total{code="404"} 1881.0
azion_cd_status_code_total{code="4xx"} 23.0
azion_cd_status_code_total{code="500"} 0.0
azion_cd_status_code_total{code="503"} 1.0
azion_cd_status_code_total{code="5xx"} 20.0
# HELP azion_scrape_collector_duration_seconds azion_exporter: Duration of a collector scrape.
# TYPE azion_scrape_collector_duration_seconds gauge
azion_scrape_collector_duration_seconds{collector="analytics"} 6.8305e-05
# HELP azion_scrape_collector_success azion_exporter: Whether a collector succeeded.
# TYPE azion_scrape_collector_success gauge
azion_scrape_collector_success{collector="analytics"} 1.0

```

## USAGE IN DOCKER

Show Azion metrics running in docker

```bash
docker run -p 9801:9801 -id mtulio/azion-exporter:latest \
    -azion.email=my@email.com -azion.password=myPass
```

* Docker Compose definition

```YAML
# Azion exporter - https://github.com/mtulio/azion-exporter
  azion:
    image: mtulio/azion-exporter:v0.1.1
    command:
        - -azion.email=my@email.com
        - -azion.password=myPass
    ports:
      - 9801:9801
    networks:
      - net
    deploy:
      resources:
        limits:
          cpus: "0.2"
          memory: 256M
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9801/"]
      interval: 5s
      timeout: 2s
      retries: 3
```

## CONTRIBUTOR

You are welcome. =)

Some issues that we need you help:

* Writing tests
* Improving the Documentation
* Split the API Library to an external repo
* Improve CLI options to enable metrics dynamicaly
* Open an issue
> [...]

Contribution guidelines:

* See [License](./LICENSE)
* Fork me
* Open an PR with enhancements, bugfixes, etc
* Request review from collavorators
