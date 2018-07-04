package sonarqube

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/alauda/bergamot/log"
	alog "github.com/alauda/bergamot/loggo"
	"github.com/alauda/bergamot/utils"
	"github.com/parnurzeal/gorequest"
)

// sonarHTTPClient for retry func
type sonarHTTPClient struct {
	superAgent       *gorequest.SuperAgent
	retryTimes       int
	authToken        string
	retryStatusCodes []int
	Logger           log.BasicLogger
}

func (httpClient *sonarHTTPClient) Get(path string) (resp gorequest.Response, body string, errs []error) {
	httpClient.Logger.Debugf("request %s with GET method...", path)
	resp, body, errs = httpClient.superAgent.
		Get(path).
		Retry(httpClient.retryTimes, 1*time.Millisecond, httpClient.retryStatusCodes...).
		Set("Authorization", httpClient.authToken).
		End()
	return
}

func (httpClient *sonarHTTPClient) Post(path string) (resp gorequest.Response, body string, errs []error) {
	httpClient.Logger.Debugf("request %s with POST method...", path)
	resp, body, errs = httpClient.superAgent.
		Post(path).
		Retry(httpClient.retryTimes, 1*time.Millisecond, httpClient.retryStatusCodes...).
		Set("Authorization", httpClient.authToken).
		End()

	return
}

func getSonarAuthToken(token string) string {
	b64 := base64.StdEncoding.EncodeToString([]byte(token + ":"))
	return fmt.Sprintf("Basic %s", b64)
}

// GetDefaultLogger get default logger for sonar
func GetDefaultLogger(prefix string) log.BasicLogger {
	logger := alog.GetLogger(prefix)
	logger.SetLogLevel(alog.DEBUG)

	return logger
}

// SonarQube client of Sonarqube
type SonarQube struct {
	Endpoint   string
	Version    string
	Token      string
	Logger     log.BasicLogger
	httpClient *sonarHTTPClient
}

// NewSonarQube return sonarqube, retrive endpoint and token from env
func NewSonarQube() (*SonarQube, error) {
	endpoint := os.Getenv("SONAR_ENDPOINT")
	token := os.Getenv("SONAR_TOKEN")
	if len(endpoint) == 0 {
		return nil, fmt.Errorf("NewSonarQube method need env SONAR_ENDPOINT, but missing")
	}
	if len(token) == 0 {
		return nil, fmt.Errorf("NewSonarQube method need env SONAR_TOKEN, but missing")
	}

	return NewSonarQubeArgs(endpoint, token)
}

// NewSonarQubeArgs return sonarqube with args init
func NewSonarQubeArgs(endpoint, token string) (*SonarQube, error) {
	var sonar SonarQube
	sonar.Endpoint = endpoint
	sonar.Token = token
	sonar.Logger = GetDefaultLogger("[alauda-sonarqube]")

	sonar.httpClient = &sonarHTTPClient{
		authToken:  getSonarAuthToken(token),
		retryTimes: 3,
		superAgent: gorequest.New(),
		retryStatusCodes: []int{
			http.StatusTooManyRequests,
			http.StatusRequestTimeout,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		},
		Logger: sonar.Logger,
	}

	version, err := sonar.GetVersion()
	if err != nil {
		sonar.Logger.Errorf("[NewSonarQubeArgs] - get sonar version error: %s", err)
		return nil, err
	}
	sonar.Version = version

	return &sonar, nil
}

// SetLogger set logger to sonar default logger
func (sonar *SonarQube) SetLogger(logger log.BasicLogger) {
	if logger == nil {
		logger = GetDefaultLogger("[alauda-sonarqube]")
	}
	sonar.Logger = logger
	if sonar.httpClient != nil {
		sonar.httpClient.Logger = logger
	}
}

// IsVersion60 true version 6.0.0
func (sonar *SonarQube) IsVersion60() bool {
	return sonar.Version == SONAR_VERSION_60
}

// IsVersion64 true version 6.4.0
func (sonar *SonarQube) IsVersion64() bool {
	return sonar.Version == SONAR_VERSION_64
}

