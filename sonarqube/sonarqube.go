package sonarqube

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	alog "github.com/alauda/bergamot/loggo"
	"github.com/alauda/bergamot/utils"
	"github.com/parnurzeal/gorequest"
)

func getSonarAuthToken(token string) string {
	b64 := base64.StdEncoding.EncodeToString([]byte(token + ":"))
	return fmt.Sprintf("Basic %s", b64)
}

// SonarQube client of Sonarqube
type SonarQube struct {
	Endpoint   string
	Version    string
	Token      string
	Logger     alog.Logger
	httpClient *gorequest.SuperAgent
	authToken  string
}

// NewSonarQube return sonarqube, retrive endpoint and token from env
func NewSonarQube() *SonarQube {
	endpoint := os.Getenv("SONAR_ENDPOINT")
	token := os.Getenv("SONAR_TOKEN")
	if len(endpoint) == 0 {
		panic("NewSonarQube method need env SONAR_ENDPOINT, but missing")
	}
	if len(token) == 0 {
		panic("NewSonarQube method need env SONAR_TOKEN, but missing")
	}

	return NewSonarQubeArgs(endpoint, token)
}

// NewSonarQubeArgs return sonarqube with args init
func NewSonarQubeArgs(endpoint, token string) *SonarQube {
	var sonar SonarQube
	sonar.Endpoint = endpoint
	sonar.Token = token
	sonar.Logger = alog.GetLogger("sonarqube")

	sonar.authToken = getSonarAuthToken(token)
	sonar.httpClient = gorequest.New()

	version, err := sonar.GetVersion()
	if err != nil {
		sonar.Logger.Errorf("[NewSonarQubeArgs] - get sonar version error: %s", err)
		panic(err)
	}
	sonar.Version = version

	return &sonar
}

// SystemStatus get sonarqube system status
func (sonar *SonarQube) SystemStatus() (map[string]string, error) {
	path, err := utils.GetURL(sonar.Endpoint, "api/system/status", nil)
	if err != nil {
		return nil, err
	}
	resp, body, errs := sonar.httpClient.Get(path).End()
	if errs != nil {
		err := fmt.Errorf("get system status error:%v", errs)
		return nil, err
	}
	defer utils.CloseResponse(resp)

	var status map[string]string
	err = json.Unmarshal([]byte(body), &status)
	if err != nil {
		return nil, err
	}

	return status, nil
}

// GetVersion get sonarqube server version like 6.0.0
func (sonar *SonarQube) GetVersion() (string, error) {
	status, err := sonar.SystemStatus()
	if err != nil {
		return "", err
	}
	versions := strings.Split(status["version"], ".")
	for len(versions) < 3 {
		versions = append(versions, "0")
	}

	return strings.Join(versions[:3], "."), nil
}

const (
	SONAR_VERSION_60 = "6.0.0"
	SONAR_VERSION_64 = "6.4.0"
)

// CreateProject create sonar project
func (sonar *SonarQube) CreateProject(name, projectKey string) error {
	query := make(url.Values)
	query["name"] = []string{name}

	if sonar.Version == SONAR_VERSION_60 {
		query["key"] = []string{projectKey}
	} else if sonar.Version == SONAR_VERSION_64 {
		query["project"] = []string{projectKey}
	} else {
		return fmt.Errorf("not support sonar version %s", sonar.Version)
	}
	path, err := utils.GetURL(sonar.Endpoint, "api/projects/create", query)
	if err != nil {
		sonar.Logger.Errorf("[GetURL] - encount error: %v", err)
		return err
	}

	resp, body, errs := sonar.httpClient.Post(path).
		Set("Authorization", sonar.authToken).
		End()
	if errs != nil {
		err = fmt.Errorf("[CreateProject] - create project error: %v", errs)
		sonar.Logger.Error(err.Error())
		return err
	}
	defer utils.CloseResponse(resp)

	if strings.Contains(body, "key already exists") {
		sonar.Logger.Infof("project key %s has existed, will not create it", projectKey)
	}

	return nil
}

// GetProjectID get project id by key, only for sonar 6.0.0
func (sonar *SonarQube) GetProjectID(projectKey string) (string, error) {
	if sonar.Version != SONAR_VERSION_60 {
		return "", fmt.Errorf("not support sonar version %s", sonar.Version)
	}

	query := url.Values{
		"key": []string{projectKey},
	}
	path, err := utils.GetURL(sonar.Endpoint, "api/projects/index", query)
	if err != nil {
		sonar.Logger.Errorf("[GetURL] - encount error: %v", err)
		return "", err
	}
	resp, body, errs := sonar.httpClient.Get(path).
		Set("Authorization", sonar.authToken).
		End()
	if errs != nil {
		err = fmt.Errorf("[GetProjectId] - get project id error: %v", errs)
		sonar.Logger.Error(err.Error())
		return "", err
	}
	defer utils.CloseResponse(resp)

	type retMap map[string]string
	var ret []retMap
	if err := json.Unmarshal([]byte(body), &ret); err != nil {
		return "", fmt.Errorf("[GetProjectId - parse response body to json error: %v", err)
	}
	if len(ret) != 1 {
		return "", fmt.Errorf("[GetProjectId - result should be one, but got %v", ret)
	}

	return ret[0]["id"], nil
}

