const fs = require('fs');
const path = require('path');

// Points to a physical file to act as our Hyperledger World State
const ledgerPath = path.join(__dirname, 'mockLedger.json');

/**
 * Helper function to read the physical ledger file.
 * If it doesn't exist, it creates it automatically.
 */
function getLedger() {
    if (!fs.existsSync(ledgerPath)) {
        fs.writeFileSync(ledgerPath, JSON.stringify({}, null, 2));
    }
    return JSON.parse(fs.readFileSync(ledgerPath, 'utf8'));
}

/**
 * Simulates committing a transaction to Hyperledger Fabric.
 */
async function storeHash(recordID, dataHash) {
    const ledger = getLedger();
    ledger[recordID] = dataHash; 
    
    // Write permanently to the hard drive
    fs.writeFileSync(ledgerPath, JSON.stringify(ledger, null, 2));
    return { success: true, txID: `tx_${Date.now()}` };
}

/**
 * Simulates querying the ledger for a specific record.
 */
async function readHash(recordID) {
    const ledger = getLedger();
    return ledger[recordID] || null; 
}

module.exports = { storeHash, readHash };