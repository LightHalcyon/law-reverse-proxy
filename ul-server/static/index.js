// const url = process.env.PROXYURL + "/ul";
// dec2hex :: Integer -> String
function dec2hex (dec) {
    return ('0' + dec.toString(16)).substr(-2)
}
  
// generateId :: Integer -> String
function generateId (len) {
    var arr = new Uint8Array((len || 40) / 2)
    window.crypto.getRandomValues(arr)
    return Array.from(arr, dec2hex).join('')
}

const host = window.location.href.replace("http://", '');

const url = 'http://localhost:23061/';

var id = ""

const form = document.querySelector('form');

form.addEventListener('submit', e => {
    e.preventDefault();

    id = generateId()
    console.log(id)
    const files = document.querySelector('[type=file]').files[0];
    const formData = new FormData();
    formData.append('file', files);
    fetch(url, {
        method: 'POST',
        headers: {
            "X-ROUTING-KEY": id,
        },
        body: formData,
    }).then(function(response) {
        return response.json();
    }).then(function(json) {
        document.getElementById("status").innerHTML = '<a href="' + json.message+ '">Download here</a>'
        console.log(JSON.stringify(json));
    });
});