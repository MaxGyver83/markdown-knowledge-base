package api

type Storage interface {
	Save(filename string, content string) error
	Load(filename string) (string, error)
	Delete(filename string) error
}
