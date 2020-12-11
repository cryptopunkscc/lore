package story

type StoryRepo interface {
	SetStoryType(id string, typ string) error
	GetStoryType(id string) (string, error)
	SetStoryRels(id string, rels []string) error
	GetStoryRels(id string) ([]string, error)
	Forget(id string) error
}
