$(document).ready(function(){
    $(".dynamic-container").on("click", ".dynamic-container-remove", function(){
        $(this).parent().parent().parent().remove();
    });

    let condition_dropdown = $("#condition-select-dropdown");

    $("#condition-select-add").on("click", function() {
        let type = condition_dropdown.val();
        console.log("Condition Select Pressed: " + type);
        let template = $("#condition-template-" + type).clone();
        template.removeAttr("id");

        let count = $("#conditions-container > .conditions-condition").length;
        replaceIndex(template, count);

        template.appendTo("#conditions-container");
    });

    let response_dropdown = $("#response-select-dropdown");

    $("#response-select-add").on("click", function() {
        let type = response_dropdown.val();
        console.log("Response Select Pressed: " + type);
        let template = $("#response-template-" + type).clone();
        template.removeAttr("id");

        let count = $("#responses-container > .responses-response").length;
        replaceIndex(template, count);

        template.appendTo("#responses-container");
    });

    function replaceIndex(objRef, index) {
        objRef.find(".row").children().each(function(){
            $(this).children().each(function(){
                let idAttr = $(this).attr("id");

                if (typeof(idAttr) !== "undefined") {
                    $(this).attr("id", idAttr.replace("{INDEX}", index));
                }

                let forAttr = $(this).attr("for");

                if (typeof(forAttr) !== "undefined") {
                    $(this).attr("for", forAttr.replace("{INDEX}", index));
                }

                let nameAttr = $(this).attr("name");

                if (typeof(nameAttr) !== "undefined") {
                    $(this).attr("name", nameAttr.replace("{INDEX}", index));
                }
            });
        });
    }
});
