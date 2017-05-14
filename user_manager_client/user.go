package user_manager_client

const (
	userUrlStr = "/user"
)

type UserService struct {
	client *Client
}

func (u *UserService) Create(user User) error {
	req, err := u.client.NewRequest(postMethod, userUrlStr, u)
	if err != nil {
		return err
	}

	_, err = u.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

type User struct {
	Name *string `json:"name"`
}

func NewUser(username string) *User {
	return &User{
		Name: &username,
	}
}
