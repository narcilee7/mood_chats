package services

type AIProvider interface {
	Chat(promt string, history []string) (string, error)
}