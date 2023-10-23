package configs

import (
	"errors"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type Env map[string]any

func NewDatabaseEnv() (Env, error) {
	var env Env = make(Env)
	// check DB_HOST environment variable
	if val, found := os.LookupEnv("DB_HOST"); !found {
		return nil, errors.New("[err]: DB_HOST not found")
	} else {
		env["DB_HOST"] = val
	}
	// check DB_PORT environment variable
	if val, found := os.LookupEnv("DB_PORT"); !found {
		return nil, errors.New("[err]: DB_PORT not found")
	} else {
		port, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.New("[err]: DB_PORT is not a valid number")
		}
		env["DB_PORT"] = port
	}
	// check DB_USER environment variable
	if val, found := os.LookupEnv("DB_USER"); !found {
		return nil, errors.New("[err]: DB_USER not found")
	} else {
		env["DB_USER"] = val
	}
	// check DB_PASS environment variable
	if val, found := os.LookupEnv("DB_PASS"); !found {
		return nil, errors.New("[err]: DB_PASS not found")
	} else {
		env["DB_PASS"] = val
	}
	// check DB_NAME environment variable
	if val, found := os.LookupEnv("DB_NAME"); !found {
		return nil, errors.New("[err]: DB_NAME not found")
	} else {
		env["DB_NAME"] = val
	}
	return env, nil
}

func NewJwtEnv() (Env, error) {
	var env Env = make(Env)
	// check SECRET_KEY environment variable
	if val, found := os.LookupEnv("SECRET_KEY"); !found {
		return nil, errors.New("[err]: SECRET_KEY not found")
	} else {
		env["SECRET_KEY"] = val
	}
	return env, nil
}
