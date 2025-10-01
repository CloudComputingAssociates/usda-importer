package usda

// SurveyFood represents USDA survey food structure
type SurveyFood struct {
	FdcID           int                `bson:"fdcId"`
	Description     string             `bson:"description"`
	FoodClass       string             `bson:"foodClass"`
	DataType        string             `bson:"dataType"`
	FoodCode        string             `bson:"foodCode"`
	PublicationDate string             `bson:"publicationDate"`
	FoodNutrients   []FoodNutrient     `bson:"foodNutrients"`
	FoodPortions    []FoodPortion      `bson:"foodPortions"`
}

// FoodNutrient represents a nutrient in USDA data
type FoodNutrient struct {
	Type     string   `bson:"type"`
	ID       int      `bson:"id"`
	Nutrient Nutrient `bson:"nutrient"`
	Amount   float64  `bson:"amount"`
}

// Nutrient represents nutrient details
type Nutrient struct {
	ID       int    `bson:"id"`
	Number   string `bson:"number"`
	Name     string `bson:"name"`
	Rank     int    `bson:"rank"`
	UnitName string `bson:"unitName"`
}

// FoodPortion represents a serving size
type FoodPortion struct {
	ID                 int     `bson:"id"`
	GramWeight         float64 `bson:"gramWeight"`
	PortionDescription string  `bson:"portionDescription"`
	SequenceNumber     int     `bson:"sequenceNumber"`
}