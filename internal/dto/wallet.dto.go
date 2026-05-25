package dto

type CreateTransferRequest struct {
	ReceiverID int    `json:"receiver_id" binding:"required,gt=0"`
	Amount     int64  `json:"amount" binding:"required,gt=0"`
	Pin        string `json:"pin" binding:"required,len=6,numeric"`
	Note       string `json:"note"`
}

type CreateTopUpRequest struct {
	PaymentMethodID int    `json:"payment_method_id" binding:"required,gt=0"`
	Amount          int64  `json:"amount" binding:"required,gt=0"`
	Note            string `json:"note"`
}

type CreateTransferResponse struct {
	TransactionID int64  `json:"transaction_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type CreateTopUpResponse struct {
	TransactionID int64  `json:"transaction_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
	Amount        int64  `json:"amount"`
	TaxPercent    int64  `json:"tax_percent"`
	TaxAmount     int64  `json:"tax_amount"`
	AdminFee      int64  `json:"admin_fee"`
	Total         int64  `json:"total"`
}
