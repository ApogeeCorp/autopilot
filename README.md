# autopilot
Cluster monitoring and automated recommendations

# prerequisites

1. go 1.9+
2. 
```
> go install github.com/libopenstorage/autopilot/cmd/autopilot
> autopilot -f /etc/config.json -d /var/run/autopilot
```

## Example API Call to autopilot
To try a the Rules API use the following cURL command.

```
curl -X GET \
  'http://localhost:9000/api/v1/providers/portworx/recommend?rules=PromPortworxStorageUsageCritical,PromPortworxVolumeUsageCritical' \
  -H 'Accept: application/json' \
  -H 'Authorization: Basic YXNkOmFzZGY=' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache'
```
