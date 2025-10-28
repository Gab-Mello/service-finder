package ports

type ProviderDirectory interface {
	GetNameByID(providerID string) (string, error)
}
