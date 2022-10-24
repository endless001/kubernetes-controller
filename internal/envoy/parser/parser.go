package parser

import (
	"kubernetes-controller/internal/store"
)

type Parser struct {
	storer store.Storer
}

func NewParser(storer store.Storer) *Parser {
	return &Parser{
		storer: storer,
	}
}
