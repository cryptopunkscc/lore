package story

type HeaderRepo interface {
	Add(id string, header *Header) error
	Remove(id string) error

	SetStoryType(id string, typ string) error
	GetStoryType(id string) (string, error)
	SetStoryRels(id string, rels []string) error
	GetStoryRels(id string) ([]string, error)
}
