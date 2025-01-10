package domain

//----------------------------------------------ENTITIES----------------------------------------------------//

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegister struct {
	Username        string `json:"username" validate:"required" example:"johndoe"`
	Email           string `json:"email" validate:"required,email" example:"example@example.com"`
	Password        string `json:"password" validate:"required" example:"password"`
	ConfirmPassword string `json:"confirmPassword" validate:"required" example:"password"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email" example:"example@example.com"`
	Password string `json:"password" validate:"required" example:"password"`
}

//----------------------------------------------RESPONSES----------------------------------------------------//

type UserResponseLogin struct {
	Code  string `json:"code" example:"200"`
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}

type Response struct {
	Code    string `json:"code" example:"201"`
	Message string `json:"message" example:"User added"`
}

type BadResponse struct {
	Code    string `json:"code" example:"400"`
	Message string `json:"message" example:"Bad request"`
}

//-----------------------------------------------OFFERS-------------------------------------------------------//

type Offer struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type OfferWithPrice struct {
	ID       int     `json:"id" validate:"required" example:"1"`
	Name     string  `json:"name" validate:"required" example:"meat"`
	Quantity int     `json:"quantity" validate:"required" example:"9"`
	Price    float64 `json:"price" validate:"required" example:"10.5"`
	Category string  `json:"category" validate:"required" example:"food"`
}

type OffersRegister struct {
	Code    string           `json:"code" example:"200"`
	Message []OfferWithPrice `json:"message"`
}

//-----------------------------------------------ORDERS-------------------------------------------------------//

type Items struct {
	ItemID   int `json:"itemID" validate:"required" example:"1"`
	Quantity int `json:"quantity" validate:"required" example:"23"`
}

type OrderCheckout struct {
	Order []Items `json:"order"`
}

type Order struct {
	Total  int    `json:"total" validate:"required" example:"100"`
	Status string `json:"status" validate:"required" example:"pending"`
}

type OrderResponse struct {
	Code    string `json:"code" example:"200"`
	Message Order  `json:"message"`
}

type UserOrderStatus struct {
	ID     int    `json:"id" example:"1"`
	User   string `json:"user" example:"johndoe"`
	Total  int    `json:"total" example:"100"`
	Status string `json:"status" example:"pending"`
}

type OrderStatusResponse struct {
	Code    string          `json:"code" example:"200"`
	Message UserOrderStatus `json:"message"`
}

//-----------------------------------------------ADMIN-------------------------------------------------------//

type Dashboard struct {
	Offers  []OfferWithPrice  `json:"offers"`
	Orders  []UserOrderStatus `json:"orders"`
	Balance int               `json:"balance" example:"100"`
}

type DashboardResponse struct {
	Code    string    `json:"code" example:"200"`
	Message Dashboard `json:"message"`
}

type OrderStatusUpdate struct {
	Status string `json:"status" validate:"required" example:"delivered"`
}

//------------------------------------THE TUTORIAL TOLD ME TO DO THAT-----------------------------------------//

func NewUser(id int, email string, password string) *User {
	return &User{
		ID:       id,
		Email:    email,
		Password: password,
	}
}

func (u *User) GetEmail() string {
	return u.Email
}
