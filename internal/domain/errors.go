package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")
var ErrExpiredDate = errors.New("date is expired")
var ErrAvailableSeatsNotInteger = errors.New("attribute 'availableSeats' doesn't match to integer type")
var ErrAvailableSeatsInvalidValue = errors.New("attribute 'availableSeats' has value less 1")
var ErrAvailableSeatsFromPointToPointEqual = errors.New("fromPoint and toPoint have the same value")
