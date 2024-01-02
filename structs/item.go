package structs

type UserCommint struct {
	Username string `josn:"username"`
	Commint  string `json:"Commint"`
	Stars    int    `json:"stars"`
}

type Model struct {
	Id          int                       `json:"id"`
	Sizes       map[string]map[string]int `json:"sizes"`
	Price       float32                   `json:"price"`
	Description string                    `json:"description"`
	Discount    int                       `json:"discount"`
	LinkesImage []string                  `json:"linkesImage"`
	Commint     []UserCommint             `json:"commint"`
}
