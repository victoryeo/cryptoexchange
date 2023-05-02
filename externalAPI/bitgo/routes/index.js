const express       = require('express')
const BitGo = require('bitgo')
const { BitGoAPI } = require('@bitgo/sdk-api');
const { Tbtc } = require('@bitgo/sdk-coin-btc');
require('dotenv').config()

// init the sdk
const newbitgo = new BitGoAPI({ env: 'test' });
//const bitgo = new BitGo.BitGo({ env: process.env.mode,
//    accessToken: process.env.accessToken4 });

// register coin
newbitgo.register('tbtc', Tbtc.createInstance);

const router = express.Router()
const app = express()

router.get('/', (req, res) => {
    console.log("Hello Bitgo")
    res.setHeader('Content-Type', 'application/json')
    res.write(JSON.stringify({title:"Hello Bitgo"}));
    res.end();
})

router.post('/', async (req, res) => {
    console.log("Init Bitgo")
    //let result = await bitgo.session();
    //console.dir(result);
    console.log(process.env.username)
    const auth_res = await newbitgo.authenticate({
      username: process.env.username,
      password: process.env.password,
      otp: "000000",
    });
    const access_token = await newbitgo.addAccessToken({
      otp: "000000",
      label: "Admin Access Token",
      scope: [
        "metamask_institutional",
        "openid",
        "pending_approval_update",
        "portfolio_view",
        "profile",
        "trade_trade",
        "trade_view",
        "wallet_approve_all",
        "wallet_create",
        "wallet_edit_all",
        "wallet_manage_all",
        "wallet_spend_all",
        "wallet_view_all",
      ],
      // Optional: Set a spending limit.
      spendingLimits: [
        {
          coin: "tbtc",
          txValueLimit: "1000000000", // 10 TBTC (10 * 1e8)
        },
      ],
    });
    console.log(access_token);

    /*const btc_params = {
      "passphrase": "hellobitgo",
      "label": "firstwallet"
    };
    // create a btc wallet
    const { wallet } = await bitgo.coin('tbtc').wallets().generateWallet(btc_params);
    console.dir(wallet);
    app.locals.wallet = wallet*/

    app.locals.holder = "testwallet"
    res.setHeader('Content-Type', 'application/json')
    res.write(JSON.stringify({title:"Init Bitgo"}));
    res.end();
})

router.post('/send/:dest', async (req, res) => {
    console.log("send order")
    console.log(req.params.dest); 
    console.log(req.body); 
  
    wallet = app.locals.wallet
    // send crypto to another address
    result = await wallet.send({
      address: req.params.dest,
      amount: 0.01 * 1e8,
      walletPassphrase:  "hellobitgo"
    });

    holder = app.locals.holder
    console.log(holder)
    res.setHeader('Content-Type', 'application/json')
    res.write(JSON.stringify({address:req.params.dest,amount:req.body.amount}));
    res.end();
})
module.exports = router
