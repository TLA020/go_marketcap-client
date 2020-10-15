# go_marketcap-client
Simple GO client to fetch prices of a cryptocurrency from coinmarketcap.

In the "basic" coinmarketcap license you're not allowed to use other endpoints. And also it's not possible to fetch by symbol. Only by coinmarketcapID.
Thats why we use the IdCurrency mapping to find the right ID by the symbol.
