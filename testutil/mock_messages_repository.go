package testutil

import (
	"github.com/satori/go.uuid"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"errors"
)

type messageRepoFilter func(message *model.Message) bool

type MockMessagesRepository struct {
	Messages map[uuid.UUID]model.Message
}

func NewMockMessagesRepository() *MockMessagesRepository {
	return &MockMessagesRepository{Messages: map[uuid.UUID]model.Message{}}
}

func (r MockMessagesRepository) Add(message *model.Message) error {
	if _, ok := r.Messages[message.Id]; ok {
		return errors.New("duplicate entry")
	}
	r.Messages[message.Id] = *message
	return nil
}

func (r MockMessagesRepository) Update(id uuid.UUID, newText string) (*model.Message, error) {
	toEdit := r.Messages[id]
	toEdit.Text = newText
	r.Messages[id] = toEdit
	return &toEdit, nil
}

func (r MockMessagesRepository) GetById(id uuid.UUID) (*model.Message, error) {
	toReturn, ok := r.Messages[id]
	if !ok {
		return nil, nil
	}
	return &toReturn, nil
}

func (r MockMessagesRepository) getMessagesWithFilter(filter messageRepoFilter) []*model.Message {
	messages := make([]*model.Message, 0)
	for messId := range r.Messages {
		message := r.Messages[messId]
		if filter(&message) {
			messages = append(messages, &message)
		}
	}
	return messages
}

func (r MockMessagesRepository) GetByUserId(userId uuid.UUID) ([]*model.Message, error) {
	return r.getMessagesWithFilter(func(message *model.Message) bool {
		return message.Creator.Id == userId
	}), nil
}

func (r MockMessagesRepository) GetTopByUserId(userId uuid.UUID, count int) ([]*model.Message, error) {
	return r.getMessagesWithFilter(func(message *model.Message) bool {
		return message.Creator.Id == userId
	})[:count], nil
}

func (r MockMessagesRepository) GetByRoomId(roomId uuid.UUID) ([]*model.Message, error) {
	return r.getMessagesWithFilter(func(message *model.Message) bool {
		return message.Room.Id == roomId
	}), nil
}

func (r MockMessagesRepository) GetTopByRoom(roomId uuid.UUID, count int) ([]*model.Message, error) {
	return r.getMessagesWithFilter(func(message *model.Message) bool {
		return message.Room.Id == roomId
	})[:count], nil
}

func (r MockMessagesRepository) GetByUserAndRoom(userId uuid.UUID, roomId uuid.UUID) ([]*model.Message, error) {
	return r.getMessagesWithFilter(func(message *model.Message) bool {
		return message.Room.Id == roomId && message.Creator.Id == userId
	}), nil
}
func (r MockMessagesRepository) GetTopByUserAndRoom(userId uuid.UUID, roomId uuid.UUID, count int) ([]*model.Message, error) {
	return r.getMessagesWithFilter(func(message *model.Message) bool {
		return message.Room.Id == roomId && message.Creator.Id == userId
	})[:count], nil
}
