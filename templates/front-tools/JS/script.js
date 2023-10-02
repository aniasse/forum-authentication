// ------------------------- JS menu drop down ----------------------------

let profileMenu = document.getElementById("profileMenu");
function toggleMenu(){
    profileMenu.classList.toggle("open-menu");
}

// ---------------------- toggle commentaire --------------------------

// Sélectionnez tous les éléments de commentaire
const commentContainers = document.querySelectorAll(".comment-container");
console.log(commentContainers)
// Ajoutez des gestionnaires d'événements à chaque commentaire
togglers = document.querySelectorAll(".comment-section-toggler")

togglers.forEach(toggler=>{
        toggler.addEventListener("click",(e)=>{
            let id = toggler.dataset.post_id
        let test =document.querySelector(`#comment-container_${id}`)
            console.log(test)
            if (test.style.display === "none" || test.style.display === "") {
                test.style.display = "block";
            } else {
                test.style.display = "none";
            }
        })
})


// ------------------- comment input js ------------------------

// Récupérez les éléments "Comment" et le champ de commentaire
const commentInput = document.querySelectorAll(".post-activity-comment-input");

// replyButtons = document.querySelectorAll(".post-activity-link-comment")

commentButtons = document.querySelectorAll(".post-activity-link-comment")
commentButtons.forEach(commentButton=>{
    commentButton.addEventListener("click",(e)=>{
        let id_post = commentButton.dataset.post_id
    let commenting = document.querySelector(`#commenting_${id_post}`)
        if (commenting.style.display === "none" || commenting.style.display === "") {
            commenting.style.display = "block";
        } else {
            commenting.style.display = "none";
        }
    })
});


// ------------------ reply input --------------------------


// Récupérez les éléments "Comment" et le champ de commentaire
const replyInput = document.querySelectorAll(".post-activity-reply-input");


document.querySelectorAll(".post-activity-link-reply").forEach(replyButton=>{
    replyButton.addEventListener("click",(e)=>{
        let id_comment = replyButton.dataset.comment_id
    
    let replying = document.querySelector(`#replying_${id_comment}`)
        if (replying.style.display === "none" || replying.style.display === "") {
            replying.style.display = "block";
        } else {
            replying.style.display = "none";
        }
    })
});

// ---------------------- JS filter underline -------------------------

function setActiveFilter() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    const filter = urlParams.get('filter');

    if (filter) {
        const links = document.querySelectorAll("#activeFilter ul li a");
        links.forEach(link => {
            if (link.getAttribute('href').includes(filter)) {
                link.classList.add("active");
            }
        });
    }
}

// Appeler la fonction au chargement de la page
window.onload = setActiveFilter;

// ----------------------------- JS loginRequest -----------------------

const postLinks = document.querySelectorAll(".post-link");
const loginRequest = document.getElementById("logRequest");
const body = document.body;

postLinks.forEach(link => {
    link.addEventListener("click", () => {
        loginRequest.classList.toggle("open-log-request"); // Utilisez toggle pour ajouter/supprimer la classe
        body.classList.toggle("open-login-overlay"); // Ajoutez/supprimez la classe sur le body
    });
});

// Lorsque le formulaire de connexion est fermé
function closeLoginRequest() {
    loginRequest.classList.remove("open-log-request");
    body.classList.remove("open-login-overlay");
}
function closeLoginRequest() {
    loginRequest.classList.remove("open-log-request");
    body.classList.remove("open-login-overlay");
};


    

/* --------------------------------dark-mode------------------------------------------ */

    
    // Sélectionnez le bouton et le corps du document
    const darkModeToggle = document.getElementById('dark-mode-toggle');
    const Body = document.body;
    
    // Fonction pour activer/désactiver le mode sombre
    function toggleDarkMode() {
        Body.classList.toggle('dark-mode'); // Active/désactive la classe dark-mode sur le corps du document
    }
    
    // Écoutez les clics sur le bouton
    darkModeToggle.addEventListener('click', toggleDarkMode);
    
    /* ---------------------------------repeatdarkmode------------------------------------- */
const darkModeToggleit = document.getElementById('dark-mode-toggle');
const Bodyb = document.body;

// Fonction pour activer le mode sombre
function enableDarkMode() {
    Bodyb.classList.add('dark-mode');
    localStorage.setItem('darkModeEnabled', 'true');
}

// Fonction pour désactiver le mode sombre
function disableDarkMode() {
    Bodyb.classList.remove('dark-mode');
    localStorage.setItem('darkModeEnabled', 'false');
}

// Fonction pour basculer entre le mode sombre et le mode clair
function toggleDarkMode() {
    if (Bodyb.classList.contains('dark-mode')) {
        disableDarkMode();
    } else {
        enableDarkMode();
    }
}

// Écoutez les clics sur le bouton
darkModeToggleit.addEventListener('click', toggleDarkMode);

// Vérifiez le localStorage pour savoir si le mode sombre est activé ou non lors du chargement de la page
document.addEventListener('DOMContentLoaded', () => {
    const darkModeEnabled = localStorage.getItem('darkModeEnabled');
    if (darkModeEnabled === 'true') {
        enableDarkMode();
    } else {
        disableDarkMode();
    }
});

    // ------------------------------- select multiple --------------------------
    
    document.addEventListener("DOMContentLoaded", function () {
        const categoriesToggle = document.getElementById("categories-toggle");
        const categoryList = document.querySelector(".category-list");
    
        // Fonction pour afficher ou masquer la liste des catégories
        function toggleCategories() {
            categoryList.classList.toggle("hidden");
        }
    
        // Ajoute un gestionnaire d'événement au clic sur le li "Categories"
        categoriesToggle.addEventListener("click", toggleCategories);
    
        // Gestionnaire d'événement pour masquer la liste si l'utilisateur clique en dehors de la liste ou du li "Categories"
        document.addEventListener("click", function (event) {
            if (!categoryList.contains(event.target) && event.target !== categoriesToggle) {
                categoryList.classList.add("hidden");
            }
        });
    });

    /* ------------------------------------adaptationatempsreelsetextarea------------------- */

const textarea = document.getElementById("postContent");

// Ajoutez un gestionnaire d'événements pour l'événement d'entrée de texte
textarea.addEventListener("input", function () {
    // Réglez la hauteur du textarea en fonction de sa taille de contenu
    this.style.height = "auto";
    this.style.height = (this.scrollHeight) + "px";
});

//  ------------------ Auto resize textarea ---------------
document.addEventListener("DOMContentLoaded", function () {
    const postContentTextarea = document.getElementById("postContent");
    const errorMessage = document.getElementById("errorMessage");
    const maxLines = 15;

    postContentTextarea.addEventListener("input", function () {
        // Séparez le texte en lignes
        const lines = postContentTextarea.value.split("\n");

        // Vérifiez le nombre de lignes
        if (lines.length > maxLines) {
            // Si le nombre de lignes dépasse la limite, affichez le message d'erreur
            errorMessage.style.display = "block";
            // Raccourcissez le texte pour le ramener à 15 lignes maximum
            postContentTextarea.value = lines.slice(0, maxLines).join("\n");
        } else {
            // Si le nombre de lignes est dans la limite, masquez le message d'erreur
            errorMessage.style.display = "none";
        }
    });
});
