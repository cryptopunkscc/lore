package pathcache

import (
	"github.com/cryptopunkscc/lore/id"
	"time"
)

type Entry struct {
	Path       string
	ID         id.ID
	VerifiedAt time.Time
}
