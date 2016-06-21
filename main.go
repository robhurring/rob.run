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
const banner = `
                           _
                          | |
 ___  ___  ___   _ __ ___ | |__   _ __ _   _ _ __
/ __|/ _ \/ _ \ | '__/ _ \| '_ \ | '__| | | | '_ \
\__ \  __/  __/_| | | (_) | |_) || |  | |_| | | | |
|___/\___|\___(_)_|  \___/|_.__(_)_|   \__,_|_| |_|
																			 run rob run!

`

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	root := iris.Party("", logger.New(iris.Logger))
	root.Get("/", func(c *iris.Context) {
		c.Text(iris.StatusOK, banner)
	})
	root.Get("/vim", redirectHandler)

	iris.Listen("0.0.0.0:" + port)
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
