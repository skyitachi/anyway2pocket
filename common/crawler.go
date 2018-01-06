package common

type Crawler interface {
	OnGetURL(url string)
	Start()
	CanNext(url string) bool
}
