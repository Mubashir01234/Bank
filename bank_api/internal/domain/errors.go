package domain

import "errors"

var (
	ErrResourceNotFound    = errors.New("resource not found")
	ErrInGhostMode         = errors.New("resource is in ghost mode")
	ErrInvalidParameter    = errors.New("invalid parameter")
	ErrSessionNotPresent   = errors.New("session not present")
	ErrSessionNotActive    = errors.New("session is not active")
	ErrInvalidAssetType    = errors.New("invalid asset type")
	ErrInvalidAssetSubType = errors.New("invalid asset sub type")
	ErrInvalidAssetID      = errors.New("invalid asset ID")
	ErrResourceGone        = errors.New("resource gone")
)
