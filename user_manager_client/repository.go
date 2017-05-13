package user_manager_client

const (
	repositoryURL = "/repository"
)

type RepositoryService struct {
	client *Client
}

func (r *RepositoryService) Create(rep Repository) error {
	req, err := r.client.NewRequest(postMethod, repositoryURL, r)
	if err != nil {
		return err
	}

	_, err = r.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryService) Delete(rep Repository) error {
	urlStr := repositoryURL + "/" + *rep.Username + "_" + *rep.Name

	req, err := r.client.NewRequest(deleteMethod, urlStr, nil)
	if err != nil {
		return err
	}

	_, err = r.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

type Repository struct {
	Username *string
	Name     *string
}

func NewRepository(username string, repName string) *Repository {
	return &Repository{
		Username: &username,
		Name:     &repName,
	}
}
