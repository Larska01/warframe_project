
function loadPage(page) {
    fetch(page) 
        .then(response => response.text())
        .then(data => {
            let pageContent = data;
            const contentStart = data.indexOf('<body>');
            const contentEnd = data.indexOf('</body>');

            if (contentStart !== -1 && contentEnd !== -1) {
                pageContent = data.slice(contentStart + 6, contentEnd);
            }

            document.getElementById('content').innerHTML = pageContent;

            if (page === 'world.html') {
                loadWorldStateJS();
            }
        })
        .catch(error => {
            document.getElementById('content').innerHTML = '<h1>Page not found!</h1>';
            console.error('Error loading page:', error);
        });
}

function loadWorldStateJS() {
    // Check if worldstate.js is already loaded
    if (!document.querySelector('script[src="worldstate.js"]')) {
        var script = document.createElement('script');
        script.src = 'worldstate.js';  // Path to worldstate.js
        script.type = 'text/javascript';

        script.onload = function() {
            console.log('worldstate.js loaded successfully!');
        };

        script.onerror = function() {
            console.error('Error loading worldstate.js');
        };

        // Append the script to the body to ensure it executes
        document.body.appendChild(script);
    } else {
        console.log('worldstate.js is already loaded');
    }
}