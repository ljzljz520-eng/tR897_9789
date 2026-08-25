# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
ok  	callcentertroubleshooter/cmd/troubleshootd	0.018s
ok  	callcentertroubleshooter/internal/app	0.032s
ok  	callcentertroubleshooter/internal/cli	0.002s
ok  	callcentertroubleshooter/internal/domain	0.002s
ok  	callcentertroubleshooter/internal/fixtures	0.002s
ok  	callcentertroubleshooter/internal/integration	0.033s
ok  	callcentertroubleshooter/internal/report	0.002s
ok  	callcentertroubleshooter/internal/security	0.002s
ok  	callcentertroubleshooter/internal/store	0.030s
--- FAIL: TestTroubleshootScopesFields (0.00s)
    troubleshoot_test.go:21: unexpected fields in scoped result: map[home_directory:\fileserver\li.wei payroll_group:CC-NORTH security_groups:[AD-STAFF AD-CC]]
FAIL
FAIL	callcentertroubleshooter/internal/troubleshoot	0.002s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/troubleshootd): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/troubleshootd): exit `0`
