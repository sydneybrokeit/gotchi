window.Twitch.ext.onAuthorized((auth) => {
    userid = auth.userId;
    userid = userid.replace(/[^\d.-]/g, '');
    let socket = new WebSocket("ws://localhost:8081/ws/" + userid)
    socket.onmessage = function(event) {
        // console.log(event);
        eventObj = JSON.parse(event.data)
        switch(eventObj.type) {
            case "REFRESH":
                handleRefresh(eventObj.state)
            case "DELTA":
                handleDelta(eventObj.delta)
        }
    };
});

function handleRefresh(state) {
    console.log("Received REFRESH STATE message");
    console.log(state);
}

function handleDelta(delta) {
    console.log("Received DELTA message");
    console.log(delta);
}



