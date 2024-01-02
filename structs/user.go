package structs

type Addr struct {
	City   string `json:"city"`
	Street string `json:"street"`
	Area   string `json:"area"`
}
type Visa struct {
	Number         string `json:"number"`
	Cvv            string `json:"cvv"`
	ExpirationData string `json:"expirationData"`
}

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserEmail string `json:"email"`
	Name      string `json:"Name"`
	Phone     string `json:"phone"`
	UserAddr  Addr   `json:"addr"`
	UserVisa  Visa `json:"visa"`
}
