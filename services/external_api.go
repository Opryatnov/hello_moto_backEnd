package services

import (
	"encoding/base64"
	"encoding/json"
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

type Categories struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type ModelWithCategory struct {
	ArticleID     *int    `json:"articleId,omitempty"`
	ModelID       *int    `json:"modelId,omitempty"`
	PriceName     *string `json:"priceName,omitempty"`
	ModelName     *string `json:"modelName,omitempty"`
	CategoryName  *string `json:"categoryName,omitempty"`
	YearName      *int    `json:"yearName,omitempty"`
	Identificator string  `json:"identificator"`
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

func (client *ExternalAPIClient) FetchCategories() ([]Categories, error) {

	// URL и параметры для запроса
	url := "https://motorcycle-specs-database.p.rapidapi.com/category"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Добавление заголовков
	req.Header.Add("x-rapidapi-key", client.APIKey)
	req.Header.Add("x-rapidapi-host", client.APIHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch categories: %s", res.Status)
	}

	var models []Categories

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("API response categories:", len(body))

	// Пытаемся распарсить ответ как массив объектов MotoModel
	if err := json.Unmarshal(body, &models); err != nil {
		fmt.Println("Error parse categories", err)
		return nil, err
	}

	fmt.Println("categories parsed", models)

	return models, nil
}

func (client *ExternalAPIClient) FetchModelByCategories(category string, makeId string) ([]ModelWithCategory, error) {
	// URL и параметры для запроса
	url := fmt.Sprintf("https://motorcycle-specs-database.p.rapidapi.com/model/make-id/%s/category/%s", makeId, category)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error req 0", err)
		return nil, err
	}

	// Добавление заголовков
	req.Header.Add("x-rapidapi-key", client.APIKey)
	req.Header.Add("x-rapidapi-host", client.APIHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error req 1", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Error req 2", err)
		return nil, fmt.Errorf("failed to fetch ModelWithCategory: %s", res.Status)
	}

	var models []ModelWithCategory

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error body", res.Body)
		return nil, err
	}

	fmt.Println("API response ModelWithCategory:", len(body))

	// Пытаемся распарсить ответ как массив объектов MotoModel
	if err := json.Unmarshal(body, &models); err != nil {
		fmt.Println("Error parse ModelWithCategory", err)
		return nil, err
	}

	for i := range models {
		models[i].Identificator = makeId
	}

	fmt.Println("ModelWithCategory parsed", models)

	return models, nil
}