// GetSettings get sonar project settings since 6.3
// 和alauda-sonar-scanner中的sonarClient兼容，200 ~ 400 返回json 数据, nil,>=400, 返回字符串, error
func (sonar *SonarQube) GetSettings(component string, keys []string) (interface{}, error) {
	if sonar.Version != SONAR_VERSION_64 {
		return nil, fmt.Errorf("method [GetSettings] does not support sonar version %s", sonar.Version)
	}
	var query = url.Values{
		"component": []string{component},
		"keys":      []string{strings.Join(keys, ",")},
	}
	path, err := utils.GetURL(sonar.Endpoint, "api/settings/values", query)
	if err != nil {
		sonar.Logger.Errorf("[GetURL] - encount error: %v", err)
		return nil, err
	}

	resp, _, errs := sonar.httpClient.Get(path).
		Set("Authorization", sonar.authToken).
		End()
	if errs != nil {
		err = fmt.Errorf("[GetSettings] - get component settings error: %v", errs)
		sonar.Logger.Error(err.Error())
		return nil, err
	}
	defer utils.CloseResponse(resp)

	return sonar.parseResponse(resp)
}

// SetSettings set sonar project settings
func (sonar *SonarQube) SetSettings(
	component string, fieldValues map[string]string,
	key string, value string, values []string,
) (interface{}, error) {

	if sonar.Version != SONAR_VERSION_64 {
		return nil, fmt.Errorf("method [GetSettings] does not support sonar version %s", sonar.Version)
	}

	fieldValuesStr := ""
	if len(fieldValues) != 0 {
		bts, err := json.Marshal(fieldValues)
		if err != nil {
			err := fmt.Errorf("[SetSettings] - marshal fieldValues: %v error: %s", fieldValues, err)
			return nil, err
		}
		fieldValuesStr = string(bts)
	}

	var query = url.Values{
		"component": []string{component},
	}
	if fieldValuesStr != "" {
		query["fieldValues"] = []string{fieldValuesStr}
	}
	if key != "" {
		query["key"] = []string{key}
	}
	if value != "" {
		query["value"] = []string{value}
	}
	if len(values) != 0 {
		query["values"] = values
	}

	path, err := utils.GetURL(sonar.Endpoint, "api/settings/set", query)
	if err != nil {
		sonar.Logger.Errorf("[GetURL] - encount error: %v", err)
		return nil, err
	}

	resp, _, errs := sonar.httpClient.Post(path).
		Set("Authorization", sonar.authToken).
		End()
	if errs != nil {
		err = fmt.Errorf("[GetSettings] - set component settings error: %v", errs)
		sonar.Logger.Error(err.Error())
		return nil, err
	}
	defer utils.CloseResponse(resp)

	return sonar.parseResponse(resp)
}

// ListQualityGates
func (sonar *SonarQube) ListQualityGates() (interface{}, error) {
	path, err := utils.GetURL(sonar.Endpoint, "api/qualitygates/list", nil)
	if err != nil {
		sonar.Logger.Errorf("[GetURL] - encount error: %v", err)
		return nil, err
	}

	resp, _, errs := sonar.httpClient.Get(path).
		Set("Authorization", sonar.authToken).
		End()
	if errs != nil {
		err = fmt.Errorf("[GetSettings] - get qualitygates error: %v", errs)
		sonar.Logger.Error(err.Error())
		return nil, err
	}
	defer utils.CloseResponse(resp)

	return sonar.parseResponse(resp)
}

// SelectQualityGates
func (sonar *SonarQube) SelectQualityGates(gateID int, projectID, projectKey string) (interface{}, error) {
	var query = url.Values{
		"gateId": []string{fmt.Sprint(gateID)},
	}
	if projectID != "" {
		query.Add("projectId", projectID)
	}
	if projectKey != "" {
		if sonar.Version == SONAR_VERSION_64 {
			query.Add("projectKey", projectKey)
		}
	}

	path, err := utils.GetURL(sonar.Endpoint, "api/qualitygates/select", query)
	if err != nil {
		sonar.Logger.Errorf("[GetURL] - encount error: %v", err)
		return nil, err
	}

	resp, _, errs := sonar.httpClient.Post(path).
		Set("Authorization", sonar.authToken).
		End()
	if errs != nil {
		err = fmt.Errorf("[GetSettings] - select qualitygates error: %v", errs)
		sonar.Logger.Error(err.Error())
		return nil, err
	}
	defer utils.CloseResponse(resp)

	return sonar.parseResponse(resp)
}

func (sonar *SonarQube) parseResponse(resp *http.Response) (interface{}, error) {
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		data, err := decodeBody(resp.Body)
		if err != nil {
			return nil, errors.New("Decode response body Error: " + err.Error())
		}
		return data, nil
	}

	if resp.StatusCode >= 500 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return data, errors.New("Read response body Error: " + err.Error())
		}

		return data, fmt.Errorf("ServerError5xx, response body is: %s", data)
	}

	dbts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read Response Body Content Error: %s", err.Error())
	}
	data := string(dbts)

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return data, fmt.Errorf("ClientError4xx, response body is: %s", data)
	}
	return data, fmt.Errorf("UnKnown Status Code, response body is: %s", data)
}

func decodeBody(reader io.Reader) (interface{}, error) {
	var data interface{}
	cts, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Read Body Content Error: %s", err.Error())
	}
	if len(cts) == 0 {
		return nil, nil
	}
	if err = json.Unmarshal(cts, &data); err != nil {
		return nil, fmt.Errorf("Unmarshall to json Error: %s , Body: %s", err.Error(), string(cts))
	}

	return data, err
}
