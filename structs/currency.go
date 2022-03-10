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
	Query *Query `json:"query"`
	Data  *Data  `json:"data"`
}

// Query is a sub-struct of ResponseCurrency struct
type Query struct {
	Apikey       string `json:"apikey"`
	BaseCurrency string `json:"base_currency"`
	Timestamp    int    `json:"timestamp"`
}

// Data is a sub-struct of ResponseCurrency struct
type Data struct {
	USD float64 `json:"USD"`
	JPY float64 `json:"JPY"`
	CNY float64 `json:"CNY"`
	CHF float64 `json:"CHF"`
	CAD float64 `json:"CAD"`
	MXN float64 `json:"MXN"`
	INR float64 `json:"INR"`
	BRL float64 `json:"BRL"`
	RUB float64 `json:"RUB"`
	KRW float64 `json:"KRW"`
	IDR float64 `json:"IDR"`
	TRY float64 `json:"TRY"`
	SAR float64 `json:"SAR"`
	SEK float64 `json:"SEK"`
	NGN float64 `json:"NGN"`
	PLN float64 `json:"PLN"`
	ARS float64 `json:"ARS"`
	NOK float64 `json:"NOK"`
	TWD float64 `json:"TWD"`
	IRR float64 `json:"IRR"`
	AED float64 `json:"AED"`
	COP float64 `json:"COP"`
	THB float64 `json:"THB"`
	ZAR float64 `json:"ZAR"`
	DKK float64 `json:"DKK"`
	MYR float64 `json:"MYR"`
	SGD float64 `json:"SGD"`
	ILS float64 `json:"ILS"`
	HKD float64 `json:"HKD"`
	EGP float64 `json:"EGP"`
	PHP float64 `json:"PHP"`
	CLP float64 `json:"CLP"`
	PKR float64 `json:"PKR"`
	IQD float64 `json:"IQD"`
	DZD float64 `json:"DZD"`
	KZT float64 `json:"KZT"`
	QAR float64 `json:"QAR"`
	CZK float64 `json:"CZK"`
	PEN float64 `json:"PEN"`
	RON float64 `json:"RON"`
	VND float64 `json:"VND"`
	BDT float64 `json:"BDT"`
	HUF float64 `json:"HUF"`
	UAH float64 `json:"UAH"`
	AOA float64 `json:"AOA"`
	MAD float64 `json:"MAD"`
	OMR float64 `json:"OMR"`
	CUC float64 `json:"CUC"`
	BYR float64 `json:"BYR"`
	AZN float64 `json:"AZN"`
	LKR float64 `json:"LKR"`
	SDG float64 `json:"SDG"`
	SYP float64 `json:"SYP"`
	MMK float64 `json:"MMK"`
	DOP float64 `json:"DOP"`
	UZS float64 `json:"UZS"`
	KES float64 `json:"KES"`
	GTQ float64 `json:"GTQ"`
	URY float64 `json:"URY"`
	HRV float64 `json:"HRV"`
	MOP float64 `json:"MOP"`
	ETB float64 `json:"ETB"`
	CRC float64 `json:"CRC"`
	TZS float64 `json:"TZS"`
	TMT float64 `json:"TMT"`
	TND float64 `json:"TND"`
	PAB float64 `json:"PAB"`
	LBP float64 `json:"LBP"`
	RSD float64 `json:"RSD"`
	LYD float64 `json:"LYD"`
	GHS float64 `json:"GHS"`
	YER float64 `json:"YER"`
	BOB float64 `json:"BOB"`
	BHD float64 `json:"BHD"`
	CDF float64 `json:"CDF"`
	PYG float64 `json:"PYG"`
	UGX float64 `json:"UGX"`
	SVC float64 `json:"SVC"`
	TTD float64 `json:"TTD"`
	AFN float64 `json:"AFN"`
	NPR float64 `json:"NPR"`
	HNL float64 `json:"HNL"`
	BIH float64 `json:"BIH"`
	BND float64 `json:"BND"`
	ISK float64 `json:"ISK"`
	KHR float64 `json:"KHR"`
	GEL float64 `json:"GEL"`
	MZN float64 `json:"MZN"`
	BWP float64 `json:"BWP"`
	PGK float64 `json:"PGK"`
	JMD float64 `json:"JMD"`
	XAF float64 `json:"XAF"`
	NAD float64 `json:"NAD"`
	ALL float64 `json:"ALL"`
	SSP float64 `json:"SSP"`
	MUR float64 `json:"MUR"`
	MNT float64 `json:"MNT"`
	NIO float64 `json:"NIO"`
	LAK float64 `json:"LAK"`
	MKD float64 `json:"MKD"`
	AMD float64 `json:"AMD"`
	MGA float64 `json:"MGA"`
	XPF float64 `json:"XPF"`
	TJS float64 `json:"TJS"`
	HTG float64 `json:"HTG"`
	BSD float64 `json:"BSD"`
	MDL float64 `json:"MDL"`
	RWF float64 `json:"RWF"`
	KGS float64 `json:"KGS"`
	GNF float64 `json:"GNF"`
	SRD float64 `json:"SRD"`
	SLL float64 `json:"SLL"`
	XOF float64 `json:"XOF"`
	MWK float64 `json:"MWK"`
	FJD float64 `json:"FJD"`
	ERN float64 `json:"ERN"`
	SZL float64 `json:"SZL"`
	GYD float64 `json:"GYD"`
	BIF float64 `json:"BIF"`
	KYD float64 `json:"KYD"`
	MVR float64 `json:"MVR"`
	LSL float64 `json:"LSL"`
	LRD float64 `json:"LRD"`
	CVE float64 `json:"CVE"`
	DJF float64 `json:"DJF"`
	SCR float64 `json:"SCR"`
	SOS float64 `json:"SOS"`
	GMD float64 `json:"GMD"`
	KMF float64 `json:"KMF"`
	STD float64 `json:"STD"`
	BTC float64 `json:"BTC"`
	XRP float64 `json:"XRP"`
	AUD float64 `json:"AUD"`
	BGN float64 `json:"BGN"`
	JOD float64 `json:"JOD"`
	GBP float64 `json:"GBP"`
	ETH float64 `json:"ETH"`
	EUR float64 `json:"EUR"`
	LTC float64 `json:"LTC"`
	NZD float64 `json:"NZD"`
}
