package handlers

import (
	"math/rand"
	"net/http"
	"sync"
	"questionary/backend/data"

	"github.com/gin-gonic/gin"
)

type Quiz struct {
	ID             int      `json:"id"`
	Questions      []gin.H  `json:"questions"`
	CurrentIndex   int      `json:"current_index"`
	Score          int      `json:"score"`
	QuestionCount  int      `json:"question_count"`
}

type QuizHandler struct {
	quizzes map[int]Quiz
	mu      sync.RWMutex
}

func NewQuizHandler() *QuizHandler {
	return &QuizHandler{
		quizzes: make(map[int]Quiz),
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (h *QuizHandler) GetTopics(c *gin.Context) {
	topics := data.GetTopics()
	c.JSON(http.StatusOK, gin.H{"topics": topics})
}

func (h *QuizHandler) GetQuestions(c *gin.Context) {
	topic := c.Query("topic")
	count := parseInt(c.Query("count"))
	if count < 1 || count > 50 {
		count = 10
	}

	var qs []data.Question
	if topic != "" {
		qs = data.GetQuestionsForTopic(topic)
	} else {
		qs = data.GetQuestions()
	}

	if count > len(qs) {
		count = len(qs)
	}
	rand.Shuffle(len(qs), func(i, j int) {
		qs[i], qs[j] = qs[j], qs[i]
	})
	qs = qs[:count]

	out := make([]gin.H, len(qs))
	for i, q := range qs {
		out[i] = gin.H{
			"id":       q.ID,
			"question": q.Question,
			"options":  q.Options,
			"answer":   q.Answer,
			"category": q.Category,
		}
	}

	c.JSON(http.StatusOK, gin.H{"questions": out})
}

func (h *QuizHandler) StartQuiz(c *gin.Context) {
	var req struct {
		QuestionCount int    `json:"questionCount"`
		Topic         string `json:"topic"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.QuestionCount < 1 {
		req.QuestionCount = 10
	}
	if req.QuestionCount > 50 {
		req.QuestionCount = 50
	}

	var qs []data.Question
	if req.Topic != "" {
		qs = data.GetQuestionsForTopic(req.Topic)
	} else {
		qs = data.GetQuestions()
	}

	if req.QuestionCount > len(qs) {
		req.QuestionCount = len(qs)
	}
	rand.Shuffle(len(qs), func(i, j int) {
		qs[i], qs[j] = qs[j], qs[i]
	})
	qs = qs[:req.QuestionCount]

	questions := make([]gin.H, req.QuestionCount)
	for i, q := range qs {
		questions[i] = gin.H{
			"id":       q.ID,
			"question": q.Question,
			"options":  q.Options,
			"answer":   q.Answer,
			"category": q.Category,
		}
	}

	quizID := rand.Intn(9000) + 1000
	quiz := Quiz{
		ID:            quizID,
		Questions:     questions,
		CurrentIndex:  0,
		Score:         0,
		QuestionCount: req.QuestionCount,
	}

	h.mu.Lock()
	h.quizzes[quizID] = quiz
	h.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"quiz_id":   quizID,
		"questions": questions,
		"total":     req.QuestionCount,
	})
}

func (h *QuizHandler) GetQuizQuestion(c *gin.Context) {
	id := parseInt(c.Param("id"))
	h.mu.RLock()
	quiz, exists := h.quizzes[id]
	h.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}
	if quiz.CurrentIndex >= quiz.QuestionCount {
		c.JSON(http.StatusOK, gin.H{"finished": true, "score": quiz.Score, "total": quiz.QuestionCount})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"index":     quiz.CurrentIndex,
		"question":  quiz.Questions[quiz.CurrentIndex],
		"score":     quiz.Score,
		"finished":  false,
	})
}

func (h *QuizHandler) SubmitQuizAnswer(c *gin.Context) {
	id := parseInt(c.Param("id"))
	var req struct {
		Answer int `json:"answer"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	h.mu.Lock()
	quiz, exists := h.quizzes[id]
	h.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}
	if quiz.CurrentIndex >= quiz.QuestionCount {
		c.JSON(http.StatusOK, gin.H{"finished": true, "score": quiz.Score, "total": quiz.QuestionCount})
		return
	}

	q := quiz.Questions[quiz.CurrentIndex]
	correct := req.Answer == q["answer"]

	if correct {
		quiz.Score++
	}

	quiz.CurrentIndex++
	h.mu.Lock()
	h.quizzes[id] = quiz
	h.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"correct":    correct,
		"correct_id": q["answer"],
		"score":      quiz.Score,
		"index":      quiz.CurrentIndex,
	})
}

func parseInt(s string) int {
	i := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			i = i*10 + int(c-'0')
		}
	}
	return i
}
