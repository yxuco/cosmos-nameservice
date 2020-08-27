# cosmos-nameservice

Detailed steps of implementation are described in [Cosmos SDK Tutorial](https://tutorials.cosmos.network/nameservice/tutorial/00-intro.html#). You may go directly to `Build app` below to build and start the app.

## Download and build `scaffold`

```bash
mkdir -p ~/work/cosmos-tutorial
cd ~/work/cosmos-tutorial
git clone git@github.com:cosmos/scaffold.git
cd scaffold
make
```

## Generate app and module

```bash
cd ~/work/cosmos-tutorial
scaffold app lvl-1 yxuco cosmos-nameservice
cd cosmos-nameservice/x
scaffold module yxuco cosmos-nameservice nameservice
cd .. && go mod tidy
```

## Define messages

- Define messages using the generated template in `x/nameservice/types/msg.go`
- implement sdk.Msg interface: <https://github.com/cosmos/cosmos-sdk/blob/master/types/tx_msg.go>
- Register messages in `x/nameservice/types/codec.go`
- Create `x/nameservice/alias.go` to make nested definitions of types and keeper available at module's top level

## Define key store and keeper functions

- In `x/nameservice/types/types.go`, define data struct for key store
- In `x/nameservice/keeper/keeper.go`, define functions for accessing key store

## Implement message handler

Implement handler functions to processe each message and interact with key store.

- Edit generated templates in `x/nameservice/handler.go`
- Define error types in `x/nameservice/types/errors.go`
- Define keeper interface for coin transfer in `x/namespace/types/expected_keepers.go`

## Implement querier to query key store

- Implement query functions in `x/nameservice/keeper/querier.go`
- Define query types in `x/nameservice/types/querier.go`

## Implement CLI

Implement CLI app to send transaction and query.

- Implement functions for creating transactions in `x/nameservice/client/cli/tx.go`
- Implement functions for creating queries in `x/nameservice/client/cli/query.go`

## Implement REST interface

- Implement REST interface in `x/nameservice/client/rest`

## Make module reusable by any app

- Edit `x/nameservice/module.go` to add bank.Keeper to AppModule struct and NowAppModule function

## Use module in app

- In `app/app.go`, follow comments of `TODO` to change name of appd appcli and add module refs
- Change folder name of `cmd/appd` and `cmd/appcli`
- Edit rootCmd description in `cmd/appcli/main.go`
- Edit rootCmd description in `cmd/appd/main.go`

## Build app

- At app root, `go mod tidy`
- Edit Makefile to set `appd`and `appcli`
- make
- Verify executable: `nsd -h`, and `nscli -h`

## Config and start a node

```bash
# create genesis.json under ~/.nsd/config
nsd init mynode --chain-id namechain

# set some options in ~/.nscli/config/config.toml
nscli config chain-id namechain
nscli config output json
nscli config indent true
nscli config trust-node true

# for test only, use test keyring
nscli config keyring-backend test

# create user jack & alice
nscli keys add jack
nscli keys add alice

# add accounts to genesis file with initial tokens for name service and stake
nsd add-genesis-account $(nscli keys show jack -a) 1000nametoken,100000000stake
nsd add-genesis-account $(nscli keys show alice -a) 1000nametoken,100000000stake

# genesis tx to create a validator signed by jack (output in ~/.nsd/config/gentx)
nsd gentx --name jack --keyring-backend test

# add genesis tx to genesis.json
nsd collect-gentxs

# validate genesis.json
nsd validate-genesis

# start node
nsd start
```

## Test nameservice

```bash
# First check the accounts to ensure they have funds
nscli query account $(nscli keys show jack -a)
nscli query account $(nscli keys show alice -a)

# Buy your first name using your coins from the genesis file
nscli tx nameservice buy-name jack.id 5nametoken --from jack

# Set the value for the name you just bought
nscli tx nameservice set-name jack.id 8.8.8.8 --from jack

# Try out a resolve query against the name you registered
nscli query nameservice resolve jack.id
# > 8.8.8.8

# Try out a whois query against the name you just registered
nscli query nameservice whois jack.id
# > {"value":"8.8.8.8","owner":"cosmos1l7k5tdt2qam0zecxrx78yuw447ga54dsmtpk2s","price":[{"denom":"nametoken","amount":"5"}]}

# Alice buys name from jack
nscli tx nameservice buy-name jack.id 10nametoken --from alice

# Alice decides to delete the name she just bought from jack
nscli tx nameservice delete-name jack.id --from alice

# Try out a whois query against the name you just deleted
nscli query nameservice whois jack.id
# > {"value":"","owner":"","price":[{"denom":"nametoken","amount":"1"}]}
```

## Start rest server

```bash
nscli rest-server --chain-id namechain --trust-node
```

## Test REST service

```bash
# Get the sequence and account numbers for jack to construct the below requests
curl -s http://localhost:1317/auth/accounts/$(nscli keys show jack -a)
# > {"type":"auth/Account","value":{"address":"cosmos127qa40nmq56hu27ae263zvfk3ey0tkapwk0gq6","coins":[{"denom":"jackCoin","amount":"1000"},{"denom":"nametoken","amount":"1010"}],"public_key":{"type":"tendermint/PubKeySecp256k1","value":"A9YxyEbSWzLr+IdK/PuMUYmYToKYQ3P/pM8SI1Bxx3wu"},"account_number":"0","sequence":"1"}}

# Get the sequence and account numbers for alice to construct the below requests
curl -s http://localhost:1317/auth/accounts/$(nscli keys show alice -a)
# > {"type":"auth/Account","value":{"address":"cosmos1h7ztnf2zkf4558hdxv5kpemdrg3tf94hnpvgsl","coins":[{"denom":"aliceCoin","amount":"1000"},{"denom":"nametoken","amount":"980"}],"public_key":{"type":"tendermint/PubKeySecp256k1","value":"Avc7qwecLHz5qb1EKDuSTLJfVOjBQezk0KSPDNybLONJ"},"account_number":"1","sequence":"2"}}

# Buy another name for jack, first create the raw transaction
# NOTE: Be sure to specialize this request for your specific environment, also the "buyer" and "from" should be the same address
curl -XPOST -s http://localhost:1317/nameservice/names --data-binary '{"base_req":{"from":"'$(nscli keys show jack -a)'","chain_id":"namechain"},"name":"jack1.id","amount":"5nametoken","buyer":"'$(nscli keys show jack -a)'"}' > unsignedTx.json

# Then sign this transaction
# NOTE: In a real environment the raw transaction should be signed on the client side. Also the sequence needs to be adjusted, depending on what the query of alice's account has shown.
nscli tx sign unsignedTx.json --from jack --offline --chain-id namechain --sequence 1 --account-number 0 > signedTx.json

# And finally broadcast the signed transaction
nscli tx broadcast signedTx.json
# > { "height": "266", "txhash": "C041AF0CE32FBAE5A4DD6545E4B1F2CB786879F75E2D62C79D690DAE163470BC", "logs": [  {   "msg_index": "0",   "success": true,   "log": ""  } ],"gas_wanted":"200000", "gas_used": "41510", "tags": [  {   "key": "action",   "value": "buy_name"  } ]}

# Set the data for that name that jack just bought
# NOTE: Be sure to specialize this request for your specific environment, also the "owner" and "from" should be the same address
curl -XPUT -s http://localhost:1317/nameservice/names --data-binary '{"base_req":{"from":"'$(nscli keys show jack -a)'","chain_id":"namechain"},"name":"jack1.id","value":"8.8.4.4","owner":"'$(nscli keys show jack -a)'"}' > unsignedTx.json
# > {"check_tx":{"gasWanted":"200000","gasUsed":"1242"},"deliver_tx":{"log":"Msg 0: ","gasWanted":"200000","gasUsed":"1352","tags":[{"key":"YWN0aW9u","value":"c2V0X25hbWU="}]},"hash":"B4DF0105D57380D60524664A2E818428321A0DCA1B6B2F091FB3BEC54D68FAD7","height":"26"}

# Again we need to sign and broadcast
nscli tx sign unsignedTx.json --from jack --offline --chain-id namechain --sequence 2 --account-number 0 > signedTx.json
nscli tx broadcast signedTx.json

# Query the value for the name jack just set
curl -s http://localhost:1317/nameservice/names/jack1.id
# 8.8.4.4

# Query whois for the name jack just bought
curl -s http://localhost:1317/nameservice/names/jack1.id/whois
# > {"value":"8.8.8.8","owner":"cosmos127qa40nmq56hu27ae263zvfk3ey0tkapwk0gq6","price":[{"denom":"STAKE","amount":"10"}]}

# Alice buys name from jack
curl -XPOST -s http://localhost:1317/nameservice/names --data-binary '{"base_req":{"from":"'$(nscli keys show alice -a)'","chain_id":"namechain"},"name":"jack1.id","amount":"10nametoken","buyer":"'$(nscli keys show alice -a)'"}' > unsignedTx.json

# Again we need to sign and broadcast
# NOTE: The account number has changed to 1 and the sequence is now 2, according to the query of alice's account
nscli tx sign unsignedTx.json --from alice --offline --chain-id namechain --sequence 2 --account-number 1 > signedTx.json
nscli tx broadcast signedTx.json
# > { "height": "1515", "txhash": "C9DCC423E10E7E5E40A549057A4AA060DA6D6A885A394F6ED5C0E40AEE984A77", "logs": [  {   "msg_index": "0",   "success": true,   "log": ""  } ],"gas_wanted": "200000", "gas_used": "42375", "tags": [  {   "key": "action",   "value": "buy_name"  } ]}

# Now, Alice no longer needs the name she bought from jack and hence deletes it
# NOTE: Only the owner can delete the name. Since she is one, she can delete the name she bought from jack
curl -XDELETE -s http://localhost:1317/nameservice/names --data-binary '{"base_req":{"from":"'$(nscli keys show alice -a)'","chain_id":"namechain"},"name":"jack1.id","owner":"'$(nscli keys show alice -a)'"}' > unsignedTx.json

# And a final time sign and broadcast
# NOTE: The account number is still 1, but the sequence is changed to 3, according to the query of alice's account
nscli tx sign unsignedTx.json --from alice --offline --chain-id namechain --sequence 3 --account-number 1 > signedTx.json
nscli tx broadcast signedTx.json

# Query whois for the name Alice just deleted
curl -s http://localhost:1317/nameservice/names/jack1.id/whois
# > {"value":"","owner":"","price":[{"denom":"STAKE","amount":"1"}]}
```
