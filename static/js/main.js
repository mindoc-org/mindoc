(function () {
    $("[data-toggle='tooltip']").tooltip();
})();

function showError($msg) {
    $("#form-error-message").addClass("error-message").removeClass("success-message").text($msg);
    return false;
}

function showSuccess($msg) {
    $("#form-error-message").addClass("success-message").removeClass("error-message").text($msg);
    return true;
}