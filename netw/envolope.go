package netw

type Envelope struct {
	Client      string      `json:"client"`
	MessageCode MessageCode `json:"msg_code"`
	Message     interface{} `json:"msg,omitempty"`
}

type EnvelopeStaging struct {
	Client      string      `json:"client"`
	MessageCode MessageCode `json:"msg_code"`
	Message     string      `json:"msg,omitempty"`
}

type Event struct {
	Id         string `json:"id"`
	InternalId string `json:"internal_id"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}
type Message struct {
	Id      string `json:"id"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Register struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	IsEnabled bool   `json:"is_enabled"`
	Result    string `json:"result"`
}

type PlayGame struct {
	Id   string `json:"id"`
	Mode string `json:"mode"`
}

type User struct {
	UserId     string  `json:"user_id"`
	Name       string  `json:"name"`
	Balance    float32 `json:"balance"`
	WinBalance float32 `json:"win_balance"`
	Win        int     `json:"win"`
	Lose       int     `json:"lose"`
	Push       int     `json:"push"`
	Blackjack  int     `json:"blackjack"`
}

type GameConfig struct {
	DeckNumber    int  `json:"deck_number"`
	DeckAmount    int  `json:"deck_amount"`
	CanDoubleDOwn bool `json:"can_double_down"`
	DealerSoft    int  `json:"dealer_soft"`
	MaxSplit      int  `json:"max_split"`
}

// MessageCode is enumarete all message types
type MessageCode int

// MessageCode is enumarete all message types
const (
	EEvent    MessageCode = iota + 0
	EMessage              // 1
	ERegister             // 2
	EPlayGame             // 3

)

var MessageCodes = []MessageCode{
	EEvent,
	ERegister,
	EPlayGame,
	EMessage,
}
