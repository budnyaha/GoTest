package main

// Courses is export on hx
type Courses struct {
	Row []struct {
		Boardid                 string `xml:"BOARDID,attr" json:"Boardid"`
		TradeDate               string `xml:"TRADEDATE,attr" json:"TradeDate"`
		ShortName               string `xml:"SHORTNAME,attr" json:"ShortName"`
		SecId                   string `xml:"SECID,attr" json:"SecId"`
		NumTrades               string `xml:"NUMTRADES,attr" json:"NumTrades"`
		Value                   string `xml:"VALUE,attr" json:"Value"`
		Open                    string `xml:"OPEN,attr" json:"Open"`
		Low                     string `xml:"LOW,attr" json:"Low"`
		High                    string `xml:"HIGH,attr" json:"Hight"`
		LegalClosePrice         string `xml:"LEGALCLOSEPRICE,attr" json:"LegalClosePrice"`
		WapPrise                string `xml:"WAPRICE,attr" json:"WapPrise"`
		Close                   string `xml:"CLOSE,attr" json:"Close"`
		Volume                  string `xml:"VOLUME,attr" json:"Volume"`
		MarketPrice2            string `xml:"MARKETPRICE2,attr" json:"MarketPrice2"`
		MarketPrice3            string `xml:"MARKETPRICE3,attr" json:"MarketPrice3"`
		AdmittedQuote           string `xml:"ADMITTEDQUOTE,attr" json:"AdmittedQuote"`
		Mp2Valtrd               string `xml:"MP2VALTRD,attr" json:"Mp2Valtrd"`
		MarketPrice3TradesValue string `xml:"MARKETPRICE3TRADESVALUE,attr" json:"MarketPrice3TradesValue"`
		AdmittedValue           string `xml:"ADMITTEDVALUE,attr" json:"AdmittedValue"`
		Waval                   string `xml:"WAVAL,attr" json:"Waval"`
		TradingSession          string `xml:"TRADINGSESSION,attr" json:"TradingSession"`
	} `xml:"data>rows>row" json:"Rows"`
}
