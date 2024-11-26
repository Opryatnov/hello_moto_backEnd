package services

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExternalAPIClient struct {
	APIKey  string
	APIHost string
}

type MotorcycleImage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	ArticleId int                `bson:"articleId" json:"articleId"`
	Image     string             `bson:"image" json:"image"` // Base64 строка изображения
}

// NewExternalAPIClient создает нового клиента для внешнего API
func NewExternalAPIClient(apiKey, apiHost string) *ExternalAPIClient {
	return &ExternalAPIClient{
		APIKey:  apiKey,
		APIHost: apiHost,
	}
}

// FetchMotoImageByArticleID запрашивает изображение мотоцикла по articleID
func (client *ExternalAPIClient) FetchMotoImageByArticleID(articleID string) (MotorcycleImage, error) {
	// Подготовка запроса
	url := fmt.Sprintf("https://motorcycle-specs-database.p.rapidapi.com/article/%s/image/media", articleID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MotorcycleImage{}, err
	}

	// Добавление заголовков
	req.Header.Add("x-rapidapi-key", client.APIKey)
	req.Header.Add("x-rapidapi-host", client.APIHost)

	// Отправка запроса
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return MotorcycleImage{}, err
	}
	defer res.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return MotorcycleImage{}, err
	}

	// Преобразование изображения в Base64
	imageBase64 := base64.StdEncoding.EncodeToString(body)

	// Преобразование articleID в int
	articleIDInt, err := strconv.Atoi(articleID)
	if err != nil {
		return MotorcycleImage{}, fmt.Errorf("invalid articleID: %v", err)
	}

	// Создание объекта MotorcycleImage
	return MotorcycleImage{
		ArticleId: articleIDInt,
		Image:     imageBase64,
	}, nil
}
