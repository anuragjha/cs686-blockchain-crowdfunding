package p3

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Start",
		"GET",
		"/start",
		Start,
	},
	Route{
		"Show",
		"GET",
		"/show",
		Show,
	},
	Route{
		"Upload",
		"GET",
		"/upload",
		Upload,
	},
	Route{
		"Upload",
		"GET",
		"/uploadpids",
		UploadPids,
	},
	Route{
		"UploadBlock",
		"GET",
		"/block/{height}/{hash}",
		UploadBlock,
	},
	Route{
		"HeartBeatReceive",
		"POST",
		"/heartbeat/receive",
		HeartBeatReceive,
	},
	Route{
		"Canonical",
		"GET",
		"/canonical",
		Canonical,
	},
	////////currency
	Route{
		"Transaction",
		"POST",
		"/transaction", //to put in tx pool
		Transaction,
	},
	Route{
		"ShowWallet",
		"GET",
		"/showWallet", //to put in tx pool
		ShowWallet,
	},
	Route{
		"ShowBalanceBook",
		"GET",
		"/showBalanceBook", //to put in tx pool
		ShowBalanceBook,
	},
	Route{
		"ShowTransactionPool",
		"GET",
		"/showTransactionPool", //to put in tx pool
		ShowTransactionPool,
	},
}
