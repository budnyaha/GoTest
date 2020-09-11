package main

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
