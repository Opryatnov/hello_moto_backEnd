package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-mongodb-app/services"
	"image"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	// Подключаем файл с константами

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	RapidAPIKey = "756a143d30msh385f180af869e40p1a22fajsnfdc87cc5b091"
	// RapidAPIKey  = "27e6cce226msha2b9adeeaf9541dp147517jsn4b00793fb267" - протух
	RapidAPIHost = "motorcycle-specs-database.p.rapidapi.com"
)

type MotoBrand struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type MotoModel struct {
	// ID        primitive.ObjectID `bson:"_id,omitempy" json:"-"`
	Moto_ID   int    `bson:"id" json:"id"`
	Name      string `json:"name"`
	BrandID   string `json:"brandid"`
	BrandName string `json:"brandName"`
}

type MotorcycleDetails struct {
	Make              string  `json:"make"`
	Model             string  `json:"model"`
	Year              string  `json:"year"`
	Type              *string `json:"type,omitempty"`
	Displacement      *string `json:"displacement,omitempty"`
	Engine            *string `json:"engine,omitempty"`
	Power             *string `json:"power,omitempty"`
	Torque            *string `json:"torque,omitempty"`
	Compression       *string `json:"compression,omitempty"`
	BoreStroke        *string `json:"bore_stroke,omitempty"`
	ValvesPerCylinder *string `json:"valves_per_cylinder,omitempty"`
	FuelSystem        *string `json:"fuel_system,omitempty"`
	FuelControl       *string `json:"fuel_control,omitempty"`
	Ignition          *string `json:"ignition,omitempty"`
	Lubrication       *string `json:"lubrication,omitempty"`
	Cooling           *string `json:"cooling,omitempty"`
	Gearbox           *string `json:"gearbox,omitempty"`
	Transmission      *string `json:"transmission,omitempty"`
	Clutch            *string `json:"clutch,omitempty"`
	Frame             *string `json:"frame,omitempty"`
	FrontSuspension   *string `json:"front_suspension,omitempty"`
	FrontWheelTravel  *string `json:"front_wheel_travel,omitempty"`
	RearSuspension    *string `json:"rear_suspension,omitempty"`
	RearWheelTravel   *string `json:"rear_wheel_travel,omitempty"`
	FrontTire         *string `json:"front_tire,omitempty"`
	RearTire          *string `json:"rear_tire,omitempty"`
	FrontBrakes       *string `json:"front_brakes,omitempty"`
	RearBrakes        *string `json:"rear_brakes,omitempty"`
	TotalWeight       *string `json:"total_weight,omitempty"`
	SeatHeight        *string `json:"seat_height,omitempty"`
	TotalHeight       *string `json:"total_height,omitempty"`
	TotalLength       *string `json:"total_length,omitempty"`
	TotalWidth        *string `json:"total_width,omitempty"`
	GroundClearance   *string `json:"ground_clearance,omitempty"`
	Wheelbase         *string `json:"wheelbase,omitempty"`
	FuelCapacity      *string `json:"fuel_capacity,omitempty"`
	Starter           *string `json:"starter,omitempty"`
	DryWeight         *string `json:"dry_weight,omitempty"`
	MakeID            string  `json:"makeID"`
	Identificator     string  `json:"identificator"`
	ID                *string `json:"id,omitempty"`
}

type MotorcycleSpecification struct {
	ArticleCompleteInfo           *ArticleCompleteInfo           `json:"articleCompleteInfo,omitempty"`
	EngineAndTransmission         *EngineAndTransmission         `json:"engineAndTransmission,omitempty"`
	ChassisSuspensionBrakesWheels *ChassisSuspensionBrakesWheels `json:"chassisSuspensionBrakesAndWheels,omitempty"`
	PhysicalMeasuresCapacities    *PhysicalMeasuresCapacities    `json:"physicalMeasuresAndCapacities,omitempty"`
	OtherSpecifications           *OtherSpecifications           `json:"otherSpecifications,omitempty"`
}

type ArticleCompleteInfo struct {
	ArticleID    *int    `json:"articleID,omitempty"`
	MakeName     *string `json:"makeName,omitempty"`
	ModelName    *string `json:"modelName,omitempty"`
	CategoryName *string `json:"categoryName,omitempty"`
	YearName     *int    `json:"yearName,omitempty"`
}

