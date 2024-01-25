package main

// Starting corpus of paths to test. You may want to extend this with a list of
// paths that you want to ensure are tested.
var initPaths = []string{
	"/transactions?order=desc",
	"/transactions?order=desc",
	"/transactions?order=desc&include_failed=false",
	"/transactions?order=desc&include_failed=true",

	"/operations?order=desc",
	"/operations?order=desc&include_failed=false",
	"/operations?order=desc&include_failed=true",

	"/operations?join=transactions&order=desc",
	"/operations?join=transactions&order=desc&include_failed=false",
	"/operations?join=transactions&order=desc&include_failed=true",

	"/payments?order=desc",
	"/payments?order=desc&include_failed=false",
	"/payments?order=desc&include_failed=true",

	"/payments?join=transactions&order=desc",
	"/payments?join=transactions&order=desc&include_failed=false",
	"/payments?join=transactions&order=desc&include_failed=true",

	"/ledgers?order=desc",
	"/effects?order=desc",
	"/trades?order=desc",

	"/accounts/GAKLCFRTFDXKOEEUSBS23FBSUUVJRMDQHGCHNGGGJZQRK7BCPIMHUC4P/transactions?limit=200",
	"/accounts/GAKLCFRTFDXKOEEUSBS23FBSUUVJRMDQHGCHNGGGJZQRK7BCPIMHUC4P/transactions?limit=200&include_failed=false",
	"/accounts/GAKLCFRTFDXKOEEUSBS23FBSUUVJRMDQHGCHNGGGJZQRK7BCPIMHUC4P/transactions?limit=200&include_failed=true",

	"/accounts/GAKLCFRTFDXKOEEUSBS23FBSUUVJRMDQHGCHNGGGJZQRK7BCPIMHUC4P/operations?limit=200",
	"/accounts/GAKLCFRTFDXKOEEUSBS23FBSUUVJRMDQHGCHNGGGJZQRK7BCPIMHUC4P/payments?limit=200",
	"/accounts/GAKLCFRTFDXKOEEUSBS23FBSUUVJRMDQHGCHNGGGJZQRK7BCPIMHUC4P/effects?limit=200",

	"/accounts/GC2ZV6KGGFLQIMDVDWBWCP6LTODUDXYBLUPTUZCFHIMDCWHR43ULZITJ/trades?limit=200",
	"/accounts/GC2ZV6KGGFLQIMDVDWBWCP6LTODUDXYBLUPTUZCFHIMDCWHR43ULZITJ/offers?limit=200",

	// Pubnet markets
	"/order_book?selling_asset_type=native&buying_asset_type=credit_alphanum4&buying_asset_code=LTC&buying_asset_issuer=GCSTRLTC73UVXIYPHYTTQUUSDTQU2KQW5VKCE4YCMEHWF44JKDMQAL23",
	"/order_book?selling_asset_type=native&buying_asset_type=credit_alphanum4&buying_asset_code=XRP&buying_asset_issuer=GCSTRLTC73UVXIYPHYTTQUUSDTQU2KQW5VKCE4YCMEHWF44JKDMQAL23",
	"/order_book?selling_asset_type=native&buying_asset_type=credit_alphanum4&buying_asset_code=BTC&buying_asset_issuer=GCSTRLTC73UVXIYPHYTTQUUSDTQU2KQW5VKCE4YCMEHWF44JKDMQAL23",
	"/order_book?selling_asset_type=native&buying_asset_type=credit_alphanum4&buying_asset_code=USD&buying_asset_issuer=GBSTRUSD7IRX73RQZBL3RQUH6KS3O4NYFY3QCALDLZD77XMZOPWAVTUK",
	"/order_book?selling_asset_type=native&buying_asset_type=credit_alphanum4&buying_asset_code=SLT&buying_asset_issuer=GCKA6K5PCQ6PNF5RQBF7PQDJWRHO6UOGFMRLK3DYHDOI244V47XKQ4GP",

	"/trade_aggregations?base_asset_type=native&counter_asset_code=USD&counter_asset_issuer=GBSTRUSD7IRX73RQZBL3RQUH6KS3O4NYFY3QCALDLZD77XMZOPWAVTUK&counter_asset_type=credit_alphanum4&end_time=1551866400000&limit=200&order=desc&resolution=900000&start_time=1514764800",
}
