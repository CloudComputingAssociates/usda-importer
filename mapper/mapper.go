package mapper

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"usda-importer/models"
	"usda-importer/usda"
)

// ImportSurveyFoods imports USDA survey foods
func ImportSurveyFoods(ctx context.Context, usdaColl, foodColl *mongo.Collection) (int, int, error) {
	cursor, err := usdaColl.Find(ctx, bson.M{})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to query USDA survey foods: %w", err)
	}
	defer cursor.Close(ctx)

	imported := 0
	skipped := 0
	nextID := 1
	batch := []interface{}{}
	batchSize := 1000

	for cursor.Next(ctx) {
		var surveyFood usda.SurveyFood
		if err := cursor.Decode(&surveyFood); err != nil {
			log.Printf("Warning: failed to decode survey food: %v", err)
			continue
		}

		// Skip foods with no nutrients
		if len(surveyFood.FoodNutrients) == 0 {
			skipped++
			continue
		}

		// Map to Food model
		food := MapSurveyFood(&surveyFood, nextID)
		batch = append(batch, food)
		nextID++

		// Insert batch when full
		if len(batch) >= batchSize {
			if _, err := foodColl.InsertMany(ctx, batch); err != nil {
				return imported, skipped, fmt.Errorf("failed to insert batch: %w", err)
			}
			imported += len(batch)
			log.Printf("Imported %d survey foods...", imported)
			batch = []interface{}{}
		}
	}

	// Insert remaining batch
	if len(batch) > 0 {
		if _, err := foodColl.InsertMany(ctx, batch); err != nil {
			return imported, skipped, fmt.Errorf("failed to insert final batch: %w", err)
		}
		imported += len(batch)
	}

	return imported, skipped, nil
}

// ImportBrandedFoods imports USDA branded foods
func ImportBrandedFoods(ctx context.Context, usdaColl, foodColl *mongo.Collection) (int, int, error) {
	cursor, err := usdaColl.Find(ctx, bson.M{})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to query USDA branded foods: %w", err)
	}
	defer cursor.Close(ctx)

	imported := 0
	skipped := 0
	nextID := 1
	batch := []interface{}{}
	batchSize := 1000

	for cursor.Next(ctx) {
		var brandedFood usda.BrandedFood
		if err := cursor.Decode(&brandedFood); err != nil {
			log.Printf("Warning: failed to decode branded food: %v", err)
			continue
		}

		// Skip foods with no nutrients
		if len(brandedFood.FoodNutrients) == 0 {
			skipped++
			continue
		}

		// Map to Food model
		food := MapBrandedFood(&brandedFood, nextID)
		batch = append(batch, food)
		nextID++

		// Insert batch when full
		if len(batch) >= batchSize {
			if _, err := foodColl.InsertMany(ctx, batch); err != nil {
				return imported, skipped, fmt.Errorf("failed to insert batch: %w", err)
			}
			imported += len(batch)
			log.Printf("Imported %d branded foods...", imported)
			batch = []interface{}{}
		}
	}

	// Insert remaining batch
	if len(batch) > 0 {
		if _, err := foodColl.InsertMany(ctx, batch); err != nil {
			return imported, skipped, fmt.Errorf("failed to insert final batch: %w", err)
		}
		imported += len(batch)
	}

	return imported, skipped, nil
}

// MapSurveyFood converts USDA survey food to Food model
func MapSurveyFood(usda *usda.SurveyFood, id int) *models.Food {
	now := time.Now()
	pubDate := parseDate(usda.PublicationDate)

	food := &models.Food{
		ID:                      id,
		Description:             usda.Description,
		FoodRequestType:         models.FoodTypeWhole,
		ANDIscore:               0,
		GlycemicIndex:           0,
		ServingSizeMultiplicand: 1.0, // Survey foods already per 100g
		DataSource:              "USDA-FNDDS",
		EnhancedAt:              &now,
		VerifiedType:            models.VerifiedTypeUSDAsurvey,
		VerifiedDate:            pubDate,
		VerifiedBy:              "USDA-Import",
		FoodImage:               "",
		FoodImageThumbnail:      "",
		NutritionFactsImage:     "",
		NutritionFactsImagePending: "",
		NutritionFactsStatus:    models.NutritionFactsStatusCompleted,
		BrandInfo:               nil,
		Recipe:                  nil,
	}

	// Map nutrition facts
	food.NutritionFacts = mapSurveyNutritionFacts(usda)

	return food
}

