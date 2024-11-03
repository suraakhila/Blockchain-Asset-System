const express = require('express');
const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

const app = express();
app.use(express.json());

const ccpPath = path.resolve(__dirname, '../connection.json'); // Path to your connection profile

app.post('/createAsset', async (req, res) => {
    try {
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        const gateway = new Gateway();
        await gateway.connect(ccpPath, {
            wallet,
            identity: 'user1',
            discovery: { enabled: true, asLocalhost: true }
        });

        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('asset_management');

        await contract.submitTransaction('CreateAsset', req.body.dealerId, req.body.msisdn, req.body.mpin, req.body.balance, req.body.status, req.body.transAmount, req.body.transType, req.body.remarks);
        res.status(200).json({ response: 'Asset created successfully' });
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

app.listen(3000, () => {
    console.log('Server is running on port 3000');
});
