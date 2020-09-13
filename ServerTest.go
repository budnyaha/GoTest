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
	"sync"
	"time"

	"github.com/maxchagin/go-memorycache-example"
)

var (
	LastResult *Courses
	rw         sync.RWMutex
)

type encryptReq struct {
	Text string `json:"text"`
	Key  string `json:"key"`
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
// func Decrypt(CriptResult interface{}) []byte {

// 	CriptText, _ := CriptResult.([]byte)

// 	block, err := aes.NewCipher(Key)
// 	if err != nil {
// 		panic(err)
// 	}

// 	Aesgcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		panic(err)
// 	}

// 	Result, err := Aesgcm.Open(nil, Nonce, CriptText, nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return Result
// }

// Метод шифровки
func Encrypt(Result []byte, Key []byte) []byte {
	Block, err := aes.NewCipher(Key)
	if err != nil {
		panic(err)
	}

	Nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, Nonce)
	if err != nil {
		panic(err)
	}

	Aesgcm, err := cipher.NewGCM(Block)
	if err != nil {
		panic(err)
	}

	Ciphertext := Aesgcm.Seal(nil, Nonce, Result, nil)

	return Ciphertext
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

		Cache.Set("myKey", output, 5*time.Minute)

	})

	http.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")
		// Result, err := json.MarshalIndent(LastResult.Data[0], "", "\t")

		// if nil != err {
		// 	fmt.Println("Error marshalling to JSON", err)
		// 	return
		// }
		CriptResult, err := Cache.Get("myKey")
		if !err {
			fmt.Fprintf(w, "ошибка")
		}
		if CriptResult == nil {
			fmt.Fprintf(w, "Хуй")
		} else {
			w.Write(CriptResult.([]byte))
		}
	})
	http.ListenAndServe(":8080", nil)
}
