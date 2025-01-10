package main

import (
	"encoding/json"
	"math/rand"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func StringToJsonObj(jsonString string) map[string]interface{} {
	var anyJson map[string]interface{}
	customJSON := []byte(jsonString)
	json.Unmarshal(customJSON, &anyJson)
	return anyJson
}

func MapStringInterfaceToString(anyJson map[string]interface{}) string {
	marshalled, _ := json.Marshal(anyJson)
	return string(marshalled)
}

func AnyStringToPayload(jsonString string) string {
	anyJson := StringToJsonObj(jsonString)
	return MapStringInterfaceToString(anyJson)
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// TODO: Remove after implementing proper mocking
func DeleteTestDb() {
	os.Remove("test_.db")
}

func TestCreateTask(t *testing.T) {
	DeleteTestDb()
	router := SetupRouter()

	testErrStr := "Must pass valid token"
	payloadStr := "{}"
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(
		"POST", "/tasks", strings.NewReader(payloadStr)))
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 401, recorder.Code, testErrStr)
	})

	testErrStr = "Empty paylaod not allowed."
	token := MakeDummyUserToken()
	recorder = httptest.NewRecorder()
	req := httptest.NewRequest(
		"POST", "/tasks", strings.NewReader(payloadStr))
	req.Header.Add("Authorization", token)
	router.ServeHTTP(recorder, req)
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 400, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Contains(
			t, respData["error"].(string), "required", testErrStr)
	})

	testErrStr = "Incorrect data type for Title."
	payloadStr = `{"title": 123}`
	recorder = httptest.NewRecorder()
	// TODO: Investigate if there is a client to set token once like django
	req = httptest.NewRequest("POST", "/tasks", strings.NewReader(payloadStr))
	req.Header.Add("Authorization", token)
	router.ServeHTTP(recorder, req)
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 400, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Contains(
			t, respData["error"].(string), "unmarshal", testErrStr)
		assert.Contains(
			t, respData["error"].(string), "Task.title of type string", testErrStr)
	})

	testErrStr = "Incorrect data type for Description."
	payloadStr = `{"description": 123}`
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/tasks", strings.NewReader(payloadStr))
	req.Header.Add("Authorization", token)
	router.ServeHTTP(recorder, req)
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 400, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Contains(
			t, respData["error"].(string), "unmarshal", testErrStr)
		assert.Contains(
			t, respData["error"].(string), "Task.description of type string", testErrStr)
	})

	testErrStr = "Incorrect data type for Description."
	payloadStr = `{"description": 123}`
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/tasks", strings.NewReader(payloadStr))
	req.Header.Add("Authorization", token)
	router.ServeHTTP(recorder, req)
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 400, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Contains(
			t, respData["error"].(string), "unmarshal", testErrStr)
		assert.Contains(
			t, respData["error"].(string), "Task.description of type string", testErrStr)
	})

	testErrStr = "Task saving failed."
	testTitle := RandomString(10)
	testDesc := RandomString(20)
	payloadStr = `{"title": "` + testTitle + `", "description": "` + testDesc + `"}`
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/tasks", strings.NewReader(payloadStr))
	req.Header.Add("Authorization", token)
	router.ServeHTTP(recorder, req)
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 201, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Contains(
			t, respData["title"].(string), testTitle, testErrStr)
		assert.Contains(
			t, respData["description"].(string), testDesc, testErrStr)
		assert.Greater(t, respData["ID"], 0.0, testErrStr)
		assert.Greater(t, respData["CreatedAt"], 0.0, testErrStr)
		assert.Greater(t, respData["UpdatedAt"], 0.0, testErrStr)
	})
	DeleteTestDb()
}

func CreateNTasks(n int) {
	// var tasks = []*Task{}
	db, _ := DBConnection()
	for i := 0; i < n; i++ {
		newTask := Task{Title: RandomString(10), Description: RandomString(20), Status: Pending}
		// tasks := append(tasks, &newTask)
		// TODO: replace with bulk create
		db.Create(&newTask)
	}
}

func TestFetchTask(t *testing.T) {
	DeleteTestDb()
	CreateNTasks(3)
	router := SetupRouter()
	testErrStr := "Task.ID type validation failed."
	token := MakeDummyUserToken()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/abc", nil)
	req.Header.Add("Authorization", token)
	router.ServeHTTP(recorder, req)
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 400, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Equal(
			t, respData["error"].(string), "Invalid Id format", testErrStr)
	})

	testErrStr = "Task.ID Not found."
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/tasks/100", nil)
	req.Header.Add("Authorization", token)
	router.ServeHTTP(recorder, req)
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 404, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Equal(
			t, respData["error"].(string), "Task not found", testErrStr)
	})

	testErrStr = "Valid Task  Not found."
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/tasks/2", nil)
	req.Header.Add("Authorization", token)
	router.ServeHTTP(recorder, req)
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 200, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		respTaskId := respData["task"].(map[string]interface{})["ID"].(float64)
		assert.Equal(t, int(respTaskId), 2, testErrStr)
	})

	DeleteTestDb()
}

func TestDeleteTask(t *testing.T) {
	DeleteTestDb()
	CreateNTasks(3)
	router := SetupRouter()
	testErrStr := "Task.ID type validation failed."
	token := MakeDummyUserToken()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/tasks/abc", nil)
	req.Header.Add("Authorization", token)
	router.ServeHTTP(recorder, req)
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 400, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Equal(
			t, respData["error"].(string), "Invalid Id format", testErrStr)
	})

	testErrStr = "Task.ID Not found."
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("DELETE", "/tasks/100", nil)
	req.Header.Add("Authorization", token)
	router.ServeHTTP(recorder, req)
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 404, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Equal(
			t, respData["error"].(string), "Task not found", testErrStr)
	})

	testErrStr = "Valid Task  Not found."
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("DELETE", "/tasks/2", nil)
	req.Header.Add("Authorization", token)
	router.ServeHTTP(recorder, req)
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 200, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		respTaskId := int(respData["task"].(map[string]interface{})["ID"].(float64))
		assert.Equal(t, respTaskId, 2, testErrStr)
		db, _ := DBConnection()
		var task Task
		result := db.Where("ID = ?", respTaskId).First(&task)
		assert.NotNil(t, result.Error)
	})

	DeleteTestDb()
}

func TestListTasks(t *testing.T) {
	DeleteTestDb()
	router := SetupRouter()
	testErrStr := "Checking list api with no data failed."
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(
		"GET", "/public/tasks", nil))
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 200, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Len(
			t, respData[`tasks`], 0, testErrStr)
	})

	CreateNTasks(3)
	testErrStr = "Checking list api with 3 tasks failed."
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(
		"GET", "/public/tasks", nil))
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 200, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Len(t, respData[`tasks`], 3, testErrStr)
	})

	testErrStr = "Checking pagination page 1 limit 2 failed."
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(
		"GET", "/public/tasks?page=1&limit=2", nil))
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 200, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Len(t, respData[`tasks`], 2, testErrStr)
	})

	testErrStr = "Checking pagination page 2 limit 2 failed."
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(
		"GET", "/public/tasks?page=2&limit=2", nil))
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 200, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Len(t, respData[`tasks`], 1, testErrStr)
	})

	testErrStr = "Checking pagination page 2 limit 10 failed."
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(
		"GET", "/public/tasks?page=2&limit=10", nil))
	t.Run(testErrStr, func(t *testing.T) {
		assert.Equal(t, 200, recorder.Code, testErrStr)
		respData := StringToJsonObj(recorder.Body.String())
		assert.Len(t, respData[`tasks`], 0, testErrStr)
	})

	DeleteTestDb()
}
