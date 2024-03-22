package lighthousedata

// Define a struct
type PageSpeedOutputStruct struct {
	LighthouseResult struct {
		Audits struct {
			FirstContentfulPaint struct {
				DisplayValue string  `json:"displayValue"`
				Score        float64 `json:"score"`
			} `json:"first-contentful-paint"`
			LargestContentfulPaint struct {
				DisplayValue string  `json:"displayValue"`
				Score        float64 `json:"score"`
			} `json:"largest-contentful-paint"`
			Interactive struct {
				DisplayValue string  `json:"displayValue"`
				Score        float64 `json:"score"`
				NumericValue float64 `json:"numericValue"`
			} `json:"interactive"`
		} `json:"audits"`
		Categories struct {
			Performance struct {
				Score float64 `json:"score"`
			} `json:"performance"`
		} `json:"categories"`
	} `json:"lighthouseResult"`
}
