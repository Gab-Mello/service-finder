package user

type Repository interface{}

type memoryRepo struct{}

func NewRepository() Repository {
	return &memoryRepo{}
}
