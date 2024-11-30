// scripts.js
document.addEventListener('DOMContentLoaded', (event) => {
    const postReplyButton = document.getElementById('postReplyButton');
    const quickReplyPopup = document.getElementById('quickReplyPopup');
    const closeBtn = document.getElementsByClassName('close')[0];
    const popupHeader = document.getElementById('popupHeader');
    const popcont = document.getElementById('popcont');

    postReplyButton.addEventListener('click', () => {
        quickReplyPopup.style.display = 'block';
        quickReplyPopup.style.top = '100px';  // Initial position
        quickReplyPopup.style.left = '100px'; // Initial position
    });

    closeBtn.addEventListener('click', () => {
        quickReplyPopup.style.display = 'none';
    });

    let isDragging = false;
    let offsetX, offsetY;

    popcont.addEventListener('mousedown', (e) => {
        isDragging = true;
        offsetX = e.clientX - quickReplyPopup.offsetLeft;
        offsetY = e.clientY - quickReplyPopup.offsetTop;
    });

    document.addEventListener('mousemove', (e) => {
        if (isDragging) {
            quickReplyPopup.style.left = `${e.clientX - offsetX}px`;
            quickReplyPopup.style.top = `${e.clientY - offsetY}px`;
        }
    });

    document.addEventListener('mouseup', () => {
        isDragging = false;
    });
});

