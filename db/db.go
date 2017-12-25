package db

type Repository interface{
	Create(interface{}) error
	Update(interface{}) error
	GetAll() ([]*interface{}, error)
	GetById(id int) (interface{}, error)
	Delete(interface{}) error
}