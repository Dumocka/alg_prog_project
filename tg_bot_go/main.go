package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Bot struct {
	api    *tgbotapi.BotAPI // Клиент Telegram API 
	states map[int64]string // Хранит состояния пользователей 
	db     *InMemoryDB // База данных в памяти 
}

type Survey struct {
	ID          int64 // Уникальный идентификатор опроса 
	CreatorID   int64 // ID создателя опроса
	Title       string // Название опроса 
	Description string  // Описание опроса 
	Questions   []Question // Список вопросов
	IsActive    bool   // Флаг активности опроса
	CreatedAt   time.Time // Время создания опроса
}

type Question struct {
	ID      int64 // ID вопроса 
	Text    string // Текст вопроса
	Options []string // Список вариантов ответа
}

type Answer struct {
	UserID     int64 // ID пользователя
	SurveyID   int64 // ID опроса
	QuestionID int64 // ID вопроса
	Answer     string // Ответ
	AnsweredAt time.Time // Время ответа
}

type InMemoryDB struct {
	sync.RWMutex // Мьютекс для безопасного доступа
	surveys map[int64]*Survey // Хранилище опросов
	answers []Answer // Хранилище ответов
	lastID  int64 // Последний использованный ID
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		surveys: make(map[int64]*Survey),
		answers: make([]Answer, 0),
	}
}

func (db *InMemoryDB) CreateSurvey(survey *Survey) int64 {
	db.Lock()
	defer db.Unlock()
	
	db.lastID++
	survey.ID = db.lastID
	survey.CreatedAt = time.Now()
	survey.IsActive = true
	db.surveys[survey.ID] = survey
	return survey.ID
}

func (db *InMemoryDB) GetSurvey(id int64) (*Survey, bool) {
	db.RLock()
	defer db.RUnlock()
	
	survey, exists := db.surveys[id]
	return survey, exists
}

func (db *InMemoryDB) ListActiveSurveys() []*Survey {
	db.RLock()
	defer db.RUnlock()
	
	var result []*Survey
	for _, survey := range db.surveys {
		if survey.IsActive {
			result = append(result, survey)
		}
	}
	return result
}

func (db *InMemoryDB) SaveAnswer(answer Answer) {
	db.Lock()
	defer db.Unlock()
	
	answer.AnsweredAt = time.Now()
	db.answers = append(db.answers, answer)
}

func (db *InMemoryDB) GetSurveyAnswers(surveyID int64) []Answer {
	db.RLock()
	defer db.RUnlock()
	
	var result []Answer
	for _, answer := range db.answers {
		if answer.SurveyID == surveyID {
			result = append(result, answer)
		}
	}
	return result
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:    api,
		states: make(map[int64]string),
		db:     NewInMemoryDB(),
	}, nil
}

func (b *Bot) handleStart(message *tgbotapi.Message) error { // Показывает главное меню с кнопками: "Создать опрос", "Мои опросы", "Доступные опросы"
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Создать опрос"),
			tgbotapi.NewKeyboardButton("Мои опросы"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Доступные опросы"),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите действие:")
	msg.ReplyMarkup = keyboard
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) handleCreateSurvey(message *tgbotapi.Message) error { // Начинает процесс создания опроса
	b.states[message.From.ID] = "waiting_survey_title"
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите название опроса:")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) handleSurveyTitle(message *tgbotapi.Message) error { // Сохраняет название опроса и запрашивает описание
	survey := &Survey{
		CreatorID: message.From.ID,
		Title:     message.Text,
		Questions: make([]Question, 0),
	}
	
	surveyID := b.db.CreateSurvey(survey)
	b.states[message.From.ID] = fmt.Sprintf("waiting_survey_description_%d", surveyID)
	
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите описание опроса:")
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) handleSurveyDescription(message *tgbotapi.Message, surveyID int64) error { // Сохраняет описание опроса и запрашивает первый вопрос
	survey, exists := b.db.GetSurvey(surveyID)
	if !exists {
		return fmt.Errorf("survey not found")
	}
	
	survey.Description = message.Text
	b.states[message.From.ID] = fmt.Sprintf("waiting_question_%d", surveyID)
	
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите вопрос:")
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) handleQuestion(message *tgbotapi.Message, surveyID int64) error { // Сохраняет вопрос и запрашивает варианты ответов
	survey, exists := b.db.GetSurvey(surveyID)
	if !exists {
		return fmt.Errorf("survey not found")
	}

	question := Question{
		ID:      int64(len(survey.Questions) + 1),
		Text:    message.Text,
		Options: make([]string, 0),
	}
	
	survey.Questions = append(survey.Questions, question)
	b.states[message.From.ID] = fmt.Sprintf("waiting_options_%d_%d", surveyID, question.ID)
	
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите варианты ответов по одному. Когда закончите, напишите 'Готово':")
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) handleOption(message *tgbotapi.Message, surveyID, questionID int64) error { // Сохраняет варианты ответов
	survey, exists := b.db.GetSurvey(surveyID)
	if !exists {
		return fmt.Errorf("survey not found")
	}

	if strings.ToLower(message.Text) == "готово" {
		if len(survey.Questions[questionID-1].Options) < 2 {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Нужно добавить минимум 2 варианта ответа. Продолжайте добавлять варианты:")
			_, err := b.api.Send(msg)
			return err
		}

		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Добавить ещё вопрос"),
				tgbotapi.NewKeyboardButton("Завершить создание"),
			),
		)

		b.states[message.From.ID] = fmt.Sprintf("confirming_survey_%d", surveyID)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Хотите добавить ещё вопрос или завершить создание опроса?")
		msg.ReplyMarkup = keyboard
		_, err := b.api.Send(msg)
		return err
	}

	survey.Questions[questionID-1].Options = append(survey.Questions[questionID-1].Options, message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Вариант '%s' добавлен. Добавьте ещё вариант или напишите 'Готово':", message.Text))
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) handleConfirmSurvey(message *tgbotapi.Message, surveyID int64) error { // Завершает процесс создания опроса
	if message.Text == "Добавить ещё вопрос" {
		b.states[message.From.ID] = fmt.Sprintf("waiting_question_%d", surveyID)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Введите вопрос:")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		_, err := b.api.Send(msg)
		return err
	}

	delete(b.states, message.From.ID)
	return b.handleStart(message)
}

