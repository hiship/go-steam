package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hiship/go-steam"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	timeTip, err := steam.GetTimeTip()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Time tip: %#v\n", timeTip)

	timeDiff := time.Duration(timeTip.Time - time.Now().Unix())
	session := steam.NewSession(&http.Client{}, "")
	if err := session.Login(os.Getenv("steamAccount"), os.Getenv("steamPassword"), os.Getenv("steamSharedSecret"), timeDiff); err != nil {
		log.Fatal(err)
	}
	log.Print("Login successful")

	key, err := session.GetWebAPIKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Key: ", key)

	identitySecret := os.Getenv("steamIdentitySecret")
	confirmations, err := session.GetConfirmations(identitySecret, time.Now().Add(timeDiff).Unix())
	if err != nil {
		log.Fatal(err)
	}

	for i := range confirmations {
		c := confirmations[i]
		log.Printf("Confirmation ID: %s, Creator: %s\n", c.ID, c.Creator)
		log.Printf("-> ID %s\n", c.ID)
		log.Printf("-> Type %d\n", c.Type)
		log.Printf("-> CreationTime %d\n", c.CreationTime)
		log.Printf("-> Nonce %s\n", c.Nonce)

		err = session.AnswerConfirmation(c, identitySecret, "allow", time.Now().Add(timeDiff).Unix())
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Accepted %d\n", c.ID)
	}

	log.Println("Bye!")
}