type EngineAndTransmission struct {
	DisplacementName           *string `json:"displacementName,omitempty"`
	EngineTypeName             *string `json:"engineTypeName,omitempty"`
	EngineDetailsName          *string `json:"engineDetailsName,omitempty"`
	PowerName                  *string `json:"powerName,omitempty"`
	TorqueName                 *string `json:"torqueName,omitempty"`
	CompressionName            *string `json:"compressionName,omitempty"`
	BoreXStrokeName            *string `json:"boreXStrokeName,omitempty"`
	ValvesPerCylinderName      *string `json:"valvesPerCylinderName,omitempty"`
	FuelSystemName             *string `json:"fuelSystemName,omitempty"`
	LubricationSystemName      *string `json:"lubricationSystemName,omitempty"`
	CoolingSystemName          *string `json:"coolingSystemName,omitempty"`
	GearboxName                *string `json:"gearboxName,omitempty"`
	TransmissionFinalDriveName *string `json:"transmissionTypeFinalDriveName,omitempty"`
	ClutchName                 *string `json:"clutchName,omitempty"`
	DrivelineName              *string `json:"drivelineName,omitempty"`
	ExhaustSystemName          *string `json:"exhaustSystemName,omitempty"`
}

type ChassisSuspensionBrakesWheels struct {
	FrameTypeName           *string `json:"frameTypeName,omitempty"`
	FrontBrakesName         *string `json:"frontBrakesName,omitempty"`
	FrontBrakesDiameterName *string `json:"frontBrakesDiameterName,omitempty"`
	FrontSuspensionName     *string `json:"frontSuspensionName,omitempty"`
	FrontTyreName           *string `json:"frontTyreName,omitempty"`
	FrontWheelTravelName    *string `json:"frontWheelTravelName,omitempty"`
	RakeName                *string `json:"rakeName,omitempty"`
	RearBrakesName          *string `json:"rearBrakesName,omitempty"`
	RearBrakesDiameterName  *string `json:"rearBrakesDiameterName,omitempty"`
	RearSuspensionName      *string `json:"rearSuspensionName,omitempty"`
	RearTyreName            *string `json:"rearTyreName,omitempty"`
	RearWheelTravelName     *string `json:"rearWheelTravelName,omitempty"`
	TrailName               *string `json:"trailName,omitempty"`
}

type PhysicalMeasuresCapacities struct {
	DryWeightName           *string `json:"dryWeightName,omitempty"`
	FuelCapacityName        *string `json:"fuelCapacityName,omitempty"`
	OverallHeightName       *string `json:"overallHeightName,omitempty"`
	OverallLengthName       *string `json:"overallLengthName,omitempty"`
	OverallWidthName        *string `json:"overallWidthName,omitempty"`
	PowerWeightRatioName    *string `json:"powerWeightRatioName,omitempty"`
	ReserveFuelCapacityName *string `json:"reserveFuelCapacityName,omitempty"`
	SeatHeightName          *string `json:"seatHeightName,omitempty"`
}

type OtherSpecifications struct {
	ColorOptionsName    *string `json:"colorOptionsName,omitempty"`
	CommentsName        *string `json:"commentsName,omitempty"`
	FactoryWarrantyName *string `json:"factoryWarrantyName,omitempty"`
	StarterName         *string `json:"starterName,omitempty"`
}

type MotorcycleModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	ModelId   int                `bson:"modelId" json:"modelId"`
	ModelName string             `bson:"modelName" json:"modelName"`
	YearName  int                `bson:"yearName" json:"yearName"`
	ArticleId int                `bson:"articleId" json:"articleId"`
	MakeID    string             `bson:"_makeid,omitempty" json:"-"`
}

type ScreenOpenCount struct {
	OpenBrandsScreenCount    int `json:"openBrandsScreenCount"`
	OpenModelsScreenCount    int `json:"openModelsScreenCount"`
	OpenModelListScreenCount int `json:"openModelListScreenCount"`
}

var collection *mongo.Collection
var client *mongo.Client