func (b *Bot) handleListSurveys(message *tgbotapi.Message) error { // Показывает список доступных опросов
	surveys := b.db.ListActiveSurveys()
	
	if len(surveys) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Нет доступных опросов.")
		_, err := b.api.Send(msg)
		return err
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, survey := range surveys {
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(survey.Title, fmt.Sprintf("survey_%d", survey.ID)),
		})
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Доступные опросы:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) handleShowSurvey(callback *tgbotapi.CallbackQuery) error { // Показывает вопросы опроса пользователю
	parts := strings.Split(callback.Data, "_")
	if len(parts) != 2 {
		return fmt.Errorf("invalid callback data")
	}

	surveyID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return err
	}

	survey, exists := b.db.GetSurvey(surveyID)
	if !exists {
		return fmt.Errorf("survey not found")
	}

	if len(survey.Questions) == 0 {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "В этом опросе нет вопросов.")
		_, err = b.api.Send(msg)
		return err
	}

	// Показываем первый вопрос
	question := survey.Questions[0]
	var keyboard [][]tgbotapi.InlineKeyboardButton
	for i, opt := range question.Options {
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(opt, fmt.Sprintf("ans_%d_%d_%d", surveyID, question.ID, i)),
		})
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, 
		fmt.Sprintf("Опрос: %s\n\n%s\n\nВопрос 1: %s", 
			survey.Title, survey.Description, question.Text))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	_, err = b.api.Send(msg)
	return err
}

func (b *Bot) handleAnswer(callback *tgbotapi.CallbackQuery) error { // Сохраняет ответ пользователя и показывает следующий вопрос
	parts := strings.Split(callback.Data, "_")
	if len(parts) != 4 {
		return fmt.Errorf("invalid callback data")
	}

	surveyID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return err
	}

	questionID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return err
	}

	optionIndex, err := strconv.Atoi(parts[3])
	if err != nil {
		return err
	}

	survey, exists := b.db.GetSurvey(surveyID)
	if !exists {
		return fmt.Errorf("survey not found")
	}

	if questionID > int64(len(survey.Questions)) {
		return fmt.Errorf("question not found")
	}

	question := survey.Questions[questionID-1]
	if optionIndex >= len(question.Options) {
		return fmt.Errorf("invalid option index")
	}

	// Сохраняем ответ
	answer := Answer{
		UserID:     callback.From.ID,
		SurveyID:   surveyID,
		QuestionID: questionID,
		Answer:     question.Options[optionIndex],
	}
	b.db.SaveAnswer(answer)

	// Показываем следующий вопрос или завершаем опрос
	nextQuestionID := questionID + 1
	if nextQuestionID > int64(len(survey.Questions)) {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Спасибо за участие в опросе!")
		_, err = b.api.Send(msg)
		return err
	}

	nextQuestion := survey.Questions[nextQuestionID-1]
	var keyboard [][]tgbotapi.InlineKeyboardButton
	for i, opt := range nextQuestion.Options {
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(opt, fmt.Sprintf("ans_%d_%d_%d", surveyID, nextQuestionID, i)),
		})
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, 
		fmt.Sprintf("Вопрос %d: %s", nextQuestionID, nextQuestion.Text))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	_, err = b.api.Send(msg)
	return err
}

func (b *Bot) handleStop(message *tgbotapi.Message) error { // Отменяет текущую операцию и возвращает в главное меню
	// Очищаем состояние пользователя
	delete(b.states, message.From.ID)

	msg := tgbotapi.NewMessage(message.Chat.ID, "Текущая операция отменена. Возвращаемся в главное меню.")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err := b.api.Send(msg)
	if err != nil {
		return err
	}

	// Возвращаемся в главное меню
	return b.handleStart(message)
}

