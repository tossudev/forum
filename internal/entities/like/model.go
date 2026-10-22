package like

type ThreadLikeRequest struct {
	ThreadID int  `json:"id" validate:"gt=0"`
	UserID   int  `validate:"gt=0"`
	Like     bool `json:"like"`
}

type CommentLikeRequest struct {
	CommentID int  `json:"id" validate:"gt=0"`
	UserID    int  `validate:"gt=0"`
	Like      bool `json:"like"`
}
