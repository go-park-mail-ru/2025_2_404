package metric

import (
	"time"

	"github.com/google/uuid"
)

type AdDetailID uuid.UUID
type SlotID		uuid.UUID

type Metric struct{
	SlotID 			SlotID
	AdDetailID			AdDetailID
	EventType		string
	CreateeTime		time.Time
}