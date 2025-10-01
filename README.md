usda-importer

takes two params, depending on which collection is being imported

used initial to load survey-foods and branded-foods from USDA

Build:
C:\> cd usda-importer
C:\git\usda-importer> go mod download
C:\git\usda-importer> go mod tidy
C:\git\usda-importer> go build



usage: 
C:\git\usda-importer> usda-importer --type=survey
C:\git\usda-importer> usda-importer --type=branded

Note:  usda db has the collections being imported (survey-foods, branded-foods) and they will go in food db  foods collection

project structure
=================

usda-importer/
├── .env                    
├── .env.example           
├── .gitignore             
├── README.md              
├── main.go                
├── go.mod                 
├── models/
│   └── food.go            ← COPIED from git\FoodsAPI\models\food.go
├── mapper/
│   └── mapper.go          
└── usda/
    ├── survey.go          
    └── branded.go         