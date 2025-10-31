package models

// CarModel represents a car model in the system
type CarModel struct {
	Name     string                     `json:"Name"`
	BodyType string                     `json:"BodyType"` // Sedan, SUV
	FuelType string                     `json:"FuelType,omitempty"` // Electric, Hybrid, x80, x90, x92, x95
	Slug     string                     `json:"Slug"`
	Brand    *RelationField[Brand]      `json:"brand,omitempty"`
	Images   *MediaCollectionField      `json:"Images,omitempty"`
}

// CarModelResponse is a convenience type for car model API responses
type CarModelResponse = StrapiResponse[CarModel]

// CarModelCollectionResponse is a convenience type for car model collection API responses
type CarModelCollectionResponse = StrapiCollectionResponse[CarModel]

// SimpleCarModel is a simplified car model response with pricing
type SimpleCarModel struct {
	ID              string      `json:"id"`
	Title           string      `json:"title"`
	Thumbnail       *MediaField `json:"thumbnail,omitempty"`
	PriceFrom       int         `json:"pricefrom"`
	PriceTo         int         `json:"priceto"`
	MarketPriceFrom int         `json:"marketpricefrom"`
	MarketPriceTo   int         `json:"marketpriceto"`
}

// DetailedCarModel is a detailed car model response with all variants and showrooms
type DetailedCarModel struct {
	ID              string                  `json:"id"`
	Title           string                  `json:"title"`
	Images          *MediaCollectionField   `json:"images,omitempty"`
	PriceFrom       int                     `json:"pricefrom"`
	PriceTo         int                     `json:"priceto"`
	MarketPriceFrom int                     `json:"marketpricefrom"`
	MarketPriceTo   int                     `json:"marketpriceto"`
	MinDownPayment  int                     `json:"mindownpayment"`
	MinInstallments int                     `json:"mininstallments"`
	Warranty        string                  `json:"warranty"`
	Variants        []SimpleVariant         `json:"variants"`
	Showrooms       []SimpleShowroom        `json:"showrooms"`
	Reviews         []ReviewItem            `json:"reviews"`
	Catalogs        []CatalogItem           `json:"catalogs"`
}

// SimpleVariant represents a simplified variant in the detailed car model response
type SimpleVariant struct {
	ID    string `json:"id"`
	Year  int    `json:"year"`
	Title string `json:"title"`
	CC    string `json:"cc"`
	Price int    `json:"price"`
}

// SimpleShowroom represents a simplified showroom in the detailed car model response
type SimpleShowroom struct {
	ID              string      `json:"id"`
	Title           string      `json:"title"`
	Thumbnail       *MediaField `json:"thumbnail,omitempty"`
	Price           int         `json:"price"`
	MinDownPayment  int         `json:"mindownpayment"`
	MinInstallments int         `json:"mininstallments"`
}

// ReviewItem represents a review link
type ReviewItem struct {
	ID         int    `json:"id"`
	YoutubeURL string `json:"youtubeurl"`
}

// CatalogItem represents a catalog/brochure
type CatalogItem struct {
	ID          int    `json:"id"`
	DownloadURL string `json:"downloadurl"`
}
