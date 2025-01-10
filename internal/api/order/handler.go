package order

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Hao1995/order-matching-system/internal/api/order/requests"
	"github.com/Hao1995/order-matching-system/internal/common/models/events"
)

var (
	now = func() time.Time {
		return time.Now()
	}
)

type Handler struct {
	producer Producer
	topic    string
}

func NewHandler(p Producer, topic string) *Handler {
	return &Handler{
		producer: p,
		topic:    topic,
	}
}

// Create handles the creation of a new order.
func (hlr *Handler) Create(c *gin.Context) {
	var request requests.CreateRequest
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

// Cancel handles the cancellation of an order.
func (hlr *Handler) Cancel(c *gin.Context) {
	var request requests.CancelRequest

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
