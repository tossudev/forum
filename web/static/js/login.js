const loginForm = document.querySelector('form.login')

loginForm.addEventListener('submit', function(e) {
    e.preventDefault();
    handleLogin();
})

// handleLogin sends user credentials to backend and handles the response from the server
async function handleLogin() {
    const formData = new FormData(loginForm);

    let response;

    try {
        // Send login form data to backend
        response = await fetch('/users/login', {
            method: 'POST',
            body: new URLSearchParams(formData)
        });
    } catch(err) {
        console.log('Error:', err)
    }

    // HTTP response code cases
    if (response.ok) {
        window.location.href = '/'; // redirect to home on successful login
        return;
    }

    if (response.status >= 500) {
        showError("Internal server error. Try again.")
        return;
    }

    showError("Invalid credentials. Try again.")
}

// showError displays an error message to the user if an error occurred during login
function showError(errMsg) {
    const errElem = document.querySelector('.error-msg')
    errElem.textContent = errMsg
}