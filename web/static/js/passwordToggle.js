// showPassword toggles the password visibility of the password field
function showPassword() {
    var pwField = document.getElementById("password");

    if (pwField.type === "password") {
        pwField.type = "text";
    } else  {
        pwField.type = "password";
    }
}