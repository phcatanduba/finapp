package pluggy

import "time"

type AuthResponse struct {
	APIKey string `json:"apiKey"`
}

type ConnectTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type Item struct {
	ID            string     `json:"id"`
	Status        string     `json:"status"`
	ConnectorID   int        `json:"connectorId"`
	Connector     Connector  `json:"connector"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	LastUpdatedAt *time.Time `json:"lastUpdatedAt,omitempty"`
}

type Connector struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AccountsResponse struct {
	Total   int       `json:"total"`
	Results []Account `json:"results"`
}

type Account struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Subtype     string   `json:"subtype"`
	Balance     float64  `json:"balance"`
	CreditData  *Credit  `json:"creditData,omitempty"`
	CurrencyCode string  `json:"currencyCode"`
}

type Credit struct {
	CreditLimit float64 `json:"creditLimit"`
}

type TransactionsResponse struct {
	Total   int           `json:"total"`
	Results []Transaction `json:"results"`
	Page    int           `json:"page"`
	// Pluggy uses cursor-based pagination
}

type Transaction struct {
	ID            string    `json:"id"`
	Description   string    `json:"description"`
	Amount        float64   `json:"amount"`
	Date          time.Time `json:"date"`
	Type          string    `json:"type"`
	AccountID     string    `json:"accountId"`
	Category      string    `json:"category,omitempty"`
	CategoryID    string    `json:"categoryId,omitempty"`
}

type WebhookRequest struct {
	Event      string `json:"event"`
	URL        string `json:"url"`
	Headers    map[string]string `json:"headers,omitempty"`
}
