package scrapeprocess

// Define the struct we want Claude to populate
type AIExtraction struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}
