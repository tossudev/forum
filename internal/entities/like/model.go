package like

type ThreadLikeRequest struct {
	ThreadID int `json:"id"`
	UserID   int
	Like     bool `json:"like"`
}

type CommentLikeRequest struct {
	CommentID int `json:"id"`
	UserID    int
	Like      bool `json:"like"`
}
