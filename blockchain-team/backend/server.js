const express = require('express');
const cors = require('cors');
const { getMockRecommendation } = require('./recommendation/engine');
const { generateSHA256 } = require('./utils/hashProvider');
const { storeHash, readHash } = require('./blockchain/fabricService');

const app = express();

// Middleware
app.use(cors());
app.use(express.json());

// ---------------------------------------------------------
// 1. RECOMMENDATION ROUTE (Frontend -> Backend -> Mock ML)
// ---------------------------------------------------------
app.post('/recommend', (req, res) => {
    // Generate ML output
    const evaluation = getMockRecommendation(req.body);
    // Generate Cryptographic Hash of that output
    const dataHash = generateSHA256(evaluation);
    
    console.log(`[API] /recommend called. Generated Hash: ${dataHash.substring(0, 15)}...`);
    res.status(200).json({ analysis: evaluation, hash: dataHash });
});

// ---------------------------------------------------------
// 2. STORE ROUTE (Backend -> Blockchain)
// ---------------------------------------------------------
app.post('/storeHash', async (req, res) => {
    const { recordID, hash } = req.body;
    
    // Failsafe for empty Postman requests
    if (!recordID || !hash) {
        console.log(`[API ERROR] /storeHash missing data. Check Postman JSON settings.`);
        return res.status(400).json({ error: "Missing recordID or hash." });
    }

    const result = await storeHash(recordID, hash);
    console.log(`[LEDGER] Successfully committed Record: ${recordID}`);
    
    res.status(200).json({ 
        success: true, 
        transactionId: result.txID,
        message: "Hash securely stored on mock ledger" 
    });
});

// ---------------------------------------------------------
// 3. VERIFY ROUTE (Auditing System)
// ---------------------------------------------------------
app.post('/verify', async (req, res) => {
    const { recordID, currentData } = req.body;
    
    const storedHash = await readHash(recordID);
    
    if (!storedHash) {
        console.log(`[VERIFY ERROR] Record ${recordID} not found in ledger.`);
        return res.status(404).json({ status: "ERROR", message: "Record not found on blockchain." });
    }

    const calculatedHash = generateSHA256(currentData);

    // X-Ray Logging for your team presentation
    console.log("\n================ BLOCKCHAIN AUDIT ================");
    console.log(`Record ID        : ${recordID}`);
    console.log(`On-Chain Hash    : ${storedHash}`);
    console.log(`Calculated Hash  : ${calculatedHash}`);
    
    if (calculatedHash === storedHash) {
        console.log(`Result           : ✅ VALID (Data is Authentic)`);
        console.log("==================================================\n");
        res.status(200).json({ status: "VALID", message: "Cryptographic signature matches blockchain state." });
    } else {
        console.log(`Result           : ❌ TAMPERED (Hashes mismatch)`);
        console.log("==================================================\n");
        res.status(200).json({ status: "TAMPERED", message: "Security Alert: Payload does not match on-chain hash." });
    }
});

// ---------------------------------------------------------
// SERVER BOOTSTRAP
// ---------------------------------------------------------
const PORT = 3000;
app.listen(PORT, () => {
    console.log(`\n🚀 Healthcare API Integration Node Online`);
    console.log(`Listening on http://localhost:${PORT}\n`);
});