// MapBrandedFood converts USDA branded food to Food model
func MapBrandedFood(usda *usda.BrandedFood, id int) *models.Food {
	now := time.Now()
	pubDate := parseDate(usda.PublicationDate)

	// Calculate serving size multiplicand (branded foods are per serving, not per 100g)
	servingSizeMultiplicand := 1.0
	if usda.ServingSize > 0 {
		// Convert serving to 100g basis
		servingSizeMultiplicand = usda.ServingSize / 100.0
	}

	food := &models.Food{
		ID:                      id,
		Description:             usda.Description,
		FoodRequestType:         models.FoodTypeBrand,
		ANDIscore:               0,
		GlycemicIndex:           0,
		ServingSizeMultiplicand: servingSizeMultiplicand,
		DataSource:              "USDA-Branded",
		EnhancedAt:              &now,
		VerifiedType:            models.VerifiedTypeUSDAbrand,
		VerifiedDate:            pubDate,
		VerifiedBy:              "USDA-Import",
		FoodImage:               "",
		FoodImageThumbnail:      "",
		NutritionFactsImage:     "",
		NutritionFactsImagePending: "",
		NutritionFactsStatus:    models.NutritionFactsStatusCompleted,
		Recipe:                  nil,
	}

	// Map brand info
	if usda.BrandOwner != "" {
		food.BrandInfo = &models.BrandInfo{
			Manufacturer: usda.BrandOwner,
		}
	}

	// Map nutrition facts
	food.NutritionFacts = mapBrandedNutritionFacts(usda)

	return food
}

// mapSurveyNutritionFacts creates NutritionFacts from survey food
func mapSurveyNutritionFacts(usda *usda.SurveyFood) *models.NutritionFacts {
	nf := &models.NutritionFacts{
		FoodName:     usda.Description,
		Ingredients:  []string{},
	}

	// Get primary serving size (sequenceNumber == 1)
	for _, portion := range usda.FoodPortions {
		if portion.SequenceNumber == 1 {
			nf.ServingSizeG = portion.GramWeight
			nf.ServingSizeHousehold = portion.PortionDescription
			break
		}
	}
	nf.ServingsPerContainer = 1

	// Map nutrients
	for _, nutrient := range usda.FoodNutrients {
		mapNutrient(nf, nutrient.Nutrient.ID, nutrient.Amount)
	}

	return nf
}

// mapBrandedNutritionFacts creates NutritionFacts from branded food
func mapBrandedNutritionFacts(usda *usda.BrandedFood) *models.NutritionFacts {
	nf := &models.NutritionFacts{
		FoodName:             usda.Description,
		ServingSizeG:         usda.ServingSize,
		ServingSizeHousehold: usda.HouseholdServingFullText,
		ServingsPerContainer: 1,
		Ingredients:          []string{},
	}

	// Parse ingredients string into array if present
	if usda.Ingredients != "" {
		nf.Ingredients = []string{usda.Ingredients}
	}

	// Map nutrients (branded foods need conversion to per 100g)
	conversionFactor := 100.0 / usda.ServingSize
	for _, nutrient := range usda.FoodNutrients {
		// Convert to per 100g
		amountPer100g := nutrient.Amount * conversionFactor
		mapNutrient(nf, nutrient.Nutrient.ID, amountPer100g)
	}

	return nf
}

// mapNutrient maps a single nutrient by ID to the appropriate field
func mapNutrient(nf *models.NutritionFacts, nutrientID int, amount float64) {
	switch nutrientID {
	case 1008: // Energy (kcal)
		nf.Calories = amount
	case 1003: // Protein
		nf.ProteinG = amount
	case 1004: // Total lipid (fat)
		nf.TotalFatG = amount
	case 1258: // Saturated fat
		nf.SaturatedFatG = amount
	case 1257: // Trans fat
		nf.TransFatG = &amount
	case 1253: // Cholesterol
		nf.CholesterolMG = amount
	case 1093: // Sodium
		nf.SodiumMG = amount
	case 1005: // Carbohydrate
		nf.TotalCarbohydrateG = amount
	case 1079: // Fiber
		nf.DietaryFiberG = amount
	case 2000: // Total sugars
		nf.TotalSugarsG = amount
	case 1114: // Vitamin D
		nf.VitaminDMcg = amount
	case 1087: // Calcium
		nf.CalciumMG = amount
	case 1089: // Iron
		nf.IronMG = amount
	case 1092: // Potassium
		nf.PotassiumMG = amount
	}
}

// parseDate parses USDA date strings (e.g., "10/31/2024")
func parseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	
	t, err := time.Parse("1/2/2006", dateStr)
	if err != nil {
		return nil
	}
	return &t
}