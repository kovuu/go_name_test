package models

type Person struct {
	Id          int    `json:"id,omitempty" db:"id"`
	Name        string `json:"name" db:"name" binding:"required"`
	Surname     string `json:"surname" db:"surname"  binding:"required"`
	Patronymic  string `json:"patronymic" db:"patronymic"`
	Age         int    `json:"age" db:"age"`
	Gender      string `json:"gender" db:"gender"`
	Nationality string `json:"nationality" db:"nationality"`
}

type PersonFailed struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
	Error      string
}