func main() {
	// Устанавливаем подключение к MongoDB
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://93.183.81.91:27017"))

	if err != nil {
		log.Fatal("Failed to create MongoDB client:", err)
	}

	// Контекст с таймаутом на подключение
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключаемся к MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Проверка подключения
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Создаем клиента для внешнего API
	externalAPIClient := services.NewExternalAPIClient(RapidAPIKey, RapidAPIHost)

	// Создаем маршруты
	http.HandleFunc("/brands", getMotorcycleBrands)
	http.HandleFunc("/models", getModelsByBrand)
	http.HandleFunc("/get-motorcycle-details", getMotorcycleDetails)
	http.HandleFunc("/save-and-fetch-motorcycle-details", saveAndFetchMotorcycleDetails)
	http.HandleFunc("/motorcycles-details", getmotorcyclesSpecifications)
	http.HandleFunc("/screen-open-count", getScreenOpenCount)

	// Обработчик для получения изображения
	http.HandleFunc("/motorcycle-image", func(w http.ResponseWriter, r *http.Request) {
		getImageByArticleID(w, r, externalAPIClient)
	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		getMotorcycleCategories(w, r, externalAPIClient)
	})

	http.HandleFunc("/models-by-categories", func(w http.ResponseWriter, r *http.Request) {
		getModelsByCategories(w, r, externalAPIClient)
	})

	// Запуск HTTPS-сервера
	port := 8181

	fmt.Printf("Server is running on HTTP port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// Получение количества открытий экрана перед показом рекламы
func getScreenOpenCount(w http.ResponseWriter, r *http.Request) {
	// Пример данных (вы можете получить их из базы данных или другой логики)

	fmt.Println("getScreenOpenCount method is called")

	response := ScreenOpenCount{
		OpenBrandsScreenCount:    9,
		OpenModelsScreenCount:    11,
		OpenModelListScreenCount: 15,
	}

	fmt.Println("getScreenOpenCount returned model", response)

	// Устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодируем данные в JSON и отправляем клиенту
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Получение списка марок
func getMotorcycleBrands(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method is called")
	var brands []MotoBrand

	// Подключение к коллекции MongoDB
	collection := client.Database("moto").Collection("brands")

	// Попытка получить данные из коллекции
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		http.Error(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer cur.Close(context.TODO())

	// Чтение данных из базы
	for cur.Next(context.TODO()) {
		var brand MotoBrand
		err := cur.Decode(&brand)
		if err != nil {
			http.Error(w, "Error decoding document: "+err.Error(), http.StatusInternalServerError)
			return
		}

		brands = append(brands, brand)
	}

	// Если данные найдены в базе, возвращаем их
	if len(brands) > 0 {
		fmt.Println("Brands fetched from database:")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(brands)
		return
	}

	// Если данных нет в базе, читаем их из файла JSON
	fmt.Println("No brands found in database, loading from JSON")
	filePath := "data/brands.json"
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening JSON file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&brands)
	if err != nil {
		http.Error(w, "Error decoding JSON file: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("Error decoding JSON file:")
		return
	}

	// Запись данных в базу
	for _, brand := range brands {
		_, err := collection.InsertOne(context.TODO(), brand)
		if err != nil {
			http.Error(w, "Error inserting data into database: "+err.Error(), http.StatusInternalServerError)
			fmt.Println("Error inserting data into database:")
			return
		}
	}

	fmt.Println("Brands loaded from JSON and inserted into database:")

	// Возвращаем данные клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(brands)
}

// Получение моделей по марке
func getModelsByBrand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type RequestBody struct {
		BrandName string `json:"brandName"`
		ID        string `json:"id"`
	}

	var requestBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed body: BrandName=%s, ID=%s\n", requestBody.BrandName, requestBody.ID)

	// Use the parsed parameters
	brandName := requestBody.BrandName
	brandID := requestBody.ID

	// Работа с MongoDB
	modelsCollection := client.Database("moto").Collection("models")
	filter := bson.M{"brandid": brandID}
	cursor, err := modelsCollection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var models []MotoModel
	for cursor.Next(context.Background()) {
		var model MotoModel
		if err := cursor.Decode(&model); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		models = append(models, model)
	}

	// Если модели найдены в базе данных, возвращаем их
	if len(models) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models)
		fmt.Println("Found models in DB", len(models))
		return
	}

	// Если моделей нет, запрашиваем их через внешний API
	fmt.Println("Models not found in DB, fetching from external API...")
	externalModels, err := fetchModelsFromAPI(brandID, brandName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error from api", err)
		return
	}

	// Сохраняем новые модели в базу данных
	for _, model := range externalModels {
		model.BrandID = brandID
		model.BrandName = brandName
		_, err = modelsCollection.InsertOne(context.Background(), model)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("Error saved models from api", err)
			return
		}
	}
	fmt.Println("Models saved in DB from external API")

	// Возвращаем только что добавленные модели
	cursor, err = modelsCollection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error только что добавленные модели", err)
		return
	}
	defer cursor.Close(context.Background())

	models = []MotoModel{}
	for cursor.Next(context.Background()) {
		var model MotoModel
		if err := cursor.Decode(&model); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("Error только что добавленные модели", len(models))
			return
		}
		models = append(models, model)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}

