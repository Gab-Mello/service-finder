package user

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ProviderProfileRequest struct {
	Bio       string `json:"bio"`
	Phone     string `json:"phone"`
	Expertise string `json:"expertise"`
	City      string `json:"city"`
	District  string `json:"district"`
}
