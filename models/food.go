// File: models/food.go
package models

import (
	"strings"
	"time"
)

type FoodRequestType string
type VerifiedType string
type NutritionFactsStatus string

const (
	FoodTypeUnknown FoodRequestType = "unknown"
	FoodTypeBrand   FoodRequestType = "brand"
	FoodTypeWhole   FoodRequestType = "whole"
	FoodTypeRecipe  FoodRequestType = "recipe"
)

const (
	VerifiedTypeUnknown      VerifiedType = "unknown"
	VerifiedTypeNutritionix  VerifiedType = "Nutritionix"
	VerifiedTypeProductLabel VerifiedType = "ProductLabel"
	VerifiedTypeAIModel      VerifiedType = "AIModel"
	VerifiedTypeUSDAsurvey   VerifiedType = "USDAsurvey"
	VerifiedTypeUSDAbrand    VerifiedType = "USDAbrand"
	VerifiedTypeBrand        VerifiedType = "brand"
)

const (
	NutritionFactsStatusPending    NutritionFactsStatus = "pending"
	NutritionFactsStatusProcessing NutritionFactsStatus = "processing"
	NutritionFactsStatusCompleted  NutritionFactsStatus = "completed"
	NutritionFactsStatusError      NutritionFactsStatus = "error"
)

// BrandInfo contains manufacturer intelligence for targeted searches
type BrandInfo struct {
	Manufacturer               string   `json:"manufacturer" bson:"manufacturer"`
	ParentCompany              string   `json:"parentCompany" bson:"parentCompany"`
	OfficialSites              []string `json:"officialSites" bson:"officialSites"`
	NutritionSiteCandidates    []string `json:"nutritionSiteCandidates" bson:"nutritionSiteCandidates"`
	ProductImageSiteCandidates []string `json:"productImageSiteCandidates" bson:"productImageSiteCandidates"`
	ProductLine                string   `json:"productLine" bson:"productLine"`
}

type Food struct {
	ID                         int                  `json:"id" bson:"id"`
	Description                string               `json:"description" bson:"description"`
	FoodRequestType            FoodRequestType      `json:"foodRequestType" bson:"foodRequestType"`
	ANDIscore                  float64              `json:"ANDIscore" bson:"ANDIscore"`
	GlycemicIndex              float64              `json:"GlycemicIndex" bson:"GlycemicIndex"`
	NutritionFacts             *NutritionFacts      `json:"nutritionFacts" bson:"nutritionFacts"`
	ServingSizeMultiplicand    float64              `json:"servingSizeMultiplicand" bson:"servingSizeMultiplicand"`
	DataSource                 string               `json:"dataSource" bson:"dataSource"`
	EnhancedAt                 *time.Time           `json:"enhancedAt" bson:"enhancedAt"`
	VerifiedType               VerifiedType         `json:"verifiedType" bson:"verifiedType"`
	VerifiedDate               *time.Time           `json:"verifiedDate" bson:"verifiedDate"`
	VerifiedBy                 string               `json:"verifiedBy" bson:"verifiedBy"`
	FoodImage                  string               `json:"foodImage" bson:"foodImage"`
	FoodImageThumbnail         string               `json:"foodImageThumbnail" bson:"foodImageThumbnail"`
	NutritionFactsImage        string               `json:"nutritionFactsImage" bson:"nutritionFactsImage"`
	NutritionFactsImagePending string               `json:"nutritionFactsImagePending" bson:"nutritionFactsImagePending"`
	NutritionFactsStatus       NutritionFactsStatus `json:"nutritionFactsStatus" bson:"nutritionFactsStatus"`
	TokensUsed                 *int                 `json:"tokensUsed" bson:"tokensUsed"`
	EstimatedCost              *float64             `json:"estimatedCost" bson:"estimatedCost"`
	BrandInfo                  *BrandInfo           `json:"brandInfo" bson:"brandInfo"`
	Recipe                     *Recipe              `json:"recipe" bson:"recipe"`
}

type NutritionFacts struct {
	FoodName             string   `json:"foodName" bson:"foodName"`
	ServingSizeHousehold string   `json:"servingSizeHousehold" bson:"servingSizeHousehold"`
	ServingSizeG         float64  `json:"servingSizeG" bson:"servingSizeG"`
	ServingsPerContainer int      `json:"servingsPerContainer" bson:"servingsPerContainer"`
	Calories             float64  `json:"calories" bson:"calories"`
	TotalFatG            float64  `json:"totalFatG" bson:"totalFatG"`
	SaturatedFatG        float64  `json:"saturatedFatG" bson:"saturatedFatG"`
	TransFatG            *float64 `json:"transFatG" bson:"transFatG"`
	CholesterolMG        float64  `json:"cholesterolMG" bson:"cholesterolMG"`
	SodiumMG             float64  `json:"sodiumMG" bson:"sodiumMG"`
	TotalCarbohydrateG   float64  `json:"totalCarbohydrateG" bson:"totalCarbohydrateG"`
	DietaryFiberG        float64  `json:"dietaryFiberG" bson:"dietaryFiberG"`
	TotalSugarsG         float64  `json:"totalSugarsG" bson:"totalSugarsG"`
	AddedSugarsG         *float64 `json:"addedSugarsG" bson:"addedSugarsG"`
	ProteinG             float64  `json:"proteinG" bson:"proteinG"`
	VitaminDMcg          float64  `json:"vitaminDMcg" bson:"vitaminDMcg"`
	CalciumMG            float64  `json:"calciumMG" bson:"calciumMG"`
	IronMG               float64  `json:"ironMG" bson:"ironMG"`
	PotassiumMG          float64  `json:"potassiumMG" bson:"potassiumMG"`
	Ingredients          []string `json:"ingredients" bson:"ingredients"`
}

// NormalizeFoodDescription converts description to lowercase for consistent searching
func NormalizeFoodDescription(desc string) string {
	return strings.ToLower(strings.TrimSpace(desc))
}

// Helper functions for verification workflows
func (f *Food) IsVerified() bool {
	return f.VerifiedType == VerifiedTypeProductLabel
}

type Recipe struct {
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description"`
	TotalWeightG float64            `json:"totalWeightG" bson:"totalweightg"`
	Servings     int                `json:"servings" bson:"servings"`
	Ingredients  []RecipeIngredient `json:"ingredients" bson:"ingredients"`
	Instructions string             `json:"instructions" bson:"instructions"`
}

type RecipeIngredient struct {
	Name        string  `json:"name" bson:"name"`
	Quantity    float64 `json:"quantity" bson:"quantity"`
	Unit        string  `json:"unit" bson:"unit"`
	Preparation string  `json:"preparation" bson:"preparation"`
	Amount      float64 `json:"amount" bson:"amount"`
	WeightG     float64 `json:"weightG" bson:"weightg"`
	Calories    float64 `json:"calories" bson:"calories"`
	TotalFatG   float64 `json:"totalFatG" bson:"totalfatg"`
	ProteinG    float64 `json:"proteinG" bson:"proteing"`
	CarbsG      float64 `json:"carbsG" bson:"carbsg"`
}