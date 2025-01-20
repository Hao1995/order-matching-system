package requests

type CreateRequest struct {
	Symbol   string  `form:"symbol" binding:"required"`
	Type     string  `form:"side" binding:"required,oneof=Buy Sell"`
	Price    float64 `form:"price" binding:"required,gt=0"`
	Quantity int64   `form:"quantity" binding:"required,gt=0"`
}
