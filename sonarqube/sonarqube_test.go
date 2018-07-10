package sonarqube_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/alauda/bergamot/sonarqube"
	"github.com/stretchr/testify/assert"
)

func TestParseMap(t *testing.T) {
	var status map[string]string
	ret := `{"id":"AWMLtw7lJgn-wo0wWz3H","version":"6.4.0.25310","status":"UP"}`
	err := json.Unmarshal([]byte(ret), &status)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, err, nil)
	assert.Equal(t, status["id"], "AWMLtw7lJgn-wo0wWz3H")
}

func getSonarEnv() (string, string, bool) {
	endpoint := os.Getenv("SONAR_ENDPOINT")
	token := os.Getenv("SONAR_TOKEN")
	if len(endpoint) == 0 || len(token) == 0 {
		log.Printf("env endpoint %s or token %s is missing", endpoint, token)
		return endpoint, token, false
	}
	log.Printf("sonar endpoint: %s, token: %s", endpoint, token)

	return endpoint, token, true
}

func getSonar(t *testing.T) *sonarqube.SonarQube {
	endpoint, token, valid := getSonarEnv()
	if !valid {
		t.Skipf("env endpoint %s or token %s is missing, skip test", endpoint, token)
	}

	sonar, err := sonarqube.NewSonarQubeArgs(endpoint, token)
	if err != nil {
		t.Errorf("new Sonarqube error: %v", err)
	}

	return sonar
}

func TestNewSonarQubeArgs(t *testing.T) {
	endpoint, token, valid := getSonarEnv()
	if !valid {
		t.Skipf("env endpoint %s or token %s is missing", endpoint, token)
	}

	sonar, err := sonarqube.NewSonarQubeArgs(endpoint, token)
	assert.Equal(t, err, nil)
	assert.NotEmpty(t, sonar.Version)
	t.Logf("sonar version is %s", sonar.Version)
}

func TestSonarQube_SystemStatus(t *testing.T) {
	sonar := getSonar(t)
	ret, err := sonar.SystemStatus()
	if err != nil {
		t.Fatalf("test SystemStatus error: %v", err)
	}

	assert.NotNil(t, ret)
	assert.Equal(t, len(ret), 3)
	t.Logf("system status is %v", ret)
}

func TestSonarQube_GetVersion(t *testing.T) {
	sonar := getSonar(t)
	assert.NotEmpty(t, sonar.Version)
	t.Logf("sonar version is %s", sonar.Version)
}

func TestSonarQube_CreateProject(t *testing.T) {
	sonar := getSonar(t)
	name := "sonar-test"
	projectKey := "sonar-test-key"

	err := sonar.CreateProject(name, projectKey)
	t.Logf("err of CreateProject is: %v", err)
}

func TestSonarQube_ListLanguages(t *testing.T) {
	sonar := getSonar(t)
	ret, err := sonar.ListLanguages()

	assert.Nil(t, err)
	assert.NotNil(t, ret)
	t.Logf("languages of sonar are: %v", ret)
}

func TestSonarQube_ListQualityGates(t *testing.T) {
	sonar := getSonar(t)
	ret, err := sonar.ListQualityGates()

	assert.Nil(t, err)
	assert.NotNil(t, ret)
	t.Logf("quality gates of sonar are: %v", ret)
}

func TestSonarQube_SelectQualityGates(t *testing.T) {
	sonar := getSonar(t)
	gateId := 1
	projectId := "1"
	projectKey := "sonar-test-key"

	ret, err := sonar.SelectQualityGates(gateId, projectId, projectKey)
	t.Logf("response body of SelectQualityGates is: %v", ret)
	t.Logf("err of SelectQualityGates is: %v", err)
}

func TestSonarQube_GetQualityGatesProjectStatusData(t *testing.T) {
	t.Skip("need analysisId or projectId or projectKey")
	sonar := getSonar(t)
	analysisId := "AWSCwIISdW2Esay5WZC2"
	projectId := ""
	projectKey := ""

	ret, err := sonar.GetQualityGatesProjectStatusData(analysisId, projectId, projectKey)
	t.Logf("response body of SelectQualityGates is: %v", ret)
	t.Logf("err of SelectQualityGates is: %v", err)
}

func TestSonarQube_GetQualityGatesProjectStatus(t *testing.T) {
	t.Skip("need analysisId or projectId or projectKey")
	sonar := getSonar(t)
	analysisId := "AWSCwIISdW2Esay5WZC2"
	projectId := ""
	projectKey := ""

	ret, err := sonar.GetQualityGatesProjectStatus(analysisId, projectId, projectKey)
	t.Logf("response body of SelectQualityGates is: %v", ret)
	t.Logf("err of SelectQualityGates is: %v", err)
}

func TestSonarQube_GetSettings(t *testing.T) {
	sonar := getSonar(t)
	component := "sonar-test-key"
	keys := []string{"sonar.dbcleaner.cleanDirectory", "sonar.webhooks.project"}

	ret, err := sonar.GetSettings(component, keys)
	t.Logf("response body of SelectQualityGates is: %v", ret)
	t.Logf("err of SelectQualityGates is: %v", err)
}

func TestSonarQube_SetSettings(t *testing.T) {
	sonar := getSonar(t)
	key := "sonar.webhooks.project"
	projectKey := "sonar-test-key"
	keys := map[string]string{
		"name": "alauda-ci",
		"url":  "http://alauda.cn/webhook",
	}

	ret, err := sonar.SetSettings(projectKey, keys, key, "", []string{})
	t.Logf("response body of SetSettings is: %v", ret)
	t.Logf("err of SetSettings is: %v", err)
}

func TestSonarQube_GetProjectID(t *testing.T) {
	sonar := getSonar(t)
	key := "sonar-key"
	ret, err := sonar.GetProjectID(key)

	t.Logf("err is %v", err)
	t.Logf("ret is %v", ret)
}

func TestSonarQube_GetAnalysisTask(t *testing.T) {
	t.Skip("need task id")
	sonar := getSonar(t)
	taskID := "AWR-rnuXELrk8W1dtIpm"
	ret, err := sonar.GetAnalysisTask(taskID)

	t.Logf("err is %v", err)
	t.Logf("ret is %v", ret)
}

func TestSonarQube_GetAnalysisTaskStatus(t *testing.T) {
	t.Skip("need task id")
	sonar := getSonar(t)
	taskID := "AWR-rnuXELrk8W1dtIpm"
	ret, err := sonar.GetAnalysisTaskStatus(taskID)

	t.Logf("err is %v", err)
	t.Logf("ret is %v", ret)
}
