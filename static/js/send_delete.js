function sendDelete(event, href){
    console.log("sendDelete called with href:", href); // Add this line
    var xhttp = new XMLHttpRequest();
    event.preventDefault();
    xhttp.onreadystatechange = function() {
        if (this.readyState !== 4) {
            return;
        }

        if (this.readyState === 4) {
            console.log("Redirecting to:", this.responseURL); // Add this line
            window.location.replace(this.responseURL);
        }
    };
    xhttp.open("DELETE", href, true);
    xhttp.send();
}