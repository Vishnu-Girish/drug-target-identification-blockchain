const crypto = require('crypto');

/**
 * Generates a SHA-256 hash. 
 * For the demo, we use a straightforward stringify to prevent 
 * minor JSON formatting differences in Postman from breaking the hash.
 */
function generateSHA256(data) {
    const stringifiedData = JSON.stringify(data);
    return crypto.createHash('sha256').update(stringifiedData).digest('hex');
}

module.exports = { generateSHA256 };