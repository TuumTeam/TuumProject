document.addEventListener("DOMContentLoaded", () => {
  const tabs = document.querySelectorAll(".tab");
  const pages = document.querySelectorAll(".page");

  tabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      const target = tab.getAttribute("data-target");

      pages.forEach((page) => {
        page.classList.remove("active");
        if (page.id === target) {
          page.classList.add("active");
        }
      });
    });
  });

  // Set the first page as active by default
  if (pages.length > 0) {
    pages[0].classList.add("active");
  }
});
