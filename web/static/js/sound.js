// sound.js

// Create an audio object and set the initial volume
var audio = new Audio('/static/sound/TOU DOUM.mp3');
audio.volume = localStorage.getItem('clickSoundVolume') || 1.0;

// Function to play the sound
function playClickSound() {
    if (localStorage.getItem('clickSoundEnabled') !== 'false') {
        audio.currentTime = 0;
        audio.play();
    }
}

// Add event listener to play sound on click
document.addEventListener("click", playClickSound);

// Function to toggle sound
function toggleSound(enable) {
    localStorage.setItem('clickSoundEnabled', enable);
}

// Function to set volume
function setVolume(volume) {
    audio.volume = volume;
    localStorage.setItem('clickSoundVolume', volume);
}

