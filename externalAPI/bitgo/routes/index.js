const express       = require('express')
const BitGo = require('bitgo')
require('dotenv').config()

// init the sdk
const bitgo = new BitGo.BitGo({ env: process.env.mode,
    accessToken: process.env.accessToken4 });
  
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
    let result = await bitgo.session();
    console.dir(result);
  
    const btc_params = {
      "passphrase": "hellobitgo",
      "label": "firstwallet"
    };
    // create a btc wallet
    const { wallet } = await bitgo.coin('tbtc').wallets().generateWallet(btc_params);
    console.dir(wallet);

    app.locals.wallet = wallet
})

router.post('/send/:dest', async (req, res) => {
    console.log("send order")
    console.log(req.body); 
  
    wallet = app.locals.wallet
    // send crypto to another address
    result = await wallet.send({
      address: req.params.dest,
      amount: 0.01 * 1e8,
      walletPassphrase:  "hellobitgo"
    });

    res.setHeader('Content-Type', 'application/json')
    res.write(req.body);
    res.end();
})
module.exports = router
