package main

import (
	"fmt"
	"net/http"									
	// http.StatusOK        // 200
	// http.StatusBadRequest // 400
	// http.StatusCreated   // 201

	//"github.com/Knetic/govaluate".  // библиотека для вычисления выражений, заданных в виде строки.
	//"github.com/google/uuid"		// создание уникальных индификаторов
	"github.com/labstack/echo/v4"	// основной пакет фреймворка Echo.
	"github.com/labstack/echo/v4/middleware"	//
)


var task string // хранение задачи

type TaskRequest struct{
	Task string `json:"task"`
}

func GetTask(c echo.Context) error{
	if task == ""{
		task = "exemple"
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("Hello, %s", task))
}
func PostTask(c echo.Context) error{
	var req TaskRequest

	if err := c.Bind(&req); err != nil{
		return c.JSON(http.StatusBadRequest, map[string]string{"error":"ivalid request"})
	}

	task = req.Task			

	return c.JSON(http.StatusCreated, map[string]string{"task" : task})
}

func main() {
	e := echo.New()								// Создаём новый сервер Echo.

	e.Use(middleware.CORS())					// Разрешает кросс-доменные запросы (полезно для фронтенда на другом домене).
	e.Use(middleware.Logger())					// Логирует все запросы на сервер (метод, путь, статус, время).

	e.GET("/task", GetTask)						// определение гет ручки
	e.POST("/task", PostTask)					// определение пост ручки

	e.Logger.Fatal(e.Start("localhost:8080"))	// старт сервера
}
