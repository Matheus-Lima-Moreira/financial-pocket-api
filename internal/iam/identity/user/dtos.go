package user

type ListRequest struct {
	Page int `form:"page" binding:"required,min=1"`
}

type DetailsRequest struct {
	ID string `uri:"id" binding:"required,min=1"`
}
