package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	OKCOIN_WEBSOCKET_USD_REALTRADES       = "ok_usd_realtrades"
	OKCOIN_WEBSOCKET_CNY_REALTRADES       = "ok_cny_realtrades"
	OKCOIN_WEBSOCKET_SPOTUSD_TRADE        = "ok_spotusd_trade"
	OKCOIN_WEBSOCKET_SPOTCNY_TRADE        = "ok_spotcny_trade"
	OKCOIN_WEBSOCKET_SPOTUSD_CANCEL_ORDER = "ok_spotusd_cancel_order"
	OKCOIN_WEBSOCKET_SPOTCNY_CANCEL_ORDER = "ok_spotcny_cancel_order"
	OKCOIN_WEBSOCKET_SPOTUSD_USERINFO     = "ok_spotusd_userinfo"
	OKCOIN_WEBSOCKET_SPOTCNY_USERINFO     = "ok_spotcny_userinfo"
	OKCOIN_WEBSOCKET_SPOTUSD_ORDER_INFO   = "ok_spotusd_order_info"
	OKCOIN_WEBSOCKET_SPOTCNY_ORDER_INFO   = "ok_spotcny_order_info"
	OKCOIN_WEBSOCKET_FUTURES_TRADE        = "ok_futuresusd_trade"
	OKCOIN_WEBSOCKET_FUTURES_CANCEL_ORDER = "ok_futuresusd_cancel_order"
	OKCOIN_WEBSOCKET_FUTURES_REALTRADES   = "ok_usd_future_realtrades"
	OKCOIN_WEBSOCKET_FUTURES_USERINFO     = "ok_futureusd_userinfo"
	OKCOIN_WEBSOCKET_FUTURES_ORDER_INFO   = "ok_futureusd_order_info"
)

type OKCoinWebsocketFutureIndex struct {
	FutureIndex float64 `json:"futureIndex"`
	Timestamp   int64   `json:"timestamp,string"`
}

type OKCoinWebsocketTicker struct {
	Timestamp float64
	Vol       string
	Buy       float64
	High      float64
	Last      float64
	Low       float64
	Sell      float64
}

type OKCoinWebsocketFuturesTicker struct {
	Buy        float64 `json:"buy"`
	ContractID string  `json:"contractId"`
	High       float64 `json:"high"`
	HoldAmount float64 `json:"hold_amount"`
	Last       float64 `json:"last,string"`
	Low        float64 `json:"low"`
	Sell       float64 `json:"sell"`
	UnitAmount float64 `json:"unitAmount"`
	Volume     float64 `json:"vol,string"`
}

type OKCoinWebsocketOrderbook struct {
	Asks      [][]float64 `json:"asks"`
	Bids      [][]float64 `json:"bids"`
	Timestamp int64       `json:"timestamp,string"`
}

type OKCoinWebsocketUserinfo struct {
	Info struct {
		Funds struct {
			Asset struct {
				Net   float64 `json:"net,string"`
				Total float64 `json:"total,string"`
			} `json:"asset"`
			Free struct {
				BTC float64 `json:"btc,string"`
				LTC float64 `json:"ltc,string"`
				USD float64 `json:"usd,string"`
				CNY float64 `json:"cny,string"`
			} `json:"free"`
			Frozen struct {
				BTC float64 `json:"btc,string"`
				LTC float64 `json:"ltc,string"`
				USD float64 `json:"usd,string"`
				CNY float64 `json:"cny,string"`
			} `json:"freezed"`
		} `json:"funds"`
	} `json:"info"`
	Result bool `json:"result"`
}

type OKCoinWebsocketFuturesContract struct {
	Available    float64 `json:"available"`
	Balance      float64 `json:"balance"`
	Bond         float64 `json:"bond"`
	ContractID   float64 `json:"contract_id"`
	ContractType string  `json:"contract_type"`
	Frozen       float64 `json:"freeze"`
	Profit       float64 `json:"profit"`
	Loss         float64 `json:"unprofit"`
}

type OKCoinWebsocketFuturesUserInfo struct {
	Info struct {
		BTC struct {
			Balance   float64                          `json:"balance"`
			Contracts []OKCoinWebsocketFuturesContract `json:"contracts"`
			Rights    float64                          `json:"rights"`
		} `json:"btc"`
		LTC struct {
			Balance   float64                          `json:"balance"`
			Contracts []OKCoinWebsocketFuturesContract `json:"contracts"`
			Rights    float64                          `json:"rights"`
		} `json:"ltc"`
	} `json:"info"`
	Result bool `json:"result"`
}

