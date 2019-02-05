package main

import (
	"errors"
	"fmt"
	"testing"
)

type ResourceDataMockGet struct {
	Arguments struct {
		Key string
	}
	Returns struct {
		Value interface{}
	}
}

type ResourceDataMockGetOk struct {
	Arguments struct {
		Key string
	}
	Returns struct {
		Value interface{}
		OK    bool
	}
}

type ResourceDataMockSet struct {
	Arguments struct {
		Key   string
		Value interface{}
	}
	Returns struct {
		Error error
	}
}

type ResourceDataMockSetId struct {
	Arguments struct {
		ID string
	}
}

type ResourceDataMock struct {
	GetCalls       []ResourceDataMockGet
	GetCallCount   int
	GetOkCalls     []ResourceDataMockGetOk
	GetOkCallCount int
	SetCalls       []ResourceDataMockSet
	SetCallCount   int
	SetIdCalls     []ResourceDataMockSetId
	SetIdCallCount int
}

func (mock *ResourceDataMock) Get(key string) interface{} {
	i := mock.GetCallCount
	mock.GetCallCount++

	mock.GetCalls[i].Arguments.Key = key
	return mock.GetCalls[i].Returns.Value
}

func (mock *ResourceDataMock) GetOk(key string) (interface{}, bool) {
	i := mock.GetOkCallCount
	mock.GetOkCallCount++

	mock.GetOkCalls[i].Arguments.Key = key
	return mock.GetOkCalls[i].Returns.Value, mock.GetOkCalls[i].Returns.OK
}

func (mock *ResourceDataMock) Set(key string, val interface{}) error {
	i := mock.SetCallCount
	mock.SetCallCount++

	mock.SetCalls[i].Arguments.Key = key
	mock.SetCalls[i].Arguments.Value = val
	return mock.SetCalls[i].Returns.Error
}

func (mock *ResourceDataMock) SetId(id string) {
	i := mock.SetIdCallCount
	mock.SetIdCallCount++

	mock.SetIdCalls[i].Arguments.ID = id
}

func TestInterpolate(t *testing.T) {
	t.Run("when no variables are set", func(t *testing.T) {
		templ := "Hello, world!"

		getMock := ResourceDataMockGet{}
		getMock.Returns.Value = templ
		setMock := ResourceDataMockSet{}
		setMock.Returns.Error = nil
		getVarsMock := ResourceDataMockGetOk{}
		getVarsMock.Returns.Value = nil
		getVarsMock.Returns.OK = false
		mock := &ResourceDataMock{}
		mock.GetCalls = append(mock.GetCalls, getMock)
		mock.GetOkCalls = append(mock.GetOkCalls, getVarsMock)
		mock.SetCalls = append(mock.SetCalls, setMock)

		interpolate(mock, nil, nil)

		if mock.GetCallCount < 1 || mock.GetCalls[0].Arguments.Key != "template" {
			t.Errorf(`it should call Get 1 time requesting "template"`)
		}
		if mock.GetOkCallCount < 1 || mock.GetOkCalls[0].Arguments.Key != "vars" {
			t.Errorf(`it should call GetOk 1 time requesting "vars"`)
		}
		if mock.SetCallCount < 1 || mock.SetCalls[0].Arguments.Key != "output" {
			t.Errorf(`it should call Set once with key "output"`)
		}
		if val, ok := mock.SetCalls[0].Arguments.Value.(string); !ok || val != templ {
			t.Errorf(`it should call Set with a string value %q`, templ)
		}
	})

	t.Run("when variables exists", func(t *testing.T) {
		type data struct {
			Template, Greeting, Name, Output string
			Err                              error
		}

		template := "{{.greeting}}{{if and .greeting .who}}, {{end}}{{if .who}}{{.who}}{{end}}!"

		for _, td := range []data{
			{template, "Hello", "Nick", "Hello, Nick!", nil},
			{template, "Greetings", "Zack", "Greetings, Zack!", nil},
			{template, "", "Sara", "Sara!", nil},
			{template, "Hello", "", "Hello!", nil},
			{"bad template {{end}}", "", "", "", errors.New("invalid template")},
		} {
			t.Run(fmt.Sprintf(`when "greeting" is %q and "who" is %q`, td.Greeting, td.Name), func(t *testing.T) {

				getTemplateMock := ResourceDataMockGet{}
				getTemplateMock.Returns.Value = template

				getVarsMock := ResourceDataMockGetOk{}
				getVarsMock.Returns.Value = map[string]interface{}{
					"greeting": td.Greeting,
					"who":      td.Name,
				}
				getVarsMock.Returns.OK = true

				setMock := ResourceDataMockSet{}
				setMock.Returns.Error = nil

				mock := &ResourceDataMock{}
				mock.GetCalls = append(mock.GetCalls, getTemplateMock)
				mock.GetOkCalls = append(mock.GetOkCalls, getVarsMock)
				mock.SetCalls = append(mock.SetCalls, setMock)

				err := interpolate(mock, nil, nil)

				// Assertions
				if td.Err == nil && err != nil {
					t.Error("it should not return an error")
				}
				if td.Err != nil && (err == nil || err.Error() != td.Err.Error()) {
					t.Errorf(`it should return error %q but it returned "%v"`, td.Err, err)
				}

				if mock.GetCallCount < 1 || mock.GetCalls[0].Arguments.Key != "template" {
					t.Error(`it should call Get with key "template"`)
				}

				if mock.GetOkCallCount < 1 || mock.GetOkCalls[0].Arguments.Key != "vars" {
					t.Error(`it should call GetOk with key "vars"`)
				}

				if td.Err != nil {
					return
				}

				// testing Set
				if mock.SetCallCount < 1 {
					t.Error(`it should call Set`)
					return
				}
				if val, ok := mock.SetCalls[0].Arguments.Value.(string); !ok ||
					mock.SetCalls[0].Arguments.Key != "output" ||
					val != td.Output {
					t.Errorf(`it should call Set with key %q a string value %v`, "output", td.Output)
					t.Logf("got: %q", val)
				}
			})
		}
	})

	t.Run("when template function cidrhost is called", func(t *testing.T) {
		t.SkipNow()

		// template := `reserved_ip_ranges: {{cidrhost .pas_subnet_cidr 0}}-{{cidrhost .pas_subnet_cidr 5}}`
		// vars := map[string]interface{}{
		// 	"pas_subnet_cidr": "10.0.0.0/16",
		// }

		// interpolate(V
	})
}
