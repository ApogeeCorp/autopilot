# autopilot
Cluster monitoring and automated recommendations

# prerequisites

1. go 1.9+
2. 
```
> go run ./cmd/autopilot/main.go collect -p prometheus -u http://70.0.69.141:9090/api/v1 -a 'query={cluster="greatdane-1914e166dc7"};start=2018-10-17T00:00:00.0Z;end=2018-10-18T00:00:00.0Z;step=15m'
```