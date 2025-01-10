package requests

type CancelRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}
