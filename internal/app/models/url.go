package models

type URLRepository interface {
	Get(short string) (origin string)
	Set(short, origin string)
}
