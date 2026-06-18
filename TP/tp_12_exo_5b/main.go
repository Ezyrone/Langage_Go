package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

var tasks = map[string]Task{}

const apiKey = "super-secret-key"

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		fmt.Printf("[LOG] %s | %s %s | %s | %v\n",
			start.Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			duration,
		)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-KEY")
		if key == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "En-tête X-API-KEY manquant"})
			c.Abort()
			return
		}
		if key != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Clé API invalide"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func getTasks(c *gin.Context) {
	taskList := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		taskList = append(taskList, t)
	}
	c.JSON(http.StatusOK, taskList)
}

func getTask(c *gin.Context) {
	id := c.Param("id")
	task, exists := tasks[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func createTask(c *gin.Context) {
	var newTask Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTask.ID = uuid.New().String()
	newTask.Done = false
	tasks[newTask.ID] = newTask
	c.JSON(http.StatusCreated, newTask)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	task, exists := tasks[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée"})
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Done        *bool   `json:"done"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Done != nil {
		task.Done = *input.Done
	}

	tasks[id] = task
	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")
	if _, exists := tasks[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée"})
		return
	}
	delete(tasks, id)
	c.Status(http.StatusNoContent)
}

func main() {
	tasks["1"] = Task{ID: "1", Title: "Apprendre Go", Description: "Suivre le cours M2", Done: false}
	tasks["2"] = Task{ID: "2", Title: "Faire le TP Gin", Description: "Implémenter l'API REST", Done: false}
	tasks["3"] = Task{ID: "3", Title: "Réviser les goroutines", Description: "Revoir channels et select", Done: true}

	r := gin.Default()
	r.Use(LoggerMiddleware())

	api := r.Group("/api/v1")
	{
		api.GET("/tasks", getTasks)
		api.GET("/tasks/:id", getTask)

		protected := api.Group("/")
		protected.Use(AuthMiddleware())
		{
			protected.POST("/tasks", createTask)
			protected.PUT("/tasks/:id", updateTask)
			protected.DELETE("/tasks/:id", deleteTask)
		}
	}

	fmt.Println("Serveur Gin démarré sur http://localhost:8080")
	fmt.Println("Endpoints :")
	fmt.Println("  GET    /api/v1/tasks       - Liste toutes les tâches")
	fmt.Println("  GET    /api/v1/tasks/:id   - Récupère une tâche")
	fmt.Println("  POST   /api/v1/tasks       - Crée une tâche (X-API-KEY requis)")
	fmt.Println("  PUT    /api/v1/tasks/:id   - Met à jour une tâche (X-API-KEY requis)")
	fmt.Println("  DELETE /api/v1/tasks/:id   - Supprime une tâche (X-API-KEY requis)")

	r.Run(":8080")
}
