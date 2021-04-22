package statics

type BetLog struct {
	BlId      int
	BId       int
	CId       int
	ReserveId string
	ReqId     string
	Balance   float64
	Endpoint  string
	XmlBody   string
	ResBody   string
	BlType    string
	Amount    float64
	CreatedAt string
}

type Bets struct {
	BId       int
	CwId      int
	ReserveId string
	Balance   float64
	CreatedAt string
	UpdatedAt string
	Status    int
}
