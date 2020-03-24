## Prepare tcpdump

1. Enter container

```
kubectl exec -it zeebe-0 bash
```

1. Install tcpdump

```
apt-get update
apt-get install -y tcpdump
```

1. Run tcpdump

```
tcpdump -i eth0 -s0 -w data/tcpdump.data
```

1. Analys missing requests after benchmark run

```
diff <(seq 0 2000) <(grep -aoe 'requestId":[0-9]\+' data/tcpdump.data | cut -d : -f 2 | sort -un)
```

Note: replace `2000` with the number of requests send

## bulk-create-instance.go

1. Setup Cloud Env

```
export ZEEBE_ADDRESS='...'
export ZEEBE_CLIENT_ID='...'
export ZEEBE_CLIENT_SECRET='...'
export ZEEBE_AUTHORIZATION_SERVER_URL='https://login.cloud.ultrawombat.com/oauth/token'
```

1. Setup Benchmark Parameters

```
export ZEEBE_INSTANCES=2000
```

1. Run Benchmark

```
go run bulk-create-instance.go
```
