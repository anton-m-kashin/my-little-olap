# My Little OLAP

Simple OLAP app with Go, ClickHouse and Grafana.

Provides simple API to upload bunch of metrics data in hardcoded format.

## Environment

Uses following environment variables:

- `MY_LITTLE_OLAP_DB_USER`: Clickhouse user that cat create database and write
    data
- `MY_LITTLE_OLAP_DB_DBNAME`: Clickhouse Database name that will be user to
    store OLAP data
- `MY_LITTLE_OLAP_DB_PASSWORD`: Clickhouse user password
- `GRAFANA_CH_USER`: Clickhouse user, that will be user by Grafana to access
    data
- `GRAFANA_CH_PASSWORD`: password for Grafana Clickhouse user
