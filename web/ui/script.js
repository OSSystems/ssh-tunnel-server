window.onload = function() {
    Terminal.applyAddon(fit);

    var term = new Terminal();
    term.open(document.getElementById("terminal"));
    term.fit();

    var ws = new WebSocket("ws://" + location.host + "/ws" + window.location.search);

    ws.onopen = function() {
    }
    
    ws.onmessage = function(e) {
        term.write(e.data);
    }

    term.on("data", function(data) {
        ws.send(data);
    });
}
