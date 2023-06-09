package requests

import (
	"fmt"
	"testing"
)

func TestGet(T *testing.T) {
	request := New("http://localhost:3000", map[string]string{
		"Content-Type": "application/json",
		"length":       "100",
	})
	err := request.Get("/a3333", map[string]interface{}{
		"name": "张三",
		"age":  19,
		"sex":  true,
	})
	if err != nil {
		T.Errorf("error: %s", err.Error())
		return
	}
	fmt.Println(string(request.ReadData()))
}
func TestPost(T *testing.T) {
	request := New("http://localhost:3000", map[string]string{
		"Content-Type": "application/json",
		"length":       "100",
	})
	err := request.Post("/a3333", map[string]interface{}{
		"name": "张三",
		"age":  19,
		"sex":  true,
	})
	if err != nil {
		T.Errorf("error: %s", err.Error())
		return
	}
	fmt.Println(string(request.ReadData()))
}