// SystemStatus get sonarqube system status
func (sonar *SonarQube) SystemStatus() (map[string]string, error) {
	path, err := utils.GetURL(sonar.Endpoint, "api/system/status", nil)
	if err != nil {
		return nil, err
	}
	resp, body, errs := sonar.httpClient.Get(path)
	if errs != nil {
		err := fmt.Errorf("request system status error:%v", errs)
		return nil, err
	}
	defer utils.CloseResponse(resp)

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("not expected status code, code is %d",
			resp.StatusCode)
		sonar.Logger.Errorf("not expected response, code is %d, body is %s",
			resp.StatusCode, body)
		return nil, err
	}

	var status map[string]string
	err = json.Unmarshal([]byte(body), &status)
	if err != nil {
		sonar.Logger.Errorf("parse response body error, body is %s, error is %v",
			body, err)
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

	if sonar.IsVersion60() {
		query["key"] = []string{projectKey}
	} else if sonar.IsVersion64() {
		query["project"] = []string{projectKey}
	} else {
		return fmt.Errorf("not support sonar version %s", sonar.Version)
	}
	path, err := utils.GetURL(sonar.Endpoint, "api/projects/create", query)
	if err != nil {
		sonar.Logger.Errorf("[GetURL] - encount error: %v", err)
		return err
	}

	resp, body, errs := sonar.httpClient.Post(path)
	if errs != nil {
		err = fmt.Errorf("[CreateProject] - create project error: %v", errs)
		return err
	}
	defer utils.CloseResponse(resp)

	switch resp.StatusCode {
	case http.StatusOK:
		sonar.Logger.Infof("create project %s success", name)
		return nil
	case http.StatusBadRequest:
		if strings.Contains(body, "key already exists") {
			sonar.Logger.Infof("project key %s has existed, will not create project", projectKey)
			return nil
		}
		sonar.Logger.Debugf("bad request, response code is %d, body is %s", resp.StatusCode, body)
		return fmt.Errorf("bad request, response code is %d", resp.StatusCode)
	default:
		sonar.Logger.Errorf("not expected response, code is %d, body is %s", resp.StatusCode, body)
		return fmt.Errorf("not expected response, code is %d", resp.StatusCode)
	}
}

