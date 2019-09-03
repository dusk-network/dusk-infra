# Testnet node monitoring
/path/to/monitor`  
## Synopsis

```sh
    $ monitor [OPTION...]
```

## Description

Runs the monitoring backend used to recurrently scan the machine parameters as well as capture any error or warning the node is generating. The data is published on a websocket in JSON format. The `monitor` is written in golang.

## Options

Two different ways of setting options are available. Passing parameters from the command line (`CLI`) or setting equivalent environment variables (`ENV`). In case both are specified `CLI` parameters have precedence.

### Generic program information

```
  -help
        this help
  -h 
        this help (shorthand)
```

### Host address of the monitoring web-based client

```
  -a string
    	http service address (shorthand) (default "localhost:8080")
  -address string
    	http service address (default "localhost:8080")
        equivalent of setting `ADDRESS` environment variable
```

### Location of the node log file 

Specify the location of the node's log file to monitor for warning and errors

```
  -l string
    	location of the node log file (shorthand) (default "/var/log/node.log")
  -logfile string
    	location of the node log file (default "/var/log/node.log")
```

### 
  -u value
    	URI of the log monitoring server (default "unix:///var/tmp/logmon.sock")(shorthand)
  -uri-logserver value
    	URI of the log monitoring server (default "unix:///var/tmp/logmon.sock")

### Latency prober address

Specify the URL of the prober used to measure latency. Normally (on testnet) this is the IP Address of the `Voucher Seeder`

```
  -s string
    	preferred voucher seeder (shorthand) (default "178.62.193.89")
  -seeder string
    	preferred voucher seeder (default "178.62.193.89")
```

#### Priviledges

The latency is measured by sending an ICMP packet to the prober and wait for a response. On a Unix system, usually this requires some heightened priviledges (root) and therefore, the latency monitoring will be activated only if the `monitor` process will be granted said priviledges.
This can be achieved in the following ways:
 - Use `setcap` to allow binary to bind to raw sockets: `setcap cap_net_raw=+ep /path/to/monitor` (*preferred*)
 - Modifying the priviledges of `ping_group_range` by running `sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"` (*discouraged*: this is system dependant)
 - Running the `monitor` as root (*discouraged*: potentially dangerous to export the ENV on superuser)

### Verbose

Set log level to `Debug`

```
  -v	start in debug mode (shorthand)
  -verbose
    	start in debug mode
``` 

## Running the monitor

### Building the monitoring

```golang
 $ go build path/to/monitor/...
```

### Running the build

```golang
 $ path/to/monitor/main
```

## System parameters

System parameters are published as json objects with the following structure:

```json
{
    "timestamp": timestamp,
    "metric": metric_name,
    "value": metric_value
}
```

 - Parameter `timestamp` is the UTC time with the following format: `yyyy-mm-ddThh:MM:ss.SSSSSS+Z" (e.g. `"2019-06-23T14:05:49.308707084+02:00"`)
 - Parameter `metric` is the ID of the parameter notified
 - Parameter `value` is the value of the parameter notified

All fields are `string`. Parameters are notified according to the following table:

| Parameter | Metric    | Value     | Description |
| --------- | ------    | -----     | ----------- |
| Disk      | disk      | 1-100     | Disk usage in percentile |
| CPU       | cpu       | 1-100     | CPU usage in percentile (aggregated on all cores) |
| Memory    | mem       | 1-100     | Memory usage in percentile |
| Latency   | latency   | seconds   | milliseconds of latency (requires priviledges) |
| Tail      | tail      | text      | output of tail process on the log file |     
| Log       | log       | text      | error notification reported as by the node |
