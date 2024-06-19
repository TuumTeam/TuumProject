const container = document.getElementById('container');
const registerBtn = document.getElementById('register');
const loginBtn = document.getElementById('login');

registerBtn.addEventListener('click', () => {
    container.classList.add("active");
});

loginBtn.addEventListener('click', () => {
    container.classList.remove("active");
});

document.getElementById('register').addEventListener('click', function() {
    document.body.classList.add('gradient-reversed');
});

document.getElementById('login').addEventListener('click', function() {
    document.body.classList.remove('gradient-reversed');
});