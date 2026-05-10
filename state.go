package main

import (
	"github.com/Dr3iundZwanzig/BlogAggregator/internal/config"
	"github.com/Dr3iundZwanzig/BlogAggregator/internal/database"
)

type state struct {
	config *config.Config
	db     *database.Queries
}
