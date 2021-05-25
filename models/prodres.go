package models

type Product struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Repertory int     `json:"repertory"`
}

type ProdRes struct {
	Total int `json:"total"`
	List []Product `json:"list"`
}