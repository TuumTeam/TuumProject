// Function to adjust the menu based on the screen width
function adjustMenuForScreenWidth() {
  var screenWidth = window.innerWidth;
  var menu = document.querySelector(".leftBar");
  var containerGrid = document.querySelector(".containerGrid");

  if (screenWidth > 800) {
    containerGrid.style.gridTemplateColumns = menu.classList.contains("open")
      ? "250px 1fr"
      : "0 1fr";
  } else {
    // Adjust as needed for screens < 800px
    containerGrid.style.gridTemplateColumns = "0 1fr";
  }
}

/* -------------------  Left menu ------------------- */

// Toggle menu click event listener
document.getElementById("toggle-menu").addEventListener("click", function () {
  var menu = document.querySelector(".leftBar");
  menu.classList.toggle("open");
  // Adjust the menu on toggle
  adjustMenuForScreenWidth();

  var icon = document.getElementById("toggle-menu");
  console.log("fliped");
  icon.classList.add("flip-animation");

  // Remove the animation class after it completes to allow re-triggering
  icon.addEventListener("animationend", function () {
    icon.classList.remove("flip-animation");
  });
});

window.addEventListener("resize", adjustMenuForScreenWidth);

/* -------------------  Popup ------------------- */
const newTuum = document.getElementById("newTuum");
const popup = document.getElementById("popup");
const closePopup = document.getElementsByClassName("popup-close-btn")[0];
const windowWrapper = document.getElementById("windowWrapper");
const leaveTuumBtn = document.getElementById("leaveTuumBtn");

newTuum.addEventListener("click", () => {
  windowWrapper.classList.add("blur");
  popup.style.display = "block";
});

// Add event listener to the creation type dropdown
document.getElementById("creationType").addEventListener("change", function () {
  var selection = this.value;
  const hashtag = document.getElementById("hashtagArea");
  const research = document.getElementById("researchBar");
  const description = document.getElementById("description");
  const title = document.getElementById("title");
  if (selection === "newTuum") {
    // Fetch the search results from the server

    research.style.display = "flex";
    description.style.gridRow = 4;
    title.style.gridRow = 3;
  } else if (selection === "newRoom") {
    research.style.display = "none";
    title.style.gridRowStart = 2;
    title.style.gridRowEnd = 4;
    console.log("newSalon");
  }
});

document.getElementById("leaveTuumBtn").addEventListener("click", function () {
  document.getElementById("popup").style.display = "none";
  // Optionally, clear the form fields if needed
  document.querySelector(".popup-title").value = "";
  document.querySelector(".popup-description").value = "";
  windowWrapper.classList.remove("blur");
});

$(document).ready(function () {
  // Vérifier que l'élément researchBar existe
  let researchBar = document.getElementById("formSearchBar");
  if (!researchBar) {
    console.error("L'élément avec l'ID 'researchBar' n'existe pas.");
    return;
  }

  // Ajouter un gestionnaire d'événements pour l'entrée de texte
  $(researchBar).on("input", function () {
    // Récupérer la valeur du champ de recherche
    let query = $(this).val();
    // La requête est envoyée à chaque entrée, vérifier que la longueur de la requête est supérieure à 1
    if (query.length > 1) {
      // Effectuer une requête AJAX GET vers le serveur
      $.ajax({
        url: "/search", // URL de l'endpoint de recherche
        method: "GET", // Méthode HTTP
        data: { q: query }, // Paramètres de la requête
        dataType: "json", // S'assurer que la réponse est traitée comme du JSON
        success: function (data) {
          $("#searchBarProposition").empty();
          if (Array.isArray(data)) {
            if (data.length === 0) {
              $("#searchBarProposition").append(
                `<option value="Aucun résultat trouvé">Aucun résultat trouvé</option>`
              );
            } else {
              data.forEach((item) => {
                $("#searchBarProposition").append(
                  `<option value="${item}">${item}</option>`
                );
              });
            }
          } else {
            console.error("La réponse n'est pas un tableau", data);
          }
        },
        error: function (xhr, status, error) {
          console.error("Erreur lors de la requête AJAX :", status, error);
        },
      });
    } else {
      // Si la requête est trop courte, vider les résultats
      $("#searchBarProposition").empty();
    }
  });
});
document.addEventListener("DOMContentLoaded", () => {
  console.log("JavaScript chargé!");
});
