package types

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ConnectionId uuid.UUID

type ConnectionAccessId uuid.UUID

type JSONStructType = datatypes.JSONType[map[string]any]
