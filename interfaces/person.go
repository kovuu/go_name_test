package interfaces

type Person struct {
	Name        string `json:"name" binding:"required"`
	Surname     string `json:"surname"  binding:"required"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type PersonFailed struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
	Error      string
}
