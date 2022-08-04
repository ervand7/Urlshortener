package models

type UrlRepository interface {
	Get(short string) (origin string)
	Set(short, origin string)
}

type Storage map[string]string

func (s Storage) Get(short string) (origin string) {
	origin, exist := s[short]
	if !exist {
		return ""
	}
	return origin
}

func (s Storage) Set(short, origin string) {
	s[short] = origin
}
