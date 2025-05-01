package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go" //nolint:depguard
)

// User stores user e-mails.
type User struct {
	Email string
}

var ErrInvalidJSON = errors.New("invalid json line")

// DomainStat stores domain statistics.
type DomainStat map[string]int

// GetDomainStat calculates domain statistics from input reader.
func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var (
		result  = make(DomainStat)
		err     error
		user    *User
		json    = jsoniter.ConfigCompatibleWithStandardLibrary
		scanner = bufio.NewScanner(r)
		suffix  = "." + domain
	)

	for i := 0; scanner.Scan(); i++ {
		if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, errors.Join(err, ErrInvalidJSON)
		}

		if strings.HasSuffix(user.Email, suffix) {
			idx := strings.LastIndex(user.Email, "@")
			if idx == -1 {
				continue
			}

			result[strings.ToLower(user.Email[idx+1:])]++
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
