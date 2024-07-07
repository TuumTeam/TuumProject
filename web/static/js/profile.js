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