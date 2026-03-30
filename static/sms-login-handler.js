// SMS login button event handler
document.addEventListener('DOMContentLoaded', function() {
    const smsBtn = document.querySelector('.social-btn.social-sms');
    if (smsBtn) {
        smsBtn.addEventListener('click', async function() {
            const email = document.getElementById('email').value.trim();
            
            if (!email) {
                showStatusMessage('Please enter your email', 'error');
                return;
            }
            
            showLoading(true);
            
            try {
                const response = await fetch('/api/email/check-email-exists', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email: email })
                });
                
                const result = await response.json();
                
                if (result.code === 200) {
                    if (result.data.exists) {
                        showStatusMessage('Redirecting to SMS login...', 'success');
                        localStorage.setItem('esg_login_email', email);
                        setTimeout(() => {
                            window.location.href = '/static/sms-login.html';
                        }, 1000);
                    } else {
                        showStatusMessage('Email not registered, please register first', 'error');
                    }
                } else {
                    showStatusMessage(result.message || 'Check email failed', 'error');
                }
            } catch (error) {
                console.error('Error:', error);
                showStatusMessage('Network error, please try again', 'error');
            } finally {
                showLoading(false);
            }
        });
    }
});