// Функция для получения моделей из внешнего API
func fetchModelsFromAPI(brandID, brandName string) ([]MotoModel, error) {
	url := fmt.Sprintf("https://motorcycle-specs-database.p.rapidapi.com/model/make-id/%s", brandID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-key", RapidAPIKey)
	req.Header.Add("x-rapidapi-host", RapidAPIHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch models: %s", res.Status)
	}

	var models []MotoModel

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("API response models:", len(body))

	// Пытаемся распарсить ответ как массив объектов MotoModel
	if err := json.Unmarshal(body, &models); err != nil {
		fmt.Println("Error parse models", err)
		return nil, err
	}

	// Set the brandName for each model, using the passed brandName parameter
	for i := range models {
		models[i].BrandName = brandName
		models[i].BrandID = brandID
	}

	fmt.Println("models parsed and modify", models)

	return models, nil
}

func getMotorcycleDetails(w http.ResponseWriter, r *http.Request) {

	type RequestBody struct {
		TempIdentifier string `json:"tempIdentifier"`
	}

	log.Printf("tempIdentifier:")

	var requestBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	collection := client.Database("moto").Collection("otherMotorcycleDetails")
	filter := bson.M{"identificator": requestBody.TempIdentifier}

	var results []MotorcycleDetails
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &results); err != nil {
		http.Error(w, "Error decoding results", http.StatusInternalServerError)
		return
	}

	// Возвращаем пустой массив, если ничего не найдено
	if len(results) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func saveAndFetchMotorcycleDetails(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Models []MotorcycleDetails `json:"models"`
	}

	var requestBody RequestBody

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		fmt.Println("Error decoding request body:", err)
		return
	}

	// Проверяем, что модели переданы
	if len(requestBody.Models) == 0 {
		http.Error(w, "No models provided", http.StatusBadRequest)
		fmt.Println("No models provided")
		return
	}

	// Логируем данные для отладки
	for _, m := range requestBody.Models {
		fmt.Printf("Make: %s, Model: %s, Year: %s, Identificator: %s, MakeID: %s\n",
			m.Make, m.Model, m.Year, m.Identificator, m.MakeID)
	}

	collection := client.Database("moto").Collection("otherMotorcycleDetails")

	// Подготавливаем данные для вставки
	docs := make([]interface{}, len(requestBody.Models))
	for i, motorcycle := range requestBody.Models {
		docs[i] = motorcycle
	}

	// Логируем документы перед сохранением
	for _, doc := range docs {
		fmt.Printf("Document to save: %+v\n", doc)
	}

	// Сохраняем данные в базу
	_, err := collection.InsertMany(context.Background(), docs)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Database error during insertion:", err)
		return
	}

	// Выполняем поиск по идентификатору первого объекта
	tempIdentifier := requestBody.Models[0].Identificator
	filter := bson.M{"identificator": tempIdentifier}

	var results []MotorcycleDetails
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Database error during search", http.StatusInternalServerError)
		fmt.Println("Database search error:", err)
		return
	}
	defer cursor.Close(context.Background())

	// Декодируем результаты поиска
	if err := cursor.All(context.Background(), &results); err != nil {
		http.Error(w, "Error decoding results", http.StatusInternalServerError)
		fmt.Println("Error decoding results:", err)
		return
	}

	// Отправляем результаты клиенту
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		fmt.Println("Error encoding response:", err)
	}
}

