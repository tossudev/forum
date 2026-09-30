package search

import (
//here goes the repo layers of the other packages
	"database/sql"
	"forum/internal/pagination"
)


type SearchRepository struct {
	
}



func (sr *SearchRepository) Search(input string) ([]Comment, []Thread) {
	return nil, nil
}


