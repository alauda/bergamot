# Alauda SonarQube Go Client
This is a simple sonarqube client for go, it just supports SonarQube 6.0 and SonarQube 6.4. bellow APIs are supported:
- SystemStatus: get api/system/status
- CreateProject: post api/projects/create
- GetSettings: get api/settings/values
- SetSettings: post api/settings/set
- ListQualityGates: get api/qualitygates/list
- SelectQualityGates: post api/qualitygates/select

## Requirements
- golang 1.8+
- make

## Example
- get system status info
```go
package main

import(
    "fmt"
    "log"

    "github.com/bergamot/alauda/sonarqube"
)

func main() {
    endpoint := "your sonar endpoit" // like http://sonar.alauda.cn
    token := "your sonar token"
    sonar := sonarqube.NewSonarQubeArgs(endpoint, token)

    status, err := sonar.SystemStatus()
    if err != nil {
        log.Fatalf("method SystemStatus error: %v", err)
    }
    fmt.Println("system status is ", status)
}
```

- create project
```go
package main

import(
    "fmt"
    "log"

    "github.com/bergamot/alauda/sonarqube"
)

func main() {
    endpoint := "your sonar endpoit" // like http://sonar.alauda.cn
    token := "your sonar token"
    sonar := sonarqube.NewSonarQubeArgs(endpoint, token)
    name := "sonar-test"
	projectKey := "sonar-test-key"

    err := sonar.CreateProject(name, projectKey)
    
    if err != nil {
        log.Fatalf("method CreateProject error: %v", err)
    }
}
```

- list quality gates
```go
package main

import(
    "fmt"
    "log"

    "github.com/bergamot/alauda/sonarqube"
)

func main() {
    endpoint := "your sonar endpoit" // like http://sonar.alauda.cn
    token := "your sonar token"
    sonar := sonarqube.NewSonarQubeArgs(endpoint, token)

    ret, err := sonar.ListQualityGates(name, projectKey)
    
    if err != nil {
        log.Fatalf("method ListQualityGates error: %v", err)
    }
    fmt.Println("quality gates of sonar are: ", ret)
}
```

## Contribute
1. clone code
```
git clone https://github.com/alauda/bergamot.git $GOPATH/src/github.com/alauda/bergamot
```
2. write code
3. test, **remember provide SONAR_ENDPOINT and SONAR_TOKEN when test**:
```
SONAR_ENDPOINT=sonar_url SONAR_TOKEN=sonar_token make test
```
4. pr on GitHub