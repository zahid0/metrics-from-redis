# Metrics From Redis
A Prometheus exporter which provides metrics be reading the keys of format "{prefix}:{metricname}:{labels}" from redis. It is useful when the process doesn't serve http requests. Process can write the metrics in redis and the exporter will make it available to be scraped by prometheus.

## Usage
```
metrics-from-redis
```
Command line arguments:
* -p : string - Redis key prefix
* -c : string - Redis connection string in the format `{host}:{port}`
* -d : Integer - Redis database number
