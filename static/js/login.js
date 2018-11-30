$("#register-submit").click(function(e) {
    e.preventDefault();
    uname = $("#register-username").val();
    pword = $("#register-password").val();
    email = $("#register-email").val();
    data = JSON.stringify({"uname": uname, "pword": pword, "email": email});
    $.post("/register", data, function(ret){
        if (ret.success) {
            window.location = "/";
        }
    })
});