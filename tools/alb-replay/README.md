# ALB log replayer

Tool that replays the successful GET requests found in an [AWS Application Load Balancer log file](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-access-logs.html).

## Install

Compile the `alb-replay` binary:

```bash
go install ./tools/alb-replay
```

## Usage

```
Usage of ./alb-replay:
  alb-replay <aws_log_file> <target_host_base_url>
  -start-from int
    	What URL number to start from (default 1)
  -timeout duration
    	HTTP request timeout (default 5s)
  -workers int
    	How many parallel workers to use (default 1)
```

## Example

```
alb-replay --workers 100 746476062914_elasticloadbalancing_us-east-1_app.diamcircle002-prd-aurora2.d65c0ca4271aa022_20210628T0000Z_54.208.185.115_567ts68u.log  https://aurora.diamcircle.org
2021/08/04 16:36:13 (4) 506.607706ms https://aurora.diamcircle.org/ledgers?limit=1&order=desc&c=0.344613637344948
2021/08/04 16:36:13 (1) 517.0601ms https://aurora.diamcircle.org/
2021/08/04 16:36:13 (5) 518.765858ms https://aurora.diamcircle.org/trades?cursor=155070734820388867-0&X-Client-Name=js-diamcircle-sdk&X-Client-Version=7.0.0
2021/08/04 16:36:13 (10) 518.778775ms https://aurora.diamcircle.org/offers?seller=GDYYFHJ34WSXDSNTPGQMS3NIS6PJR5WPZKSVZPADAR43RKDHZRIU7A5V
2021/08/04 16:36:13 (3) 519.966962ms https://aurora.diamcircle.org/order_book?selling_asset_type=credit_alphanum12&selling_asset_code=DOGET&selling_asset_issuer=GDOEVDDBU6OBWKL7VHDAOKD77UP4DKHQYKOKJJT5PR3WRDBTX35HUEUX&buying_asset_type=native&X-Client-Name=js-diamcircle-sdk&X-Client-Version=8.0.0
2021/08/04 16:36:13 (2) 520.147353ms https://mainnet.diamcircle.io/GDNXSZSF7HIGVRL2LG6VWXN5PWV3KHI77DQTHLLPKNPLUZFKRRDQJBXP?c=0.020849836853811032
2021/08/04 16:36:13 (6) 520.645528ms https://aurora.diamcircle.org/trades?base_asset_type=native&counter_asset_type=credit_alphanum12&counter_asset_code=DOGET&counter_asset_issuer=GDOEVDDBU6OBWKL7VHDAOKD77UP4DKHQYKOKJJT5PR3WRDBTX35HUEUX&limit=50&order=desc&c=0.433731649711667
2021/08/04 16:36:13 (8) 523.850289ms https://aurora.diamcircle.org/order_book?selling_asset_type=native&buying_asset_type=credit_alphanum12&buying_asset_code=Falcon9&buying_asset_issuer=GCHG35QMNQ6MOZEIQNHKGABUWOEVF7STLOBYEPQXARI7QAIV6ZVVPNKQ&limit=200&c=0.7329191653797233
[...]
```
