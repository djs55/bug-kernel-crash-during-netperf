
This follows the repro instructions from https://github.com/docker/for-mac/issues/5280

First install `netperf`:
```
brew install netperf
```

Run a server on the host with
```
netserver -D
```

Run Docker Desktop and switch to the VIrtualization Framework.

Build and run the test harness with:
```
go build
./bug-kernel-crash-during-netperf
```

Every 10 minutes it will restart the Desktop process. This is because I observed the crash seemed to happen
more often soon after the VM started. On one occasion it ran for about 2 days without issue, then I rebooted
the VM and it crashed.
