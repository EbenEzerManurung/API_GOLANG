package models

type Customer struct {
    Custcd          string `json:"custcd"`
    NamaCustomer    string `json:"nama_customer"`
    Address         string `json:"address"`
    Phone           string `json:"phone"`
    InactiveCustomer string `json:"inactive_customer"`
}