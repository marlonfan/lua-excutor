package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	luaHttp "github.com/cjoudrey/gluahttp"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	luaJson "github.com/glendc/gopher-json"
	"github.com/robfig/cron/v3"
	lua "github.com/yuin/gopher-lua"
	"gorm.io/gorm"
)

//go:embed build
var staticFiles embed.FS

var (
	username     = os.Getenv("HTTP_AUTH_USERNAME")
	password     = os.Getenv("HTTP_AUTH_PASSWORD")
	webhookToken = os.Getenv("HTTP_WEBHOOK_TOKEN")
)

func init() {
	if username == "" {
		username = "admin"
	}
	if password == "" {
		password = "admin"
	}
	if webhookToken == "" {
		webhookToken = "token"
	}
}

type Script struct {
	Name        string `json:"name" gorm:"primary_key"`
	Code        string `json:"code"`
	Schedule    string `json:"schedule"`
	Description string `json:"description"`
	Alias       string `json:"alias"`
}

type KV struct {
	Key        string `json:"key" gorm:"primary_key"`
	Value      string `json:"value"`
	ScriptName string `json:"script_name"`
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("scripts.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	createTable()

	r := gin.Default()

	// 设置HTTP基本认证中间件
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		username: password,
	}))

	authorized.Any("/", func(c *gin.Context) {
		c.FileFromFS("build/", http.FS(staticFiles))
	})

	authorized.Any("/static/*filepath", func(c *gin.Context) {
		c.FileFromFS("build/static/"+c.Param("filepath"), http.FS(staticFiles))
	})

	authorized.POST("/api/scripts", submitScript)
	authorized.Any("/api/scripts/:name/execute", executeScript)
	authorized.GET("/api/scripts", getScripts)
	authorized.POST("/api/scripts/:name/schedule", scheduleScript)
	authorized.PUT("/api/update/:name", updateScript)
	authorized.DELETE("/api/scripts/:name", deleteScript)

	r.Any("/api/scripts/:name/:token/execute", executeScript)

	go runScheduler()

	r.Run(":8080")
}

func createTable() {
	if !db.Migrator().HasTable(&Script{}) {
		db.Migrator().CreateTable(&Script{})
	}
	if !db.Migrator().HasTable(&KV{}) {
		db.Migrator().CreateTable(&KV{})
	}
}

func submitScript(c *gin.Context) {
	var script Script
	if err := c.ShouldBindJSON(&script); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查脚本名称是否已存在
	var existingScript Script
	err := db.Where("name = ?", script.Name).First(&existingScript).Error
	if err == nil {
		// 如果没有错误，说明脚本名称已存在
		c.JSON(http.StatusConflict, gin.H{"error": "Script name already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		// 如果是其他错误，返回内部服务器错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 如果脚本名称不存在，则插入新脚本
	err = db.Create(&script).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Script saved successfully"})
}

func executeScript(c *gin.Context) {
	name := c.Param("name")
	token := c.Param("token")

	if strings.TrimSpace(token) != "" && token != webhookToken {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
		return
	}

	params := c.Request.URL.Query()

	bodyParams := make(map[string]interface{})
	c.ShouldBindJSON(&bodyParams)

	mergedParams := make(map[string]interface{})

	for key, values := range params {
		if len(values) > 0 {
			mergedParams[key] = values[0]
		}
	}

	for key, value := range bodyParams {
		mergedParams[key] = value
	}

	mergedParamsJSON, err := json.Marshal(mergedParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var script Script
	err = db.Where("name = ? OR alias = ?", name, name).First(&script).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Script not found"})
		return
	}

	result, err := runLuaScript(script.Name, script.Code, string(mergedParamsJSON))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func scheduleScript(c *gin.Context) {
	var script Script
	if err := c.ShouldBindJSON(&script); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.Model(&Script{}).Where("name = ? OR alias = ?", script.Name, script.Alias).Update("schedule", script.Schedule).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Script scheduled successfully"})
}

func getScripts(c *gin.Context) {
	var scripts []Script
	err := db.Find(&scripts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scripts)
}

func runScheduler() {
	for {
		var scripts []Script
		err := db.Where("schedule != ?", "").Find(&scripts).Error
		if err != nil {
			log.Println("Error querying scheduled scripts:", err)
			time.Sleep(time.Second)
			continue
		}

		for _, script := range scripts {
			if shouldRun(script.Schedule) {
				go func(code string) {
					_, err := runLuaScript(script.Name, code)
					if err != nil {
						log.Printf("Error running scheduled script: %v\n", err)
					}
				}(script.Code)
			}
		}

		time.Sleep(time.Second)
	}
}

func shouldRun(schedule string) bool {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	sched, err := parser.Parse(schedule)
	if err != nil {
		log.Printf("Error parsing schedule: %v\n", err)
		return false
	}

	now := time.Now()
	next := sched.Next(now.Add(-time.Second))
	return now.Truncate(time.Second).Equal(next)
}

func runLuaScript(name, code string, params ...string) (string, error) {
	L := lua.NewState()
	defer L.Close()

	luaJson.Preload(L)
	L.PreloadModule("http", luaHttp.NewHttpModule(httpClient).Loader)

	if len(params) > 0 {
		L.SetGlobal("params", lua.LString(params[0]))
	}

	L.SetGlobal("script_name", lua.LString(name))
	L.SetGlobal("kv_set", L.NewFunction(kvSet))
	L.SetGlobal("kv_get", L.NewFunction(kvGet))

	err := L.DoString(code)
	if err != nil {
		return "", err
	}

	return "Script executed successfully", nil
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	},
}

func kvSet(L *lua.LState) int {
	key := L.ToString(1)
	value := L.ToString(2)
	scriptName := L.GetGlobal("script_name").String() // 修正获取脚本名称的方式

	kv := KV{
		Key:        key,
		Value:      value,
		ScriptName: scriptName,
	}

	err := db.Model(&KV{}).Where("key = ? AND script_name = ?", key, scriptName).FirstOrCreate(&kv).Debug().Error
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	fmt.Println("kv", kv)

	L.Push(lua.LBool(true))
	return 1
}

func kvGet(L *lua.LState) int {
	key := L.ToString(1)
	scriptName := L.GetGlobal("script_name").String()

	kv := KV{}
	err := db.Model(&KV{}).Where("key = ? AND script_name = ?", key, scriptName).First(&kv).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(kv.Value))
	return 1
}

func updateScript(c *gin.Context) {
	name := c.Param("name")
	var script Script
	if err := c.ShouldBindJSON(&script); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.Model(&Script{}).Where("name = ? OR alias = ?", name, name).Updates(script).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Script updated successfully"})
}

func deleteScript(c *gin.Context) {
	name := c.Param("name")

	err := db.Where("name = ? OR alias = ?", name, name).Delete(&Script{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Script deleted successfully"})
}
