import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
    vus: 100, // Number of virtual users
    duration: '30s', // Duration of the test
};

export default function () {
    let uri = 'http://192.168.2.103:8080/user/' + generateRandomUsername();
    http.get(uri);
    sleep(0.01);
}

function generateRandomUsername() {
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let username = 'User_';
    let length = 2
    for (let i = 0; i < length; i++) {
        username += characters.charAt(Math.floor(Math.random() * characters.length));
    }
    return username;
}