func (b *Bot) handleShowResults(message *tgbotapi.Message) error { // Показывает статистику ответов для опросов, созданных пользователем
	surveys := b.db.ListActiveSurveys()
	
	var userSurveys []*Survey
	for _, survey := range surveys {
		if survey.CreatorID == message.From.ID {
			userSurveys = append(userSurveys, survey)
		}
	}

	if len(userSurveys) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "У вас нет созданных опросов.")
		_, err := b.api.Send(msg)
		return err
	}

	for _, survey := range userSurveys {
		answers := b.db.GetSurveyAnswers(survey.ID)
		
		var response strings.Builder
		response.WriteString(fmt.Sprintf("Результаты опроса '%s':\n\n", survey.Title))

		// Группируем ответы по вопросам
		answersByQuestion := make(map[int64]map[string]int)
		for _, answer := range answers {
			if _, exists := answersByQuestion[answer.QuestionID]; !exists {
				answersByQuestion[answer.QuestionID] = make(map[string]int)
			}
			answersByQuestion[answer.QuestionID][answer.Answer]++
		}

		// Выводим статистику по каждому вопросу
		for _, question := range survey.Questions {
			response.WriteString(fmt.Sprintf("Вопрос: %s\n", question.Text))
			if answers, exists := answersByQuestion[question.ID]; exists {
				for answer, count := range answers {
					response.WriteString(fmt.Sprintf("- %s: %d голосов\n", answer, count))
				}
			} else {
				response.WriteString("Пока нет ответов\n")
			}
			response.WriteString("\n")
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, response.String())
		_, err := b.api.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) Run() error { // Основной цикл бота
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	log.Println("Started receiving updates...")

	for update := range updates {
		log.Printf("Received update: %+v", update)
		var err error

		if update.Message != nil {
			log.Printf("Received message: %s from user %s", update.Message.Text, update.Message.From.UserName)
			// Обработка команд и текстовых сообщений
			switch {
			case update.Message.Command() == "start":
				log.Println("Handling /start command")
				err = b.handleStart(update.Message)

			case update.Message.Command() == "stop":
				log.Println("Handling /stop command")
				err = b.handleStop(update.Message)

			case update.Message.Text == "Создать опрос":
				log.Println("Handling create survey")
				err = b.handleCreateSurvey(update.Message)

			case update.Message.Text == "Мои опросы":
				log.Println("Handling show results")
				err = b.handleShowResults(update.Message)

			case update.Message.Text == "Доступные опросы":
				log.Println("Handling list surveys")
				err = b.handleListSurveys(update.Message)

			default:
				// Обработка состояний
				state, exists := b.states[update.Message.From.ID]
				if !exists {
					log.Printf("No state for user %d, sending unknown command message", update.Message.From.ID)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нет такой команды")
					_, err = b.api.Send(msg)
					continue
				}
				log.Printf("Handling state: %s", state)

				switch {
				case state == "waiting_survey_title":
					err = b.handleSurveyTitle(update.Message)

				case strings.HasPrefix(state, "waiting_survey_description_"):
					parts := strings.Split(state, "_")
					surveyID, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
					err = b.handleSurveyDescription(update.Message, surveyID)

				case strings.HasPrefix(state, "waiting_question_"):
					parts := strings.Split(state, "_")
					surveyID, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
					err = b.handleQuestion(update.Message, surveyID)

				case strings.HasPrefix(state, "waiting_options_"):
					parts := strings.Split(state, "_")
					surveyID, _ := strconv.ParseInt(parts[len(parts)-2], 10, 64)
					questionID, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
					err = b.handleOption(update.Message, surveyID, questionID)

				case strings.HasPrefix(state, "confirming_survey_"):
					parts := strings.Split(state, "_")
					surveyID, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
					err = b.handleConfirmSurvey(update.Message, surveyID)
				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нет такой команды")
					_, err = b.api.Send(msg)
				}
			}
		} else if update.CallbackQuery != nil {
			log.Printf("Received callback query: %s", update.CallbackQuery.Data)
			// Обработка callback-запросов
			switch {
			case strings.HasPrefix(update.CallbackQuery.Data, "survey_"):
				err = b.handleShowSurvey(update.CallbackQuery)

			case strings.HasPrefix(update.CallbackQuery.Data, "ans_"):
				err = b.handleAnswer(update.CallbackQuery)
			default:
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Нет такой команды")
				b.api.Request(callback)
				continue
			}

			if err == nil {
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
				b.api.Request(callback)
			}
		}

		if err != nil {
			log.Printf("Error: %v", err)
			if update.Message != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка. Попробуйте еще раз.")
				b.api.Send(msg)
			}
		}
	}

	return nil
}

func main() { // Инициализация бота
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set in .env file")
	}
	log.Printf("Using token: %s", token)

	bot, err := NewBot(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.api.Self.UserName)
	log.Println("Bot started...")
	
	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
