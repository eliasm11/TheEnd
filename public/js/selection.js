// Object.getOwnPropertyNames(data).forEach((val, idx, array) => {
//     console.log(`${val} -> ${data[val]}`);
// });
fetch("/AllContainerAndKind", {
    method: "GET"
}).then(res => res.json()).then(data => {
    var listallContainer = Object.getOwnPropertyNames(data);
    listallContainer.forEach(element => {
        if (document.getElementById("m").innerText === element) {
            return
        }
        var li = document.createElement("li")
        var a = document.createElement("a")
        a.innerText = element
        a.setAttribute("href", "/" + element)
        li.appendChild(a)
        document.getElementById("navlist").appendChild(li)
    });
})


