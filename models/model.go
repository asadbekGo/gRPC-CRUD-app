package models

type User struct {
	Id        int32
	FirstName string
	LastName  string
	Username  string
	Age       int32
}

type PostUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Age       int32  `json:"age"`
}

type Id struct {
	Id int32 `json:"id"`
}
