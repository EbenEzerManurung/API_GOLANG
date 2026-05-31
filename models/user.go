package models

type User struct {
    IDUser   int    `json:"id_user"`
    NamaUser string `json:"nama_user"`
    RoleUser string `json:"role_user"`
}