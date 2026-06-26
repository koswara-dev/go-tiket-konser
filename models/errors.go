package models

import "errors"

var (
	ErrConcertAlreadyExists = errors.New("concert already exists")
	ErrConcertNotFound      = errors.New("concert not found")
	ErrConcertDeleted       = errors.New("concert deleted")
	ErrConcertNotActive     = errors.New("concert not active")
	ErrInvalidDate          = errors.New("invalid date")
	ErrInvalidVenue         = errors.New("invalid venue")
	ErrInvalidStatus        = errors.New("invalid status")
	ErrBookingNotFound      = errors.New("booking not found")
	ErrUnauthorize          = errors.New("unauthorized")
	ErrUserNotFound         = errors.New("user not found")
	ErrCustomerNotFound     = errors.New("customer not found")
)
