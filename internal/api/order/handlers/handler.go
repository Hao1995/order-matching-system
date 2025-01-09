package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Hao1995/order-matching-system/internal/api/order/usecases"
	"github.com/Hao1995/order-matching-system/pkg/models/events"
)

var (
	now = func() time.Time {
		return time.Now()
	}
)

type OrderCreateRequest struct {
	Symbol   string  `form:"symbol" binding:"required"`
	Side     string  `form:"side" binding:"required,oneof=buy sell"`
	Price    float64 `form:"price" binding:"required,gt=0"`
	Quantity int64   `form:"quantity" binding:"required,gt=0"`
}

type OrderCancelRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type Handler struct {
	producer usecases.Producer
	topic    string
}

// CreateOrder handles the creation of a new order.
func (hlr *Handler) CreateOrder(c *gin.Context) {
	var request OrderCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
		return
	}

	id := uuid.New().String()
	side, err := events.ParseSide(request.Side)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid order side"})
		return
	}

	orderEvent := events.OrderEvent{
		EventType: events.OrderEventTypeCREATE,
		Data: events.OrderCreateEvent{
			ID:                id,
			Symbol:            request.Symbol,
			Side:              side,
			Price:             request.Price,
			Quantity:          request.Quantity,
			RemainingQuantity: request.Quantity,
			CanceledQuantity:  0,
			CreatedAt:         now(),
			UpdatedAt:         now(),
		},
	}

	err = hlr.producer.Publish(c.Request.Context(), hlr.topic, &orderEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed create a Create Order request"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "The Create Order request has been accepted"})
}

// CancelOrder handles the cancellation of an order.
func (hlr *Handler) CancelOrder(c *gin.Context) {
	var request OrderCancelRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request data"})
		return
	}

	orderEvent := events.OrderEvent{
		EventType: events.OrderEventTypeCANCEL,
		Data: events.OrderCancelEvent{
			ID: request.ID,
		},
	}

	err := hlr.producer.Publish(c.Request.Context(), hlr.topic, &orderEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed create a Cancel Order request"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "The Cancel Order request has been accepted"})
}
