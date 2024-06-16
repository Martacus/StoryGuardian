package project

type Entity struct {
	Id          string
	Name        string
	ProjectId   string
	Description string
	CreatedAt   string
	Tags        []string
}
type EntityManager struct {
	ApplicationManager *ApplicationManager
}

func NewEntityManager(appManager *ApplicationManager) *EntityManager {
	return &EntityManager{
		ApplicationManager: appManager,
	}
}

func (e *EntityManager) GetEntities() {
	//e.ApplicationManager
}
