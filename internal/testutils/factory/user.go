package factory

import (
	"github.com/brice-74/golang-base-api/internal/domains/user"
	"github.com/ventu-io/go-shortid"
)

func (f Factory) CreateUserAccount(props *user.User) *user.User {
	model := user.Model{DB: f.DB}

	u := &user.User{}
	if props != nil {
		u = props
	}

	if u.Email == "" {
		u.Email = f.faker.Internet().Email()
	}

	if u.Password == "" {
		u.Password = f.faker.Internet().Password()
	}

	if len(u.Roles) == 0 {
		u.Roles = user.Roles{user.RoleAnonymous}
	}

	if u.ProfilName == "" {
		u.ProfilName = f.faker.Beer().Name()
	}

	if u.ShortId == "" {
		u.ShortId = shortid.MustGenerate()
	}

	if err := model.InsertRegisteredUserAccount(u); err != nil {
		f.T.Fatalf("error during user factory insertion: %s", err)
	}

	return u
}

func (f Factory) CreateUserSession(props *user.Session) *user.Session {
	model := user.Model{DB: f.DB}

	s := &user.Session{}
	if props != nil {
		s = props
	}

	if s.ID == "" {
		s.ID = f.faker.UUID().V4()
	}

	if s.IP == "" {
		s.IP = f.faker.Internet().Ipv4()
	}

	if s.Name == "" {
		s.Name = f.faker.UserAgent().UserAgent()
	}

	if s.Location == "" {
		s.Location = f.faker.Address().Country()
	}

	if s.UserID == "" {
		s.UserID = f.CreateUserAccount(nil).ID
	}

	if err := model.InsertOrUpdateUserSession(s); err != nil {
		f.T.Fatalf("error during session factory insertion: %s", err)
	}

	return s
}
