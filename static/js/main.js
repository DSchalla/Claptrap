$(document).ready(function(){
   $(".dynamic-container-remove").on("click", function(){
       $(this).parent().parent().parent().remove();
   });
});
