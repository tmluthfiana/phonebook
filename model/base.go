package model

type BaseModel interface {
	PreSave() error
}

type Model struct {
	BaseModel
}

func (m *Model) PreSave() error {
	return nil
}
