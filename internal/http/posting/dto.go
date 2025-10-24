package posting

type CreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Category    string `json:"category"`
	City        string `json:"city"`
	District    string `json:"district"`
}
