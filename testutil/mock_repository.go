package testutil

import "github.com/aaronaaeng/chat.connor.fun/db"

type MockRepository struct {
	UsersRepo *MockUserRepository
	RolesRepo *MockRolesRepository
	RoomsRepo *MockRoomsRepository
	MessagesRepo *MockMessagesRepository
	VerificationsRepo *MockVerificationsRepo
}

func (r *MockRepository) NewFromSource(source db.DataSource) db.Repository {
	return r
}

func (r *MockRepository) Messages() db.MessagesRepository {
	return r.MessagesRepo
}

func (r *MockRepository) Users() db.UserRepository {
	return r.UsersRepo
}

func (r *MockRepository) Roles() db.RolesRepository {
	return r.RolesRepo
}

func (r *MockRepository) Rooms() db.RoomsRepository {
	return r.RoomsRepo
}

func (r *MockRepository) Verifications() db.VerificationCodeRepository {
	return r.VerificationsRepo
}
