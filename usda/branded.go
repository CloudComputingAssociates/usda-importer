package usda

// BrandedFood represents USDA branded food structure
type BrandedFood struct {
	FdcID                     int            `bson:"fdcId"`
	Description               string         `bson:"description"`
	FoodClass                 string         `bson:"foodClass"`
	DataType                  string         `bson:"dataType"`
	PublicationDate           string         `bson:"publicationDate"`
	BrandOwner                string         `bson:"brandOwner"`
	GtinUpc                   string         `bson:"gtinUpc"`
	Ingredients               string         `bson:"ingredients"`
	ServingSize               float64        `bson:"servingSize"`
	ServingSizeUnit           string         `bson:"servingSizeUnit"`
	HouseholdServingFullText  string         `bson:"householdServingFullText"`
	BrandedFoodCategory       string         `bson:"brandedFoodCategory"`
	FoodNutrients             []FoodNutrient `bson:"foodNutrients"`
}