# Sonar cmd
A cmd client for sonarqube and Linux. More subcommands are on going or making PRs.

## Install
```
go install github.com/alauda/bergamot/sonarqube/cmdclient
```

## Usage
### taskmonitor subcommand
just run, it will monitor sonarqube analysis task and judge result by quality gate
```
cmdclient --host sonar_host --token sonar token
```

## Exit code
If cmdclient exit unnormally, please note the exit code, they have following meanings:
```golang
// exit code
const (
	ExitCodeNormal             int = 0
	ExitCodeProgramError       int = 131
	ExitCodeTimeout            int = 132
	ExitCodeQualityGateError   int = 133
	ExitCodeQualityGateWarn    int = 134
	ExitCodeQualityGateUnknown int = 135
)
```