package domain

import "errors"

var (
	ErrTripNotFound                     = errors.New("trip not found")
	ErrForbidden                        = errors.New("forbidden")
	ErrConflict                         = errors.New("conflict")
	ErrAlreadyExists                    = errors.New("already exists")
	ErrNotAllowedCurrentStatusToPublish = errors.New("current status isn't allowed to publish")
	ErrClientNotDriver                  = errors.New("client isn't driver of this trip")
	ErrStatusIsPublishedAlready         = errors.New("status is published already")
)
