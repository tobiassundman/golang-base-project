package logging

import "go.uber.org/zap"

// NewProductionLogger creates a new production logger.
func NewProductionLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}
