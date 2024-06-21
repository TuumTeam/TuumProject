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

console.log("Hello from index.js");