func getValue(v *string) string {
	if v == nil {
		return "null"
	}
	return *v
}

func EscapeMongoRegex(input string) string {
	return regexp.QuoteMeta(input)
}

func getmotorcyclesSpecifications(w http.ResponseWriter, r *http.Request) {

	if client == nil {
		http.Error(w, "MongoDB connection is not initialized", http.StatusInternalServerError)
		return
	}

	// Получаем параметры из запроса
	makename := r.URL.Query().Get("makeName")
	modelname := r.URL.Query().Get("modelName")

	// Декодируем параметры
	decodedMakeName, err := url.QueryUnescape(makename)
	if err != nil {
		log.Printf("Error decoding makeName: %v", err)
		decodedMakeName = makename // используем оригинал, если ошибка
	}

	decodedModelName, err := url.QueryUnescape(modelname)
	if err != nil {
		log.Printf("Error decoding modelName: %v", err)
		decodedModelName = modelname // используем оригинал, если ошибка
	}

	// Экранируем modelname для безопасного использования в MongoDB фильтре
	escapedModelName := EscapeMongoRegex(decodedModelName)

	// Формируем фильтр для MongoDB
	filter := bson.M{
		"articlecompleteinfo.makename":  bson.M{"$regex": decodedMakeName, "$options": "i"},
		"articlecompleteinfo.modelname": bson.M{"$regex": "^" + escapedModelName + "$", "$options": "i"},
	}

	log.Printf("Final MongoDB Filter: %+v", filter)

	// Указываем коллекцию
	motorcyclesCollection := client.Database("moto").Collection("motorcyclesDetails")

	// Функция для поиска в базе данных
	findInDatabase := func() ([]MotorcycleSpecification, error) {
		cursor, err := motorcyclesCollection.Find(context.Background(), filter)
		if err != nil {
			return nil, fmt.Errorf("failed to query database: %w", err)
		}
		defer cursor.Close(context.Background())

		var motorcycleSpecifications []MotorcycleSpecification
		for cursor.Next(context.Background()) {
			var spec MotorcycleSpecification
			if err := cursor.Decode(&spec); err != nil {
				return nil, fmt.Errorf("failed to decode database result: %w", err)
			}
			motorcycleSpecifications = append(motorcycleSpecifications, spec)
		}
		return motorcycleSpecifications, nil
	}

	// 1. Ищем в базе данных
	motorcycleSpecifications, err := findInDatabase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Если модели найдены, возвращаем их
	if len(motorcycleSpecifications) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(motorcycleSpecifications)
		return
	}

	// 3. Если ничего не найдено, вызываем внешний API
	externalModels, err := fetchMotoModelsFromAPI(makename, modelname)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch data from API: %v", err), http.StatusInternalServerError)
		return
	}

	// 4. Сохраняем результат внешнего API в базу

	log.Printf("Сохраняем результат внешнего API в базу -----------------: %s", externalModels)

	for _, model := range externalModels {
		// Преобразуем структуру в BSON перед сохранением
		bsonData, err := bson.Marshal(model)
		if err != nil {
			log.Printf("failed to marshal model: %v", err)
			http.Error(w, fmt.Sprintf("failed to marshal model: %v", err), http.StatusInternalServerError)
			return
		}

		_, err = motorcyclesCollection.InsertOne(context.Background(), bsonData)
		if err != nil {
			log.Printf("failed to save model to database: %v", err)
			http.Error(w, fmt.Sprintf("failed to save model to database: %v", err), http.StatusInternalServerError)
			return
		}
	}

	log.Printf("Снова ищем в базе данных после вставки -----------------: %s", motorcycleSpecifications)

	// 5. Снова ищем в базе данных после вставки
	motorcycleSpecifications, err = findInDatabase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Показываем что нашли -----------------: %s", motorcycleSpecifications)

	// 6. Возвращаем данные, если они найдены, или пустой массив в противном случае
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(motorcycleSpecifications)
	log.Printf("Возвращаем данные, если они найдены, или пустой массив в противном случае -----------------: %s", motorcycleSpecifications)
}
func escapeRegexString(input string) string {

	log.Printf("input -----------------: %s", input)
	// Экранируем специальные символы, такие как '(', ')', '#', и любые другие, если нужно
	escaped := strings.ReplaceAll(input, `\`, `\\`)  // Экранируем обратную косую черту
	escaped = strings.ReplaceAll(escaped, `#`, `\#`) // Экранируем символ '#'
	escaped = strings.ReplaceAll(escaped, `(`, `\(`) // Экранируем открывающую скобку
	escaped = strings.ReplaceAll(escaped, `)`, `\)`) // Экранируем закрывающую скобку

	// Можно добавить другие символы, если они могут вызвать проблемы, например, '.', '[', ']', и т.д.

	return escaped
}

func fetchMotoModelsFromAPI(make, model string) ([]MotorcycleSpecification, error) {
	// Формирование URL
	escapedMake := url.PathEscape(make)
	escapedModel := strings.ReplaceAll(model, " ", "%20")

	url := fmt.Sprintf("https://motorcycle-specs-database.p.rapidapi.com/make/%s/model/%s", escapedMake, escapedModel)

	fmt.Println("Generated URL:", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Add("x-rapidapi-key", RapidAPIKey)
	req.Header.Add("x-rapidapi-host", RapidAPIHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch models: %s", res.Status)
	}

	var responseBody []map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Парсинг моделей
	var models []MotorcycleSpecification
	for _, item := range responseBody {
		spec, err := parseMotorcycleSpecification(item)
		if err != nil {
			fmt.Println("Error parsing motorcycle specification:", err)
			continue
		}
		models = append(models, *spec)
	}

	return models, nil
}

func parseMotorcycleSpecification(data map[string]interface{}) (*MotorcycleSpecification, error) {
	spec := &MotorcycleSpecification{}

	// Парсинг ArticleCompleteInfo
	if articleData, ok := data["articleCompleteInfo"].(map[string]interface{}); ok {
		var article ArticleCompleteInfo
		if err := mapToStruct(articleData, &article); err != nil {
			fmt.Println("Error parsing ArticleCompleteInfo:", err)
		} else {
			spec.ArticleCompleteInfo = &article
		}
	}

	// Парсинг EngineAndTransmission
	if engineData, ok := data["engineAndTransmission"].(map[string]interface{}); ok {
		var engine EngineAndTransmission
		if err := mapToStruct(engineData, &engine); err != nil {
			fmt.Println("Error parsing EngineAndTransmission:", err)
		} else {
			spec.EngineAndTransmission = &engine
		}
	}

	// Парсинг ChassisSuspensionBrakesWheels
	if chassisData, ok := data["chassisSuspensionBrakesAndWheels"].(map[string]interface{}); ok {
		var chassis ChassisSuspensionBrakesWheels
		if err := mapToStruct(chassisData, &chassis); err != nil {
			fmt.Println("Error parsing ChassisSuspensionBrakesWheels:", err)
		} else {
			spec.ChassisSuspensionBrakesWheels = &chassis
		}
	}

	// Парсинг PhysicalMeasuresCapacities
	if measuresData, ok := data["physicalMeasuresAndCapacities"].(map[string]interface{}); ok {
		var measures PhysicalMeasuresCapacities
		if err := mapToStruct(measuresData, &measures); err != nil {
			fmt.Println("Error parsing PhysicalMeasuresCapacities:", err)
		} else {
			spec.PhysicalMeasuresCapacities = &measures
		}
	}

	// Парсинг OtherSpecifications
	if otherData, ok := data["otherSpecifications"].(map[string]interface{}); ok {
		var other OtherSpecifications
		if err := mapToStruct(otherData, &other); err != nil {
			fmt.Println("Error parsing OtherSpecifications, setting to nil:", err)
			spec.OtherSpecifications = nil // Зануляем, если не удалось распарсить
		} else {
			spec.OtherSpecifications = &other
		}
	} else {
		// Если otherSpecifications - массив или отсутствует
		spec.OtherSpecifications = nil
	}

	return spec, nil
}

// Универсальная функция для маппинга данных в структуру
func mapToStruct(input map[string]interface{}, output interface{}) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input map: %w", err)
	}
	if err := json.Unmarshal(jsonData, output); err != nil {
		return fmt.Errorf("failed to unmarshal into struct: %w", err)
	}
	return nil
}

