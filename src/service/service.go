package service

import (
	"github.com/kalpit-sharma-dev/chat-service/src/repository"
)

type ChatServiceImpl struct {
	DatabaseRepository repository.DatabaseRepository
}

// Login implements ChatService.
func (c *ChatServiceImpl) Login() {
	panic("unimplemented")
}

// RegisterUser implements ChatService.
func (c *ChatServiceImpl) RegisterUser() {
	panic("unimplemented")
}

// ServeWs implements ChatService.
func (c *ChatServiceImpl) ServeWs() {
	panic("unimplemented")
}

// UploadFile implements ChatService.
func (c *ChatServiceImpl) UploadFile() {
	panic("unimplemented")
}

// VerifyUser implements ChatService.
func (c *ChatServiceImpl) VerifyUser() {
	panic("unimplemented")
}

// func NewChat1Service(repository repository.DatabaseRepository) ChatService {
// 	return &ChatServiceImpl{DatabaseRepository: repository}
// }
