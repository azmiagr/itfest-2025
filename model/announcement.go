package model

import (
	"errors"
)

var ErrUserRecordNotFound = errors.New("Not Found data user")

type RequestAnnouncement struct {
	Message string `json:"message" binding:"required"`
}

type ResponseAnnouncement struct {
	AnnouncementID string `json:"id_announcement"`
	Message        string `json:"message_announcement"`
}