func getImageByArticleID(w http.ResponseWriter, r *http.Request, apiClient *services.ExternalAPIClient) {
	if client == nil {
		http.Error(w, "MongoDB connection is not initialized", http.StatusInternalServerError)
		return
	}

	// Получаем параметр articleId из запроса
	articleIDStr := r.URL.Query().Get("articleId")
	if articleIDStr == "" {
		http.Error(w, "Required parameter 'articleId' is missing", http.StatusBadRequest)
		fmt.Println("article parameter is missing -------")
		return
	}

	fmt.Println("articleIDStr:", articleIDStr)

	// Преобразуем articleId из строки в число
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		http.Error(w, "Invalid articleId parameter", http.StatusBadRequest)
		fmt.Println("article converted to int error -------")
		return
	}
	fmt.Println("article converted to int:", articleID)

	// Конструируем фильтр для поиска по articleId
	filter := bson.M{"articleId": articleID}

	// Определяем коллекцию
	motorcyclesCollection := client.Database("moto").Collection("motorcyclesImages")

	// Ищем запись в базе данных
	var motoImage services.MotorcycleImage
	err = motorcyclesCollection.FindOne(context.Background(), filter).Decode(&motoImage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("article Если запись не найдена, вызываем внешний сервис для получения изображения", articleID)
			// Если запись не найдена, вызываем внешний сервис для получения изображения
			motoImage, err = apiClient.FetchMotoImageByArticleID(articleIDStr)
			if err != nil {
				http.Error(w, "Failed to fetch image from external API", http.StatusInternalServerError)
				fmt.Println("Failed to fetch image from external API")
				return
			}

			// // Проверяем валидность поля Image
			// if !isValidBase64Image(motoImage.Image) {
			// 	http.Error(w, "Invalid image data from external API", http.StatusInternalServerError)
			// 	fmt.Println("Invalid image data from external API")
			// 	return
			// }

			// Сохраняем изображение в базу данных
			_, err := motorcyclesCollection.InsertOne(context.Background(), motoImage)
			if err != nil {
				http.Error(w, "Failed to save image to the database", http.StatusInternalServerError)
				fmt.Println("Failed to save image to the database")
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Возвращаем модель с изображением в формате JSON
	fmt.Println("модель с изображением", motoImage)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(motoImage)
}

// Проверяет, является ли строка Image валидным изображением
func isValidBase64Image(base64Str string) bool {
	// Убираем пробелы и новые строки
	base64Str = strings.TrimSpace(base64Str)

	// Минимальная длина Base64 строки (например, пустые или короткие строки сразу исключаем)
	if len(base64Str) < 20 {
		fmt.Println("Base64 string too short to be valid:", base64Str)
		return false
	}

	// Декодируем строку Base64
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		// Пробуем URL-safe декодирование, если стандартное не подходит

		fmt.Println("Base64:", base64Str)
		data, err = base64.URLEncoding.DecodeString(base64Str)
		if err != nil {
			fmt.Println("Base64 decode error:", err)
			return false
		}
	}

	// Пробуем определить формат изображения
	_, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		fmt.Println("Image decode error:", err)
		return false
	}

	// Проверяем, является ли формат поддерживаемым
	fmt.Println("Image format detected:", format)
	if format == "" {
		return false
	}

	return true
}

