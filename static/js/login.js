var userRe = /[^a-zA-Z\d]/
var emailRe = /.+\@.+\..+/

$("#close").click(function(e) {
    e.preventDefault();
    $("#incorrect").hide();
})

$("#login-submit").click(function(e) {
    e.preventDefault();
    uname = $("#login-username").val().trim();
    pword = $("#login-password").val();
    data = JSON.stringify({"uname": uname, "pword": pword});
    $.post("/login", data, function(ret) {
        if (ret.success){
            if ($("#login-check").prop("checked")) {
                d = new Date();
                d.setTime(d.getTime() + (7 * 24 * 60 * 60 * 1000));
                document.cookie = "key=" + ret.key + ";expires=" + d.toUTCString()+";path=/";
            } else {
                document.cookie = "key=" + ret.key + ";path=/";
            }
            window.location = "/";
        } else {
            $("#incorrect").show();
        }
    });
});

$("#register-submit").click(function(e) {
    e.preventDefault();
    valid = true;
    uname = $("#register-username").val().trim();
    pword = $("#register-password").val();
    email = $("#register-email").val().trim();
    if (userRe.exec(uname) || uname.length < 4 || uname.length > 20){
        $("#register-username").addClass("is-invalid");
        valid = false;
    } else {
        $("#register-username").removeClass("is-invalid");
        $("#register-username").addClass("is-valid");
    }
    if (pword.length < 8 || pword.length > 256) {
        $("#register-password").addClass("is-invalid");
        valid = false;
    } else {
        $("#register-password").removeClass("is-invalid");
        $("#register-password").addClass("is-valid");
    }
    if (email.length > 320 || !emailRe.exec(email)) {
        $("#register-email").addClass("is-invalid");
        valid = false;
    } else {
        $("#register-email").removeClass("is-invalid");
        $("#register-email").addClass("is-valid");
    }
    if (!valid) return;
    data = JSON.stringify({"uname": uname, "pword": pword, "email": email});
    $.post("/register", data, function(ret){
        if (ret.success) {
            if ($("#register-check").prop("checked")) {
                d = new Date();
                d.setTime(d.getTime() + (7 * 24 * 60 * 60 * 1000));
                document.cookie = "key=" + ret.key + ";expires=" + d.toUTCString()+";path=/";
            } else {
                document.cookie = "key=" + ret.key + ";path=/";
            }
            window.location = "/";
        }
        if (ret.error == "invu" || ret.error == "ulen" || ret.error == "utake") {
            $("#register-username").removeClass("is-valid");
            if (ret.error == "utake") {
                $("#uinv").text("This username is already taken");
                $("#register-username").addClass("is-invalid");
            } else {
                $("#uinv").text("Usernames are alphanumeric and must be between 4-20 characters inclusive")
                $("#register-username").addClass("is-invalid");
            }
        } else {
            $("#register-username").removeClass("is-invalid");
            $("#register-username").addClass("is-valid");
        }
        if (ret.error == "invp" || ret.error == "plen") {
            $("#register-password").removeClass("is-valid");
            $("#register-password").addClass("is-invalid");
        } else {
            $("#register-password").removeClass("is-invalid");
            $("#register-password").addClass("is-valid");
        }
        if (ret.error == "inve" || ret.error == "elen") {
            $("#register-email").removeClass("is-valid");
            $("#register-email").addClass("is-invalid");
        } else {
            $("#register-email").removeClass("is-invalid");
            $("#register-email").addClass("is-valid");
        }
    });
});