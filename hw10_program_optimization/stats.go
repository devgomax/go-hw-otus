package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

func getUsers(r io.Reader) ([]User, error) {
	var (
		user  User
		users []User
	)

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	scanner := bufio.NewScanner(r)

	for i := 0; scanner.Scan(); i++ {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, errors.Join(err, ErrInvalidJSON)
		}

		users = append(users, user)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func countDomains(u []User, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.HasSuffix(user.Email, "."+domain) {
			idx := strings.LastIndex(user.Email, "@")
			if idx == -1 {
				continue
			}

			result[strings.ToLower(user.Email[idx+1:])]++
		}
	}

	return result, nil
}
