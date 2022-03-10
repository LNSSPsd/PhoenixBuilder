// example05.js
// 本脚本演示了fetch功能

// fetch
FB_Println("Start...")
var x = fetch('https://storage.fastbuilder.pro').then(function(r) {
    r.text().then(function(d) {
        FB_Println(r.statusText)
        for (var k in r.headers._headers) {
            FB_Println(k + ':', r.headers.get(k))
        }
        FB_Println(d)
    });
});

FB_Println("Awaiting...")

// setTimeout, clearTimeout, setInterval and clearInterval
setTimeout(function (){
    FB_Println("Timeout 10s")
},1000)

//  atob and btoa
base64encodedString=btoa("raw string")
recoveredString=atob(base64encodedString)
FB_Println(base64encodedString)
FB_Println(recoveredString)

// URL and URLSearchParams
// URL.revokeObjectURL()