type OKCoinWebsocketOrder struct {
	Amount      float64 `json:"amount"`
	AvgPrice    float64 `json:"avg_price"`
	DateCreated float64 `json:"create_date"`
	TradeAmount float64 `json:"deal_amount"`
	OrderID     float64 `json:"order_id"`
	OrdersID    float64 `json:"orders_id"`
	Price       float64 `json:"price"`
	Status      int64   `json:"status"`
	Symbol      string  `json:"symbol"`
	OrderType   string  `json:"type"`
}

type OKCoinWebsocketFuturesOrder struct {
	Amount         float64 `json:"amount"`
	ContractName   string  `json:"contract_name"`
	DateCreated    float64 `json:"createdDate"`
	TradeAmount    float64 `json:"deal_amount"`
	Fee            float64 `json:"fee"`
	LeverageAmount int     `json:"lever_rate"`
	OrderID        float64 `json:"order_id"`
	Price          float64 `json:"price"`
	AvgPrice       float64 `json:"avg_price"`
	Status         int     `json:"status"`
	Symbol         string  `json:"symbol"`
	TradeType      int     `json:"type"`
	UnitAmount     float64 `json:"unit_amount"`
}

type OKCoinWebsocketRealtrades struct {
	AveragePrice         float64 `json:"averagePrice,string"`
	CompletedTradeAmount float64 `json:"completedTradeAmount,string"`
	DateCreated          float64 `json:"createdDate"`
	ID                   float64 `json:"id"`
	OrderID              float64 `json:"orderId"`
	SigTradeAmount       float64 `json:"sigTradeAmount,string"`
	SigTradePrice        float64 `json:"sigTradePrice,string"`
	Status               int64   `json:"status"`
	Symbol               string  `json:"symbol"`
	TradeAmount          float64 `json:"tradeAmount,string"`
	TradePrice           float64 `json:"buy,string"`
	TradeType            string  `json:"tradeType"`
	TradeUnitPrice       float64 `json:"tradeUnitPrice,string"`
	UnTrade              float64 `json:"unTrade,string"`
}

type OKCoinWebsocketFuturesRealtrades struct {
	Amount         float64 `json:"amount,string"`
	ContractID     float64 `json:"contract_id,string"`
	ContractName   string  `json:"contract_name"`
	ContractType   string  `json:"contract_type"`
	TradeAmount    float64 `json:"deal_amount,string"`
	Fee            float64 `json:"fee,string"`
	OrderID        float64 `json:"orderid"`
	Price          float64 `json:"price,string"`
	AvgPrice       float64 `json:"price_avg,string"`
	Status         int     `json:"status,string"`
	TradeType      int     `json:"type,string"`
	UnitAmount     float64 `json:"unit_amount,string"`
	LeverageAmount int     `json:"lever_rate,string"`
}

type OKCoinWebsocketEvent struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}

