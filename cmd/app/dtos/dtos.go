package dtos

type Currency struct {
	Symbol string `json:"symbol"`
	// Можноиспользовать не строку, а число с плавающей точкой, но
	// насколько я знаю, когда имеешь дело с валютой, лучше этого избегать (из-за погрешности)
	// поэтому любой сервис, который получает от нас данные сможет сам решить
	// какой тип ему нужен и перевести price
	Price string `json:"price"`
}
