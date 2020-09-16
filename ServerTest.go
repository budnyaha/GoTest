package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"

	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/maxchagin/go-memorycache-example"
)

var (
	LastResult *Courses
	rw         sync.RWMutex
	Nonce, _   = hex.DecodeString("64a9433eae7ccceee2fc0eda")
)

type encryptReq struct {
	Text string `json:Text`
	Key  string `json:"Key"`
}

// Структура xml файла
type Courses struct {
	Data []struct {
		Rows []struct {
			Row []struct {
				Boardid                 string  `xml:"BOARDID,attr"`
				TradeDate               string  `xml:"TRADEDATE,attr"`
				ShortName               string  `xml:"SHORTNAME,attr"`
				SecId                   string  `xml:"SECID,attr"`
				NumTrades               int64   `xml:"NUMTRADES,attr"`
				Value                   float64 `xml:"VALUE,attr"`
				Open                    float64 `xml:"OPEN,attr"`
				Low                     float64 `xml:"LOW,attr"`
				High                    float64 `xml:"HIGH,attr"`
				LegalClosePrice         float64 `xml:"LEGALCLOSEPRICE,attr"`
				WapPrise                float64 `xml:"WAPRICE,attr"`
				Close                   float64 `xml:"CLOSE,attr"`
				Volume                  int64   `xml:"VOLUME,attr"`
				MarketPrice2            float64 `xml:"MARKETPRICE2,attr"`
				MarketPrice3            float64 `xml:"MARKETPRICE3,attr"`
				AdmittedQuote           float64 `xml:"ADMITTEDQUOTE,attr"`
				Mp2Valtrd               float64 `xml:"MP2VALTRD,attr"`
				MarketPrice3TradesValue float64 `xml:"MARKETPRICE3TRADESVALUE,attr"`
				AdmittedValue           float64 `xml:"ADMITTEDVALUE,attr"`
				Waval                   int64   `xml:"WAVAL,attr"`
				TradingSession          int64   `xml:"TRADINGSESSION,attr"`
			} `xml:"row" json:"Row"`
		} `xml:"rows"`
	} `xml:"data"`
}

// Метод дешифровки
func Decrypt(cryptoText []byte, key []byte) []byte {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	if len(cryptoText) < aes.BlockSize {
		panic("cryptoText too short")
	}
	iv := cryptoText[:aes.BlockSize]
	cryptoText = cryptoText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(cryptoText, cryptoText)

	return cryptoText
}

// Метод шифровки
func Encrypt(Result []byte, Key []byte) []byte {
	plaintext := []byte(Result)

	block, err := aes.NewCipher(Key)
	if err != nil {
		panic(err)
	}

	Сiphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := Сiphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(Сiphertext[aes.BlockSize:], plaintext)

	return Сiphertext
}

// Метод получения данных с сайта
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

// Метод обработки данных в Json
func Content() {
	for {
		xmlFile, err := getXML("https://iss.moex.com/iss/history/engines/stock/markets/shares/boards/tqbr/securities.xml")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		data := &Courses{}
		err = xml.Unmarshal(xmlFile, data)
		if nil != err {
			fmt.Println("Error unmarshalling from XML", err)
			panic(err)
		}

		rw.Lock()
		LastResult = data
		rw.Unlock()

		<-time.After(5 * time.Second)
	}
}

func main() {

	Cache := memorycache.New(5*time.Minute, 10*time.Minute)

	go Content()

	http.HandleFunc("/encrypt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if r.Body == nil {
			w.WriteHeader(500)
			return
		}

		var req encryptReq
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(data, &req)

		key, err := hex.DecodeString(req.Key)
		if err != nil {
			panic(err)
		}
		output := Encrypt([]byte(req.Text), key)
		w.Write(output)
	})
	http.HandleFunc("/decrypt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if r.Body == nil {
			w.WriteHeader(500)
			return
		}

		var req encryptReq
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(data, &req)

		key, err := hex.DecodeString(req.Key)
		if err != nil {
			panic(err)
		}
		output := Encrypt([]byte(req.Text), key)
		w.Write(output)

	})

	http.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")
		Result, err := json.MarshalIndent(LastResult.Data[0], "", "\t")
		if nil != err {
			fmt.Println("Error marshalling to JSON", err)
			w.WriteHeader(500)
			return
		}
		CriptResult, erro := Cache.Get("myKey")

		if !erro {
			Cache.Set("myKey", Result, 5*time.Minute)
			CriptResult, _ = Cache.Get("myKey")
		}

		fmt.Fprintf(w, "%s\n", CriptResult)

	})
	http.ListenAndServe(":8080", nil)
}
