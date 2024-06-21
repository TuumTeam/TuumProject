function openNav() {
    document.getElementById("mySidenav").style.width = "250px";
}

function closeNav() {
    document.getElementById("mySidenav").style.width = "0";
}

function showContent(section) {
    // Hide all content sections
    var contentSections = document.getElementsByClassName('content-section');
    for (var i = 0; i < contentSections.length; i++) {
        contentSections[i].style.display = 'none';
    }

    // Remove 'selected' class from all nav links
    var navLinks = document.getElementsByClassName('navLink');
    for (var i = 0; i < navLinks.length; i++) {
        navLinks[i].classList.remove('selected');
    }

    // Show the selected content section
    document.getElementById(section + 'Content').style.display = 'block';
    // Add 'selected' class to the clicked nav link
    document.getElementById(section + 'Link').classList.add('selected');
}

// Show the profile content by default when the page loads
document.addEventListener('DOMContentLoaded', function() {
    showContent('profile');
});
