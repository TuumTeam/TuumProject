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
  leaveTuumBtn.style.display = "block";
  windowWrapper.classList.add("blur"); // Add blur to windowWrapper

  // Show the popup
  popup.style.display = "block";
});

document.getElementById("leaveTuumBtn").addEventListener("click", function () {
  document.getElementById("popup").style.display = "none";
  // Optionally, clear the form fields if needed
  document.querySelector(".popup-title").value = "";
  document.querySelector(".popup-description").value = "";
  windowWrapper.classList.remove("blur");
  leaveTuumBtn.style.display = "none";
});

// Add event listener to close the popup
// closePopup.addEventListener("click", () => {
//   windowWrapper.classList.remove("blur"); // Remove blur from windowWrapper
// });
