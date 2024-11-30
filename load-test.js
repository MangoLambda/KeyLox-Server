import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
    vus: 100, // Number of virtual users
    duration: '30s', // Duration of the test
};

export default function () {
    http.post('http://192.168.2.103:8080/register', JSON.stringify({
        clientSalt: generateRandomBase64String(32),
        key: generateRandomBase64String(32),
        username: generateRandomUsername()
    }), {
        headers: { 'Content-Type': 'application/json' }
    });
    sleep(0.01);
}

function generateRandomUsername() {
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let username = 'User_';
    let length = Math.floor(Math.random() * 64)
    for (let i = 0; i < length; i++) {
        username += characters.charAt(Math.floor(Math.random() * characters.length));
    }
    return username;
}

function generateRandomBase64String(length) {
    // Create a random byte array of the desired length
    const bytes = new Uint8Array(length);
    for (let i = 0; i < length; i++) {
        bytes[i] = Math.floor(Math.random() * 256); // Random byte between 0 and 255
    }

    // Convert byte array to Base64 manually
    return base64Encode(bytes);
}

// Custom Base64 encoding function
function base64Encode(array) {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
    let result = '';
    let i = 0;
    while (i < array.length) {
        let byte1 = array[i++];
        let byte2 = i < array.length ? array[i++] : 0;
        let byte3 = i < array.length ? array[i++] : 0;

        result += chars[byte1 >> 2];
        result += chars[((byte1 & 3) << 4) | (byte2 >> 4)];
        result += chars[((byte2 & 15) << 2) | (byte3 >> 6)];
        result += chars[byte3 & 63];

        // Handle padding
        if (i > array.length) {
            result = result.slice(0, -1) + '=';
        } else if (i === array.length - 1) {
            result = result.slice(0, -2) + '==';
        }
    }
    return result;
}
