package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	//"MyPackges/crypt"

	"github.com/maxchagin/go-memorycache-example"
)

func Decrypt(CriptResult interface{}, Key []byte, Nonce []byte) []byte {

	CriptText, _ := CriptResult.([]byte)

	block, err := aes.NewCipher(Key)
	if err != nil {
		panic(err.Error())
	}

	Aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	Result, err := Aesgcm.Open(nil, Nonce, CriptText, nil)
	if err != nil {
		panic(err.Error())
	}

	return Result
}

func Encrypt(Result []byte) ([]byte, []byte, []byte) {
	Key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")

	Block, err := aes.NewCipher(Key)
	if err != nil {
		panic(err.Error())
	}

	Nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, Nonce)
	if err != nil {
		panic(err.Error())
	}

	Aesgcm, err := cipher.NewGCM(Block)
	if err != nil {
		panic(err.Error())
	}

	Ciphertext := Aesgcm.Seal(nil, Nonce, Result, nil)

	return Ciphertext, Key, Nonce
}

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func Content() []byte {
	xmlFile, err := getXML("https://iss.moex.com/iss/history/engines/stock/markets/shares/boards/tqbr/securities.xml")
	if err != nil {
		fmt.Println(err)
	}

	data := &Courses{}
	err = xml.Unmarshal(xmlFile, data)
	if nil != err {
		fmt.Println("Error unmarshalling from XML", err)
		panic(err)
	}

	Result, err := json.MarshalIndent(data.Data[0], "", "\t")
	if nil != err {
		fmt.Println("Error marshalling to JSON", err)
		panic(err)
	}

	return Result
}

func main() {

	Cache := memorycache.New(5*time.Minute, 10*time.Minute)

	var Key, Nonce []byte

	http.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {

		CriptResult, err := Cache.Get("myKey")

		if !err {
			CriptResult, Key, Nonce = Encrypt(Content())
			Cache.Set("myKey", CriptResult, 5*time.Minute)
		}

		fmt.Fprintf(w, "%s\n", Decrypt(CriptResult, Key, Nonce))

	})
	http.ListenAndServe(":8080", nil)
}
