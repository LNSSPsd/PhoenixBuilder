// example05.js
// 本脚本演示了fetch功能

// fetch
engine.message("Start...")
var x = fetch('https://storage.fastbuilder.pro').then(function(r) {
    r.text().then(function(d) {
        FB_Println(r.statusText)
        for (var k in r.headers._headers) {
            FB_Println(k + ':', r.headers.get(k))
        }
        FB_Println(d)
    });
});

engine.message("Awaiting...")

// setTimeout, clearTimeout, setInterval and clearInterval
setTimeout(function (){
    engine.message("Timeout 10s")
},1000)

//  atob and btoa
let base64encodedString=btoa("raw string")
let recoveredString=atob(base64encodedString)
engine.message(base64encodedString)
engine.message(recoveredString)

// URL and URLSearchParams
// URL.revokeObjectURL()