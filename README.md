# autopilot

[![Build Status](https://travis-ci.org/libopenstorage/autopilot.svg?branch=master)](https://travis-ci.org/libopenstorage/autopilot)

Cluster monitoring and automated recommendations

## Building

## Vendoring

This repo uses [go dep](https://golang.github.io/dep/) for vendoring. Following make rule will update the vendor.

```shell
make vendor
```

## Code generation

Once you make changes to the CRD, use the following make rule to update the generated code. When committing changes, keep the generated code separate.

```shell
make codegen
go run ./cmd/autopilot/*.go  -f ./etc/config-example.yaml policy test ./etc/policy-example.yaml
```


## Running

autopilot expects to run in cluster, of running standalone be sure to set `KUBERNETES_CONFIG` to your cluster configuration file path.