package search

import (
	"fmt"
	"net/http"

	"forum/internal/errs"
	"github.com/go-playground/validator/v10"
)

type SearchHandler struct {
	threadService	*thread.ThreadService
	commentService 	*comment.CommentService
	renderer 		*render.Renderer
}

func (sh *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	//the input will come from the body somewhere, (from the textbox it goes somewhere and then we fetch it here)
	// get the session so we can get the user. What do we need the user for? Everyone can do a search, maybe just skip the user then

	input := r.PathValue("input")
	newSearch := SearchResult{
		Query: input,
	}

	//if length of input is x amount then give message and error, search input too long
	
	// in the repo layer the query will be.... what? need to figure that out first. 




}




