package order

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Hao1995/order-matching-system/internal/api/order/requests"
	"github.com/Hao1995/order-matching-system/internal/common/models/events"
	"github.com/Hao1995/order-matching-system/pkg/logger"
	"github.com/Hao1995/order-matching-system/pkg/mqkit"
)

var (
	now = func() time.Time {
		return time.Now()
	}
)

type Handler struct {
	producer mqkit.Producer
	topic    string
}

func NewHandler(p mqkit.Producer, topic string) *Handler {
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

	id := uuid.NewString()
	orderType, err := events.ParseOrderType(request.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid order orderType"})
		return
	}

	orderEvent := events.Event{
		EventType: events.EventTypeCreateOrder,
		Data: events.OrderEvent{
			ID:        id,
			Symbol:    request.Symbol,
			Type:      orderType,
			Price:     request.Price,
			Quantity:  request.Quantity,
			CreatedAt: now(),
		},
	}

	val, err := json.Marshal(orderEvent)
	if err != nil {
		logger.Error("failed to json marshal event", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to json marshal"})
	}

	err = hlr.producer.Publish(c.Request.Context(), val)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed create a Create Order request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "The Create Order request has been accepted"})
}

// Cancel handles the cancellation of an order.
func (hlr *Handler) Cancel(c *gin.Context) {
	var request requests.CancelRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request data"})
		return
	}

	event := events.Event{
		EventType: events.EventTypeCancelOrder,
		Data: events.OrderEvent{
			ID: request.ID,
		},
	}

	val, err := json.Marshal(event)
	if err != nil {
		logger.Error("failed to json marshal event", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to json marshal"})
	}

	err = hlr.producer.Publish(c.Request.Context(), val)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed create a Cancel Order request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "The Cancel Order request has been accepted"})
}
