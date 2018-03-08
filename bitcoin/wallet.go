package bitcoin

import (
	"github.com/OpenBazaar/multiwallet/client"
	"github.com/OpenBazaar/multiwallet/keys"
	wi"github.com/OpenBazaar/wallet-interface"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/OpenBazaar/multiwallet/wallet"
	"github.com/tyler-smith/go-bip39"
	hd "github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/btcutil"
)

type BitcoinWallet struct {
	db     wi.Datastore
	km     *keys.KeyManager
	params *chaincfg.Params
	client client.APIClient
	wm     *wallet.WalletManager
	mnemonic string
}

func NewBitcoinWallet(db wi.Datastore, mnemonic string, client client.APIClient, params *chaincfg.Params) (*BitcoinWallet, error) {
	seed := bip39.NewSeed(mnemonic, "")

	mPrivKey, err := hd.NewMaster(seed, params)
	if err != nil {
		return nil, err
	}
	km, err := keys.NewKeyManager(db.Keys(), params, mPrivKey, wi.Bitcoin)
	if err != nil {
		return nil, err
	}
	wm := wallet.NewWalletManager(db, km, client, params, wi.Bitcoin)

	return &BitcoinWallet{db, km, params, client, wm, mnemonic}, nil
}

func (w *BitcoinWallet) Start() {
	w.wm.Start()
}

func (w *BitcoinWallet) Mnemonic() string {
	return w.mnemonic
}


func (w *BitcoinWallet) CurrentAddress(purpose wi.KeyPurpose) btcutil.Address {
	key, _ := w.km.GetCurrentKey(purpose)
	addr, _ := key.Address(w.params)
	return btcutil.Address(addr)
}

func (w *BitcoinWallet) NewAddress(purpose wi.KeyPurpose) btcutil.Address {
	i, _ := w.db.Keys().GetUnused(purpose)
	key, _ := w.km.GenerateChildKey(purpose, uint32(i[1]))
	addr, _ := key.Address(w.params)
	w.db.Keys().MarkKeyAsUsed(addr.ScriptAddress())
	return btcutil.Address(addr)
}