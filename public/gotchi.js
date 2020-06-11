window.Twitch.ext.onAuthorized((auth) => {
    userid = auth.userId;
    userid = userid.replace(/[^\d.-]/g, '');
    let socket = new WebSocket("ws://localhost:8081/ws/" + userid)
    socket.onmessage = function(event) {
        console.log(event);
    };
});

