const registerForm = document.querySelector('form.register')

registerForm.addEventListener('submit', function(e) {
    e.preventDefault();
    handleRegister();
})

// handleRegister sends user credentials to backend and handles the response from the server
async function handleRegister() {
    const formData = new FormData(registerForm);

    let response;

    try {
        // Send register form data to backend
        response = await fetch('/users/register', {
            method: 'POST',
            body: new URLSearchParams(formData)
        });
    } catch(err) {
        console.log('Error:', err)
    }

    // HTTP response code cases
    if (response.ok) {
        window.location.href = '/'; // redirect to home on successful registeration
        return;
    }

    if (response.status >= 500) {
        showError("Internal server error. Try again.")
        return;
    }

    showError("Invalid credentials. Try again.") // TODO: Show more specific error messages
}

// showError displays an error message to the user if an error occurred during registration
function showError(errMsg) {
    const errElem = document.querySelector('.error-msg')
    errElem.textContent = errMsg
}