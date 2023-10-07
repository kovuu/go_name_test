package types

type Person struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname"  binding:"required"`
	Patronymic string `json:"patronymic,omitempty"`
}

type PersonFailed struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
	Error      string
}
