
















/*
function Getorders(){
    fetch("/admin/orders", {
        method: "GET",
    }).then(Response => Response.json()).then(data => {
    }) 
}
// link products.js with manage.html ?? whete <script src="js/products.js"></script>
function posproduct(){
    // Define the product data
const productData = {
    "idModel": 4,
    "color": "red",
    "size": "xl",
    "Container": "newContainer",
    "Kind": "short"
}; 
// Convert the product data into an array of objects
const productArray = Object.entries(productData).map(([key, value]) => ({ id: key, value }));

// Use the fetch function to make a POST request
fetch("/api/buy", {
    method: "POST",
    headers: {
        "Content-Type": "application/json",
    },
    body: JSON.stringify(productArray),
})
    .then(response => {return response.json();
    }).then(data => {
        // Handle the response data here, if needed
        console.log(data);
    })
}
// Needs to be linked to a file manage.html file
const containerValue = document.getElementById('container').value;
  const kindValue = document.getElementById('kind').value;
function PostAdminProducts(){
    const productData = {
        "idModel": 4,
        "color": "red",
        "size": "xl",
        "Container": "newContainer",
        "Kind": "short"
    };
    // Convert the product data into an array of objects
    
    const productArray = Object.entries(productData).map(([key, value]) => ({ id: key, value }));
    fetch("/admin/product", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(productArray),
    })
        .then(response => {return response.json();
        }).then(data => {
            // Handle the response data here, if needed
            console.log(data);
        })

// Needs to be linked to a file manage.html file
       // post kind
 const containerValue = document.getElementById('container').value;
  const kindValue = document.getElementById('kind').value;
  // Construct the JSON object
//   const data = {
//     container: containerValue,
//     kind: kindValue
//   };
  // Make a fetch request with the constructed JSON object in the body
  fetch('/admin/product/kind', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: `kind=${kindValue}`,
  })
  .then(response => response.json())
  .then(result => {
    // Handle the response result
    console.log(result);
  });

fetch("/admin/product/container", {
    method: "POST",
    headers: {
        "Content-Type": "application/x-www-form-urlencoded",
    },
    body: `container=${containerValue}`,
}).then(Response => Response.text()).then(data => {
    console.log(data)
})


  // POST MODEL 

  const postData = {
    container: "newContainer",
    kind: "short",
    model: {
      sizes: { xl: { red: 1 } },
      price: 12,
      description: "dswadas",
      discount: 0,
      linkesImage: ["dasd"]
    }
  };
  
  fetch("/admin/product/model", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(postData),
  })
    .then(response => response.json())
    .then(data => {
      console.log(data);
    })
  
}
// GET CONTAINER
function getcontainer(){
    fetch("/admin/product/container", {
        method: "GET"
    }).then(re => re.ok).then(d => {
            window.location.reload()
    })
}

//GET MODEL
function getmodel(){
    const deleteData = {
        container: "newContainer",
        kind: "short",
        id: 2
      };
      // Construct the URL with parameters
      fetch("/admin/product/model", {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(deleteData),
      })
      .then(response => response.json())
      .then(result => {
        // Handle the response result
        console.log(result);
      })
       
}
function deleteKind(){
  // delete kind 
  fetch("/admin/product/kind", {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json'
    },
    body: `kind=${kindValue}`,
  })
  .then(response => response.json())
  .then(result => {
    // Handle the response result
    console.log(result);
  })

}
function deletecontiners(){

// DELETE container
fetch("/admin/product/container", {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json'
    },
    body: `container=${containerValue}`,
  })
  .then(response => response.json())
  .then(result => {
    // Handle the response result
    console.log(result);
  })
}*/
