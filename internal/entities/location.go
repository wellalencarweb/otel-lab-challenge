package entities

type Location struct {
	Zipcode      string `json:"cep"`
	AddressLine1 string `json:"logradouro"`
	AddressLine2 string `json:"complemento"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
	IBGECode     string `json:"ibge"`
	GIACode      string `json:"gia"`
	AreaCode     string `json:"ddd"`
	SIAFICode    string `json:"siafi"`
}