func getMotorcycleCategories(w http.ResponseWriter, r *http.Request, apiClient *services.ExternalAPIClient) {
	if client == nil {
		http.Error(w, "MongoDB connection is not initialized", http.StatusInternalServerError)
		fmt.Println("Error MongoDB NOT connected")
		return
	}

	// Работа с MongoDB
	fmt.Println("MongoDB is connected")
	modelsCollection := client.Database("moto").Collection("categories")

	// Используем пустой фильтр для выборки всех документов
	cursor, err := modelsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var models []services.Categories
	for cursor.Next(context.Background()) {
		var model services.Categories
		if err := cursor.Decode(&model); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		models = append(models, model)
	}

	// Если модели найдены в базе данных, возвращаем их
	if len(models) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models)
		fmt.Println("Found models in DB", len(models))
		return
	}

	// Если моделей нет, запрашиваем их через внешний API
	fmt.Println("Models not found in DB, fetching from external API...")
	externalModels, err := apiClient.FetchCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error from api", err)
		return
	}

	// Сохраняем новые модели в базу данных
	for _, model := range externalModels {
		_, err = modelsCollection.InsertOne(context.Background(), model)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("Error saved models from api", err)
			return
		}
	}
	fmt.Println("Models saved in DB from external API")

	// Возвращаем только что добавленные модели
	cursor, err = modelsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error только что добавленные модели", err)
		return
	}
	defer cursor.Close(context.Background())

	models = []services.Categories{}
	for cursor.Next(context.Background()) {
		var model services.Categories
		if err := cursor.Decode(&model); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("Error только что добавленные модели", len(models))
			return
		}
		models = append(models, model)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}

