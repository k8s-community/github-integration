package models

import "time"

//go:generate reform

//reform:builds
type Build struct {
	ID         int64  `reform:"id,pk"`
	UUID       string `reform:"uuid"`
	Username   string `reform:"username"`
	Repository string `reform:"repository"`
	Commit     string `reform:"commit"`
	Passed     bool   `reform:"passed"`
	Log        string `reform:"log"`

	CreatedAt time.Time `reform:"created_at"`
	UpdatedAt time.Time `reform:"updated_at"`
}

// BeforeInsert set CreatedAt and UpdatedAt.
func (b *Build) BeforeInsert() error {
	b.CreatedAt = time.Now().UTC().Truncate(time.Second)
	b.UpdatedAt = b.CreatedAt
	return nil
}

// BeforeUpdate set UpdatedAt.
func (b *Build) BeforeUpdate() error {
	b.UpdatedAt = time.Now().UTC().Truncate(time.Second)
	return nil
}
