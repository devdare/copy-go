package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {

	go server()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go watch(watcher)

	err = watcher.Add("src")
	if err != nil {
		log.Fatal(err)
	}

	<-done

}

func copy(filePath string, dst string) {

	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatalln(err)
		return
	}

	key := []byte("strongpasswordha")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(file))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], file)

	filename := dst + "/" + strings.Split(filePath, "/")[1]

	if err := ioutil.WriteFile(filename, ciphertext, 0777); err != nil {
		log.Fatalln(err)
	}

}

func watch(watcher *fsnotify.Watcher) {

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				copy(event.Name, "dst")
				log.Println("copied file:", event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}

}
