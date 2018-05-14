package utils_test

import (
	"net/url"
	"testing"

	"github.com/alauda/bergamot/utils"
)

func TestGetURL(t *testing.T) {
	testTables := []struct {
		Endpoint string
		Path     string
		Values   url.Values
		Expected string
	}{
		{
			Endpoint: "http://alauda.cn",
			Path:     "api/ping",
			Values: url.Values{
				"name": []string{"wqlu"},
				"age":  []string{"18"},
			},
			Expected: "http://alauda.cn/api/ping?age=18&name=wqlu",
		},
		{
			Endpoint: "https://alauda.cn",
			Path:     "api/projects/create",
			Values: url.Values{
				"name":    []string{"sonar"},
				"project": []string{"anbc"},
			},
			Expected: "https://alauda.cn/api/projects/create?name=sonar&project=anbc",
		},
		{
			Endpoint: "http://alauda.cn",
			Path:     "api/ping/",
			Values: url.Values{
				"name": []string{"wqlu"},
				"age":  []string{"12"},
			},
			Expected: "http://alauda.cn/api/ping/?age=12&name=wqlu",
		},
		{
			Endpoint: "http://alauda.cn",
			Path:     "/api/list/",
			Values: url.Values{
				"name": []string{"wqlu"},
				"age":  []string{"18"},
			},
			Expected: "http://alauda.cn/api/list/?age=18&name=wqlu",
		},
		{
			Endpoint: "http://alauda.cn",
			Path:     "/api/list/",
			Values:   nil,
			Expected: "http://alauda.cn/api/list/",
		},
	}

	for _, table := range testTables {
		ret, err := utils.GetURL(table.Endpoint, table.Path, table.Values)
		if err != nil {
			t.Errorf("test case %v error:%v", table, err)
		}
		if ret != table.Expected {
			t.Errorf("result is expected %s, but got %s", table.Expected, ret)
		}
	}
}
