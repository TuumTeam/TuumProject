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
