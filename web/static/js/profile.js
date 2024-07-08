function deleteAccount() {
    const email = document.getElementById("email").innerText;
    if (confirm("Are you sure you want to delete your account? This action cannot be undone.")) {
        fetch('/deleteAccount', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email: email }),
        })
            .then(response => {
                if (response.ok) {
                    alert('Account deleted successfully.');
                    window.location.href = '/logout';
                } else {
                    alert('Failed to delete account.');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('An error occurred while deleting the account.');
            });
    }
}
// Sound controls
document.getElementById('soundToggle').addEventListener('change', function() {
    toggleSound(this.checked);
});

document.getElementById('soundVolume').addEventListener('input', function() {
    setVolume(this.value);
});

// Initialize sound controls based on local storage
document.getElementById('soundToggle').checked = localStorage.getItem('clickSoundEnabled') !== 'false';
document.getElementById('soundVolume').value = localStorage.getItem('clickSoundVolume') || 1.0;