//  --------------------------------------upload manager--------------------------
function checkFile() {
    var fileInput = document.getElementById('image');
    var errorSpan = document.getElementById('error');
    var maxSize = 1*1024; // max file size in KB
    if (fileInput.files.length > 0) {
        var fileSize = fileInput.files[0].size / 1024; // Converting in KB
        if (fileSize > maxSize) {
            console.log(fileSize)
            errorSpan.textContent = "❌ file must not exceed 20MB";
            fileInput.value = null; // Réinitializing the field
            console.log("size checked")
        } else {
            errorSpan.textContent = null; // deleting the error message
        }
    }
}