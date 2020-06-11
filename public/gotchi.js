
console.log("in gotchi.js");
window.Twitch.ext.onAuthorized(function(auth) {
    console.log("authenticated");
    console.log(auth);
    window.Twitch.ext.listen('broadcast', function (topic, contentType, message) {
        console.log("listener got a message!");
        handleMessage(message);
    });
});



function handleMessage(message) {
    console.log(message);
}

function handleRefresh(state) {
    console.log("Received REFRESH STATE message");
    console.log(state);
}

function handleEvent(event) {
    console.log("Received EVENT message");
    console.log(event);
}



