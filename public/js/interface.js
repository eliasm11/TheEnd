const login = document.querySelector('.login');
const loginLink = document.querySelector('.login-link');
const registerLink = document.querySelector('.login-register');
const headerloginIcon = document.querySelector('.headerlogin-icon');
const closeLogin = document.querySelector('.close-login');

registerLink.addEventListener('click', ()=> {
    login.classList.add('active');
});

loginLink.addEventListener('click', ()=> {
    login.classList.remove('active');
});

headerloginIcon.addEventListener('click', ()=> {
    login.classList.add('active-popup');
});

closeLogin.addEventListener('click', ()=> {
    login.classList.remove('active-popup');
});