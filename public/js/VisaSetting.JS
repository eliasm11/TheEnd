function PostVisaSetting(){
    const Visanumber = document.getElementById("Visa_number").value
    console.log(Visanumber)
    const cvv = document.getElementById("cvv").value
    console.log(cvv)
    const expirydate = document.getElementById("expirydate").value
    console.log(expirydate)
    const cardholdername1 = document.getElementById("cardholder-name").value   
 fetch("/user/visa", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: `{"Visa_number":"${Visanumber}","cvv":"${cvv}","expirydate":"${expirydate}","cardholder-name":"${cardholdername1}}`,
    }).then(Response => Response.text()).then(data => {

    })

}