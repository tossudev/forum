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

type Like struct {
	Exists bool // Is a like/dislike reaction present at all
	Liked  bool // True = like, false = dislike
}
