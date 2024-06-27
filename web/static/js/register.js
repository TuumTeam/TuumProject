const container = document.getElementById('container');
const registerBtn = document.getElementById('register');
const loginBtn = document.getElementById('login');

document.addEventListener('keydown', function(event) {
    if (event.keyCode === 13) {
        event.preventDefault();
    }
});

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

async function hash(string) {
    const utf8 = new TextEncoder().encode(string);
    const hashBuffer = await crypto.subtle.digest('SHA-256', utf8);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray
        .map((bytes) => bytes.toString(16).padStart(2, '0'))
        .join('');
    return hashHex;
}

async function hashSubmit(logType) {
    if (!validateForm()) {
        alert("Ca marche pas.");
        return;
    }
    document.getElementsByName("hash")[logType].value = await hash(document.getElementsByName("password")[logType].value);
    switch (logType) {
        case 0:
            document.getElementById("registerForm").submit();
            break;
        case 1:
            document.getElementById("loginForm").submit();
            break;
    }
}

function validateForm() {
    // Get form elements
    var username = document.getElementById('username').value;
    var email = document.getElementById('email').value;
    var password = document.getElementById('password').value;

    // Check if all fields are filled
    if (username === "" || email === "" || password === "") {
        alert("All fields must be filled out");
        return false;
    }

    // Check if email is valid
    var emailRegex = /^[a-zA-Z0-9._%+-]+@.+(\.com|\.fr)$/;
    if (!emailRegex.test(email)) {
        alert("Email address is not valid");
        return false;
    }

    // Check if password is longer than 2 characters
    if (password.length <= 2) {
        alert("Password is too short or incorrect");
        return false;
    }

    // If all checks pass, return true
    return true;
}

document.getElementById('registerForm').addEventListener('submit', function(event) {
    if (!validateForm()) {
        event.preventDefault();
    }
});

document.getElementById('loginForm').addEventListener('submit', function(event) {
    if (!validateForm()) {
        event.preventDefault();
    }
});