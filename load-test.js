import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
    vus: 1000, // Number of virtual users
    duration: '10s', // Duration of the test
};

export default function () {
    http.get('http://localhost:8080/user');
    sleep(0.01);
}