package model

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type TodoAIInterface interface {
	GetResponseAPI(c echo.Context, key string, todoAI TodoAI) (openai.ChatCompletionResponse, error)
}

type TodoAI struct {
	Todo string    `json:"todo" form:"todo"`
	Time time.Time `json:"time" form:"time"`
}

type TodoAIModel struct {
	db *gorm.DB
}

func (tm *TodoAIModel) InitTodo(db *gorm.DB) {
	tm.db = db
}

func NewTodoAIModel(db *gorm.DB) TodoAIInterface {
	return &TodoAIModel{
		db: db,
	}
}

func (tm *TodoAIModel) GetResponseAPI(c echo.Context, key string, todo TodoAI) (openai.ChatCompletionResponse, error) {
	content := fmt.Sprintf("Data ini berdasarkan pengisian dari data json. Berikan rekomendasi kegiatan sesuai dengan kriteria dan waktu yang ditentukan. Jika salah satu inputan kosong tetap buatkan rekomendasi dengan inputan yang ada. Berikan jawaban terbaik. `{'kegiatan':%s, 'waktu':%s }`", todo.Todo, todo.Time)
	client := openai.NewClient(key)

	res, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	return res, err
}