type OKCoinWebsocketResponse struct {
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type OKCoinWebsocketEventAuth struct {
	Event      string            `json:"event"`
	Channel    string            `json:"channel"`
	Parameters map[string]string `json:"parameters"`
}

type OKCoinWebsocketEventAuthRemove struct {
	Event      string            `json:"event"`
	Channel    string            `json:"channel"`
	Parameters map[string]string `json:"parameters"`
}

type OKCoinWebsocketTradeOrderResponse struct {
	OrderID int64 `json:"order_id,string"`
	Result  bool  `json:"result,string"`
}

func (o *OKCoin) PingHandler(message string) error {
	err := o.WebsocketConn.WriteControl(websocket.PingMessage, []byte("{'event':'ping'}"), time.Now().Add(time.Second))

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (o *OKCoin) AddChannel(channel string) {
	event := OKCoinWebsocketEvent{"addChannel", channel}
	json, err := JSONEncode(event)
	if err != nil {
		log.Println(err)
		return
	}
	err = o.WebsocketConn.WriteMessage(websocket.TextMessage, json)

	if err != nil {
		log.Println(err)
		return
	}

	if o.Verbose {
		log.Printf("%s Adding channel: %s\n", o.GetName(), channel)
	}
}

func (o *OKCoin) RemoveChannel(channel string) {
	event := OKCoinWebsocketEvent{"removeChannel", channel}
	json, err := JSONEncode(event)
	if err != nil {
		log.Println(err)
		return
	}
	err = o.WebsocketConn.WriteMessage(websocket.TextMessage, json)

	if err != nil {
		log.Println(err)
		return
	}

	if o.Verbose {
		log.Printf("%s Removing channel: %s\n", o.GetName(), channel)
	}
}

func (o *OKCoin) WebsocketSpotTrade(symbol, orderType string, price, amount float64) {
	values := make(map[string]string)
	values["symbol"] = symbol
	values["type"] = orderType
	values["price"] = strconv.FormatFloat(price, 'f', -1, 64)
	values["amount"] = strconv.FormatFloat(amount, 'f', -1, 64)
	channel := ""

	if o.WebsocketURL == OKCOIN_WEBSOCKET_URL_CHINA {
		channel = OKCOIN_WEBSOCKET_SPOTCNY_TRADE
	} else {
		channel = OKCOIN_WEBSOCKET_SPOTUSD_TRADE
	}

	o.AddChannelAuthenticated(channel, values)
}

func (o *OKCoin) WebsocketFuturesTrade(symbol, contractType string, price, amount float64, orderType, matchPrice, leverage int) {
	values := make(map[string]string)
	values["symbol"] = symbol
	values["contract_type"] = contractType
	values["price"] = strconv.FormatFloat(price, 'f', -1, 64)
	values["amount"] = strconv.FormatFloat(amount, 'f', -1, 64)
	values["type"] = strconv.Itoa(orderType)
	values["match_price"] = strconv.Itoa(matchPrice)
	values["lever_rate"] = strconv.Itoa(orderType)
	o.AddChannelAuthenticated(OKCOIN_WEBSOCKET_FUTURES_TRADE, values)
}

func (o *OKCoin) WebsocketSpotCancel(symbol string, orderID int64) {
	values := make(map[string]string)
	values["symbol"] = symbol
	values["order_id"] = strconv.FormatInt(orderID, 10)
	channel := ""

	if o.WebsocketURL == OKCOIN_WEBSOCKET_URL_CHINA {
		channel = OKCOIN_WEBSOCKET_SPOTCNY_CANCEL_ORDER
	} else {
		channel = OKCOIN_WEBSOCKET_SPOTUSD_CANCEL_ORDER
	}

	o.AddChannelAuthenticated(channel, values)
}

func (o *OKCoin) WebsocketFuturesCancel(symbol, contractType string, orderID int64) {
	values := make(map[string]string)
	values["symbol"] = symbol
	values["order_id"] = strconv.FormatInt(orderID, 10)
	values["contract_type"] = contractType
	o.AddChannelAuthenticated(OKCOIN_WEBSOCKET_FUTURES_CANCEL_ORDER, values)
}

func (o *OKCoin) WebsocketSpotOrderInfo(symbol string, orderID int64) {
	values := make(map[string]string)
	values["symbol"] = symbol
	values["order_id"] = strconv.FormatInt(orderID, 10)
	channel := ""

	if o.WebsocketURL == OKCOIN_WEBSOCKET_URL_CHINA {
		channel = OKCOIN_WEBSOCKET_SPOTCNY_ORDER_INFO
	} else {
		channel = OKCOIN_WEBSOCKET_SPOTUSD_ORDER_INFO
	}

	o.AddChannelAuthenticated(channel, values)
}

func (o *OKCoin) WebsocketFuturesOrderInfo(symbol, contractType string, orderID int64, orderStatus, currentPage, pageLength int) {
	values := make(map[string]string)
	values["symbol"] = symbol
	values["order_id"] = strconv.FormatInt(orderID, 10)
	values["contract_type"] = contractType
	values["status"] = strconv.Itoa(orderStatus)
	values["current_page"] = strconv.Itoa(currentPage)
	values["page_length"] = strconv.Itoa(pageLength)
	o.AddChannelAuthenticated(OKCOIN_WEBSOCKET_FUTURES_ORDER_INFO, values)
}

func (o *OKCoin) ConvertToURLValues(values map[string]string) url.Values {
	urlVals := url.Values{}
	for i, x := range values {
		urlVals.Set(i, x)
	}
	return urlVals
}

func (o *OKCoin) WebsocketSign(values map[string]string) string {
	values["api_key"] = o.PartnerID
	urlVals := o.ConvertToURLValues(values)
	return strings.ToUpper(HexEncodeToString(GetMD5([]byte(urlVals.Encode() + "&secret_key=" + o.SecretKey))))
}

func (o *OKCoin) AddChannelAuthenticated(channel string, values map[string]string) {
	values["sign"] = o.WebsocketSign(values)
	event := OKCoinWebsocketEventAuth{"addChannel", channel, values}
	json, err := JSONEncode(event)
	if err != nil {
		log.Println(err)
		return
	}
	err = o.WebsocketConn.WriteMessage(websocket.TextMessage, json)

	if err != nil {
		log.Println(err)
		return
	}

	if o.Verbose {
		log.Printf("%s Adding authenticated channel: %s\n", o.GetName(), channel)
	}
}

func (o *OKCoin) RemoveChannelAuthenticated(conn *websocket.Conn, channel string, values map[string]string) {
	values["sign"] = o.WebsocketSign(values)
	event := OKCoinWebsocketEventAuthRemove{"removeChannel", channel, values}
	json, err := JSONEncode(event)
	if err != nil {
		log.Println(err)
		return
	}
	err = o.WebsocketConn.WriteMessage(websocket.TextMessage, json)

	if err != nil {
		log.Println(err)
		return
	}

	if o.Verbose {
		log.Printf("%s Removing authenticated channel: %s\n", o.GetName(), channel)
	}
}

func (o *OKCoin) WebsocketClient() {
	klineValues := []string{"1min", "3min", "5min", "15min", "30min", "1hour", "2hour", "4hour", "6hour", "12hour", "day", "3day", "week"}
	currencyChan, userinfoChan := "", ""

	if o.WebsocketURL == OKCOIN_WEBSOCKET_URL_CHINA {
		currencyChan = OKCOIN_WEBSOCKET_CNY_REALTRADES
		userinfoChan = OKCOIN_WEBSOCKET_SPOTCNY_USERINFO
	} else {
		currencyChan = OKCOIN_WEBSOCKET_USD_REALTRADES
		userinfoChan = OKCOIN_WEBSOCKET_SPOTUSD_USERINFO
	}

	for o.Enabled && o.Websocket {
		var Dialer websocket.Dialer
		var err error
		o.WebsocketConn, _, err = Dialer.Dial(o.WebsocketURL, http.Header{})

		if err != nil {
			log.Printf("%s Unable to connect to Websocket. Error: %s\n", o.GetName(), err)
			continue
		}

		if o.Verbose {
			log.Printf("%s Connected to Websocket.\n", o.GetName())
		}

		o.WebsocketConn.SetPingHandler(o.PingHandler)

		if o.AuthenticatedAPISupport {
			if o.WebsocketURL == OKCOIN_WEBSOCKET_URL {
				o.AddChannelAuthenticated(OKCOIN_WEBSOCKET_FUTURES_REALTRADES, map[string]string{})
				o.AddChannelAuthenticated(OKCOIN_WEBSOCKET_FUTURES_USERINFO, map[string]string{})
			}
			o.AddChannelAuthenticated(currencyChan, map[string]string{})
			o.AddChannelAuthenticated(userinfoChan, map[string]string{})
		}

		for _, x := range o.EnabledPairs {
			currency := StringToLower(x)
			currencyUL := currency[0:3] + "_" + currency[3:]
			if o.AuthenticatedAPISupport {
				o.WebsocketSpotOrderInfo(currencyUL, -1)
			}
			if o.WebsocketURL == OKCOIN_WEBSOCKET_URL {
				o.AddChannel(fmt.Sprintf("ok_%s_future_index", currency))
				for _, y := range o.FuturesValues {
					if o.AuthenticatedAPISupport {
						o.WebsocketFuturesOrderInfo(currencyUL, y, -1, 1, 1, 50)
					}
					o.AddChannel(fmt.Sprintf("ok_%s_future_ticker_%s", currency, y))
					o.AddChannel(fmt.Sprintf("ok_%s_future_depth_%s_60", currency, y))
					o.AddChannel(fmt.Sprintf("ok_%s_future_trade_v1_%s", currency, y))
					for _, z := range klineValues {
						o.AddChannel(fmt.Sprintf("ok_future_%s_kline_%s_%s", currency, y, z))
					}
				}
			} else {
				o.AddChannel(fmt.Sprintf("ok_%s_ticker", currency))
				o.AddChannel(fmt.Sprintf("ok_%s_depth60", currency))
				o.AddChannel(fmt.Sprintf("ok_%s_trades_v1", currency))

				for _, y := range klineValues {
					o.AddChannel(fmt.Sprintf("ok_%s_kline_%s", currency, y))
				}
			}
		}

		for o.Enabled && o.Websocket {
			msgType, resp, err := o.WebsocketConn.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			switch msgType {
			case websocket.TextMessage:
				response := []interface{}{}
				err = JSONDecode(resp, &response)

				if err != nil {
					log.Println(err)
					continue
				}

				for _, y := range response {
					z := y.(map[string]interface{})
					channel := z["channel"]
					data := z["data"]
					success := z["success"]
					errorcode := z["errorcode"]
					channelStr, ok := channel.(string)

					if !ok {
						log.Println("Unable to convert channel to string")
						continue
					}

					if success != "true" && success != nil {
						errorCodeStr, ok := errorcode.(string)
						if !ok {
							log.Printf("%s Websocket: Unable to convert errorcode to string.\n", o.GetName)
							log.Printf("%s Websocket: channel %s error code: %s.\n", o.GetName(), channelStr, errorcode)
						} else {
							log.Printf("%s Websocket: channel %s error: %s.\n", o.GetName(), channelStr, o.WebsocketErrors[errorCodeStr])
						}
						continue
					}

					dataJSON, err := JSONEncode(data)

					if err != nil {
						log.Println(err)
						continue
					}

					switch true {
					case StringContains(channelStr, "ticker") && !StringContains(channelStr, "future"):
						tickerValues := []string{"buy", "high", "last", "low", "sell", "timestamp"}
						tickerMap := data.(map[string]interface{})
						ticker := OKCoinWebsocketTicker{}
						ticker.Vol = tickerMap["vol"].(string)

						for _, z := range tickerValues {
							result := reflect.TypeOf(tickerMap[z]).String()
							if result == "string" {
								value, err := strconv.ParseFloat(tickerMap[z].(string), 64)
								if err != nil {
									log.Println(err)
									continue
								}

								switch z {
								case "buy":
									ticker.Buy = value
								case "high":
									ticker.High = value
								case "last":
									ticker.Last = value
								case "low":
									ticker.Low = value
								case "sell":
									ticker.Sell = value
								case "timestamp":
									ticker.Timestamp = value
								}

							} else if result == "float64" {
								switch z {
								case "buy":
									ticker.Buy = tickerMap[z].(float64)
								case "high":
									ticker.High = tickerMap[z].(float64)
								case "last":
									ticker.Last = tickerMap[z].(float64)
								case "low":
									ticker.Low = tickerMap[z].(float64)
								case "sell":
									ticker.Sell = tickerMap[z].(float64)
								case "timestamp":
									ticker.Timestamp = tickerMap[z].(float64)
								}
							}
						}
					case StringContains(channelStr, "ticker") && StringContains(channelStr, "future"):
						ticker := OKCoinWebsocketFuturesTicker{}
						err = JSONDecode(dataJSON, &ticker)

						if err != nil {
							log.Println(err)
							continue
						}
					case StringContains(channelStr, "depth"):
						orderbook := OKCoinWebsocketOrderbook{}
						err = JSONDecode(dataJSON, &orderbook)

						if err != nil {
							log.Println(err)
							continue
						}
					case StringContains(channelStr, "trades_v1") || StringContains(channelStr, "trade_v1"):
						type TradeResponse struct {
							Data [][]string
						}

						trades := TradeResponse{}
						err = JSONDecode(dataJSON, &trades.Data)

						if err != nil {
							log.Println(err)
							continue
						}
						// to-do: convert from string array to trade struct
					case StringContains(channelStr, "kline"):
						klines := []interface{}{}
						err := JSONDecode(dataJSON, &klines)

						if err != nil {
							log.Println(err)
							continue
						}
					case StringContains(channelStr, "spot") && StringContains(channelStr, "realtrades"):
						if string(dataJSON) == "null" {
							continue
						}
						realtrades := OKCoinWebsocketRealtrades{}
						err := JSONDecode(dataJSON, &realtrades)

						if err != nil {
							log.Println(err)
							continue
						}
					case StringContains(channelStr, "future") && StringContains(channelStr, "realtrades"):
						if string(dataJSON) == "null" {
							continue
						}
						realtrades := OKCoinWebsocketFuturesRealtrades{}
						err := JSONDecode(dataJSON, &realtrades)

						if err != nil {
							log.Println(err)
							continue
						}
					case StringContains(channelStr, "spot") && StringContains(channelStr, "trade") || StringContains(channelStr, "futures") && StringContains(channelStr, "trade"):
						tradeOrder := OKCoinWebsocketTradeOrderResponse{}
						err := JSONDecode(dataJSON, &tradeOrder)

						if err != nil {
							log.Println(err)
							continue
						}
					case StringContains(channelStr, "cancel_order"):
						cancelOrder := OKCoinWebsocketTradeOrderResponse{}
						err := JSONDecode(dataJSON, &cancelOrder)

						if err != nil {
							log.Println(err)
							continue
						}
					case StringContains(channelStr, "spot") && StringContains(channelStr, "userinfo"):
						userinfo := OKCoinWebsocketUserinfo{}
						err = JSONDecode(dataJSON, &userinfo)

						if err != nil {
							log.Println(err)
							continue
						}
					case StringContains(channelStr, "futureusd_userinfo"):
						userinfo := OKCoinWebsocketFuturesUserInfo{}
						err = JSONDecode(dataJSON, &userinfo)

						if err != nil {
							log.Println(err)
							continue
						}
					case StringContains(channelStr, "spot") && StringContains(channelStr, "order_info"):
						type OrderInfoResponse struct {
							Result bool                   `json:"result"`
							Orders []OKCoinWebsocketOrder `json:"orders"`
						}
						var orders OrderInfoResponse
						err := JSONDecode(dataJSON, &orders)

						if err != nil {
							log.Println(err)
							continue
						}
					case StringContains(channelStr, "futureusd_order_info"):
						type OrderInfoResponse struct {
							Result bool                          `json:"result"`
							Orders []OKCoinWebsocketFuturesOrder `json:"orders"`
						}
						var orders OrderInfoResponse
						err := JSONDecode(dataJSON, &orders)

						if err != nil {
							log.Println(err)
							continue
						}
					case StringContains(channelStr, "future_index"):
						index := OKCoinWebsocketFutureIndex{}
						err = JSONDecode(dataJSON, &index)

						if err != nil {
							log.Println(err)
							continue
						}
					}
				}
			}
		}
		o.WebsocketConn.Close()
		log.Printf("%s Websocket client disconnected.", o.GetName())
	}
}

func (o *OKCoin) SetWebsocketErrorDefaults() {
	o.WebsocketErrors = map[string]string{
		"10001": "Illegal parameters",
		"10002": "Authentication failure",
		"10003": "This connection has requested other user data",
		"10004": "This connection did not request this user data",
		"10005": "System error",
		"10009": "Order does not exist",
		"10010": "Insufficient funds",
		"10011": "Order quantity too low",
		"10012": "Only support btc_usd/btc_cny ltc_usd/ltc_cny",
		"10014": "Order price must be between 0 - 1,000,000",
		"10015": "Channel subscription temporally not available",
		"10016": "Insufficient coins",
		"10017": "WebSocket authorization error",
		"10100": "User frozen",
		"10216": "Non-public API",
		"20001": "User does not exist",
		"20002": "User frozen",
		"20003": "Frozen due to force liquidation",
		"20004": "Future account frozen",
		"20005": "User future account does not exist",
		"20006": "Required field can not be null",
		"20007": "Illegal parameter",
		"20008": "Future account fund balance is zero",
		"20009": "Future contract status error",
		"20010": "Risk rate information does not exist",
		"20011": `Risk rate bigger than 90% before opening position`,
		"20012": `Risk rate bigger than 90% after opening position`,
		"20013": "Temporally no counter party price",
		"20014": "System error",
		"20015": "Order does not exist",
		"20016": "Liquidation quantity bigger than holding",
		"20017": "Not authorized/illegal order ID",
		"20018": `Order price higher than 105% or lower than 95% of the price of last minute`,
		"20019": "IP restrained to access the resource",
		"20020": "Secret key does not exist",
		"20021": "Index information does not exist",
		"20022": "Wrong API interface",
		"20023": "Fixed margin user",
		"20024": "Signature does not match",
		"20025": "Leverage rate error",
	}
}
