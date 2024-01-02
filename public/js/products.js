let cartIcon = document.querySelector('#cart-icon');
let cart = document.querySelector('.cart');
let closeCart = document.querySelector('#close-cart');

// open cart
cartIcon.onclick = () => {
    cart.classList.add('active');
}
// close cart
closeCart.onclick = () => {
    cart.classList.remove('active');
}

let com = document.querySelector('#com');
let comments = document.querySelector('.comments');
let closeCom = document.querySelector('#close-com');
// open comments
com.onclick = () => {
    comments.classList.add('active');
}
// close comments
closeCom.onclick = () => {
    comments.classList.remove('active');
}

function Getproducts(kind , container , id , callfunction){
  
    fetch(`/product?kind=${kind}&container=${container}&id=${id}`, {
        method: "GET",
    }).then(Response => Response.json()).then(data => {
      callfunction(data)
    })
}
