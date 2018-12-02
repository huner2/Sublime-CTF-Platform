$("#logout").click(function(e) {
    e.preventDefault();
    document.cookie = "key=;expires=Thu, 01 Jan 1970 00:00:01 GMT;";
    document.location = "/";
});