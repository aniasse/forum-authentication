//  --------------------------------------upload manager--------------------------
function checkFile() {
    var fileInput = document.getElementById('image');
    var errorSpan = document.getElementById('error');
    var maxSize = 20*1024; // max file size in KB
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


// Fonction pour mettre à jour la prévisualisation de l'image
function previewImage() {
    var input = document.getElementById('image');
    var imagePreview = document.getElementById('imagePreview');

    if (input.files && input.files[0]) {
        var reader = new FileReader();

        reader.onload = function (e) {
            imagePreview.src = e.target.result;
            imagePreview.style.display = 'block';
        };

        reader.readAsDataURL(input.files[0]);
    } else {
        imagePreview.style.display = 'none';
    }
}

// Écouteur d'événements pour détecter les changements dans l'input file
var imageInput = document.getElementById('image');
imageInput.addEventListener('change', previewImage);



    
