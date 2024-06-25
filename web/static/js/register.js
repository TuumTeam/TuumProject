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
    let password = document.getElementsByName("password")[logType].value;
    document.getElementsByName("password")[logType].value = await hash(password);
    switch (logType) {
        case 0:
            document.getElementById("registerForm").submit();
            break;
        case 1:
            document.getElementById("loginForm").submit();
            break;
    }
}