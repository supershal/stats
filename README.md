# Stats - Go middleware to collect HTTP stats

## Features
 - Stats Middleware that can be plugged into your http stack
 - Collects HTTP request and repsponse stats with default implementation.
 - Support to provide your own implementation of stats collection to suit your application needs
 - HttpHandler that serves "/stats" endpoint.
 - build on top of coda-hale's metrics library - http://github.com/codahale/metrics
 - Outputs metrics in Influxdb Line protocol, JSON, Graphite format. (Currently only supports influxdb line protocol: https://github.com/influxdata/influxdb/blob/master/tsdb/README.md)
 - Examples to demonstrate application and request level metrics collection.
 - `github.com/supershal/stats/metrics` package can be used to collect metrics for non-http apps. For example, Database stats can be collected using `metrics` package.

Installation
------------
1. Install Go Version >=1.5.1
2. Make sure you have installed git, hg etc. 
3. Your workspace is setup with GOPATH. http://golang.org/doc/code.html
4. Download it:
	```
    go get github.com/supershal/stats
    ```

## GoDoc

## Usage
HTTP middleware provided to collect stats. Currently it supports [alice](https://github.com/justinas/alice) and Negroni(https://github.com/codegangsta/negroni) compatible middlewares. However, its very easy to create your own middleware function for your framework.(Checkout `middleware.go`)
 
1. Add HttpStats middleware in your HTTP app
   Following example is using 'alice' middleware chaining. Checkout `/example` for more details.
	```
		func middlewares() alice.Chain {
			host, _ := os.Hostname()
			tags := map[string]string{
			"host": host,
		}
		s := stats.NewHTTPStats(tags)
		return alice.New(s.HTTPStatsHandler)
		}
	```
2. Serve metrics on separate HTTP server. 
	``` 
		stats.ServeMetrics(5555, "/metrics") 
	```
A metric collector agent [collectd](https://github.com/collectd/collectd) or [telegraf](https://github.com/influxdata/telegraf>) can invoke `localhost:5555/metrics` periodically and send metrics back to TSDB (influxdb or graphite).

## Output Influxdb example
https://github.com/influxdata/influxdb/blob/master/tsdb/README.md
```
http_request,host=localhost,foo=bar GET=10
http_request,host=localhost,foo=bar POST=5 
http_response,host=localhost,foo=bar 200=5
http_response,host=localhost,foo=bar 503=2
http_response,host=localhost,foo=bar 403=3
http_response,host=localhost,foo=bar total=10
http_response,host=localhost,foo=bar latency.P50=10
http_response,host=localhost,foo=bar latency.P75=15
http_response,host=localhost,foo=bar latency.P90=20
http_response,host=localhost,foo=bar latency.P95=25
http_response,host=localhost,foo=bar latency.P99=40
http_response,host=localhost,foo=bar latency.P999=50
```

Create PR or new [issue](https://github.com/supershal/stats/issues) for any feature request or bugs.