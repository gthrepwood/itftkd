package main

import (
	"log"
	"os"

	"questionary/backend/handlers"

	"github.com/gin-gonic/gin"
)

func frontendPath() string {
	if _, err := os.Stat("./frontend"); err == nil {
		return "./frontend"
	}
	if _, err := os.Stat("../frontend"); err == nil {
		return "../frontend"
	}
	return "./frontend"
}

func main() {
	r := gin.Default()

	h := handlers.NewQuizHandler()

	api := r.Group("/api")
	{
		api.GET("/topics", h.GetTopics)
		api.GET("/questions", h.GetQuestions)
		api.POST("/quiz/start", h.StartQuiz)
		api.GET("/quiz/:id/question", h.GetQuizQuestion)
		api.POST("/quiz/:id/answer", h.SubmitQuizAnswer)
	}

	fp := frontendPath()
	r.Static("/assets", fp)
	r.GET("/", func(c *gin.Context) {
		c.File(fp + "/index.html")
	})

	log.Println("Server running on http://localhost:9000")
	if err := r.Run(":9000"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
