package handlers

import (
	"github.com/sirupsen/logrus"
)

type DebugResponder struct {
	base   Responder
	logger *logrus.Entry
}

func NewDebugResponder(r Responder, logger *logrus.Entry) *DebugResponder {

	return &DebugResponder{
		base:   r,
		logger: logger,
	}
}

func (d *DebugResponder) InvalidSession() ErrorResponse {

	d.logger.Debug("invalid session")

	return d.base.InvalidSession()
}

func (d *DebugResponder) Unauthorized() ErrorResponse {
	d.logger.Debug("unauthroized")

	return d.base.Unauthorized()
}

func (d *DebugResponder) InternalError(err error) ErrorResponse {

	d.logger.WithFields(logrus.Fields{
		"err": err,
	}).Debug("internal server error")

	return d.base.InternalError(err)
}

func (d *DebugResponder) OK(headers ...map[string]string) Response {

	return d.base.OK(headers...)
}
