package db


type RepositoryImpl struct {
	Source DataSource
	MessagesRepo MessagesRepository
	UsersRepo UserRepository
	RolesRepo RolesRepository
	RoomsRepo RoomsRepository
	VerificationsRepo VerificationCodeRepository
}

func (r *RepositoryImpl) NewFromSource(source DataSource) Repository {
	return &RepositoryImpl{
		Source: source,
		MessagesRepo: 		r.MessagesRepo.NewFromSource(source),
		UsersRepo: 			r.UsersRepo.NewFromSource(source),
		RolesRepo: 			r.RolesRepo.NewFromSource(source),
		RoomsRepo: 			r.RoomsRepo.NewFromSource(source),
		VerificationsRepo: 	r.VerificationsRepo.NewFromSource(source),
	}
}

func (r *RepositoryImpl) Messages() MessagesRepository {
	return r.MessagesRepo
}

func (r *RepositoryImpl) Users() UserRepository {
	return r.UsersRepo
}

func (r *RepositoryImpl) Roles() RolesRepository {
	return r.RolesRepo
}

func (r *RepositoryImpl) Rooms() RoomsRepository {
	return r.RoomsRepo
}

func (r *RepositoryImpl) Verifications() VerificationCodeRepository {
	return r.VerificationsRepo
}