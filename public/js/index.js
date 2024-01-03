fetch("/AllContainerAndKind", {
    method: "GET"
}).then(res => res.json()).then(data => {
    if (data != "no data") {
        var listallContainer = Object.getOwnPropertyNames(data);
        listallContainer.forEach(element => {
            var li = document.createElement("li")
            var a = document.createElement("a")
            a.innerText = element
            a.setAttribute("href", "/" + element)
            li.appendChild(a)
            document.getElementById("navlist").appendChild(li)

            li = document.createElement("li")
            a = document.createElement("a")
            a.innerText = element
            a.setAttribute("href", "/" + element)
            li.appendChild(a)
            document.getElementById("main-contact").appendChild(li)
        });
    }
})


function logoutUser() {

    fetch("/user/logout", {
        method: "GET"
    }).then(re => re.ok).then(d => {
        window.location.reload()
    })

}


function loginIndex() {
    const username = document.getElementById("usernameLogin").value
    const password = document.getElementById("passwordLogin").value
    fetch("/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: `{"username":"${username}","password":"${password}"}`,
    }).then(Response => Response.json()).then(data => {

        if (data != "ok") {
            document.getElementById("spanErrorLogin").innerText = data
        } else {
            window.location.reload()
        }
    })

}

function ClearErrorSpan() {
    document.getElementById("spanErrorLogin").innerText = ""
    document.getElementById("spanErrorResister").innerText = ""
}



function PostSettings() {

    const firstname = document.getElementById("first-name").value
    const phone = document.getElementById("last-name").value
    const oldPassowrd = document.getElementById("current-password").value
    const newpassword = document.getElementById("new-password").value
    if (firstname.length != 0) {
        fetch("/user/name", {
            method: "POST",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            },
            body: `Name=${firstname}`,
        }).then(Response => Response.json()).then(data => {
            if (data != "update") {
                alert(data)
            } else {
                alert("Name Update")
            }

        })
    }
    if (phone.length != 0) {
        fetch("/user/phone", {
            method: "POST",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            },
            body: `phone=${phone}`,
        }).then(Response => Response.json()).then(data => {
            if (data != "update") {
                alert(data)
            } else {
                alert("Phone Update")
            }
        })
    }
    if (oldPassowrd.length == 0 || newpassword.length == 0) return
    fetch("/user/password", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: `{"oldPassowrd":"${oldPassowrd}","newPassowrd":"${newpassword}"}`,
    }).then(Response => Response.json()).then(data => {
        if (data != "update") {
            alert(data)
        } else {

            alert("Passowrd Update")
        }
    })
}


