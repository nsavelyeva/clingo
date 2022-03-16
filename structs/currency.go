package structs

// DetailsCurrency is a struct to store currency details for a certain country
type DetailsCurrency struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	SymbolNative  string  `json:"symbol_native"`
	DecimalDigits int     `json:"decimal_digits"`
	Rounding      float32 `json:"rounding"`
	Code          string  `json:"code"`
	NamePlural    string  `json:"name_plural"`
}

// ResponseCurrency is a struct to store successful HTTP response from currency API
type ResponseCurrency struct {
	Meta *Meta `json:"query"`
	Data *Data `json:"data"`
}

// Meta is a sub-struct of ResponseCurrency struct
type Meta struct {
	LastUpdatedAt int `json:"last_updated_at"`
}

// Rate is struct of every item in Data sub-struct of ResponseCurrency struct
type Rate struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

// Data is a sub-struct of ResponseCurrency struct
type Data struct {
	USD *Rate `json:"USD"`
	JPY *Rate `json:"JPY"`
	CNY *Rate `json:"CNY"`
	CHF *Rate `json:"CHF"`
	CAD *Rate `json:"CAD"`
	MXN *Rate `json:"MXN"`
	INR *Rate `json:"INR"`
	BRL *Rate `json:"BRL"`
	RUB *Rate `json:"RUB"`
	KRW *Rate `json:"KRW"`
	IDR *Rate `json:"IDR"`
	TRY *Rate `json:"TRY"`
	SAR *Rate `json:"SAR"`
	SEK *Rate `json:"SEK"`
	NGN *Rate `json:"NGN"`
	PLN *Rate `json:"PLN"`
	ARS *Rate `json:"ARS"`
	NOK *Rate `json:"NOK"`
	TWD *Rate `json:"TWD"`
	IRR *Rate `json:"IRR"`
	AED *Rate `json:"AED"`
	COP *Rate `json:"COP"`
	THB *Rate `json:"THB"`
	ZAR *Rate `json:"ZAR"`
	DKK *Rate `json:"DKK"`
	MYR *Rate `json:"MYR"`
	SGD *Rate `json:"SGD"`
	ILS *Rate `json:"ILS"`
	HKD *Rate `json:"HKD"`
	EGP *Rate `json:"EGP"`
	PHP *Rate `json:"PHP"`
	CLP *Rate `json:"CLP"`
	PKR *Rate `json:"PKR"`
	IQD *Rate `json:"IQD"`
	DZD *Rate `json:"DZD"`
	KZT *Rate `json:"KZT"`
	QAR *Rate `json:"QAR"`
	CZK *Rate `json:"CZK"`
	PEN *Rate `json:"PEN"`
	RON *Rate `json:"RON"`
	VND *Rate `json:"VND"`
	BDT *Rate `json:"BDT"`
	HUF *Rate `json:"HUF"`
	UAH *Rate `json:"UAH"`
	AOA *Rate `json:"AOA"`
	MAD *Rate `json:"MAD"`
	OMR *Rate `json:"OMR"`
	CUC *Rate `json:"CUC"`
	BYN *Rate `json:"BYN"`
	BYR *Rate `json:"BYR"`
	AZN *Rate `json:"AZN"`
	LKR *Rate `json:"LKR"`
	SDG *Rate `json:"SDG"`
	SYP *Rate `json:"SYP"`
	MMK *Rate `json:"MMK"`
	DOP *Rate `json:"DOP"`
	UZS *Rate `json:"UZS"`
	KES *Rate `json:"KES"`
	GTQ *Rate `json:"GTQ"`
	URY *Rate `json:"URY"`
	HRV *Rate `json:"HRV"`
	MOP *Rate `json:"MOP"`
	ETB *Rate `json:"ETB"`
	CRC *Rate `json:"CRC"`
	TZS *Rate `json:"TZS"`
	TMT *Rate `json:"TMT"`
	TND *Rate `json:"TND"`
	PAB *Rate `json:"PAB"`
	LBP *Rate `json:"LBP"`
	RSD *Rate `json:"RSD"`
	LYD *Rate `json:"LYD"`
	GHS *Rate `json:"GHS"`
	YER *Rate `json:"YER"`
	BOB *Rate `json:"BOB"`
	BHD *Rate `json:"BHD"`
	CDF *Rate `json:"CDF"`
	PYG *Rate `json:"PYG"`
	UGX *Rate `json:"UGX"`
	SVC *Rate `json:"SVC"`
	TTD *Rate `json:"TTD"`
	AFN *Rate `json:"AFN"`
	NPR *Rate `json:"NPR"`
	HNL *Rate `json:"HNL"`
	BIH *Rate `json:"BIH"`
	BND *Rate `json:"BND"`
	ISK *Rate `json:"ISK"`
	KHR *Rate `json:"KHR"`
	GEL *Rate `json:"GEL"`
	MZN *Rate `json:"MZN"`
	BWP *Rate `json:"BWP"`
	PGK *Rate `json:"PGK"`
	JMD *Rate `json:"JMD"`
	XAF *Rate `json:"XAF"`
	NAD *Rate `json:"NAD"`
	ALL *Rate `json:"ALL"`
	SSP *Rate `json:"SSP"`
	MUR *Rate `json:"MUR"`
	MNT *Rate `json:"MNT"`
	NIO *Rate `json:"NIO"`
	LAK *Rate `json:"LAK"`
	MKD *Rate `json:"MKD"`
	AMD *Rate `json:"AMD"`
	MGA *Rate `json:"MGA"`
	XPF *Rate `json:"XPF"`
	TJS *Rate `json:"TJS"`
	HTG *Rate `json:"HTG"`
	BSD *Rate `json:"BSD"`
	MDL *Rate `json:"MDL"`
	RWF *Rate `json:"RWF"`
	KGS *Rate `json:"KGS"`
	GNF *Rate `json:"GNF"`
	SRD *Rate `json:"SRD"`
	SLL *Rate `json:"SLL"`
	XOF *Rate `json:"XOF"`
	MWK *Rate `json:"MWK"`
	FJD *Rate `json:"FJD"`
	ERN *Rate `json:"ERN"`
	SZL *Rate `json:"SZL"`
	GYD *Rate `json:"GYD"`
	BIF *Rate `json:"BIF"`
	KYD *Rate `json:"KYD"`
	MVR *Rate `json:"MVR"`
	LSL *Rate `json:"LSL"`
	LRD *Rate `json:"LRD"`
	CVE *Rate `json:"CVE"`
	DJF *Rate `json:"DJF"`
	SCR *Rate `json:"SCR"`
	SOS *Rate `json:"SOS"`
	GMD *Rate `json:"GMD"`
	KMF *Rate `json:"KMF"`
	STD *Rate `json:"STD"`
	BTC *Rate `json:"BTC"`
	XRP *Rate `json:"XRP"`
	AUD *Rate `json:"AUD"`
	BGN *Rate `json:"BGN"`
	JOD *Rate `json:"JOD"`
	GBP *Rate `json:"GBP"`
	ETH *Rate `json:"ETH"`
	EUR *Rate `json:"EUR"`
	LTC *Rate `json:"LTC"`
	NZD *Rate `json:"NZD"`
}
