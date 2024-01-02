function postadmin(){
    const username = document.getElementById("usernameLogin").value
    const password = document.getElementById("passwordLogin").value
    fetch("/admin/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
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