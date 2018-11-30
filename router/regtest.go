package router

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/btcsuite/btcd/btcjson"

	"github.com/btcsuite/btcd/wire"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/gorilla/mux"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
)

const (
	host     = "localhost:19001"
	user     = "admin1"
	password = "123"
)

// RegTest handles communication with both the regtest daemon and the local database
// @param Client <*Client>: handles remote API calls to the regtest daemon
// @param DB <*Database>: BBolt database is used to keep track of the utxos.
// 												This handles read/write/list/delete operations to the db
type RegTest struct {
	Client *rpcclient.Client
	DB     *Database
}

// New configures and creates a Client instance
func (r *RegTest) New() error {
	// Configure regtest client to connect the daemon
	connConfig := &rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         host,
		User:         user,
		Pass:         password,
	}

	client, err := rpcclient.New(connConfig, nil)
	if err != nil {
		return err
	}

	// Setup db
	db := &Database{}
	err = db.New()
	if err != nil {
		return err
	}

	r.Client = client
	r.DB = db

	// Since regtest chain should be empty at this point, lets mine some blocks to increase miner balance
	mine(r, 200)

	return nil
}

// Shutdown disconnect from rpc server and stop all goroutines
func (r *RegTest) Shutdown() {
	r.Client.Shutdown()
	r.DB.Close()
}

type txResponse struct {
	TxHash string `json:"tx_hash"`
}

