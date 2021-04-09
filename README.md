# vigilant-engine

Dynamic DNS updater


## Deploy

To work it needs a secret

```
apiVersion: v1
kind: Secret
metadata:
  name: vigilant-engine
data:
  USERNAME: abcdefghi==
  PASSWORD: abcdefghi==
  DNS_RECORD: abcdefghi==
```
