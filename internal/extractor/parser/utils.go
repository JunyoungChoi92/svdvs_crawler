package parser

import (
	"log"
	"net/url"
)

// combine string slice
func CombineSlice(sli1 []string, sli2 []string) []string {
	sli1 = append(sli1, sli2...)
	return checkDuplicate(sli1)
}

// duplicated check string slice
func checkDuplicate(src []string) []string {
	dupMap := map[string]bool{}
	var dupSlice []string
	for _, d := range src {
		if dupMap[d] {
			continue
		}
		dupMap[d] = true
		dupSlice = append(dupSlice, d)
	}
	return dupSlice
}

func ExtractDomain(fullUrl string) (string, error) {
	u, err := url.Parse(fullUrl)
	if err != nil {
		log.Println(err)
		return "", err
	}
	domain := u.Scheme + "://" + u.Hostname()

	return domain, nil
}