// SendTo sends 1 btc to the given address from the miner account (faucet service)
// Here the db is updated adding the new utxo to the "unpent" bucket
func (r *RegTest) SendTo(w http.ResponseWriter, req *http.Request) {
	// send request through regtest client
	address := mux.Vars(req)["address"]
	txHash, blockHash, err := sendTo(r, address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// from the transaction we retrieve the hash of the utxo by double_hash(txHash + outputTxIndex)
	// this hash is used as the id of the utxo in the db, the value is the stringified json representation of the utxo
	key, value, err := prepareUtxo(r, txHash, address, blockHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = r.DB.Update(address, "unspent", txHash.String(), key, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := txResponse{txHash.String()}
	json.NewEncoder(w).Encode(resp)
}

// Broadcast publishes a raw transaction to the network
// Here we need to do things both for tx inputs and outputs:
// - for each input, move the corresponding entry from the "unpent" bucket to the "spent" one
// - for each output, create an entry in the "unspent" bucket
func (r *RegTest) Broadcast(w http.ResponseWriter, req *http.Request) {
	tx := mux.Vars(req)["tx"]
	rawTx, err := hex.DecodeString(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	txHash, blockHash, err := broadcast(r, rawTx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// handle inputs
	addresses, keys, txHashes, err := getInputsFromTx(r, txHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, k := range keys {
		v, err := r.DB.Get(addresses[i], "unspent", txHashes[i], k)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = r.DB.Delete(addresses[i], "unspent", txHashes[i], k)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = r.DB.Update(addresses[i], "spent", txHashes[i], k, v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// handle outputs
	addresses, keys, values, err := getOutputsFromTx(r, txHash, blockHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, k := range keys {
		err = r.DB.Update(addresses[i], "unspent", txHash.String(), k, values[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp := txResponse{txHash.String()}
	json.NewEncoder(w).Encode(resp)
}

// GetUtxos returns the list of unpsent tx output for a given address
func (r *RegTest) GetUtxos(w http.ResponseWriter, req *http.Request) {
	address := mux.Vars(req)["address"]
	rawList, err := r.DB.List(address, "unspent")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	list := []*utxo{}
	for _, raw := range rawList {
		vout := &utxo{}
		err := json.Unmarshal(raw, vout)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		list = append(list, vout)
	}

	json.NewEncoder(w).Encode(list)
}

// EstimateFees queries the blockchain to get the min relay fee per KB
func (r *RegTest) EstimateFees(w http.ResponseWriter, req *http.Request) {
	fees, err := estimateFees()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(fees)
}

func mine(r *RegTest, num int) ([]*chainhash.Hash, error) {
	return r.Client.Generate(uint32(num))
}

func sendTo(r *RegTest, address string) (*chainhash.Hash, *chainhash.Hash, error) {
	receiver, err := btcutil.DecodeAddress(address, &chaincfg.RegressionNetParams)
	if err != nil {
		return nil, nil, err
	}
	txHash, err := r.Client.SendToAddress(receiver, btcutil.Amount(100000000))
	if err != nil {
		return nil, nil, err
	}

	blockHash, err := mine(r, 1)
	if err != nil {
		return nil, nil, err
	}

	return txHash, blockHash[0], nil
}

func broadcast(r *RegTest, tx []byte) (*chainhash.Hash, *chainhash.Hash, error) {
	rawTx := &wire.MsgTx{}
	err := rawTx.Deserialize(bytes.NewReader(tx))
	if err != nil {
		return nil, nil, err
	}

	txHash, err := r.Client.SendRawTransaction(rawTx, true)
	if err != nil {
		return nil, nil, err
	}

	blockHash, err := mine(r, 1)
	if err != nil {
		return nil, nil, err
	}

	return txHash, blockHash[0], nil
}

type estimateFeeResponse struct {
	High   int `json:"high_fee_per_kb"`
	Medium int `json:"medium_fee_per_kb"`
	Low    int `json:"low_fee_per_kb"`
}

// estimateFees makes an HTTP API call to the regtest deamon without using the RegTest client.
// This is because in the actual implementation of rpcclient does not support the `getnetworkinfo` method,
// where `relayFee` can be found. More on this can be found here http://chainquery.com/bitcoin-api/getnetworkinfo
func estimateFees() (*estimateFeeResponse, error) {
	url, headers, body := getEstimationRequestParams()

	statusCode, resp, err := httpPOST(url, body, headers)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, fmt.Errorf("An unexpected error occured while estimating fees, check the response: %s", resp)
	}

	type response struct {
		Result struct {
			Fee float64 `json:"relayFee"`
		} `json:"result"`
	}
	out := &response{}
	json.Unmarshal([]byte(resp), out)

	fee := int(out.Result.Fee * math.Pow10(8))
	return &estimateFeeResponse{fee, fee, fee}, nil
}

// decodeTx takes a tx hash and returns the decoded tx object
func decodeTx(r *RegTest, txHash *chainhash.Hash) (*btcjson.TxRawResult, error) {
	rawTx, err := r.Client.GetRawTransaction(txHash)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	err = rawTx.MsgTx().Serialize(buf)
	if err != nil {
		return nil, err
	}

	return r.Client.DecodeRawTransaction(buf.Bytes())
}

// prepareUtxo takes a tx object, finds the utxo of to the given address in the tx outputs
//  and returns the key/value pair to be saved in the db
func prepareUtxo(r *RegTest, txHash *chainhash.Hash, address string, blockHash *chainhash.Hash) (string, string, error) {
	txObject, err := decodeTx(r, txHash)
	if err != nil {
		return "", "", err
	}

	var key string
	var value string
	for _, v := range txObject.Vout {
		if v.ScriptPubKey.Addresses[0] == address {
			key = getUtxoKey(txHash.String(), int(v.N))
			vout, _ := marshalUtxo(v, txHash.String(), blockHash.String())
			value = string(vout)
		}
	}

	return key, value, nil
}

// getInputsFromTx reconstructs the tx object from the hash and then, for each input:
// - gets the address who owns the input
// - calculates the utxo id used in the database (key)
// - gets the input tx hash
// this functions returns 3 arrays containg input addresses, keys and txHashes, or an error
func getInputsFromTx(r *RegTest, txHash *chainhash.Hash) ([]string, []string, []string, error) {
	txObject, err := decodeTx(r, txHash)
	if err != nil {
		return nil, nil, nil, err
	}

	keys := []string{}
	addresses := []string{}
	txHashes := []string{}
	for _, v := range txObject.Vin {
		key := getUtxoKey(v.Txid, int(v.Vout))
		txHash, _ := chainhash.NewHashFromStr(v.Txid)
		tx, _ := decodeTx(r, txHash)
		address, _ := getAddressFromUtxo(tx, v.Vout)

		keys = append(keys, key)
		addresses = append(addresses, address)
		txHashes = append(txHashes, v.Txid)
	}

	return addresses, keys, txHashes, nil
}

// getOutpytsFromTx similarly to the above function deserializes the tx from the hash and for each output:
// - gets the receiving address
// - calculates the utxo id (key)
// - stringifies the entire utxo JSON object (value)
// This returns 3 arrays of addresses, keys and values, or an error
func getOutputsFromTx(r *RegTest, txHash *chainhash.Hash, blockHash *chainhash.Hash) ([]string, []string, []string, error) {
	txObject, err := decodeTx(r, txHash)
	if err != nil {
		return nil, nil, nil, err
	}

	keys := []string{}
	values := []string{}
	addresses := []string{}
	for _, vout := range txObject.Vout {
		if vout.Value > 0 {
			address := vout.ScriptPubKey.Addresses[0]
			key := getUtxoKey(txHash.String(), int(vout.N))
			value, _ := marshalUtxo(vout, txHash.String(), blockHash.String())

			addresses = append(addresses, address)
			keys = append(keys, key)
			values = append(values, string(value))
		}
	}

	return addresses, keys, values, nil
}

// getUtxoKey calculates the double_hash(utxoTxHash + string(utxoTxIndex * 100))
// This hash is used as the utxo unique identifier in the database
func getUtxoKey(hash string, nout int) string {
	message, _ := hex.DecodeString(hash + strconv.Itoa(nout*100))
	return hex.EncodeToString(chainhash.DoubleHashB(message))
}

// getAddressFromUtxo returns the address that is the receiver for the NOUTth output of the given transaction
func getAddressFromUtxo(tx *btcjson.TxRawResult, nout uint32) (string, error) {
	if utxo := tx.Vout[nout]; utxo.N == nout {
		return utxo.ScriptPubKey.Addresses[0], nil
	}

	for _, utxo := range tx.Vout {
		if utxo.N == nout {
			return utxo.ScriptPubKey.Addresses[0], nil
		}
	}

	return "", errors.New("Error while getting utxo: Address not found")
}

type utxo struct {
	Index        uint32 `json:"index"`
	Value        uint32 `json:"value"`
	ScriptPubKey string `json:"script_pubkey"`
	Hash         string `json:"hash"`
	BlockHash    string `json:"block_hash"`
}

func marshalUtxo(v btcjson.Vout, hash string, blockHash string) ([]byte, error) {
	satoshis := uint32(v.Value * math.Pow10(8))
	c := utxo{v.N, satoshis, v.ScriptPubKey.Hex, hash, blockHash}
	return json.Marshal(c)
}

func getEstimationRequestParams() (string, map[string]string, string) {
	url := fmt.Sprintf("http://%s:%s@%s", user, password, host)
	header := map[string]string{"Content-Type": "application/json"}
	body := `{"jsonrpc": "1.0", "id": "2", "method": "getnetworkinfo", "params": []}`

	return url, header, body
}
