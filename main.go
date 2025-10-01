package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"usda-importer/mapper"
)

func main() {
	// Parse command line flags
	importType := flag.String("type", "", "Type of import: 'survey' or 'branded'")
	flag.Parse()

	if *importType != "survey" && *importType != "branded" {
		log.Fatal("Error: --type flag must be 'survey' or 'branded'")
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// MongoDB connection strings
	foodMongoURI := os.Getenv("MONGO_URI")
	foodDB := os.Getenv("MONGO_DB")
	foodCollection := os.Getenv("MONGO_COLLECTION")

	usdaMongoURI := os.Getenv("USDA_MONGO_URI")
	usdaDB := os.Getenv("USDA_MONGO_DB")
	
	var usdaCollection string
	if *importType == "survey" {
		usdaCollection = os.Getenv("USDA_SURVEY_COLLECTION")
	} else {
		usdaCollection = os.Getenv("USDA_BRANDED_COLLECTION")
	}

	// Validate environment variables
	if foodMongoURI == "" || foodDB == "" || foodCollection == "" {
		log.Fatal("Error: MONGO_URI, MONGO_DB, and MONGO_COLLECTION must be set")
	}
	if usdaMongoURI == "" || usdaDB == "" {
		log.Fatal("Error: USDA_MONGO_URI and USDA_MONGO_DB must be set")
	}

	ctx := context.Background()

	// Connect to Food database
	log.Println("Connecting to Food database...")
	foodClient, err := mongo.Connect(ctx, options.Client().ApplyURI(foodMongoURI))
	if err != nil {
		log.Fatal("Failed to connect to Food database:", err)
	}
	defer foodClient.Disconnect(ctx)

	// Connect to USDA database
	log.Println("Connecting to USDA database...")
	usdaClient, err := mongo.Connect(ctx, options.Client().ApplyURI(usdaMongoURI))
	if err != nil {
		log.Fatal("Failed to connect to USDA database:", err)
	}
	defer usdaClient.Disconnect(ctx)

	// Get collections
	foodColl := foodClient.Database(foodDB).Collection(foodCollection)
	usdaColl := usdaClient.Database(usdaDB).Collection(usdaCollection)

	// Count USDA records
	totalCount, err := usdaColl.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to count USDA records:", err)
	}
	log.Printf("Found %d %s foods to import\n", totalCount, *importType)

	// Start import
	startTime := time.Now()
	log.Printf("Starting import of %s foods...\n", *importType)

	var importFunc func(context.Context, *mongo.Collection, *mongo.Collection) (int, int, error)
	if *importType == "survey" {
		importFunc = mapper.ImportSurveyFoods
	} else {
		importFunc = mapper.ImportBrandedFoods
	}

	imported, skipped, err := importFunc(ctx, usdaColl, foodColl)
	if err != nil {
		log.Fatal("Import failed:", err)
	}

	duration := time.Since(startTime)
	log.Printf("\n=== Import Complete ===")
	log.Printf("Total records processed: %d", totalCount)
	log.Printf("Successfully imported: %d", imported)
	log.Printf("Skipped (empty nutrients): %d", skipped)
	log.Printf("Duration: %s", duration)
	log.Printf("Rate: %.2f records/second", float64(imported)/duration.Seconds())
}