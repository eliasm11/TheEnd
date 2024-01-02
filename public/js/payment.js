function PostVisa() {
    const Visanumber = document.getElementById("Visa_number").value
    console.log(Visanumber)
    const cvv = document.getElementById("cvv").value
    console.log(cvv)
    const day = document.getElementById("day").value
    console.log(day)
    const month = document.getElementById("month").value
    console.log(month)
    const cardholdername = document.getElementById("cardholder-name").value
    fetch("/user/visa", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: `{"Visanumber":"${Visanumber}","cvv":"${cvv}","day":"${day}" , month":"${month}"}`,
    }).then(Response => Response.text()).then(data => {

    })
    const mobilenumber = document.getElementById("mobile-number").value
    console.log(mobilenumber)
    fetch("/user/phone", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `mobilenumber=${mobilenumber}`,
    }).then(Response => Response.text()).then(data => {

    })
    
    console.log(cardholdername)
    fetch("/user/name", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `cardholdername=${cardholdername}`,
    }).then(Response => Response.text()).then(data => {

    })
    // Get values from the input fields
const city = document.getElementById('city').value;
const street = document.getElementById('street').value;
const apartment = document.getElementById('apartment').value;
const building = document.getElementById('building').value;

// Construct the request body
const requestBody = `city=${encodeURIComponent(city)}&street=${encodeURIComponent(street)}&apartment=${encodeURIComponent(apartment)}&building=${encodeURIComponent(building)}`;

// Use the fetch function to make a POST request
fetch("/user/addr", {
    method: "POST",
    headers: {
        "Content-Type": "application/x-www-form-urlencoded",
    },
    body: requestBody,

}).then(Response => Response.text()).then(data => {

})
    
    

}
