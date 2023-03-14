import {randomIntBetween, randomString} from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import ws from 'k6/ws';
import {check} from 'k6';

const sessionDuration = randomIntBetween(10000, 60000); // user session between 10s and 1m

// export const options = {
//     vus: 100,
//     iterations: 10,
// };

export default function () {
    const url = `ws://ws.localhost/ws/null`;
    const params = {tags: {my_tag: 'my ws session'}};

    const res = ws.connect(url, params, function (socket) {
        socket.on('open', function open() {
            console.log(`VU ${__VU}: connected`);

            socket.send(JSON.stringify({event: 'SET_NAME', new_name: `Croc ${__VU}`}));

            socket.setInterval(function timeout() {
                socket.send(JSON.stringify({event: 'SAY', message: `I'm saying ${randomString(5)}`}));
            }, randomIntBetween(1000, 3000)); // say something every 1-3seconds
        });

        socket.on('close', function () {
            console.log(`VU ${__VU}: disconnected`);
        });

        socket.on('message', function (message) {
        });

        socket.setTimeout(function () {
            console.log(`VU ${__VU}: ${sessionDuration}ms passed, leaving the chat`);
            socket.send(JSON.stringify({event: 'LEAVE'}));
        }, sessionDuration);

        socket.setTimeout(function () {
            console.log(`Closing the socket forcefully 3s after graceful LEAVE`);
            socket.close();
        }, sessionDuration + 3000);
    });

    check(res, {'Connected successfully': (r) => r && r.status === 101});
}