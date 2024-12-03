// Пакет api реализует HTTP API сервера
package api

import (
	"context"
	"net/http"
	"strconv"
	"testing-platform/core/internal/config"
	"testing-platform/core/internal/model"
	"testing-platform/core/internal/service"

	"github.com/gin-gonic/gin"
)

// Server представляет HTTP сервер приложения
type Server struct {
	cfg     *config.Config    // Конфигурация сервера
	service *service.Service  // Сервисный слой
	router  *gin.Engine      // Маршрутизатор HTTP запросов
	srv     *http.Server     // HTTP сервер
}

// NewServer создает новый экземпляр сервера
func NewServer(cfg *config.Config, svc *service.Service) *Server {
	s := &Server{
		cfg:     cfg,
		service: svc,
		router:  gin.Default(),
	}

	s.setupRoutes()
	s.srv = &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: s.router,
	}

	return s
}

// setupRoutes настраивает маршруты API
func (s *Server) setupRoutes() {
	// Добавляем эндпоинт для healthcheck
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	api := s.router.Group("/api")
	{
		tests := api.Group("/tests")
		{
			tests.POST("/", s.createTest)           // Создание теста
			tests.GET("/:id", s.getTest)           // Получение теста по ID
			tests.POST("/:id/questions", s.addQuestion)    // Добавление вопроса к тесту
			tests.POST("/:id/answers", s.submitAnswer)     // Отправка ответов на тест
			tests.GET("/:id/results", s.getResults)        // Получение результатов теста
		}
	}
}

// Start запускает сервер
func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

// Shutdown выполняет корректное завершение работы сервера
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// createTest обрабатывает создание нового теста
func (s *Server) createTest(c *gin.Context) {
	var test model.Test
	if err := c.ShouldBindJSON(&test); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := s.getUserID(c)
	err := s.service.CreateTest(c.Request.Context(), userID, &test)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, test)
}

// getTest возвращает информацию о тесте по его ID
func (s *Server) getTest(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID теста"})
		return
	}

	userID := s.getUserID(c)
	test, err := s.service.GetTest(c.Request.Context(), userID, id)
	if err != nil {
		if err == service.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, test)
}

// addQuestion добавляет новый вопрос к тесту
func (s *Server) addQuestion(c *gin.Context) {
	testID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID теста"})
		return
	}

	var question model.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question.TestID = testID
	userID := s.getUserID(c)
	err = s.service.AddQuestion(c.Request.Context(), userID, &question)
	if err != nil {
		if err == service.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, question)
}

// submitAnswer обрабатывает отправку ответов на тест
func (s *Server) submitAnswer(c *gin.Context) {
	testID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID теста"})
		return
	}

	var answer model.Answer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := s.getUserID(c)
	answer.UserID = userID
	answer.TestID = testID // Используем testID
	err = s.service.SubmitAnswer(c.Request.Context(), userID, &answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, answer)
}

// getResults возвращает результаты теста
func (s *Server) getResults(c *gin.Context) {
	testID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID теста"})
		return
	}

	userID := s.getUserID(c)
	results, err := s.service.GetResults(c.Request.Context(), testID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if results == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "результаты не найдены"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// getUserID извлекает ID пользователя из заголовка Authorization
// В реальном приложении здесь должна быть проверка JWT токена
func (s *Server) getUserID(c *gin.Context) int64 {
	// Временная реализация для примера
	return 1
}