// GetProjectID get project id by key, only for sonar 6.0.0
func (sonar *SonarQube) GetProjectID(projectKey string) (string, error) {
	if !sonar.IsVersion60() {
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
	resp, body, errs := sonar.httpClient.Get(path)
	if errs != nil {
		err = fmt.Errorf("[GetProjectId] - get project id error: %v", errs)
		sonar.Logger.Errorf("%v", err)
		return "", err
	}
	defer utils.CloseResponse(resp)

	switch resp.StatusCode {
	case http.StatusOK:
		type retMap map[string]string
		var ret []retMap
		if err := json.Unmarshal([]byte(body), &ret); err != nil {
			sonar.Logger.Errorf("parse response body error, response code is %d, body is %s, error is %v",
				resp.StatusCode, body, err)
			return "", fmt.Errorf("parse response body error, response code is %d, error is %v",
				resp.StatusCode, err)
		}
		if len(ret) != 1 {
			return "", fmt.Errorf("result of GetProjectID should be one, but got %v", ret)
		}
		return ret[0]["id"], nil
	default:
		return "", fmt.Errorf("not expected response, code is %d", resp.StatusCode)
	}
}

// GetSettings get sonar project settings since 6.3
// 和alauda-sonar-scanner中的sonarClient兼容，200 ~ 400 返回json 数据,nil;>=400, 返回字符串,error
func (sonar *SonarQube) GetSettings(component string, keys []string) (interface{}, error) {
	switch sonar.Version {
	case SONAR_VERSION_64:
		var query = url.Values{
			"component": []string{component},
			"keys":      []string{strings.Join(keys, ",")},
		}
		path, err := utils.GetURL(sonar.Endpoint, "api/settings/values", query)
		if err != nil {
			sonar.Logger.Errorf("[GetURL] - encount error: %v", err)
			return nil, err
		}

		resp, _, errs := sonar.httpClient.Get(path)
		if errs != nil {
			err = fmt.Errorf("[GetSettings] - get component settings error: %v", errs)
			sonar.Logger.Errorf("%v", err)
			return nil, err
		}
		defer utils.CloseResponse(resp)

		return sonar.parseResponse(resp)
	case SONAR_VERSION_60:
		return nil, fmt.Errorf("method [GetSettings] does not support sonar version %s", sonar.Version)
	default:
		return nil, fmt.Errorf("method [GetSettings] does not support sonar version %s", sonar.Version)
	}
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

	resp, _, errs := sonar.httpClient.Post(path)
	if errs != nil {
		err = fmt.Errorf("[GetSettings] - set component settings error: %v", errs)
		sonar.Logger.Errorf("%v", err)
		return nil, err
	}
	defer utils.CloseResponse(resp)

	return sonar.parseResponse(resp)
}

// ListLanguages list sonar languages
func (sonar *SonarQube) ListLanguages() (interface{}, error) {
	path, err := utils.GetURL(sonar.Endpoint, "api/languages/list", nil)
	if err != nil {
		sonar.Logger.Errorf("[GetURL] - encount error: %v", err)
		return nil, err
	}

	resp, _, errs := sonar.httpClient.Get(path)
	if errs != nil {
		err = fmt.Errorf("[GetSettings] - get languages error: %v", errs)
		sonar.Logger.Errorf("%v", err)
		return nil, err
	}
	defer utils.CloseResponse(resp)

	return sonar.parseResponse(resp)
}

// ListQualityGates list sonar quality gates
func (sonar *SonarQube) ListQualityGates() (interface{}, error) {
	path, err := utils.GetURL(sonar.Endpoint, "api/qualitygates/list", nil)
	if err != nil {
		sonar.Logger.Errorf("[GetURL] - encount error: %v", err)
		return nil, err
	}

	resp, _, errs := sonar.httpClient.Get(path)
	if errs != nil {
		err = fmt.Errorf("[GetSettings] - get qualitygates error: %v", errs)
		sonar.Logger.Errorf("%v", err)
		return nil, err
	}
	defer utils.CloseResponse(resp)

	return sonar.parseResponse(resp)
}

// SelectQualityGates select sonar quality gates
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

	resp, _, errs := sonar.httpClient.Post(path)
	if errs != nil {
		err = fmt.Errorf("[GetSettings] - select qualitygates error: %v", errs)
		sonar.Logger.Errorf("%v", err)
		return nil, err
	}
	defer utils.CloseResponse(resp)

	return sonar.parseResponse(resp)
}

func (sonar *SonarQube) parseResponse(resp *http.Response) (interface{}, error) {
	if resp == nil {
		return nil, fmt.Errorf("response is nil, skip encode")
	}

	var (
		err        error
		statusCode int
	)

	statusCode = resp.StatusCode
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return string(data), fmt.Errorf("Read response body error: %v", err)
	}

	if statusCode >= http.StatusInternalServerError { // 500
		sonar.Logger.Errorf("ServerError5xx, response code is %d, response body is: %s", statusCode, data)
		err = fmt.Errorf("ServerError5xx, response code is %d", statusCode)
	} else if statusCode >= http.StatusBadRequest { // 400
		sonar.Logger.Debugf("ClientError4xx, response code is %d, response body is: %s", statusCode, data)
		err = fmt.Errorf("ClientError4xx, response code is %d", statusCode)
	} else if statusCode >= http.StatusOK { // 200
		var ret interface{}
		if len(data) > 0 {
			if err = json.Unmarshal(data, &ret); err != nil {
				return nil, fmt.Errorf("Unmarshall to json error: %v, body: %s", err, data)
			}
		}

		return ret, nil
	}

	err = fmt.Errorf("UnKnown response status code %d", statusCode)
	return data, err
}
