# Sonar cmd
A cmd client for sonarqube. More subcommands are on going or making PRs.

## Install
### On OSX/Linux
```
go install github.com/alauda/bergamot/sonarqube/cmdclient
```
### On container
1. clone code to GOPATH
```
git clone https://github.com/alauda/bergamot $PATH/src/github.com/alauda/bergamot
```
2. static build binary, [click me for details](https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/)
```
CGO_ENABLED=0 go build -a -installsuffix cgo -o cmdclient github.com/alauda/bergamot/sonarqube/cmdclient
```
## Usage
### taskmonitor subcommand
This command aims to monitor analysis task and project status.The help information are:
```
Sonarqube analysis task monitor sub command

Usage:
  sonarclient taskmonitor [flags]

Flags:
  -h, --help         help for taskmonitor
      --i duration   check sonar analysis task status interval (default 4s)
      --id string    sonar analysis task id
      --t duration   sonar analysis timeout (default 30m0s)
      --w string     sonar scanner work directory (default "./")
      --warn         whether set analysis status 'warn' to error

Global Flags:
      --host string    sonarqube server url (default "localhost:9000")
      --token string   sonarqube api token
```
demo usage:
```
cmdclient --host sonar_host --token sonar_token --w /sonar_work/
```
### qualitygate subcommand
`qualitygate` provides `select` command.Help information are:
```
cmdclient qualitygate select --help
Select sonarqube qualitygate

Usage:
  sonarclient qualitygate select [flags]

Flags:
      --gate-id int         gate id for qualitygate
  -h, --help                help for select
      --name string         sonarqube project name
      --properties string   project settings file path (default "./sonar-project.properties")

Global Flags:
      --host string    sonarqube server url (default "localhost:9000")
      --token string   sonarqube api token
```
All flags are needed:
```
cmdclient qualitygate select --host sonar_host --token sonar_token --gate-id 2 --name sonar_project_name --properties ./sonar-project.properties
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
	ExitCodeNeedMoreInfo       int = 136
)
```