func getModelsByCategories(w http.ResponseWriter, r *http.Request, apiClient *services.ExternalAPIClient) {

	if client == nil {
		http.Error(w, "MongoDB connection is not initialized", http.StatusInternalServerError)
		return
	}

	// Получаем параметры из запроса
	tempCategory := r.URL.Query().Get("category")
	category := strings.ReplaceAll(tempCategory, "%20", " ")

	// Получаем параметр articleId из запроса
	makenameStr := r.URL.Query().Get("makename")
	if makenameStr == "" {
		http.Error(w, "Required parameter 'makenameStr' is missing", http.StatusBadRequest)
		fmt.Println("makenameStr parameter is missing -------")
		return
	}

	fmt.Println("makenameStr:", makenameStr)

	// Преобразуем articleId из строки в число
	make, err := strconv.Atoi(makenameStr)
	if err != nil {
		http.Error(w, "Invalid makenameStr parameter", http.StatusBadRequest)
		fmt.Println("makename converted to int error -------")
		return
	}
	fmt.Println("makeName converted to int:", make)

	if category == "" {
		http.Error(w, "makeName and category are required parameters", http.StatusBadRequest)
		return
	}

	// filter := bson.M{"categoryname": category}
	// filter = bson.M{"identificator": makenameStr}

	filter := bson.M{
		"$and": []bson.M{
			{"categoryname": bson.M{"$regex": category, "$options": "i"}},
			{"identificator": fmt.Sprintf("%s", makenameStr)},
		},
	}

	log.Printf("Filter: %+v", filter)

	// Указываем коллекцию
	motorcyclesCollection := client.Database("moto").Collection("modelsbyCategories")

	// Функция для поиска в базе данных
	findInDatabase := func() ([]services.ModelWithCategory, error) {
		cursor, err := motorcyclesCollection.Find(context.Background(), filter)
		if err != nil {
			return nil, fmt.Errorf("failed to query database: %w", err)
		}
		defer cursor.Close(context.Background())

		var models []services.ModelWithCategory
		for cursor.Next(context.Background()) {
			var model services.ModelWithCategory
			if err := cursor.Decode(&model); err != nil {
				return nil, fmt.Errorf("failed to decode database result: %w", err)
			}
			models = append(models, model)
		}
		return models, nil
	}

	// 1. Ищем в базе данных
	models, err := findInDatabase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("models not found in DB:", filter)
		return
	}

	// 2. Если модели найдены, возвращаем их
	if len(models) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models)
		fmt.Println("parsing models:", models)
		return
	}

	fmt.Println("BD models by categories: - ничего не найдено, вызываем внешний API")
	// 3. Если ничего не найдено, вызываем внешний API
	externalModels, err := apiClient.FetchModelByCategories(category, makenameStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch data from API: %v", err), http.StatusInternalServerError)
		return
	}

	// 4. Сохраняем результат внешнего API в базу
	for _, model := range externalModels {
		// Преобразуем структуру в BSON перед сохранением
		bsonData, err := bson.Marshal(model)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to marshal model: %v", err), http.StatusInternalServerError)
			return
		}

		_, err = motorcyclesCollection.InsertOne(context.Background(), bsonData)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to save model to database: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// 5. Снова ищем в базе данных после вставки
	models, err = findInDatabase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. Возвращаем данные, если они найдены, или пустой массив в противном случае
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}
