# Stats - A middleware and http handler to collect metrics from golang HTTP applications

## Features
 - Middleware that can be plugged into your http stack
 - HttpHandler that serves "/stats" endpoint.
 - build on top of coda-hale's metrics library. [Link]
 - Record HTTP stats across all requests
 - Record HTTP stats for each URL
 - Outputs metrics in Influxdb Line protocol, JSON, Graphite format.

## Output Influxdb example
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

 ## Output Json example
```
```

## Installation

## Usage