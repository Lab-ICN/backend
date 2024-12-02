package logging

import (
	"go.uber.org/zap"
)

func NewDevelopment() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func NewProduction() (*zap.Logger, error) {
	return zap.NewProduction()
}
