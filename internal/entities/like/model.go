package like

type ThreadLikeRequest struct {
	ThreadID int  `json:"thread_id" validate:"gt=0"`
	UserID   int  `json:"user_id" validate:"gt=0"`
	Like     bool `json:"like"`
}

type CommentLikeRequest struct {
	CommentID int  `json:"comment_id" validate:"gt=0"`
	UserID    int  `json:"user_id" validate:"gt=0"`
	Like      bool `json:"like"`
}
