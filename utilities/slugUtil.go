package utilities

import (
	"strings"
)

type ISlugUtil interface {
	GetSlug(name string) string
	GetName(slug string) string
}

type SlugUtil struct {

}

func NewSlugUtil() ISlugUtil {
	util := SlugUtil{}
	return &util
}

func (util SlugUtil) GetSlug(name string) string {
	slug := strings.Replace(name, " ", "-", -1)
	return slug
}

func (util SlugUtil) GetName(slug string) string {
	name := strings.Replace(slug, "-", " ", -1)
	return name
}