//Intl = Internationalization. This js object holds tools for handling timezone related things. 
//DateTimeFormat = Creates a formatter. Since here we dont have parameters in the func, it uses the browsers own settings (language, timezone etc)
//resolvedOptions = asks formatter which settings did we end up using, returning the stuff that are in use, one of them being the timezone
const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone


//document = js object, which has a property of cookie, which reads and writes the pages cookies.
//tz = will be the timezone const above. Here in Helsinki it'll be Europe/Helsinki
//path sets the cookies on all paths in this page, not just the homepage
//max-age is the cookies lifetime in seconds, this being one year
document.cookie = `tz=${timezone}; path=/; max-age=31536000;`
