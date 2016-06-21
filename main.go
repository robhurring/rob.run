package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
)

const envBase = "REDIRECT"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log := logger.New(iris.Logger)
	root := iris.Party("", log)

	root.Get("/vim", redirectHandler)

	iris.Listen(":8080")
}

func redirectHandler(c *iris.Context) {
	to, err := redirectForPath(c.PathString())

	if err != nil {
		c.Text(iris.StatusInternalServerError, err.Error())
	} else {
		c.Redirect(to, iris.StatusTemporaryRedirect)
	}
}

func pathToEnv(path string) string {
	root := strings.TrimLeft(path, "/")
	root = strings.Replace(root, "/", "_", -1)
	root = envBase + "_" + root

	return strings.ToUpper(root)
}

func redirectForPath(path string) (string, error) {
	key := pathToEnv(path)
	val := os.Getenv(key)

	if val == "" {
		return "", errors.New("Unknown key: " + key)
	}

	return val, nil
}
