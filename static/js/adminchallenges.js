var cat;
var challid = 0;

$("#categories").on('click', '.cat-selector', function(e){
    e.preventDefault();
    var name = $(this).attr('id').substring(4);
    $(".cat-selector").each(function() {
        $(this).removeClass("active");
    });
    $(this).addClass("active");
    var challs = $(`#cat-challs-${name}`);
    $(".chall-group").each(function() {
        $(this).css("display", "none");
    })
    challs.css("display", "flex");
    $("#newchallform").css("display", "none");
})

$("#createCat").on("click", function(e){
    e.preventDefault();
    var name = $("#newCatName");
    name.removeClass("is-invalid");
    $.post("/admin/challenges", JSON.stringify({operation: "create", name: name.val(), type: "category"}), function(){
        $("#categories").append(
            "<a class=\"list-group-item list-group-item-action cat-selector\" href=\"#\" id=\"cat-" + name.val() + "\">" + name.val()
            + "<button type=\"button\" class=\"close\" id=\"delete-" + name.val() + "\"><span>&times;</span></button>" 
            + "</a>");
        $("#newModal").modal('hide');
    }).fail(function(res) {
        $("#incorrectText").text(res.responseText);
        $("#incorrect").show();
    });
})

$("#categories").on('click', '.cat-selector>.close', function(e){
    e.preventDefault();
    cat = $(e.target).parent().attr("id").split("-")[1];
    $("#catToDelete").text(cat);
    $("#deleteModal").modal('show');
})

$("#deleteCat").on("click", function(e) {
    e.preventDefault();
    $.post("/admin/challenges", JSON.stringify({operation: "delete", name: cat, type: "category"}), function(){
        $("#cat-" + cat).remove();
        $("#deleteModal").modal('hide');
    }).fail(function() {
        $("#error").show();
    });
})

$("#newchall").on("click", function(e) {
    e.preventDefault();
    $("#challenge-name").val("");
    $("#challenge-desc").val("");
    $("#challenge-flag").val("");
    $("#challenge-points").val(0);
    $("#newchallform").css("display", "flex");
});

$("#save-challenge").on("click", function(e) {
    e.preventDefault();
    $.post("/admin/challenges", JSON.stringify({
        operation: "update",
        name: $("#challenge-name").val(),
        desc: $("#challenge-desc").val(),
        flag: $("#challenge-flag").val(),
        points: $("#challenge-points").val(),
        cat: cat,
        id: challid,
        type: "challenge"
    })).fail(function(res) {
        $("#incorrectChallText").text(res.responseText);
        $("#incorrectChall").show();
    })
})