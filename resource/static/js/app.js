window.onload = function () {
    console.log("Protocol: " + location.protocol);
    var wsURL = "ws://" + document.location.host + "/ws"
    if (location.protocol == 'https:') {
        wsURL = "wss://" + document.location.host + "/ws"
    }
    console.log("WS URL: " + wsURL);

    var panel = document.getElementById("character");
    if (panel) {
        sock = new WebSocket(wsURL);

        var connDiv = document.getElementById("connection-status");
        connDiv.innerText = "closed";

        sock.onopen = function () {
            console.log("connected to " + wsURL);
            connDiv.innerText = "open";
        };

        sock.onclose = function (e) {
            console.log("connection closed (" + e.code + ")");
            connDiv.innerText = "closed";
        };

        sock.onmessage = function (e) {
            console.log(e);
            var p = JSON.parse(e.data);
            console.log(p);

            /*
            {
                "name": "C-3PO",
                "height": "167",
                "mass": "75",
                "skin_color": "gold",
                "eye_color": "yellow",
                "birth_year": "112BBY",
            }
            */

            panel.innerHTML = `${p.name}, born: ${p.birth_year}, height: ${p.height}, eyes: ${p.eye_color}`;
        };
    } // if panel
};