package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	var err error
	db, err = gorm.Open(sqlite.Open("test_scripts.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open test database: %v", err)
	}
	defer func() {
		os.Remove("./test_scripts.db")
	}()

	createTable()
	m.Run()
}

func TestSubmitScript(t *testing.T) {
	r := gin.Default()
	r.POST("/submit", submitScript)

	script := Script{
		Name: "testScript",
		Code: `
		local http = require("http")
		local json = require("json")

		resp, err = http.get("https://httpbin.org/ip")
		if err == nil then
			print(err)
			ip = json.decode(resp.body).origin
			print(ip)
		end
		
		resp, err = http.get("https://httpbin.org/status/404")
		if err == nil then
			print(resp.status_code)
			print(json.decode(resp.body))
		end
		
		success, err = kv_set("test111", "123123")
		print(success, err)

		data, err = kv_get("test111")
		print(data, err)
		`,
		Schedule:    "* * * * *",
		Description: "A test script",
		Alias:       "test",
	}
	scriptJSON, err := json.Marshal(script)
	if err != nil {
		t.Fatalf("Failed to marshal script: %v", err)
	}

	req, err := http.NewRequest("POST", "/submit", bytes.NewBuffer(scriptJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("Response: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "Script saved successfully"}`, w.Body.String())
}

func TestExecuteScript(t *testing.T) {
	r := gin.Default()
	r.GET("/execute/:name", executeScript)

	req, _ := http.NewRequest("GET", "/execute/testScript", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"result": "Script executed successfully"}`, w.Body.String())
}

func TestScheduleScript(t *testing.T) {
	r := gin.Default()
	r.POST("/schedule", scheduleScript)

	script := Script{
		Name:     "testScript",
		Schedule: "* * * * *",
	}
	scriptJSON, _ := json.Marshal(script)

	req, _ := http.NewRequest("POST", "/schedule", bytes.NewBuffer(scriptJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "Script scheduled successfully"}`, w.Body.String())
}
