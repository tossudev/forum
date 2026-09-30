package search

import (
	"fmt"
	"net/http"

	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/entities/thread"
	"forum/internal/entities/comment"

	"github.com/go-playground/validator/v10"
)

type SearchHandler struct {
	SearchService 	*SearchService
	Renderer 		*render.Renderer
	Validator		*validator.Validate
}

func NewHandler(service *SearchService, renderer *renderer.Renderer, validator *validator.Validate) *SearchService {
	return &SearchService{SearchService: service,
						  Renderer: renderer,
						  Validator: validator}
}

func (sh *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	query := r.URL.Query()
	pagination, err := pagination.Parse(query)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	input := r.URL.Query().Get("q")
	if input == "" || len(input) > 1000 {
		errs.WriteError(w, fmt.Errorf("%w: invalid search length", errs.ErrInvalidUserInput)
		return
	}



	comments, threads, err := sh.SearchService(ctx, input, pagination)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	data := SearchResult{
		Query: input,

	}
	// do the rendering stuff here mf


	//if length of input is x amount then give message and error, search input too long
	
	// in the repo layer the query will be.... what? need to figure that out first. 

// we need DATA to send to the frontend. The data is a struct of stuff the frontend needs in the render func at the end of the handler. 


}




