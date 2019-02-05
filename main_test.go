package main

import (
	"bytes"
	"errors"
	"testing"

	mocking "github.com/golang/mock/gomock"
)

// ************************************
// Mocks

//go:generate mockgen -source main.go -destination resource_data_mock_test.go -package main

// ************************************
// Unit Tests

func TestExecute(t *testing.T) {
	testGenerateID := func(_ string) string { return "some-hash" }

	t.Run("when no variables are set", func(t *testing.T) {
		template := "Hello, world!"

		mockCtrl := mocking.NewController(t)
		defer mockCtrl.Finish()
		dataResource := NewMockResourceData(mockCtrl)

		dataResource.EXPECT().Get(mocking.Eq("template")).Return(template)
		dataResource.EXPECT().GetOk(mocking.Eq("vars")).Return(nil, false)
		dataResource.EXPECT().Set(mocking.Eq("output"), mocking.Eq(template)).Return(nil)
		dataResource.EXPECT().SetId(mocking.Eq("some-hash"))

		err := execute(dataResource, nil, testGenerateID, defaultTemplate())

		if !gotExpectedErr(t, err, nil) {
			t.Fail()
		}
	})

	t.Run("when template is invalid", func(t *testing.T) {
		template := "Hello, world!"

		mockCtrl := mocking.NewController(t)
		defer mockCtrl.Finish()
		dataResource := NewMockResourceData(mockCtrl)

		dataResource.EXPECT().Get(mocking.Eq("template")).Return(template)
		// SetId, Get("vars"), and Set("output", output) should not be called

		err := execute(dataResource, nil, testGenerateID, defaultTemplate())

		if !gotExpectedErr(t, err, errors.New("invalid template")) {
			t.Fail()
		}
	})

	t.Run("when a valid template executes with variables", func(t *testing.T) {
		type testdata struct {
			Output, Template string
			Vars             map[string]interface{}
		}

		greetingTemplate := "{{.greeting}}{{if and .greeting .who}}, {{end}}{{if .who}}{{.who}}{{end}}!"
		var idSet Set

		for _, td := range []testdata{
			{"Hello, Nick!", greetingTemplate, map[string]interface{}{"greeting": "Hello", "who": "Nick"}},
			{"Greetings, Zach!", greetingTemplate, map[string]interface{}{"greeting": "Greetings", "who": "Zach"}},
			{"Sara!", greetingTemplate, map[string]interface{}{"greeting": "", "who": "Sara"}},
			{"Farewell!", greetingTemplate, map[string]interface{}{"greeting": "Farewell", "who": ""}},
		} {

			mockCtrl := mocking.NewController(t)
			defer mockCtrl.Finish()
			dataResource := NewMockResourceData(mockCtrl)

			dataResource.EXPECT().Get(mocking.Eq("template")).Return(td.Template)
			dataResource.EXPECT().GetOk(mocking.Eq("vars")).Return(td.Vars, true)
			dataResource.EXPECT().Set(mocking.Eq("output"), mocking.Eq(td.Output)).Return(nil)

			// each call to SetID should be passed a unique ID
			dataResource.EXPECT().SetId(mocking.Not(idSet))

			err := execute(dataResource, nil, testGenerateID, defaultTemplate())

			if !gotExpectedErr(t, err, nil) {
				t.Fail()
			}
		}
	})
}

func TestCIDRHost(t *testing.T) {
	t.SkipNow()

	type testdata struct {
		Output, Template string
	}

	template := defaultTemplate()
	vars := map[string]interface{}{
		"cidr0": "10.0.0.0/24",
		"cidr1": "10.0.1.0/24",
		"cidr2": "10.2.2.128/25",
	}

	tmpl := "" +
		"host0: {{cidrhost .cidr0 0}}\n" +
		"host1: {{cidrhost .cidr1 1}}\n" +
		"host2: {{cidrhost .cidr2 -1}}"

	out := "" +
		"host0: 10.0.0.0\n" +
		"host1: 10.0.1.1\n" +
		"host2: 10.2.2.255"

	var (
		err    error
		buffer bytes.Buffer
	)
	template, err = template.Parse(tmpl)
	if err != nil {
		t.Error("it should not error")
		t.Log(err)
		return
	}
	if err := template.Execute(&buffer, vars); err != nil {
		t.Error("it should not error")
		t.Log(err)
		return
	}

	got := string(buffer.Bytes())
	if got != out {
		t.Errorf("it expected output %q but got %q", out, got)
	}
}

func TestGenerateID(t *testing.T) {
	if generateID("hello") != generateID("hello") {
		t.Fail()
	}
	if generateID("") != generateID("") {
		t.Fail()
	}
	if generateID("hello") == generateID("world") {
		t.Fail()
	}
}

// ************************************
// Test Helpers and Data Structures

func gotExpectedErr(t *testing.T, gotErr, expectedErr error) bool {
	t.Helper()
	errMsg := "it should set expected error value\n\nexpected: %v\n\tgot: %v"

	if expectedErr == nil {
		if gotErr != nil {
			t.Logf(errMsg, nil, gotErr.Error())
			return false
		}
		return true
	}

	if gotErr == nil {
		t.Logf(errMsg, expectedErr.Error, nil)
		return false
	}

	if gotErr.Error() != expectedErr.Error() {
		t.Logf(errMsg, gotErr, expectedErr)
		return false
	}

	return true
}

type Set []string

func (set Set) Contains(elem string) bool {
	for _, current := range set {
		if elem == current {
			return true
		}
	}
	return false
}

// Insert returns true if the elem was inserted.
// If the elem will be unique, it will be inserted.
func (set *Set) Insert(elem string) bool {
	if set.Contains(elem) {
		return false
	}
	(*set) = append((*set), elem)
	return true
}

// Matches returns true if the elem is already in the set.
func (set *Set) Matches(x interface{}) bool {
	str, correctType := x.(string)
	if !correctType {
		return false
	}
	return !set.Insert(str)
}

func (set Set) String() string {
	return "some value in the string set"
}

func TestSet(t *testing.T) {
	var set Set

	if set.Matches("hello") {
		t.Fail()
	}
	if !set.Matches("hello") {
		t.Fail()
	}
}
